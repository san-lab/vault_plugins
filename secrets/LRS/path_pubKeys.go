package mock

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathPubKeys(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: showPubKeysPath + "$",
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

	var rawData = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store["/ringSignature/publicKeys"], &rawData); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	resp.Data["pubKey1_X"] = rawData["pubKey1_X"]
	resp.Data["pubKey1_Y"] = rawData["pubKey1_Y"]
	resp.Data["pubKey2_X"] = rawData["pubKey2_X"]
	resp.Data["pubKey2_Y"] = rawData["pubKey2_Y"]
	resp.Data["pubKey3_X"] = rawData["pubKey3_X"]
	resp.Data["pubKey3_Y"] = rawData["pubKey3_Y"]
	resp.Data["pubKey4_X"] = rawData["pubKey4_X"]
	resp.Data["pubKey4_Y"] = rawData["pubKey4_Y"]
	resp.Data["pubKey5_X"] = rawData["pubKey5_X"]
	resp.Data["pubKey5_Y"] = rawData["pubKey5_Y"]

	
	
	return resp, nil
}

const objectHelpSyn = `
TODO
`
const objectHelpDesc = `
TODO
`
