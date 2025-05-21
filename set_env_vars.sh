export PROJECT_ID=<the project id> # Change add add project id
export REGION=<your infra and models region> # Change region as required

# Optional
export GENKIT_FEEDBACK_REGION=<your region where the genkit feedback collection service is located> # Change region as required

# Leave this variable empty to disable genkit feedback. Only enable it if your genkit project has this feature enabled.
export GENKIT_FEEDBACK_URL= "" # Change to https://${GENKIT_FEEDBACK_REGION}-${PROJECT_ID}.cloudfunctions.net/ext-firebase-ai-user-engagement-collectEngagement  Leave empty to disable feedback 

# Don't change these variables
export FIREBASE_GCP_ID=${PROJECT_ID}
export USE_AUTH=false 
export STRICT_CORS=false 

