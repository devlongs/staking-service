# Staking Service (Section 1)

This repository contains the implementation for Section 1 of the Stakeway Technical Assessment. The service is a lightweight backend API built in Go that manages staking operations using a SQLite database. It demonstrates best practices such as dependency injection, the repository pattern, proper logging, input validation, graceful shutdown, unit testing, and containerization.

## Features

- **API Endpoints**
  - **GET `/v1/health`**: Returns the service status.
  - **POST `/v1/stake`**: Accepts a JSON payload with a wallet address and amount to simulate staking.
  - **GET `/v1/rewards/{wallet_address}`**: Calculates and returns rewards (5% of the staked amount) for a given wallet.
- **SQLite Integration**: Uses a SQLite database to store staking transactions.
- **Input Validation**: Validates Ethereum addresses and ensures non-negative staking amounts.
- **Structured Logging**: Uses the standard library’s logging package with a custom prefix.
- **Graceful Shutdown**: Listens for OS signals and gracefully shuts down the HTTP server.
- **Unit Tests**: Basic unit tests are provided for each endpoint.
- **Containerization**: Dockerfile and Docker Compose configuration for local development.


## Prerequisites

- **Go 1.20** or higher
- **Docker** (for containerization)
- **Docker Compose** (for multi-container orchestration)

## Setup and Configuration

1. **Clone the Repository**
 ```bash
   git clone https://github.com/devlongs/staking-service.git
   cd staking-service
```
2. **Environment Variables**
 ```dotenv
   ADDR=:8080
   DB_PATH=stakeway.db
```
These variables configure the HTTP server address and SQLite database file

## Running Locally
1. **Install Dependencies**
Make sure Go modules are up-to-date:
 ```bash
   go mod tidy
```

2. **Run the backend**
 ```bash
   go run ./cmd/api/main.go
```

2. **Test Endpoints with Curl**
   - Health Check:
   ```bash
     curl http://localhost:8080/v1/health
   ```
   Expected response:
   ```json
   {"alive": true}
   ```

    - Stake Endpoint:
   ```bash
     curl -X POST http://localhost:8080/v1/stake \
     -H "Content-Type: application/json" \
     -d '{"wallet_address": "0x1234567890abcdef1234567890abcdef12345678", "amount": 100}'
   ```
   Expected response:
   ```json
      {"message": "Staking successful"}
   ```

    - Rewards Endpoint::
   ```bash
   curl http://localhost:8080/v1/rewards/0x1234567890abcdef1234567890abcdef12345678
   ```
   Expected response (assuming 100 staked and 5% reward):
   ```json
      {"wallet_address": "0x1234567890abcdef1234567890abcdef12345678", "rewards": 5}
   ```


## Using Docker and Docker Compose
Docker
1. Build the Docker Image
```bash
sudo docker build -t stakeway-backend .
```

1. Run the Docker Container
```bash
sudo docker run -p 8080:8080 stakeway-backend
```

Docker Compose
1. Start the service
```bash
docker compose up --build
```

1. Stop the Service
```bash
docker compose down
```

## Running Unit Tests
Unit tests are provided in the test/v1/ directory. To run all tests:
```bash
go test ./test/v1/...
```
This command runs tests for health, stake, and rewards endpoints.

### Logging and Graceful Shutdown
The service logs are prefixed with "stakeway: " and include timestamps. The application listens for OS termination signals like SIGTERM and gracefully shuts down the HTTP server.