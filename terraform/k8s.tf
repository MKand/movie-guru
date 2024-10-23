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

data "http" "locustfile" {
  url = var.locust_file
}

resource "helm_release" "movie_guru" {
  name  = "movie-guru"
  chart = var.helm_chart
  set {
    name  = "Config.Image.Repository"
    value = "manaskandula"
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
      data.http.locustfile.response_body
    )
  }
  depends_on = [helm_release.movie_guru]
}

resource "helm_release" "locust" {
  name      = "locust"
  chart     = "deliveryhero/locust"
  namespace = "locust"

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
    name  = "service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "worker.replicas"
    value = "3"
  }

  set {
    name  = "master.environment.CHAT_SERVER"
    value = "http://server-service.movie-guru.svc.cluster.local"
  }

  set {
    name  = "master.environment.MOCK_USER_SERVER"
    value = "http://mockuser-service.movie-guru.svc.cluster.local"
  }

  set {
    name  = "worker.environment.CHAT_SERVER"
    value = "http://server-service.movie-guru.svc.cluster.local"
  }

  set {
    name  = "worker.environment.MOCK_USER_SERVER"
    value = "http://mockuser-service.movie-guru.svc.cluster.local"
  }
    depends_on = [helm_release.movie_guru, kubernetes_config_map.loadtest_locustfile]
}

data "kubernetes_service" "locust" {
  metadata {
    name      = "locust"  
    namespace = "locust"   
  }
}

data "kubernetes_service" "backend" {
  metadata {
    name      = "server-service"  
    namespace = "movie-guru"   
  }
}