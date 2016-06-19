package AEDAcrypt

import "testing"

func TestEncrypter(t *testing.T) {
	var str string = "Hello World!"
	strenc := Encrypter([]byte(str))
	_, err := Decrypter(strenc)

	if err != nil {
		t.Error("Expected 1.5, got ")
	}
}

func TestDecrypter(t *testing.T) {
	var str string = "Hello World!"
	strenc := Encrypter([]byte(str))
	_, err := Decrypter(strenc)

	if err != nil {
		t.Error("Expected 1.5, got ")
	}
}
