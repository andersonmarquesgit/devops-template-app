# Documentation for Devops Template App

This application is a template for test and create components using Backstage and standard for pipeline creation using GitHub Actions, Helm, Kubernetes and ArgoCD.
Is possible explore all the flow using CI/CD and configurations needs for this.
Exist a endpoint for test the application **http://devops-template-app.localtest.me/api/v1/details**.

```json
{
	hostname: "devops-template-app-55c794f65b-dg4lj",
	message: "New attribute message for test cicd!!",
	time: "2025-09-17T13:58:31.212528637Z"
}
```

Este repositório contém um guia passo a passo para configurar um ambiente de desenvolvimento utilizando **Minikube**, **ArgoCD**, **Actions Runner Controller (ARC)** e **Backstage**.

# 📦 Pré-requisitos

-   Docker
-   kubectl
-   Helm
-   Minikube

## 1. Minikube

Inicie o cluster:
```bash
minikube start --driver=docker
```
Se o túnel do Minikube ficar preso em outra execução:

```bash
# Liste túneis ativos
sudo pgrep -alf "minikube.*tunnel"

# Finalize processos
sudo pkill -TERM -f "minikube.*tunnel" || true
sleep 1
sudo pkill -KILL -f "minikube.*tunnel" || true

# Limpe rotas/ips criados pelo túnel anterior
sudo minikube tunnel --cleanup || true

# Verifique se não sobrou túnel
sudo pgrep -alf "minikube.*tunnel" || echo "Sem tunnel ativo"
```
Garanta que o **Ingress** esteja configurado como LoadBalancer:
```bash
kubectl -n ingress-nginx get svc ingress-nginx-controller -o wide
```

## ## 2. ArgoCD

Instale o ArgoCD:
```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm upgrade --install argocd argo/argo-cd -n argocd --create-namespace -f values-argo.yaml --wait
```

Recupere a senha inicial:
```bash
kubectl get secrets -n argocd argocd-initial-admin-secret -o yaml
echo "<ARGOCD_SECRET_BASE64>" | base64 -d
```
Faça login:
```bash
argocd login argocd.localtest.me --insecure --grpc-web \
  --username <ARGOCD_USER> \
  --password <ARGOCD_PASSWORD>
```
## 3. Actions Runner Controller (ARC)

```bash
# Cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.18.2/cert-manager.yaml

# Helm repo
helm repo add actions-runner-controller https://actions-runner-controller.github.io/actions-runner-controller

# Deploy do controller
helm upgrade --install actions-runner-controller \
  actions-runner-controller/actions-runner-controller \
  -n actions-runner-system \
  --create-namespace \
  --set=authSecret.create=true \
  --set=authSecret.github_token="<GITHUB_PAT>" \
  --wait
```
🔄 O self-hosted dura cerca de 1h. Caso perca o runner, atualize o token:
```bash
kubectl -n actions-runner-system create secret generic controller-manager \
  --from-literal=github_token=<NOVO_GITHUB_PAT> \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl -n actions-runner-system rollout restart deploy/actions-runner-controller
```
Runner Deployment
```yaml
kubectl -n actions-runner-system apply -f - <<'YAML'
apiVersion: actions.summerwind.dev/v1alpha1
kind: RunnerDeployment
metadata:
  name: org-runners
spec:
  replicas: 1
  template:
    spec:
      organization: <ORG_NAME>
      group: Default
      labels:
        - self-hosted
        - Linux
        - X64
YAML
```

Valide o status:
```bash
kubectl -n actions-runner-system get runnerdeployments,runnerreplicasets,runners
kubectl -n actions-runner-system get pods -l actions.github.com/runner-deployment=org-runners
kubectl -n actions-runner-system logs -l actions.github.com/runner-deployment=org-runners -c runner --tail=200 -f
```

# Backstage

## Setup inicial

