package command

import (
	"fmt"
	"github.com/kamilczerw/tokens-cli/lib/store"
)


const (
	VERSION = "0.1.0"
)

type Version struct{}

func (l *Version) Run(store store.Store) error {
	fmt.Printf("Tokens v%s\n", VERSION)
	return nil
}