## Make sure the env variable PROJECT_ID is set.
echo $PROJECT_ID

## Get an access token
ACCESS_TOKEN=`gcloud auth print-access-token`

## Give the service a name
SERVICE_DISPLAY_NAME="chat-service"

## Create a custom service definition
CREATE_SERVICE_POST_BODY=$(cat <<EOF
{
  "displayName": "${SERVICE_DISPLAY_NAME}",
   "custom": {},
   "telemetry": {}
}
EOF
)


## POST to create the service
curl  --http1.1 --header "Authorization: Bearer ${ACCESS_TOKEN}" --header "Content-Type: application/json" -X POST -d "${CREATE_SERVICE_POST_BODY}" https://monitoring.googleapis.com/v3/projects/${PROJECT_ID}/services?service_id=${SERVICE_DISPLAY_NAME}
