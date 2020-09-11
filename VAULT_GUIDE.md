# Quick guide Hashicorp-Vault:

## Initialize and unseal the vault:
First we will need to specify the base configuration for our vault. Let's take this config file from their documentation as reference.
```
storage "raft" {
  path    = "./vault/data"
  node_id = "node1"
}

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = 1
}

api_addr = "http://127.0.0.1:8200"
cluster_addr = "https://127.0.0.1:8201"
ui = true
```
First we are setting were are we going to set the raft storage and giving it the id of "node1". Second we are setting the tcp listening IP and port (as you can see we are doing this demo locally). Then, we set the specify the address that we will use to reach the API and the address on which the raft cluster will be listening. And lastly we just say that we want to enable the User Interface for our vault.
The only remaining thing before starting the vault would be to create the actual storage folder.
```
 mkdir -p vault/data
```

Once we have this config file ready we can start the vault by executing:
```
vault server -config=config.hcl
```
This will start the vault and keep running on that console where you can see any log that it produces
Now we need to set the enviroment variable
```
export VAULT_ADDR="http://127.0.0.1:8200"
```
The next step is to initialize the vault for this we need to open a new terminal tab and type.
```
vault operator init
```
This will give us 5 unseal keys that we need to store somewhere because they are needed both for sealing and unsealing the vault 3 out of those 5 keys need to be used in order to unseal the vault. 
By default when a vault is initialized it starts always as sealed and we cannot interact with it until we unseal it. At the same time the vault could be sealed again in the future in case we want it to keep all the information it already had but being imposible to operate with for some time.
It is also important to keep the root token since that is the token that will allow as to log in with root permission after we have unsealed the vault.
Now we need to unseal the vault by executing the following command with 3 out of the 5 unseal keys.
```
vault operator unseal
```
Finally the vault is unsealed and we can start interacting with it. Do not forget the root login token since is what you need to use in order to access the vault and start configurating it.
```
vault login <root-token>
```

## Policies and profiles:
Policies are used to specify different profiles and which routes can be accesed by each of the profiles. Everything is expressed as paths from the secret engines, to the custom plugins and even the different API verbs of a given plugin.
Policies are defined by creating different “.hcl” files. Here we have an example of a very simple one.
```
path "LRS/signMsg/pedro" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

path "LRS/showPubKeys" {
  capabilities = [ "read"]
}
```
In this case we are creating a user which will have rights to operate on those 2 paths: “LRS/signMsg/pedro” and “LRS/showPubKeys”. We can see that under “LRS/signMsg/pedro” he will have all permissions (create, read, update, delete and list) while on “LRS/showPubKeys” he will have permissions only for reading.

Once we have this configuration file all we need is to execute this command to create an actual policy on the vault.
```
vault policy write pedro user1.hcl
```
In this case we are loading the config from "user1.hcl" that is the name we gave the upper configuration file and we are creating a policy named "pedro".

### Activate userpass module and log in using it
Now that we have a policy created we can enable the userpass module on the vault, which allows us to log in using a username/password.
```
vault auth enable userpass
```
After that we need to create the new user/password with a associated policy.
```
vault write /auth/userpass/users/pedro password=prueba policies=pedro
```
Here we are specifying those 3 things: the username at the end of the path and after that the password and the set of policies.
Once we have a user created the user knowing its user/password can execute:
```
vault write /auth/userpass/login/guille password=prueba
```
which will return many information, for now lets just take the "token" that which is always a string starting by "s."
```
s.SxHmNObbq0KTTydqWM5vxsfD
```
having this token finally the user can login executing:
```
vault login <token>
```

## Vault Plugin system
### Quick thoughts on the plugin system.
1. Having Golang as the language to program the plugins is a very good feature, since it is a language that has many libraries for security and cryptography handling. This allows you to use any preexisting code, we did so in the transaction signer plugin where we are reusing the actual code  that is used for signing inside Go-ethereum.
2. Now a days the documentation is inexistent, all the information that we used to program the pluging was either obtained talking to actual Hashicorp Engineers or by looking at the actual code from Hashicorp Vault.

### Quick guide to build you own plugin
First you need to create the main.go file that is what links you plugin to the rest of vault, this code is always the same and does not depend on your plugin.
```go
package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	mock "github.com/hashicorp/vault-guides/secrets/mock"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: mock.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
```

Now we need to create the file where we define Factory that we used on main.go.
This is the point where we define how many different funcitonalities will our plugin have and which will be the path for each of them.

```go
package ethPlugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const generatePath string = "genKey"
const showAddressPath string = "showAddr"
const signTxPath string = "signTx"

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{
		store: make(map[string][]byte),
	}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(ethereumPluginHelp),
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathGenerate(b),
				pathAddress(b),
				pathSignTx(b),
			},
		),
	}

	//b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

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

const ethereumPluginHelp = `
The ethereumPlugin backend is a plugin that allows you to input a ethereum transaction and returns it signed.
`
```
The main parts of the code are:
1. After the imports we are defining 3 constants. Each one of them is one of the verbs that will be used on the app to access the different functionalities of the plugin(in this case generating a key pair, showing the address related to the key pair and the actual signing of a transaction using that key pair).
2. On the definition of the Factory function take a look on the definition of the different paths. Those are the 3 function that will be called depending on the API verb used (pathGenerate, pathAddress and pathSignTx) we will see how they are defined on the next step.

