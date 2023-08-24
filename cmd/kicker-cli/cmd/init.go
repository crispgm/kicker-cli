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
	initPath string
	orgName  string
)

func init() {
	initCmd.Flags().StringVarP(&initPath, "path", "p", ".", "Path to folder")
	initCmd.Flags().StringVarP(&orgName, "name", "n", "Foosball", "Organization name")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a folder path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(initPath) == 0 || initPath == "." {
			pterm.Info.Println("Initializing current folder...")
		}
		// check `.kicker.yaml`
		instance := app.NewApp(initPath, app.DefaultName)
		err := instance.LoadConf()
		if err != nil {
			pterm.Info.Printfln("Initializing `%s` ...", instance.Path)
			pterm.Info.Printfln("Creating `%s` ...", instance.Name)
			instance.Conf.Organization = *entity.NewOrganization(orgName)
			err = instance.WriteConf()
			if err != nil {
				pterm.Error.Println(err)
				os.Exit(1)
			}
			err = os.Mkdir(filepath.Join(instance.Path, "data"), os.FileMode(0o755))
			if err != nil {
				pterm.Error.Println(err)
				os.Exit(1)
			}
		} else {
			pterm.Error.Println("Found existing `.kicker.yaml`")
			os.Exit(1)
		}
	},
}
