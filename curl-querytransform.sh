curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
  "history": [
    {
      "sender": "user",
      "message": "I want to watch a movie"
    }
  ],
  "userProfile": {
    "likes": {
      "actors": [
        ""
      ],
      "director": [
        ""
      ],
      "genres": [
        ""
      ],
      "other": [
        ""
      ]
    },
    "dislikes": {
      "actors": [
        ""
      ],
      "director": [
        ""
      ],
      "genres": [
        ""
      ],
      "other": [
        ""
      ]
    }
  },
  "userMessage": "I want to watch a movie"
}' \
  http://localhost:3403/queryTransformFlow