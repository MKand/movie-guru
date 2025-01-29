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

# Verify that the scripts are being run from Linux and not Mac
if [[ $OSTYPE != "linux-gnu" ]]; then
    echo -e "\e[91mERROR: This script has only been tested on Linux. Currently, only Linux (debian) is supported. Please run in Cloud Shell or in a VM running Linux".
    exit;
fi

# Export a SCRIPT_DIR var and make all links relative to SCRIPT_DIR
export SCRIPT_DIR=$(dirname $(readlink -f $0 2>/dev/null) 2>/dev/null || echo "${PWD}/$(dirname $0)")

usage()
{
   echo ""
   echo "Usage: $0"
   echo -e "\t--backend | -b Must be one of 'genkit-go', 'genkit-python' or 'langchain'. Default is 'genkit-go'."
   echo -e "\tExample usage: /deploy.sh -b genkit-go"
   exit 1 # Exit script after printing help
}

skip_infra()
{
    echo -e "\e[95mSetting Skip Infra var to 'true'...\e[0m"
    INFRA=false
}

skip_app()
{
    echo -e "\e[95mSetting Skip App var to 'true'...\e[0m"
    APP=false
}


# Setting default value
INFRA=true
APP=true
APP_VERSION=v1_go

# Define bash args
while [ "$1" != "" ]; do
    case $1 in
        --skipinfra  | -i )      shift
                                skip_infra
                                ;;
        --skipapp  | -a )      shift
                                skip_app
                                ;;
        --help | -h )           usage
                                exit
                                ;;

        
    esac
    shift
done


# Set project to PROJECT_ID or exit
[[ ! "${PROJECT_ID}" ]] && echo -e "Please export PROJECT_ID variable (\e[95mexport PROJECT_ID=<YOUR PROJECT ID>\e[0m)\nExiting." && exit 0
echo -e "\e[95mPROJECT_ID is set to ${PROJECT_ID}\e[0m"
gcloud config set core/project ${PROJECT_ID}

BUCKET_NAME="gs://${PROJECT_ID}"
REGION="europe-west4"

      # Check if the bucket exists.  gsutil ls exits with 0 if it exists, non-zero otherwise.
      gsutil ls "$BUCKET_NAME" > /dev/null 2>&1  # Redirect output to suppress it

      if [[ $? -ne 0 ]]; then  # Check the exit code of gsutil ls
        echo -e "\e[95mBucket $BUCKET_NAME for terraform state does not exist. Creating...\e[0m"
        gsutil mb -l "$REGION" "$BUCKET_NAME"
      else
        echo -e "\e[95mBucket $BUCKET_NAME for terraform state already exists.\e[0m"
      fi

      gsutil versioning set on "$BUCKET_NAME" # Enable versioning (always do this)
      echo -e "\e[95mVersioning enabled on $BUCKET_NAME\e[0m"
      

# Enable Cloudbuild API
echo -e "\e[95mEnabling Cloudbuild API in ${PROJECT_ID}\e[0m"
gcloud services enable servicenetworking.googleapis.com \
cloudbuild.googleapis.com \
serviceusage.googleapis.com \
cloudresourcemanager.googleapis.com \
aiplatform.googleapis.com \
artifactregistry.googleapis.com \
cloudresourcemanager.googleapis.com \
storage-api.googleapis.com \
run.googleapis.com \
firebase.googleapis.com  \

while true; do
    SHORT_SHA=$(head -c 64 /dev/urandom | tr -dc 'a-z0-9' | grep -E '^[a-z]' | head -n 1 | cut -c1-63)
    if [[ -n $SHORT_SHA ]]; then  # Check if SHORT_SHA is not empty
        break
    fi
done
echo -e "\e[95mSHORT_SHA is set to ${SHORT_SHA}\e[0m"

# Start main build
echo -e "\e[95mStarting Cloudbuild to CREATE infrastructure using terraform...\e[0m"

gcloud builds submit --config=deploy/setup-infra.yaml --async --ignore-file=.gcloudignore --substitutions=_PROJECT_ID=${PROJECT_ID},\
_SHORT_SHA=${SHORT_SHA},\
_APP_VERSION=${APP_VERSION}
