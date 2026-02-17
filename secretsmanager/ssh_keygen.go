package secretsmanager

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SSHKeyType represents the supported SSH key types
type SSHKeyType string

const (
	SSHKeyTypeED25519  SSHKeyType = "ssh-ed25519"
	SSHKeyTypeRSA      SSHKeyType = "ssh-rsa"
	SSHKeyTypeECDSA256 SSHKeyType = "ecdsa-sha2-nistp256"
	SSHKeyTypeECDSA384 SSHKeyType = "ecdsa-sha2-nistp384"
	SSHKeyTypeECDSA521 SSHKeyType = "ecdsa-sha2-nistp521"
)

// SSHKeyPairResult holds the generated key pair output
type SSHKeyPairResult struct {
	PublicKey  string // OpenSSH authorized_keys format
	PrivateKey string // OpenSSH PEM format
}

// GenerateSSHKeyPair generates an SSH key pair of the specified type.
// If passphrase is non-empty, the private key PEM will be encrypted.
// keyBits is only used for RSA keys.
func GenerateSSHKeyPair(keyType SSHKeyType, keyBits int, passphrase string) (*SSHKeyPairResult, error) {
	var privateKey interface{}
	var err error

	switch keyType {
	case SSHKeyTypeED25519:
		_, privateKey, err = ed25519.GenerateKey(rand.Reader)
	case SSHKeyTypeRSA:
		if keyBits < 2048 {
			keyBits = 4096
		}
		privateKey, err = rsa.GenerateKey(rand.Reader, keyBits)
	case SSHKeyTypeECDSA256:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case SSHKeyTypeECDSA384:
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case SSHKeyTypeECDSA521:
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("unsupported SSH key type: %s", keyType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to generate %s key: %w", keyType, err)
	}

	// Marshal public key to OpenSSH authorized_keys format
	pubKey, err := ssh.NewPublicKey(publicKeyFromPrivate(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH public key: %w", err)
	}
	publicKeyStr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pubKey)))

	// Marshal private key to PEM
	privateKeyPEM, err := marshalPrivateKeyPEM(privateKey, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key to PEM: %w", err)
	}

	return &SSHKeyPairResult{
		PublicKey:  publicKeyStr,
		PrivateKey: string(privateKeyPEM),
	}, nil
}

// publicKeyFromPrivate extracts the public key from a private key
func publicKeyFromPrivate(key interface{}) interface{} {
	switch k := key.(type) {
	case ed25519.PrivateKey:
		return k.Public()
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

// marshalPrivateKeyPEM encodes a private key to PEM format,
// optionally encrypting with a passphrase.
func marshalPrivateKeyPEM(key interface{}, passphrase string) ([]byte, error) {
	if passphrase != "" {
		pemBlock, err := ssh.MarshalPrivateKeyWithPassphrase(key, "", []byte(passphrase))
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(pemBlock), nil
	}

	pemBlock, err := ssh.MarshalPrivateKey(key, "")
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(pemBlock), nil
}
