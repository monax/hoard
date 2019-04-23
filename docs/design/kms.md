# Key Management

Hoard acts much like a conventional password store; we symmetrically encrypt the given data with it's original hash then store it at an address determined by the hash of the encrypted data. Compare this with a typical setup where a user's password is hashed and authentication is based on whether the user can prove knowledge of the hash; both systems require you to have some cognizance of the plaintext input.

An obvious extension would be to build an integration to [Hashicorp Vault](https://www.vaultproject.io/), where Hoard would act as a back-end encrypted key-value store. However, it may also be possible to use it in isolation...

With symmetric grants we are able to encrypt this reference in a closed eco-system with a password like derivative. This means that we can readily share this object with all internal actors in plain sight of external entities. Alternatively we can explicitly share it with a specified party using an asymmetric grant - which locks the original reference with their public key. Though we may consider defining this more formally in the future, this gives us one form of access control which allows us to securely share keys.

## Kubernetes

It is typical to manage a wide range of secrets in a typical Kubernetes environment, from API keys to cloud credentials and database logins. Bitnami's [Sealed-Secrets](https://github.com/bitnami-labs/sealed-secrets) addresses this issue; "I can manage all my K8s config in git, except Secrets.". This system essentially allows an operator to encrypt a secret with the server's public key and commit it to version control. When the Custom Resource Definition (CRD) is created in the cluster, the server will decrypt it in-place. With access grants, we could easily achieve a similar level of functionality with Hoard. 
