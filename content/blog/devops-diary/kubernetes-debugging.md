---
title: "Kubernetes Debugging Cheat Sheet"
description: "Commands and techniques for troubleshooting K8s issues"
date: 2025-02-03
---

When things go wrong in Kubernetes, here's how to investigate.

## Pod Issues

```bash
# Get pod status
kubectl get pods -o wide

# Describe pod (shows events)
kubectl describe pod <pod-name>

# View logs
kubectl logs <pod-name>
kubectl logs <pod-name> --previous  # Crashed container

# Exec into container
kubectl exec -it <pod-name> -- /bin/sh
```

> [!TIP]
> Add `-f` to `logs` to follow output in real-time.

## Common Pod States

| State | Meaning | First Check |
|-------|---------|-------------|
| Pending | Not scheduled | Node resources, affinity |
| ImagePullBackOff | Can't pull image | Image name, registry auth |
| CrashLoopBackOff | Keeps crashing | Application logs |
| OOMKilled | Out of memory | Memory limits |

## Network Debugging

```bash
# Check service endpoints
kubectl get endpoints <service-name>

# Test DNS resolution
kubectl run debug --rm -it --image=busybox -- nslookup <service-name>

# Check network policies
kubectl get networkpolicies
```

## Resource Issues

```bash
# Node resource usage
kubectl top nodes

# Pod resource usage
kubectl top pods

# Check resource quotas
kubectl describe resourcequota
```

> [!WARNING]
> `kubectl top` requires metrics-server to be installed.

## Quick Wins

```bash
# See all events (sorted by time)
kubectl get events --sort-by='.lastTimestamp'

# Check if pods can reach each other
kubectl run test --rm -it --image=busybox -- wget -qO- http://<service>:<port>
```
