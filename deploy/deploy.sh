#!/usr/bin/env bash

# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
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

# Default region
DEFAULT_REGION="europe-west4"
REGION="$DEFAULT_REGION"

# Usage function
usage() {
   echo ""
   echo "Usage: $0 [--region <region>]"
   echo -e "\t--region, -r : Specify a region (default: europe-west4)"
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

echo -e "\e[95mUsing region: $REGION\e[0m"

# Check if PROJECT_ID is set, or exit
[[ ! "${PROJECT_ID}" ]] && echo -e "Please export PROJECT_ID variable (\e[95mexport PROJECT_ID=<YOUR PROJECT ID>\e[0m)\nExiting." && exit 0
echo -e "\e[95mPROJECT_ID is set to ${PROJECT_ID}\e[0m"
gcloud config set core/project "${PROJECT_ID}"

BUCKET_NAME="gs://${PROJECT_ID}"

# Check if the bucket exists
gsutil ls "$BUCKET_NAME" > /dev/null 2>&1  # Suppress output

if [[ $? -ne 0 ]]; then  # Check exit code of gsutil ls
    echo -e "\e[95mBucket $BUCKET_NAME for Terraform state does not exist. Creating...\e[0m"
    gsutil mb -l "$REGION" "$BUCKET_NAME"
else
    echo -e "\e[95mBucket $BUCKET_NAME for Terraform state already exists.\e[0m"
fi

# Enable versioning
gsutil versioning set on "$BUCKET_NAME"
echo -e "\e[95mVersioning enabled on $BUCKET_NAME\e[0m"

# Enable Cloud APIs
echo -e "\e[95mEnabling required Cloud APIs in ${PROJECT_ID}\e[0m"
gcloud services enable servicenetworking.googleapis.com \
    cloudbuild.googleapis.com \
    serviceusage.googleapis.com \
    cloudresourcemanager.googleapis.com \
    


# Make cloudbuild SA roles/owner for PROJECT_ID
echo -e "\e[95mAssigning Cloudbuild Service Account roles/owner in ${PROJECT_ID}\e[0m"
export PROJECT_NUMBER=$(gcloud projects describe ${PROJECT_ID} --format 'value(projectNumber)')

gcloud projects add-iam-policy-binding ${PROJECT_ID} --member serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --role roles/owner


# Start Cloud Build
echo -e "\e[95mStarting Cloud Build to CREATE infrastructure using Terraform...\e[0m"

gcloud builds submit --config=deploy/setup-infra.yaml --async --ignore-file=.gcloudignore --substitutions=_PROJECT_ID="${PROJECT_ID}",\
_REGION="${REGION}"
