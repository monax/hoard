'use strict'

const path = require('path')

const LOCAL_PATH = path.join(__dirname, './protobuf');
const PROTO_PATH = path.join(__dirname, '../protobuf');
const PROTO_FILE = 'services.proto';

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
    const hoard_proto = grpc.loadPackageDefinition(packageDefinition).services;
    this.cleartextClient = new hoard_proto.Cleartext(address,
        grpc.credentials.createInsecure());
    this.encryptionClient = new hoard_proto.Encryption(address,
        grpc.credentials.createInsecure());
    this.storageClient = new hoard_proto.Storage(address,
        grpc.credentials.createInsecure());
    this.grantClient = new hoard_proto.Grant(address,
        grpc.credentials.createInsecure());
};

HoardClient.prototype.get = function (reference) {
    const client = this.cleartextClient;
    return new Promise(function (resolve, reject) {
        client.get(reference, function (err, plaintext) {
            if (err) {
                reject(err);
            } else {
                resolve(plaintext);
            }
        });
    });
};

HoardClient.prototype.put = function (plaintext) {
    const client = this.cleartextClient;
    return new Promise(function (resolve, reject) {
        client.put(plaintext, function (err, reference) {
            if (err) {
                reject(err);
            } else {
                resolve(reference);
            }
        });
    });
};

HoardClient.prototype.unsealget = function (grant) {
    const client = this.grantClient;
    return new Promise(function (resolve, reject) {
        client.unsealGet(grant, function (err, plaintext) {
            if (err) {
                reject(err);
            } else {
                resolve(plaintext);
            }
        });
    });
};

HoardClient.prototype.putseal = function (plaintextAndGrantSpec) {
    const client = this.grantClient;
    return new Promise(function (resolve, reject) {
        client.putSeal(plaintextAndGrantSpec, function (err, grant) {
            if (err) {
                reject(err);
            } else {
                resolve(grant);
            }
        });
    });
};

HoardClient.prototype.encrypt = function (plaintext) {
    const client = this.encryptionClient;
    return new Promise(function (resolve, reject) {
        client.encrypt(plaintext, function (err, referenceAndCiphertext) {
            if (err) {
                reject(err);
            } else {
                resolve(referenceAndCiphertext);
            }
        });
    });
};

HoardClient.prototype.decrypt = function (referenceAndCiphertext) {
    const client = this.encryptionClient;
    return new Promise(function (resolve, reject) {
        client.decrypt(referenceAndCiphertext, function (err, plaintext) {
            if (err) {
                reject(err);
            } else {
                resolve(plaintext);
            }
        });
    });
};

HoardClient.prototype.push = function (ciphertext) {
    const client = this.storageClient;
    return new Promise(function (resolve, reject) {
        client.push(ciphertext, function (err, address) {
            if (err) {
                reject(err);
            } else {
                resolve(address);
            }
        });
    });
};

HoardClient.prototype.pull = function (address) {
    const client = this.storageClient;
    return new Promise(function (resolve, reject) {
        client.pull(address, function (err, ciphertext) {
            if (err) {
                reject(err);
            } else {
                resolve(ciphertext);
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

HoardClientDynamic.prototype = Object.create(HoardClient.prototype);
module.exports.Client = HoardClientDynamic;