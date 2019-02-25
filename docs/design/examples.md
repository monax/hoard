# Examples

To configure hoard with a fully-functioning secret back-end use the following example config:

```toml
ListenAddress = "tcp://:53431"

[Storage]
  StorageType = "memory"
  AddressEncoding = "base64"

[Logging]
  LoggingType = "json"
  Channels = ["info", "trace"]

[Secrets.OpenPGP]
  ID = "10449759736975846181"
  File = "${HOME}/go/src/github.com/monax/hoard/grant/private.key.asc"

[[Secrets.Symmetric]]
  ID = "test"
  Passphrase = "test"
```

Remember to change the provided values before deploying to production.