# Movie Guru

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
There are 3 agents used in this repo and are part of the backend. While they differ slightly in configuration from each backend type, they are mostly similar. 
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
* *fake-movies-table*: This contains the information about the fake movies and their embeddings. The data for the table is found in dataset/movies_with_posters.csv
* *user-preferences-table*: This contains the user's long term preferences profile information. 
* *app-metadata*: This is used to configure the backend and has information about the model version, cors setting etc.


## Getting Started
WIP

## License

This code of the repo is licensed under the Apache 2.0 License. To view a copy of this license, visit https://opensource.org/licenses/Apache-2.0 
This AI generated movie data and posters in the repo are licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/   

