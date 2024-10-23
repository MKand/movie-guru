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
* **Database:** The movies with their embeddings and user preferences dbs are deployed in a docker postgres db.

### Flow

There are several flows used in this repo and are part of the backend. While they differ greatly in their roles, they are mostly similar in structure. All agents use a Gemini model through VertexAI APIs.

This describes how the Go-Genkit backend agents works.

* **The User Profile / User Preferences Flow**: Used to analyse the user message and extract any long-lasting likes and dislikes from the conversation.
* **The Query Transform Flow**: Analyses the last (max 10) messages in the history to extract the context and understand the user's latest message. For example, if the if the agent mentions, that it knows of 3 horror movies (movies A, B, C) and the user then asks to know more about "the last one", the query transform agent analyses this and states that the user's query is to know more about "movie C". The output of this agent is passed onto the retriever to retrieve relevant documents.
* **The Movie Flow**: Takes the information about the user's conversation, their profile, and the documents related to the context of the conversation and returns a response. The response consists of the answer, the justfication of hte answer, and finally a list of relevant movies that are related to the answer.
* **The Doc Retriever Flow**: Takes a user query and returns relevant documents from the vector database. This flow is responsible for generating a vector representation of the query, and returning the relevant documents from the PGVector database.
* **Indexer Flow**: This flow is run to add the *movies* data into the vector database. This flow parses each entry in the *dataset/movies_withy_posters.csv* file, restructures it, creates a vector embedding and adds the resturctured data with the embeddings to the Postgres PGVector database. This flow is only invoked once during the setup of the application.

### Data

* The data about the movies is stored in local pgVector database. There are around 600 movies, with a plot, list of actors, director, rating, genre, and poster link. The posters are stored in a cloud storage bucket.
* The user's profile data (their likes and dislikes) are stored in the database.
* The user's conversation history is stored in memory store for Redis. Only the most recent 10 messages are stored. This number is configurable. The session info for the webserver is also stored in memory store.

### PostgreSQL

There are 2 tables:

* *movies*: This contains the information about the AI Generated movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv. If you choose to host your own posters, replace the links in this file.

* *user_preferences*: This contains the user's long term preferences profile information. 

## Getting Started

Make sure you have deployed the core infrastructure (mainly postgres db, service accounts, and DB users), and enabled the APIs for this application. We'll perform the first step *Steps for backend infra* in the **main** branch (instructions are added here for convenience).

- Clone the repo.

    ```sh
    git clone https://github.com/MKand/movie-guru.git --branch ghack-v2
    cd movie-guru
    ```

Step 2:

- Open a terminal.
- Check if the basic tools we need are installed. Run the following command.

    ```sh
    docker compose version
    ```

- If it prints out a version number (>= 2.29) you are good to go.

Step 3:

- Set project ID as environment variable

    ```sh
    export PROJECT_ID="<enter project id>"
    ```

Step 4:

- Go to the project in the GCP console. Go to **IAM > Service Accounts**.
- Select the service account (movie-guru-chat-server-sa@##########.iam.gserviceaccount.com).

![IAM](images/IAM.png)

- Select **Create a new JSON key**.

![CreateKey](images/createnewkey.png)

- Download the key and store it as **.key.json** in the root of this repo (make sure you use the filename exactly).

> **Warning**: In production it is BAD practice to store keys in file. Applications running in GoogleCloud use serviceaccounts attached to the platform to perform authentication. The setup used here is simply for convenience.

- Create a shared network for all the containers. We will be running containers across different docker compose files so we want to ensure the db is reachable to all of the containers.

     ```sh
    docker network create db-shared-network
    docker compose -f docker-compose-setup.yaml up -d
     ```

Step 5:

- Lets setup the app.
  
  ```sh
  docker compose up --build -d
  ```

- Navigate to <http://localhost:5173>. This is your app.

Step 6:

- Lets set up the **Genkit Developer UI**.  From the root of the project directory run the following command.

    ```sh
    docker compose -f docker-compose-setup.yaml exec genkit sh
    ```

- We are going to *exec* into the genkit container we created in the **docker-compose-setup.yaml file**. The reason we are not using **genkit start** as a startup command for the container is that it has an interactive step at startup that cannot be bypassed. So, we will exec into the container and then run the command **genkit start**.

- This should open up a shell inside the container at the location **/app**.

> **Note**: In the docker compose file, we mount the local directory **js/flows-js** into the container at **/app**, so that we can make changes in the local file system, while still being able to execute genkit tools from a container.

- Inside the container, run

    ```sh
    npm install 
    genkit start
    ```

- You should see something like this in your terminal

    ```text
    Genkit CLI and Developer UI use cookies and similar technologies from Google
    to deliver and enhance the quality of its services and to analyze usage.
    Learn more at https://policies.google.com/technologies/cookies
    Press "Enter" to continue
    ```

- Then press **ENTER** as instructed (this an interactive step that needs to be performed to start the Genkit UI).
- This should start the genkit server inside the container at port 4000 which we forward to port **4000** to your host machine (in the docker compose file).

> **Note**: Wait till you see an output that looks like this. This basically means that all the Genkit has managed to: (1) load the necessary go dependencies, (2) build the go module and (3) load the genkit actions. This might take 30-60 seconds for the first time, and the process might pause output for several seconds before proceeding.
**Please be patient**.

```sh
> flow@1.0.0 build
> tsc
Starting app at `lib/index.js`...
Genkit Tools API: http://localhost:4000/api
Registering plugin vertexai...
[TRUNCATED]
Registering retriever: movies
Registering flow: movieDocFlow
Starting flows server on port 3400
    - /userProfileFlow
    - /queryTransformFlow
    - /movieQAFlow
    - /movieDocFlow
Reflection API running on http://localhost:3100
Flows server listening on port 3400
Initializing plugin vertexai:
[TRUNCATED]
Registering embedder: vertexai/textembedding-gecko@001
Registering embedder: vertexai/text-embedding-004
Registering embedder: vertexai/textembedding-gecko-multilingual@001
Registering embedder: vertexai/text-multilingual-embedding-002
Initialized local file trace store at root: /tmp/.genkit/8931f61ceb1c88e84379f345e686136a/traces
Genkit Tools UI: http://localhost:4000
```

- Once up and running, navigate to **<http://localhost:4000>** in your browser. This will open up the **Genkit UI**. It will look something like this:

    ![Genkit UI JS](images/genkit-js.png)

- This is the developer interface of Genkit. Using this interface, you can test out the Flows you have created, the prompts you have created, etc.

> **WARNING: Potential error message**: At first, the genkit ui might show an error message and have no flows or prompts loaded. This might happen if genkit has yet had the time to detect and load the necessary go files. If that happens, go to **js/flows-js/src/index.ts**, make a small change (add a newline) and save it. This will cause the files to be detected and reloaded.