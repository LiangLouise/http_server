# CSCD58 Final Project Report

## Project Goals

As we stated in our proposal, we intially planned to 

* support http 1.0 + 1.1 and pipelining;
* serve multiple types of files from the host machine, including .html, .css, .js, .txt, and .jpg;
* support http 2 over TCP if time allows.

We managed to accomplished the first two. Thus this is a HTTP server can serve the files over HTTP/1.0 or 1.1 protocol.
When choosing the HTTP/1.1, the server is capable to handle persistent connections and HTTP pipelining.

[![http1_x_connections](imgs/http1_x_connections.png)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Connection_management_in_HTTP_1.x)

> Source: [MDN Web doc](https://developer.mozilla.org/en-US/docs/Web/HTTP/Connection_management_in_HTTP_1.x)

The Project is implemented in [Golang](https://go.dev/), one of C family language and with great built-in concurrency and third-party modules management support.

## Project Implementation and Documentation

### Project Sturcture

The structure of the project is as below.

The entry of point the server or main program is under `/cmd/http1_server/`.

The packages in `/pkg` directory contain the logics including:

* `/pkg/config`: The codes used to load configurations for server from yaml file;
* `/pkg/fsService`: Support services to help manipulate file in the host machines;
* `/pkg/httpParser`: Help convert between raw HTTP message to Go objects
* `/pkg/httpProto`: Store the const data such as HTTP status code, text and HTTP version text;
* `/pkg/router`: The handler logics for handling incomming HTTP request;
* `/pkg/server`: API to create a new Server instance and run it.

```shell
$ tree .
.
├── Makefile
├── cmd
│   └── http1_server
│       └── main.go
├── pkg
│   ├── config
│   │   └── config.go
│   ├── fsService
│   │   └── fsService.go
│   ├── httpParser
│   │   ├── requestParser.go
│   │   └── responseParser.go
│   ├── httpProto
│   │   └── httpProto.go
│   ├── router
│   │   └── router.go
│   └── server
│       └── sever.go
├── go.mod
├── go.sum
├── main
├── mainconfig.yml
```

For more details of each methods and struct, please refer the doc in each of *.go file

### Server Work Flow

The work flow of the server can be illustrated as below:

![project_workflow](imgs/project_workflow.svg)

## Work Contribution

- Roy
  - Full implementation of `main.go`
  - Created skeleton of  `router.go`.
  - Full implementation of `server.go`
  - Created configuration files: `go.mod`, `go.sum`, `mainconfig.yml`, `config.go`.
  - Full implementation `fsService.go` 
  - Implementation of `httpProto.go` (http protocol part)
  - Helped in debugging the `router.go` and the `requestParser.go`
  - Created workflow image
- Cheng
  - Full implementation of `requestParser.go`
  - Full implementation of `responseParser.go`
  - Implementation of `httpProto.go` (response status code and text parts)
  - Major implementation of `router.go`

## Setup And Running

#### Prerequisite

Download and Install [Golang](https://go.dev/dl/)

#### Compile

Under root of project directory, run:

```shell
http_server user$ make
```

It will create a folder `bin` and an executable file `http1_server` inside it.

#### Configuration

Here is a sample mainconfig.yml, you can change the configuration to test server's performance and features

```yaml
Server:
  SERVER_PORT: 8080
  SERVER_HOST: "127.0.0.1"
  HTTP_VERSION: "HTTP/1.1"

RunTime:
  MAX_CONCURRENT_CONNECTIONS: 1024
  ENABLE_PESISTANT: True
  ENABLE_PIPELINING: True
  MAX_PIPELINING_NUMBER: 10
  TIMEOUT_DURATION: 10
```



#### Run

```shell
user$ ./path_to_executable ./path_to_mainconfig.yml
```



## Implementation Testing