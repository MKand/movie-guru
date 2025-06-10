


1. App works.
2. Push helm update using a wrong image tag.
    1. Get the gcp project id

    ```sh
    export gcp_project_id=<project id>
    ```

    2. Connect to cluster

    ```sh
    gcloud container clusters get-credentials movie-guru-gke --region us-central1 --project $gcp_project_id
    ```

    ```sh
        helm upgrade movie-guru-sre oci://us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru/movie-guru-sre \
        --install \
        --namespace movieguru \
        --version "3.0.0" \
        --create-namespace \
        --set Config.Image.Repository=us-central1-a-docker.pkg.dev/o11y-movie-guru/movie-guru \
        --set Config.Image.Tag="sre-e72ca10" \
        --set Config.gatewayAddress="movieguru.endpoints.${gcp_project_id}.cloud.goog" \
        --set Config.projectID=${gcp_project_id} \
        --set Config.geminiApiLocation=us-central1
    ```

    3. They use GKE playbooks to understand the image pull backoff error and fix the helm deploy command.
        Go to cluster/AppErrors -> Select error
        Select troubleshooting options. Enable cloud assist -> start chatting

    4. You realise that there is something wrong with the name of the image.
    5. You rollback to the previous version

        ```sh
        helm rollback movie-guru-sre 1
        ```
        



3. App is running smoothly. 
4. Structrued logging.
5. Logs analytics.