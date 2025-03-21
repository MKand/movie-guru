## Unique Service ID of an existing service
SERVICE_ID=${UNIQUE_SERVICE_ID}
## Get an access token
ACCESS_TOKEN=`gcloud auth print-access-token`

SLO_POST_BODY=$(cat <<EOF
{
  "displayName": "90% - Chat Success Rate - Calendar Day",
  "goal": 0.90,
  "calendarPeriod": "DAY",
  "serviceLevelIndicator": {
    "requestBased": {
      "goodTotalRatio": {
        "goodServiceFilter": "metric.type=\"prometheus.googleapis.com/movieguru_chat_calls_success_total/counter\" resource.type=\"prometheus_target\"",
        "totalServiceFilter": "metric.type=\"prometheus.googleapis.com/movieguru_chat_calls_total/counter\" resource.type=\"prometheus_target\""
      }
    }
  }
}
EOF
)
curl  --http1.1 --header "Authorization: Bearer ${ACCESS_TOKEN}" --header "Content-Type: application/json" -X POST -d "${SLO_POST_BODY}" https://monitoring.googleapis.com/v3/projects/${PROJECT_ID}/services/${SERVICE_ID}/serviceLevelObjectives
