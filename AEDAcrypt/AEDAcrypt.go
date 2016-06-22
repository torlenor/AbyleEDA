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

type CryptCfg struct {
	Key   []byte
	Nonce []byte
}

var ccfg CryptCfg

func init() {
	// TODO: IMPLEMENT SOME KIND OF AUTHENTICATION FOR NONCE AND MAKE KEY CONFIGURABLE
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	ccfg.Key = []byte("AES256Key-32Characters1234567890")
	ccfg.Nonce, _ = hex.DecodeString("bb8ef84243d2ee95a41c6c57")
}

func Encrypter(encmsg []byte, ccfg CryptCfg) []byte {
	plaintext := encmsg

	block, err := aes.NewCipher(ccfg.Key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, ccfg.Nonce, plaintext, nil)

	return ciphertext
}

func Decrypter(cryptmsg []byte) ([]byte, error) {
	block, err := aes.NewCipher(ccfg.Key)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, ccfg.Nonce, cryptmsg, nil)
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
