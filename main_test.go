package main

import (
	"bytes"
	"testing"
)

func TestJWTSigningKey(t *testing.T) {
	const secret = "a-real-test-secret-with-more-than-32-bytes"
	k1, err := jwtSigningKey("first-password", secret)
	if err != nil {
		t.Fatal(err)
	}
	k2, err := jwtSigningKey("first-password", secret)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(k1, k2) {
		t.Fatal("permanent signing key changed unexpectedly")
	}
	k3, err := jwtSigningKey("second-password", secret)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(k1, k3) {
		t.Fatal("password change did not revoke tokens")
	}
	random1, err := jwtSigningKey("first-password", "")
	if err != nil {
		t.Fatal(err)
	}
	random2, err := jwtSigningKey("first-password", "")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(random1, random2) {
		t.Fatal("random key repeated")
	}
	if _, err := jwtSigningKey("password", "weak"); err == nil {
		t.Fatal("weak manually supplied JWT secret was accepted")
	}
}
