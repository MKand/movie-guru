# Table of Contents

- [Table of Contents](#table-of-contents)
  - [Movie Guru](#movie-guru)
  - [Description](#description)
  - [Overall Architecture](#overall-architecture)
    - [Components](#components)
  - [Deployment](#deployment)
    - [Containers](#containers)
  - [Flows](#flows)
    - [Data](#data)
    - [Postgres](#postgres)
- [Movie Guru App Deployment Guide](#movie-guru-app-deployment-guide)
  - [Prerequisites](#prerequisites)
  - [Step 1: Set Environment Variables](#step-1-set-environment-variables)
  - [Step 2: Deploy Infrastructure](#step-2-deploy-infrastructure)
  - [Step 3: Get Firebase Web app Configuration](#step-3-get-firebase-web-app-configuration)
  - [Step 4: Update Environment Variables](#step-4-update-environment-variables)
  - [Step 5: Build and Push Containers](#step-5-build-and-push-containers)
  - [Step 6: Run the Cloud Run Job](#step-6-run-the-cloud-run-job)
  - [Step 7: Connect to GKE Cluster](#step-7-connect-to-gke-cluster)
  - [Step 8: Deploy Application Using Helm](#step-8-deploy-application-using-helm)
  - [Step 9: Register Workloads and Services with App Hub](#step-9-register-workloads-and-services-with-app-hub)
  - [Final Step: Verify Deployment](#final-step-verify-deployment)
  - [Appendix](#appendix)
    - [Original repo](#original-repo)
    - [Indexer](#indexer)
    - [Support](#support)
  - [License](#license)

## Movie Guru

**Genkit version**: 0.9.12 for Node.js

[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/l_KhN3RJ8qA)

## Description

Movie Guru is a website that helps users find movies to watch through an RAG powered chatbot. The movies are all fictional and are generated using GenAI.
The goal of this repo is to explore the best practices when building AI powered applications.

## Overall Architecture

### Components

- **Frontend (Vue.js):** User interface for interacting with the chatbot.
- **Web Backend (Go):** Handles API requests and communicates with the Flows Backend.
- **Flows Backend (Genkit for Node):** Orchestrates AI tasks, connects to GenAI models, and interacts with a vector database.
- **Database:** Stores movie data, embeddings, and user profiles in a Postgres databse with `pgvector`.
- **Cache (Redis):** Caches conversation history and session data.

## Deployment

### Containers

- **Frontend:** Vue.js application.
- **Web Backend:** Go-based API server.
- **Flows Backend:** Node.js-based AI task orchestrator.
- **Cache:** Cloud Memorystore for Redis for caching chat history and sessions.
- **Database:** Cloud SQL with `pgvector` extension.

## Flows

1. **User Profile Flow:** Extracts user preferences from conversations.
2. **Query Transform Flow:** Maps vague user queries to specific database queries.
3. **Movie Flow:** Combines user data and relevant documents to provide responses.
4. **Movie Doc Flow:** Retrieves relevant documents from the vector database. Perform a keyword based, vector based, or mixed search based on the type of query.
5. **Indexer Flow:** Parses movie data and adds it to the vector database.

### Data

- The data about the movies is stored in a pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
- The user's profile data (their likes and dislikes) are stored in the CloudSQL database.
- The user's conversation history is stored in a local redis cache. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.

### Postgres

There are multiple tables:

- *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
- *user_preferences*: This contains the user's long term preferences profile information.

# Movie Guru App Deployment Guide

## Prerequisites

Ensure you have the following installed and configured before proceeding:

- **Google Cloud SDK**: <https://cloud.google.com/sdk/docs/install>
- **Helm**: <https://helm.sh/docs/intro/install/>
- **Firebase Account**: Access to Firebase Console <https://console.firebase.google.com/>
- **GCP Account** with necessary permissions to create and manage GKE, Cloud Build, and IAM resources.
- [Optional for debugging] **kubectl**: <https://kubernetes.io/docs/tasks/tools/install-kubectl/>

## Step 1: Set Environment Variables

Before starting, set the following environment variables in your terminal:

```bash
export PROJECT_ID=<your-gcp-project-id>
export REGION=<your-desired-gcp-region>

```

## Step 2: Deploy Infrastructure

Run the following script to deploy the infrastructure using Cloud Build:

```bash
./deploy/infra.sh --region $REGION
```

This will trigger a pipeline that creates the required infrastructure on GCP. The process will take approximately **10-15 minutes** to complete.

**Note**: The setup of the Identity Platform and Firebase Auth through Terraform is a bit fragile. You might need to manually reconfigure it if the deployed application is creating auth issues.

## Step 3: Get Firebase Web app Configuration

Once the infrastructure setup is complete,

1. Perform the steps [here](./firebase-setup/README.md)
2. Go to the **Firebase Console**:

- A new project with the **Display Name: "Movie Guru App"** should be created.
- Navigate to the **Authentication** section and **enable Google Auth** for the web app.

Next, **copy the Firebase configuration parameters** (e.g., API key, auth domain, etc.) from the Firebase web app settings. You will need these values in the next step.

## Step 4: Update Environment Variables

Create a file **`set_env_vars.sh`** and **replace the placeholder values** with the Firebase parameters and the external IP address obtained in the previous steps.

```bash
export REGION="us-central1"
export PROJECT_ID="change"
export FIREBASE_API_KEY=""
export FIREBASE_AUTH_DOMAIN=""
export FIREBASE_GCP_ID=""
export FIREBASE_STORAGE_BUCKET=""
export FIREBASE_MESSAGING_SENDERID=""
export FIREBASE_APPID=""
export GATEWAY_ADDRESS="movie-guru.endpoints.${PROJECT_ID}$.cloud.goog"
export SERVER_URL="https://${GATEWAY_ADDRESS}"
```

After updating the file, run the script to apply the environment variables:

```bash
source set_env_vars.sh
```

## Step 5: Build and Push Containers

Run the following script to build and push the application containers using Cloud Build:

```bash
source ./deploy/ci.sh --region $REGION
```

This should take around 10 minutes

## Step 6: Run the Cloud Run Job

```bash
gcloud run job execute movie-guru-db-init --region $REGION --project $PROJECT_ID
```

## Step 7: Connect to GKE Cluster

Go to the **GKE page** in the **GCP Console** and find the connection string for your cluster.

```bash
gcloud container clusters get-credentials movie-guru-cluster --region ${REGION} --project ${PROJECT_ID}
```

## Step 8: Deploy Application Using Helm

Deploy the application to GKE using Helm:

Update the [helm](./deploy/app/helm/movieguru/values.yaml) file with values obtained in the previous step. Ensure the (helm)(./deploy/app.sh#L90) file has the right file

```bash
./deploy/app.sh --region $REGION
```

## Step 9: Register Workloads and Services with App Hub

```bash
./deploy/register.sh --region $REGION
```

## Final Step: Verify Deployment

Once the Helm deployment is complete, verify that the application is running correctly on your GKE cluster.
You can go to `https://movie-guru.endpoints.${PROJECT_ID}.cloud.goog` to interact with the app

___

## Appendix

### Original repo

This repo is a fork the original [repo](https://github.com/MKand/movie-guru)

### Indexer

The index loads data in the database through a cloud run job. The dataset CSV must be modified and embedded. To embed the file, run the following steps:

1. Update the CSV file as needed
2. Run the command `cd indexer && go-bindata  -o pkg/dataset/dataset.go ../dataset`
3. Fix the package name in  `pkg/dataset/dataset.go`
4. Re-run `./deploy/ci.sh --region $REGION`
5. Re-run `./deploy/deploy.sh --region $REGION` to deploy the cloud run job (if not already done)

**NOTE**: Get `go-bindata` by running the command `go install -a -v github.com/go-bindata/go-bindata/...@latest`

### Support

This demo is *NOT* endorsed by Google or Google Cloud.  
The repo is intended for educational/hobbyists use only.

## License

The AI generated movie data and posters in the repo are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit <http://creativecommons.org/licenses/by/4.0/>
