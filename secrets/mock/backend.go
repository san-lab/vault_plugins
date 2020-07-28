package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
    "math/big"
    "encoding/hex"
    //"strconv"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/common"
    "github.com/btcsuite/btcd/btcec"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{
		store: make(map[string][]byte),
	}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(mockHelp),
		BackendType: logical.TypeLogical,
	}

	b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	b.Backend.Setup(ctx, conf)

	return b, nil
}

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend

	store map[string][]byte
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: framework.MatchAllRegex("path"),

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Specifies the path of the secret.",
				},
				"tx": {
					Type:        framework.TypeString,
					Description: "Specifies the tx to be signed.",
				},
				"user": {
					Type:        framework.TypeString,
					Description: "Specifies the user to show its address.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRead,
					Summary:  "Retrieve the secret from the map.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
				/*logical.ListOperation: &framework.PathOperation{
					Callback: b.handleList,
					Summary:  "Lists the generated key",
				},*/
			},

			ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

func (b *backend) handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
			
	}
	path := data.Get("path").(string)
	tx := data.Get("tx").(string)
	user := data.Get("user").(string)

	if(tx != ""){
		transaction := new(types.Transaction)
	
		err := transaction.UnmarshalJSON([]byte(tx))

		if err != nil{
			resp.Data["error"] = fmt.Sprint(err)
			return resp, nil
		}

		// Decode the data
		var rawData = map[string]string{}
		if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/"+path], &rawData); err != nil {
			return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
		}

		resp.Data["result"] = signTransaction(rawData["ethKey"], transaction)
	}else if(user != ""){
		var rawData = map[string]string{}
		if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/"+path], &rawData); err != nil {
			return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
		}

		resp.Data["addressOfSigner"] = rawData["address"]
	}

	
	
	return resp, nil
}

/*func (b *backend) handleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}

	path := data.Get("path").(string)

	// Decode the data
	var rawData = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/"+path], &rawData); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	// Generate the response
	resp.Data["ethKey"] = "AAA"//rawData["ethKey"]
	return resp, nil
}*/

func (b *backend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("data must be provided to store in secret")
	}

	path := data.Get("path").(string)

	ethKeyGen, _ := crypto.GenerateKey()
	publicKey := ethKeyGen.PublicKey
    address := crypto.PubkeyToAddress(publicKey).Hex()

	// JSON encode the data
	req.Data["ethKey"] = fmt.Sprintf("%x", ethKeyGen.D.Bytes())
	req.Data["address"] = address
	buf, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}
	//buf = ethKeyGen.D.Bytes()

	// Store kv pairs in map at specified path
	b.store[req.ClientToken+"/"+path] = buf

	return nil, nil
}

func (b *backend) handleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	path := data.Get("path").(string)

	// Remove entry for specified path
	delete(b.store, path)

	return nil, nil
}

func signTransaction(PrivKeyHex string, tx *types.Transaction) (string){

    bts, err := hex.DecodeString(PrivKeyHex[2:])
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

type Txdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

const mockHelp = `
The Mock backend is a dummy secrets backend that stores kv pairs in a map.
`
