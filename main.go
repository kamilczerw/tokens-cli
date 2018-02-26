package main

import (
	"os"
	"github.com/kamilczerw/tokens-cli/lib"
  "fmt"
  "github.com/kamilczerw/tokens-cli/lib/store"
  "log"
)


func main() {
  fileStore, err := store.NewFileStore()
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
	command, err := tokens.ParseArgs(os.Args[1:], fileStore)

  if err == nil {
    err = command.Run(fileStore)
  }

	if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
	}
}
