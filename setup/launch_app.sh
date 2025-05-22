#!/usr/bin/env bash

# Function to display usage instructions
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  --app <app_name>  Specify an application profile. Currently supports 'genkit-monitoring'."
    echo "                    If 'genkit-monitoring' is specified, 'docker-compose.ghack.genkitmonitoring.yaml' will be used."
    echo "  --stop            Stop and remove containers (runs 'docker compose down')."
    echo "  --silent          Run the app in silent mode."

    echo "  -h, --help        Display this help message."
    exit 1
}

source .env

# Default values
APP_NAME="default"
STOP_ACTION=false
DOCKER_COMPOSE_FILE_OPT=()
SILENT=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        --app)
        if [[ -z "$2" || "$2" == --* ]]; then
            echo -e "\e[91mERROR: --app option requires a value.\e[0m" >&2
            usage
        fi
        APP_NAME="$2"
        shift # past argument
        shift # past value
        ;;
        --stop)
        STOP_ACTION=true
        shift # past argument
        ;;
        --silent)
        SILENT=true
        shift # past argument
        ;;
        -h|--help)
        usage
        ;;
        *)    # unknown option
        echo -e "\e[91mERROR: Unknown option: $1\e[0m" >&2
        usage
        ;;
    esac
done

# Determine Docker Compose file based on --app
if [[ "$APP_NAME" == "genkit-monitoring" ]]; then
    echo -e "\e[96mUsing custom Docker Compose file: docker-compose.ghack.genkitmonitoring.yaml\e[0m"
    DOCKER_COMPOSE_FILE_OPT=("-f" "docker-compose.ghack.genkitmonitoring.yaml")
elif [[ "$APP_NAME" == "sre" ]]; then
    echo -e "\e[96mUsing custom Docker Compose file: docker-compose.ghack.practical-sre.yaml\e[0m"
    DOCKER_COMPOSE_FILE_OPT=("-f" "docker-compose.ghack.practical-sre.yaml")
elif [[ -n "$APP_NAME" ]]; then
    echo -e "\e[93mWarning: App name '$APP_NAME' provided, but no specific docker-compose file is configured for it. Using default.\e[0m"
fi

# Handle --stop action
if [ "$STOP_ACTION" = true ] ; then
    echo -e "\e[95mStopping application variant: ${APP_NAME} with docker compose...\e[0m"
    docker compose "${DOCKER_COMPOSE_FILE_OPT[@]}" down
    exit $?
else
    # Proceed with starting the application
    # Verify that the script is being run on Linux
    if [[ $OSTYPE != "linux-gnu" ]]; then
        echo -e "\e[91mERROR: This script is supported only on Linux. Please run it in a Linux environment.\e[0m"
        exit 1
    fi

    # Check if PROJECT_ID is set (only needed for starting)
    if [[ -z "$PROJECT_ID" ]]; then
        echo -e "\e[91mERROR: Please set the PROJECT_ID environment variable (e.g., export PROJECT_ID=<YOUR_PROJECT_ID>).\e[0m"
        exit 1
    fi
    echo -e "\e[93mUsing Project: $PROJECT_ID\e[0m"

    SERVICE_ACCOUNT_EMAIL="movie-guru-chat-server-sa@$PROJECT_ID.iam.gserviceaccount.com"

    if [[ ! -f "key.json" ]]; then
        echo -e "\e[95mCreating service account local key as key.json for $SERVICE_ACCOUNT_EMAIL\e[0m"
        gcloud iam service-accounts keys create key.json \
            --iam-account="$SERVICE_ACCOUNT_EMAIL"
        if [ $? -ne 0 ]; then
            echo -e "\e[91mERROR: Failed to create service account key.json. Please check permissions and if the service account exists.\e[0m"
            exit 1
        fi
    else
        echo -e "\e[95mService account local key 'key.json' already exists. Skipping creation.\e[0m"
    fi

    echo -e "\e[95mStarting application variant: ${APP_NAME} with docker compose...\e[0m"
    if [[ "$SILENT" == "true" ]]; then
        docker compose "${DOCKER_COMPOSE_FILE_OPT[@]}" up --build -d
    else
        docker compose "${DOCKER_COMPOSE_FILE_OPT[@]}" up --build
    fi
    exit $?
fi
