# Sigstore GCP KMS Plugin

This is an example repo that implements a plugin to Sigstore for GCP's KMS.
See: https://github.com/sigstore/sigstore/pull/1901

It basically imports the existing GCP KMS code, but you will build this as a plugin binary.

Example, with the pending https://github.com/sigstore/cosign/pull/3954 and https://github.com/sigstore/sigstore/pull/1901:

```shell
# Compile the plugin program, adding it to PATH

PATH=$PATH:$(pwd)/bin
go -C sigstore-kms-mygcpkms build -o $(pwd)/bin

# Prepare our KeyRef for use in cosign

PROJECT="my-project"
LOCATION="global"
KEYRING="my-test-keyring"
KEY="my-test-key-2"
KEY_VERSION="1"
KEY_RESOURCE_ID="mygcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION"

# Create a key with cosign and our KMS plugin

go -C cosign/cmd/cosign run ./ generate-key-pair --kms "$KEY_RESOURCE_ID"

# Confirm the correct key is in your GCP account

gcloud kms keys versions get-public-key "$KEY_VERSION" \
    --key "$KEY" \
    --keyring "$KEYRING" \
    --location "$LOCATION"

```
