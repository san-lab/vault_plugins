package ethPlugin

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

	reqDataCopy2 := make(map[string]interface{})
    for key, value := range req.Data {
	  reqDataCopy2[key] = value
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

	reqDataCopy2["pubKey"] = publicKey.X.String()+","+publicKey.Y.String()
	bufPubKey, err := json.Marshal(reqDataCopy2)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	// Store kv pairs in map at specified path
	b.store[req.ClientToken+"/address/"+user] = bufAddr
	b.store[req.ClientToken+"/key/"+user] = bufKey
	b.store[req.ClientToken+"/pubKey/"+user] = bufPubKey

	return nil, nil
}

const confHelpSyn = `
Command that generates a new private
key for the specified user
`
const confHelpDesc = `
This command creates the private key
and stores it on the user private path 
so it's only accesible to that user
`
