package selfupdate

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"testing"
)

func TestSHA2Validator(t *testing.T) {
	validator := &SHA2Validator{}
	data, err := ioutil.ReadFile("testdata/foo.zip")
	if err != nil {
		t.Fatal(err)
	}
	hashData, err := ioutil.ReadFile("testdata/foo.zip.sha256")
	if err != nil {
		t.Fatal(err)
	}
	if err := validator.Validate(data, hashData); err != nil {
		t.Fatal(err)
	}
}

func TestSHA2ValidatorFail(t *testing.T) {
	validator := &SHA2Validator{}
	data, err := ioutil.ReadFile("testdata/foo.zip")
	if err != nil {
		t.Fatal(err)
	}
	hashData, err := ioutil.ReadFile("testdata/foo.zip.sha256")
	if err != nil {
		t.Fatal(err)
	}
	hashData[0] = '0'
	if err := validator.Validate(data, hashData); err == nil {
		t.Fatal(err)
	}
}

func TestECDSAValidator(t *testing.T) {
	pemData, err := ioutil.ReadFile("testdata/Test.crt")
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate")
	}

	pubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Errorf("PublicKey is not ECDSA")
	}

	validator := &ECDSAValidator{
		PublicKey: pubKey,
	}
	data, err := ioutil.ReadFile("testdata/foo.zip")
	if err != nil {
		t.Fatal(err)
	}
	signatureData, err := ioutil.ReadFile("testdata/foo.zip.sig")
	if err != nil {
		t.Fatal(err)
	}
	if err := validator.Validate(data, signatureData); err != nil {
		t.Fatal(err)
	}
}

func TestECDSAValidatorFail(t *testing.T) {
	pemData, err := ioutil.ReadFile("testdata/Test.crt")
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate")
	}

	pubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Errorf("PublicKey is not ECDSA")
	}

	validator := &ECDSAValidator{
		PublicKey: pubKey,
	}
	data, err := ioutil.ReadFile("testdata/foo.tar.xz")
	if err != nil {
		t.Fatal(err)
	}
	signatureData, err := ioutil.ReadFile("testdata/foo.zip.sig")
	if err != nil {
		t.Fatal(err)
	}
	if err := validator.Validate(data, signatureData); err == nil {
		t.Fatal(err)
	}
}

func TestValidatorSuffix(t *testing.T) {
	for _, test := range []struct {
		v      Validator
		suffix string
	}{
		{
			v:      &SHA2Validator{},
			suffix: ".sha256",
		},
		{
			v:      &ECDSAValidator{},
			suffix: ".sig",
		},
	} {
		want := test.suffix
		got := test.v.Suffix()
		if want != got {
			t.Errorf("Wanted %q but got %q", want, got)
		}
	}
}
