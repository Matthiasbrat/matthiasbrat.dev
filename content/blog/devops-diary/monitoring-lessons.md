---
title: "Monitoring Lessons from Production Incidents"
description: "What I wish I knew about observability before my first on-call"
date: 2025-02-08
---

Good monitoring is the difference between "we detected an issue" and "customers are complaining."

## The Four Golden Signals

From Google's SRE book:

1. **Latency** - Time to serve requests
2. **Traffic** - Demand on your system
3. **Errors** - Rate of failed requests
4. **Saturation** - How "full" your system is

> [!TIP]
> If you only track four things, track these.

## Alert Fatigue is Real

Bad alerts:

```yaml
- alert: HighCPU
  expr: cpu_usage > 80
  # Fires constantly, gets ignored
```

Better alerts:

```yaml
- alert: HighCPU
  expr: cpu_usage > 90
  for: 15m
  labels:
    severity: warning
```

> [!IMPORTANT]
> Every alert should be actionable. If you can't do anything about it, don't alert.

## Log Levels Matter

| Level | When to Use |
|-------|-------------|
| ERROR | Something failed, needs attention |
| WARN | Degraded but working |
| INFO | Normal operations |
| DEBUG | Development only |

## Structured Logging

```json
{
    "timestamp": "2025-02-08T10:30:00Z",
    "level": "error",
    "message": "Payment failed",
    "user_id": "123",
    "amount": 99.99,
    "error_code": "CARD_DECLINED"
}
```

> [!NOTE]
> Structured logs are searchable. `Payment failed for user 123` is not.

## Dashboards

Every service should have:

- Request rate
- Error rate
- Latency percentiles (p50, p95, p99)
- Resource usage
