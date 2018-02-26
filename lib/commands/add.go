package command

import (
	"github.com/fatih/color"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/zalando/go-keyring"
	"syscall"
	"crypto/x509"
	"os"
	"crypto/rand"
	"github.com/kamilczerw/tokens-cli/lib/store"
	"log"
	"encoding/pem"
	"encoding/base64"
	"strings"
)

var (
	green = color.New(color.FgGreen)
	cyan  = color.New(color.FgCyan)
	blue  = color.New(color.FgBlue)

	faintColor   = color.New(color.Faint)
	menuColor    = color.New(color.FgHiBlue)
	warningColor = color.New(color.FgHiYellow)

	ErrExists      = errors.New("App exists.")
	ErrEmptySecret = errors.New("Please provide the secret")
	ErrEmptyPassword = errors.New("Please provide the password")
	ErrNotConfirmedPassword = errors.New("Provided passwords was not the same")
)

type Add struct {
	AppName  string
}

type Creds struct {
	Secret string
	Password string
}



func (add *Add) Run(store store.Store) error {
	creds, err := GetSecrets()
	if err != nil {
		log.Fatal(err)
		return err
	}

	encrypted, err := encrypt(creds)
	if err != nil {
		log.Fatal(err)
		return err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(encrypted))

	err = keyring.Set(add.AppName, "tokens", encoded)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = store.AddDevice(add.AppName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("App '%s' successfully saved!\n", add.AppName)

	return nil
}

func GetSecrets() (Creds, error)  {
	creds := Creds{}

	secret, err := GetSecret("Enter mfa device secret: ", ErrEmptySecret)
	if err != nil {
		return creds, err
	}
	password, err := GetSecret("Enter password: ", ErrEmptyPassword)
	if err != nil {
		return creds, err
	}
	confirmPassword, err := GetSecret("Confirm password: ", ErrEmptyPassword)
	if err != nil {
		return creds, err
	}

	if password != confirmPassword {
		return creds, ErrNotConfirmedPassword
	}

	creds.Secret = secret
	creds.Password = password

	return creds, nil
}

func GetSecret(message string, err error) (string, error) {
	fmt.Print(message)
	byteSecret, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\r %s \r", strings.Repeat(" ", len(message)))

	if err != nil {
		return "", err
	}

	secret := string(byteSecret)

	if len(secret) == 0 {
		return "", err
	}

	return secret, nil
}

func encrypt(creds Creds) (string, error) {
	blockType := "RSA PRIVATE KEY"

	cipherType := x509.PEMCipherAES256

	EncryptedPEMBlock, err := x509.EncryptPEMBlock(rand.Reader,
		blockType,
		[]byte(creds.Secret),
		[]byte(creds.Password),
		cipherType)

	if err != nil {
		return "", err
	}

	// check if encryption is successful or not

	if !x509.IsEncryptedPEMBlock(EncryptedPEMBlock) {
		fmt.Println("PEM Block is not encrypted!")
		os.Exit(1)
	}


	if EncryptedPEMBlock.Type != blockType {
		fmt.Println("Block type is wrong!")
		os.Exit(1)
	}

	encoded := pem.EncodeToMemory(EncryptedPEMBlock)

	return string(encoded), nil
}

type AddHelp struct {}

func (e *AddHelp) Run(store store.Store) error {

	fmt.Printf(`Add new mfa device to generate one time codes.
Usage:
  tokens add DEVICE_NAME
`)

	return nil
}


