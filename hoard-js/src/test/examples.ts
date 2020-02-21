import example = require('../examples');

describe('Examples should run', function () {
  it('buffer based plaintext', async function () {
    const data = Buffer.from('some stuff', 'utf8');
    const salt = Buffer.from('foo', 'ascii');
    await example.example(data, salt);
  })
  it('base64 based plaintext', async function () {
    const data = Buffer.from('some stuff', 'utf8').toString('base64');
    const salt = Buffer.from('foo', 'ascii').toString('base64');
    await example.example(data, salt);
  })
})
