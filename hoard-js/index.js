const PROTO_PATH = __dirname + '/hoard.proto';

const GRPC = require('grpc')

const HoardClient = function () {
};

// HoardClient using statically generated javascript types
const HoardClientStatic = function (address) {
    const services = require('./protobuf/hoard_grpc_pb');
    this.cleartextClient = new services.CleartextClient(address,
        GRPC.credentials.createInsecure());
    this.encryptionClient = new services.EncryptionClient(address,
        GRPC.credentials.createInsecure());
    this.storageClient = new services.StorageClient(address,
        GRPC.credentials.createInsecure());
};

// HoardClient using dynamically types and mapping (the default)
const HoardClientDynamic = function (address) {
    const hoard_proto = GRPC.load(PROTO_PATH).core;
    this.cleartextClient = new hoard_proto.Cleartext(address,
        GRPC.credentials.createInsecure());
    this.encryptionClient = new hoard_proto.Encryption(address,
        GRPC.credentials.createInsecure());
    this.storageClient = new hoard_proto.Storage(address,
        GRPC.credentials.createInsecure());
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

HoardClientStatic.prototype = Object.create(HoardClient.prototype);
HoardClientDynamic.prototype = Object.create(HoardClient.prototype);

module.exports.Client = HoardClientDynamic;
module.exports.DynamicClient = HoardClientDynamic;
module.exports.StaticClient = HoardClientStatic;
