#!/bin/bash

# Script to configure environment variables for the Movie Guru Webserver

# Function to prompt for a general value
# Usage: prompt_for_value VAR_NAME "Prompt message" [DEFAULT_VALUE] [IS_MANDATORY_FLAG]
# IS_MANDATORY_FLAG can be "mandatory" or empty/anything else for optional.
prompt_for_value() {
    local var_name="$1"
    local prompt_message="$2"
    local default_value="$3"
    local is_mandatory="$4"
    local current_value="${!var_name}" # Get current value if already set from environment
    local input

    local prompt_display="$prompt_message"
    local details=""
    if [ -n "$current_value" ]; then
        details="current: $current_value"
    fi
    if [ -n "$default_value" ]; then
        if [ -n "$details" ]; then details+=", "; fi
        details+="default: $default_value"
    fi
    if [ -n "$details" ]; then
        prompt_display+=" ($details)"
    fi

    while true; do
        read -p "$prompt_display: " input

        if [ -z "$input" ]; then # User pressed Enter
            if [ -n "$default_value" ]; then
                export "$var_name"="$default_value"
                echo "  -> $var_name set to (default): '$default_value'"
                break
            elif [ -n "$current_value" ]; then
                export "$var_name"="$current_value"
                echo "  -> $var_name kept current value: '$current_value'"
                break
            elif [ "$is_mandatory" == "mandatory" ]; then
                echo "  Error: $var_name is mandatory. Please provide a value."
                # Loop again to re-prompt
            else
                export "$var_name"=""
                echo "  -> $var_name set to empty (optional, no default/input)."
                break
            fi
        else # User provided input
            export "$var_name"="$input"
            echo "  -> $var_name set to: '$input'"
            break
        fi
    done
}

# Function to prompt for a boolean value (yes/no)
# Usage: prompt_for_boolean VAR_NAME "Prompt message" [DEFAULT_BOOLEAN_VALUE_AS_STRING_true_or_false]
# Sets the environment variable to "true" or "false".
prompt_for_boolean() {
    local var_name="$1"
    local prompt_message="$2"
    local default_bool_str="$3" # "true" or "false"
    local current_value_env="${!var_name}"
    local answer

    local default_display_char=""
    if [ "$default_bool_str" == "true" ]; then
        default_display_char="y"
    elif [ "$default_bool_str" == "false" ]; then
        default_display_char="n"
    fi

    local prompt_display="$prompt_message"
    local details_bool=""

    if [ -n "$current_value_env" ]; then
        local current_display_val="<$current_value_env>" # Show raw value if not true/false
        if [ "$current_value_env" == "true" ]; then current_display_val="yes";
        elif [ "$current_value_env" == "false" ]; then current_display_val="no"; fi
        details_bool="current: $current_display_val"
    fi

    if [ -n "$default_display_char" ]; then
        if [ -n "$details_bool" ]; then details_bool+=", "; fi
        details_bool+="default: $default_display_char"
    fi

    local yn_prompt="[y/n]"
    if [ -n "$details_bool" ]; then
        prompt_display+=" ($details_bool) $yn_prompt"
    else
        prompt_display+=" $yn_prompt"
    fi

    while true; do
        read -p "$prompt_display: " answer
        local answer_lower=$(echo "$answer" | tr '[:upper:]' '[:lower:]')

        if [ -z "$answer_lower" ]; then # User pressed enter
            if [ "$default_bool_str" == "true" ]; then
                export "$var_name"="true"
                echo "  -> $var_name set to true (default)."
                break
            elif [ "$default_bool_str" == "false" ]; then
                export "$var_name"="false"
                echo "  -> $var_name set to false (default)."
                break
            elif [ "$current_value_env" == "true" ] || [ "$current_value_env" == "false" ]; then
                export "$var_name"="$current_value_env"
                echo "  -> $var_name kept current value: '$current_value_env'"
                break
            else
                # Go's ParseBool treats empty or invalid as false.
                export "$var_name"="false"
                echo "  -> $var_name set to false (no default/input, defaulting to false)."
                break
            fi
        elif [ "$answer_lower" == "y" ] || [ "$answer_lower" == "yes" ]; then
            export "$var_name"="true"
            echo "  -> $var_name set to true."
            break
        elif [ "$answer_lower" == "n" ] || [ "$answer_lower" == "no" ]; then
            export "$var_name"="false"
            echo "  -> $var_name set to false."
            break
        else
            echo "  Invalid input. Please enter 'y' or 'n'."
        fi
    done
}

USE_AUTH=false
STRICT_CORS=false
TOKEN_AUDIENCE=""
CORS_ORIGINS=""
MAX_USER_MESSAGE_LENGTH=500
HISTORY_LENGTH=10
FIREBASE_API_KEY=""
FIREBASE_APP_ID=""
FIREBASE_AUTH_DOMAIN=""
GENKIT_FEEDBACK_URL=""
FIREBASE_API_KEY=""
FIREBASE_APP_ID=""
FIREBASE_AUTH_DOMAIN=""
POSTER_BUCKET_NAME="generated_posters"
FLOWS_URL="http://flows:3400"
ENABLE_METRICS=false

