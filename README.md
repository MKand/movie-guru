# Movie Guru
[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/YOUR_VIDEO_ID)

 **BRANCH: ghack-version**: This version is a *minimal version* of the frontend and backend that doesn't have complex login logic like the version in **main**. It is meant to be run locally (except for the postgresDB, which is running in a cloudSQL instance on GCP). You could also choose to run the postgresdb locally.

## Description

Movie Guru is a website that helps users find movies to watch through an RAG powered chatbot. The movies are all fictional and are generated using GenAI. 
The goal of this repo is to explore the best practices when building AI powered applications.

This demo is *NOT* endorsed by Google or Google Cloud.  
The repo is intended for educational/hobbyists use only.

Refer to the readme in the **main** branch for more information. 


## Overall Architecture

The application follows a standard client-server model:

* **Frontend (Vue.js):**
    * Handles user interactions, displays movie suggestions, and manages the chat interface.
* **Backend (Go/Node.js/Python):**
    * There are 3 options for the backend. All backends have nearly identical functionality and can be used interchangibly:
        * Go-Genkit 
        * JS-Genkit (WIP)
    * The backend hosts the AI component of the application and the webserver.
    * Provides an API for the frontend to interact with.
    * Responsible for searching through the movie database, handling user requests, and managing user data and sessions.
    * Connects to GenAI models (through VertexAI APIs) to chat with users.
    * Connects to the VectorDB (CloudSQL with pgvector) to search for movies and information about movies.

## Deployment
* **Frontend**: Deployed in a docker container.
* **Backend:** Deployed in a docker container.
* **Cache:**  Redis cache is used as a cache to improve performance and reduce latency for frequently accessed data like chat history. Deployed in a docker container.
* **Database:** The movies with their embeddings and user preferences dbs are deplyed on a CloudSQL postgres db. Can also be deployed in a docker container.


### Flow
There are 3 agents used in this repo and are part of the backend. While they differ greatly in their roles, they are mostly similar in structure. All agents use a Gemini model through VertexAI APIs. 

This describes how the Go-Genkit backend agents works.
* **The User Profile / User Preferences Flow**: Used to analyse the user message and extract any long-lasting likes and dislikes from the conversation. 
* **The Query Transform Flow**: Analyses the last (max 10) messages in the history to extract the context and understand the user's latest message. For example, if the if the agent mentions, that it knows of 3 horror movies (movies A, B, C) and the user then asks to know more about "the last one", the query transform agent analyses this and states that the user's query is to know more about "movie C". The output of this agent is passed onto the retriever to retrieve relevant documents.
* **The Movie Flow**: Takes the information about the user's conversation, their profile, and the documents related to the context of the conversation and returns a response. The response consists of the answer, the justfication of hte answer, and finally a list of relevant movies that are related to the answer.

### Data
* The data about the movies is stored in CloudSQL pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
* The user's conversation history is stored in memory store for Redis. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.
* The user's profile data (their likes and dislikes) are stored in the CloudSQL database.

### CloudSQL
There are 2 tables:
* *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
* *user-preferences*: This contains the user's long term preferences profile information. 

## Getting Started
Make sure you have deployed the core infrastructure (mainly postgres db), and enabled the APIs for this application. We'll perform the first step *Steps for backend infra* in the **main** branch (instructions are added here for convenience).  
```sh
git checkout main
export PROJECT_ID=<set project id>
./deploy/deploy.sh --skipapp --backend genkit-go 
```
Once finished, go back to the ghack-version branch.

```sh
git checkout genkit-version
```

## Create and populate the database
### Create tables
Connect to the sql db through the cloud sql studio.
For ease of use, the terraform script when creating the db allows all IPs to access the db. This is not good practice in production.
The DB password for user **main** and user **minimal-user** is stored in the secret manager in the project.

```SQL
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS movies (
    tconst VARCHAR PRIMARY KEY,
    embedding VECTOR(768),
    title VARCHAR,
    runtime_mins INTEGER,
    genres VARCHAR,
    rating NUMERIC(3, 1),
    released INTEGER,
    actors VARCHAR,
    director VARCHAR,
    plot VARCHAR,
    poster VARCHAR,
    content VARCHAR
);

CREATE TABLE user_preferences (
  "user" VARCHAR(255) NOT NULL, 
  preferences JSON NOT NULL,
  PRIMARY KEY ("user")
);

```
## Reduce the permissions on the minimal-user 
The minimal user can read data from the movies table and read and write to the user_preferences table.

```SQL
GRANT SELECT ON movies TO "minimal-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON user_preferences TO "minimal-user";
```
## Insert data into movies tables
We'll be running the indexer flow chat_server_go/cmd/indexer/main.go to insert data into the datapace.
Go to *set_env_vars.sh*. 
Replace the values of the environment variables there and then run. (We use 2 db users, one is the *main* user and the other is a *minimal-user*. You can decide to skip the *minimal-user*, and use the *main* user for everything.)
```sh
source set_env_vars.sh
```
Next, go to the project in the GCP console. Go to **IAM > Service Accounts**. Select the movie guru service account (movie-guru-chat-server-sa@<project id>.iam.gserviceaccount.com). And create a new JSON key. Download the key and store it as **.key.json** in the root of this repo (make sure you use the filename exactly). This key is needed to authorise the indexer to use the VertexAI APIs needed to generate the embeddings.
Let's run the indexer so it can add movies data into the database.
```sh
docker compose -f docker-compose-indexer.yaml up --build -d 
```
This takes a about 5 minutes to run, so be patient. The embedding creation process is slowed down intentionally to ensure we stay under the rate limit.
You can run the command below to ensure there are **652** entries in the db.

```sql
SELECT COUNT(*)
FROM "movies";
```
Lets turn down the container. 
```sh
docker compose -f docker-compose-indexer.yaml down
```
Once all the required data is added, now it is time to run the application that consists of the **frontend**, the **webserver**, the **genkit flows** server and the **redis cache**. These will be running locally in containers. The webserver communicates with the postgres DB running in the cloud.
## Run the other the app

```sh
source set_env_vars.sh # Rerurn this if you are using a new terminal. If not, it doesn't hurt to rerun it to ensure the env variables are present.
docker compose -f docker-compose.yaml up --build
```
This should start the application. You can go to **http://localhost:5173** to visit the front end and interact with the application.

Once finished, run
```sh
docker compose -f docker-compose.yaml down
```

