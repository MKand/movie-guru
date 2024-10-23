# Walkthrough

```sh
docker network create shared-network
docker compose -f docker-compose-data.yaml up -d
```

```sh
docker compose up --build 
```


Start chatting "I want to watch a horror film" and "I am in the mood for something scary"

View embeddings and movies

```sql
SELECT title, plot, genres, embeddings FROM movies
```

```sql
SELECT title, plot, actors, director, runtime_mins, rating, released, poster FROM movies
```

## Demo part 1

- Go to <http://localhost:5173>
- say "Hi"
- say "I love horror movies" and look at the movies
- say "I feel like watching something scary" and look at the movies
- Switch to slide 11

## Demo part 2

- Setup genkit

    ```sh
    docker compose -f docker-compose-genkit.yaml up -d
    docker compose -f docker-compose-genkit.yaml exec genkit sh
    ```

    ```sh
    npm install --force # I know!!!
    genkit start
    ```

- Go to <http://localhost:4000>

### Vector based searches

- Go VectorSearchFlow
- Search for "horror movies"
- Search for "cat movies"
- View trace and view the embedding
- Search for "movies released after 2005". The results will be terrible.
- Switch to slide 26

## Demo part 3

### Mixed searches

- Go to MixedSearchFlow in genkitUI
- Search for "movies with cats"
- Show results and traces>embedding
- Search for "movies with rating less than 3"
- Show results and traces>no embedding
- Switch to slide 36

## Demo part 4

### Mixed search agent

- Go to Prompts/MixedSearchPrompt
- Explain the elements of the view.
- Search for "movies with cats"
- Show the prompt and output
- Search for "movies with rating less than 3"
- Show the prompt and output
- Switch to slide 43

## Demo part 5

### RAG: Part 1

- Go to Flow/RAGFlow
- Say "show me horror movies"
- Show output and it should be good.
- Switch to slide 51

### RAG: Part 2

- Go to Flow/RAGFlow
- Copy the following into the input

```json
{
  "history": [
    {
      "role": "agent",
      "content": "hi"
    },
    {
      "role": "user",
      "content": "I want to watch a movie, but not sure what"
    },
    {
      "role": "agent",
      "content": "I can help you. What do you like? I can look for movies with those themes"
    },
    {
      "role": "user",
      "content": "cats.. i guess"
    }
  ],
  "userPreferences": {
  },
  "userMessage": "cats.. i guess"
}
```
