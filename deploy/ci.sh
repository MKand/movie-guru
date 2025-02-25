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
   echo "Usage: $0 [--region <region>] [--skip-gcs=true]"
   echo -e "\tExample: ./deploy.sh --region us-central1"
   exit 1
}

# Parse command-line arguments
while [ "$1" != "" ]; do
    case $1 in
        --region | -r ) shift
                        REGION=$1
                        ;;
        --skip-gcs | -s ) shift
                        SKIP_GCS="true"
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
    echo -e "Please set it using: \e[95mexport REGION=<your-gcp-region>\e[0m"
    exit 1
fi

echo -e "\e[95mUsing PROJECT_ID: $PROJECT_ID\e[0m"
echo -e "\e[95mUsing REGION: $REGION\e[0m"

# Set GCP project
gcloud config set project "$PROJECT_ID"

# Generate a short SHA identifier
SHORT_SHA=$(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | fold -w 10 | head -n 1)
echo -e "\e[95mGenerated SHORT_SHA: $SHORT_SHA\e[0m"

echo -e "\e[95mSubstituting env variables in init.sql\e[0m"

envsubst < pgvector/init.sql > pgvector/init_substituted.sql
envsubst < pgvector/py_init.sql > pgvector/py_init_substituted.sql

# Start Cloud Build
echo -e "\e[95mStarting Cloud Build...\e[0m"
gcloud builds submit --config=deploy/ci.yaml --region=${REGION} --async --ignore-file=.gcloudignore --worker-pool="projects/${PROJECT_ID}/locations/${REGION}/workerPools/movie-guru" --project=${PROJECT_ID} \
  --substitutions=_PROJECT_ID=$PROJECT_ID,_SHORT_SHA=$SHORT_SHA,_REGION=$REGION,_VITE_FIREBASE_API_KEY=$FIREBASE_API_KEY,_VITE_FIREBASE_AUTH_DOMAIN=$FIREBASE_AUTH_DOMAIN,_VITE_GCP_PROJECT_ID=$PROJECT_ID,_VITE_FIREBASE_STORAGE_BUCKET=$FIREBASE_STORAGE_BUCKET,_VITE_FIREBASE_MESSAGING_SENDERID=$FIREBASE_MESSAGING_SENDERID,_VITE_FIREBASE_APPID=$FIREBASE_APPID,_VITE_CHAT_SERVER_URL="https://movie-guru.endpoints.${PROJECT_ID}.cloud.goog/server"

# Check if SKIP_GCS is set
if [[ -n "$SKIP_GCS" ]]; then
    exit 0
fi

echo -e "\e[92mCloud Build submitted successfully!\e[0m"

echo -e "\e[92mDownloading and unzipping posters from the external archive..\e[0m"

# Download the zip file with posters
curl -o dataset/posters_small.zip https://storage.googleapis.com/movie-guru-posters/posters_small.zip

# Unzip into dataset/posters
unzip dataset/posters_small.zip -d .

# Upload posters to bucket
echo -e "\e[92mUploading posters to bucket and deleting local posters\e[0m"

# Delete zip file
rm dataset/posters_small.zip

gcloud storage cp ./dataset/posters_small/* "gs://${PROJECT_ID}_posters/"

rm -rf dataset/posters_small

echo -e "\e[ Making posters publicly readable\e[0m"

gcloud storage buckets add-iam-policy-binding "gs://${PROJECT_ID}_posters/" \
  --member="allUsers" \
  --role="roles/storage.objectViewer"

