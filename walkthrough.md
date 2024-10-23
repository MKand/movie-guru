# Walkthrough

```sh
docker network create shared-network
docker compose -f docker-compose-data.yaml up -d
```

```sh
docker compose up --build 
```

Go to http://localhost:5173

Start chatting "I want to watch a horror film" and "I am in the mood for something scary"

View embeddings and movies

```sql
SELECT title, plot, genres, embeddings FROM movies
```

```sql
SELECT title, plot, actors, director, runtime_mins, rating, released, poster FROM movies
```

```sh
docker compose -f docker-compose-genkit.yaml up -d
docker compose -f docker-compose-genkit.yaml exec genkit sh
```

```sh
npm install --force # I know!!!
genkit start
```

Go to http://localhost:4000


## Vector based searches

Go VectorSearchFlow

Search for "horror movies"

Search for "cat movies"

View trace and view the embedding

Search for "movies released after 2005"

The results will be terrible.

## Mixed searches

Go to MixedSearchFlow

Search for "movies with cats"

Show results and traces>embedding

Search for "movies with rating less than 3"

Show results and traces>no embedding


## Mixed search agent

Go to Prompts/MixedSearchPrompt

Search for "movies with cats"

Show the prompt and output

Search for "movies with rating less than 3"

Show the prompt and output

## RAG

Go to Flow/RAGFlow

Say "show me horror movies"