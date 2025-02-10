# Staking Service

This repository contains the implementation for the Stakeway Technical Assessment.

- **Section 1: Backend API** – A lightweight Go backend that manages staking operations using SQLite. It provides endpoints for health, staking, and rewards. This section features dependency injection, structured logging, graceful shutdown, unit testing, and containerization.  
  **(Section 1 is available on the master branch.)**
  
- **Section 2: Deployment and Monitoring** – Kubernetes manifests to deploy the backend service (with at least 2 replicas) plus instructions for setting up monitoring with Prometheus and Grafana. The service exposes a `/metrics` endpoint for monitoring (showing total requests per endpoint and average response time) along with sample alert rules.  
  **(Section 2 is implemented in the branch `feature/section2`.)**

---

## Section 1: Backend API

### Features

- **API Endpoints**
  - **GET `/v1/health`** – Returns `{"alive": true}`.
  - **POST `/v1/stake`** – Accepts JSON with `wallet_address` (string) and `amount` (float). Validates input and simulates staking by storing the transaction in SQLite.
  - **GET `/v1/rewards/{wallet_address}`** – Calculates rewards as 5% of the staked amount for the given wallet.
- **SQLite Integration** – Stores staking transactions.
- **Input Validation** – Validates Ethereum addresses and ensures non-negative amounts.
- **Structured Logging** – Logs using the standard library with a custom prefix (`"stakeway: "`).
- **Graceful Shutdown** – Listens for OS signals and shuts down the server gracefully.
- **Unit Tests** – Basic tests are provided for each endpoint.
- **Containerization** – Dockerfile and Docker Compose configuration are provided.


### Setup and Running Locally

1. **Run the Backend**
```bash
   go run ./cmd/api/main.go
```

2. **Test Endpoints**
- Health:
```bash
   curl http://localhost:8080/v1/health
```
Expected: {"alive": true}

- Stake:
```bash
   ccurl -X POST http://localhost:8080/v1/stake \
     -H "Content-Type: application/json" \
     -d '{"wallet_address": "0x1234567890abcdef1234567890abcdef12345678", "amount": 100}'
```
Expected: {"message": "Staking successful"}

- Rewards:
```bash
   curl http://localhost:8080/v1/rewards/0x1234567890abcdef1234567890abcdef12345678
```
Expected (after staking): {"wallet_address": "0x1234567890abcdef1234567890abcdef12345678", "rewards": 5}


### Using Docker and Docker Compose
**Docker:**
1. **Build the image**
```bash
sudo docker build -t stakeway-backend .
```

2. **Run the Container**
```bash
sudo docker run -p 8080:8080 stakeway-backend
```

**Docker Compose:**
1. **Start the Service**
```bash
docker compose up --build
```

2. **Stop the Service**
```bash
docker compose down
```

### Running Unit Tests
Run all tests from the project root:
```bash
go test ./test/v1/...
```



---

### Section 2: Deployment and Monitoring
#### Kubernetes Deployment
All Kubernetes manifests are located in the k8s/ directory.

Here is the link to the public container registery: https://hub.docker.com/r/devlongs/stakeway-backend

1. **Apply the ConfigMap**
```bash
kubectl apply -f k8s/configmap.yaml
```

2. **Apply the Deployment**
```bash
kubectl apply -f k8s/deployment.yaml
```

3. **Apply the Service**
```bash
kubectl apply -f k8s/service.yaml
```

4. **Verify**
```bash
kubectl get pods
kubectl get svc
```

To access the service via Minikube:
```bash
minikube service stakeway-backend
```

Or use port-forwarding:
```bash
kubectl port-forward svc/stakeway-backend 8080:80
```

Test the health endpoint:
```bash
curl http://localhost:8080/v1/health
```

### Monitoring Setup
#### Deploy Prometheus

1. **Create a Namespace**
```bash
kubectl create namespace monitoring
```

2. **Add the Prometheus Helm Repository and Update**
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

3. **Install Prometheus**
```bash
helm install prometheus prometheus-community/prometheus --namespace monitoring
```

4. **Access Prometheus UI**
```bash
kubectl port-forward svc/prometheus-server -n monitoring 9090:80
```
Then visit http://localhost:9090/targets in your browser to verify your /metrics endpoint is being scraped.

#### Deploy Grafana

1. **Add the Grafana Helm Repository and Update**
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

2. **Install Grafana**
```bash
helm install grafana grafana/grafana --namespace monitoring
```

3. **Access Grafana UI**
Port-forward the Grafana service:
```bash
kubectl port-forward svc/grafana -n monitoring 3000:80
```

1. **Add the Grafana Helm Repository and Update**
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```
Then open http://localhost:3000 in your browser.
Default credentials are usually admin/admin.


### Configure Grafana
- **Add Prometheus as a Data Source:**
  - In Grafana, go to Configuration → Data Sources → Add data source.
  - Select Prometheus and set the URL to http://localhost:9090 (if port-forwarding) or http://prometheus-server.monitoring.svc.cluster.local.
  - Click Save & Test.
- **Create a Dashboard:**
  - Add panels using queries such as:
    - Total Requests per Endpoint:
    ```bash
    sum by (path)(http_requests_total)
    ```

     - Average Response Time (over 5 minutes):
    ```bash
    avg by (path)(rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m]))
    ```
