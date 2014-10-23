package set

import (
	"flag"
	"github.com/franela/vault/gpg"
	"github.com/franela/vault/ui"
	"github.com/franela/vault/vault"
	"github.com/mitchellh/cli"
	"path"
	"path/filepath"
)

const setHelpText = `
Usage: vault set [options] vaultpath [text]

   Sets something into the vault. Either plain text or files using the -f 
   parameter are allowed.
`

func Factory() (cli.Command, error) {
	return setCommand{}, nil
}

type setCommand struct {
}

func (setCommand) Help() string {
	return setHelpText
}

func (setCommand) Run(args []string) int {

	vaultFile, err := vault.LoadVaultfile()

	if err != nil {
		ui.Printf("%s", err)
		return 1
	}

	if len(vaultFile.Recipients) == 0 {
		ui.Printf("Cannot set in vault if Vaultfile has no recipients. Please add at least one recipient.\n")
		return 3
	}

	cmdFlags := flag.NewFlagSet("set", flag.ContinueOnError)

	var fileName string

	cmdFlags.StringVar(&fileName, "f", "", "specify the file to encrypt")

	if err := cmdFlags.Parse(args); err != nil {
		ui.Printf(setHelpText)
		return 1
	}

	args = cmdFlags.Args()

	if len(args) == 0 || len(args) > 2 {
		ui.Printf(setHelpText)
		return 1
	}

	if len(fileName) > 0 {
		vaultPath := args[0]
		if filepath.Ext(vaultPath) != ".asc" {
			vaultPath = vaultPath + ".asc"
		}
		err := gpg.EncryptFile(path.Join(vault.GetHomeDir(), vaultPath), fileName, vaultFile.Recipients)
		if err != nil {
			ui.Printf("%s", err)
			return 1
		}
	} else {
		if len(args) != 2 {
			ui.Printf(setHelpText)
			return 1
		}

		text := args[0]
		vaultPath := args[1]
		if filepath.Ext(vaultPath) != ".asc" {
			vaultPath = vaultPath + ".asc"
		}

		err := gpg.Encrypt(path.Join(vault.GetHomeDir(), vaultPath), text, vaultFile.Recipients)

		if err != nil {
			ui.Printf("%s", err)
			return 1
		}
	}

	return 0
}

func (setCommand) Synopsis() string {
	return "Store files or plain-text intro your vault."
}
