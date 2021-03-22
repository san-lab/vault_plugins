package ethPlugin

import (
	"context"
	"fmt"
    "math/big"
	"encoding/hex"
	"crypto/ecdsa"
    //"strconv"


    "github.com/ethereum/go-ethereum/core/types"
    //"github.com/btcsuite/btcd/btcec"
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
    priv, _ := PrivKeyFromBytes(S256(), bts)
    privateKey := priv.ToECDSA()
    txN := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), nil)
    signTx, _ := SignTx(txN, types.NewEIP155Signer(big.NewInt(4)),privateKey)
    marshalledTXSigned, _ := signTx.MarshalJSON()
    return string(marshalledTXSigned)
}

func SignTx (tx *types.Transaction, s types.Signer, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	h := s.Hash(tx)
	sig, err := Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

func Sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := SignCompact(S256(), (*PrivateKey)(prv), hash, false)
	if err != nil {
		return nil, err
	}
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

const queryHelpSyn = `
TODO
`
const queryHelpDesc = `
TODO
`
