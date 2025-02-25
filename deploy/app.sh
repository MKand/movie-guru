#!/usr/bin/env bash

# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Verify that the script is running on Linux (not macOS)
if [[ $OSTYPE != "linux-gnu" ]]; then
    echo -e "\e[91mERROR: This script has only been tested on Linux. Currently, only Linux (Debian) is supported. Please run in Cloud Shell or in a VM running Linux.\e[0m"
    exit 1
fi

# Export a SCRIPT_DIR var and make all links relative to SCRIPT_DIR
export SCRIPT_DIR=$(dirname "$(readlink -f "$0" 2>/dev/null)" 2>/dev/null || echo "${PWD}/$(dirname "$0")")

# Usage function
usage() {
   echo ""
   echo "Usage: $0 [--region <region>]"
   echo -e "\tExample: ./deploy.sh --region us-central1"
   exit 1
}

# Parse command-line arguments
while [ "$1" != "" ]; do
    case $1 in
        --region | -r ) shift
                        REGION=$1
                        ;;
        --help | -h )   usage
                        ;;
        * )             echo -e "\e[91mUnknown parameter: $1\e[0m"
                        usage
                        ;;
    esac
    shift
done


# Check if PROJECT_ID is set
if [[ -z "$PROJECT_ID" ]]; then
    echo -e "\e[91mERROR: PROJECT_ID environment variable is required.\e[0m"
    echo -e "Please set it using: \e[95mexport PROJECT_ID=<your-gcp-project-id>\e[0m"
    exit 1
fi

# Check if REGION is set
if [[ -z "$REGION" ]]; then
    echo -e "\e[91mERROR: REGION environment variable is required.\e[0m"
    echo -e "Please set it using: --region "
    exit 1
fi

echo -e "\e[95mUsing PROJECT_ID: $PROJECT_ID\e[0m"
echo -e "\e[95mUsing REGION: $REGION\e[0m"

if [[ -z "$SHORT_SHA" ]]; then
  SHORT_SHA="latest"
  echo -e "\e[95m Missing env var SHORT_SHA, using tag latest \e[0m"
else
  echo -e "\e[95mFound SHORT_SHA: $SHORT_SHA\e[0m"
  echo -e "\e[95m Using Image tag: $SHORT_SHA\e[0m"
fi

# Set GCP project
gcloud config set project "$PROJECT_ID"

echo -e "\e[95m Substituting env variables in init.sql\e[0m"

envsubst < pgvector/init.sql > pgvector/init_substituted.sql

gsutil cp pgvector/init_substituted.sql "gs://fb-webapp-${PROJECT_ID}/sql/init.sql"
rm pgvector/init_substituted.sql 


echo -e "\e[95mConnecting to GKE cluster\e[0m"
gcloud container clusters get-credentials movie-guru-cluster --region $REGION --project $PROJECT_ID
echo -e "\e[95m Starting Helm deploy for app...\e[0m"

helm upgrade --install movieguru \
./deploy/app/helm/movieguru \
--namespace movieguru \
--create-namespace \
-f ./deploy/app/helm/movieguru/values.simple.yaml \
--set projectID=${PROJECT_ID} \
--set Image.tag=$SHORT_SHA \
--set region=${REGION}


echo -e "\e[95m Creating ns and configmap for locust.\e[0m"

kubectl create ns locust

kubectl delete configmap loadtest-locustfile -n locust
kubectl create configmap loadtest-locustfile --from-file=locust/locustfile.py -n locust

echo -e "\e[95m Starting Helm deploy for locust.\e[0m"

# Need to delete to make sure the pods mount the updated config from cm
helm delete locust -n locust
helm upgrade --install locust \
  deliveryhero/locust \
  --namespace locust \
  --set loadtest.name=movieguru-loadtest \
  --set loadtest.locust_locustfile_configmap=loadtest-locustfile \
  --set loadtest.locust_locustfile=locustfile.py \
  --set loadtest.locust_host=http://server.movieguru.svc.cluster.local:8080 \
  --set service.type=ClusterIP \
  --set worker.replicas=5


echo -e "\e[95m Starting Helm deploy for mock user ...\e[0m"

helm upgrade --install mockuser \
./deploy/app/helm/mockuser \
--namespace mockuser \
--create-namespace \
--set projectID=${PROJECT_ID} \
--set Image.tag=${SHORT_SHA} \
--set region=${REGION} \
--set modelLocation=us-central1


echo -e "\e[95m Starting Helm deploy for otel collector ...\e[0m"

helm delete otel -n otel-collector
helm upgrade --install otel \
./deploy/app/helm/otel \
--namespace otel-collector \
--create-namespace \
--set projectID=${PROJECT_ID} 

echo -e "\e[95m Creating IAM bindings for the KSAs.\e[0m"

gcloud iam service-accounts add-iam-policy-binding movie-guru-chat-server-sa@${PROJECT_ID}.iam.gserviceaccount.com \
 --role roles/iam.workloadIdentityUser     --member "serviceAccount:${PROJECT_ID}.svc.id.goog[movieguru/movieguru-sa]"

gcloud iam service-accounts add-iam-policy-binding movie-guru-chat-server-sa@${PROJECT_ID}.iam.gserviceaccount.com \
 --role roles/iam.workloadIdentityUser     --member "serviceAccount:${PROJECT_ID}.svc.id.goog[mockuser/mockuser-sa]"

gcloud iam service-accounts add-iam-policy-binding movie-guru-chat-server-sa@${PROJECT_ID}.iam.gserviceaccount.com \
 --role roles/iam.workloadIdentityUser     --member "serviceAccount:${PROJECT_ID}.svc.id.goog[otel-collector/otel-sa]"

echo -e "\e[95m Port forwarding Locust to localhost:8089.\e[0m"

kubectl --namespace locust port-forward service/locust 8089:8089