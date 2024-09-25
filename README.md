# Movie Guru


[![Movie Guru](https://img.youtube.com/vi/l_KhN3RJ8qA/0.jpg)](https://youtu.be/YOUR_VIDEO_ID)


**NOTE**: the repo is still in development.

## Description

Movie Guru is a website that helps users find movies to watch through an RAG powered chatbot. The movies are all fictional and are generated using GenAI. 
The goal of this repo is to explore the best practices when building AI powered applications.

This demo is *NOT* endorsed by Google or Google Cloud.  
The repo is intended for educational/hobbyists use only.


## Overall Architecture

The application follows a standard client-server model:

* **Frontend (Vue.js, Firebase Hosting):**
    * Handles user interactions, displays movie suggestions, and manages the chat interface.
* **Backend (Go/Node.js/Python):**
    * There are 3 options for the backend. All backends have nearly identical functionality and can be used interchangibly:
        * Go-Genkit 
        * JS-Genkit (WIP)
        * Langchain
    * The backend hosts the AI component of the application and the webserver.
    * Provides an API for the frontend to interact with.
    * Responsible for searching through the movie database, handling user requests, and managing user data and sessions.
    * Connects to GenAI models (through VertexAI APIs) to chat with users.
    * Connects to the VectorDB (CloudSQL with pgvector) to search for movies and information about movies.

## Deployment
* **Frontend**: The frontend is deployed on Firebase Hosting for easy deployment and scalability.
* **Backend Hosting:** The backend is deployed on Cloud Run for serverless execution and auto-scaling.
* **Caching:** Memorystore for Redis is used as a cache to improve performance and reduce latency for frequently accessed data.

## Agents
There are 3 agents used in this repo and are part of the backend. While they differ slightly in configuration from each backend type, they are mostly similar. All agents use a Gemini model through VertexAI APIs. 

This describes how the Go-Genkit backend agents works.
* **The User Profile / User Preferences Agent**: Used to analyse the user message and extract any long-lasting likes and dislikes from the conversation. 
* **The Query Transform Agent**: Analyses the last (max 10) messages in the history to extract the context and understand the user's latest message. For example, if the if the agent mentions, that it knows of 3 horror movies (movies A, B, C) and the user then asks to know more about "the last one", the query transform agent analyses this and states that the user's query is to know more about "movie C". The output of this agent is passed onto the retriever to retrieve relevant documents.
* **The Movie Agent**: Takes the information about the user's conversation, their profile, and the documents related to the context of the conversation and returns a response. The response consists of the answer, the justfication of hte answer, and finally a list of relevant movies that are related to the answer.


## Data
* The data about the movies is stored in CloudSQL pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
* The user's conversation history is stored in memory store for Redis. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.
* The user's profile data (their likes and dislikes) are stored in the CloudSQL database.

## CloudSQL
There are 3 tables:
* *fake-movies-table*: This contains the information about the fake movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.

* *user-preferences-table*: This contains the user's long term preferences profile information. 
* *app-metadata*: This is used to configure the backend and has information about the model version, cors setting etc.
* *User logins*: Keeps track of users that have logged in.
* *Invite codes*: Keeps track of valid invite codes.

## Getting Started
 
Set project ID
```sh
export PROJECT_ID=<set project id>
```
If you are using Langchain, go to Langsmith, create an account and get an API key. Set the following environment variables. You can also choose to not use langsmith. 
In case set LANGCHAIN_TRACING_V2 to false.
You can skip this step if you are using GenKIT.
Otherwise,

```sh
export LANGSMITH_API_KEY=<api key>
export LANGCHAIN_TRACING_V2="true"
export LANGCHAIN_PROJECT=<project name>
export LANGCHAIN_ENDPOINT="https://api.smith.langchain.com" # Double check with your project.
```

Clone the project

```sh
git clone https://github.com/manasakandula/movie-guru.git
cd movie-guru
```

### Steps for the backend infra
Start the Deploy
```sh
./deploy/deploy.sh --skipapp --backend genkit-go  # or --backend langchain or --backend genkit-js (WIP)
```
We add --skipapp to make sure we wait for the db and the data are created before we deploy the application. 

# Create and populate the database

## Create tables
Connect to the sql db through the cloud sql studio (the db is running on a private IP and hence cannot be reached directly without the use of cloudsql proxy). The [CloudSQL studio](https://cloud.google.com/sql/docs/mysql/manage-data-using-studio) is the is the easiest way to connect to it. Another option while testing locally is to set [Authorized Networks](https://cloud.google.com/sql/docs/mysql/authorize-networks) and allow list the IP address of the machine you are working on. 
For ease of use, the terraform script when creating the db allows all IPs to access the db. **Make sure** you delete that setting after you finish inserting data.

The DB password for user **main** is stored in the secret manager in the project under the name **postgres-main-user-secret**.


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

CREATE TABLE IF NOT EXISTS invite_codes (
    code VARCHAR(255) PRIMARY KEY,
    valid BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS user_logins (
    email VARCHAR(255) PRIMARY KEY,
    login_count INT NOT NULL DEFAULT 0,
    special_login VARCHAR(255) DEFAULT NULL
);


CREATE TABLE IF NOT EXISTS app_metadata (
    app_version VARCHAR(255) PRIMARY KEY,
    token_audience VARCHAR(255) NOT NULL,
    history_length INT NOT NULL,
    max_user_message_len INT NOT NULL,
    cors_origin VARCHAR(255) NOT NULL, 
    retriever_length INT NOT NULL,
    google_chat_model_name VARCHAR(255) NOT NULL,
    google_embedding_model_name VARCHAR(255) NOT NULL,
    front_end_domain VARCHAR(255) NOT NULL
);

CREATE TABLE user_preferences (
  "user" VARCHAR(255) NOT NULL, 
  preferences JSON NOT NULL,
  PRIMARY KEY ("user")
);

```
## Insert data into the tables

Insert some data into the tables. Make changes and add the right values where required. You can play around with the values. 
The **CORS origin** should be the allowed front end domains (comma seperated list) from which your backend recieves calls. 
The **token audience** is the firebase project id from which the auth tokens are generated. 
The **history length** is the number of the most recent messages that are stored in cache.
The **max_user_message_len** is the max number of allowed characters in the user's message.
The current setup uses the same model for the 3 agents and is passed along as **google_chat_model name**.

```SQL
INSERT INTO app_metadata (app_version, token_audience, history_length, max_user_message_len, cors_origin, retriever_length, google_chat_model_name, google_embedding_model_name, front_end_domain)
VALUES ('v1', '<project id>' , 10, 1000, 'https://<PROJECT ID>.web.app', 10, 'gemini-1.5-flash-001', 'text-embedding-004', 'https://<project id>.web.app/');

INSERT INTO invite_codes (code, valid)
VALUES ('<secret invite code>', TRUE);

GRANT SELECT ON movies TO "minimal-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON user_preferences TO "minimal-user";
```

## Insert data into movies tables
### GENKIT GO ###
If using genkit-go do the following:

From this folder, run the command 
```sh
cd chat_server_go/cmd/indexer

export PROJECT_ID=<project id>
export POSTGRES_DB_USER_PASSWORD=<password>
export POSTGRES_HOST=<db public ip>
export GCLOUD_PROJECT=<project id>

export GCLOUD_LOCATION="europe-west4" # or another region
export POSTGRES_DB_INSTANCE="movie-guru-db-instance"
export POSTGRES_DB_USER="main"
export TABLE_NAME="movies"
export APP_VERSION="v1"
export POSTGRES_DB_NAME="fake-movies-db"

go run main.go
```
This takes a about 20 minutes to run, so be patient. The embedding creation process is slowed down intentionally to ensure we stay under the rate limit.

You can run the command below to ensure there are **652** entries in the db.

```sql
SELECT COUNT(*)
FROM "movies";
```

### GENKIT JS ###
WIP
### LANGCHAIN ###
WIP


**IMPORTANT**: The terraform script allows the postgres DB access from all IPs 0.0.0.0/0. This is bad practice in production. So, after inserting the data, make sure you remove the aurthorized networks portion in the definition of the postgresdb (deploy/terraform/go-server-infra/postgres.tf or deploy/terraform/langchain-server-infra/postgres.tf). Remove the section below and rerun the deploy pipeline. Or you can also remove this setting from the DB from google cloud console.
```tf
authorized_networks {
        name            = "All Networks"
        value           = "0.0.0.0/0"
        expiration_time = "3021-11-15T16:19:00.094Z"
      }

```

## Build and deploy the app

```sh
./deploy/deploy.sh --skipinfra --backend genkit-go  # or --backend langchain or --backend genkit-js (WIP)
```

### Steps for the frontend hosted on firebase
Create a firebase project. And create a webapp. Navigate to the project settings and find the firebase configuration variables. You should see something that looks like this:
```sh
  apiKey: "abcdefghijklmnkopqrstuvwxyz12345890",
  authDomain: "<firebase project name>.firebaseapp.com",
  projectId: "<firebase project name>",
  storageBucket: "<firebase project name>.appspot.com",
  messagingSenderId: "1234567890",
  appId: "1:234567890:web:1234567890"
```
Navigate to **chat_client_vue/movie-agent** and create a **.env** file.
Create the following env variables to the file.

```.env
VITE_FIREBASE_API_KEY=<apiKey>
VITE_FIREBASE_AUTH_DOMAIN=<authDomain>
VITE_GCP_PROJECT_ID=<projectId>
VITE_FIREBASE_STORAGE_BUCKET=<storageBucket>
VITE_FIREBASE_MESSAGING_SENDERID=<messagingSenderId>
VITE_FIREBASE_APPID=<appId>
VITE_CHAT_SERVER_URL=<address

```


## License

This code of the repo is licensed under the Apache 2.0 License. To view a copy of this license, visit https://opensource.org/licenses/Apache-2.0 
This AI generated movie data and posters in the repo are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/   

