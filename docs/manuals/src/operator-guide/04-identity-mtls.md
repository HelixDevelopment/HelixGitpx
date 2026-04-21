## 4. Identity and mTLS

HelixGitpx uses OIDC for user identity (Keycloak) and SPIFFE/SPIRE for
workload identity (east-west mTLS). This chapter covers both.

### 4.1 Keycloak

Keycloak runs in-cluster by default. The `helixgitpx` realm is imported
from `impl/helixgitpx-platform/helm/keycloak/realms/helixgitpx.json`
during initial sync. Customize it in a GitOps overlay, never via the
admin UI — changes made through the UI are not reflected in git.

**Federate to an external IdP** (Okta, Azure AD, Google Workspace) via
SAML or OIDC federation in the realm config. See the operator-guide
chapter 5 for worked examples.

**Clients included at GA:**

- `helixgitpx-web` — the Angular app (public, PKCE).
- `helixgitpx-cli` — the CLI (public, device code).
- `helixgitpx-mobile` — Android/iOS (public, PKCE).
- `helixgitpx-desktop` — desktop apps (public, PKCE).
- `helixgitpx-api-m2m` — machine-to-machine (confidential, client-credentials).

### 4.2 SPIRE

SPIRE issues SVIDs to every HelixGitpx workload. The trust domain is
`helixgitpx.io`. Attestation is Kubernetes-native (ServiceAccount name
+ Pod label + node identity).

**Bootstrap sequence:**

1. `spire-server` deploys in namespace `spire`.
2. `spire-agent` DaemonSet runs one pod per node.
3. On first launch, every HelixGitpx Deployment's init container
   executes `spire-agent api fetch -write /var/run/spire/` to populate
   a workload SVID.

**Rotation** is automatic (default 1 hour TTL).

### 4.3 mTLS between services

Istio Ambient picks up the SVID and terminates mTLS in the ztunnel
layer. No sidecars. Services speak plain HTTP/gRPC to localhost; the
ztunnel adds TLS at the node boundary.

**Verify mTLS is on** for a given service:

```bash
istioctl authn tls-check $(kubectl get pod -l app=orgteam -o jsonpath='{.items[0].metadata.name}').helixgitpx
```

Expected output includes `STATUS: OK` and `AUTHN POLICY: default.istio-system`.

### 4.4 External (north-south) TLS

cert-manager + Let's Encrypt issue TLS for every `Ingress`. The
production issuer is `letsencrypt-prod`; staging uses
`letsencrypt-staging`. Rotation is automatic at 60 days.

Private CA option: for air-gapped installs, swap the ClusterIssuer to
the Vault-backed `helixgitpx-private-ca`.

### 4.5 Break-glass

If OIDC is broken and you can't log in:

1. `kubectl -n keycloak exec deploy/keycloak -- kcadm.sh login --user admin --password <break-glass>`
2. Re-issue the admin secret from Vault.
3. Recreate the realm-import config if it was overwritten.

A runbook for this exact scenario is at
[`docs/operations/runbooks/cert-expiry.md`](../../../operations/runbooks/cert-expiry.md).

---
