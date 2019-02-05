## Access Grants
Access grants are encrypted envelopes containing a content address and secret key (i.e. `reference.Ref`). They are intended to provide a means of sharing access to the underlying data.

A `GrantSpec` describes how the encryption is performed - its content depends on the type of grant but the essential content is the identifier of some secret to which the Hoard daemon `hoard` has access. This may be direct access to secret key matter as in a symmetric grant or via a proxy such as the system's gpg-agent or a remote secret store.

### Grant types
Below are some planned Hoard grants types (not all are implemented).

#### Symmetric
This is the first grant type to be implemented. The grants are encrypted with AES256-GCM using

#### OpenPGP (WIP)
Use local system keys to asymmetrically encrypt grants

#### Vault (WIP)

#### ECDH (WIP)
Based on curve25519 keys/station-to-station protocol.

### Metadata
A Grant contains sufficient metadata to allow for its decryption in the form of a contained GrantSpec.

### Secrets (WIP)

The GrantSpec may entail accessing a remote/local secret store outside of Hoard. We could support self-hosted secrets in:

- Hoard config
- In-memory/provided via environment variables
- Compiled-in secrets
- Hoard document accessible via secret embedded by any of the above
- Flat file