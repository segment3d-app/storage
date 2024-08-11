# Segment3D Storage

This project is a part of the storage service for the Segment3D App. It is built using Go, a statically typed, compiled programming language designed for simplicity and performance. The API provides core functionality for both the Segment3D Mobile and Web applications. For database management, the project uses PostgreSQL, and it leverages Docker for containerization, ensuring consistency across different development and production environments. Additionally, the project utilizes Air for live reloading during development, streamlining the development process by automatically restarting the server on code changes.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Usage](#usage)
  - [Running the Server](#running-the-server)
  - [Running with Docker Compose](#running-with-docker-compose)
  - [Running All Segment3d Services](#running-all-segment3d-services)
  - [API Documentation](#api-documentation)
- [License](#license)

## Getting Started

### Prerequisites

Ensure you have the following tools installed on your machine:

- [Go v1.21.6](https://go.dev/dl/)
- [Docker](https://hub.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- [Air](https://github.com/cosmtrek/air) (for live reloading during development)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/segment3d-app/storage.git
    cd storage
    ```

2.  **Install dependencies:**

    ```bash
    go get -u ./...
    ```

### Configuration

1. **Copy the `.env.example` file to `.env`:**

   ```bash
   cp .env.example .env
   ```

## Usage

### Running the Server

```bash
make server-dev
```

### Running with Docker Compose

You can also run the entire application using Docker Compose. This will set up the necessary containers for the backend and any other services you may have.

1. **Build and start the application using Docker Compose:**

   ```bash
   docker-compose up --build
   ```

   This command will build the Docker image and start the containers, and you can access the API at `http://localhost:8080`.

2. **Stop the application:**

   To stop the running containers, use:

   ```bash
   docker-compose down
   ```

### Running All Services

If you want to run all services, you can visit [Deployment Master Services](https://github.com/segment3d-app/deployment-master)

### API Documentation

To access the API documentation, visit the Swagger documentation at `http://localhost:8080/swagger/index.html` after starting the server.

## License

This project is licensed under the [MIT License](LICENSE).
