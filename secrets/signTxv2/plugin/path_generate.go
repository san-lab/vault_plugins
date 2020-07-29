package ccpsecrets

import (
	"context"
	"encoding/json"
	"fmt"

    "github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	//"github.com/mitchellh/mapstructure"
)

func pathGenerate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: generatePath + "$",
		Fields: map[string]*framework.FieldSchema{
			"user": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathGenerateWrite,
			logical.UpdateOperation: b.pathGenerateWrite,
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

// pathConfigRead handles read commands to the config
func (b *backend) pathGenerateWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("data must be provided to store in secret")
	}

	user := data.Get("user").(string)

	ethKeyGen, _ := crypto.GenerateKey()
	publicKey := ethKeyGen.PublicKey
    address := crypto.PubkeyToAddress(publicKey).Hex()

	// JSON encode the data
	req.Data["address"] = address
	buf, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	// Store kv pairs in map at specified path
	b.store[req.ClientToken+"/address/"+user] = buf

	return nil, nil
}

const confHelpSyn = `
TODO.
`
const confHelpDesc = `
TODO.
`