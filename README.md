# Zwitter
zwitter is a micro blogging platform built with microservice architecture powered by go, grpc, envoy and reactjs.


### Running locally
To run zwitter locally you will need to have go installed in your system first. Go to the [official documentation](https://golang.org/doc/install) for installing golang in your system and download.
Also [install nodejs]() and [install yarn](https://classic.yarnpkg.com/en/docs/install#debian-stable) in your system for js package management.

After that clone and build this project locally.
```bash
git clone https://github.com/dpakach/zwitter.git

cd zwitter

## Create default config files for the services
make config

make
```
This will pull some docker images and build the images locally so it might take some time.

This will clone the project on your local machine, build the binaries of the microservices and build docker images for the services.

With this you can run the project using docker-compose
For that just run the following command

``` bash
make run
```

Open your browser and visit http://localhost:8080.

You should be see the home page of zwitter.

### Setting up dev environment
To setup dev environment, just run following command
```
make dev
```

This will run webpack-dev-server so hot reloading is enabled in all the frontend files.

To reload after the changes made on go files, first build the service after making the changes and then restart the service.

for eg. If you do some changes in the posts service, first build the service with following commands
```
cd posts
make
```
Now restart he service with following command
```
make restart-service SERVICE=posts
```

You will also need to install go protobuf dependency along with [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) or the code generation on the server side. First you will need to have the protoc compiler for protobuf files. You can installit following [this](https://developers.google.com/protocol-buffers/docs/gotutorial#compiling-your-protocol-buffers) guide. After that install grpc-gateway and its required dependencies from [here](https://github.com/grpc-ecosystem/grpc-gateway). With all the dependencies out of the way, you should be ready to go.

### Architecture
Zwitter is composed of five different microservice components. These services are containerized using docker and can be run using docker-composer to start the whole system.

Each of the services have an instance of [envoy proxy](https://www.envoyproxy.io/) running inside their containers, which provides the http entrypoint for the services. All the services then connect to an seperate instance of envoy proxy which provides access to all the services through an single http entrypoint. The web service is mounted at root whereas all the other services are mounted at the endpoint starting their own name.

![zwitter architecture](https://github.com/dpakach/zwitter/blob/master/assets/zwitter.jpg?raw=true)


### Commands
Run `make help` to get the list of other available make commands
```
➜ make help

Zwitter
------------------------------------------

Available commands

all                            Build the binaries and docker-images for all the services
build                          Build the binaries and docker-images for all the services
clean                          Clean the build products for all the services
config                         Create default configuraion files
dev                            Start the development environment
docker-run                     Run zwitter with vanilla docker-compose setup
envoy-run                      Run zwitter using docker-compose with envoy proxy
fmt                            Format the go code with gofmt
frontend-watch                 Start frontend dev environment
help                           Display this help screen
js-deps                        Install javascript dependencies
logs                           Show docker compose logs for service $SERVICE, shows all logs if $SERVICE is not set
restart-service                Restart a service $SERVICE
service                        build the binary and docker image of $SERVICE service
stop                           Stop the services  
```

To get the list of available make commands for each service, go into the service folder and run `make help`

```
➜ make help

Zwitter:Auth 1.0
------------------------------------------

Available commands

api                            Auto-generate grpc go sources
clean                          Clean all build products
client                         Build the client binary
dep                            Download go dependencies
docker-build                   Build the docker image of the service
docker-run                     Build and run the service in a docker container
docs                           Generate swagger documentation files
help                           List available make commands
server                         Build the server binary   
```
## License

Copyright (c) 2020 Dipak Acharya

Licensed under [MIT License](https://github.com/dpakach/zwitter/blob/master/LICENSE)