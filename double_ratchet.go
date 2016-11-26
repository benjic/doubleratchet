package doubleratchet

// A DoubleRatchet provides the state nessecary to maintain a double ratchet
// algorithm.
type DoubleRatchet struct {
	dhs      DiffieHellmanPair
	dhr      DiffieHellmanPublic
	rk       []byte
	cks      []byte
	ckr      []byte
	ns       int
	nr       int
	pn       int
	skipped  skipPersistor
	provider DiffieHellmanProvider
}

// Init configures a DoubleRatchet state
func Init(dhr DiffieHellmanPublic, sharedKey []byte, provider DiffieHellmanProvider) *DoubleRatchet {
	dhs := provider.Generate()
	do := provider.DeffieHellmanOutput(dhs, dhr)
	rk, cks := provider.KdfRootKey(sharedKey, do)

	return &DoubleRatchet{
		provider: provider,
		cks:      cks,
		dhs:      dhs,
		dhr:      dhr,
		rk:       rk,
		skipped:  make(mappedSkipPersistor, 0),
	}
}

// Encrypt performs internal ratchetting and produces a ciphertext for the
// given plaintext.
func (d *DoubleRatchet) Encrypt(plaintext []byte, ad []byte) (Header, []byte) {
	cks, mk := d.provider.KdfChainKey(d.cks)
	d.ns++
	d.cks = cks
	header := Header{Key: d.dhs, Pn: d.pn, N: d.ns}
	ciphertext := d.provider.Encrypt(mk, plaintext, d.provider.Concat(ad, header))

	return header, ciphertext
}

// Decrypt perfroms internall ratchetting and produces plain text from a given
// ciphertext.
func (d *DoubleRatchet) Decrypt(header Header, ciphertext []byte, ad []byte) ([]byte, error) {
	var mk []byte

	if mk = d.skipped.Get(header.Key, header.N); mk != nil {
		return d.provider.Decrypt(mk, ciphertext, ad), nil
	}

	if header.Key != d.dhr {
		if err := d.skipMessageKeys(header.Pn); err != nil {
			return nil, err
		}

		d.pn = d.ns
		d.ns = 0
		d.nr = 0
		d.dhr = header.Key

		do := d.provider.DeffieHellmanOutput(d.dhs, d.dhr)
		d.rk, d.ckr = d.provider.KdfRootKey(d.rk, do)
		d.dhs = d.provider.Generate()
		do = d.provider.DeffieHellmanOutput(d.dhs, d.dhr)
		d.rk, d.ckr = d.provider.KdfRootKey(d.rk, do)
	}

	if err := d.skipMessageKeys(header.N); err != nil {
		return nil, err
	}

	d.ckr, mk = d.provider.KdfChainKey(d.ckr)
	d.nr++

	return d.provider.Decrypt(mk, ciphertext, ad), nil
}

func (d *DoubleRatchet) skipMessageKeys(until int) error {
	if d.ckr != nil {
		for ; d.nr < until; d.nr++ {
			var mk []byte

			d.ckr, mk = d.provider.KdfChainKey(d.ckr)
			err := d.skipped.Put(d.dhr, d.nr, mk)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
