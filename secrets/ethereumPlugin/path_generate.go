package mock

import (
	"context"
	"encoding/json"
	"fmt"

    "github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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

	user := data.Get("user").(string)

	ethKeyGen, _ := crypto.GenerateKey()
	publicKey := ethKeyGen.PublicKey
    address := crypto.PubkeyToAddress(publicKey).Hex()

    reqDataCopy := make(map[string]interface{})
    for key, value := range req.Data {
	  reqDataCopy[key] = value
	}

	// JSON encode the data
	req.Data["address"] = address
	bufAddr, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	reqDataCopy["ethkey"] = fmt.Sprintf("%x", ethKeyGen.D.Bytes())
	bufKey, err := json.Marshal(reqDataCopy)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	// Store kv pairs in map at specified path
	b.store[req.ClientToken+"/address/"+user] = bufAddr
	b.store[req.ClientToken+"/key/"+user] = bufKey

	return nil, nil
}

const confHelpSyn = `
TODO.
`
const confHelpDesc = `
TODO.
`
