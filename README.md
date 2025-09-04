# ssh-check
Check ssh connectivity to nodes within a cluster

This check requires two environment variables in the kube spec used to deploy it.  An ssh key, and the user.

**Example yaml**

```yaml
---
apiVersion: kuberhealthy.github.io/v2
kind: KuberhealthyCheck
metadata:
  name: ssh-check
  namespace: kuberhealthy
spec:
  runInterval: 5m
  timeout: 2m
  extraAnnotations:
    comcast.com/testAnnotation: test.annotation
  extraLabels:
    testLabel: testLabel
  podSpec:
    containers:
      - name: ssh-check
        image: rjacks161/ssh-check:v3.0.0
        imagePullPolicy: IfNotPresent
        env:
          - name: SSH_PRIVATE_KEY
            value: "CHANGE_ME"
          - name: SSH_USERNAME
            value: "CHANGEME"
          - name: SSH_EXCLUDE_LIST
            value: "CHANGEME1 CHANGEME2"
        resources:
          requests:
            cpu: 10m
            memory: 50Mi
```

#### How-to

Apply a `.yaml` file similar to the one shown above with `kubectl apply -f`

#### SSH Exclude List
The exclude list is a space delimited list of node names that the user would like to exclude from the ssh check.
