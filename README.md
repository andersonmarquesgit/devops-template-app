
# Documentation for Devops Template App

This application is a template for test and create components using Backstage and standard for pipeline creation using GitHub Actions, Helm, Kubernetes and ArgoCD.  
Is possible explore all the flow using CI/CD and configurations needs for this.  
Exist a endpoint for test the application **http://devops-template-app.localtest.me/api/v1/details**.

```json  
{  
 hostname: "devops-template-app-55c794f65b-dg4lj", message: "New attribute message for test cicd!!", time: "2025-09-17T13:58:31.212528637Z"}  
```  

This repository contains a step-by-step guide to configure a development environment using **Minikube**, **ArgoCD**, **Actions Runner Controller (ARC)** and **Backstage**.

# 📦 Pré-requisitos

- Docker
- kubectl
- Helm
- Minikube

## 1. Minikube

Start cluster::
```bash  
minikube start --driver=docker```  
If tunnel gets stuck::  
```bash  
  
sudo pgrep -alf "minikube.*tunnel"  
  
sudo pkill -TERM -f "minikube.*tunnel" || true  
sleep 1  
sudo pkill -KILL -f "minikube.*tunnel" || true  
  
sudo minikube tunnel --cleanup || true  
  
sudo pgrep -alf "minikube.*tunnel" || echo "Sem tunnel ativo"  
```  

Check Ingress as LoadBalancer:

```bash  
kubectl -n ingress-nginx get svc ingress-nginx-controller -o wide
```  

## 2. ArgoCD

Install ArgoCD:
```bash  
helm repo add argo https://argoproj.github.io/argo-helmhelm upgrade --install argocd argo/argo-cd -n argocd --create-namespace -f values-argo.yaml --wait
```  

Recover initial password:
```bash  
kubectl get secrets -n argocd argocd-initial-admin-secret -o yamlecho "<ARGOCD_SECRET_BASE64>" | base64 -d
```  

Login:

```bash  
argocd login argocd.localtest.me --insecure --grpc-web \ --username <ARGOCD_USER> \ --password <ARGOCD_PASSWORD>
```  

## 3. Actions Runner Controller (ARC)

```bash  
# Cert-manager  
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.18.2/cert-manager.yaml  
  
# Helm repo  
helm repo add actions-runner-controller https://actions-runner-controller.github.io/actions-runner-controller  
  
# Deploy do controller  
helm upgrade --install actions-runner-controller \  
 actions-runner-controller/actions-runner-controller \ -n actions-runner-system \ --create-namespace \ --set=authSecret.create=true \ --set=authSecret.github_token="<GITHUB_PAT>" \ --wait  
```  

🔄 Runners expire after ~1h. To refresh:

```bash  
kubectl -n actions-runner-system create secret generic controller-manager \ --from-literal=github_token=<NOVO_GITHUB_PAT> \  
 --dry-run=client -o yaml | kubectl apply -f -  
kubectl -n actions-runner-system rollout restart deploy/actions-runner-controller
```  

Runner Deployment example:

```yaml  
kubectl -n actions-runner-system apply -f - <<'YAML'  
apiVersion: actions.summerwind.dev/v1alpha1  
kind: RunnerDeployment  
metadata:  
 name: org-runnersspec:  
 replicas: 1 template: spec: organization: <ORG_NAME> group: Default labels: - self-hosted - Linux - X64YAML  
```  

Validate:
```bash  
kubectl -n actions-runner-system get runnerdeployments,runnerreplicasets,runnerskubectl -n actions-runner-system get pods -l actions.github.com/runner-deployment=org-runnerskubectl -n actions-runner-system logs -l actions.github.com/runner-deployment=org-runners -c runner --tail=200 -f
```  

# Backstage

## Setup inicial

```bash  
docker pull node:22-bookworm-slim  
docker run --rm \ -e AUTH_GITHUB_CLIENT_ID=<GITHUB_CLIENT_ID> \ -e AUTH_GITHUB_CLIENT_SECRET=<GITHUB_CLIENT_SECRET> \ -p 3000:3000 -p 7007:7007 \ -v /path/to/backstage-app:/app \ -w /app \ node:22-bookworm-slim bash
```  

Create app:

