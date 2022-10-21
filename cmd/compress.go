package cmd

import (
	"context"
	"github.com/fatih/color"
	"github.com/saracen/fastzip"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress SOURCE ZIP_FILE",
	Short: "Compress a folder/file to a ZIP archive.",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("compress called")
		source := args[0]
		sourceFile := args[0]
		archive := args[1]

		if _, err := os.Stat(source); os.IsNotExist(err) {
			color.Red("'%s' does not exist.", source)
			os.Exit(1)
		}

		// Create archive file
		w, err := os.Create(archive)
		if err != nil {
			color.Red("Could not create archive file '%s', error: %s", archive, err)
			os.Exit(1)
		}
		defer w.Close()

		// If source is a file, get the parent folder to use in the archiver.
		sourceInfo, err := os.Stat(source)
		if !sourceInfo.IsDir() {
			source = filepath.Dir(source)
		}

		// Create new Archiver
		a, err := fastzip.NewArchiver(w, source)
		if err != nil {
			color.Red("Could not create archiver, error: %s", err)
			os.Exit(1)
		}
		defer a.Close()

		// Register a non-default level compressor if required
		// a.RegisterCompressor(zip.Deflate, fastzip.FlateCompressor(1))

		files := make(map[string]os.FileInfo)
		if sourceInfo.IsDir() {
			// Walk directory, adding the files we want to add
			err = filepath.Walk(source, func(pathname string, info os.FileInfo, err error) error {
				files[pathname] = info
				return nil
			})
		} else {
			files[sourceFile] = sourceInfo
		}

		// Archive
		if err = a.Archive(context.Background(), files); err != nil {
			color.Red("Could not archive '%s' to '%s', error: %s", source, archive, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
