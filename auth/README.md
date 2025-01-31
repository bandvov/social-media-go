# Auth service
## Build
```bash
docker build -t your-dockerhub-username/auth-service:latest .
```

## Push to dockerhub
```bash
docker push your-dockerhub-username/auth-service:latest
```
## Deploy to kubernetes
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/secret.yaml

```