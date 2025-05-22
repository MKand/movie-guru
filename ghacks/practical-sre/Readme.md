

```sh
source .env
gcloud auth print-access-token | helm registry login -u oauth2accesstoken \
--password-stdin https://${REGION}-docker.pkg.dev
```

```sh
    cd ./ghacks/practical-sre/deploy/app/helm && helm package movie-guru
```

```sh
helm push movie-guru-0.3.0.tgz oci://${REGION}-docker.pkg.dev/${PROJECT_ID}/movie-guru/movie-guru
```