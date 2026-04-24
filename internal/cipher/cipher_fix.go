package cipher

// Ensure compile-time interface satisfaction.
// Other pipeline stages expose an Apply method; Cipher follows the same
// contract but returns an error because cryptographic operations are
// fallible. This file holds any supplementary helpers.

// Levels returns all recognised cipher levels.
func Levels() []Level {
	return []Level{LevelNone, LevelEncrypt, LevelDecrypt}
}
