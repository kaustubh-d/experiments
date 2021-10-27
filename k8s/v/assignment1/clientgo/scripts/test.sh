#!/bin/sh -x

# Update this to use whatever namespace you want.
APPNAMESPACE="default"

go run main.go -kubeconfig ~/.kube/config  -list -namespace "$APPNAMESPACE"

go run main.go -kubeconfig ~/.kube/config  -create -namespace "$APPNAMESPACE"

# Give some time to start the pods
kubectl get pods -n $APPNAMESPACE -w -l app=rgb

go run main.go -kubeconfig ~/.kube/config  -list -namespace "$APPNAMESPACE"

go run main.go -kubeconfig ~/.kube/config  -update -namespace "$APPNAMESPACE"

go run main.go -kubeconfig ~/.kube/config  -list -namespace "$APPNAMESPACE"

go run main.go -kubeconfig ~/.kube/config  -delete -namespace "$APPNAMESPACE"

# Give some time to start the pods
kubectl get pods -n $APPNAMESPACE -w -l app=rgb

go run main.go -kubeconfig ~/.kube/config  -list -namespace "$APPNAMESPACE"

