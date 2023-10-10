package cmd

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
	"github.com/crispgm/kicker-cli/internal/entity"
)

var (
	initPath    string
	initOrgName string
)

func init() {
	initCmd.Flags().StringVarP(&initPath, "path", "p", ".", "path to folder")
	initCmd.Flags().StringVarP(&initOrgName, "name", "n", "Foosball", "organization name")
	_ = initCmd.MarkFlagDirname("path")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a folder path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(initPath) == 0 || initPath == "." {
			pterm.Println("Initializing current folder...")
		}
		// check `.kicker.yaml`
		instance := app.NewApp(initPath, app.DefaultName)
		err := instance.LoadConf()
		if err != nil {
			pterm.Printfln("Initializing `%s` ...", instance.Path)
			pterm.Printfln("Creating `%s` ...", instance.Name)
			instance.Conf.Organization = *entity.NewOrganization(initOrgName)
			err = instance.WriteConf()
			if err != nil {
				errorMessageAndExit(err)
			}
			err = os.Mkdir(filepath.Join(instance.Path, "data"), os.FileMode(0o755))
			if err != nil {
				errorMessageAndExit(err)
			}
		} else {
			errorMessageAndExit("Found existing `.kicker.yaml`")
		}
	},
}
