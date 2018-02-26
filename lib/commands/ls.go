package command

import (
  "fmt"
  "github.com/kamilczerw/tokens-cli/lib/store"
  "log"
  "strings"
)



type List struct {
  QuietMode bool
}

func (list *List) Run(store store.Store) error {
  devices, err := store.ListDevices()
  if err != nil {
    log.Fatal(err)
    return err
  }

  separator := "\n"
  if !list.QuietMode {
    fmt.Print("MFA devices: \n - ")
    separator = "\n - "
  }
  fmt.Println(strings.Join(devices, separator))

  return nil
}


type ListHelp struct {}

func (help *ListHelp) Run(store store.Store) error {

  fmt.Printf(`List all mfa devices.
Usage:
  tokens ls [options]

Options:
  -q, --quiet    Quiet mode 
`)

  return nil
}
