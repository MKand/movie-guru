Critical User Journeys:

## Getting Movie Information:

SLI:
- Query Success Rate (percentage of successful movie information requests)
  - Chat Success Counter
- User Engagement Rate (percentage of chat turns where the user is engaged in the conversation)
  - Chat Outcome Counter (ENGAGED).
- User Sentiment Rate (percentage of chat turns where the user is positive in the conversation)
  - Chat Sentiment Counter (POSITIVE).
- Query Latency (time taken to return movie information): Done
- Example SLO: 99% of movie information queries should succeed within 500ms and have 90% POSITIVE User sentiment, and 60% ENGAGED outcome.

## Managing User Preferences:

SLI:
- Preference Update Success Rate (percentage of successful preference update actions)
- Preference Retrieval Latency (time taken to retrieve user preferences)
- Example SLO: 99.9% of preference update actions should succeed within 200ms. 99% of preference retrieval should occur within 100ms.


https://cloud.google.com/trace/docs/setup/go-ot#run-sample
