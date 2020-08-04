package mock

import (
	"context"
	"fmt"
	//"time"
	//"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"math/big"
	//"math/rand"
	//"golang.org/x/crypto/sha3"
	//"crypto/rand"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
    "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// pathQuery executes query operations against the CCP Web Service
func pathSignMsg(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: signMsgPath + "$",
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
	//Gx y Gy coordernadas generador.
	//P is equivalent to _p on the python code
	//N is equivalent to _q on the python code
	pubKeys_x := []string{pubk1_x, pubk2_x, pubk3_x, pubk4_x, pubk5_x}
	pubKeys_y := []string{pubk1_y, pubk2_y, pubk3_y, pubk4_y, pubk5_y}
    
    L := pubk1_x + pubk2_x + pubk3_x + pubk4_x + pubk5_x
    L_bytes := []byte(L)
    h_0 := crypto.Keccak256(L_bytes)
    h_x, h_y := koblitz.ScalarMult(koblitz.Gx, koblitz.Gy, h_0)

    privKey_bytes, _ := hex.DecodeString(PrivKeyHex)
    y_tilde_x, y_tilde_y := koblitz.ScalarMult(h_x, h_y, privKey_bytes)

    u := big.NewInt(5) //This should be random
    s_list := make([]*big.Int, 5)
    for i := 0; i < 5; i++{
    	s_list[i] = big.NewInt(5) //This should be random
    }
    c_list := make([]*big.Int, 5)
    //c_list[pos] = H1(u, pubKeys_x[pos], pubKeys_y[pos], big.NewInt(0), y_tilde_x, y_tilde_y, h_x, h_y, L, msg, koblitz)
    ret := H1(u, pubKeys_x[pos], pubKeys_y[pos], big.NewInt(0), y_tilde_x, y_tilde_y, h_x, h_y, L, msg, koblitz)
    return koblitz.Gx.Text(16) + "      " + ret

    j := ((pos+1) % 5)

    for j != pos {
    	prev_j := (j-1) % 5
    	if prev_j == -1 {
    		prev_j = 4
    	}
    	//c_list[j] = H1(s_list[j], pubKeys_x[j], pubKeys_y[j], c_list[prev_j], y_tilde_x, y_tilde_y, h_x, h_y, L, msg, koblitz)
    	j = (j+1) % 5
    }

    parsed := new(big.Int)
    parsed.SetString(PrivKeyHex, 16)
    s_last_mul := new(big.Int)
    s_last_mul.Mul(parsed, c_list[(j-1) % 5])
    s_last_sub := new(big.Int)
    s_last_sub.Sub(u,s_last_mul)
    s_last := new(big.Int)
    s_last.Mod(s_last_sub,koblitz.N)
    s_list[pos] = s_last

    //c_last := H1(s_last, pubKeys_x[pos], pubKeys_y[pos], c_list[(j-1) % 5], y_tilde_x, y_tilde_y, h_x, h_y, L, msg, koblitz)

    signature_str := "{'C': " + c_list[0].Text(16) + " 'S_list': [" + strings.Trim(strings.Replace(fmt.Sprint(s_list), " ", ",", -1), "[]") + "] 'Y_tilde':{ 'x': " + y_tilde_x.Text(16) +" , 'y' : " + y_tilde_y.Text(16) + "}, 'msg' : " + msg + "}"

    return signature_str
}

func H1(s *big.Int, pubK_x, pubK_y string, c, y_tilde_x, y_tilde_y , h_x, h_y *big.Int, L, msg string, koblitz *secp256k1.BitCurve) (string){
	Y_x := new(big.Int)
	Y_x.SetString(pubK_x, 16)
	Y_y := new(big.Int)
	Y_y.SetString(pubK_y, 16)

	s_bytes := s.Bytes()
	c_bytes := c.Bytes()
	if len(c_bytes) == 0 {
		c_bytes =  []byte{1}
	} 

	t1_x, t1_y := koblitz.ScalarMult(koblitz.Gx, koblitz.Gy, s_bytes)
	t2_x, t2_y := koblitz.ScalarMult(Y_x, Y_y, c_bytes)
	return t1_x.Text(16) + " " + t1_y.Text(16) + " " + t2_x.Text(16) + " " + t2_y.Text(16)
	t3_x, _ := koblitz.Add(t1_x, t1_y, t2_x, t2_y)
	
	v1_x, v1_y := koblitz.ScalarMult(h_x, h_y, s_bytes)
	v2_x, v2_y := koblitz.ScalarMult(y_tilde_x, y_tilde_y, c_bytes)
	v3_x, _ := koblitz.Add(v1_x, v1_y, v2_x, v2_y)

    ytildaxToHash := y_tilde_x.Text(16)
    t3ToHash := t3_x.Text(16)
    v3ToHash := v3_x.Text(16)

    str_Tohash := L + ytildaxToHash + msg + t3ToHash + v3ToHash
    ToHash_bytes := []byte(str_Tohash)
    nextC_bytes := crypto.Keccak256(ToHash_bytes)
    nextC := new(big.Int)
    nextC.SetBytes(nextC_bytes)

    return "aa"
    //return nextC
}

const queryHelpSyn = `
TODO
`
const queryHelpDesc = `
TODO
`
