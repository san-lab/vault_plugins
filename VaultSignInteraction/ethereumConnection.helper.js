const ethers = require('ethers');
const rpcUrl = 'https://rinkeby.infura.io/v3/f2a8581c640340758bead17199084148'
const mainPrivateKey = '0x2659d295cf455bc033e5b5ec59afc67057425af8a71a694a5f59ad0e6b333f0c';
const provider = new ethers.providers.JsonRpcProvider(rpcUrl);
const wallet = new ethers.Wallet(mainPrivateKey, provider);


module.exports = {
    wallet
} 