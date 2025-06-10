## Unique Service ID of an existing service
SERVICE_ID=${UNIQUE_SERVICE_ID}
## Get an access token
ACCESS_TOKEN=`gcloud auth print-access-token`


SLO_POST_BODY=$(cat <<EOF
{
  "displayName": "70% - Chat Sentiment Rate - Calendar day",
  "goal": 0.7,
  "calendarPeriod": "DAY",
  "serviceLevelIndicator": {
    "requestBased": {
      "goodTotalRatio": {
        "goodServiceFilter": "metric.type=\"prometheus.googleapis.com/movieguru_chat_sentiment_counter_total/counter\" resource.type=\"prometheus_target\" metric.labels.Sentiment=monitoring.regex.full_match(\"Positive|Neutral\")",
        "totalServiceFilter": "metric.type=\"prometheus.googleapis.com/movieguru_chat_sentiment_counter_total/counter\" resource.type=\"prometheus_target\""
      }
    }
  }
}
EOF
)

curl  --http1.1 --header "Authorization: Bearer ${ACCESS_TOKEN}" --header "Content-Type: application/json" -X POST -d "${SLO_POST_BODY}" https://monitoring.googleapis.com/v3/projects/${PROJECT_ID}/services/${SERVICE_ID}/serviceLevelObjectives