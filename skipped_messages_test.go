package doubleratchet

import (
	"bytes"
	"testing"
)

var notInPersistorKey DiffieHellmanPublic = "This key sucks"
var inPersistorKey DiffieHellmanPublic = "This key rules"
var notInPersistorN = 0
var inPersistorN = 1
var stubMk = []byte{0x01}

func createStubPersistor() mappedSkipPersistor {
	return mappedSkipPersistor{
		inPersistorKey: map[int][]byte{inPersistorN: stubMk},
	}
}

func TestMappedSkipPersistorGet(t *testing.T) {
	cases := []struct {
		key    DiffieHellmanPublic
		n      int
		want   []byte
		reason string
	}{
		{notInPersistorKey, notInPersistorN, nil, "expected nil result if key is not present"},
		{notInPersistorKey, inPersistorN, nil, "expected nil result if key is not present but n is"},
		{inPersistorKey, notInPersistorN, nil, "expected nil result if key is present but n is not"},
		{inPersistorKey, inPersistorN, stubMk, "expected real result if key and n is present"},
	}

	for _, c := range cases {
		persistor := createStubPersistor()
		got := persistor.Get(c.key, c.n)

		if !bytes.Equal(got, c.want) {
			t.Errorf("%s:\n\twant:%+v\n\t got:%+v", c.reason, c.want, got)
		}

		if got != nil && len(persistor) != 0 {
			t.Error("Expected Get to cleanup map when value is returend")
		}
	}
}

func TestMappedSkipPersistorPut(t *testing.T) {
	cases := []struct {
		dhCount int
		nCount  int
		want    error
		reason  string
	}{
		{1, 1, nil, "Expected no error when inserting a single value"},
		{maxSkippedKeys - 1, maxSkippedKeys - 1, nil, "Expected no errors when persistor is filled"},
		{maxSkippedKeys + 2, 1, errPersistorFull, "Expected error when inserting too many dh values"},
		{1, maxSkippedKeys + 11, errPersistorFull, "Expected error when inserting too many n values"},
	}

	for _, c := range cases {
		var got error
		persistor := mappedSkipPersistor{}
		for i := 0; i < c.dhCount; i++ {
			for j := 0; j < c.nCount; j++ {
				if v := persistor.Put(i, j, stubMk); v != nil && got == nil {
					got = v
				}

				if got == nil && !bytes.Equal(persistor[i][j], stubMk) {
					t.Errorf("Expected to find stub value at %d, %d; Got %+v", i, j, persistor[i][j])
				}
			}
		}

		if got != c.want {
			t.Errorf("%s:\n\twant:%+v\n\t got:%+v\n\t%+v", c.reason, c.want, got, persistor)
		}
	}
}
