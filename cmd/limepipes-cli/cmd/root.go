package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "limepipes-cli",
	Short: "A command line tool managing bagpipe tunes",
	Long: `Some things are cumbersome doing with the REST API like importing many tunes in 
one run. With this command line tool it is possible to things like that.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	opts := &Options{}
	rootCmd.AddCommand(NewParseCmd(opts))
	rootCmd.AddCommand(NewImportCmd(opts))
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
