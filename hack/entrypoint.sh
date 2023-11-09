#!/bin/bash

echo "Expecting: <artifact-input> <artifact-output> <command>..."
echo "Full provided set of arguments are $@"

# First get the args we need for locations for the artifact
artifactPipe="${1}"
shift 

artifactInput="${1}"
shift

artifactOutput="${1}"
shift

# The command is the remainder of the script $@
echo "Command is $@"
echo "Pipe to is ${artifactPipe}"
echo "Artifact input is ${artifactInput}"
echo "Artifact output is ${artifactOutput}"

# Wait for the sidecar to finish, indicated by the file indicator we wait for
wget -q https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs

# Wait for the indicator from the sidecar that artifact is ready
goshare-wait-fs -p /mnt/oras/oras-operator-init.txt
		
# We expect to be in the working directory needed for the container
# The artifact inputs can either be extracted here, or elsewhere
if [[ "${artifactInput}" == "NA" ]]; then
    cp -R /mnt/oras/inputs/ .
else
    cp -R /mnt/oras/inputs/ ${artifactInput}
fi

# Run the original command
if [[ "${artifactPipe}" == "NA" ]]; then
$@
else
$@ > ${artifactPipe}
fi

# indicate we are done
mkdir -p /mnt/oras/outputs

# Output requires the directory to be defined, otherwise we assume none
if [[ "${artifactOutput}" != "NA" ]]; then
    cp -R ${artifactOutput} /mnt/oras/outputs/
fi

# We are done!
touch /mnt/oras/oras-operator-done.txt

# The script (container) should exit here, and the pod will finish
# when the sidecar is done.