Finally we need to define the functionality for each of the paths that we defined on the factory. We will have a look at how we defined the "pathGenerate" function since the process to generate all of them is pretty similar.
```go
package ethPlugin

import (
	"context"
	"encoding/json"
	"fmt"

    "github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathGenerate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: generatePath + "$",
		Fields: map[string]*framework.FieldSchema{
			"user": {
				Type:        framework.TypeString,
				Description: "Specifies the user of the secret.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathGenerateWrite,
			logical.UpdateOperation: b.pathGenerateWrite,
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

func (b *backend) pathGenerateWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	user := data.Get("user").(string)

	ethKeyGen, _ := crypto.GenerateKey()
	publicKey := ethKeyGen.PublicKey
    address := crypto.PubkeyToAddress(publicKey).Hex()

    reqDataCopy := make(map[string]interface{})
    for key, value := range req.Data {
	  reqDataCopy[key] = value
	}

	// JSON encode the data
	req.Data["address"] = address
	bufAddr, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	reqDataCopy["ethkey"] = fmt.Sprintf("%x", ethKeyGen.D.Bytes())
	bufKey, err := json.Marshal(reqDataCopy)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	// Store kv pairs in map at specified path
	b.store[req.ClientToken+"/address/"+user] = bufAddr
	b.store[req.ClientToken+"/key/"+user] = bufKey
	
	resp := &logical.Response{
		Data: map[string]interface{}{},
			
	}
	resp.Data["address"] = address

	return resp, nil
}

const confHelpSyn = `
TODO.
`
const confHelpDesc = `
TODO.
`
```
Inside pathGenerate first we need to specify the different input parameters that will be accepted by the command. In this case we only have one that is "user" of type string.

After that we need to specify the different callbacks that our command will have. Each command can have 5 different callbacks: CreateOperation, ReadOperation, UpdateOperation, DeleteOperation and ListOperation. As you can see in this case we are only defining behaviour for create and update, the rest of them will be undefined and will show a warning if they are called.

Finally we need to write the function with the whole logic for each of those different callbacks. In this case let's take a look at our "pathGenerateWrite" function.
- Inside data we have all the input parameters that we specified previously. We just need to cast them to their actual data type since they are stored as a generic interface without datatype specified.
    ```
    user := data.Get("user").(string)
    ```
- Another important part of the code is the *Data* field from the *req* variable passed as a parameter. This variable will contain any paraneter that was passed to the command that was not defined explicitly. In our example we only defined user as an expected input variable, so imagine that we call our command as.
    ```shell
    <command> user=Isaac city=London country=England
    ```
    In this case *req.Data["city"]* would contain "London" and *req.Data["country"]* would contain England.
    In the example we are analising *req.Data* is empty but we are reusing it as an intermidiate datastructure to store the address and get it marshalled.
    ```go
    req.Data["address"] = address
    bufAddr, err := json.Marshal(req.Data)
    if err != nil {
        return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
    }
    ```
    
- Next almost at the end of the function you can observe that we are saving both the generated address and key inside "b.store".
    ```go
        b.store[req.ClientToken+"/address/"+user] = bufAddr
        b.store[req.ClientToken+"/key/"+user] = bufKey
    ```
    This corresponds to the storage part of the backend that will persist between executions, the store is as simple as a mapping from string to bytes. This is very usefull because you can store any kind of information as long as you Marshall it before to get its bytes representation.
    
- Finally let's take a look at how to return anything that we want to show on the screen as an ouput to the user. First we need to initialize the response as.
    ```go
    resp := &logical.Response {
        Data: map[string]interface{}{},
            
    }
    ```
    After that we can insert insert any key/ value pair that we want. In this example we are saving under the key address the actual generated address that is a string.
    ```go
    resp.Data["address"] = address
    ```
    The whole structure will be printed as a table of key value to the end user.

Once you have your code ready you just need to compile, but we recommend to specify where the output of the compilation and save all your plugins under a well known folder.

### How to use your custom plugin once compiled
If you want to have a first idea of how to write a plugin you can take a look at the plugins that are present on this repo. They are under secrets, both EthereumPlugin and LRS.
Once you have writen and compiled all your plugins in a known folder all you need to do is run the following command to start testing them.
```
vault server -dev -dev-root-token-id=root -dev-plugin-dir=<your_plugins_folder>
```
As you can see we are running the vault in dev mode which means we dont need to unseal or seal it. 
At the same time we are setting the "dev-root-token" to "root" which means that we can get root permissions with the following command. Instead needing to remember a complex root token.
```
vault login root
```
Once we have already logged in, all we need to do is to enable our plugin secret engine with the following command.
```
vault secrets enable <plugin_name>
```
After this we can use our custom plugin normally as if it were part of the standard functionalities from vault.

### Register plugin for non-dev mode Vault
First you need to add an extra line on the config.hcl file from the vault at the begining.
```
plugin_directory = <your_plugin_directory>
```
After that you need to calculate the sha256 checksum of each of your plugins by executing.
On linux
```
sha256sum <your_plugin_directory>/<your_plugin>
```
On mac
```
shasum -a 256 <your_plugin_directory>/<your_plugin>
```
You will get the hash of the plugin so now you need to register it on the new running vault.
For that you need to execute
```
vault plugin register -sha256=<plugin_hash> <plugin_type> <plugin_name>
```
Where plugin type can be auth (for an authentication plugin) or secret (for a secret engine plugin)
After this you will have already registered the plugin on you vault and you just need to enable it by executing.
```
vault secrets enable -path=<plugin_name> LRS <plugin_name>
```
After that you can use the plugin as normal.