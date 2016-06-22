package licensekey

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	block, _ := pem.Decode(privateKeyPEM)
	privateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)

	block, _ = pem.Decode(publicKeyPEM)
	publicKeyIface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey = publicKeyIface.(*rsa.PublicKey)
}

func TestGenerateKeys(t *testing.T) {
	_, err := GenerateKeys(4096)
	assertNoError(t, err)
}

func TestStorePrivateKey(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tmpfile.Name())
	tmpfilename := tmpfile.Name()
	tmpfile.Close()

	err := StorePrivateKey(privateKey, tmpfilename)
	assertNoError(t, err)
	pem, err := ioutil.ReadFile(tmpfilename)
	assertNoError(t, err)
	assertEqualString(t, string(privateKeyPEM), string(pem))
}

func TestStorePublicKey(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tmpfile.Name())
	tmpfilename := tmpfile.Name()
	tmpfile.Close()

	err := StorePublicKey(publicKey, tmpfilename)
	assertNoError(t, err)
	pem, err := ioutil.ReadFile(tmpfilename)
	assertNoError(t, err)
	assertEqualString(t, string(publicKeyPEM), string(pem))
}

func TestLoadPrivateKey(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tmpfile.Name())
	io.WriteString(tmpfile, string(privateKeyPEM))
	tmpfile.Close()

	key, err := LoadPrivateKey(tmpfile.Name())
	assertNoError(t, err)
	assertEqualPrivateKey(t, privateKey, key)
}

func TestLoadPublicKey(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tmpfile.Name())
	io.WriteString(tmpfile, string(publicKeyPEM))
	tmpfile.Close()

	key, err := LoadPublicKey(tmpfile.Name())
	assertNoError(t, err)
	assertEqualPublicKey(t, publicKey, key)
}

func TestPrivateKeyFromBytes(t *testing.T) {
	key, err := PrivateKeyFromBytes(privateKeyPEM)
	assertNoError(t, err)
	assertEqualPrivateKey(t, privateKey, key)
}

func TestPublicKeyFromBytes(t *testing.T) {
	key, err := PublicKeyFromBytes(publicKeyPEM)
	assertNoError(t, err)
	assertEqualPublicKey(t, publicKey, key)
}

func TestSign(t *testing.T) {
	signature, err := Sign(privateKey, tMessage)
	assertNoError(t, err)
	assertEqualString(t, string(tSignature), string(signature))
}

func TestVerify(t *testing.T) {
	err := Verify(publicKey, tMessage, tSignature)
	assertNoError(t, err)
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("nil != %v", err)
	}
}

func assertEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("%s != %s", expected, actual)
	}
}

func assertEqualPrivateKey(t *testing.T, expected, actual *rsa.PrivateKey) {
	if expected.D.Cmp(actual.D) != 0 || expected.N.Cmp(actual.N) != 0 || expected.E != actual.E {
		t.Fatalf("%v != %v", expected, actual)
	}
}

func assertEqualPublicKey(t *testing.T, expected, actual *rsa.PublicKey) {
	if expected.N.Cmp(actual.N) != 0 || expected.E != actual.E {
		t.Fatalf("%v != %v", expected, actual)
	}
}

