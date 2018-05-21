package AEDAcrypt

import (
	"encoding/hex"
	"testing"
)

func TestEncrypterDecrypter(t *testing.T) {
	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	ccfg := CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}

	var str = "Hello World!"
	byteenc, encerr := Encrypter([]byte(str), ccfg)
	if encerr != nil {
		t.Error("Got an error from Encrypter")
	}

	bytedec, decerr := Decrypter(byteenc, ccfg)
	strdec := string(bytedec[:])

	if decerr != nil && str != strdec {
		t.Error("Got an error from Decrypter")
	}

	if str != strdec {
		t.Errorf("Expected %s, got %s", str, strdec)
	}
}

func TestWrongEncryptionKey(t *testing.T) {
	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	ccfg := CryptCfg{Key: []byte("keytoshort"),
		Nonce: nonce}

	var str = "Hello World!"
	_, encerr := Encrypter([]byte(str), ccfg)
	if encerr == nil {
		t.Error("Did not get an error from Encrypter despite invalid keys")
	}
}

func TestWrongDecryptionKey(t *testing.T) {
	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
	ccfgencrypt := CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}
	ccfgdecrypt := CryptCfg{Key: []byte("WRONGKEY--32Characters1234567890"),
		Nonce: nonce}

	var str = "Hello World!"
	byteenc, encerr := Encrypter([]byte(str), ccfgencrypt)
	if encerr != nil {
		t.Error("Got an error from Encrypter")
	}

	bytedec, decerr := Decrypter(byteenc, ccfgdecrypt)
	strdec := string(bytedec[:])

	if decerr == nil {
		t.Error("Did not get error from Decrypter despite wrong key")
	}

	if strdec != "" {
		t.Error("Decrypted string is not empty")
	}
}

func TestGetMD5HashFromString(t *testing.T) {
	var str = "Hello World!"
	var expectedHash = "ed076287532e86365e841e92bfc50d8c" // echo -n "Hello World!" | md5sum
	hash := GetMD5HashFromString(str)

	if hash != expectedHash {
		t.Errorf("Expected %s, got %s", expectedHash, hash)
	}
}

func TestGetMD5HashFromByte(t *testing.T) {
	var str = "Hello World!"
	var expectedHash = "ed076287532e86365e841e92bfc50d8c" // echo -n "Hello World!" | md5sum
	hash := GetMD5HashFromByte([]byte(str))

	if hash != expectedHash {
		t.Errorf("Expected %s, got %s", expectedHash, hash)
	}
}
