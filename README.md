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
  - [Step 4: Retrieve External IP Address](#step-4-retrieve-external-ip-address)
  - [Step 5: Update Environment Variables](#step-5-update-environment-variables)
  - [Step 6: Build and Push Containers](#step-6-build-and-push-containers)
  - [Step 7: Connect to GKE Cluster](#step-7-connect-to-gke-cluster)
  - [Step 8: Deploy Application Using Helm](#step-8-deploy-application-using-helm)
  - [Final Step: Verify Deployment](#final-step-verify-deployment)
  - [\[OPTIONAL\] Steps to run app locally for testing](#optional-steps-to-run-app-locally-for-testing)

## Movie Guru

**Genkit version**: 0.9.12 for Node.js

[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/l_KhN3RJ8qA)

 This version is a *minimal version* of the frontend and backend that doesn't have complex login logic like the version in **main**. It is meant to be run fully locally while using VertexAI APIs.

## Description

Movie Guru is a website that helps users find movies to watch through an RAG powered chatbot. The movies are all fictional and are generated using GenAI.
The goal of this repo is to explore the best practices when building AI powered applications.

This demo is *NOT* endorsed by Google or Google Cloud.  
The repo is intended for educational/hobbyists use only.

Refer to the readme in the **main** branch for more information.

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
- **Cache:** Redis for caching chat history and sessions.
- **Database:** Postgres with `pgvector`.

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

- **Google Cloud SDK**: https://cloud.google.com/sdk/docs/install
- **Helm**: https://helm.sh/docs/intro/install/
- **Firebase Account**: Access to Firebase Console https://console.firebase.google.com/
- **GCP Account** with necessary permissions to create and manage GKE, Cloud Build, and IAM resources.
- [Optional for debugging] **kubectl**: https://kubernetes.io/docs/tasks/tools/install-kubectl/

## Step 1: Set Environment Variables

Before starting, set the following environment variables in your terminal:

```bash
export PROJECT_ID=<your-gcp-project-id>
export REGION=<your-desired-gcp-region>

```

## Step 2: Deploy Infrastructure

Run the following script to deploy the infrastructure using Cloud Build:

```bash
./deploy/deploy.sh --region $REGION
```

This will trigger a pipeline that creates the required infrastructure on GCP. The process will take approximately **10-15 minutes** to complete.

**Note**: The setup of the Identity Platform and Firebase Auth through Terraform is a bit fragile. You might need to manually reconfigure it if the deployed application is creating auth issues.

## Step 3: Get Firebase Web app Configuration

Once the infrastructure setup is complete, go to the **Firebase Console**:

- A new project with the **Display Name: "Movie Guru App"** should be created.

Next, **copy the Firebase configuration parameters** (e.g., API key, auth domain, etc.) from the Firebase web app settings. You will need these values in the next step.

## Step 4: Retrieve External IP Address

Go to the **GCP Console** and search for **"IP addresses"**:

- Look for an IP address labeled **"movie-guru-external-ip"**.
- Copy this IP address for use in the environment variables setup.

## Step 5: Update Environment Variables

Open the file **`set_env_vars.sh`** and **replace the placeholder values** with the Firebase parameters and the external IP address obtained in the previous steps.

After updating the file, run the script to apply the environment variables:

```bash
./set_env_vars.sh
```

## Step 6: Build and Push Containers

Run the following script to build and push the application containers using Cloud Build:

```bash
./deploy/ci.sh --region $REGION
```

This should take around 10 minutes

## Step 7: Connect to GKE Cluster

Go to the **GKE page** in the **GCP Console** and find the connection string for your cluster.

- Copy the connection string.
- Run the command in your terminal to connect to the cluster.

## Step 8: Deploy Application Using Helm

Deploy the application to GKE using Helm:

```bash
helm upgrade --install movieguru \
./deploy/app/helm/movieguru \
--namespace movieguru \
--create-namespace \
--set PROJECT_ID=${PROJECT_ID} \
--set IMAGE.TAG=latest \
--set REGION=${REGION}
```

## Final Step: Verify Deployment

Once the Helm deployment is complete, verify that the application is running correctly on your GKE cluster.
You can go to http://$GATEWAY_IP to interact with the app
Make sure you use the invite code **0000** the first time you log into the app with your Google credentials. If you plan to leave the app running on the cloud for longer, make sure you change this value in the **invite_codes** table in the database. You can port forward Adminer to localhost and use the credentials **username: main, password: main** to login the database through Adminer.

## [OPTIONAL] Steps to run app locally for testing

- Make sure you've run **/deploy/deploy.sh** (to enable apis and create service accounts) and **/deploy/ci.sh** (to upload posters).
- Make sure you have substituted the environment variables in **set_env_vars.sh**

- Create the **pgvector/init_substituted.sql** file.

  ```bash
  ./set_env_vars.sh
  envsubst < pgvector/init.sql > pgvector/init_substituted.sql
  ```

- Download the service account key for movie-guru-chat-server-sa and store it as **.key.json** at the project root.

- Run the docker compose file

  ```bash
  docker compose up --build 
  ```

**Note**: If the server exits prematurely, it is because it is waiting for the db table to be populated with the necessary data. There is a race condition. Kill the container and re-run **docker compose up**

- The frontend is running at http://localhost:4001

- After teardown, run this to delete the temporary sql file

  ```bash
  rm pgvector/init_substituted.sql
  ```
