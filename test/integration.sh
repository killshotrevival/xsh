#! /bin/bash

set -e

XSH_BINARY=$1
CONFIG_PATH=$2

if [[ -z $CONFIG_PATH ]]
then
  CONFIG_PATH="$PWD/test"
  echo "No Config path provided, taking $CONFIG_PATH as default"
fi

if [[ $CONFIG_PATH == "/" ]]
then
  echo "Dangerous config path provided $CONFIG_PATH"
  exit 1
fi

mkdir -p $CONFIG_PATH

# Removing resources if any
rm -f $CONFIG_PATH/xsh.db*

echo "Using XSH binar : $XSH_BINARY"
echo "Using config directory: $CONFIG_PATH"

# Setting the config directory path for XSH
export XSH_CONFIG_PATH=$CONFIG_PATH

$XSH_BINARY init

echo "Creating new identity"
cat <<EOF > mock_idfile
-----BEGIN OPENSSH PRIVATE KEY-----
This is a random string, XSH never reads the private file it just check the path provided holds a valid file.
asdasdasdasdasdasdasdasdaasdasdasdasdasdasdasdasasdasdasdasdasdasdasdasasdasdasdasdasdasdasdasasdasdasdasdasdasdasdas
-----END OPENSSH PRIVATE KEY-----
EOF

# Adding identity to database
$XSH_BINARY put identity id1 "$PWD/mock_idfile"

# Fetching identity id
$XSH_BINARY get identity -o json -f identity.json
IDENTITY_ID=$(cat identity.json | jq -r '.[] | select(.name == "id1") .id')
echo "New identity created successfully with id $IDENTITY_ID"

echo "Creating new region"
# Adding region to database
$XSH_BINARY put region "us-east-1"

# Fetching region id
$XSH_BINARY get region -o json -f region.json
REGION_ID=$(cat region.json | jq -r '.[0].id')
echo "New region created successfully with id: $REGION_ID"

echo "Creating new jumphost"
# Creating a jumphost
cat <<EOF > host.json
{
  "name": "example-test",
  "address": "exmaple.test.com",
  "port": 202,
  "user": "test",
  "region_id": "$REGION_ID",
  "identity_id": "$IDENTITY_ID",
  "jumphost_id": null
}
EOF

$XSH_BINARY put host -f host.json

#Fetch Host ID
$XSH_BINARY get host -o json -f host.json
JUMPHOST_ID=$(cat host.json | jq -r '.[0].id')
echo "New jumphost created successfully with id: $JUMPHOST_ID"

echo "Creating new normal host"
# Creating a host using jumphost Id from above
cat <<EOF > host.json
[{
  "name": "example-test-2",
  "address": "exmaple.test-2.com",
  "port": 202,
  "user": "test",
  "region_id": "$REGION_ID",
  "identity_id": "$IDENTITY_ID",
  "jumphost_id": "$JUMPHOST_ID"
}]
EOF

$XSH_BINARY put host -f host.json

# Checking connection command
CONNECT_COMMAND=$($XSH_BINARY connect example-test-2 -p)
echo "New normal host created successfully with connect command"

