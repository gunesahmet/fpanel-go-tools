package CrypterWin

import (
	"Helpers"
	"encoding/base64"
	"github.com/gunesahmet/openssl"
	"regexp"
)

type String string

type Crypter struct {
	key []byte
	iv  []byte
}

func NewCrypter(key []byte, iv []byte) (*Crypter, error) {
	return &Crypter{key, iv}, nil
}

func (s *String) replaceAll(find string, replace string) *String {
	re, _ := regexp.Compile(find)
	*s = String(re.ReplaceAllString(string(*s), replace))
	return s
}

func (c *Crypter) Encrypt(input []byte) ([]byte, error) {
	enc, err := openssl.AesCFBEncrypt(input, c.key, c.iv, openssl.ZEROS_PADDING)
	if err != nil {
		return nil, err
	}
	return enc, nil
}
func (c *Crypter) Decrypt(input []byte) ([]byte, error) {
	dec, err := openssl.AesCFBDecrypt(input, c.key, c.iv, openssl.ZEROS_PADDING)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func safeBase64Encode(str string) string {
	data := []byte(str)
	s := base64.StdEncoding.EncodeToString(data)
	s = Helpers.ReplaceAllOld(s, "+", "-")
	s = Helpers.ReplaceAllOld(s, "/", "_")
	s = Helpers.RemoveAll(s, "=")
	return s
}
func safeBase64Decode(str string) string {
	var result string
	s := String(str)
	s.replaceAll("-", "+")
	s.replaceAll("_", "/")
	result = string(s)
	strLen := len(s)
	mod4 := strLen % 4
	comma := "===="
	if mod4 > 0 {
		result += comma[mod4:]
	}

	//fmt.Println(result)
	return result
}

func (c *Crypter) EncodeToString(str string) string {
	encrypted, err := c.Encrypt([]byte(str))
	if err != nil {
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	return encoded
}

func (c *Crypter) DecodeToString(str string) string {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		//Logger.Error("Decode status is false. ", "1005")
		return ""
	}
	decrypted, err := c.Decrypt(decoded)
	if err != nil {
		//Logger.Error("Decrypt status is fail.", "1006")
		return ""
	}
	return string(decrypted)
}

func (c *Crypter) Encode(data string) (string, bool) {
	encrypted, err := c.Encrypt([]byte(data))
	if err != nil {
		return "Error enc", false
	}
	encoded := safeBase64Encode(string(encrypted))
	return encoded, true
}

func DecodeBase64(str string) ([]byte, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte(nil), false
	}
	return data, true
}
func EncodeBase64(str string) (string, bool) {
	data := base64.StdEncoding.EncodeToString([]byte(str))
	return data, true
}

func (c *Crypter) Decode(data string) (string, bool) {
	dsb := safeBase64Decode(data)
	decoded, status := DecodeBase64(dsb)

	if status != true {
		return "", false
	}
	decrypted, err := c.Decrypt(decoded)
	if err != nil {
		return "", false
	}
	return string(decrypted), true
}
