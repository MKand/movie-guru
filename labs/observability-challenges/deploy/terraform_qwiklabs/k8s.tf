data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  client_certificate     = base64decode(google_container_cluster.primary.master_auth.0.client_certificate)
  client_key             = base64decode(google_container_cluster.primary.master_auth.0.client_key)
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}

provider "helm" {
  kubernetes = {
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

resource "helm_release" "otel" {
  name             = "otel"
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart            = "opentelemetry-collector"
  namespace        = "otel"
  create_namespace = true

  set = [{
    name  = "image.repository"
    value = "otel/opentelemetry-collector-contrib"
    },
    {
      name  = "mode"
      value = "deployment"
  }]

  values = [
    data.http.otel_file.response_body
  ]
}

resource "helm_release" "movie_guru" {
  name             = "movie-guru"
  chart            = var.helm_chart
  namespace        = "movieguru"
  version          = var.helm_chart_version
  wait             = true
  create_namespace = true

  set = [
    {
      name  = "Config.Image.Repository"
      value = var.repo_prefix
    },
    {
      name  = "Config.Image.Tag"
      value = var.image_tag
    },
    {
      name  = "Config.gatewayAddress"
      value = "movieguru.endpoints.${var.gcp_project_id}.cloud.goog"
    },
    {
      name  = "Config.projectID"
      value = var.gcp_project_id
    },
        {
      name  = "Config.projectID"
      value = var.gcp_project_id
    },
    {
      name  = "Config.geminiApiLocation"
      value = var.vertexAI_model_location
  }]
}

resource "kubernetes_config_map" "loadtest_locustfile" {
  metadata {
    name = "loadtest-locustfile"
  }
  data = {
    "locustfile.py" = (
      data.http.locust_py_file.response_body
    )
  }

}

resource "helm_release" "locust" {
  name             = "locust"
  chart            = "oci://ghcr.io/deliveryhero/helm-charts/locust"
  namespace        = "default"
  version          = "0.31.6"
  create_namespace = true
  depends_on       = [kubernetes_config_map.loadtest_locustfile]
  wait             = true
  atomic           = true
  set = [
    {
      name  = "loadtest.name"
      value = "movieguru-loadtest"
    },
    {
      name  = "loadtest.locust_locustfile_configmap"
      value = "loadtest-locustfile"
    },
    {
      name  = "loadtest.locust_locustfile"
      value = "locustfile.py"
    },
    {
      name  = "loadtest.locust_host"
      value = "http://movieguru.endpoints.${var.gcp_project_id}.cloud.goog/server"
    },
    {
      name  = "service.type"
      value = "LoadBalancer"
    },
    {
      name  = "worker.replicas"
      value = "3"
    }
  ]
}

data "kubernetes_service" "locust" {
  metadata {
    name      = "locust"
    namespace = "default"
  }
  depends_on = [helm_release.locust]
}
