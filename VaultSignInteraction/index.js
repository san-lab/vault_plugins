const ETHEREUM_URL = 'https://rinkeby.infura.io/v3/f2a8581c640340758bead17199084148';
const VAULT_URL = 'http://127.0.0.1:8200';

const Web3 = require('web3');
const rp = require('request-promise');
const web3 = new Web3(ETHEREUM_URL);
const Tx = require('ethereumjs-tx').Transaction;
const username="pedro"

const ETH_NETWORK = 'rinkeby';

const getVault = (urlEncodedTx, user) => {
    return new Promise((resolve, reject) => {
        const options = {
            method: 'GET',
            headers: {
                'X-Vault-Token': 'root'
            },
            uri: `${VAULT_URL}/v1/ethereumPlugin/signTx?tx=${urlEncodedTx}&user=${user}`,
            json: true,
        };

        rp(options)
            .then(res => {
                console.log("Transaction successfully sent")
                resolve(res);
            })
            .catch(err => {
                console.error(err)
                reject();
            });
    });
};

const urlEncoder = (txObject) => {
    return JSON.stringify(txObject).split('').map(el => {
        switch (el) {
            case '{':
                return '%7B';
            case '"':
                return '%22';
            case ':':
                return '%3A';
            case ',':
                return '%2C';
            case '}':
                return '%7D';
            default:
                return el;
        }
    }).join('')
}

const getAddress = (user) => {
    return new Promise((resolve, reject) => {
        const options = {
            method: 'GET',
            headers: {
                'X-Vault-Token': 'root'
            },
            uri: `${VAULT_URL}/v1/ethereumPlugin/showAddr?user=${user}`,
            json: true,
        };

        rp(options)
            .then(res => {
                console.log("Transaction successfully sent")
                resolve(res);
            })
            .catch(err => {
                console.error(err)
                reject();
            });
    });
};

(async function () {

    const accountAddress = await getAddress(username);
    const tnonce = await web3.eth.getTransactionCount(accountAddress.data.address)

    const tnonceHex = `0x${tnonce.toString(16)}`;

    const txParams = {
        nonce: tnonceHex,
        gasPrice: '0x3B9ACA00',
        gas: '0x7530',
        to: '0xf17f52151EbEF6C7334FAD080c5704D77216b732',
        value: '0x1',
        input: '0x',
        v: '',
        r: '',
        s: ''
    }

    const urlEncodedTx = urlEncoder(txParams)

    console.log(`URL encoded transaction: ${urlEncodedTx}`);

    try {
        const signedTx = await getVault(urlEncodedTx, username);

        console.log(`Transaction with signature: ${signedTx.data.result}`)

        const tx = new Tx(JSON.parse(signedTx.data.result), { 'chain': ETH_NETWORK });
        const serializedTx = tx.serialize();
        const res = await web3.eth.sendSignedTransaction('0x' + serializedTx.toString('hex'));
        console.log("View on Etherscan:");
        console.log(`https://rinkeby.etherscan.io/tx/${res.transactionHash}`);
    }
    catch (err) {
        console.error(err);
    }

})();