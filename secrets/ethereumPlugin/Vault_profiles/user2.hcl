path "ethereumPlugin/genKey" {
  capabilities = [ "create", "update"]
}

path "ethereumPlugin/showAddr/guille" {
  capabilities = ["read"]
}

path "ethereumPlugin/signTx/guille" {
  capabilities = [ "create", "update"]
}
