#!/usr/bin/env bash

# Verify that the script is being run on Linux
if [[ $OSTYPE != "linux-gnu" ]]; then
    echo -e "\e[91mERROR: This script is supported only on Linux. Please run it in a Linux environment.\e[0m"
    exit 1
fi

# Check if PROJECT_ID is set
if [[ -z "$PROJECT_ID" ]]; then
    echo -e "\e[91mERROR: Please set the PROJECT_ID environment variable (e.g., export PROJECT_ID=<YOUR_PROJECT_ID>).\e[0m"
    exit 1
fi

echo -e "\e[95mEnabling required APIs for project: $PROJECT_ID\e[0m"

gcloud config set core/project "$PROJECT_ID"
gcloud services enable \
    storage.googleapis.com \
    serviceusage.googleapis.com \
    cloudresourcemanager.googleapis.com \
    aiplatform.googleapis.com \
    storage-api.googleapis.com \
    firebase.googleapis.com \
    monitoring.googleapis.com

echo -e "\e[95mAPIs have been enabled successfully.\e[0m"

# Create the service account
SERVICE_ACCOUNT_NAME="movie-guru-local-sa"
SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"

echo -e "\e[95mCreating service account: $SERVICE_ACCOUNT_NAME\e[0m"
gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME" \
    --description="Service account for Movie Guru application" \
    --display-name="Movie Guru Local SA"

# Assign roles to the service account
echo -e "\e[95mAssigning roles to service account: $SERVICE_ACCOUNT_EMAIL\e[0m"
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/aiplatform.user"
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/logging.logWriter"
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/monitoring.metricWriter"
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/cloudtrace.agent"

echo -e "\e[95mCreating service account local key as .key.json\e[0m"
gcloud iam service-accounts keys create ./.key.json \
    --iam-account=movie-guru-local-sa@$PROJECT_ID.iam.gserviceaccount.com

echo -e "\e[95mService account $SERVICE_ACCOUNT_NAME has been created and configured successfully.\e[0m"
