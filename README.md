# Table of Contents

- [Table of Contents](#table-of-contents)
  - [About Movie Guru](#about-movie-guru)
  - [Description](#description)
    - [Overall Architecture](#overall-architecture)
      - [Components](#components)
    - [Deployment](#deployment)
      - [Docker Containers](#docker-containers)
    - [Genkit Flows and Prompts](#genkit-flows-and-prompts)
    - [Data](#data)
      - [Postgres](#postgres)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Environment setup](#environment-setup)
    - [Firebase setup](#firebase-setup)
    - [Database Setup](#database-setup)
      - [Run the database service](#run-the-database-service)
      - [Populate the database (Optional)](#populate-the-database-optional)
    - [Run the Application](#run-the-application)
    - [Clean up](#clean-up)

## About Movie Guru

**Genkit version**: 0.9.12 for Node.js

[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/l_KhN3RJ8qA)

 This version is a *minimal version* of the frontend and backend that doesn't have complex login logic like the version in **cloud-movieguru**. It is meant to be run fully locally while using VertexAI APIs.
 If you want to run this demo entirely in the cloud use the **cloud-movieguru** branch.

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

### Genkit Flows and Prompts

1. **Safety Prompt:** Checks each user statement to validate whether or not it is safe to proceed with.
1. **Query Transform Prompt:** Maps vague user queries to specific database queries.
1. **Movie Prompt:** Combines user data and relevant documents to provide responses.
1. **Movie Doc Flow:** Retrieves relevant documents from the vector database. Perform a keyword based, vector based, or mixed search based on the type of query.
1. **Chat Flow:** Combines all the aforementioned prompts flows into a single flow that is used by the chat server.
1. **User Profile Flow:** Additional flow that extracts user preferences from conversations.
1. **Indexer Flow:** Parses movie data and adds it to the vector database.

### Data

- The data about the movies is stored in a pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
- The user's profile data (their likes and dislikes) are stored in the CloudSQL database.
- The user's conversation history is stored in a local redis cache. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.

#### Postgres

There are 2 important tables:

- *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
- *user_preferences*: This contains the user's long term preferences profile information.

## Getting Started

### Prerequisites

- A Google Cloud project with owner permissions.
- Tools:
  - [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)
  - Docker and Docker Compose
- Required APIs enabled (will be performed in `setup_local.sh`).

### Get into Cloud Shell

- Navigate to the Cloud Console.
- Type Cloud Shell into the search bar at the top
- Open the Cloud Shell Editor - this will provision a VM for you and give you a canned environment to operate in.

### Clone the Repository

Open a terminal in the Cloud Shell Editor.

```sh
git clone https://github.com/MKand/movie-guru.git
cd movie-guru
git checkout schnecle-genkit-monitoring-dogfood
```

### Firebase setup

1. Navigate to the Firebase Console. 
1. Using the same project you are using in Cloud Shell, set up a new web app in the Firebase Console. Follow the steps [here](https://firebase.google.com/docs/projects/use-firebase-with-existing-cloud-project#how-to-add-firebase_console).
1. Create a new firebase web app and copy the firebase config variables into **set_env_vars.sh**.
1. Update your project to Blaze (pay-as-you-go) plan.

### Environment setup

1. Authenticate with Google Cloud (unnecessary if running from Cloud Shell)

    ```sh
    gcloud auth login
    gcloud config set project <YOUR_PROJECT_ID>
    ```

1. Run setup script.

    ```sh
    chmod +x setup_cloud.sh
    ./setup_cloud.sh
    ```

This enables the required APIs and creates the necessary service account with roles.

### Run the Application

1. Start the all docker containers

```sh
  chmod +x start_app.sh # if the first time running
  ./start_app.sh
```

1. Access the Frontend Application Open http://localhost:8080 in your browser.

### Clean up

Run the following commands:

```sh
  chmod +x stop_app.sh # if the first time running
  ./stop_app.sh
```
