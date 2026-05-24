---
title: "Getting Started"
weight: 10
---

## 🚀 Installation

### 1. The Easy Way (Docker Compose)
The fastest way to test Searchbase locally is using the native DuckDuckGo demo Compose example:

```bash
git clone https://github.com/coolapso/searchbase.git
cd searchbase
docker compose -f examples/compose/docker-compose.native.yml up
```
The Gateway will be available at `http://localhost:8080`.

{{< hint warning >}}
The Compose files are demonstration examples only and are not intended for production deployments. See [Self Hosting](/docs/self-hosting/) for backend-specific examples and deployment notes.
{{< /hint >}}

### 2. The Production Way (Kubernetes)
Searchbase is designed to run on Kubernetes. You can scale `search-gateway`, `crawl-worker`, and search backends independently based on traffic, rendering load, and search workload.

```bash
kubectl apply -f k8s/
```

For backend choices, environment variables, and production notes, see [Self Hosting](/docs/self-hosting/).