```bash  
npx @backstage/create-app@latestcd my-backstage-appyarn start
```  
Access: [http://localhost:3000](http://localhost:3000)

## Backstage Tech Docs

Dependencies:

```bash  
apt-get update && \apt-get install -y python3 python3-pip python3-venv && \rm -rf /var/lib/apt/lists/*  
python3 -m venv /opt/venvexport PATH="/opt/venv/bin:$PATH"pip3 install mkdocs-techdocs-core
```  
Create docs:

```bash  
mkdir docsecho "# Documentação inicial" > docs/index.md
```  

## Software Templates

For new projects:
-   **Dockerfile** for build
-   **GitHub** for repository
-   **GitHub Actions** for CI/CD
-   **Docker Hub** for images
-   **Helm + ArgoCD + Kubernetes** for deploy

## Defining CI/CD integration into backstage

1.  Build Backstage image with TechDocs:

```dockerfile  
FROM node:22-bookworm-slim  
  
RUN apt-get update && apt-get install -y python3 python3-pip python3-venv && rm -rf /var/lib/apt/lists/*  
  
RUN python3 -m venv /opt/venv  
ENV PATH="/opt/venv/bin:$PATH"  
RUN pip3 install mkdocs-techdocs-core  
```  

2.  Add GitHub Actions plugin:
 ```bash
 yarn --cwd packages/app add @backstage-community/plugin-github-actions 
 ```

3.  Edit `packages/app/src/components/catalog/EntityPage.tsx` to include the plugin.  
    (Use `nano` if needed: `apt-get update && apt-get install nano -y`)
- https://backstage.io/plugins
- https://github.com/backstage/community-plugins/tree/main/workspaces/github-actions/plugins/github-actions

## Backstage in Production Mode
### With Docker
- Create network:
 ```bash
docker network create backstage 
 ```

- Run Postgres (https://hub.docker.com/_/postgres):
 ```bash
docker run -d --name psql \
  -e POSTGRES_DB=backstagedb \
  -e POSTGRES_USER=*USER_DB* \
  -e POSTGRES_PASSWORD=*PASSWORD_DB* \
  -e PGDATA=/var/lib/postgresql/data/pgdata \
  -v /tmp/psql:/var/lib/postgresql/data \
  --network backstage \
  postgres:16
```

-   Create `app-config.production.yaml`
-   Build backend (https://backstage.io/docs/deployment/):
```bash
yarn build:backend --config ../../app-config.yaml --config ../../app-config.production.yaml
```
- Build production image:
```bash
docker build -t backstage_production .
```

- Build production image:
```bash
docker run --network backstage --name backstage -d \
  -e AUTH_GITHUB_CLIENT_ID=<CLIENT_ID> \
  -e AUTH_GITHUB_CLIENT_SECRET=<CLIENT_SECRET> \
  -e GITHUB_TOKEN=<TOKEN> \
  -e NODE_OPTIONS=--no-node-snapshot \
  -p 3000:3000 -p 7007:7007 backstage_production
```
Access: [http://localhost:7007](http://localhost:7007)

### With Kubernetes
- Create Helm repo:
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
```

- Install Postgres:
```bash
helm install psql bitnami/postgresql --version 16 -n backstage --create-namespace -f values-postgres.yaml
```

- Build & push Backstage image:
```bash
docker build -t backstage:1.0 .
docker tag backstage:1.0 andersonmarquesdocker/backstage:1.0
docker push andersonmarquesdocker/backstage:1.0
```

- Apply manifests:
```bash
kubectl apply -f backstage-k8s.yaml -n backstage
kubectl get pods -n backstage
kubectl get ingress -n backstage
```

Update GitHub OAuth App redirect URL to:  
**[https://backstage.localtest.me/](https://backstage.localtest.me/)**


## 📚 Useful Docs
-   [ArgoCD Installation](https://argo-cd.readthedocs.io/en/stable/operator-manual/installation/)
-   [Actions Runner Controller Quickstart](https://github.com/actions/actions-runner-controller/blob/master/docs/quickstart.md)
-   [Backstage Getting Started](https://backstage.io/docs/next/getting-started/)
-   [Backstage Auth](https://backstage.io/docs/auth/)
-   [Backstage Auth GitHub](https://backstage.io/docs/auth/github/provider)
-   [Backstage GitHub Repo](https://github.com/backstage/backstage)
-   [Catalog Example](https://github.com/backstage/backstage/blob/master/packages/catalog-model/examples/acme/team-a-group.yaml)
-   [TechDocs](https://backstage.io/docs/features/techdocs/)