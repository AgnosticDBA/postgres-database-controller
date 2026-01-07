# PostgresDatabase Controller

Kubernetes operator that manages PostgresDatabase custom resources by automatically generating and maintaining PerconaPGCluster instances.

## Overview

This controller watches for PostgresDatabase resources and translates the simple developer-facing API into full PerconaPGCluster configurations with:

- Automatic HA configuration with pod anti-affinity
- Sensible defaults for image versions and storage classes
- Built-in monitoring with PMM integration
- Automated backup setup with pgBackRest
- Resource limits and requests

## Features

- **Automatic Translation**: Converts 8-line PostgresDatabase manifests to 50+ line PerconaPGCluster
- **Best Practices**: Applies platform defaults (affinity, storage classes, security)
- **Status Management**: Updates PostgresDatabase status based on cluster state
- **Reconciliation**: Handles updates and maintains desired state
- **Extensible**: Configurable platform defaults

## Installation

### Prerequisites

- Kubernetes cluster (v1.20+)
- Percona PostgreSQL Operator installed
- PostgresDatabase CRD deployed

### Deploy the Controller

1. **Build the Docker image:**

```bash
docker build -t agnosticdba/postgres-database-controller:latest .
```

2. **Deploy with Helm or Kubernetes manifests:**

```bash
kubectl apply -f deploy/
```

## Configuration

### Platform Defaults

The controller uses sensible defaults that can be customized:

```go
PlatformConfig{
    DefaultStorageClass:    "standard",
    DefaultImageRegistry:   "docker.io/percona",
    DefaultCRVersion:       "2.8.2",
    DefaultPGBouncerImage:  "percona-pgbouncer:1.25.0-1",
    DefaultPGBackRestImage: "percona-pgbackrest:2.57.0-1",
    DefaultPMMImage:        "pmm-client:3.5.0",
    DefaultPMMHost:         "prometheus.monitoring",
}
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DEFAULT_STORAGE_CLASS` | "standard" | Kubernetes storage class |
| `DEFAULT_IMAGE_REGISTRY` | "docker.io/percona" | Container registry |
| `DEFAULT_PMM_HOST` | "prometheus.monitoring" | PMM server host |

## Development

### Local Development

1. **Run locally:**

```bash
go run main.go --leader-elect=false
```

2. **Install dependencies:**

```bash
go mod tidy
go mod vendor
```

3. **Run tests:**

```bash
go test ./...
```

### Project Structure

```
postgres-database-controller/
├── api/v1/                          # Custom resource definitions
│   ├── postgresdatabase_types.go    # API types
│   ├── groupversion_info.go          # API version info
│   └── register.go                   # Scheme registration
├── internal/controller/              # Controller logic
│   └── postgresdatabase_controller.go # Reconciliation logic
├── main.go                          # Manager entrypoint
├── go.mod                           # Go module definition
├── go.sum                           # Dependency checksums
├── Dockerfile                       # Container build
└── README.md                        # This file
```

## How It Works

### Reconciliation Loop

1. **Watch for Changes**: Monitors PostgresDatabase resources
2. **Generate Cluster**: Creates PerconaPGCluster from simple spec
3. **Apply Best Practices**: Adds affinity, security, and monitoring
4. **Update Status**: Reports cluster state back to PostgresDatabase
5. **Handle Updates**: Detects spec changes and updates clusters

### Example Transformation

**Input (PostgresDatabase):**
```yaml
apiVersion: databases.mycompany.com/v1
kind: PostgresDatabase
metadata:
  name: billing-db
  namespace: production
spec:
  version: 17
  replicas: 3
  storage: 100Gi
  backup: true
  monitoring: true
```

**Output (PerconaPGCluster):**
```yaml
apiVersion: pgv2.percona.com/v2
kind: PerconaPGCluster
metadata:
  name: billing-db
  namespace: production
  labels:
    created-by: postgres-database-controller
spec:
  crVersion: 2.8.2
  image: docker.io/percona/percona-distribution-postgresql:17.7-2
  postgresVersion: 17
  instances:
  - name: instance1
    replicas: 3
    affinity:
      podAntiAffinity:
        preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 1
          podAffinityTerm:
            labelSelector:
              matchLabels:
                postgres-operator.crunchydata.com/data: postgres
            topologyKey: kubernetes.io/hostname
    dataVolumeClaimSpec:
      accessModes: [ReadWriteOnce]
      resources:
        requests:
          storage: 100Gi
      storageClassName: standard
  proxy:
    pgBouncer:
      replicas: 3
      image: docker.io/percona/percona-pgbouncer:1.25.0-1
      affinity: # ... auto-configured
  backups:
    pgbackrest:
      image: docker.io/percona/percona-pgbackrest:2.57.0-1
      repos: [{name: "repo1"}]
  pmm:
    enabled: true
    image: docker.io/percona/pmm-client:3.5.0
    serverHost: prometheus.monitoring
```

## RBAC

The controller requires the following permissions:

```yaml
# PostgresDatabase management
- groups: ["databases.mycompany.com"]
  resources: ["postgresdatabases", "postgresdatabases/status", "postgresdatabases/finalizers"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# PerconaPGCluster management
- groups: ["pgv2.percona.com"]
  resources: ["perconapgclusters", "perconapgclusters/status"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

## Monitoring

The controller exposes metrics on `:8080/metrics`:

- `controller_runtime_reconcile_total` - Total reconciliations
- `controller_runtime_reconcile_errors_total` - Reconciliation errors
- `workqueue_depth` - Queue depth
- `workqueue_latency` - Queue processing latency

## Troubleshooting

### Common Issues

1. **Controller fails to start**: Check RBAC permissions
2. **Clusters not created**: Verify PerconaPGCluster CRD exists
3. **Status not updating**: Check controller logs for reconciliation errors

### Debug Mode

```bash
export LOG_LEVEL=debug
go run main.go
```

## License

Apache License 2.0

## Related Repositories

- **[postgres-database](https://github.com/AgnosticDBA/postgres-database)** - CRD definition and examples
- **[percona-postgresql-operator](https://github.com/percona/percona-postgresql-operator)** - Underlying PostgreSQL operator