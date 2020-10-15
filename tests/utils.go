package tests

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

const (
	awsbiImageTag = "epiphanyplatform/awsbi:0.0.1"
	awsksImageTag = "epiphanyplatform/awsks:0.0.1"
)

func setup(suffix string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := path.Join(wd, fmt.Sprintf("%s-%s", "shared", suffix))
	return p, os.MkdirAll(p, os.ModePerm)
}

func cleanup(sharedPath string) error {
	return os.RemoveAll(sharedPath)
}

func normStr(s string) string {
	return strings.TrimSpace(s)
}

func getLastLineFromMultilineString(s string) (string, error) {
	in := strings.NewReader(s)
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF {
			return string(line), nil
		}
	}
}

func generateRsaKeyPair(directory, name string) error {
	privateRsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateRsaKey)}
	privateKeyBytes := pem.EncodeToMemory(pemBlock)

	publicRsaKey, err := ssh.NewPublicKey(&privateRsaKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	err = ioutil.WriteFile(path.Join(directory, name), privateKeyBytes, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(directory, fmt.Sprintf("%s.pub", name)), publicKeyBytes, 0644)
}

func getAwsCreds(t *testing.T) (awsAccessKey, awsSecretKey string) {
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
	if len(awsAccessKey) == 0 {
		t.Fatalf("expected non-empty AWS_ACCESS_KEY environment variable")
	}

	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
	if len(awsSecretKey) == 0 {
		t.Fatalf("expected non-empty AWS_SECRET_KEY environment variable")
	}

	return
}
