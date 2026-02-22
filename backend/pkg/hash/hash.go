package hash

import (
	"backend/pkg/config"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Hash interface {
	Hash(password string) (string, error)
	Verify(password string, hashPassword string) (bool, error)
}

type Argon2 struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func (a *Argon2) Hash(password string) (string, error) {
	salt, err := generateSalt(int(a.saltLen))
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, a.time, a.memory, a.threads, a.keyLen)

	b64Salt := base64.RawStdEncoding.Strict().EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.Strict().EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.time, a.threads, b64Salt, b64Hash)

	return encoded, nil
}

func (a *Argon2) Verify(password string, hashPassword string) (bool, error) {
	parsed, salt, hash, err := decodeHash(hashPassword)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	otherHash := argon2.IDKey([]byte(password), salt, parsed.time, parsed.memory, parsed.threads, parsed.keyLen)

	return subtle.ConstantTimeCompare(otherHash, hash) == 1, nil
}

func NewArgon2(cfg config.Hash) *Argon2 {
	return &Argon2{
		time:    cfg.Time,
		memory:  cfg.Memory,
		threads: cfg.Threads,
		keyLen:  cfg.KeyLen,
		saltLen: cfg.SaltLen,
	}
}

var _ Hash = (*Argon2)(nil)

func generateSalt(saltLen int) ([]byte, error) {
	salt := make([]byte, saltLen)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return salt, nil
}

func decodeHash(encoded string) (p *Argon2, salt, hash []byte, err error) {
	vals := strings.Split(encoded, "$")
	if len(vals) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format")
	}

	var version int

	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse version: %w", err)
	}

	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("unsupported hash version: %d", version)
	}

	p = &Argon2{}

	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.threads); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	if len(salt) > math.MaxUint32 {
		return nil, nil, nil, fmt.Errorf("invalid salt length: %d", len(salt))
	}

	p.saltLen = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode hash: %w", err)
	}

	if len(hash) > math.MaxUint32 {
		return nil, nil, nil, fmt.Errorf("invalid hash length: %d", len(hash))
	}

	p.keyLen = uint32(len(hash))

	return p, salt, hash, nil
}
