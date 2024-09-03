package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes/cmd/limepipes-cli/importtype"
	"strings"
)

const DefaultOutputDir = "./parser_success"

// Options that can be passed to the command via command line flags
type Options struct {
	Recursive       bool
	DryRun          bool
	ImportTypes     []string
	SkipFailedFiles bool
	Verbose         bool
	OutputDir       string
}

func addImportFileTypes(cmd *cobra.Command, opts *Options) {
	cmd.Flags().StringSliceVarP(&opts.ImportTypes,
		"import-file-types",
		"i",
		[]string{
			importtype.FromFileFormat(fileformat.Format_BWW),
		},
		fmt.Sprintf("only import files of a specific type. Can be specified with a comma separated list like "+
			"-i=bww,xml. Valid file types are (%s).", strings.Join(importtype.AllTypes(), ", ")),
	)
}

func addRecursive(cmd *cobra.Command, opts *Options) {
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false,
		"import files recursively for all given directory paths",
	)
}

func addDryRun(cmd *cobra.Command, opts *Options) {
	cmd.Flags().BoolVarP(&opts.DryRun, "dry-run", "d", false,
		"import files but only convert them and print infos. Don't import them to the database and don't stop on parsing errors.",
	)
}

func addVerbose(cmd *cobra.Command, opts *Options) {
	cmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false,
		"verbose output",
	)
}

func addSkipFailedFiles(cmd *cobra.Command, opts *Options) {
	cmd.PersistentFlags().BoolVarP(&opts.SkipFailedFiles, "skip-failed", "s", false,
		"only print error when file couldn't be parsed. If flag is false, the command returns an error on first failed parsing",
	)
}

func addOutputDir(cmd *cobra.Command, opts *Options) {
	cmd.Flags().StringVarP(
		&opts.OutputDir,
		"output-dir",
		"o",
		DefaultOutputDir,
		"Output directory where to move the successful parsed files into",
	)
}
