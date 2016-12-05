# pod-doublecheck
Scheduled double check to ensure all the pods from a kubernetes cluster and namespace are registred properly to a Netflix Euereka service discovery.
The app will expose an endpoint with the last check result, basically number of diferences, service name and origin from it is missing.


Run unit tests
+ go test

Compiled with runtime with: 
+ GOOS=windows GOARCH=386 go build -o pod-doublecheck.exe
+ GOOS=linux GOARCH=386 go build -o pod-doublecheck
+ GOOS=darwin GOARCH=386 go build -o pod-doublecheck

Build Docker image with
+ cp /source_cfg_files/*env* .
+ docker build -f Dockerfile . -tag pod-doublecheck
+ docker run --publish 8000:8000 --name pod-doublecheck --rm pod-doublecheck --restart=always pod-doublecheck 

Kubernetes
+ docker build -t docker-registry.oneboxtickets.com/oneboxtm/pod-doublecheck:version .
+ docker push docker-registry.oneboxtickets.com/oneboxtm/pod-doublecheck:version


