package tokens

import (
	"fmt"
	"errors"
	"github.com/spf13/pflag"
	"github.com/kamilczerw/tokens-cli/lib/commands"
	"github.com/kamilczerw/tokens-cli/lib/store"
)

const (
	EX_USAGE_ERROR     = 64
	EX_DATA_ERROR      = 65
	EX_UNAVAILABLE     = 69
	EX_TEMPORARY_ERROR = 79
)

type ErrorWithExitCode struct {
	error
	ExitCode int
}

type Command interface {
	Run(store store.Store) error
}

func ParseArgs(args []string, store store.Store) (Command, error) {
	cmd, err := parseArgs(args, store)
	if err == pflag.ErrHelp {
		return nil, fmt.Errorf("arg0: %s", args[0])
		//if HelpAliases[args[0]] == "" {
		//	return parseHelpArgs(nil)
		//} else {
		//	return parseHelpArgs(args)
		//}
	}

	// If arguments fail to parse for any reason, it's a usage error
	if err != nil {
		if _, ok := err.(ErrorWithExitCode); !ok {
			err = ErrorWithExitCode{err, EX_USAGE_ERROR}
		}
	}

	return cmd, err
}

var (
	ErrTooManyArguments            = ErrorWithExitCode{errors.New("too many arguments provided"), EX_USAGE_ERROR}
	ErrNotEnoughArguments          = ErrorWithExitCode{errors.New("not enough arguments provided"), EX_USAGE_ERROR}
	//ErrVaultNameRequired           = ErrorWithExitCode{errors.New("A vault name must be specified"), EX_USAGE_ERROR}
	//ErrMixingCommandAndInteractive = ErrorWithExitCode{errors.New("Cannot mix an interactive shell with command arguments"), EX_USAGE_ERROR}
	//
	//ErrUnknownShell = errors.New("Unknown shell")
)

func parseArgs(args []string, store store.Store) (Command, error) {
	flag := spawnFlagSet()
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.Changed("version") {
		return &command.Version{}, nil
	}

	if flag.Changed("help") {
		return &command.Help{}, nil
	}

	// Parse command
	commandArgs := flag.Args()
	if len(commandArgs) == 0 || flag.ArgsLenAtDash() == 0 {
		return &command.Help{}, nil
	}

	if flag.ArgsLenAtDash() > -1 {
		commandArgsWithDash := append([]string{}, commandArgs[:flag.ArgsLenAtDash()]...)
		commandArgsWithDash = append(commandArgsWithDash, "--")
		commandArgsWithDash = append(commandArgsWithDash, commandArgs[flag.ArgsLenAtDash():]...)
		commandArgs = commandArgsWithDash
	}

	devices, err := store.ListDevices()

	switch commandArgs[0] {
	case "add":
		return parseAddArgs(commandArgs[1:])

	case "help":
		return &command.Help{}, nil

	case "ls":
		return parseListArgs(commandArgs[1:])

	case "version":
		return &command.Version{}, nil

	default:
		cmd := commandArgs[0]
		if stringInSlice(cmd, devices) {
			return parseGenerateArgs(commandArgs)
		}

		return nil, fmt.Errorf("'%s' is not a tokens command, or '%s' device doesn't exist.\n " +
			"Use: 'tokens add %s' to add the device", commandArgs[0], commandArgs[0], commandArgs[0])
	}
}

func spawnFlagSet() *pflag.FlagSet {
	flag := pflag.NewFlagSet("tokens", pflag.ContinueOnError)
	flag.Usage = func() {}
	flag.SetInterspersed(false)
	flag.StringP("name", "n", "", "Name of the vault to use")
	flag.BoolP("version", "v", false, "Specify current version of Tokens")
	flag.BoolP("help", "h", false, "Show help")
	return flag
}

func parseAddArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("add", pflag.ContinueOnError)
	flag.BoolP("help", "h", false, "Show help")
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.Changed("help") {
		return &command.AddHelp{}, nil
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	cmd := &command.Add{}
	cmd.AppName = flag.Arg(0)

	return cmd, nil
}

func parseListArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("ls", pflag.ContinueOnError)
	flag.BoolP("help", "h", false, "Show help")
	flag.BoolP("quiet", "q", false, "Quiet mode")
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.Changed("help") {
		return &command.ListHelp{}, nil
	}

	if flag.NArg() > 0 {
		return nil, ErrTooManyArguments
	}

	quiet, err := flag.GetBool("quiet")
	if err != nil {
		return nil, err
	}

	cmd := &command.List{}
	cmd.QuietMode = quiet

	return cmd, nil
}

func parseGenerateArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("generate", pflag.ContinueOnError)
	flag.BoolP("help", "h", false, "Show help")
	flag.BoolP("copy", "c", false, "Copy to clipboard")
	flag.Usage = func() {}
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.Changed("help") {
		return &command.GenerateHelp{}, nil
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	copyMode, err := flag.GetBool("copy")
	if err != nil {
		return nil, err
	}

	cmd := &command.Generate{}
	cmd.CopyMode = copyMode
	cmd.AppName = flag.Arg(0)

	return cmd, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}