## Steps to test the plugin
```
cd vault-guides/secrets/mock
go build -o vault/plugins/signTx cmd/mock/main.go
```

Start vault on a different terminal with 
```
vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins
```
Go back to the first terminal
```
export VAULT_ADDR="http://127.0.0.1:8200"
vault login root
vault secrets enable signTx
vault policy write pedro Vault_profiles/user1.hcl 
vault policy write guillermo Vault_profiles/user2.hcl 
vault policy write przemek Vault_profiles/user3.hcl
```

Now you can try to log in as different users and see that each of them can only acces the key they insert.
```
vault token create -policy=pedro (to create a login token)
vault login <token>
vault write signTx/ethKeypedro ethKey="0xC87509A1C067BBDE78BEB793E6FA76530B6382A4C0241E5E4A9EC0A0F44DC0D3"
vault read signTx/ethKeypedro
```


Now we will try to log in as other user for example przemek and we will see that it cannot access the key from pedro
First we need root privileges to ask for a token for przemek profile
```
vault login root
```
```
vault token create -policy=przemek
vault login <token>
vault read signTx/ethKeypedro
```
z