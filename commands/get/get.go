package get

import (
	"flag"
	"github.com/franela/vault/gpg"
	"github.com/franela/vault/ui"
	"github.com/franela/vault/vault"
	"github.com/mitchellh/cli"
	"path"
)

const getHelpText = `
`

func Factory() (cli.Command, error) {
	return getCommand{}, nil
}

type getCommand struct {
}

func (getCommand) Help() string {
	return getHelpText
}

func (getCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("get", flag.ContinueOnError)

	var outputFile string

	cmdFlags.StringVar(&outputFile, "o", "", "specify the output file to store decrypted text")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()
	file := args[0]

	if len(outputFile) > 0 {
		if err := gpg.DecryptFile(outputFile, path.Join(vault.GetHomeDir(), file)); err != nil {
			ui.Printf("Error decrypting file %s %s", file, err)
			return 1
		}
	} else {
		if text, err := gpg.Decrypt(path.Join(vault.GetHomeDir(), file)); err != nil {
			ui.Printf("Error decrypting file %s %s", file, err)
			return 1
		} else {
			ui.Printf("%s", text)
		}
	}
	return 0
}

func (getCommand) Synopsis() string {
	return ""
}
