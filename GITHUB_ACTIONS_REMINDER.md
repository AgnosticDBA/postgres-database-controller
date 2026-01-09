# GitHub Actions for PostgreSQL Microservices Cluster

Yes! You have excellent GitHub Actions examples for microservices PostgreSQL clusters. Here are the ones I found:

## 🚀 **Available Examples:**

### 1. **Microservice CI/CD** (`common/github-actions/deploy-microservice.yml`)
**Purpose**: Test and deploy Node.js/Go microservices
**Features**:
- Multi-branch support (main, develop)
- Environment-specific deployments (dev, staging, prod)
- Docker Buildx with multi-platform support
- PostgreSQL service integration with health checks
- Automated testing with coverage reports
- Container registry push to GitHub Container Registry

**Workflow**:
```yaml
name: Test and Deploy Microservice
on:
  push:
    branches: [ main, develop ]
  workflow_dispatch:
    inputs:
      environment:
        description: 'Target environment'
        required: true
        default: 'dev'
        type: choice
        options:
          - dev
          - staging
          - prod

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: ${{ github.event.inputs.environment || 'dev' }}
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
        
      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
          
      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: ./examples/microservice
          push: true
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-
            type=raw,value=latest,enable={{is_default_branch}}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha,mode=max
          
      - name: Deploy to Kubernetes
        run: |
          # Update image in deployment
          sed -i 's|image: .*|image: ${{ needs.build.outputs.image }}|g' examples/k8s/microservice/deployment.yaml
          
          # Apply deployment
          kubectl apply -f examples/k8s/microservice/
```

---

### 2. **Database Migration** (`common/github-actions/database-migration.yml`)
**Purpose**: Run database migrations as part of deployment
**Features**:
- Migration job with rollback support
- PostgreSQL connection validation
- Integration with deployment workflows

### 3. **PostgreSQL Cluster Deployment** (`common/github-actions/deploy-postgres.yml`)
**Purpose**: Deploy complete PostgreSQL clusters for microservices
**Features**:
- Environment-aware deployments
- Percona PostgreSQL Operator integration
- Cluster health verification
- Service configuration and port-forwarding
- Database connectivity testing

---

## 🎯 **How to Use for Your Project:**

### **For Microservices:**
```bash
# 1. Add your microservice to the examples/microservice/ directory
# 2. Update the Dockerfile and package.json as needed
# 3. Push to trigger the workflow
git push origin main
# 4. Monitor deployment via GitHub Actions
```

### **For Database Clusters:**
```bash
# 1. Configure your PostgreSQL cluster specs in examples/k8s/cluster/
# 2. Trigger deployment workflow
git push origin main  # or use workflow_dispatch
# 3. Monitor cluster creation via GitHub Actions
```

## 📋 **Integration Points:**

### **CI/CD Pipeline:**
- Microservice builds and deploys to different environments
- Database clusters provisioned automatically
- Health checks and service discovery

### **PostgreSQL DBaaS Integration:**
- Your microservices can connect to PostgreSQL clusters created by these workflows
- Service discovery through Kubernetes DNS
- Configurable database credentials and connection strings

### **Customization:**
- All workflows use environment variables for configuration
- Docker Buildx supports ARM64/AMD64 for your development setup
- Kubernetes manifests follow K8s best practices

## 🔗 **Related Actions:**

These examples show how to build a complete CI/CD pipeline that:
1. **Builds and tests** your applications
2. **Deploys infrastructure** (PostgreSQL clusters)  
3. **Manages configurations** for different environments
4. **Integrates services** with your database infrastructure

Your project has excellent automation examples! You can adapt these patterns for your specific microservices and PostgreSQL setup.