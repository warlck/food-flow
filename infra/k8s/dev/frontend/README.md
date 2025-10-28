# Frontend Kubernetes Deployment

This directory contains the Kubernetes manifests for deploying the Food Flow Online Hub frontend application to a Kubernetes cluster.

## Directory Structure

```
infra/k8s/
├── base/
│   └── frontend/
│       ├── base-frontend.yaml      # Base deployment and service definitions
│       └── kustomization.yaml      # Base kustomization
└── dev/
    └── frontend/
        ├── kustomization.yaml                    # Dev overlay kustomization
        ├── dev-frontend-patch-deploy.yaml        # Dev deployment patches (resources, replicas)
        └── dev-frontend-patch-service.yaml       # Dev service patches (NodePort)
```

## Resources Defined

### Deployment
- **Name**: `frontend`
- **Namespace**: `sales-system`
- **Image**: `localhost/food-flow/frontend:0.0.1`
- **Replicas**: 1 (dev environment)
- **Container Port**: 8080
- **Strategy**: RollingUpdate (zero downtime deployments)

### Service
- **Name**: `frontend-service`
- **Type**: NodePort (dev environment)
- **Port**: 8080
- **NodePort**: 30080 (accessible from outside the cluster)

### Health Checks
- **Liveness Probe**: `/health` endpoint on port 8080
  - Initial delay: 10s
  - Period: 10s
  - Timeout: 3s
- **Readiness Probe**: `/health` endpoint on port 8080
  - Initial delay: 5s
  - Period: 5s
  - Timeout: 3s

### Resource Limits (Dev)
- **Requests**: 50m CPU, 32Mi memory
- **Limits**: 100m CPU, 64Mi memory

## Deployment

### Prerequisites
1. Build the frontend Docker image:
   ```bash
   make frontend
   ```

2. Have a running Kind cluster:
   ```bash
   make dev-up
   ```

### Deploy Frontend to Kubernetes

**Option 1: Deploy everything (recommended)**
```bash
# Build, load images, and deploy all services including frontend
make build
make dev-load
make dev-apply
```

**Option 2: Deploy only frontend**
```bash
# Build and load frontend image
make frontend
kind load docker-image localhost/food-flow/frontend:0.0.1 --name food-flow-cluster

# Apply frontend manifests
kustomize build infra/k8s/dev/frontend | kubectl apply -f -

# Wait for frontend to be ready
kubectl wait pods --namespace=sales-system --selector app=frontend --timeout=120s --for=condition=Ready
```

**Option 3: Update existing frontend deployment**
```bash
# Build new image
make frontend

# Load into cluster
kind load docker-image localhost/food-flow/frontend:0.0.1 --name food-flow-cluster

# Restart deployment
kubectl rollout restart deployment frontend --namespace=sales-system

# Or use make command
make dev-restart
```

## Accessing the Frontend

### From within the cluster
```bash
# Using service name
curl http://frontend-service.sales-system.svc.cluster.local:8080/health
```

### From your local machine (NodePort)
Since the service is exposed as NodePort on port 30080:

```bash
# Access via localhost (with Kind's port mapping)
curl http://localhost:30080/health

# Open in browser
open http://localhost:30080
```

### Port Forward (alternative method)
```bash
# Forward local port 8080 to pod port 8080
kubectl port-forward -n sales-system deployment/frontend 8080:8080

# Access at
open http://localhost:8080
```

## Monitoring & Troubleshooting

### Check deployment status
```bash
# View all resources
kubectl get all -n sales-system -l app=frontend

# Describe deployment
make dev-describe-frontend
# Or:
kubectl describe deployment frontend -n sales-system
```

### View logs
```bash
# Follow logs
make dev-logs-frontend

# Or manually:
kubectl logs -n sales-system -l app=frontend -f --tail=100
```

### Check pod status
```bash
# Get pod details
kubectl get pods -n sales-system -l app=frontend

# Describe pod
kubectl describe pod -n sales-system -l app=frontend

# Get events
kubectl get events -n sales-system --sort-by='.lastTimestamp' | grep frontend
```

### Test health endpoint
```bash
# Get pod name
POD_NAME=$(kubectl get pods -n sales-system -l app=frontend -o jsonpath='{.items[0].metadata.name}')

# Test health endpoint
kubectl exec -n sales-system $POD_NAME -- wget -qO- http://localhost:8080/health
```

### Debug container
```bash
# Get shell access to frontend container
kubectl exec -it -n sales-system deployment/frontend -- sh

# Inside container:
ls -la /usr/share/nginx/html/     # Check built files
cat /etc/nginx/conf.d/default.conf # Check nginx config
ps aux                             # Check running processes
```

## Scaling

### Manual scaling
```bash
# Scale to 3 replicas
kubectl scale deployment frontend -n sales-system --replicas=3

# Verify
kubectl get pods -n sales-system -l app=frontend
```

### Auto-scaling (HPA)
```yaml
# Create HPA (example)
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: frontend-hpa
  namespace: sales-system
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: frontend
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## Configuration Updates

### Update API URL (if needed)
The frontend is built at Docker image creation time. To change the API URL:

1. Rebuild image with VITE_API_URL:
   ```bash
   docker build \
     -f infra/docker/dockerfile.frontend \
     -t localhost/food-flow/frontend:0.0.1 \
     --build-arg VITE_API_URL=http://sales-service.sales-system.svc.cluster.local:3000 \
     .
   ```

2. Load and restart:
   ```bash
   kind load docker-image localhost/food-flow/frontend:0.0.1 --name food-flow-cluster
   kubectl rollout restart deployment frontend -n sales-system
   ```

### Update Resource Limits
Edit `infra/k8s/dev/frontend/dev-frontend-patch-deploy.yaml`:
```yaml
resources:
  requests:
    cpu: "100m"      # Increase as needed
    memory: "64Mi"
  limits:
    cpu: "200m"
    memory: "128Mi"
```

Then reapply:
```bash
kustomize build infra/k8s/dev/frontend | kubectl apply -f -
```

## Production Considerations

For production deployment:

1. **Change Service Type**: Use `LoadBalancer` or `Ingress` instead of NodePort
2. **Increase Resources**: Allocate more CPU/memory based on traffic
3. **Add Ingress**: Configure proper ingress with SSL/TLS
4. **Enable Monitoring**: Add Prometheus annotations for metrics
5. **Configure Autoscaling**: Set up HPA based on metrics
6. **Use CDN**: Consider serving static assets via CDN
7. **Multiple Replicas**: Run at least 3 replicas for high availability
8. **Pod Disruption Budget**: Prevent simultaneous pod evictions

Example Ingress:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: frontend-ingress
  namespace: sales-system
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - app.example.com
    secretName: frontend-tls
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 8080
```

## Cleanup

### Remove frontend deployment
```bash
kustomize build infra/k8s/dev/frontend | kubectl delete -f -
```

### Remove entire cluster
```bash
make dev-down
```

## Integration with Backend Services

The frontend can connect to backend services using Kubernetes service DNS:

- **Sales API**: `http://sales-service.sales-system.svc.cluster.local:3000`
- **Auth API**: `http://auth-service.sales-system.svc.cluster.local:6000`

These URLs should be configured in the frontend build-time environment variables.
