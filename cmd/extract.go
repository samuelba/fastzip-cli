package cmd

import (
	"context"
	"github.com/fatih/color"
	"github.com/saracen/fastzip"
	"github.com/spf13/cobra"
	"os"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract ZIP_FILE [DESTINATION]",
	Short: "Extract a ZIP archive.",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("extract called")
		archive := args[0]
		destination := "."
		if len(args) == 2 {
			destination = args[1]
		}

		if _, err := os.Stat(archive); os.IsNotExist(err) {
			color.Red("'%s' does not exist.", archive)
			os.Exit(1)
		}

		e, err := fastzip.NewExtractor(archive, destination)
		if err != nil {
			color.Red("Could not create extractor, error: %s", err)
			os.Exit(1)
		}
		defer e.Close()

		// Extract archive files
		if err = e.Extract(context.Background()); err != nil {
			color.Red("Could not extract '%s' to '%s', error: %s", archive, destination, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//extractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
