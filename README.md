## Steps to test the plugin
```
cd vault-guides/secrets/ethereumPlugin
go build -o vault/plugins/ethereumPlugin cmd/ethereumPlugin/main.go
```

Start vault on a different terminal with 
```
vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins
```
Go back to the first terminal
```
export VAULT_ADDR="http://127.0.0.1:8200"
vault login root
vault secrets enable ethereumPlugin
vault policy write pedro Vault_profiles/user1.hcl 
vault policy write guillermo Vault_profiles/user2.hcl 
vault policy write przemek Vault_profiles/user3.hcl
```

Now you can try to log in as different users and see that each of them can only acces the key they insert.
```
vault token create -policy=pedro (to create a login token)
vault login <token>
vault write ethereumPlugin/genKey user=pedro
vault read ethereumPlugin/showAddr user=pedro
vault read ethereumPlugin/signTx user=pedro tx="{\"nonce\":\"0x33\",\"gasPrice\":\"0x0\",\"gas\":\"0x989680\",\"to\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"value\":\"0xbd3580\",\"input\":\"0x\",\"v\":\"\",\"r\":\"\",\"s\":\"\",\"hash\":\"0xc0bacd35d3ea25a130696336dd6b1d811e9f5defdeb28530d0222b7ff2c979cb\"}"
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

If you want to have a look at the backend code its under /vault-guides/secrets/mock/backend.go


WIP: Command to call the plugin using http

curl -H "X-Vault-Token: root" -X GET  http://127.0.0.1:8200/v1/signTx/ethKeypedro?tx=%7B%22nonce%22%3A%220x33%22%2C%22gasPrice%22%3A%220x0%22%2C%22gas%22%3A%220x989680%22%2C%22to%22%3A%220x627306090abab3a6e1400e9345bc60c78a8bef57%22%2C%22value%22%3A%220xbd3580%22%2C%22input%22%3A%220x%22%2C%22v%22%3A%22%22%2C%22r%22%3A%22%22%2C%22s%22%3A%22%22%2C%22hash%22%3A%220xc0bacd35d3ea25a130696336dd6b1d811e9f5defdeb28530d0222b7ff2c979cb%22%7D

Command using CLI

vault read signTx/ethKeypedro tx="{\"nonce\":\"0x33\",\"gasPrice\":\"0x0\",\"gas\":\"0x989680\",\"to\":\"0x627306090abab3a6e1400e9345bc60c78a8bef57\",\"value\":\"0xbd3580\",\"input\":\"0x\",\"v\":\"\",\"r\":\"\",\"s\":\"\",\"hash\":\"0xc0bacd35d3ea25a130696336dd6b1d811e9f5defdeb28530d0222b7ff2c979cb\"}"
