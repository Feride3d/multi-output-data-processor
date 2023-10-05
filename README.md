# Multi-Output Data Processor
## Overview
The Multi-Output Data Processor is a versatile HTTP API service designed to process and distribute JSON data with various output options. It allows you to specify the tag (type) of data, and it handles the distribution to different channels. This tool is suitable for scenarios where data needs to be routed to multiple destinations efficiently.


## How It Works

The service accepts JSON input with two key attributes:

```
{
"tag": "data tag",
"data": "data"
}
```
The service processes the `tag` and output the `data` to one or many of the following channels:
* stdout (standard output)
* stderr (standard error output)
* file (writes to a local file named `./examples/file.txt`)
* null (discards the data `/dev/null`)

The choice of the output channel is determined by a configuration defined in a local YAML file. The configuration specifies which tag corresponds to which output channels. For example:
```
config:
tag: 
    - info
outputs:
    - stdout

tag: 
    - debug
outputs:
    - file

tag: 
    - error
outputs: 
    - stderr
    - file

tag: 
    - trace
outputs:
    - "null"
```

## Input Validation
The service enforces input validation according to the following rules:
* The `tag` parameter is required and should be one of the tags available in the YML configuration file.
* The `data` parameter is required and should not be empty.

## Non-Blocking and Fault Tolerant
The output process is non-blocking (asynchronous) and does not impact the response code or time of the API. If there are issues with writing data to an output channel, the data is not lost. Instead, it's stored in a `./examples/dead-letter-queue.txt` file. The service also retries writing to output channels up to three times before failing.

## Running the Service
You can run the service in three different ways: (1) on your local machine, (2) in a Docker container, or (3) in a Docker container using Docker Compose (recommended if you plan to add more containers to the application). Follow the instructions below for each scenario:

### 1. Run the Service on Your Local Machine
To run the service on your local machine, use the following command in the project directory:
```
go run cmd/multi-output-data-processor/main.go
```

### 2. Run the Service in a Docker Container
To run the service in a Docker container, follow these steps:

Step 1: Build the Docker image for the service
```
docker build -t multi-output-data-processor .
```
Step 2: Run the service in a Docker container
```
docker run -p 8080:8080 multi-output-data-processor
```

### 3. Run the Service in a Docker Container using Docker Compose
To run the service in a Docker container using Docker Compose, follow these steps:

Step 1: Build the Docker image for the service
```
docker build -t multi-output-data-processor .
```
Step 2: Create and start containers using Docker Compose
```
docker-compose up -d
```

Step 3: Enter the Docker container's shell
```
docker exec -ti multi-output-data-processor bash
```
Step 4: Change the current directory to `/app` inside the container
```
cd /app
```
Step 5: Run the service in the running container
```
go run cmd/multi-output-data-processor/main.go
```
Now you can try making requests as described in the `Example of the request` section.

### Stopping and Removing the Docker Container Using Docker Compose
To stop and remove the Docker container, along with associated networks, images, and volumes, use the following command:
```
docker-compose down
```

## Run Unit Tests
To run unit tests, use the following command in the project directory:
```
go test -race -cover -count 100 ./internal/service
```

## Example of the request
___________
You can try an API to send a request (e.g., using Postman): 

* Method: POST
* Path: http://127.0.0.1:8080/process
* Body: 
```
{
"tag": "info",
"data": "This is a test data"
}
```
Feel free to explore the Multi-Output Data Processor by sending different requests.
