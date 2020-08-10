package mock

import (
	"context"
	"fmt"
    "math/big"
	"encoding/hex"
    //"strconv"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/btcsuite/btcd/btcec"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const signTxPathRegExp = signTxPath + "/(?P<user>[^/]+)$"
// pathQuery executes query operations against the CCP Web Service
func pathSignTx(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: signTxPathRegExp,
		Fields: map[string]*framework.FieldSchema{
			"tx": {
				Type:        framework.TypeString,
				Description: "Specifies the tx of the secret.",
			},
			"user": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathSignWrite,
			logical.UpdateOperation: b.pathSignWrite,
		},

		HelpSynopsis:    queryHelpSyn,
		HelpDescription: queryHelpDesc,
	}
}

func (b *backend) pathSignWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
			
	}
	tx := data.Get("tx").(string)
	user := data.Get("user").(string)

	transaction := new(types.Transaction)

	err := transaction.UnmarshalJSON([]byte(tx))

	if err != nil{
		resp.Data["error"] = fmt.Sprint(err)
		return resp, nil
	}

	
	var rawData = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/key/"+user], &rawData); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	resp.Data["result"] = signTransaction(rawData["ethkey"], transaction)
	

	
	
	return resp, nil
}

func signTransaction(PrivKeyHex string, tx *types.Transaction) (string){

    bts, err := hex.DecodeString(PrivKeyHex[:])
    if err !=nil{
        fmt.Println(err)
                return "error"
    }
    priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), bts)
    privateKey := priv.ToECDSA()
    //publicKey := privateKey.PublicKey
    //address := crypto.PubkeyToAddress(publicKey).Hex()

    //i, err := strconv.Atoi(nonce)
    //nonceUint := uint64(i)
    txN := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), nil)
    signTx, _ := types.SignTx(txN, types.NewEIP155Signer(big.NewInt(4)),privateKey)
    marshalledTXSigned, _ := signTx.MarshalJSON()
    return string(marshalledTXSigned)
}

const queryHelpSyn = `
TODO
`
const queryHelpDesc = `
TODO
`
