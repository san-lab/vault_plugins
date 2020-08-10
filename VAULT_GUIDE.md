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

## How to use your own custom plugin