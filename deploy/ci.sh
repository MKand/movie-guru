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

# Default region
REGION="europe-west4"

# Check if PROJECT_ID is set
if [[ -z "$PROJECT_ID" ]]; then
    echo -e "\e[91mERROR: PROJECT_ID environment variable is required.\e[0m"
    echo -e "Please set it using: \e[95mexport PROJECT_ID=<your-gcp-project-id>\e[0m"
    exit 1
fi

echo -e "\e[95mUsing PROJECT_ID: $PROJECT_ID\e[0m"
echo -e "\e[95mUsing REGION: $REGION\e[0m"

# Set GCP project
gcloud config set project "$PROJECT_ID"

# Generate a short SHA identifier
SHORT_SHA=$(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | fold -w 10 | head -n 1)
echo -e "\e[95mGenerated SHORT_SHA: $SHORT_SHA\e[0m"


# Start Cloud Build
echo -e "\e[95mStarting Cloud Build...\e[0m"
gcloud builds submit --config=deploy/ci.yaml --async --ignore-file=.gcloudignore \
  --substitutions=_PROJECT_ID="${PROJECT_ID}",_SHORT_SHA="${SHORT_SHA}",_REGION="${REGION}"

echo -e "\e[92mCloud Build submitted successfully!\e[0m"
