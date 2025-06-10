# Pushing Movie Guru Helm Chart

This guide describes how to package the `movie-guru` Helm chart (from the `movie-guru-observability-challenge` directory), push it to Google Artifact Registry, and set up an OpenTelemetry collector.

## Prerequisites

- Ensure you have an `.env` file in your current directory with `PROJECT_ID` and `REGION` variables defined.
- Google Cloud SDK (`gcloud`) installed and authenticated.
- Helm CLI installed.

## 1. Log in to Google Artifact Registry for Helm

Authenticate Helm with your Google Artifact Registry. This command uses the credentials from `gcloud`.

```sh
source .env
gcloud auth print-access-token | helm registry login -u oauth2accesstoken \
--password-stdin https://${REGION}-docker.pkg.dev
```

## 2. Package the Helm Chart

Navigate to the Helm chart's parent directory and package the chart. The chart is assumed to be in a directory named movie-guru-observability-challenge, and its Chart.yaml should define the chart name (e.g., movie-guru) and version (e.g., 1.0.0).

```sh
   cd ./labs/observability-challenges/deploy/app/helm && helm package movie-guru
```

This command will create a [chart-name]-[chart-version].tgz file (e.g., movie-guru-observability-lab-1.0.0.tgz) in the helm directory.

## 3. Push Helm Chart to Artifact Registry

Push the packaged Helm chart to your Google Artifact Registry. The chart will be pushed to an OCI repository named movie-guru-observability-challenge.

```sh
helm push movie-guru-observability-lab-1.0.0.tgz oci://${REGION}-docker.pkg.dev/${PROJECT_ID}/movie-guru

```

Note: Ensure the filename movie-guru-observability-lab-1.0.0.tgz matches the output of the helm package command.
