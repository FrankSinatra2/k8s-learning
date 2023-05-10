
# Purpose
Gain more experience with k8s by learning how deployments communicate within the same k8s namespace.

# Requirements
- redis sesson storage
- go backend that provides:
  - log in page (username only)
  - home page that displays your login username and the time you last visited


# Takeaways

- Remember that specifying protocol in a connection string actually means something. Sometimes the error is right in front of you, but you are thinking too hard.
- Running commands in a pod is a very useful debug step, particuallarly running `nslookup`
- Ways to communicate with pods in the same namespace:
    - direct IP
    - using the name of the service that exposes the pods
    - using a longer version of the above: <service-name>.<namespace>.svc.cluster.local
