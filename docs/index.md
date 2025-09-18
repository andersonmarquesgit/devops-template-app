# Documentation for Devops Template App

This application is a template for test and create components using Backstage and standard for pipeline creation using GitHub Actions, Helm, Kubernetes and ArgoCD.
Is possible explore all the flow using CI/CD and configurations needs for this.
Exist a endpoint for test the application:
- `/api/v1/details`

You can access using URL **http://devops-template-app.localtest.me/api/v1/details**

```json
{
	hostname: "devops-template-app-55c794f65b-dg4lj",
	message: "New attribute message for test cicd!!",
	time: "2025-09-17T13:58:31.212528637Z"
}
```
