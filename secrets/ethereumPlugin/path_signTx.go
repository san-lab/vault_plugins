package ethPlugin

import (
	"context"
	"fmt"
    "math/big"
	"encoding/hex"
	"crypto/ecdsa"
	"strings"
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
	transaction := data.Get("tx").(string)
	user := data.Get("user").(string)

	tx := new(types.Transaction)

	err := tx.UnmarshalJSON([]byte(transaction))

	if err != nil{
		resp.Data["error"] = fmt.Sprint(err)
		return resp, nil
	}
	
	var rawDataPriv = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/key/"+user], &rawDataPriv); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	var rawDataPub = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/pubKey/"+user], &rawDataPub); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	//priv, pub and hash
	privKeyStr := rawDataPriv["ethkey"]
	privateKeyECDSA, _ := converToECDSA(privKeyStr)
	pubKeyStr := rawDataPub["pubKey"]
	pubKeyCoordinates := strings.Split(pubKeyStr, ",")
	pubKeyY := new(big.Int)
	pubKeyY.SetString(pubKeyCoordinates[1], 10)

	signer := types.NewEIP155Signer(big.NewInt(4))
	txN := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), nil)
   	hash := signer.Hash(txN)

	//sign(priv, hash) -> barray
	hsm_signature, _ := privateKeyECDSA.Sign(hash[:])
	hsm_signature_BigInt := new(big.Int)
	hsm_signature_BigInt.SetBytes(hsm_signature)

	//barray -> (R,S)
	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(hsm_signature[0:32])
	s.SetBytes(hsm_signature[32:64])

	//use pubk to determine v -> (y coordinate) bigger than half of the prime of the field or even/odd
	curve := S256()
	var v byte
	//if(pubKeyY.Cmp(curve.HalfOrder()) == -1){
	if(pubKeyY.Bit(0) == 0){
		v = 27
	} else {
		v = 28
	}
	signatureBytes := VRS_to_bytes(curve, v, r, s)

	//append to signedTx
	signedTxStr := addSignatureToTransaction(signatureBytes, signer, txN)

	resp.Data["intermidiateSign"] = hsm_signature_BigInt.String()
	resp.Data["result"] = signedTxStr
	resp.Data["resultPrev"] = signTransaction(privKeyStr, tx)
	
	return resp, nil
}

func addSignatureToTransaction(signatureBytes []byte, signer types.Signer, tx *types.Transaction) (string) {
	signTx, _ := tx.WithSignature(signer, signatureBytes)
	marshalledTXSigned, _ := signTx.MarshalJSON()
	return string(marshalledTXSigned)
}

func VRS_to_bytes(curve *KoblitzCurve, v byte, r *big.Int, s *big.Int) ([]byte) {
	result := make([]byte, 1, 2*curve.byteSize+1)
	result[0] = v
	// Not sure this needs rounding but safer to do so.
	curvelen := (curve.BitSize + 7) / 8

	// Pad R and S to curvelen if needed.
	bytelen := (r.BitLen() + 7) / 8
	if bytelen < curvelen {
		result = append(result,
			make([]byte, curvelen-bytelen)...)
	}
	result = append(result, r.Bytes()...)

	bytelen = (s.BitLen() + 7) / 8
	if bytelen < curvelen {
		result = append(result,
			make([]byte, curvelen-bytelen)...)
	}
	result = append(result, s.Bytes()...)

	v_aux := result[0] - 27
	copy(result, result[1:])
	result[64] = v_aux

	return result
}

func converToECDSA (privKeyStr string) (*PrivateKey, error) {
	bts, err := hex.DecodeString(privKeyStr[:])
    if err !=nil{
        fmt.Println(err)
                return nil, fmt.Errorf("incorrect priv key")
    }
    priv, _ := PrivKeyFromBytes(S256(), bts)
    privateKeyECDSA := priv.ToECDSA()

    return (*PrivateKey)(privateKeyECDSA), nil
}

//METHODS USE ON PREVIOUS VERSION

func signTransaction(PrivKeyHex string, tx *types.Transaction) (string){

    bts, err := hex.DecodeString(PrivKeyHex[:])
    if err !=nil{
        fmt.Println(err)
                return "error"
    }
    priv, _ := PrivKeyFromBytes(S256(), bts)
    privateKey := priv.ToECDSA()
    txN := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), nil)
    _, signTx, _ := SignTx(txN, types.NewEIP155Signer(big.NewInt(4)),privateKey)
    marshalledTXSigned, _ := signTx.MarshalJSON()
    return string(marshalledTXSigned)
}

func SignTx(tx *types.Transaction, s types.Signer, prv *ecdsa.PrivateKey) (*big.Int, *types.Transaction, error) {
	h := s.Hash(tx)
	sig, precomputeSig, err := Sign(h[:], prv)
	if err != nil {
		return nil, nil, err
	}
	sigRet, errRet := tx.WithSignature(s, sig)

	return precomputeSig, sigRet, errRet
}

func Sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, *big.Int, error) {
	if len(hash) != 32 {
		return nil, nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != S256() {
		return nil, nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, precomputeSig, err := SignCompact(S256(), (*PrivateKey)(prv), hash, false)
	if err != nil {
		return nil, nil, err
	}
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, precomputeSig, nil
}

const queryHelpSyn = `
TODO
`
const queryHelpDesc = `
TODO
`
