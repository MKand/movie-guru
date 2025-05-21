export PROJECT_ID=<the project id> # Change add add project id
export REGION="us-central1" # Change region as required
export POSTER_REGION="" # The region where the poster bucket is. Change region as required 
export GENKIT_FEEDBACK_REGION="" # Change region as required

export FIREBASE_GCP_ID="some value"
export FIREBASE_API_KEY="some value"
export FIREBASE_AUTH_DOMAIN="some value"
export FIREBASE_APPID="some value"

export USE_AUTH=true # Make this false if no Auth is required
export STRICT_CORS=false # App accepts any origin if false
export CORS_ORIGINS="localhost:8080" # if strict cors is false, then give a comma seperated list of accepted origins
export TOKEN_AUDIENCE=${PROJECT_ID}
