package ccpsecrets

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	ccp "github.com/liviusnl/go-ccp"
)

const configPath string = "config"
const generatePath string = "genKey"
const showAddressPath string = "showAddr"
const signTxPath string = "signTx"

type backend struct {
	*framework.Backend

	lock   sync.Mutex
	client *ccp.Client
	store map[string][]byte
}

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := newBackend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend implements the CCP Secrets Engine.
func newBackend() *backend {
	var b = &backend{}

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				generatePath,
			},
		},

		Paths: framework.PathAppend(
			[]*framework.Path{

				pathGenerate(b),
				pathAddress(b),
				pathSignTx(b),
			},
		),

		InitializeFunc: b.initialize,
		Invalidate:     b.invalidate,

		Clean: b.cleanup,
	}

	return b
}

// initialize the plugin.
func (b *backend) initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.Client(ctx, req.Storage)

	return nil
}

// invalidate resets the plugin. This is called when a key is updated via
// replication.
func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case generatePath:
		b.ResetClient(nil)
	}
}

func (b *backend) cleanup(ctx context.Context) {
	b.ResetClient(nil)
}

const backendHelp = `
The CyberArk Credentials Provider (CCP) allows users to access secrets stored
in the CyberArk Enterprise Password Vault (EPV) without the need to directly 
access the CCP Web Service.

After mounting this secrets engine, you can configure the secrets eingine
using the "config/" endpoints.
`
