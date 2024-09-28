package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"yarastore/pkg/yarastore"
)

var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Match compiled ruleset against files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		compilerState, err := yarastore.NewCompilerState()
		if err != nil {
			log.Fatal(err)
		}

		ruleset := args[0]
		if err := compilerState.ReadCompiled(ruleset); err != nil {
			log.Fatal(err)
		}

		matches, err := compilerState.MatchConfig(&conf)
		if err != nil {
			log.Fatal(err)
		}

		output := viper.GetString("target.output")
		yarastore.RuleMatchAsJson(matches, output)
	},
}

func init() {
	matchCmd.PersistentFlags().StringSliceP("files", "f", nil, "Files to match rules against")
	viper.BindPFlag("target.files", matchCmd.PersistentFlags().Lookup("files"))

	matchCmd.PersistentFlags().StringSliceP("dirs", "d", nil, "Directories to match rules against")
	viper.BindPFlag("target.dirs", matchCmd.PersistentFlags().Lookup("dirs"))

	matchCmd.PersistentFlags().StringSliceP("exclude", "x", nil, "Files or Directories (suffixed with /) to exclude")
	viper.BindPFlag("target.exclude", matchCmd.PersistentFlags().Lookup("exclude"))

	matchCmd.PersistentFlags().StringP("include-pattern", "i", "", "Files matching this pattern will be matched against")
	viper.BindPFlag("target.include_pattern", matchCmd.PersistentFlags().Lookup("include-pattern"))

	matchCmd.PersistentFlags().StringP("exclude-pattern", "e", "", "Files matching this pattern will not be matched against")
	viper.BindPFlag("target.exclude_pattern", matchCmd.PersistentFlags().Lookup("exclude-pattern"))

	matchCmd.PersistentFlags().BoolP("recursive", "r", false, "Recursively search for files in directories")
	viper.BindPFlag("target.recursive", matchCmd.PersistentFlags().Lookup("recursive"))

	matchCmd.PersistentFlags().StringP("output", "o", "yarastore_report.json", "Path to store the report in")
	viper.BindPFlag("target.output", matchCmd.PersistentFlags().Lookup("output"))

	rootCmd.AddCommand(matchCmd)
}
