const ETHEREUM_URL = 'https://rinkeby.infura.io/v3/f2a8581c640340758bead17199084148';
const VAULT_URL = 'http://127.0.0.1:8200';

const Web3 = require('web3');
const rp = require('request-promise');
const web3 = new Web3(ETHEREUM_URL);
const Tx = require('ethereumjs-tx').Transaction;
const ETH_PRIV_KEY = '2659d295cf455bc033e5b5ec59afc67057425af8a71a694a5f59ad0e6b333f0c';
const ETH_ACCOUNT = '0xc8dfCA661A53bC05EC1BC76d20Ba77C34F8facAb';
const ETH_NETWORK = 'rinkeby';
const username = "pedro";


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

    const privateKey = Buffer.from(
        ETH_PRIV_KEY,
        'hex',
    )



    const tnonce = await web3.eth.getTransactionCount(ETH_ACCOUNT)
    const tnonceHex = `0x${tnonce.toString(16)}`;

    const txParams = {
        nonce: tnonceHex,
        gasPrice: '0x2540BE400',
        gas: '0x7530',
        to: null,
        value: '0x16345785D8A0000',
        input: '0x',
    }

    try {
        const addressToBeFunded = await getAddress(username);

        console.log(`Address to be funded: ${addressToBeFunded.data.address}`)

        txParams.to = addressToBeFunded.data.address;

        const tx = new Tx(txParams, { 'chain': ETH_NETWORK });
        tx.sign(privateKey);
        const serializedTx = tx.serialize();
        const res = await web3.eth.sendSignedTransaction('0x' + serializedTx.toString('hex'));
        console.log("View on Etherscan:");
        console.log(`https://rinkeby.etherscan.io/tx/${res.transactionHash}`);
    }
    catch (err) {
        console.error(err);
    }

})();