package AEDAcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("AEDAlogger")

type CryptCfg struct {
	Key   []byte
	Nonce []byte
}

var ccfg CryptCfg

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

func Decrypter(cryptmsg []byte, ccfg CryptCfg) ([]byte, error) {
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

func GetMD5HashFromString(text string) string {
    return getMD5Hash([]byte(text))
}

func GetMD5HashFromByte(by []byte) string {
    return getMD5Hash(by)
}

func getMD5Hash(by []byte) string {
    hasher := md5.New()
    hasher.Write(by)
    return hex.EncodeToString(hasher.Sum(nil))
}
