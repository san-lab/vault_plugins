const VAULT_URL = 'http://127.0.0.1:8200';

const rp = require('request-promise');

const ETH_ACCOUNT = '0xc8dfCA661A53bC05EC1BC76d20Ba77C34F8facAb';
const ETH_PRIV_KEY = '0x2659d295cf455bc033e5b5ec59afc67057425af8a71a694a5f59ad0e6b333f0c';

const setPrivateKey = () => {
    return new Promise((resolve, reject) => {
        const options = {
            method: 'POST',
            headers: {
                'X-Vault-Token': 'root'
            },
            body: {
                ethKey: ETH_PRIV_KEY
            },
            uri: `${VAULT_URL}/v1/signTx/ethKeypedro`,
            json: true,
        };

        rp(options)
            .then(res => {
                console.log("Private key set")
                resolve(res);
            })
            .catch(err => {
                console.error(err)
                reject();
            });
    });
};

(async function () {
    try {
        const response = await setPrivateKey();
        console.log(response);
    }
    catch (err) {
        console.error(err);
    }

})();