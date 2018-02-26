package command

import (
  "fmt"
  "github.com/kamilczerw/tokens-cli/lib/store"
  "strings"
)

type Help struct {}

func (e *Help) Run(store store.Store) error {

  devices, err := store.ListDevices()
  if err != nil {
    return err
  }

  fmt.Printf(`Tokens is a cli app for generating one time codes. 
Usage:
  tokens COMMAND
  tokens DEVICE_NAME

Commands:
  add    Add new mfa device
  rm     Remove mfa device
  ls     List mfa devices

Apps:
  - ` +
  strings.Join(devices, "\n  - ") +
  `

Run 'tokens COMMAND --help' for more information on a command.
`)

  return nil
}
