# Kubernetes Manifests

Here reside kubernetes manifests to use k3supdater in your cluster. 
Since `k3supdater` is not web service/microservice but rather a job in itself, it is defined here as  a cronjob that runs `daily`.

## Github Access Token Secret

You may notice that inside the `k3supdater-cron` cronjob, it reads a secret called `github-access-token`, but we don't define it.
That's because you'll have to define it and create it in your cluster before creating the cronjob, like so (make sure to replace the actual token value by yours):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: github-access-token
  namespace: k3supdater
type: Opaque
data:
  GITHUB_ACCESS_TOKEN: <YOUR_ACCESS_TOKEN_HERE>
```
