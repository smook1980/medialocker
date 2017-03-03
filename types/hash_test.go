package types_test

import (
	"strings"
	"testing"

	"github.com/smook1980/medialocker/types"
)

func TestHashReaderWithSha1Hasher(t *testing.T) {
	rd := strings.NewReader("Test Data")
	hash, err := types.HashReader(types.SHA1Hasher, rd)

	if err != nil {
		t.Errorf("Expected no errors, got %s.", err)
	}

	if hash.HashMethod() != "sha1" {
		t.Errorf("Expected sha1 HashMethod(), got ", hash.HashMethod())
	}
}

func TestHashReaderWithBlake2bHasher(t *testing.T) {
	rd := strings.NewReader("Test Data")
	hash, err := types.HashReader(types.Blake2bHasher, rd)

	if err != nil {
		t.Errorf("Expected no errors, got %s.", err)
	}

	if hash.HashMethod() != "blake2b" {
		t.Errorf("Expected blake2b HashMethod(), got ", hash.HashMethod())
	}

	expectedResult := "blake2b://ftsJHeXSytHmX5wSTT8_2miV7DexuwJxqteN9kF6AeI="
	if hash.URI() != expectedResult {
		t.errorF("Expected %s, got %s", expectedResult, hash.URI())
	}
}
