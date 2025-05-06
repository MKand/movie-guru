#!/usr/bin/env bash

# Verify that the script is being run on Linux
if [[ $OSTYPE != "linux-gnu" ]]; then
    echo -e "\e[91mERROR: This script is supported only on Linux. Please run it in a Linux environment.\e[0m"
    exit 1
fi

echo -e "\e[95m Stopping app with docker compose \e[0m"

docker compose down
