---
name: Bug report
about: Create a reproducible bug report.
title: "[BUG] - "
labels: bug
assignees: UrielCohen456

---

## Checklist

<!-- Do NOT open an issue until you have: --> 

* [ ] Double-checked my configuration.
* [ ] Tested using the latest version.

## Summary

What happened/what you expected to happen?

What version are you running?

## Diagnostics

Paste the smallest workflow that reproduces the bug. We must be able to run the workflow.

```yaml

```

```bash
# Logs from the workflow controller:
kubectl logs -n argo deploy/workflow-controller | grep ${workflow} 

# The workflow's pods that are problematic:
kubectl get pod -o yaml -l workflows.argoproj.io/workflow=${workflow},workflow.argoproj.io/phase!=Succeeded

# Logs from in your workflow's side container (the plugin logs):
kubectl logs -c sidecar -l workflows.argoproj.io/workflow=${workflow},workflow.argoproj.io/phase!=Succeeded
```
