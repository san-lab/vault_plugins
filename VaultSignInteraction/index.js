const Web3 = require('web3');
const rp = require('request-promise');
const web3 = new Web3('https://rinkeby.infura.io/v3/f2a8581c640340758bead17199084148');
var Tx = require('ethereumjs-tx').Transaction;
const queryString = require('query-string');


(async function () {

    const privateKey = Buffer.from(
        '2659d295cf455bc033e5b5ec59afc67057425af8a71a694a5f59ad0e6b333f0c',
        'hex',
    )

    const tnonce = await web3.eth.getTransactionCount('0xc8dfCA661A53bC05EC1BC76d20Ba77C34F8facAb')
    const tnonceHex = `0x${tnonce.toString(16)}`;


    const txParams = {
        nonce: tnonceHex,
        gasPrice: '0x161E70F600',
        gas: '0x989680',
        to: '0xf17f52151EbEF6C7334FAD080c5704D77216b732',
        value: '0x1',
        input: '0x',
        v: '',
        r: '',
        s: ''
    }

    const txParams2 = {
        nonce: tnonceHex,
        gasPrice: '0x161E70F600',
        gas: '0x989680',
        to: '0xf17f52151EbEF6C7334FAD080c5704D77216b732',
        value: '0x1',
        input: '0x',
    }


    // The second parameter is not necessary if these values are used
    tx2 = new Tx(txParams2, { 'chain': 'rinkeby' })
    tx2.sign(privateKey)
    const serializedTx2 = tx2.serialize();

    console.log(tx2.toJSON())



    const urlEncodedTx = urlEncoder(txParams)

    console.log(urlEncodedTx);

    try {
        const foo = await getVault(urlEncodedTx);
        console.log(foo.data.result);
        const tx = new Tx(JSON.parse(foo.data.result), { 'chain': 'rinkeby' });
        console.log(tx.toJSON())
        const serializedTx = tx.serialize();
        const res = await web3.eth.sendSignedTransaction('0x' + serializedTx.toString('hex'));
        console.log(res);
    }
    catch (err) {
        console.error(err);
    }

})();

const getVault = (urlEncodedTx) => {
    return new Promise((resolve, reject) => {
        const options = {
            method: 'GET',
            headers: {
                'X-Vault-Token': 'root'
            },
            uri: `http://127.0.0.1:8200/v1/signTx/ethKeypedro?tx=${urlEncodedTx}`,
            json: true,
        };

        rp(options)
            .then(res => {
                console.log("Bien")
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