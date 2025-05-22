#!/bin/bash

set -euo pipefail

# === CONFIGURABLE VARIABLES ===
MOVIEGURU_IP=$(gcloud compute addresses describe movieguru-address --global --project "$GCP_PROJECT_ID" --format="value(address)")
MOCKSERVER_IP=$(gcloud compute addresses describe mockerserver-address --global --project "$GCP_PROJECT_ID" --format="value(address)")
MOVIEGURU_CHART="oci://us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru/movie-guru"
LOCUST_CHART="oci://ghcr.io/deliveryhero/helm-charts/locust"
REPO_PREFIX="us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru"
# === INSTALL MOVIEGURU CHART ===
echo "Installing MovieGuru chart..."
helm upgrade --install movie-guru "$MOVIEGURU_CHART" \
  --version "0.2.0" \
  --namespace movieguru \
  --set Config.Image.Repository="$REPO_PREFIX" \
  --set Config.serverAddress="http://movieguru.endpoints.${GCP_PROJECT_ID}.cloud.goog/server" \
  --set Config.mockserverIP="$MOCKSERVER_IP" \
  --set Gateway.IP="$MOVIEGURU_IP" \
  --set Config.projectID="$GCP_PROJECT_ID" 


# === INSTALL LOCUST CHART ===
echo "Installing Locust chart..."
helm upgrade --install locust "$LOCUST_CHART" \
  --version "0.31.6" \
  --namespace locust \
  --set loadtest.name="movieguru-loadtest" \
  --set loadtest.locust_locustfile_configmap="loadtest-locustfile" \
  --set loadtest.locust_locustfile="locustfile.py" \
  --set loadtest.locust_host="http://server-service.movie-guru.svc.cluster.local" \
  --set service.type="LoadBalancer" \
  --set worker.replicas=3 
  
  
# === WAIT FOR LOCUST SERVICE TO GET EXTERNAL IP ===
echo "Waiting for Locust service external IP..."
for i in {1..30}; do
  LOCUST_IP=$(kubectl get svc locust -n locust -o jsonpath="{.status.loadBalancer.ingress[0].ip}" 2>/dev/null || true)
  if [[ -n "$LOCUST_IP" ]]; then
    break
  fi
  echo "Waiting for external IP... ($i/30)"
  sleep 10
done

# === OUTPUTS ===
echo ""
echo "==================== Deployment Outputs ===================="
echo "locust_address:           http://${LOCUST_IP:-<pending>}:8089"
echo "movieguru_ip:            ${MOVIEGURU_IP}"
echo "movieguru_backend_address: http://movieguru.endpoints.${GCP_PROJECT_ID}.cloud.goog/server"
echo "============================================================"
