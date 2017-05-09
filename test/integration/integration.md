create a pod without deployment and /badpodswo should show it.

```
apiVersion: v1
kind: Pod
metadata:
  name: single-pod
spec:
  containers:
    - name: nginx
      image: nginx:1.7.9
      ports:
        - containerPort: 80
```

create a bad deployment:

```
kubectl run bad-deployment --namespace=dummy-test --image=gcr.io/google_containers/echoserver:1.4
```

create a bad ingress:

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: bad-ingress
  labels:
    name: bad-ingress
spec:
  rules:
  - host: pizza.example.com
    http:
      paths:
      - backend:
          serviceName: some-service
          servicePort: 80
        path: /
```

create a bad job

```
apiVersion: batch/v1
kind: Job
metadata:
  name: bad-job
spec:
  template:
    metadata:
      name: bad-job
    spec:
      containers:
      - name: bad-job
        image: bash
        command: ["sleep 300"]
      restartPolicy: Never
```

create a bad cronjob

```
apiVersion: batch/v2alpha1
kind: CronJob
metadata:
  name: bad-cronjob
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```