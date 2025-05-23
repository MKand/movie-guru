
This Readme tells you how to create a hosting docker and helm repos for the movieguru helm charts and docker images.




```sh
source .env
gcloud auth print-access-token | helm registry login -u oauth2accesstoken \
--password-stdin https://${REGION}-docker.pkg.dev
```

```sh
    cd ./ghacks/practical-sre/deploy/app/helm && helm package movie-guru
```

```sh
helm push movie-guru-0.3.0.tgz oci://${REGION}-docker.pkg.dev/${PROJECT_ID}/movie-guru
```

helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts


helm upgrade --install otel open-telemetry/opentelemetry-collector \
  --namespace otel --create-namespace \
  --set image.repository="otel/opentelemetry-collector-contrib" \
  --set mode="deployment" \
  -f utils/metrics/otel.values.yaml
