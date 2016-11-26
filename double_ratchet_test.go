package doubleratchet

import (
	"testing"
)

func TestInit(t *testing.T) {
	provider := stubDiffieHellmanProvider{}
	sharedKey := []byte{0x01, 0x02, 0x03}
	dhr := 0

	dr := Init(dhr, sharedKey, provider)

	if dr.skipped == nil {
		t.Error("Expected non nil value for skipped")
	}
}
