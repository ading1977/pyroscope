# OTLP Debug Tool

This tool helps debug the OTLP v1.7.0 issue where unsymbolized frames disappear from profiles.

## Quick Setup

1. **Build and Deploy**:
```bash
./build-and-deploy.sh
```

2. **Access the Service**:
```bash
kubectl port-forward service/otlp-debug-service 4040:4040
```

3. **Send Your OTLP Data**:
```bash
# Redirect your OTLP client to: http://localhost:4040/v1/profiles
# Or use curl to test with existing OTLP data
curl -X POST http://localhost:4040/v1/profiles \
  -H "Content-Type: application/x-protobuf" \
  --data-binary @your-profile.pb
```

4. **Check Debug Output**:
```bash
kubectl logs -l app=otlp-debug -f
```

## What It Shows

The debug tool will output:
- **HasFunctions status** for each mapping (TRUE/FALSE)
- **Memory ranges** and **filenames** 
- **Build IDs** from attributes
- **Location counts** and symbolization status
- **Summary statistics**

## Debugging Your Issue

### Step 1: Compare OTLP v1.4.0 vs v1.7.0

1. **Deploy this debug tool** in your environment
2. **Temporarily redirect** your OTLP client to send data to the debug tool
3. **Check the logs** to see if `HasFunctions=true` for mappings that should be unsymbolized

### Step 2: Key Things to Look For

Look for patterns like this in the output:
```
Mapping 0: HasFunctions=true, HasFilenames=false, HasLineNumbers=false
  MemoryStart=0x400000, MemoryLimit=0x500000, Filename=/usr/bin/myapp
  BuildID=abc123...
```

**Expected for unsymbolized**: `HasFunctions=false`
**Your issue**: `HasFunctions=true` (incorrectly set)

### Step 3: Root Cause Identification

If you see `HasFunctions=true` for unsymbolized mappings, the issue is in your OTLP client/collector, not Pyroscope. The client is incorrectly setting this flag in OTLP v1.7.0.

## Manual Testing

### Build Locally
```bash
cd /path/to/pyroscope
cp go.mod.debug go.mod
go mod tidy
go run debug_otlp_input.go
```

### Test with Sample Data
```bash
# Create test OTLP data and send it
curl -X POST http://localhost:4040/v1/profiles \
  -H "Content-Type: application/x-protobuf" \
  --data-binary @sample.pb
```

## Configuration for Your Environment

### Change the Service Type
Edit `k8s-debug-deployment.yaml` and change service type:
```yaml
spec:
  type: LoadBalancer  # or NodePort
```

### Update the Ingress
Change the host in `k8s-debug-deployment.yaml`:
```yaml
rules:
- host: your-domain.com
```

### Scale for High Traffic
```yaml
spec:
  replicas: 3  # Increase replicas if needed
```

## Cleanup

```bash
kubectl delete -f k8s-debug-deployment.yaml
docker rmi otlp-debug:latest
```

## Next Steps After Debugging

Once you identify the root cause:

1. **If HasFunctions=true incorrectly**: Update your OTLP client/collector configuration
2. **If HasFunctions=false correctly**: The issue is in Pyroscope's handling (investigate symbolizer logic)
3. **Check OTLP spec differences**: Compare v1.4.0 vs v1.7.0 default behaviors

## Troubleshooting

- **Image not found**: Run `docker build` manually or check your Kubernetes image pull policy
- **Service not accessible**: Verify port-forward and check pod status with `kubectl get pods`
- **No logs**: Ensure your OTLP client is sending data to the correct endpoint

```
I've created a comprehensive debugging solution with:

  Files Created:

  1. debug_otlp_input.go - Main debug program
  2. go.mod.debug - Go module file
  3. Dockerfile.debug - Docker image
  4. k8s-debug-deployment.yaml - Kubernetes deployment + service + ingress
  5. build-and-deploy.sh - Automated build and deployment script
  6. debug-README.md - Complete usage instructions

  How to Use:

  1. Deploy Everything:
  ./build-and-deploy.sh

  2. Access the Debug Service:
  kubectl port-forward service/otlp-debug-service 4040:4040

  3. Redirect Your OTLP Client to send data to http://localhost:4040/v1/profiles
  4. Watch the Debug Output:
  kubectl logs -l app=otlp-debug -f

  What This Will Reveal:

  The debug tool will show you exactly what's in your OTLP data:

  === OTLP DEBUG [2025-01-21 10:30:45] ===
  ResourceProfile 0:
    ScopeProfile 0:
      Profile 0:
        Mapping 0: HasFunctions=true, HasFilenames=false, HasLineNumbers=false
          MemoryStart=0x400000, MemoryLimit=0x500000, Filename=/usr/bin/myapp
          BuildID=abc123...
        Unsymbolized locations: 150/200
  SUMMARY: Total mappings=1, Unsymbolized mappings=0

  Key Diagnostic:

  - If you see HasFunctions=true for unsymbolized mappings → The problem is in your OTLP client/collector
  - If you see HasFunctions=false correctly → The problem is in Pyroscope's processing

  This will definitively identify whether the issue is in the OTLP data itself or in Pyroscope's handling of it.
```