```bash
docker pull node:22-bookworm-slim

docker run --rm \
  -e AUTH_GITHUB_CLIENT_ID=<GITHUB_CLIENT_ID> \
  -e AUTH_GITHUB_CLIENT_SECRET=<GITHUB_CLIENT_SECRET> \
  -p 3000:3000 -p 7007:7007 \
  -v /path/to/backstage-app:/app \
  -w /app \
  node:22-bookworm-slim bash
```
Crie o app:
```bash
npx @backstage/create-app@latest
cd my-backstage-app
yarn start
```
Acesse: [http://localhost:3000](http://localhost:3000)

## Backstage Tech Docs

Dependências:

```bash
apt-get update && \
apt-get install -y python3 python3-pip python3-venv && \
rm -rf /var/lib/apt/lists/*

python3 -m venv /opt/venv
export PATH="/opt/venv/bin:$PATH"
pip3 install mkdocs-techdocs-core
```
Crie documentação:

```bash
mkdir docs
echo "# Documentação inicial" > docs/index.md
```

## Software Templates


Definições adotadas para novos projetos via Backstage:
-   **Dockerfile** para build da aplicação
-   **GitHub** como repositório de código
-   **GitHub Actions** para pipelines de CI/CD
-   **Docker Hub** para imagens
-   **Helm + ArgoCD + Kubernetes** para deploy automatizado


👉 Agora você tem um passo a passo unificado para configurar **Minikube + ArgoCD + Runners + Backstage** no mesmo ambiente.

## Docs
- https://argo-cd.readthedocs.io/en/stable/operator-manual/installation/
<<<<<<< Updated upstream
- https://github.com/actions/actions-runner-controller/blob/master/docs/quickstart.md
- https://backstage.io/docs/next/getting-started/
- https://backstage.io/docs/auth/
- https://backstage.io/docs/auth/github/provider
- https://github.com/backstage/backstage
- https://github.com/backstage/backstage/blob/master/packages/catalog-model/examples/acme/team-a-group.yaml
- https://backstage.io/docs/features/techdocs/
=======
- "kubectl get secrets -n argocd argocd-initial-admin-secret -o yaml"
- "echo "SUA SECRET" | base64 -d"
- "argocd login argocd.localtest.me --insecure --grpc-web --username [ARGOCD USER] --password [ARGOCD PASSWORD]"
Actions Runner: 
- https://github.com/actions/actions-runner-controller/blob/master/docs/quickstart.md

Backstage (Spotify)
- https://backstage.io/docs/next/getting-started/
- docker pull node:22-bookworm-slim
- docker run --rm -e AUTH_GITHUB_CLIENT_ID=Ov23lieP29lt6UcbyUnM -e AUTH_GITHUB_CLIENT_SECRET=58bb78966ecd60bd4657e13132956eb54f2e7ace -p 3000:3000 -ti -p 7007:7007 -v /Users/anderson/GolandProjects/backstage-app:/app -w /app node:22-bookworm-slim bash
- npx @backstage/create-app@latest
- root@69e84d903d95:/app/backstage# apt-get install nano
- ls app-config-app
  more app-config-app
- add listen: 0.0.0.0 a app-config usando nano (apt-get install -y nano)
- cd my-backstage-app # your app name
  yarn start
- Access http://localhost:3000
- https://backstage.io/docs/auth/
- https://backstage.io/docs/auth/github/provider

- Usei OAuth Apps do GitHub adicionando o trecho abaixo em app-config.local.yaml

- Parte Backend

```yaml
auth:
  environment: development
  providers:
    github:
      development:
        clientId: ${AUTH_GITHUB_CLIENT_ID}
        clientSecret: ${AUTH_GITHUB_CLIENT_SECRET}
        ## uncomment if using GitHub Enterprise
        # enterpriseInstanceUrl: ${AUTH_GITHUB_ENTERPRISE_INSTANCE_URL}
        ## uncomment to set lifespan of user session
        # sessionDuration: { hours: 24 } # supports `ms` library format (e.g. '24h', '2 days'), ISO duration, "human duration" as used in code
        signIn:
          resolvers:
            # See https://backstage.io/docs/auth/github/provider#resolvers for more resolvers
            - resolver: usernameMatchingUserEntityName
```

e no diretório do /app/backstage execute a instalação do plugin para adicionar o provider
- yarn --cwd packages/backend add @backstage/plugin-auth-backend-module-github-provider
- Add this lines in packages/backend/src/index.ts
```typescript
backend.add(import('@backstage/plugin-auth-backend'));
backend.add(import('@backstage/plugin-auth-backend-module-github-provider'));
```
- Parte Frontend https://backstage.io/docs/auth/#sign-in-configuration



- A última coisa a fazer é configurar o Resolver
- https://github.com/backstage/backstage
- https://github.com/backstage/backstage/blob/master/packages/catalog-model/examples/acme/team-a-group.yaml
>>>>>>> Stashed changes
