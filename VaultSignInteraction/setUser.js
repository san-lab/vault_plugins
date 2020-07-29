const VAULT_URL = 'http://127.0.0.1:8200';

const rp = require('request-promise');

const setPrivateKey = () => {
    return new Promise((resolve, reject) => {
        const options = {
            method: 'POST',
            headers: {
                'X-Vault-Token': 'root'
            },
            body: {
                user: 'user'
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