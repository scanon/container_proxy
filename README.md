Container Proxy Service

This service runs outside the container and listens on a socket for request.
A client utility should be accessible in the container and any commands should be symlinked to the client.


**Usage**

Note this feature isn't fully integrated into Shifter yet so some manual steps are required.

Here is an example usage.  This would be inside the Slurm Job prior to starting the Container do the following...

```
# Set the Proxy Filename.  Make sure this is uniqueue
export PROXY_SOCKET=/tmp/${USER}.sock
# Start the server and background it
/global/common/shared/das/container_proxy/server.py &
CPID=$!
# Start the container process
shifter --image=centos:8 bash
> # Add proxied commands to the path
> PATH=/global/common/shared/das/container_proxy/:$PATH
> # Now you can call some slurm tools
> squeue
> exit
kill $CPID
```

