package identity

// Regression test for the ECDSA r||s width bug: big.Int.Bytes() drops leading
// zeros, so ~1/128 of signatures used to come out 63 bytes and fail Verify.

import "testing"

func TestSignFixedWidth(t *testing.T) {
	signer := NewIdentitySigner(MakeIdentity())
	msg := make([]byte, 32) // stand-in for a block digest

	for i := 0; i < 5000; i++ {
		msg[0] = byte(i)
		msg[1] = byte(i >> 8)
		sig, err := signer.Sign(msg)
		if err != nil {
			t.Fatalf("sign failed: %v", err)
		}
		if len(sig) != 64 {
			t.Fatalf("iteration %d: expected 64-byte signature, got %d", i, len(sig))
		}
	}
}
