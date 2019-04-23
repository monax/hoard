# Hoard

## Properties
- Resistance to chosen plaintext (AES with one-time key SHA-256 of content)
- Resistance to chosen ciphertext (GCM over AES)
- Saltable (GCM additional data)
- Confidentiality, integrity
- Authenticity (ciphertext came from a party actually holding key and therefore plaintext)

## Encryption scheme
Hoard implements an encryption scheme based off the SHA256 cryptographic hash function and the symmetric block cipher AES256-GCM (Galois Counter Mode is an authenticated mode of AES). It is an example of envelope encryption where an object is encrypted with a specific one-time key and where that secret key can itself be shared by encrypting it (asymmetrically or otherwise) and publishing it to a recipient. It is motivated by and possesses the following features:

- Usable on a publicly-accessible storage backend (strongly encrypted by AES256-GCM)
- Addressable by the hash of the ciphertext (the bytes of encrypted object) so can be ciphertext can de-duplicated and object existence can be queried without the secret key
- Ciphertext is recoverable from plaintext (if you have a copy of the plaintext you can check whether the ciphertext is stored and recover the secret key)
- Permits sharing of access grants in public via asymmetric encryption of the secret key and address (can be used to implement a decentralised equivalent of access control lists)
- Optional salt allows one to simulate random encryption keys or to add entropy to short objects

The encryption scheme relies on well-known cryptographic primitives. Given an object (just a sequence of bytes) and a salt the encryption proceeds as follows:

1. Compute the SHA256 of the object.
2. Prepend the salt (just arbitrary bytes) to the object.
3. Encrypt the salted object with AES256-GCM with an empty IV (this is okay since our keys are one-time since they are based on a semantically secure hash of the object) with authenticated additional data describing the presence of the salt (authenticated additional data is a feature of GCM).
4. Compute the SHA256 of the encrypted object.
5. Return the output of (1) as the `secretKey`, the output of (3) as the `encryptedData`, and the output of (4) as the `address`.

Given an encrypted object, secret key, and salt the decryption proceeds as follows:

1. Decrypt the encrypted object with AES256-GCM with an empty IV with `secretKey` and additional authenticated data associated with `salt`.
2. Trim the prefix `salt` from the output of (1).
3. Return the object as `data` from the output of (2).


### Security
By design this scheme is trivially vulnerable to known-plaintext attacks (if you know the plaintext you can find the key).

It is also vulnerable to a ciphertext-only attack where a rainbow table of ciphertexts may be computed where the space of variation in a set of plaintexts is small. This occurs when the plaintexts are short or when the plaintexts are mostly identical template with a small amount of variation (an 8-digit account number for example). This is only a problem if revealing that an object has been stored at all leaks sensitive information.

The scheme should not be vulnerable to general chosen-ciphertext attacks (CCA) through the use of Galois Counter Mode with AES that is CCA-secure.

If you want known-plaintext or ciphertext-only security you can provide a salt. The salt can be of any length and can be used a number of ways. If you are encrypting short plaintexts you can agree a sufficiently long salt to be pre-shared amongst a set of parties that adds entropy to the encryption so that a rainbow table of possible values cannot be feasibly built (semantic security of SHA256 means the salt will induce an unpredictable variation of ciphertexts). In this case amongst the parties in possession of the salt the advantages of deterministic encryption are preserved. Alternatively a random hash can be used for the salt, this will induce a different secret key and address each time you encrypt the same bytes. It is effectively the same as using a random key.

The encryption is furthermore vulnerable to the same timing and length attacks that to which AES is susceptible, but for most purposes these attacks are not usually considered an issue.

### Maturity
Hoard is still young, but the API should be mostly stable. The cryptographic libraries used are standard Go libraries (and Go's NACL implementation) so should be of reasonable quality and are widely deployed. The encryption scheme is straight-forward and has an isolated implementation. However there may be bugs in the implementation.