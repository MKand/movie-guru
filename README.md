# Table of Contents

- [Table of Contents](#table-of-contents)
  - [Movie Guru](#movie-guru)
  - [Description](#description)
  - [Overall Architecture](#overall-architecture)
    - [Components](#components)
  - [Deployment](#deployment)
    - [Docker Containers](#docker-containers)
  - [Flows](#flows)
    - [Data](#data)
    - [Postgres](#postgres)
  - [Environment setup](#environment-setup)
    - [Prerequisites](#prerequisites)
    - [Steps](#steps)
  - [Database Setup](#database-setup)
  - [Run the Application](#run-the-application)

## Movie Guru

**Genkit version**: 0.9.12 for Node.js

[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/YOUR_VIDEO_ID)

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
- **Database:** Stores movie data, embeddings, and user profiles in a local Postgres databse with `pgvector`.
- **Cache (Redis):** Caches conversation history and session data.

## Deployment

### Docker Containers

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

There are 2 tables:

- *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
- *user_preferences*: This contains the user's long term preferences profile information.

## Environment setup

### Prerequisites

- A Google Cloud project with owner permissions.
- Tools:
  - [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)
  - Docker and Docker Compose
- Required APIs enabled (will be performed in `setup_local.sh`).

### Steps

1. **Clone the Repository**

   ```sh
   git clone https://github.com/MKand/movie-guru.git
   cd movie-guru
   git checkout <current-branch> # Replace with branch name
   ```

1. Authenticate with Google Cloud

    ```sh
    gcloud auth login
    gcloud config set project <YOUR_PROJECT_ID>
    ```

1. Set the require environment variables

    ```sh
    export PROJECT_ID=<YOUR_PROJECT_ID>
    export LOCATION=<YOUR_DESIRED_GCLOUD_REGION> # defaults to us-central1 if this is not set
    ```

1. Run setup script.

    ```sh
    chmod +x setup_local.sh
    ./setup_local.sh
    ```

This enables the required APIs and creates the necessary service account with roles.

## Database Setup

1. Create a shared network for all the app containers we will use

    ```sh
    docker network create db-shared-network
    ```

1. Setup local DB
We'll setup a local *pgvecto*r db and an *Adminer* instance

    ```sh
    docker compose -f docker-compose-pgvector.yaml up -d
    ```

Navigate to *localhost:8082*, to access the db via *Adminer*. Use the main user credentials (user name: main, password: mainpassword).

At this stage, there will be 2 tables, with no data. We will populate the table in the next steps.

1. Grant the necessary permissions to minimal-user.

    ```SQL
    GRANT SELECT ON movies TO "minimal-user";
    GRANT SELECT, INSERT, UPDATE, DELETE ON user_preferences TO "minimal-user";
    ```

1. Populate the movie table

    ```sh
    source set_env_vars.sh
    export PROJECT_ID=<YOUR_PROJECT_ID>
    export LOCATION=<YOUR_DESIRED_GCLOUD_REGION> # defaults to us-central1 if this is not set
    ```

1. Download the JSON key for the service account

    This is required for you to be able to grant your docker containers access to the Vertex APIs

    **Note:** This step requires you to have the ability to create JSON keys for a service account. It may be disabled by some organizations.

    - Go to the project in the GCP console. Go to **IAM > Service Accounts**.
    - Select the movie guru service account (movie-guru-local-sa@<project id>.iam.gserviceaccount.com).
    - Create a new JSON key.
    - Download the key and store it as **.key.json** in the root of this repo (make sure you use the filename exactly).

1. Run the javascript indexer so it can add movies data into the database. The execution of this intentionally slowed down to stay below the rate-limits.

    ```sh
    docker compose -f docker-compose-indexer.yaml up --build -d 
    ```

    This takes about 10-15 minutes to run, so be patient. The embedding creation process is slowed down intentionally to ensure we stay under the rate limit.

1. Verify the number of entries in the DB.
There should be **652** entries in the movies table.

    ```sql
    SELECT COUNT(*)
    FROM "movies";
    ```

1. Shut down the indexer container.

    ```sh
    docker compose -f docker-compose-indexer.yaml down
    ```

Once all the required data is added, it is time to run the application that consists of the **frontend**, the **webserver**, the **genkit flows** server and the **redis cache**. These will be running locally in containers. The servers communicate with the **postgres DB** also running locally in a container.

## Run the Application

1. Make sure the env variables are in the execution context of docker compose.

    ```sh
    source set_env_vars.sh
    export PROJECT_ID=<YOUR_PROJECT_ID>
    export LOCATION=<YOUR_DESIRED_GCLOUD_REGION> # defaults to us-central1 if this is not set
    ```

1. Start the application.

    ```sh
    docker compose up --build
    ```

1. Access the Application Open http://localhost:5173 in your browser.

1. Once finished, stop the application.

    ```sh
    docker compose down
    ```
