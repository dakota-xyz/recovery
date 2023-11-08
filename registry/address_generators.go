package registry

import "io"

// KeyGenerator describes a function that can produce
// an address and a private key
type KeyGenerator func(seed io.Reader) (string, string, error)

// Generators holds the available generators
var Generators = make(map[string]KeyGenerator)

// RegisterGenerator allows the registration of new generators
func RegisterGenerator(id string, generator KeyGenerator) {
	Generators[id] = generator
}
