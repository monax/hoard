'use strict'

const path = require('path')

const LOCAL_PATH = path.join(__dirname, './protobuf');
const PROTO_PATH = path.join(__dirname, '../protobuf');
const PROTO_FILE = 'api.proto';

const protoLoader = require('@grpc/proto-loader');
const grpc = require('grpc')

const options = {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
    includeDirs: [
        LOCAL_PATH, PROTO_PATH
    ]
}

const packageDefinition = protoLoader.loadSync(PROTO_FILE, options);

const HoardClient = function () {};

// HoardClient uses dynamic types and mapping
const HoardClientDynamic = function (address) {
    const hoard_proto = grpc.loadPackageDefinition(packageDefinition).api;
    this.cleartextClient = new hoard_proto.Cleartext(address,
        grpc.credentials.createInsecure());
    this.encryptionClient = new hoard_proto.Encryption(address,
        grpc.credentials.createInsecure());
    this.storageClient = new hoard_proto.Storage(address,
        grpc.credentials.createInsecure());
    this.grantClient = new hoard_proto.Grant(address,
        grpc.credentials.createInsecure());
    this.chunkSize = 64 * 1024;
};

HoardClient.prototype.get = function (reference) {
    const client = this.cleartextClient;
    return new Promise(function (resolve, reject) {
        const call = client.get(reference);
        getPlaintext(call, resolve, reject);
    });
};

HoardClient.prototype.put = function (plaintext) {
    const client = this.cleartextClient;
    const size = this.chunkSize;
    return new Promise(function (resolve, reject) {
        const call = client.put(function(error, reference) {
            if (error) {
                reject(error);
            } else {
                resolve(reference);
            }
        });
        putPlaintext(call, plaintext, size);
        call.end();
    });
};

HoardClient.prototype.unsealget = function (grant) {
    const client = this.grantClient;
    return new Promise(function (resolve, reject) {
        const call = client.unsealGet(grant);
        getPlaintext(call, resolve, reject);
    });
};

HoardClient.prototype.unsealdelete = function (grant) {
    const client = this.grantClient;
    return new Promise(function (resolve, reject) {
        client.unsealDelete(grant, function (err, address) {
            if (err) {
                reject(err);
            } else {
                resolve(address);
            }
        });
    });
};

HoardClient.prototype.putseal = function (plaintextAndGrantSpec) {
    const client = this.grantClient;
    const size = this.chunkSize;
    return new Promise(function (resolve, reject) {
        const call = client.putSeal(function(error, grant) {
            if (error) {
                reject(error);
            } else {
                resolve(grant);
            }
        });

        const spec = plaintextAndGrantSpec.GrantSpec;
        const salt = plaintextAndGrantSpec.Plaintext.Salt;
        const data = plaintextAndGrantSpec.Plaintext.Data;

        call.write({GrantSpec: spec});
        call.write({Plaintext: {Salt: salt}});
        for (var i=0; i<data.length; i+=size) {
            if (i+size>data.length) {
                call.write({Plaintext: {Data: data.slice(i, data.length)}});
            } else {
                call.write({Plaintext: {Data: data.slice(i, i+size)}});
            }
        }
        call.end();
    });
};

HoardClient.prototype.encrypt = function (plaintext) {
    const client = this.encryptionClient;
    const size = this.chunkSize;
    return new Promise(function (resolve, reject) {
        const call = client.encrypt(function(error, referenceAndCiphertext) {
            if (error) {
                reject(error);
            } else {
                resolve(referenceAndCiphertext);
            }
        });
        putPlaintext(call, plaintext, size);
        call.end();
    });
};

HoardClient.prototype.decrypt = function (referenceAndCiphertext) {
    const client = this.encryptionClient;
    return new Promise(function (resolve, reject) {
        const call = client.decrypt(referenceAndCiphertext);
        getPlaintext(call, resolve, reject);
    });
};

HoardClient.prototype.push = function (ciphertext) {
    const client = this.storageClient;
    const size = this.chunkSize;
    return new Promise(function (resolve, reject) {
        const call = client.push(function(error, address) {
            if (error) {
                reject(error);
            } else {
                resolve(address);
            }
        });
        const data = ciphertext.EncryptedData;
        for (var i=0; i<data.length; i+=size) {
            if (i+size>data.length) {
                call.write({EncryptedData: data.slice(i, data.length)});
            } else {
                call.write({EncryptedData: data.slice(i, i+size)});
            }
        }
        call.end();
    });
};

HoardClient.prototype.pull = function (address) {
    const client = this.storageClient;
    return new Promise(function (resolve, reject) {
        const call = client.pull(address);

        var ciphertext = {EncryptedData: Buffer.alloc(0)};
        call.on('data', function(data) {
            ciphertext.EncryptedData = Buffer.concat([ciphertext.EncryptedData, data.EncryptedData]);
        });

        call.on('error', function(e) {
            reject(e);
        });

        call.on('end', function() {
            resolve(ciphertext);
        });
    });
};

HoardClient.prototype.delete = function (address) {
    const client = this.storageClient;
    return new Promise(function (resolve, reject) {
        client.delete(address, function (err, address) {
            if (err) {
                reject(err);
            } else {
                resolve(address);
            }
        });
    });
};

HoardClient.prototype.stat = function (address) {
    const client = this.storageClient;
    return new Promise(function (resolve, reject) {
        client.stat(address, function (err, statInfo) {
            if (err) {
                reject(err);
            } else {
                resolve(statInfo);
            }
        });
    });
};

// Walk over the given object and base64 encode any buffers
HoardClient.prototype.base64ify = function (obj) {
    let newObj = {};
    for (let key of Object.keys(obj)) {
        let value = obj[key];
        if (value instanceof Buffer) {
            newObj[key] = value.toString(`base64`)
        } else {
            newObj[key] = value
        }
    }
    return newObj;
};

function putPlaintext(call, plaintext, size) {
    call.write({Salt: plaintext.Salt})
    const data = plaintext.Data;
    for (var i=0; i<data.length; i+=size) {
        if (i+size>data.length) {
            call.write({Data: data.slice(i, data.length)});
        } else {
            call.write({Data: data.slice(i, i+size)});
        }
    }
}

function getPlaintext(call, resolve, reject) {
    var plaintext = {Data: Buffer.alloc(0), Salt: Buffer.alloc(0)};
    call.on('data', function(data) {
        if (data.input == 'Salt') {
            plaintext.Salt = Buffer.concat([plaintext.Salt, data.Salt]);
        } else if (data.input == 'Data') {
            plaintext.Data = Buffer.concat([plaintext.Data, data.Data]);
        }
    });

    call.on('error', function(e) {
        reject(e);
    });

    call.on('end', function() {
        resolve(plaintext);
    });
}

HoardClientDynamic.prototype = Object.create(HoardClient.prototype);
module.exports.Client = HoardClientDynamic;
