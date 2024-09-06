# Movie Guru


[Movie Guru](https://www.youtube.com/watch?v=l_KhN3RJ8qA)

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

### Components
* **Frontend**: The frontend is deployed on Firebase Hosting for easy deployment and scalability.
* **Backend Hosting:** The backend is deployed on Cloud Run for serverless execution and auto-scaling.
* **Caching:** Memorystore for Redis is used as a cache to improve performance and reduce latency for frequently accessed data.

### Agents
There are 3 agents used in this repo and are part of the backend. While they differ slightly in configuration from each backend type, they are mostly similar. All agents use a Gemini model through VertexAI APIs. 

This describes how the Go-Genkit backend agents works.
* **The User Profile / User Preferences Agent**: Used to analyse the user message and extract any long-lasting likes and dislikes from the conversation. 
* **The Query Transform Agent**: Analyses the last (max 10) messages in the history to extract the context and understand the user's latest message. For example, if the if the agent mentions, that it knows of 3 horror movies (movies A, B, C) and the user then asks to know more about "the last one", the query transform agent analyses this and states that the user's query is to know more about "movie C". The output of this agent is passed onto the retriever to retrieve relevant documents.
* **The Movie Agent**: Takes the information about the user's conversation, their profile, and the documents related to the context of the conversation and returns a response. The response consists of the answer, the justfication of hte answer, and finally a list of relevant movies that are related to the answer.

### Data
* The data about the movies is stored in CloudSQL pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
* The user's conversation history is stored in memory store for Redis. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.
* The user's profile data (their likes and dislikes) are stored in the CloudSQL database.

### CloudSQL
There are 5 tables:
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

### Function calling

#### GenKit
No additional steps required.

#### LangChain
```sh
export LANGCHAIN_TRACING_V2="false"
```

#### LangChain + LangSmith
Create an account and get an API key, then set the following environment variables:
```sh
export LANGSMITH_API_KEY=<api key>
export LANGCHAIN_TRACING_V2="true"
export LANGCHAIN_PROJECT=<project name>
export LANGCHAIN_ENDPOINT="https://api.smith.langchain.com"  # Double check with your project.
```

### Clone the project

```sh
git clone https://github.com/manasakandula/movie-guru.git
cd movie-guru
```

### Start the Deploy
```sh
./deploy/deploy.sh --skipapp --backend genkit-go  # or --backend langchain or --backend genkit-js (WIP)
```
We add --skipapp to make sure we wait for the db and the data are created before we deploy the application. 

### Create and populate the database

#### Create tables
Connect to the sql db through the cloud sql studio (the db is running on a private IP and hence cannot be reached directly without the use of cloudsql proxy). The [CloudSQL studio](https://cloud.google.com/sql/docs/mysql/manage-data-using-studio) is the is the easiest way to connect to it. Another option while testing locally is to set [Authorized Networks](https://cloud.google.com/sql/docs/mysql/authorize-networks) and allow list the IP address of the machine you are working on.


```SQL
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
    context VARCHAR
);

CREATE TABLE invite_codes (
    code VARCHAR(255) PRIMARY KEY,   
    valid BOOLEAN NOT NULL DEFAULT TRUE, 
);

CREATE TABLE user_logins (
    email VARCHAR(255) PRIMARY KEY,
    login_count INT NOT NULL DEFAULT 0,
    special_login VARCHAR(255)  DEFAULT NULL; 
);

CREATE TABLE app_metadata (
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

```
#### Insert data into the tables

Insert some data into the tables. Make changes and add the right values where required. You can play around with the values. 
The **CORS origin** should be the allowed front end domains (comma seperated list) from which your backend recieves calls. 
The **token audience** is the firebase project id from which the auth tokens are generated. 
The **history length** is the number of the most recent messages that are stored in cache.
The **max_user_message_len** is the max number of allowed characters in the user's message.
The current setup uses the same model for the 3 agents and is passed along as **google_chat_model name**.

```SQL
INSERT INTO app_metadata (app_version, token_audience, history_length, max_user_message_len, cors_origin, retriever_length, google_chat_model_name, google_embedding_model_name, front_end_domain)
VALUES ('v1', <project id> , 10, 1000, "https://<PROJECT ID>.web.app", 10, 'gemini-1.5-flash-001', 'text-embedding-004', "https://<project id>.web.app/");

INSERT INTO invite_codes (code, valid)
VALUES (<secret invite code>, TRUE);
```

#### Insert data into movies tables

WIP

### Build and deploy the app

```sh
./deploy/deploy.sh --skipinfra --backend genkit-go  # or --backend langchain or --backend genkit-js (WIP)
```

## License

This code of the repo is licensed under the Apache 2.0 License. To view a copy of this license, visit https://opensource.org/licenses/Apache-2.0 
This AI generated movie data and posters in the repo are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/   

