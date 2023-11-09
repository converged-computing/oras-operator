#!/bin/bash

echo "Expecting: <pull-from> <push-to>"
echo "Full provided set of arguments are $@"

# The command is the remainder of the script $@
pushto="${1}"
shift

# We will get an unknown number of inputs
pullfrom="${1}"
shift

echo "Artifact URI to push to is: ${pushto}"
echo "Artifact URI to pull from is: ${pullfrom}"

# Create inputs artifact directory
mkdir -p /mnt/oras/inputs /mnt/oras/outputs

while [ "${pullfrom}" != "NA" ]; do
    echo "Artifact URI to retrieve is: ${pullfrom}"
    cd /mnt/oras/inputs
    oras pull ${pullfrom} --plain-http
    echo "Pulled ${pullfrom} to /mnt/oras/inputs"
    pullfrom="${1}"
    shift
    if [[ "${pullfrom}" == "" ]]; then
        echo "Hit last artifact to pull."
        pullfrom="NA"
    fi
    ls -l
done

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