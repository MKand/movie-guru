kubectl create secret generic initdb-sql --from-file=init.sql=./pgvector/init/init.sql 

https://github.com/deliveryhero/helm-charts/tree/master/stable/locust

```sh
gcloud container clusters get-credentials movie-guru-gke --region $LOCATION --project $PROJECT_ID
```

```sh
kubectl create namespace locust
kubectl create namespace movie-guru
kubectl create namespace otel-collector
kubectl create configmap loadtest-locustfile \
--from-file locust/locustfile.py \
--namespace locust
```

```sh
helm upgrade movie-guru ./k8s/movie-guru \
--set Config.Image.Repository=manaskandula \
--set Config.projectID=$PROJECT_ID \
--namespace movie-guru \
--install
```

```sh
helm repo add deliveryhero https://charts.deliveryhero.io/
helm repo update

helm upgrade locust deliveryhero/locust \
--set loadtest.name=movieguru-loadtest \
--set loadtest.locust_locustfile_configmap=loadtest-locustfile \
--set loadtest.locust_locustfile=locustfile.py \
--set service.type=LoadBalancer \
--set master.environment.CHAT_SERVER=http://server-service.movie-guru.svc.cluster.local \
--set master.environment.MOCK_USER_SERVER=http://mockuser-service.movie-guru.svc.cluster.local \
--set worker.environment.CHAT_SERVER=http://server-service.movie-guru.svc.cluster.local \
--set worker.environment.MOCK_USER_SERVER=http://mockuser-service.movie-guru.svc.cluster.local \
--namespace locust \
--install
```

```sh
 export SERVICE_IP=$(kubectl get svc --namespace locust locust --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo http://$SERVICE_IP:8089
```