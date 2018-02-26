package command

import (
	"fmt"
	"github.com/kamilczerw/tokens-cli/lib/store"
	"github.com/zalando/go-keyring"
	"log"
	"errors"
	"encoding/base64"
)



type Remove struct {
	AppName string
	CopyMode bool
}

func (remove *Remove) Run(store store.Store) error {
	encoded, err := keyring.Get(remove.AppName, "tokens")
	if err != nil {
		log.Fatal(err)
		return err
	}

	password, err := GetSecret("Enter password: ", ErrEmptyPassword)
	if err != nil {
		return err
	}

	encrypted, _ := base64.StdEncoding.DecodeString(encoded)
	_, err = decrypt(encrypted, password)
	if err != nil {
		return err
	}

	err = keyring.Delete(remove.AppName, "tokens")
	if err != nil {
		return errors.New(fmt.Sprintf("cannot remove '%s' from keyring\ncouse: %s", remove.AppName, err))
	}

	err = store.RemoveDevice(remove.AppName)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot remove '%s'\ncouse: %s", remove.AppName, err))
	}

	fmt.Printf("remove '%s' has been sucesfully deleted\n", remove.AppName)
	return nil
}


type RemoveHelp struct {}

func (help *RemoveHelp) Run(store store.Store) error {

	fmt.Printf(`Remove mfadevice.
Usage:
  tokens rm DEVICE_NAME [options]

Options:
  -h, --help    Show help 
`)

	return nil
}
