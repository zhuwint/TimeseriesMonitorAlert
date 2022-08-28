
kubectl apply -f ../deploy/namespace.yaml
kubectl apply -f ../deploy/influxdb-pvc.yaml

kubectl create configmap influxdb-config --from-file ../deploy/influxdb.conf -n time-series

kubectl apply -f ../deploy/influxdb-svc.yaml