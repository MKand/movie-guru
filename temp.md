kubectl create secret generic initdb-sql --from-file=init.sql=./pgvector/init/init.sql

kubectl create configmap my-loadtest-locustfile --from-file locust/locustfile.py
https://github.com/deliveryhero/helm-charts/tree/master/stable/locust