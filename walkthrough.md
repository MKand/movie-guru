# Walkthrough

```sh
docker network create shared-network
docker compose -f docker-compose-data.yaml up -d
```

```sh
docker compose up --build 
```

View embeddings and movies

```sql
SELECT title, plot, genres, embeddings FROM movies
```