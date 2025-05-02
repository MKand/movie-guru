#!/usr/bin/env bash

# Set environment variables
source set_env_vars.sh

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
echo -e "\e[93mUsing Project: $PROJECT_ID\e[0m"


SERVICE_ACCOUNT_NAME="movie-guru-chat-server-sa"
SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"

echo -e "\e[95mCreating service account local key as key.json\e[0m"
gcloud iam service-accounts keys create key.json \
    --iam-account=$SERVICE_ACCOUNT_EMAIL

echo -e "\e[95m Substituting env variables in init.sql\e[0m"

envsubst < pgvector/init.sql > pgvector/init_substituted.sql
