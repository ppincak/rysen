package security

import (
	"testing"
)

func TestLoadAndCreateAccounts(t *testing.T) {
	_, err := LoadAndCreateAccounts("./accounts_test_encrypted.json", []byte("aesBasicTestKey1"))
	if err != nil {
		t.Error(err)
		t.Fatalf("There shouldnt be an error during loading the test json")
	}
}
