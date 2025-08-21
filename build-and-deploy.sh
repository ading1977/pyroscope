#!/bin/bash

set -e

echo "Building OTLP Debug Tool..."

# Build Docker image
echo "Building Docker image..."
docker build -f Dockerfile.debug -t otlp-debug:latest .

# If using kind cluster, load the image
if command -v kind &> /dev/null; then
    echo "Loading image into kind cluster..."
    kind load docker-image otlp-debug:latest --name=pyroscope-dev || echo "kind cluster not found, skipping..."
fi

# If using minikube, load the image
if command -v minikube &> /dev/null && minikube status &> /dev/null; then
    echo "Loading image into minikube..."
    minikube image load otlp-debug:latest || echo "minikube not running, skipping..."
fi

echo "Deploying to Kubernetes..."
kubectl apply -f k8s-debug-deployment.yaml

echo "Waiting for deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/otlp-debug

echo "Getting service information..."
kubectl get pods -l app=otlp-debug
kubectl get svc otlp-debug-service

echo ""
echo "=== DEPLOYMENT COMPLETE ==="
echo ""
echo "To access the debug service:"
echo "1. Port forward: kubectl port-forward service/otlp-debug-service 4040:4040"
echo "2. Then send OTLP data to: http://localhost:4040/v1/profiles"
echo "3. Check logs with: kubectl logs -l app=otlp-debug -f"
echo ""
echo "To check service endpoint in cluster:"
kubectl get svc otlp-debug-service -o jsonpath='{.spec.clusterIP}' && echo ":4040/v1/profiles"