package doubleratchet

type DiffieHellmanPair interface{}

type DiffieHellmanPublic interface{}

type Header struct {
	Key DiffieHellmanPublic
	Pn  int
	N   int
}

// The DiffieHellmanProvider provides an abstraction over the cryptographic
// functions that facilitate the double ratchet implementation.
type DiffieHellmanProvider interface {

	// Generate produces a new DiffieHellmanPair
	Generate() DiffieHellmanPair

	// The DeffieHellmanOutput returns the result of a DiffieHellman exchange
	DeffieHellmanOutput(dhPair DiffieHellmanPair, dhPub DiffieHellmanPublic) []byte

	// The KdfRootKey function ratchets the root KDF chain.
	KdfRootKey(rootKey, dhOut []byte) (newRootKey, chainKey []byte)

	// The KdfChainKey ratchets the sending or reciving KDF chain.
	KdfChainKey(chainKey []byte) (newChainKey, messageKey []byte)

	// Encrypt returns the result of a cryptographic encryption of the given
	// plaintext via the given key.
	Encrypt(messageKey, plaintext, associatedData []byte) (ciphertext []byte)

	// Decrypt returns the result of a cryptographic decryption of the given
	// ciphertext via the given key.
	Decrypt(messageKey, ciphertext, associatedData []byte) (plaintext []byte)

	// Concat prepends the given associated data to a byte encoding of the
	// given header.
	//
	// TODO: This may be useless on the diffieHellmanProvider interface.
	Concat(associatedData []byte, header Header) []byte
}
