package ethPlugin

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	//"strconv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
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

	if err != nil {
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
	privBytes, _ := hex.DecodeString(privKeyStr)
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privBytes)
	pubKeyStr := rawDataPub["pubKey"]
	pubKeyCoordinates := strings.Split(pubKeyStr, ",")
	pubKeyX := new(big.Int)
	pubKeyX.SetString(pubKeyCoordinates[0], 10)
	pubKeyY := new(big.Int)
	pubKeyY.SetString(pubKeyCoordinates[1], 10)

	signer := types.NewEIP155Signer(big.NewInt(4))
	txN := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), nil)
	hash := signer.Hash(txN)

	//sign(priv, hash) -> barray
	hsm_signature, _ := privateKey.Sign(hash[:])

	hsm_resp := hsm_signature.Serialize()

	hsmR := new(big.Int).SetBytes(hsm_resp[4:36])
	hsmS := new(big.Int).SetBytes(hsm_resp[38:70])

	//barray -> (R,S)
	r := hsmR //hsm_signature.R
	s := hsmS //hsm_signature.S

	//use pubk to determine v -> (y coordinate) bigger than half of the prime of the field or even/odd

	var v byte

	//Calculate inverse of s
	s1 := new(big.Int)
	s1.ModInverse(s, btcec.S256().N)

	//R=s1*(r*Pb + h*G)
	Rx, Ry := btcec.S256().ScalarBaseMult(hash[:])
	Ax, Ay := btcec.S256().ScalarMult(pubKeyX, pubKeyY, r.Bytes())
	Rx, Ry = btcec.S256().Add(Rx, Ry, Ax, Ay)
	Rx, Ry = btcec.S256().ScalarMult(Rx, Ry, s1.Bytes())

	v = byte(Ry.Bit(0))

	signatureBytes := VRS_to_bytes(v, r, s)

	//append to signedTx
	signedTxStr := addSignatureToTransaction(signatureBytes, signer, txN)

	//resp.Data["intermidiateSign"] = hsm_signature_BigInt.String()
	resp.Data["result"] = signedTxStr
	resp.Data["resultPrev"] = signTransaction(privKeyStr, tx)

	return resp, nil
}

func addSignatureToTransaction(signatureBytes []byte, signer types.Signer, tx *types.Transaction) string {
	signTx, _ := tx.WithSignature(signer, signatureBytes)
	marshalledTXSigned, _ := signTx.MarshalJSON()
	return string(marshalledTXSigned)
}

//Returns a slice []byte with r, s (padded) and v (as 0,1)
func VRS_to_bytes(v byte, r *big.Int, s *big.Int) []byte {
	result := make([]byte, 0, 2*32+1)
	// Not sure this needs rounding but safer to do so.
	curvelen := 32

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

	result = append(result, v)
	return result
}

//METHODS USE ON PREVIOUS VERSION

func signTransaction(PrivKeyHex string, tx *types.Transaction) string {

	bts, err := hex.DecodeString(PrivKeyHex[:])
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), bts)
	txN := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), nil)
	signTx, _ := SignTx(txN, types.NewEIP155Signer(big.NewInt(4)), priv)
	marshalledTXSigned, _ := signTx.MarshalJSON()
	return string(marshalledTXSigned)
}

func SignTx(tx *types.Transaction, s types.Signer, prv *btcec.PrivateKey) (*types.Transaction, error) {
	h := s.Hash(tx)
	sig, err := Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	sigRet, errRet := tx.WithSignature(s, sig)

	return sigRet, errRet
}

func Sign(hash []byte, prv *btcec.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != btcec.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := btcec.SignCompact(btcec.S256(), prv, hash, false)
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
