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

# Check if REGION is set
if [[ -z "$REGION" ]]; then
    echo -e "\e[91mERROR: Please set the REGION environment variable (e.g., export REGION=<YOUR REGION>).\e[0m"
    exit 1
fi
echo -e "\e[93mUsing REGION: $REGION\e[0m"

SERVICE_ACCOUNT_NAME="movie-guru-chat-server-sa"
SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"
POSTER_BUCKET_NAME="gs://${PROJECT_ID}_posters"

# Check if the bucket exists
gsutil ls "$POSTER_BUCKET_NAME" > /dev/null 2>&1  # Suppress output

if [[ $? -ne 0 ]]; then  # Check exit code of gsutil ls
    echo -e "\e[95mBucket $POSTER_BUCKET_NAME for posters does not exist. Creating...\e[0m"
    gsutil mb -l "$REGION" "$POSTER_BUCKET_NAME"
else
    echo -e "\e[95mBucket $POSTER_BUCKET_NAME for posters already exists.\e[0m"
fi
fi

echo -e "\e[92mDownloading and unzipping posters from the external archive..\e[0m"

# Download the zip file with posters
curl -o dataset/posters_small.zip https://storage.googleapis.com/movie-guru-posters/posters_small.zip

# Unzip into dataset/posters
unzip dataset/posters_small.zip -d . > /dev/null 2>&1  # Suppress output

# Upload posters to bucket
echo -e "\e[92mUploading posters to bucket and deleting local posters\e[0m"

# Delete zip file
rm dataset/posters_small.zip 

gcloud storage cp ./dataset/posters_small/* $POSTER_BUCKET_NAME > /dev/null 2>&1  # Suppress output

rm -rf dataset/posters_small

echo -e "\e[93mMaking posters publicly readable\e[0m"

gcloud storage buckets add-iam-policy-binding $POSTER_BUCKET_NAME \
  --member="allUsers" \
  --role="roles/storage.objectViewer"
