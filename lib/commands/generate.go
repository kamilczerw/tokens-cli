package command

import (
	"fmt"
	"github.com/kamilczerw/tokens-cli/lib/store"
	"github.com/zalando/go-keyring"
	"log"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"encoding/base64"
  "github.com/pquerna/otp/totp"
  "time"
)



type Generate struct {
	AppName string
	CopyMode bool
	QuietMode bool
}

func (generate *Generate) Run(store store.Store) error {
	encoded, err := keyring.Get(generate.AppName, "tokens")
	if err != nil {
		log.Fatal(err)
		return err
	}

	password, err := GetSecret("Enter password: ", ErrEmptyPassword)
	if err != nil {
		return err
	}

	encrypted, _ := base64.StdEncoding.DecodeString(encoded)
	secret, err := decrypt(encrypted, password)
	if err != nil {
		return err
	}

  t := time.Now()
  code, _ := totp.GenerateCode(secret, t)
  seconds := 30 - int32(t.Unix() % 30)

  if generate.QuietMode {
  	fmt.Print(code)
	} else {
		fmt.Printf("Token: %s\nExpires in: %d seconds\n", code, seconds)
	}

	return nil
}

func decrypt(message []byte, passphrase string) (string, error) {
	block, _ := pem.Decode(message)
	if block == nil {
		return "", errors.New("empty message")
	}
	if !x509.IsEncryptedPEMBlock(block) {
		return "", errors.New("block is not PEM block")

	}
	decrypted, err := x509.DecryptPEMBlock(block, []byte(passphrase))

	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}


type GenerateHelp struct {}

func (help *GenerateHelp) Run(store store.Store) error {

	fmt.Printf(`Generate one time code for the device.
Usage:
  tokens DEVICE_NAME [options]

Options:
  -c, --copy    Copy to clipboard 
  -q, --quiet   Quiet mode 
`)

	return nil
}
