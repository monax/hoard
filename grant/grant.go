package grant

type Grant struct {
	GrantSpec
	EncryptedReference []byte
}
