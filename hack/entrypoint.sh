#!/bin/bash

# First get the args we need for locations for the artifact
artifactInput="${1}"
shift

artifactOutput="${2}"
shift

# The command is the remainder of the script $@
echo "Command is $@"

# Wait for the sidecar to finish, indicated by the file indicator we wait for
wget -q https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs

# Wait for the indicator from the sidecar that artifact is ready
goshare-wait-fs -p /mnt/oras/oras-operator-init.txt
		
# We expect to be in the working directory needed for the container
# The artifact inputs can either be extracted here, or elsewhere
if [[ "${artifactInput}" == "NA" ]]; then
    cp -R /mnt/oras/inputs/* .
else
    cp -R /mnt/oras/inputs/* ${artifactInput}
fi

# Run the original command
$@

# indicate we are done
mkdir -p /mnt/oras/outputs

# Same with output - either copy from working directory, or as indicated
if [[ "${artifactInput}" == "NA" ]]; then
    cp -R . /mnt/oras/outputs/
else
    cp -R ${artifactOutput} /mnt/oras/outputs/
fi

# We are done!
touch /mnt/oras/oras-operator-done.txt

# The script (container) should exit here, and the pod will finish
# when the sidecar is done.