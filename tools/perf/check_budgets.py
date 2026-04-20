#!/usr/bin/env python3
"""Enforce perf budgets against k6 JSON output.

Exit 1 if any metric breaches a budget defined in budgets.json.
"""
import json
import sys
from pathlib import Path


def load_k6(path: Path) -> dict:
    metrics = {}
    with path.open() as f:
        for line in f:
            try:
                rec = json.loads(line)
            except json.JSONDecodeError:
                continue
            if rec.get("type") != "Point":
                continue
            m = rec.get("metric")
            metrics.setdefault(m, []).append(rec["data"]["value"])
    return metrics


def scenario_from_filename(p: Path) -> str:
    stem = p.stem.replace("k6-out-", "")
    return {
        "baseline": "api_baseline",
        "git": "git_push_pull",
        "ws": "websocket_fanout",
        "ai": "ai_chat",
    }.get(stem, stem)


def percentile(values, pct):
    if not values:
        return 0
    values = sorted(values)
    idx = int(len(values) * pct / 100)
    return values[min(idx, len(values) - 1)]


def main(argv):
    if len(argv) < 3:
        print("usage: check_budgets.py <k6-out-*.json...> <budgets.json>", file=sys.stderr)
        return 2
    budgets = json.loads(Path(argv[-1]).read_text())
    failures = []
    for jf in argv[1:-1]:
        path = Path(jf)
        scn = scenario_from_filename(path)
        budget = budgets.get(scn)
        if not budget:
            continue
        metrics = load_k6(path)
        dur = metrics.get("http_req_duration", [])
        p95 = percentile(dur, 95)
        p99 = percentile(dur, 99)
        if "p95_ms" in budget and p95 > budget["p95_ms"]:
            failures.append(f"{scn}: p95 {p95:.1f}ms > {budget['p95_ms']}ms")
        if "p99_ms" in budget and p99 > budget["p99_ms"]:
            failures.append(f"{scn}: p99 {p99:.1f}ms > {budget['p99_ms']}ms")
    if failures:
        for f in failures:
            print("FAIL:", f)
        return 1
    print("All perf budgets met.")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
