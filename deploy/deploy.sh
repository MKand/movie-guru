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
BACKEND=genkit-go
APP_VERSION=v1_go

# Define bash args
while [ "$1" != "" ]; do
    case $1 in
        --backend | -b )      shift
                                BACKEND=$1
                                ;;
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
echo -e "\e[95mBACKEND is set to ${BACKEND}\e[0m"

# Enable Cloudbuild API
echo -e "\e[95mEnabling Cloudbuild API in ${PROJECT_ID}\e[0m"
gcloud services enable servicenetworking.googleapis.com \
cloudbuild.googleapis.com \
storage.googleapis.com \
serviceusage.googleapis.com \
cloudresourcemanager.googleapis.com \
aiplatform.googleapis.com \
artifactregistry.googleapis.com \
cloudresourcemanager.googleapis.com
sqladmin.googleapis.com \
storage-api.googleapis.com \
sql-component.googleapis.com \
run.googleapis.com \
redis.googleapis.com \
firebase.googleapis.com   

# Make cloudbiuld SA roles/owner for PROJECT_ID
# TODO: Make these permissions more granular to precisely what is required by cloudbuild
echo -e "\e[95mAssigning Cloudbuild Service Account roles/owner in ${PROJECT_ID}\e[0m"
export PROJECT_NUMBER=$(gcloud projects describe ${PROJECT_ID} --format 'value(projectNumber)')
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com --role roles/owner

while true; do
    SHORT_SHA=$(head -c 64 /dev/urandom | tr -dc 'a-z0-9' | grep -E '^[a-z]' | head -n 1 | cut -c1-63)
    if [[ -n $SHORT_SHA ]]; then  # Check if SHORT_SHA is not empty
        break
    fi
done
echo -e "\e[95mSHORT_SHA is set to ${SHORT_SHA}\e[0m"

# Start main build
echo -e "\e[95mStarting Cloudbuild to CREATE infrastructure using terraform...\e[0m"

[[ ${BACKEND} == "genkit-go" ]] && gcloud builds submit --config=deploy/setup-go-backend.yaml --async --ignore-file=.gcloudignore --substitutions=_PROJECT_ID=${PROJECT_ID},\
_SHORT_SHA=${SHORT_SHA},\
_APP_VERSION=v1_go,\
_INFRA=${INFRA},\
_APP=${APP}

[[ ${BACKEND} == "langchain" ]] && gcloud builds submit --config=deploy/setup-langchain-backend.yaml --async --ignore-file=.gcloudignore --substitutions=_PROJECT_ID=${PROJECT_ID},\
_SHORT_SHA=${SHORT_SHA},\
_LANGSMITH_API_KEY=${LANGSMITH_API_KEY},\
_LANGCHAIN_TRACING_V2=${LANGCHAIN_TRACING_V2},\
_LANGCHAIN_PROJECT=${LANGCHAIN_PROJECT},\
_APP_VERSION=v1,\
_INFRA=${INFRA},\
_APP=${APP}