const example = require("../examples")

describe('Examples should run', function () {
    it('buffer based plaintext', async function () {
        let plaintext = {
            Data: Buffer.from('some stuff', 'utf8'),
            Salt: Buffer.from('foo', 'ascii')
        };
        let result = await example.example(plaintext)
        console.log(result)
    });
    it('base64 based plaintext', async function () {
        let plaintext = {
            Data: Buffer.from('some stuff', 'utf8').toString('base64'),
            Salt: Buffer.from('foo', 'ascii').toString('base64')
        };
        let result = await example.example(plaintext)
        console.log(result)
    });
});