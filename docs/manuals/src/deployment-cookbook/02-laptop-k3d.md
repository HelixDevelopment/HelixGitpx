## 2. Laptop install via k3d + Argo CD

The fastest way to see HelixGitpx running end-to-end is a single-node
[k3d](https://k3d.io) cluster on your laptop. Bring-up takes ~10 minutes.

### 2.1 Prerequisites

| Tool | Min version | Install |
|------|-------------|---------|
| k3d | 5.7.x | `brew install k3d` (macOS), `curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh \| bash` (Linux) |
| kubectl | 1.30 | `brew install kubectl` |
| helm | 3.16 | `brew install helm` |
| podman *or* docker | latest | OS package |

### 2.2 Create the cluster

```bash
k3d cluster create helixgitpx \
    --servers 1 --agents 2 \
    --port "8080:80@loadbalancer" \
    --port "8443:443@loadbalancer" \
    --k3s-arg "--disable=traefik@server:*"
```

Verify:

```bash
kubectl get nodes
# NAME                           STATUS   ROLES                  AGE   VERSION
# k3d-helixgitpx-server-0        Ready    control-plane,master   …
# k3d-helixgitpx-agent-0         Ready    <none>                 …
# k3d-helixgitpx-agent-1         Ready    <none>                 …
```

### 2.3 Install Argo CD

```bash
kubectl create namespace argocd
helm repo add argo https://argoproj.github.io/argo-helm
helm install argocd argo/argo-cd -n argocd \
    -f impl/helixgitpx-platform/helm/argocd/values.yaml
kubectl -n argocd wait deploy/argocd-server --for=condition=Available --timeout=5m
```

### 2.4 Apply the app-of-apps

```bash
kubectl apply -f impl/helixgitpx-platform/argocd/app-of-apps.yaml
```

Watch the rollout:

```bash
watch 'kubectl -n argocd get applications -o wide'
```

All 53 applications should reach `Synced + Healthy` within 15 minutes on
a typical laptop. CNPG's first-boot takes longest.

### 2.5 Seed an admin

```bash
kubectl -n keycloak exec deploy/keycloak -- /opt/keycloak/bin/kc.sh \
    create-user --realm helixgitpx --username you --password changeme
```

### 2.6 First push

```bash
# Port-forward git-ingress.
kubectl -n helixgitpx port-forward svc/git-ingress 8080:8080 &

# Create a local repo against it.
mkdir -p /tmp/smoke && cd /tmp/smoke
git init && echo hi > README.md
git add . && git commit -m init
git remote add hx http://localhost:8080/me/smoke.git
git push -u hx main
```

### 2.7 Tear down

```bash
k3d cluster delete helixgitpx
```

All data is volatile by design — k3d clusters are disposable.

---
