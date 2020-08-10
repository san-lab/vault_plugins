package mock

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const showAddressPathRegExp = showAddressPath + "/(?P<user>[^/]+)$"
func pathAddress(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: showAddressPathRegExp,
		Fields: map[string]*framework.FieldSchema{
			"user": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathAddressRead,
		},

		HelpSynopsis:    objectHelpSyn,
		HelpDescription: objectHelpDesc,
	}
}

// pathObjectRead executes a CCP Object request
func (b *backend) pathAddressRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
			
	}
	user := data.Get("user").(string)

	var rawData = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/address/"+user], &rawData); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	resp.Data["address"] = rawData["address"]
	
	return resp, nil
}

const objectHelpSyn = `
TODO
`
const objectHelpDesc = `
TODO
`
