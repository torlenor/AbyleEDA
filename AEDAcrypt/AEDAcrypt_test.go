package AEDAcrypt

import (
	"encoding/hex"
)

import "testing"

func TestEncrypterDecrypter(t *testing.T) {
    nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")
    ccfg := CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
        Nonce: nonce}
        
	var str string = "Hello World!"
	byteenc := Encrypter([]byte(str), ccfg)
	bytedec, err := Decrypter(byteenc, ccfg)

	if err != nil && str != string(bytedec[:]) {
		t.Error("Got an error from Decrypter")
	}
    
    strdec := string(bytedec[:])
    if str != strdec {
        t.Errorf("Expected %s, got %s", str, strdec)
    }
}

func TestGetMD5HashFromString(t *testing.T) {
	var str string = "Hello World!"
    var expectedHash string = "ed076287532e86365e841e92bfc50d8c" // echo -n "Hello World!" | md5sum
    hash := GetMD5HashFromString(str)
    
    if hash != expectedHash {
        t.Errorf("Expected %s, got %s", expectedHash, hash)
    }
}

func TestGetMD5HashFromByte(t *testing.T) {
    var str string = "Hello World!"
    var expectedHash string = "ed076287532e86365e841e92bfc50d8c" // echo -n "Hello World!" | md5sum
    hash := GetMD5HashFromByte([]byte(str))
    
    if hash != expectedHash {
        t.Errorf("Expected %s, got %s", expectedHash, hash)
    }
}
