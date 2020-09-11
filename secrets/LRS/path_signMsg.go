package LRS

import (
	"context"
	"fmt"
	"encoding/hex"
	"strconv"
	"strings"
	"math/big"
	"math/rand"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
    "github.com/ethereum/go-ethereum/crypto"
	"github.com/HcashOrg/hcashd/hcashec/secp256k1"
)

const signMsgPathRegExp = signMsgPath + "/(?P<user>[^/]+)$"
// pathQuery executes query operations against the CCP Web Service
func pathSignMsg(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: signMsgPathRegExp,
		Fields: map[string]*framework.FieldSchema{
			"msg": {
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
	msg := data.Get("msg").(string)
	user := data.Get("user").(string)

	
	var rawData = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store["/ringSignature/private/"+user], &rawData); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	var rawData2 = map[string]string{}
	if err := jsonutil.DecodeJSON(b.store["/ringSignature/publicKeys"], &rawData2); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}
	
	resp.Data["result"] = LRS(rawData["privKey"], rawData["pos"], msg, rawData2["pubKey1_X"], rawData2["pubKey2_X"], rawData2["pubKey3_X"], rawData2["pubKey4_X"], rawData2["pubKey5_X"], rawData2["pubKey1_Y"], rawData2["pubKey2_Y"], rawData2["pubKey3_Y"], rawData2["pubKey4_Y"], rawData2["pubKey5_Y"])
	
	return resp, nil
}

//We assume that PrivKeyHex has no 0x prefix
func LRS(PrivKeyHex, pos_str, msg, pubk1_x,pubk2_x,pubk3_x,pubk4_x,pubk5_x,pubk1_y,pubk2_y,pubk3_y,pubk4_y,pubk5_y string) (string){
	koblitz := secp256k1.S256()
	pos, _ := strconv.Atoi(pos_str)

	pubKeys_x := []string{pubk1_x, pubk2_x, pubk3_x, pubk4_x, pubk5_x}
	pubKeys_y := []string{pubk1_y, pubk2_y, pubk3_y, pubk4_y, pubk5_y}
    
    L := pubk1_x + pubk2_x + pubk3_x + pubk4_x + pubk5_x
    L_bytes, _ := hex.DecodeString(L)
    h_0 := crypto.Keccak256(L_bytes)
    h_x, h_y := koblitz.ScalarMult(koblitz.Gx, koblitz.Gy, h_0)

    privKey_bytes, _ := hex.DecodeString(PrivKeyHex)
    y_tilde_x, y_tilde_y := koblitz.ScalarMult(h_x, h_y, privKey_bytes)


    u := randBigInt()
    s_list := make([]*big.Int, 5)
    for i := 0; i < 5; i++{
    	s_list[i] = randBigInt()
    }
    c_list := make([]*big.Int, 5)
    c_list[pos] = H1(u, pubKeys_x[pos], pubKeys_y[pos], big.NewInt(0), y_tilde_x, y_tilde_y, h_x, h_y, L, msg, koblitz)	


    j := ((pos+1) % 5)
    prev_j := (j-1) % 5
    for j != pos {
    	prev_j = (j-1) % 5
    	if prev_j == -1 {
    		prev_j = 4
    	}
    	c_list[j] = H1(s_list[j], pubKeys_x[j], pubKeys_y[j], c_list[prev_j], y_tilde_x, y_tilde_y, h_x, h_y, L, msg, koblitz)
    	j = (j+1) % 5
    }

    parsed := new(big.Int)
    parsed.SetString(PrivKeyHex, 16)
    s_last_mul := new(big.Int)
 	prev_j = (j-1) % 5
 	if prev_j == -1 {
		prev_j = 4
	}
    s_last_mul.Mul(parsed, c_list[prev_j])
    s_last_sub := new(big.Int)
    s_last_sub.Sub(u,s_last_mul)
    s_last := new(big.Int)
    s_last.Mod(s_last_sub,koblitz.N)
    s_list[pos] = s_last

    signature_str := "{'C': 0x" + c_list[0].Text(16) + ", 'S_list': [" + strings.Trim(strings.Replace(fmt.Sprint(s_list), " ", ",", -1), "[]") + "], 'Y_tilde':{ 'x': 0x" + y_tilde_x.Text(16) +" , 'y' : 0x" + y_tilde_y.Text(16) + "}, 'msg' : \"" + msg + "\"}"

    return signature_str
}

func H1(s *big.Int, pubK_x, pubK_y string, c, y_tilde_x, y_tilde_y , h_x, h_y *big.Int, L, msg string, koblitz *secp256k1.KoblitzCurve) (*big.Int){
	Y_x := new(big.Int)
	Y_x.SetString(pubK_x, 16)
	Y_y := new(big.Int)
	Y_y.SetString(pubK_y, 16)

	s_bytes := s.Bytes()
	c_bytes := c.Bytes()
	if len(c_bytes) == 0 {
		c_bytes =  []byte{0}
	} 

	t1_x, t1_y := koblitz.ScalarMult(koblitz.Gx, koblitz.Gy, s_bytes)
	t2_x, t2_y := koblitz.ScalarMult(Y_x, Y_y, c_bytes)
	t3_x, _ := koblitz.Add(t1_x, t1_y, t2_x, t2_y)
	
	v1_x, v1_y := koblitz.ScalarMult(h_x, h_y, s_bytes)
	v2_x, v2_y := koblitz.ScalarMult(y_tilde_x, y_tilde_y, c_bytes)
	v3_x, _ := koblitz.Add(v1_x, v1_y, v2_x, v2_y)

    ytildaxToHash := y_tilde_x.Text(16)
    t3ToHash := t3_x.Text(16)
    v3ToHash := v3_x.Text(16)

    for len(ytildaxToHash) < 64 {
        ytildaxToHash = "0" + ytildaxToHash
    }

	messageToHash := msg
    for len(messageToHash) < 64 {
        messageToHash = "0" + messageToHash
    }

    for len(t3ToHash) < 64 {
        t3ToHash = "0" + t3ToHash
    }

    for len(v3ToHash) < 64 {
        v3ToHash = "0" + v3ToHash
    }

    str_Tohash := L + ytildaxToHash + messageToHash + t3ToHash + v3ToHash
    ToHash_bytes, _ := hex.DecodeString(str_Tohash)
    nextC_bytes := crypto.Keccak256(ToHash_bytes)
    nextC := new(big.Int)
    nextC.SetBytes(nextC_bytes)

    return nextC
}

func randBigInt() (*big.Int){
	randomSeed := make([]byte, 32)
    rand.Read(randomSeed)
	randomBigInt := new(big.Int)
    return randomBigInt.SetBytes(randomSeed)
}

const queryHelpSyn = `
TODO
`
const queryHelpDesc = `
TODO
`
