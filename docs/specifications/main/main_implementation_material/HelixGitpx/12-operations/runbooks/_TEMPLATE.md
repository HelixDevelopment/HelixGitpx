# RB-NNN — <Title>

> **Severity default**: <P1|P2|P3>
> **Owner**: <team/role>
> **Last tested**: <YYYY-MM-DD>

## 1. Detection

Alert: `<AlertName>` — brief one-liner or metric condition.

Supporting signals:
- Metric / log signal 1.
- Metric / log signal 2.

---

## 2. First 2 Minutes

1. Acknowledge the alert.
2. Confirm it's real (check dashboard, query, synthetic probes).
3. Gather context (recent deployments, known issues).

---

## 3. Diagnose

Step-by-step commands, queries, and dashboards to isolate the root cause.

```bash
# Example commands
kubectl get pod -n <namespace>
```

---

## 4. Mitigations

1. Immediate action (stop the bleed).
2. Rollback / failover / scale.
3. Stabilization.

---

## 5. Verify Recovery

- Service responds to health checks.
- Error rate returns to baseline.
- Queues draining.

---

## 6. Post-Incident

- [ ] Incident timeline documented
- [ ] RCA identified
- [ ] Action items filed with owners + due dates
- [ ] Runbook updated with learnings

---

## 7. Related

- Related runbooks: RB-XXX, RB-YYY
- Dashboards: Grafana link
- Alerts: AlertManager rules
