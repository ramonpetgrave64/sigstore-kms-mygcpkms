package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sigstore/sigstore/pkg/signature/kms"

	"github.com/sigstore/sigstore/pkg/signature/kms/cliplugin/common"
	"github.com/sigstore/sigstore/pkg/signature/kms/cliplugin/handler"
	"github.com/sigstore/sigstore/pkg/signature/kms/gcp"
)

const expectedProtocolVersion = "1"

func newSignerVerifier(initOptions *common.InitOptions) (kms.SignerVerifier, error) {
	ctx := context.TODO()
	// cliplugin will strip the part up to [plugin name]://[key ref],
	// but the existing GCP code expects a specific regex, so we reconstruct.
	fullKeyResourceID := "gcpkms://" + initOptions.KeyResourceID
	return gcp.LoadSignerVerifier(ctx, fullKeyResourceID)
}

func main() {
	// we log to stderr, not stdout. stdout is reserved for the plugin return value.
	spew.Fdump(os.Stderr, os.Args)
	if protocolVersion := os.Args[1]; protocolVersion != expectedProtocolVersion {
		err := fmt.Errorf("expected protocol version: %s, got %s", expectedProtocolVersion, protocolVersion)
		handler.WriteErrorResponse(os.Stdout, err)
		panic(err)
	}

	pluginArgs, err := handler.GetPluginArgs(os.Args)
	if err != nil {
		handler.WriteErrorResponse(os.Stdout, err)
		panic(err)
	}
	spew.Fdump(os.Stderr, pluginArgs)

	signerVerifier, err := newSignerVerifier(pluginArgs.InitOptions)
	if err != nil {
		handler.WriteErrorResponse(os.Stdout, err)
		panic(err)
	}

	resp, err := handler.Dispatch(os.Stdout, os.Stdin, pluginArgs, signerVerifier)
	if err != nil {
		// Dispatch() will have already called WriteResponse() with the error.
		panic(err)
	}
	spew.Fdump(os.Stderr, resp)
}
