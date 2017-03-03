package types

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"regexp"

	"github.com/minio/blake2b-simd"
)

// Hasher is the hash algorithm and string id
type Hasher struct {
	name   string
	hasher func() hash.Hash
}

// Name is the hash algorithm
func (h *Hasher) Name() string {
	return h.name
}

// From hashes a io.Reader
func (h *Hasher) From(in io.Reader) ([]byte, error) {
	hasher := h.hasher()
	_, err := io.Copy(hasher, in)
	if err != nil {
		return nil, err
	}

	// Hash and print. Pass nil since
	// the data is not coming in as a slice argument
	// but is coming through the writer interface
	return hasher.Sum(nil), nil
}

var (
	// SHA1Hasher is the default sha1 method
	SHA1Hasher = Hasher{name: "sha1", hasher: sha1.New}
	// Blake2bHasher is an optimized fast hash
	Blake2bHasher = Hasher{name: "blake2b", hasher: func() hash.Hash { hsh, _ := blake2b.New(&blake2b.Config{Size: 32}); return hsh }}
)

// Hash represents a hash method, the hash value and serializing to
// a string URI representation.
type Hash struct {
	hash   []byte
	method string
}

// Hash is the []byte representation
func (h *Hash) Hash() []byte {
	return h.hash
}

// Base64Hash is the Base64 encoded representation
func (h *Hash) Base64Hash() string {
	return base64.URLEncoding.EncodeToString(h.hash)
}

// HashMethod is hash algorithm used
func (h *Hash) HashMethod() string {
	return h.method
}

// URI is the string URI representation
func (h *Hash) URI() string {
	return fmt.Sprintf("%s://%s", h.method, h.Base64Hash())
}

// String is the URI representation
func (h *Hash) String() string {
	return h.URI()
}

// HashFromReader returns a Hash from an io.Reader using the given Hasher
func HashReader(hasher Hasher, data io.Reader) (*Hash, error) {
	hsh, err := hasher.From(data)
	if err != nil {
		return nil, err
	}

	return &Hash{hash: hsh, method: hasher.Name()}, nil
}

// HashFromURI desearlizes a Hash from a string URI
func HashFromURI(uri string) (*Hash, error) {
	comp := regexp.MustCompile("://").Split(uri, -1)

	if len(comp) != 2 {
		msg := fmt.Sprintf("Malformed Hash URI: %s\nComponents: %v", uri, comp)
		return nil, errors.New(msg)
	}

	hsh, err := base64.URLEncoding.DecodeString(comp[1])

	if err != nil {
		return nil, err
	}

	return &Hash{method: comp[0], hash: hsh}, nil
}
