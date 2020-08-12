package mock

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathQuery executes query operations against the CCP Web Service
func pathSignMsgExt(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: signMsgPathExt,
		Fields: map[string]*framework.FieldSchema{
			"msg": {
				Type:        framework.TypeString,
				Description: "Specifies the tx of the secret.",
			},
			"privKey": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pos": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK1_x": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK1_y": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK2_x": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK2_y": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK3_x": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK3_y": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK4_x": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK4_y": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK5_x": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
			"pubK5_y": {
				Type:        framework.TypeString,
				Description: "WIP",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathSignExtWrite,
			logical.UpdateOperation: b.pathSignExtWrite,
		},

		HelpSynopsis:    queryHelpSyn,
		HelpDescription: queryHelpDesc,
	}
}

func (b *backend) pathSignExtWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
			
	}
	msg := data.Get("msg").(string)
	privKey := data.Get("privKey").(string)
	pos := data.Get("pos").(string)
	pubKey1_x := data.Get("pubK1_x").(string)
	pubKey1_y := data.Get("pubK1_y").(string)
	pubKey2_x := data.Get("pubK2_x").(string)
	pubKey2_y := data.Get("pubK2_y").(string)
	pubKey3_x := data.Get("pubK3_x").(string)
	pubKey3_y := data.Get("pubK3_y").(string)
	pubKey4_x := data.Get("pubK4_x").(string)
	pubKey4_y := data.Get("pubK4_y").(string)
	pubKey5_x := data.Get("pubK5_x").(string)
	pubKey5_y := data.Get("pubK5_y").(string)
	
	resp.Data["result"] = LRS(privKey, pos, msg, pubKey1_x, pubKey2_x, pubKey3_x, pubKey4_x, pubKey5_x, pubKey1_y, pubKey2_y, pubKey3_y, pubKey4_y, pubKey5_y)
	
	return resp, nil
}