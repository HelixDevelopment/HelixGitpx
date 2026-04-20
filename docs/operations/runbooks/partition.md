# Runbook — Network partition

**Alert:** `HelixGitpxServiceUnreachable` or sudden error-rate spike isolated to one namespace.

## Immediate checks (1 min)

- `kubectl -n helixgitpx get pods -o wide` — note node distribution.
- Check Cilium Hubble: `hubble observe --from-label app=<svc> --not --to-label app=<svc>` for blackholes.
- Check NetworkPolicy recently applied: `kubectl get netpol -A --sort-by=.metadata.creationTimestamp | tail -5`.

## Diagnosis

Most partitions are caused by (a) bad NetworkPolicy (deny-all landed before allow), (b) Istio Ambient sidecar issues, (c) underlay node NIC saturation.

## Remediation

1. Roll back the most recent NetworkPolicy/Kyverno change if timing fits.
2. If Istio: `istioctl analyze` + restart ztunnel on affected nodes.
3. If underlay: drain the hot node.

## Escalation

Partition persists > 5 min → page on-call SRE + network lead.
