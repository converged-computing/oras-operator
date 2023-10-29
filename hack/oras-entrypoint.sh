#!/bin/bash

# The command is the remainder of the script $@
pullfrom="${1}"
pushto="${2}"
echo "Artifact URI to retrieve is: ${pullfrom}"
echo "Artifact URI to push to is: ${pushto}"

# Create inputs artifact directory
mkdir -p /mnt/oras/inputs /mnt/oras/outputs
if [[ "${pullfrom}" != "NA" ]]; then
    cd /mnt/oras/inputs
    oras pull ${pushto} --plain-http
    echo "Pulled inputs to /mnt/oras/inputs"
    ls -l
fi

# indicate to application we are ready to run!
touch /mnt/oras/oras-operator-init.txt

# Wait for the application to finish, indicated by the file indicator we wait for
wget -q https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs

# Wait for the indicator from the sidecar that artifact is ready
goshare-wait-fs -p /mnt/oras/oras-operator-done.txt

# If we don't have a place to push, we are done
if [[ "${pushto}" == "NA" ]]; then
    exit 0
fi
	
# Push the contents to the location
cd /mnt/oras/outputs
oras push ${pushto} --plain-http .

# Now we are done and can exit