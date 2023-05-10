
# relevent k8s page
- https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough

# k8s resources used
- Namespace
- StatefulSet
- HoriontalPodAutoScaler

## Namespace
- Just used to keep this sub-project organized in the cluster

## Service
- required to access the fib-app stateful set
- Learned:
    - more about how a service selects the pods to expose. In this case, `spec.selector.run: fib-app` 
      will select pods with `template.metadata.labls.run: fib-app`

## StatefulSet
- For this sub-project, a stateful set was used as I had not used one prior.
- Learned:
    - `resources.limits` defines the maximum allowance for a given resource
    - `resources.requests` defines the expected allowace for a give resource


## HoriontalPodAutoScaler
- Main point of this sub-project
- Learned:
    - How to select the resource to autoscale 
    - Simple autoscaling by monitoring the cpu usage of the specified target

# Helpful Commands

## Running a load generator
`kubectl run -i --tty <name> --rm --image=busybox --restart=Never -- /bin/sh -c "<script>"`

## Monitoring HoriontalPodAutoScaler
`kubectl get hpa <name> --watch`


# Extra Notes
- It can take a few minutes for pods to be automatically scaled down
