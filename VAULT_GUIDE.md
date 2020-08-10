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

## How to use your custom plugin once compiled
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

## Register plugin for non-dev mode Vault
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