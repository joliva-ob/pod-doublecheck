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
| ENV                        | Example value           | 
| -------------------------- |:----------------------- |
| LOG_FORMAT                 | "%{color}%{time:0102 15:04:05.000} %{level:.4s} %{id:03x} ▶ %{shortfunc}: %{color:reset} %{message}" | 
| EUREKA_APP_NAME            | pod-doublecheck | 
| spring_cloud_config_label  | dev |
| EUREKA_PUBLIC_HOST         | 10.200.2.28 |
| server_port                | 8080 |
| ENV                        | dev |
| spring_cloud_config_uri    | http://localhost:8888 |
| eureka_instance_ip_address | 10.1.51.167 |
| spring_profiles_active     | pro |
| spring_application_name    | pod-doublecheck |
| LOG_FILE                   | /opt/pod-doublecheck/log/pod-doublecheck.log |
| CONF_PATH                  | /opt/pod-doublecheck/conf/ |

## 