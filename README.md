# Table of Contents

- [Table of Contents](#table-of-contents)
  - [About Movie Guru](#about-movie-guru)
  - [Description](#description)
    - [Overall Architecture](#overall-architecture)
      - [Components](#components)
    - [Deployment](#deployment)
      - [Docker Containers](#docker-containers)
    - [Genkit Flows](#genkit-flows)
    - [Data](#data)
      - [Postgres](#postgres)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Environment setup](#environment-setup)
    - [Database Setup](#database-setup)
      - [Run the database service](#run-the-database-service)
      - [Populate the database (Optional)](#populate-the-database-optional)
    - [Run the Application](#run-the-application)
    - [Create a Firebase project to view AI monitoring (Optional)](#create-a-firebase-project-to-view-ai-monitoring-optional)
    - [Clean up](#clean-up)

## About Movie Guru

**Genkit version**: 0.9.12 for Node.js

[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/YOUR_VIDEO_ID)

 This version is a *minimal version* of the frontend and backend that doesn't have complex login logic like the version in **main**. It is meant to be run fully locally while using VertexAI APIs.

## Description

Movie Guru is a website that helps users find movies to watch through an RAG powered chatbot. The movies are all fictional and are generated using GenAI.
The goal of this repo is to explore the best practices when building AI powered applications.

This demo is *NOT* endorsed by Google or Google Cloud.  
The repo is intended for educational/hobbyists use only.

Refer to the readme in the **main** branch for more information.

### Overall Architecture

#### Components

- **Frontend (Vue.js):** User interface for interacting with the chatbot.
- **Web Backend (Go):** Handles API requests and communicates with the Flows Backend.
- **Flows Backend (Genkit for Node):** Orchestrates AI tasks, connects to GenAI models, and interacts with a vector database.
- **Database:** Stores movie data, embeddings, and user profiles in a local Postgres databse with `pgvector`.
- **Cache (Redis):** Caches conversation history and session data.

### Deployment

#### Docker Containers

- **Frontend:** Vue.js application.
- **Web Backend:** Go-based API server.
- **Flows Backend:** Node.js-based AI task orchestrator.
- **Cache:** Redis for caching chat history and sessions.
- **Database:** Postgres with `pgvector`.

### Genkit Flows

1. **User Profile Flow:** Extracts user preferences from conversations.
2. **Query Transform Flow:** Maps vague user queries to specific database queries.
3. **Movie Flow:** Combines user data and relevant documents to provide responses.
4. **Movie Doc Flow:** Retrieves relevant documents from the vector database. Perform a keyword based, vector based, or mixed search based on the type of query.
5. **Indexer Flow:** Parses movie data and adds it to the vector database.

### Data

- The data about the movies is stored in a pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
- The user's profile data (their likes and dislikes) are stored in the CloudSQL database.
- The user's conversation history is stored in a local redis cache. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.

#### Postgres

There are 2 tables:

- *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
- *user_preferences*: This contains the user's long term preferences profile information.

## Getting Started

### Prerequisites

- A Google Cloud project with owner permissions.
- Tools:
  - [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)
  - Docker and Docker Compose
- Required APIs enabled (will be performed in `setup_local.sh`).

### Environment setup

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

### Database Setup

#### Run the database service

1. Create a shared network for all the app containers we will use

    ```sh
    docker network create db-shared-network
    ```

2. Setup local DB
We'll setup a local *pgvector* db and an *Adminer* instance

    ```sh
    docker compose -f docker-compose-pgvector.yaml up -d
    ```

Navigate to *localhost:8082*, to access the db via *Adminer*. Use the main user credentials (user name: main, password: mainpassword).
Make sure you set `System` as `PostgresSQL` and `Server` as `db`, and `Database` as `fake-movies-db`.

#### Populate the database (Optional)

At this stage, there will be 2 tables, with data pre-loaded.
You can either choose to either reload the data into the table again or skip ahead to the [Run the application](#run-the-application) step.
Skipping ahead will save you approx. 20 minutes of the setup time.

1. Populate the movie table

    ```sh
    source set_env_vars.sh
    export PROJECT_ID=<YOUR_PROJECT_ID>
    export LOCATION=<YOUR_DESIRED_GCLOUD_REGION> # defaults to us-central1 if this is not set
    ```

2. Run the javascript indexer so it can add movies data into the database. The database comes pre-populated with the required data, but you can choose to re-add the data. The execution of this intentionally slowed down to stay below the rate-limits.

    ```sh
    docker compose -f docker-compose-indexer.yaml up --build -d 
    ```

    This takes about 10-15 minutes to run, so be patient. The embedding creation process is slowed down intentionally to ensure we stay under the rate limit.

3. Shut down the indexer container.

    ```sh
    docker compose -f docker-compose-indexer.yaml down
    ```

4. Verify the number of entries in the DB.
There should be **652** entries in the movies table.

    ```sql
    SELECT COUNT(*)
    FROM "movies";
    ```

Once all the required data is added, it is time to run the application that consists of the **frontend**, the **webserver**, the **genkit flows** server and the **redis cache**. These will be running locally in containers. The servers communicate with the **postgres DB** also running locally in a container.

### Run the Application

1. Make sure the env variables are in the execution context of docker compose.

    ```sh
    source set_env_vars.sh
    export PROJECT_ID=<YOUR_PROJECT_ID>
    export LOCATION=<YOUR_DESIRED_GCLOUD_REGION> # defaults to us-central1 if this is not set
    ```

1. Start the application services.

    ```sh
    docker compose up --build
    ```

1. Access the Application Open http://localhost:5173 in your browser.

### Create a Firebase project to view AI monitoring (Optional)

To view the traces of the application Genkit flows, you will need to add firebase to the Project.
Follow the steps [here](https://firebase.google.com/docs/projects/use-firebase-with-existing-cloud-project#how-to-add-firebase_console) and use the value of the PROJECT_ID of the Google Cloud project.
Once created, you can navigate to the Genkit section of the firebase console to view the traces of your AI flows in the application.

### Clean up

Run the following commands:

```sh
  docker compose down
  docker compose -f docker-compose-pgvector.yaml down
  docker network rm db-shared-network
```
