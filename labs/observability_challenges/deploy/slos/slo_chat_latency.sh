## Unique Service ID of an existing service
SERVICE_ID=${UNIQUE_SERVICE_ID}
## Get an access token
ACCESS_TOKEN=`gcloud auth print-access-token`

CHAT_LATENCY_SLO_POST_BODY=$(cat <<EOF
{
  "displayName": "99% - Chat Latency - Rolling Day",
  "goal": 0.99,
  "rollingPeriod": "86400s",
  "serviceLevelIndicator": {
    "requestBased": {
      "distributionCut": {
            "distributionFilter": "metric.type=\"prometheus.googleapis.com/movieguru_chat_latency/histogram\" resource.type=\"prometheus_target\"",
            "range": {
              "min": -1000,
              "max": 8000
            }
          
        },
      }
    }
  }
EOF
)

curl  --http1.1 --header "Authorization: Bearer ${ACCESS_TOKEN}" --header "Content-Type: application/json" -X POST -d "${CHAT_LATENCY_SLO_POST_BODY}" https://monitoring.googleapis.com/v3/projects/${PROJECT_ID}/services/${SERVICE_ID}/serviceLevelObjectives
