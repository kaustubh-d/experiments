#!/bin/sh

# Update this to use whatever namespace you want.
APPNAMESPACE="rgb"

/usr/bin/example-clientgo -list -namespace "$APPNAMESPACE"

/usr/bin/example-clientgo -create -namespace "$APPNAMESPACE"

# Give some time to start the pods
sleep 30

/usr/bin/example-clientgo -list -namespace "$APPNAMESPACE"

/usr/bin/example-clientgo -update -namespace "$APPNAMESPACE"

/usr/bin/example-clientgo -list -namespace "$APPNAMESPACE"

/usr/bin/example-clientgo -delete -namespace "$APPNAMESPACE"

# Give some time to stop the pods
sleep 30

/usr/bin/example-clientgo -list -namespace "$APPNAMESPACE"

# sleep for long so we can check pods logs.
sleep 10000