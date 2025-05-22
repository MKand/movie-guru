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
    - [Simple Local Setup](#simple-local-setup)
      - [Prerequisites](#prerequisites)
      - [Instructions](#instructions)
    - [Simple local setup with Firebase Authentication](#simple-local-setup-with-firebase-authentication)
    - [Cloud setup with Firebase Authentication](#cloud-setup-with-firebase-authentication)
      - [Populate the database (Optional)](#populate-the-database-optional)

## About Movie Guru

**Genkit version**: 1.7.0 for Node.js

**Gemini Models**: Gemini 2.0 Flash amd Flash Lite

**Embedding Models**: textEmbedding005

[![Watch the Movie Guru Demo](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/l_KhN3RJ8qA "Movie Guru Demo")

> This version is a *minimal version* of the frontend and backend that doesn't have complex login logic like the version in the **cloud-movieguru** branch. It is meant to be run fully locally while using Vertex AI APIs.
> If you want to run this demo entirely in the cloud, please use the **cloud-movieguru** branch.

## Description

Movie Guru is a website that helps users find movies to watch through an RAG powered chatbot. The movies are all fictional and are generated using GenAI.
The goal of this repo is to explore the best practices when building AI powered applications.

The repo is intended for educational/hobbyists use only.

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

1. **Safety Flow:** Checks each user statement to validate whether or not it is safe to proceed with.
2. **Query Transform Flow:** Maps vague user queries to specific database queries.
3. **Movie QA Flow:** Combines user data and relevant documents to provide responses.
4. **Movie Doc Flow:** Retrieves relevant documents from the vector database. Performs a keyword-based, vector-based, or mixed search based on the type of query.
5. **Chat Flow:** Combines all the aforementioned prompts and flows into a single flow that is used by the chat server.
6. **User Preferences Flow:** Additional flow that extracts user preferences from conversations.
7. **Indexer Flow:** Parses movie data and adds it to the vector database.

### Data

- Movie data is stored in a PostgreSQL database with the `pgvector` extension. This includes details for approximately 600 fictional movies (plot, actors, director, rating, genre, poster link). Posters are stored in a Cloud Storage bucket.
- User profile data (likes and dislikes) is stored in the PostgreSQL database.
- User conversation history (most recent 10 messages, configurable) and webserver session info are stored in a local Redis cache.

#### Postgres

There are 2 important tables:

- *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
- *user_preferences*: This contains the user's long term preferences profile information.
- **movies**: Contains information about the AI-generated movies and their embeddings. Initial data is sourced from `dataset/movies_with_posters.csv`. If hosting your own posters, update the links in this file.
- **user_preferences**: Stores users' long-term preference profiles.

## Getting Started

### Simple Local Setup

#### Prerequisites

- A Google Cloud project with owner permissions.
- Tools:
  - [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)
  - Docker and Docker Compose
- Required APIs enabled (will be performed in `setup_local.sh`).

#### Instructions

1. **Clone the Repository**

   ```sh
   git clone https://github.com/MKand/movie-guru.git
   cd movie-guru
   git checkout main
   ```

1. Authenticate with Google Cloud

    ```sh
    gcloud auth login
    gcloud config set project <YOUR_PROJECT_ID>
    ```

1. Run the interactive environment setup script where you will be prompted to enter values like **PROJECT_ID** and **REGION** (which region do you want to consume the models from)

    ```sh
    chmod +x ./startup/configure_env_simple.sh
    ./startup/configure_env_simple.sh
    ```

1. Make sure the required APIs are enabled and create a service account. You will need owner level access to the project. This enables the required APIs and creates the necessary service account with roles.

    ```sh
    chmod +x ./startup/setup_cloud_simple.sh
    ./startup/setup_cloud_simple.sh
    ```

1. Start the app.

    ```sh
      chmod +x ./startup/launch_app.sh
      ./startup/launch_app.sh
    ```

1. Access the Frontend Application Open <http://localhost:8080> in your browser.

1. To stop the app, press **Ctrl+C** in the terminal. Then run
  
    ```sh
      ./startup/launch_app.sh --stop
    ```

### Simple local setup with Firebase Authentication

This uses firebase authentication for the frontend of the application.

1. **Clone the Repository**

   ```sh
   git clone https://github.com/MKand/movie-guru.git
   cd movie-guru
   git checkout main
   ```

2. Authenticate with Google Cloud

    ```sh
    gcloud auth login
    gcloud config set project <YOUR_PROJECT_ID>
    ```

3. Setup a firebase web app.

   1. Go to the firebase console. Follow the steps [here](https://firebase.google.com/docs/projects/use-firebase-with-existing-cloud-project#how-to-add-firebase_console).
   2. Create a new firebase web app and note down the values of the **APP ID**, **API KEY** and **AUTH DOMAIN**. You'll need them in the next step.

4. Run the interactive environment setup script where you will be prompted to enter values like **PROJECT_ID** and **REGION** (which region do you want to consume the models from) and the firebase web app's config values that you noted in the previous step.

    ```sh
    chmod +x ./startup/configure_env.sh
    ./startup/configure_env.sh
    ```

5. Make sure the required APIs are enabled and create a service account. You will need owner level access to the project. This enables the required APIs and creates the necessary service account with roles.


    ```sh
    chmod +x ./startup/setup_cloud_simple.sh
    ./startup/setup_cloud_simple.sh
    ```


6. Start the app.

    ```sh
    chmod +x ./startup/launch_app.sh
    ./startup/launch_app.sh
    ```

7. Access the Frontend Application Open <http://localhost:8080> in your browser. Use **0000** as the invite code.

8. To stop the app, press **Ctrl+C** in the terminal. Then run
  
    ```sh
    ./startup/launch_app.sh --stop
    ```

### Cloud setup with Firebase Authentication

WIP

#### Populate the database (Optional)

To update the data in the database, you can run the indexer.

1. Start the local application by following the steps outlined here: [Run the application](#simple-local-setup) step. (The indexer requires the database to be running which this step does for you).

2. Run the javascript indexer so it can add movies data into the database. The database comes pre-populated with the required data, but you can choose to re-add the data. The execution of this intentionally slowed down to stay below the rate-limits.

    ```sh
    docker compose -f docker-compose.indexer.yaml up --build -d 
    ```

    This takes about 10-15 minutes to run, so be patient. The embedding creation process is slowed down intentionally to ensure we stay under the rate limit.

3. Shut down the indexer container.

    ```sh
    docker compose -f docker-compose-indexer.yaml down
    ```

4. Verify the number of entries in the DB by using **Adminer** that is running at <http://localhost:8082>. Use the **minimal** user credentials (user name: minimal-user, password: minimal). Make sure you set System as PostgresSQL and Server as **db**, and Database as **fake-movies-db**.

There should be **652** entries in the movies table.

  ```sql
  SELECT COUNT(*)
  FROM "movies";
  ```

This ensures that all the movies are updated in the db.
