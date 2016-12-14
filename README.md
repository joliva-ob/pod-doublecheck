# pod-doublecheck
Scheduled double check to ensure all the pods from a kubernetes cluster and namespace are registred properly to a Netflix Euereka service discovery.
The app will expose an endpoint with the last check result, basically number of diferences, service name and origin from it is missing.


## Run unit tests
+ go test

## Compiled with runtime with:
+ GOOS=windows GOARCH=386 go build -o pod-doublecheck.exe
+ GOOS=linux GOARCH=386 go build -o pod-doublecheck
+ GOOS=darwin GOARCH=386 go build -o pod-doublecheck

## Build Docker image with
+ cp /source_cfg_files/*env* .
+ docker build -f Dockerfile . -tag pod-doublecheck
+ docker run --publish 8000:8000 --name pod-doublecheck --rm pod-doublecheck --restart=always pod-doublecheck

## Kubernetes
+ docker build -t docker-registry.oneboxtickets.com/oneboxtm/pod-doublecheck:version .
+ docker push docker-registry.oneboxtickets.com/oneboxtm/pod-doublecheck:version

## Environment variables

| Var name                        | Example value           |
| -------------------------- |:----------------------- |
| LOG_FORMAT                 | "%{color}%{time:0102 15:04:05.000} %{level:.4s} %{id:03x} ▶ %{shortfunc}: %{color:reset} %{message}" |
| LOG_FILE                   | /opt/pod-doublecheck/log/pod-doublecheck.log |
| spring_application_name    | pod-doublecheck |
| spring_profiles_active     | pro |
| spring_cloud_config_uri    | http://localhost:8888 |
| spring_cloud_config_label  | dev |
| server_port                | 8080 |
| EUREKA_APP_NAME            | pod-doublecheck |
| EUREKA_PUBLIC_HOST         | 10.200.2.28 |
| REFRESH_TIME_SECONDS       | 300 |
| eureka_instance_ip_address | 10.1.51.167 |
| ENV                        | dev |

## Managing the application
There a couple of endpoint to manage the application, see the insights below:
+ GET http://localhost:8080/monitoring
Just to get the last metrics info like kubernetes pods retrieved, or eureka apps listed and the differences found if exists.
```
{
  "Eureka apps": {
    "Name": "Eureka apps",
    "Value": 61,
    "Threshold": 300,
    "Alert": false
  },
  "Kubernetes pods": {
    "Name": "Kubernetes pods",
    "Value": 0,
    "Threshold": 300,
    "Alert": false
  },
  "Pods not found": {
    "Name": "Pods not found",
    "Value": 0,
    "Threshold": 0,
    "Alert": false
  }
}
```
+ PUT http://localhost:8080/refreshtime?-1
This is to re-schedule the checking time interval by setting a number of seconds > 0, or a number <= 0 to disable the service.
+ On both cases an Authorization header is needed to set with
```
Bear 1736cc7f-7c60-4576-b851-b7b3630cfeab
```
