package ethPlugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const generatePath string = "genKey"
const showAddressPath string = "showAddr"
const preSignPath string = "preSign"
const signTxPath string = "signTx"

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{
		store: make(map[string][]byte),
	}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(ethereumPluginHelp),
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathGenerate(b),
				pathAddress(b),
				pathPreSign(b),
				pathSignTx(b),
			},
		),
	}

	//b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	b.Backend.Setup(ctx, conf)

	return b, nil
}

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend

	store map[string][]byte
}



const ethereumPluginHelp = `
The ethereumPlugin backend is a plugin that allows you to input a ethereum transaction and returns it signed.
`