echo "--- Movie Guru Webserver Environment Configuration ---"
echo "This script will guide you through setting necessary environment variables."
echo "If a variable is already set, its 'current' value will be shown."
echo "Press Enter to accept the default value (if shown), or the current value if no default."
echo ""

# General Settings
echo "General Settings:"
prompt_for_value "PROJECT_ID" "1. Google Cloud Project ID" "" "mandatory"
prompt_for_value "REGION" "2. Google Cloud Region. defaults to us-central1" "us-central1"


echo ""
echo "--- Configuration Summary ---"
echo "The following environment variables have been configured in this script's session:"
echo "POSTER_BUCKET_NAME='${POSTER_BUCKET_NAME}'"
echo "FLOWS_URL='${FLOWS_URL}'"
echo "GENKIT_FEEDBACK_URL='${GENKIT_FEEDBACK_URL}'"
echo "ENABLE_METRICS='${ENABLE_METRICS}'"
echo "USE_AUTH='${USE_AUTH}'"

if [ "${USE_AUTH}" == "true" ]; then
    echo "TOKEN_AUDIENCE='${TOKEN_AUDIENCE}'"
    echo "CORS_ORIGINS='${CORS_ORIGINS}'"
    echo "STRICT_CORS='${STRICT_CORS}'"
    echo "MAX_USER_MESSAGE_LENGTH='${MAX_USER_MESSAGE_LENGTH}'"
    echo "HISTORY_LENGTH='${HISTORY_LENGTH}'"
fi
echo "-----------------------------"
echo ""
echo "To apply these settings in your current terminal session, run:"
echo "  source ./configure_env.sh"
echo ""
echo "Alternatively, you can prepend them to your command, for example:"
echo "  POSTER_BUCKET_NAME='${POSTER_BUCKET_NAME}' FLOWS_URL='${FLOWS_URL}' \\"
echo -n "  ENABLE_METRICS='${ENABLE_METRICS}' USE_AUTH='${USE_AUTH}'"
AUTH_EXAMPLE_PART=""
if [ "${USE_AUTH}" == "true" ]; then
    AUTH_EXAMPLE_PART=" \\\n  TOKEN_AUDIENCE='${TOKEN_AUDIENCE}' CORS_ORIGINS='${CORS_ORIGINS}'"
    AUTH_EXAMPLE_PART+=" \\\n  STRICT_CORS='${STRICT_CORS}' MAX_USER_MESSAGE_LENGTH='${MAX_USER_MESSAGE_LENGTH}' HISTORY_LENGTH='${HISTORY_LENGTH}'"
fi
echo -n -e "${AUTH_EXAMPLE_PART}" # -e to interpret \n if present
echo " \\"
echo "  go run code/webserver/cmd/webserver/main.go  # Adjust path if needed"
echo ""
echo "You can also save these settings to a .env file for use with tools like Docker Compose or some auto-loaders:"
echo "POSTER_BUCKET_NAME='${POSTER_BUCKET_NAME}'"
echo "FLOWS_URL='${FLOWS_URL}'"
echo "GENKIT_FEEDBACK_URL='${GENKIT_FEEDBACK_URL}'"
echo "ENABLE_METRICS='${ENABLE_METRICS}'"
echo "USE_AUTH='${USE_AUTH}'"
if [ "${USE_AUTH}" == "true" ]; then
    echo "TOKEN_AUDIENCE='${TOKEN_AUDIENCE}'"
    echo "CORS_ORIGINS='${CORS_ORIGINS}'"
    echo "STRICT_CORS='${STRICT_CORS}'"
    echo "MAX_USER_MESSAGE_LENGTH='${MAX_USER_MESSAGE_LENGTH}'"
    echo "HISTORY_LENGTH='${HISTORY_LENGTH}'"
fi

# Actual .env file creation
cat <<EOF > .env
PROJECT_ID='${PROJECT_ID}'
REGION='${REGION}'
POSTER_BUCKET_NAME='${POSTER_BUCKET_NAME}'
FLOWS_URL='${FLOWS_URL}'
ENABLE_METRICS='${ENABLE_METRICS}'
USE_AUTH='${USE_AUTH}'
TOKEN_AUDIENCE='${TOKEN_AUDIENCE}'
CORS_ORIGINS='${CORS_ORIGINS}'
STRICT_CORS='${STRICT_CORS}'
MAX_USER_MESSAGE_LENGTH='${MAX_USER_MESSAGE_LENGTH}'
HISTORY_LENGTH='${HISTORY_LENGTH}'
GENKIT_FEEDBACK_URL='${GENKIT_FEEDBACK_URL}'
FIREBASE_API_KEY='${FIREBASE_API_KEY}'
FIREBASE_APP_ID='${FIREBASE_APP_ID}'
FIREBASE_AUTH_DOMAIN='${FIREBASE_AUTH_DOMAIN}'
EOF
echo "A .env file has been generated/updated with these settings."
