package handlers

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	cryptorand "crypto/rand"

	"github.com/sfarosu/go-tooling-portal/internal/tmpl"
	"golang.org/x/crypto/ssh"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	sshKeyGenCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ssh_key_generated_total",
		Help: "The total number of generated ssh keys",
	})
)

func sshKeyGen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/sshkeygen", http.StatusSeeOther)
	}
	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "sshkeygen.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func sshProcessKeypair(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sshkeygen", http.StatusSeeOther)
	}

	privateKeyContent, publicKeyContent, err := generateSSHKeyPair(r.FormValue("algorithm-select"), r.FormValue("ecdsa-bitsize-select"), r.FormValue("rsa-bitsize-select"), strings.TrimSpace(r.FormValue("email")))
	if err != nil {
		log.Println("error generating ssh key pair: ", err)
	}

	data := struct {
		Algorithm         string
		BitSize           string
		Email             string
		UsePass           string
		Pass              string
		PrivateKeyContent string
		PublicKeyContent  string
	}{
		Algorithm:         r.FormValue("algorithm-select"),
		BitSize:           "",
		Email:             r.FormValue("email"),
		UsePass:           r.FormValue("usepass"),
		Pass:              r.FormValue("password"),
		PrivateKeyContent: privateKeyContent,
		PublicKeyContent:  publicKeyContent,
	}

	switch r.FormValue("algorithm-select") {
	case "ed25519":
		data.BitSize = "N/A"
	case "ecdsa":
		data.BitSize = r.FormValue("ecdsa-bitsize-select")
	case "rsa":
		data.BitSize = r.FormValue("rsa-bitsize-select")
	default:
		log.Println("unsupported algorithm; supported: ed25519, ecdsa, rsa")
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))

	err = tmpl.Tpl.ExecuteTemplate(w, "sshkeygen-process.html", data)
	if err != nil {
		log.Println("error executing template: ", err)
	}

	sshKeyGenCounter.Inc()
}

// generateSSHKeyPair creates SSH key pairs based on the specified algorithm
func generateSSHKeyPair(algorithm string, ecdsaBits string, rsaBits string, email string) (string, string, error) {
	var privateKeyPEM, publicKeySSH string
	var err error

	switch algorithm {
	case "ed25519":
		privateKeyPEM, publicKeySSH, err = generateEd25519Key(email)
	case "ecdsa":
		ecdsaBitsInt, errConv := strconv.Atoi(ecdsaBits)
		if errConv != nil {
			log.Printf("error converting [%v] string to int: %v", ecdsaBits, errConv)
			return "", "", errConv
		}
		privateKeyPEM, publicKeySSH, err = generateECDSAKey(ecdsaBitsInt, email)
	case "rsa":
		rsaBitsInt, errConv := strconv.Atoi(rsaBits)
		if errConv != nil {
			log.Printf("error converting [%v] string to int: %v", rsaBits, errConv)
			return "", "", errConv
		}
		privateKeyPEM, publicKeySSH, err = generateRSAKey(rsaBitsInt, email)
	default:
		return "", "", errors.New("unsupported algorithm; supported: ed25519, rsa, ecdsa")
	}

	if err != nil {
		return "", "", err
	}

	return privateKeyPEM, publicKeySSH, nil
}

// generateEd25519Key creates an Ed25519 key pair and returns the private and public keys
func generateEd25519Key(email string) (string, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		return "", "", err
	}

	// Create an openSSH formated public key from the generated Ed25519 public key
	sshPubKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	// Marshal the openSSH public key to the authorized key format and add a comment (email) at the end
	publicKeyWithComment := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPubKey))) + " " + email

	// Create a PEM block to encode the private key
	// The PEM block contains the type "OPENSSH PRIVATE KEY" and the marshaled private key bytes
	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: marshalED25519PrivateKey(privateKey), // marshals ed25519 correctly
	}

	// Encode the PEM block to a byte slice
	privKey := pem.EncodeToMemory(pemKey)

	return string(privKey), string(publicKeyWithComment), nil
}

// generateECDSAKey creates a ECDSA key pair and returns the private and public keys
func generateECDSAKey(bits int, email string) (string, string, error) {
	var curve elliptic.Curve
	switch bits {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return "", "", errors.New("unsupported curve size; supported: 256, 384, 521")
	}

	// Generate ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(curve, cryptorand.Reader)
	if err != nil {
		return "", "", err
	}

	// Marshal the ECDSA private key to ASN.1 DER format
	privKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	// Encode the DER-encoded private key to PEM format
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	// Create the public key from the private key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	// Format the public key with a comment (email)
	publicKeyWithComment := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey))) + " " + email

	return string(privKeyPEM), publicKeyWithComment, nil
}

// generateRSAKey creates a RSA key pair and returns the private and public keys
func generateRSAKey(bits int, email string) (string, string, error) {
	privateKey, err := rsa.GenerateKey(cryptorand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// Encode the RSA private key to PKCS#1 format and then to PEM format
	marshalledRSAPrivateKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))

	// Create an openSSH public key from the RSA public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	// Format the public key with a comment (email)
	publicKeyWithComment := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey))) + " " + email

	return marshalledRSAPrivateKey, publicKeyWithComment, nil
}

// The x509 package does not support marshaling ed25519 key types in the format used by openSSH
// source: https://github.com/mikesmitty/edkey/blob/master/edkey.go#L10
func marshalED25519PrivateKey(key ed25519.PrivateKey) []byte {
	// Add our key header (followed by a null byte)
	magic := append([]byte("openssh-key-v1"), 0)

	var w struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	}

	// Fill out the private key fields
	pk1 := struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}{}

	// Set our check ints
	ci := rand.Uint32()
	pk1.Check1 = ci
	pk1.Check2 = ci

	// Set our key type
	pk1.Keytype = ssh.KeyAlgoED25519

	// Add the pubkey to the optionally-encrypted block
	pk, ok := key.Public().(ed25519.PublicKey)
	if !ok {
		//fmt.Fprintln(os.Stderr, "ed25519.PublicKey type assertion failed on an ed25519 public key. This should never ever happen.")
		return nil
	}
	pubKey := []byte(pk)
	pk1.Pub = pubKey

	// Add our private key
	pk1.Priv = []byte(key)

	// Might be useful to put something in here at some point
	pk1.Comment = ""

	// Add some padding to match the encryption block size within PrivKeyBlock (without Pad field)
	// 8 doesn't match the documentation, but that's what ssh-keygen uses for unencrypted keys. *shrug*
	bs := 8
	blockLen := len(ssh.Marshal(pk1))
	padLen := (bs - (blockLen % bs)) % bs
	pk1.Pad = make([]byte, padLen)

	// Padding is a sequence of bytes like: 1, 2, 3...
	for i := 0; i < padLen; i++ {
		pk1.Pad[i] = byte(i + 1)
	}

	// Generate the pubkey prefix "\0\0\0\nssh-ed25519\0\0\0 "
	prefix := []byte{0x0, 0x0, 0x0, 0x0b}
	prefix = append(prefix, []byte(ssh.KeyAlgoED25519)...)
	prefix = append(prefix, []byte{0x0, 0x0, 0x0, 0x20}...)

	// Only going to support unencrypted keys for now
	w.CipherName = "none"
	w.KdfName = "none"
	w.KdfOpts = ""
	w.NumKeys = 1
	w.PubKey = append(prefix, pubKey...)
	w.PrivKeyBlock = ssh.Marshal(pk1)

	magic = append(magic, ssh.Marshal(w)...)

	return magic
}
