package secretsmanager

import (
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGenerateSSHKeyPair_AllTypes(t *testing.T) {
	tests := []struct {
		name       string
		keyType    SSHKeyType
		keyBits    int
		wantPrefix string
	}{
		{"ed25519", SSHKeyTypeED25519, 0, "ssh-ed25519 "},
		{"rsa-4096", SSHKeyTypeRSA, 4096, "ssh-rsa "},
		{"rsa-2048", SSHKeyTypeRSA, 2048, "ssh-rsa "},
		{"rsa-3072", SSHKeyTypeRSA, 3072, "ssh-rsa "},
		{"ecdsa-256", SSHKeyTypeECDSA256, 0, "ecdsa-sha2-nistp256 "},
		{"ecdsa-384", SSHKeyTypeECDSA384, 0, "ecdsa-sha2-nistp384 "},
		{"ecdsa-521", SSHKeyTypeECDSA521, 0, "ecdsa-sha2-nistp521 "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateSSHKeyPair(tt.keyType, tt.keyBits, "")
			if err != nil {
				t.Fatalf("GenerateSSHKeyPair(%s) error = %v", tt.keyType, err)
			}
			if !strings.HasPrefix(result.PublicKey, tt.wantPrefix) {
				t.Errorf("public key prefix = %q, want prefix %q", result.PublicKey[:30], tt.wantPrefix)
			}
			if !strings.Contains(result.PrivateKey, "-----BEGIN OPENSSH PRIVATE KEY-----") {
				t.Errorf("private key missing PEM header")
			}
			if !strings.Contains(result.PrivateKey, "-----END OPENSSH PRIVATE KEY-----") {
				t.Errorf("private key missing PEM footer")
			}
		})
	}
}

func TestGenerateSSHKeyPair_WithPassphrase(t *testing.T) {
	passphrase := "test-passphrase-123"
	result, err := GenerateSSHKeyPair(SSHKeyTypeED25519, 0, passphrase)
	if err != nil {
		t.Fatalf("GenerateSSHKeyPair with passphrase error = %v", err)
	}

	// Verify private key can be parsed with correct passphrase
	_, err = ssh.ParseRawPrivateKeyWithPassphrase([]byte(result.PrivateKey), []byte(passphrase))
	if err != nil {
		t.Errorf("failed to parse private key with correct passphrase: %v", err)
	}

	// Verify private key cannot be parsed without passphrase
	_, err = ssh.ParseRawPrivateKey([]byte(result.PrivateKey))
	if err == nil {
		t.Errorf("expected error parsing encrypted private key without passphrase")
	}
}

func TestGenerateSSHKeyPair_PublicKeyMatchesPrivate(t *testing.T) {
	tests := []struct {
		name    string
		keyType SSHKeyType
		keyBits int
	}{
		{"ed25519", SSHKeyTypeED25519, 0},
		{"rsa", SSHKeyTypeRSA, 4096},
		{"ecdsa-256", SSHKeyTypeECDSA256, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateSSHKeyPair(tt.keyType, tt.keyBits, "")
			if err != nil {
				t.Fatalf("GenerateSSHKeyPair error = %v", err)
			}

			// Parse the private key back
			privKey, err := ssh.ParseRawPrivateKey([]byte(result.PrivateKey))
			if err != nil {
				t.Fatalf("failed to parse private key: %v", err)
			}

			// Extract public key from parsed private key
			signer, err := ssh.NewSignerFromKey(privKey)
			if err != nil {
				t.Fatalf("failed to create signer: %v", err)
			}
			derivedPubKey := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey())))

			if derivedPubKey != result.PublicKey {
				t.Errorf("public key mismatch:\n  got:  %s\n  want: %s", derivedPubKey, result.PublicKey)
			}
		})
	}
}

func TestGenerateSSHKeyPair_InvalidType(t *testing.T) {
	_, err := GenerateSSHKeyPair("ssh-dss", 0, "")
	if err == nil {
		t.Error("expected error for unsupported key type")
	}
	if !strings.Contains(err.Error(), "unsupported SSH key type") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateSSHKeyPair_RSADefaultBits(t *testing.T) {
	// When keyBits < 2048, should default to 4096
	result, err := GenerateSSHKeyPair(SSHKeyTypeRSA, 0, "")
	if err != nil {
		t.Fatalf("GenerateSSHKeyPair error = %v", err)
	}
	if !strings.HasPrefix(result.PublicKey, "ssh-rsa ") {
		t.Errorf("expected ssh-rsa public key")
	}
}
