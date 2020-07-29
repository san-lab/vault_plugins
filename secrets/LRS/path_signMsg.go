package mock

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	//"github.com/ethereum/go-ethereum/crypto/secp256k1"
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
	if err := jsonutil.DecodeJSON(b.store["/ringSignature/publicKeys/"], &rawData2); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	resp.Data["result"] = LRS(rawData["privKey"], msg, rawData2["pubKey1_X"], rawData2["pubKey2_X"], rawData2["pubKey3_X"], rawData2["pubKey4_X"], rawData2["pubKey5_X"])
	
	return resp, nil
}

/*def createSignature(myPrivKey, publicKeys, msg, pos, curveUsed):
    print(curveUsed)
    curv = Curve.get_curve(curveUsed)
    _g = curv.generator
    _p = curv.field
    _q = curv.order

    modulo = len(publicKeys)

    msg = msg[2:]

    kec1 = keccak.new(digest_bits=256)
    L = ""
    for key in publicKeys:
        keyToHash = str(hex(key['x']))[2:]
        while len(keyToHash) < 64:
            keyToHash = '0' + keyToHash
        L = L + keyToHash

    L_byte_count = ceil(len(L)/2)
    L_bytes = int(L, 16).to_bytes(L_byte_count, "big")

    h_0 = kec1.update(L_bytes).hexdigest()

    #print("HBefore")
    #print(int(h_0,16))

    #print("Generator")
    #print(g)

    h = curv.mul_point(int(h_0,16), _g)

    #print("HAfter")
    #print(h)
    
    y_tilde = curv.mul_point(myPrivKey, h)

    c_list = [None] * len(publicKeys)
    u = getRandomNumber()
    s_list = list()
    for i in range (0, len(publicKeys)):
        s_list.append(getRandomNumber())

    P_ini = Point(publicKeys[pos]['x'], publicKeys[pos]['y'], curv)
    c_list[pos] = H1(u,  P_ini, 0, y_tilde, h, L, msg, curv, _g)
    j = ((pos+1) % modulo)

    while j != pos:
        P = Point(publicKeys[j]['x'], publicKeys[j]['y'], curv)
        c_list[j] = H1(s_list[j], P, c_list[(j-1) % modulo], y_tilde, h, L, msg, curv, _g)
        j = (j+1) % modulo

    s_last = (u - myPrivKey * c_list[(j-1) % modulo]) % _q
    s_list[pos] = s_last
    P_last = Point(publicKeys[pos]['x'], publicKeys[pos]['y'], curv)

    c_last = H1(s_last, P_last, c_list[(j-1) % modulo], y_tilde, h, L, msg, curv, _g)


    return {'C': c_list[0], 'S_list': s_list, 'Y_tilde':{ 'x': y_tilde.x, 'y' : y_tilde.y }, 'msg' : msg }
*/

func LRS(PrivKeyHex, msg, pubk1,pubk2,pubk3,pubk4,pubk5 string) (string){
	//koblitz := secp256k1.S256() 
	//Gx y Gy coordernadas generador.
	//P is equivalent to _p on the python code
	//N is equivalent to _q on the python code
    
    return "aa"
}

/*func (c *H1ForLRS) Run(input []byte) ([]byte, error) {
	koblitz := secp256k1.S256()

	Gx := new(big.Int).SetBytes(getData(input, 0, 32))
	Gy := new(big.Int).SetBytes(getData(input, 32, 32))

	Y_x := new(big.Int).SetBytes(getData(input, 96, 32))
	Y_y := new(big.Int).SetBytes(getData(input, 128, 32))

	t1_x, t1_y := koblitz.ScalarMult(Gx, Gy, input[64:96])
	t2_x, t2_y := koblitz.ScalarMult(Y_x, Y_y, input[160:192])
	t3_x, t3_y := koblitz.Add(t1_x, t1_y, t2_x, t2_y)

	Hx := new(big.Int).SetBytes(getData(input, 192, 32))
	Hy := new(big.Int).SetBytes(getData(input, 224, 32))

	YT_x := new(big.Int).SetBytes(getData(input, 256, 32))
	YT_y := new(big.Int).SetBytes(getData(input, 288, 32))

	v1_x, v1_y := koblitz.ScalarMult(Hx, Hy, input[64:96])
	v2_x, v2_y := koblitz.ScalarMult(YT_x, YT_y, input[160:192])
	v3_x, v3_y := koblitz.Add(v1_x, v1_y, v2_x, v2_y)

	first_half := koblitz.Marshal(t3_x, t3_y)
	second_half := koblitz.Marshal(v3_x, v3_y)
	return_bytes := []byte{}

	return_bytes = append(first_half, second_half...)

	return return_bytes, nil
}*/

const queryHelpSyn = `
TODO
`
const queryHelpDesc = `
TODO
`
