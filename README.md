# Go Shell Command API

## About

The GO Shell Command API is a simple HTTP REST API built with GO that allows users to execute shell commands remotely. It provides a secure and convenient way to execute various shell commands and retrieve output as an API response.

## Architecture

This implements `Clean Architecture` which helps to separate concerns by organizing code into several layers with a very explicit rule which enables us to create a testable and maintainable project.  [The Clean Architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html).
For this project we only have the presentation, usecase, entity and configs layers.

### 1. Presentation

This represents logic that consume the business logic from the `Usecase Layer`
and renders to the view. Here you can choose to render the view in e.g `rest`

### 2. Usecases

The code in this layer contains application specific business rules.
This represents the pure business logic of the application.
The rules of the application also shouldn't rely on the UI or the persistence frameworks being used.

### 3. Entities

 Defines the core domain entities, which represents the user entity with its properties and behavior.

### 4. Configs

Holds configuration files and related code, such as config.go, which handles the parsing and retrieval of application configuration settings.

## Technologies

1. Golang

## How to use it

1. First, clone the code.

    ```bash
        dev@dev:~$ git clone git@github.com:VicOsewe/go-shell-command-api.git
    ```

2. The next step is setting up your local environment by creating an `env.sh` and add your envs

    ``` bash
        export PORT="8080"
        export AUTH_PASSWORD="basic-auth-password"
        export AUTH_USERNAME="basic-auth-username"
    ```

3. Install the dependencies and generally ensure code health by running these commands.

    ```bash
        dev@dev:~$ go mod tidy
        dev@dev:~$ go generate ./...
    ```

4. Run the server at `http://localhost:8080/`

    ```bash
        dev@dev:~$ go run main.go
    ```

## Features

- Accepts shell commands via query parameters or JSON body.
- Returns the output of the executed command along with a message and status code in the response body.
- Implements basic authentication to secure access to the API.
- Supports a wide range of shell commands to cater to different use cases.

### 1.1 API Specification

> **PS**: THis feature exposes a `RESTFul` API.

**Authorization** Basic Auth

| Key          | Value           |
| ------------ | --------------- |
| Username     | "username"      |
| Password     | "password"      |

**POST** create lead

{{baseURL}}/api/cmd

**Body** raw (json)

```json
    {
    "command": "pwd",
   }
```

**Response** raw (json)

```json
    {
        "message": "command retrieved successfully",
        "status_code": 200,
        "body":  "/home"
    }
```
