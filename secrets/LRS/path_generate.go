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
			"user1": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
			"user2": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
			"user3": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
			"user4": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
			"user5": {
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

	user1 := data.Get("user1").(string)
	user2 := data.Get("user2").(string)
	user3 := data.Get("user3").(string)
	user4 := data.Get("user4").(string)
	user5 := data.Get("user5").(string)

	ringSK1, _ := crypto.GenerateKey()
	ringPK1 := ringSK1.PublicKey
	ringSK2, _ := crypto.GenerateKey()
	ringPK2 := ringSK2.PublicKey
	ringSK3, _ := crypto.GenerateKey()
	ringPK3 := ringSK3.PublicKey
	ringSK4, _ := crypto.GenerateKey()
	ringPK4 := ringSK4.PublicKey
	ringSK5, _ := crypto.GenerateKey()
	ringPK5 := ringSK5.PublicKey

	baseDataPublicKeys := make(map[string]interface{})

	for key, value := range req.Data {
	  baseDataPublicKeys[key] = value
	}

	baseDataPublicKeys["pubKey1_X"] = fmt.Sprintf("%x", ringPK1.X.Bytes())
	baseDataPublicKeys["pubKey1_Y"] = fmt.Sprintf("%x", ringPK1.Y.Bytes())
	baseDataPublicKeys["pubKey2_X"] = fmt.Sprintf("%x", ringPK2.X.Bytes())
	baseDataPublicKeys["pubKey2_Y"] = fmt.Sprintf("%x", ringPK2.Y.Bytes())
	baseDataPublicKeys["pubKey3_X"] = fmt.Sprintf("%x", ringPK3.X.Bytes())
	baseDataPublicKeys["pubKey3_Y"] = fmt.Sprintf("%x", ringPK3.Y.Bytes())
	baseDataPublicKeys["pubKey4_X"] = fmt.Sprintf("%x", ringPK4.X.Bytes())
	baseDataPublicKeys["pubKey4_Y"] = fmt.Sprintf("%x", ringPK4.Y.Bytes())
	baseDataPublicKeys["pubKey5_X"] = fmt.Sprintf("%x", ringPK5.X.Bytes())
	baseDataPublicKeys["pubKey5_Y"] = fmt.Sprintf("%x", ringPK5.Y.Bytes())
	
	bufPubKeys, err := json.Marshal(baseDataPublicKeys)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	// Store kv pairs in map at specified path
	b.store["/ringSignature/publicKeys"] = bufPubKeys

	baseDataPriv1 := make(map[string]interface{})

	for key, value := range req.Data {
	  baseDataPriv1[key] = value
	}

	baseDataPriv1["privKey"] = fmt.Sprintf("%x", ringSK1.D.Bytes())
	baseDataPriv1["pos"] = "0"
	bufPriv1, err := json.Marshal(baseDataPriv1)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}
	path1 := "/ringSignature/private/"+user1
	b.store[path1] = bufPriv1

	baseDataPriv2 := make(map[string]interface{})

	for key, value := range req.Data {
	  baseDataPriv2[key] = value
	}

	baseDataPriv2["privKey"] = fmt.Sprintf("%x", ringSK2.D.Bytes())
	baseDataPriv2["pos"] = "1"
	bufPriv2, err := json.Marshal(baseDataPriv2)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}
	path2 := "/ringSignature/private/"+user2
	b.store[path2] = bufPriv2

	baseDataPriv3 := make(map[string]interface{})

	for key, value := range req.Data {
	  baseDataPriv3[key] = value
	}

	baseDataPriv3["privKey"] = fmt.Sprintf("%x", ringSK3.D.Bytes())
	baseDataPriv3["pos"] = "2"
	bufPriv3, err := json.Marshal(baseDataPriv3)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}
	path3 := "/ringSignature/private/"+user3
	b.store[path3] = bufPriv3

	baseDataPriv4 := make(map[string]interface{})

	for key, value := range req.Data {
	  baseDataPriv4[key] = value
	}

	baseDataPriv4["privKey"] = fmt.Sprintf("%x", ringSK4.D.Bytes())
	baseDataPriv4["pos"] = "3"
	bufPriv4, err := json.Marshal(baseDataPriv4)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}
	path4 := "/ringSignature/private/"+user4
	b.store[path4] = bufPriv4

	baseDataPriv5 := make(map[string]interface{})

	for key, value := range req.Data {
	  baseDataPriv5[key] = value
	}

	baseDataPriv5["privKey"] = fmt.Sprintf("%x", ringSK5.D.Bytes())
	baseDataPriv5["pos"] = "4"
	bufPriv5, err := json.Marshal(baseDataPriv5)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}
	path5 := "/ringSignature/private/"+user5
	b.store[path5] = bufPriv5

	//resp := &logical.Response{
	//	Data: map[string]interface{}{},
	//}

	//resp.Data["user1"] = baseDataPriv1["privKey"].(string)
	//resp.Data["user2"] = baseDataPriv2["privKey"].(string)
	//resp.Data["user3"] = baseDataPriv3["privKey"].(string)
	//resp.Data["user4"] = baseDataPriv4["privKey"].(string)
	//resp.Data["user5"] = baseDataPriv5["privKey"].(string)

	return nil, nil
}

const confHelpSyn = `
TODO.
`
const confHelpDesc = `
TODO.
`
