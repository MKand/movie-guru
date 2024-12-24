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
* **Web Backend (Go):**
    * The backend hosts the webserver.
    * Provides an API for the frontend to interact with.
    * Communicates with the Flows backend to execute the AI tasks that are a part of the application.
* **Flows Backend (Go-Genkit/ JS-Genkit):** 
    * There are 2 options for the *flows*. The two versions have identical functionality and can be used interchangibly:
        * Go-Genkit 
        * JS-Genkit (WIP)
    * The flows run the AI flows using Genkit. See the **FLOW** section below for more information.
    * Connects to GenAI models (through VertexAI APIs) to chat with users.
    * Connects to the VectorDB (CloudSQL with pgvector) to search for movies and information about movies.

## Deployment
* **Frontend**: Deployed in a docker container.
* **Web Backend:** Deployed in a docker container.
* **Flows Backend:** Deployed in a docker container.
* **Cache:**  Redis cache is used as a cache to improve performance and reduce latency for frequently accessed data like chat history. Deployed in a docker container.
* **Database:** The movies with their embeddings and user preferences dbs are deplyed on a CloudSQL postgres db. Can also be deployed in a docker container.


### Flow
There are 3 agents used in this repo and are part of the backend. While they differ greatly in their roles, they are mostly similar in structure. All agents use a Gemini model through VertexAI APIs. 

This describes how the Go-Genkit backend agents works.
* **The User Profile / User Preferences Flow**: Used to analyse the user message and extract any long-lasting likes and dislikes from the conversation. 
* **The Query Transform Flow**: Analyses the last (max 10) messages in the history to extract the context and understand the user's latest message. For example, if the if the agent mentions, that it knows of 3 horror movies (movies A, B, C) and the user then asks to know more about "the last one", the query transform agent analyses this and states that the user's query is to know more about "movie C". The output of this agent is passed onto the retriever to retrieve relevant documents.
* **The Movie Flow**: Takes the information about the user's conversation, their profile, and the documents related to the context of the conversation and returns a response. The response consists of the answer, the justfication of hte answer, and finally a list of relevant movies that are related to the answer.
* **The Doc Retriever Flow**: Takes a user query and returns relevant documents from the vector database. This flow is responsible for generating a vector representation of the query, and returning the relevant documents from the PGVector database.
* **Indexer Flow**: This flow is run to add the *movies* data into the vector database. This flow parses each entry in the *dataset/movies_withy_posters.csv* file, restructures it, creates a vector embedding and adds the resturctured data with the embeddings to the Postgres PGVector database. This flow is only invoked once during the setup of the application.

### Data

* The data about the movies is stored in CloudSQL pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
* The user's profile data (their likes and dislikes) are stored in the CloudSQL database.
* The user's conversation history is stored in memory store for Redis. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.

### CloudSQL

There are 2 tables:

* *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.
* *user_preferences*: This contains the user's long term preferences profile information. 

## Getting Started

Make sure you have deployed the core infrastructure (mainly postgres db, service accounts, and DB users), and enabled the APIs for this application. We'll perform the first step *Steps for backend infra* in the **main** branch (instructions are added here for convenience).

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

### Create tables (if using Cloud SQL)

Connect to the sql db through the **Cloud Sql Studio**.
For ease of use, the terraform script when creating the db allows all IPs to access the db. This is not good practice in production.
We will use the **main** user for this step. If you followed the getting started step, you would have a postgres DB with a **main** user. The DB password for **main** user (and the **minimal-user**) is stored in the secret manager in the project.

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

The **minimal user** should be allowed to only read data from the movies table and read and write to the user_preferences table. This is the user we will use for the application. You can choose to just use the **main** user instead for all db activities. If that is the case, update the docker-compose files and the set_env_vars.sh script.

```SQL
GRANT SELECT ON movies TO "minimal-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON user_preferences TO "minimal-user";
```

### Create tables (Local DB)

There is a local version of the db with pre-populated data.
To use that instead run the following

```sh
docker compose -f docker-compose-pgvector.yaml up -d
```

If you navigate to *localhost:8082*, you can access the db via *Adminer*. Use the main user credentials (user name: main, password: mainpassword).

At this stage, there will be 2 tables, with no data. We will populate the table in the next steps.

We'll grant the necessary permissions to minimal-user.

```SQL
GRANT SELECT ON movies TO "minimal-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON user_preferences TO "minimal-user";
```

### Populating the movie table

First, edit the *set_env_vars.sh* file. Update the project id, to reflect the ID of the GCP project you are using. If you are using a CloudSQL db, then you'll need to also update the postgres host, user name, and passwords to the correct ones. If you are running the local db, then you can leave the DB related env variables as is.
Then run the following command.

```sh
source set_env_vars.sh
```

Next:

* Go to the project in the GCP console. Go to **IAM > Service Accounts**.
* Select the movie guru service account (movie-guru-chat-server-sa@<project id>.iam.gserviceaccount.com). * Create a new JSON key.
* Download the key and store it as **.key.json** in the root of this repo (make sure you use the filename exactly).

This key is used to authorise the **indexer** and the **flows** backend to use the VertexAI APIs needed to generate the embeddings and query the model.

Let's run the javascript indexer so it can add movies data into the database. The execution of this intentionally slowed down to stay below the ratelimits.

```sh
docker compose -f docker-compose-indexer.yaml up --build -d 
```

This takes about 10-15 minutes to run, so be patient. The embedding creation process is slowed down intentionally to ensure we stay under the rate limit.
You can run the command below to ensure there are **652** entries in the db.

```sql
SELECT COUNT(*)
FROM "movies";
```

Lets turn down the container.

```sh
docker compose -f docker-compose-indexer.yaml down
```

Once all the required data is added, it is time to run the application that consists of the **frontend**, the **webserver**, the **genkit flows** server and the **redis cache**. These will be running locally in containers. The servers communicate with the **postgres DB** running in the cloud.

## Run the app

Let us make sure the env variables are in the execution context of docker-compose.
Re-run this if you are using a new terminal. If you've run it once before, its ok to re-run it.
Make sure you replace the dummy values in the file with real ones before you run it.

```sh
source set_env_vars.sh 
```

Now, let's get the application running.

```sh
docker compose up --build
```

This should start the application. You can go to **http://localhost:5173** to visit the front end and interact with the application.

Once finished, run the following to take the application down.

```sh
docker compose down
```
