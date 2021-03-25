package ethPlugin

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func HsmCallMock(b *backend, user string, req *logical.Request, hashToSign []byte) (signature []byte, err error) {
	var rawDataPriv = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/key/"+user], &rawDataPriv); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}
	privKeyStr := rawDataPriv["ethkey"]
	privBytes, _ := hex.DecodeString(privKeyStr)
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privBytes)
	sign, err := privateKey.Sign(hashToSign)
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/key/"+user], &rawDataPriv); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}
	return sign.Serialize(), nil
}
