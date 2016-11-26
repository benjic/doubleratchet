package doubleratchet

import "errors"

var errPersistorFull = errors.New("Persistor has reached its max size")

const maxSkippedKeys int = 10

// A skippedStore allows the double ratchet to maintain a set of message
// keys that have not been encountered.
type skipPersistor interface {
	Get(DiffieHellmanPublic, int) []byte
	Put(DiffieHellmanPublic, int, []byte) error
}

type mappedSkipPersistor map[DiffieHellmanPublic]map[int][]byte

func (m mappedSkipPersistor) Get(dh DiffieHellmanPublic, n int) []byte {
	dhMap := m[dh]

	if dhMap == nil {
		return nil
	}

	mk := dhMap[n]

	delete(m[dh], n)

	if len(m[dh]) == 0 {
		delete(m, dh)
	}

	return mk
}

func (m mappedSkipPersistor) Put(dh DiffieHellmanPublic, n int, mk []byte) error {

	// Prevent a baddie from flooding the missed key store.
	//
	// TODO: It may be better to make a circular buffer so the store isn't
	// plugged up forever with baddie turds.

	if len(m) == maxSkippedKeys {
		return errPersistorFull
	}

	if m[dh] == nil {
		m[dh] = make(map[int][]byte, 0)
	}

	if len(m[dh]) == maxSkippedKeys {
		return errPersistorFull
	}

	m[dh][n] = mk
	return nil
}
