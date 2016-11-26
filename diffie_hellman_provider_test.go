package doubleratchet

type stubDiffieHellmanProvider struct{}

func (stub stubDiffieHellmanProvider) Generate() DiffieHellmanPair {
	return nil
}

func (stub stubDiffieHellmanProvider) DeffieHellmanOutput(dhPair DiffieHellmanPair, dhPub DiffieHellmanPublic) []byte {
	return nil
}

func (stub stubDiffieHellmanProvider) KdfRootKey(rootKey, dhOut []byte) ([]byte, []byte) {
	return nil, nil
}

func (stub stubDiffieHellmanProvider) KdfChainKey(chainKey []byte) ([]byte, []byte) { return nil, nil }

func (stub stubDiffieHellmanProvider) Encrypt(messageKey, plaintext, associatedData []byte) []byte {
	return nil
}

func (stub stubDiffieHellmanProvider) Decrypt(messageKey, ciphertext, associatedData []byte) []byte {
	return nil
}

func (stub stubDiffieHellmanProvider) Concat(assoicatedData []byte, header Header) []byte { return nil }
