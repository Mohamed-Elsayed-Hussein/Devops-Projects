# Kubernetes Deployment – 3-Tier App (Step-by-Step Documentation, Updated)

This document provides a **chronological step-by-step guide** for deploying the 3-tier application on Kubernetes, reflecting **all actual changes** including Docker image pushes, Deployment adjustments, and storage configuration.

---

## 1. Namespace and Context

* **Create namespace for the project:**

```bash
kubectl create namespace 3-tier-app
```

* **Switch current context to the new namespace:**

```bash
kubectl config set-context --current --namespace=3-tier-app
kubectl config view --minify | grep namespace
```

> From this point, all commands are executed in the `3-tier-app` namespace.

---

## 2. Secrets

### Database Password Secret

* **Create secret for database password:**

```bash
kubectl create secret generic db-password --from-file=backend/db-password --dry-run=client -o yaml > K8S/backend-secret.yaml
kubectl apply -f K8S/backend-secret.yaml
```

### Database Environment Variables Secret

* **Create secret for database environment variables:**

```bash
kubectl create secret generic database-sec-vars --from-env-file=database/env --dry-run=client -o yaml > K8S/db-secrets.yaml
kubectl apply -f K8S/db-secrets.yaml
```

---

## 3. Persistent Storage for Database

* **Define PersistentVolume (PV) and PersistentVolumeClaim (PVC) for MySQL:**

```yaml
# Example: PV and PVC YAML added to database deployment
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mysql-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /mnt/data/mysql

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

* **Attach PVC in the MySQL Deployment YAML** after dry-run:

```yaml
volumes:
  - name: mysql-storage
    persistentVolumeClaim:
      claimName: mysql-pvc
```

---

## 4. Deployments

> **Note:** Before applying the deployments, Docker images were **pushed to Docker Hub** from the first project and reused here.

### 4.1 Database Deployment (MySQL)

* **Create deployment YAML and customize volumes & strategy after dry-run:**

```bash
kubectl create deployment db --image=mysql:8 --dry-run=client -o yaml > K8S/database_deployment.yaml
```

* **Manually edit YAML** to:

  * Add PV/PVC volume mount
  * Set `strategy` to `RollingUpdate` with `maxSurge` and `maxUnavailable` if needed

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 25%
    maxUnavailable: 25%
```

* **Apply deployment:**

```bash
kubectl apply -f K8S/database_deployment.yaml
```

### 4.2 Backend Deployment (Go Application)

* **Create deployment YAML for backend using pushed Docker image:**

```bash
kubectl create deployment backend-go --image=mohamedelsayed22/go-backend-senior-projec1:v1.0 --port=8000 --dry-run=client -o yaml > K8S/backend_deployment.yaml
```

* **Edit YAML** to customize:

  * Volume mounts (if any)
  * Deployment strategy (RollingUpdate)

* **Apply deployment:**

```bash
kubectl apply -f K8S/backend_deployment.yaml
```

### 4.3 Proxy Deployment (Nginx)

* **Create deployment YAML for proxy using pushed Docker image:**

```bash
kubectl create deployment proxy --image=mohamedelsayed22/nginx:v1.0 --port=443 --dry-run=client -o yaml > K8S/proxy_deployment.yaml
```

* **Edit YAML** to adjust volumes or strategy if required
* **Apply deployment:**

```bash
kubectl apply -f K8S/proxy_deployment.yaml
```

---

## 5. Services

### 5.1 Backend Service (ClusterIP)

```bash
kubectl expose deployment backend-go --type=ClusterIP --port=8000 --target-port=8000 --dry-run=client -o yaml > K8S/backend_service.yaml
kubectl apply -f K8S/backend_service.yaml
```

### 5.2 Database Service (ClusterIP)

```bash
kubectl expose deployment db --port=3306 --dry-run=client -o yaml > K8S/db-service.yaml
kubectl apply -f K8S/db-service.yaml
```

### 5.3 Proxy Service (NodePort)

```bash
kubectl expose deployment proxy --type=NodePort --port=443 --dry-run=client -o yaml > K8S/proxy_nodeport.yaml
kubectl apply -f K8S/proxy_nodeport.yaml
```

* **Access proxy externally**:

```bash
minikube ip  # e.g., 192.168.49.2
curl -k https://192.168.49.2:31087
```

> Use `-k` to bypass self-signed certificate verification.

---

## 6. Verification

* **Check all resources in namespace:**

```bash
kubectl get all -n 3-tier-app
kubectl get secrets -n 3-tier-app
kubectl get pv,pvc -n 3-tier-app
```

* **Inspect individual pods:**

```bash
kubectl describe pod <pod-name> -n 3-tier-app
kubectl logs <pod-name> -n 3-tier-app
```

