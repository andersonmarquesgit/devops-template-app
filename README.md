# devops-template-app

Sources:
Argocd: 
- "helm upgrade --install argocd argo/argo-cd -n argocd -f values-argo.yaml --wait" Or " helm upgrade --install argocd argo/argo-cd -n argocd --create-namespace -f values-argo.yaml"
- https://argo-cd.readthedocs.io/en/stable/operator-manual/installation/
- "kubectl get secrets -n argocd argocd-initial-admin-secret -o yaml"
- "echo "SUA SECRET" | base64 -d"
- "argocd login argocd.localtest.me --insecure --grpc-web --username [ARGOCD USER] --password [ARGOCD PASSWORD]"
Actions Runner: 
- https://github.com/actions/actions-runner-controller/blob/master/docs/quickstart.md