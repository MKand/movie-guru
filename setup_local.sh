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
echo -e "\e[93mUsing Project: $PROJECT_ID\e[0m"

# Check if REGION is set
if [[ -z "$REGION" ]]; then
    echo -e "\e[91mERROR: Please set the REGION environment variable (e.g., export REGION=<YOUR REGION>).\e[0m"
    exit 1
fi
echo -e "\e[93mUsing REGION: $REGION\e[0m"


SKIP_INFRA=false
for arg in "$@"; do
  if [[ "$arg" == "--skip-infra" ]]; then
    SKIP_INFRA=true
    echo -e "\e[93mSkipping infrastructure setup as requested.\e[0m"
    break
  fi
done

SERVICE_ACCOUNT_NAME="movie-guru-chat-server-sa"
SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"
POSTER_BUCKET_NAME="gs://${PROJECT_ID}_posters"


if [[ "$SKIP_INFRA" == false ]]; then

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
    echo -e "\e[95mService account $SERVICE_ACCOUNT_NAME has been created and configured successfully.\e[0m"
    
    # Check if the bucket exists
    gsutil ls "$POSTER_BUCKET_NAME" > /dev/null 2>&1  # Suppress output

    if [[ $? -ne 0 ]]; then  # Check exit code of gsutil ls
        echo -e "\e[95mBucket $POSTER_BUCKET_NAME for posters does not exist. Creating...\e[0m"
        gsutil mb -l "$REGION" "$POSTER_BUCKET_NAME"
    else
        echo -e "\e[95mBucket $POSTER_BUCKET_NAME for posters already exists.\e[0m"
    fi
fi


echo -e "\e[95mCreating service account local key as .key.json\e[0m"
gcloud iam service-accounts keys create ./.key.json \
    --iam-account=$SERVICE_ACCOUNT_EMAIL

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

echo -e "\e[95m Substituting env variables in init.sql\e[0m"

envsubst < pgvector/init.sql > pgvector/init_substituted.sql