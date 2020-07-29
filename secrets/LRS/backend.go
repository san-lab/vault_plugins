package mock

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const generatePath string = "genKeys"
const showPubKeysPath string = "showPubKeys"
const signMsgPath string = "signMsg"

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{
		store: make(map[string][]byte),
	}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(mockHelp),
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathGenerate(b),
				pathPubKeys(b),
				pathSignMsg(b),
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



const mockHelp = `
The Mock backend is a dummy secrets backend that stores kv pairs in a map.
`
