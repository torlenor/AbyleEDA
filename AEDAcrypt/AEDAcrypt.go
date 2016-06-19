package AEDAcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("example")

type cryptcfg struct {
	key   []byte
	nonce []byte
}

var ccfg cryptcfg

func init() {
	// TODO: IMPLEMENT SOME KIND OF AUTHENTICATION FOR NONCE AND MAKE KEY CONFIGURABLE
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	ccfg.key = []byte("AES256Key-32Characters1234567890")
	ccfg.nonce, _ = hex.DecodeString("bb8ef84243d2ee95a41c6c57")
}

func Encrypter(encmsg []byte) []byte {
	plaintext := encmsg

	block, err := aes.NewCipher(ccfg.key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, ccfg.nonce, plaintext, nil)

	return ciphertext
}

func Decrypter(cryptmsg []byte) ([]byte, error) {
	block, err := aes.NewCipher(ccfg.key)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, ccfg.nonce, cryptmsg, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return plaintext, nil
}

func CheckError(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(0)
	}
}

func GetMD5HashFromString(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetMD5HashFromByte(by []byte) string {
	hasher := md5.New()
	hasher.Write(by)
	return hex.EncodeToString(hasher.Sum(nil))
}