var (
	privateKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`)
	publicKeyPEM = []byte(`-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA4s9fb3akEQtq1uoYLjK5
g3K3iZ7L4FxWhxmNEGpppSDdXtMBdAM0vx+wwdtgFRomc7zJ4k/+59Uy/sHQ7zdJ
mltZ8/bGdUjYs5/R7wo9nsZ2SN6RuPDocc0kMVQydwU2PossbmilJWcfwr5mSD48
xQCxcPq4YYVS1vuvM+ALZt97khmQKYpTVpcCaZaPxkVvrIZ3UK9Ea940OEp+n0h9
zMYiU/DDOBluZ9v8fiIWZo5Xx+xsKY74HpZAamDPwmJwtT/uR9jXKjmuLCRUA1+A
4vdm1YxFL78Q2E/4Ye9nZQ4J6KiZX3HxMajeTwKzTFxvLThMJOVeK+OOGzIYGPq3
rJSLk1rmSheUM8i24ygar2t1k/vE/73C+JZhUCU4QjOdqXqlMnt4CkBSkbjKc5OL
y3XMGHcdDTEWpXWhaXg2Wgbw1ULXu+B9BkHxfcx/cx563jOzbXfAGXsiAuJuFSW4
JrhyCVH2bIqMwVI5G5mEQd1TZXmgvnvCez4daD+cJKiQytsO7EMQd8aafJEYGpcZ
ByIP3dypvQL+v+oHwEelskfh+Nq8ptJyR0lC5n043z/VEUHQHpVCq6SBlzR1lRkM
d3BOiyBeOWDCGQAk9rBT39HqdsMGw8PNGlFZ+XVxfd4bThntZuSfCtaDNCCtRt0x
Yz0hycBaTZyViPOgXllYdFECAwEAAQ==
-----END RSA PUBLIC KEY-----
`)
	tMessage   = []byte("message")
	tSignature = []byte{0xc6, 0x42, 0x42, 0xd9, 0x46, 0x10, 0xb6, 0xb4, 0x8b,
		0xa8, 0x85, 0x13, 0xf1, 0xaa, 0xc4, 0xbd, 0x3a, 0xe0, 0x9a, 0xe8, 0x7e, 0xe6,
		0xf, 0x48, 0x40, 0x22, 0x7f, 0x8b, 0x2c, 0xa7, 0x69, 0x34, 0x97, 0x65, 0x8e,
		0x7b, 0xa9, 0xe7, 0x89, 0x3b, 0x58, 0x41, 0xe3, 0xa2, 0x8a, 0x81, 0xf8, 0x71,
		0xb0, 0x90, 0x8c, 0xa4, 0x12, 0x0, 0x2c, 0x4f, 0xca, 0xb5, 0x90, 0x6f, 0x12,
		0x55, 0xc7, 0xf6, 0x45, 0xd9, 0x47, 0xcc, 0xe7, 0x1b, 0x7, 0x2d, 0x37, 0xc0,
		0x8c, 0xc0, 0x2, 0x40, 0x2c, 0x71, 0xfc, 0x2e, 0xc7, 0x3e, 0xef, 0x55, 0x82,
		0x87, 0xe7, 0xdb, 0x4a, 0x2f, 0xe, 0xf0, 0xc2, 0xa7, 0x79, 0xbf, 0xfd, 0xaf,
		0x30, 0x28, 0xe0, 0x92, 0x13, 0x9c, 0x6b, 0xd, 0x95, 0x40, 0x93, 0xdf, 0xa,
		0xb3, 0x22, 0x6a, 0x6f, 0x3c, 0xb9, 0x48, 0xdc, 0xb, 0x91, 0x94, 0x63, 0xd9,
		0x5, 0x6, 0xe9, 0xf7, 0x85, 0x95, 0x62, 0x54, 0x43, 0x56, 0xf2, 0x59, 0x90,
		0xc3, 0x19, 0xe0, 0xb, 0xda, 0x95, 0x5d, 0x73, 0x18, 0xba, 0xf3, 0x87, 0x15,
		0xdb, 0x5d, 0xcb, 0x30, 0x45, 0x8a, 0xfc, 0x9, 0xa6, 0x8f, 0xcc, 0xdd, 0xee,
		0x26, 0xed, 0x1c, 0xec, 0xa2, 0x2c, 0x49, 0x98, 0xa6, 0x9d, 0x78, 0xc6, 0xec,
		0x27, 0x0, 0x94, 0x4, 0x4, 0x6f, 0x70, 0x77, 0x7c, 0xd0, 0xb3, 0xe3, 0x4b,
		0x53, 0x44, 0x2a, 0x28, 0x4c, 0x9c, 0x3f, 0x1f, 0x1f, 0x9f, 0x93, 0x2a, 0xc8,
		0xba, 0x93, 0xc7, 0x1a, 0xe0, 0x84, 0x7c, 0x52, 0x26, 0x5d, 0xb8, 0xa, 0xad,
		0xfb, 0x2c, 0x3d, 0xbf, 0x46, 0xa, 0x41, 0xe7, 0x30, 0xe0, 0xca, 0x4d, 0xf6,
		0x37, 0xf4, 0xd8, 0xb6, 0xa8, 0xed, 0xfc, 0x9a, 0x14, 0xff, 0x81, 0xa5, 0x83,
		0xb3, 0xed, 0x7f, 0x86, 0xb2, 0x90, 0xb9, 0x12, 0x28, 0x3b, 0xa2, 0xda, 0xb9,
		0xaf, 0xbf, 0xd2, 0x12, 0x1c, 0xc6, 0xc0, 0x7f, 0xa2, 0x4c, 0xe8, 0x7, 0x8d,
		0xa4, 0x2e, 0xbc, 0xf0, 0xbb, 0xbe, 0xb7, 0x87, 0xa6, 0x5b, 0x8, 0xac, 0x42,
		0xeb, 0xc, 0x11, 0x6c, 0x42, 0x25, 0x82, 0x73, 0x96, 0x4a, 0xfc, 0xbb, 0x93,
		0x24, 0x74, 0x6a, 0x4a, 0xc1, 0xdd, 0x73, 0xa5, 0x6d, 0x62, 0x53, 0x1c, 0xc1,
		0x98, 0x14, 0x23, 0x3, 0xc5, 0x3, 0x3a, 0x1c, 0x61, 0xb3, 0x57, 0xbf, 0x82,
		0xad, 0xb9, 0x5a, 0x84, 0x2b, 0x3c, 0x40, 0xf, 0xd2, 0xc1, 0x6c, 0xf1, 0xd7,
		0x2a, 0x98, 0x59, 0x12, 0x31, 0xba, 0xd8, 0xd4, 0x80, 0x41, 0xa1, 0xc6, 0x69,
		0x36, 0x2f, 0x4d, 0xf2, 0xa5, 0x76, 0x91, 0x9d, 0xe0, 0xad, 0x99, 0xcb, 0x97,
		0xf1, 0x28, 0xff, 0xd1, 0x54, 0xd1, 0xf5, 0xe4, 0x1e, 0x32, 0x3, 0x10, 0xd5,
		0xc5, 0xc4, 0x85, 0xa0, 0x5d, 0x33, 0x2, 0xfe, 0xb0, 0x6, 0xf3, 0xa9, 0x5d,
		0x18, 0x5a, 0xf3, 0x13, 0x73, 0x46, 0x52, 0x40, 0x76, 0x21, 0x64, 0xff, 0x59,
		0x13, 0x61, 0xa0, 0xd1, 0x61, 0xc2, 0x11, 0x50, 0x8f, 0xb2, 0xad, 0x89, 0x9,
		0x75, 0xef, 0x77, 0x6d, 0x36, 0x22, 0x5e, 0x7a, 0x94, 0x15, 0x23, 0x15, 0x89,
		0xd3, 0xc1, 0x80, 0xaf, 0xfe, 0x5c, 0x46, 0x9d, 0x99, 0x5b, 0xc8, 0xae, 0xc6,
		0x99, 0x85, 0x4c, 0xb5, 0x3e, 0x16, 0x37, 0x12, 0x25, 0x2e, 0xa7, 0xf4, 0x24,
		0xb0, 0xe0, 0xce, 0x5c, 0x2c, 0x56, 0xd5, 0x11, 0x6c, 0x9d, 0xd1, 0x4c, 0xff,
		0xb3, 0x32, 0x80, 0xab, 0x56, 0xc3, 0x27, 0x25, 0xa5, 0x6a, 0x16, 0x84, 0x18,
		0x28, 0x94, 0x4f, 0x52, 0xc3, 0x30, 0x33, 0x2c, 0x64, 0x35, 0x82, 0x45, 0x92,
		0xfe, 0x37, 0x7c, 0xdc, 0xa7, 0x95, 0x7d, 0x23, 0x5e, 0xda, 0xaa, 0x7b, 0x49,
		0xba, 0xb5, 0xf1, 0x4b, 0xce, 0x12, 0x92, 0x9e, 0x81}
)
