package cmd

import (
	"log"
	"yarastore/pkg/yarastore"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile YARA rules from various files into single ruleset",
	Run: func(cmd *cobra.Command, args []string) {
		compilerState, err := yarastore.NewCompilerState()
		if err != nil {
			log.Fatal(err)
		}

		if err := compilerState.ReadConfig(&conf); err != nil {
			log.Fatal(err)
		}

		if err := compilerState.Compile(); err != nil {
			log.Fatal(err)
		}

		output := viper.GetString("rules.output")
		if err := compilerState.Save(output); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	compileCmd.PersistentFlags().StringSliceP("files", "f", nil, "Files to read rules from")
	viper.BindPFlag("rules.files", compileCmd.PersistentFlags().Lookup("files"))

	compileCmd.PersistentFlags().StringSliceP("dirs", "d", nil, "Directories to read rules from")
	viper.BindPFlag("rules.dirs", compileCmd.PersistentFlags().Lookup("dirs"))

	compileCmd.PersistentFlags().StringSliceP("exclude", "x", nil, "Files or Directories (suffixed with /) to exclude")
	viper.BindPFlag("rules.exclude", compileCmd.PersistentFlags().Lookup("exclude"))

	compileCmd.PersistentFlags().StringP("include-pattern", "i", "", "Files matching this pattern will be read from")
	viper.BindPFlag("rules.include_pattern", compileCmd.PersistentFlags().Lookup("include-pattern"))

	compileCmd.PersistentFlags().StringP("exclude-pattern", "e", "", "Files matching this pattern will not be read from")
	viper.BindPFlag("rules.exclude_pattern", compileCmd.PersistentFlags().Lookup("exclude-pattern"))

	compileCmd.PersistentFlags().BoolP("recursive", "r", false, "Recursively search for files in directories")
	viper.BindPFlag("rules.recursive", compileCmd.PersistentFlags().Lookup("recursive"))

	compileCmd.PersistentFlags().StringP("output", "o", "ruleset", "Path to store the compiled ruleset in")
	viper.BindPFlag("rules.output", compileCmd.PersistentFlags().Lookup("output"))

	rootCmd.AddCommand(compileCmd)
}
