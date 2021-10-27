#!/bin/sh

/usr/bin/example-clientgo -list

/usr/bin/example-clientgo -create

# Give some time to start the pods
sleep 30

/usr/bin/example-clientgo -list

/usr/bin/example-clientgo -update

/usr/bin/example-clientgo -list

/usr/bin/example-clientgo -delete

# Give some time to stop the pods
sleep 30

/usr/bin/example-clientgo -list

# sleep for long so we can check pods logs.
sleep 10000