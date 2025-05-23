data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  client_certificate     = base64decode(google_container_cluster.primary.master_auth.0.client_certificate)
  client_key             = base64decode(google_container_cluster.primary.master_auth.0.client_key)
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host                   = "https://${google_container_cluster.primary.endpoint}"
    token                  = data.google_client_config.default.access_token
    client_certificate     = base64decode(google_container_cluster.primary.master_auth.0.client_certificate)
    client_key             = base64decode(google_container_cluster.primary.master_auth.0.client_key)
    cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  }
}

resource "google_compute_global_address" "movieguru-address" {
  name         = "movieguru-address"
  address_type = "EXTERNAL"
  project      = var.gcp_project_id
}

resource "google_compute_global_address" "mockserver-address" {
  name         = "mockerserver-address"
  address_type = "EXTERNAL"
  project      = var.gcp_project_id
}

resource "google_endpoints_service" "openapi_service" {
  service_name = "movieguru.endpoints.${var.gcp_project_id}.cloud.goog"
  project      = var.gcp_project_id
  openapi_config = yamlencode({
    swagger = "2.0"
    info = {
      description = "Cloud Endpoints service for MovieGuru"
      title       = "MovieGuru"
      version     = "1.0.0"
    }
    paths = {}
    host  = "movieguru.endpoints.${var.gcp_project_id}.cloud.goog"
    x-google-endpoints = [
      {
        name   = "movieguru.endpoints.${var.gcp_project_id}.cloud.goog"
        target = google_compute_global_address.movieguru-address.address
      },
    ]
  })
}

data "http" "locust_py_file" {
  url = var.locust_py_file
}

data "http" "sql_file" {
  url = var.sql_file
}

data "http" "otel_file" {
  url = var.otel_file
}

resource "kubernetes_namespace" "locust" {
  metadata {
    name = "locust"
  }
}


resource "helm_release" "otel"{
  name      = "otel"
  repository = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart     = "opentelemetry-collector"
  namespace = "otel"
  create_namespace = true
  set{
    name = "image.repository"
    value = "otel/opentelemetry-collector-k8s"
  }
  set{
    name = "mode"
    value = "deployment"
  }

    values = [
    data.http.otel_file.response_body
  ]
}

resource "helm_release" "movie_guru" {
  name      = "movie-guru"
  chart     = var.helm_chart
  namespace = "movieguru"
  version   = "0.3.0"
  wait      = false
  create_namespace = true

  set {
    name  = "Config.Image.Repository"
    value = var.repo_prefix
  }
  set {
    name  = "Config.serverAddress"
    value = "http://movieguru.endpoints.${var.gcp_project_id}.cloud.goog/server"
  }

  set {
    name  = "Config.mockserverIP"
    value = google_compute_global_address.mockserver-address.address
  }

  set {
    name  = "Gateway.IP"
    value = google_compute_global_address.movieguru-address.address
  }

  set {
    name  = "Config.projectID"
    value = var.gcp_project_id
  }
}




resource "kubernetes_config_map" "loadtest_locustfile" {
  metadata {
    name      = "loadtest-locustfile"
    namespace = "locust"
  }
  data = {
    "locustfile.py" = (
      data.http.locust_py_file.response_body
    )
  }

  depends_on = [kubernetes_namespace.locust]
}


resource "helm_release" "locust" {
  name      = "locust"
  chart     = "oci://ghcr.io/deliveryhero/helm-charts/locust"
  namespace = "locust"
  version   = "0.31.6"
  create_namespace = false

  set {
    name  = "loadtest.name"
    value = "movieguru-loadtest"
  }

  set {
    name  = "loadtest.locust_locustfile_configmap"
    value = "loadtest-locustfile"
  }

  set {
    name  = "loadtest.locust_locustfile"
    value = "locustfile.py"
  }
  set {
    name  = "loadtest.locust_host"
    value = "http://movieguru.endpoints.${var.gcp_project_id}.cloud.goog/server"
  }

  set {
    name  = "service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "worker.replicas"
    value = "3"
  }

  depends_on = [kubernetes_config_map.loadtest_locustfile]
}

data "kubernetes_service" "locust" {
  metadata {
    name      = "locust"
    namespace = "locust"
  }
  depends_on = [helm_release.locust]
}
