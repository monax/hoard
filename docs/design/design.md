TODO:
- envelope encryption
- access grants

## Hoard 

- Resistance to chosen plaintext (AES with one-time key SHA-256 of content)
- Resistance to chosen ciphertext (GCM over AES)
- Saltable (GCM additional data)
- Confidentiality, integrity
- Authenticity (ciphertext came from a party actually holding key and therfore plaintext)

