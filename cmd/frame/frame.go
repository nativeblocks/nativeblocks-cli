package frame

import (
	"encoding/json"
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

func FrameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frame",
		Short: "Manage frames",
	}

	cmd.AddCommand(genCommand())
	return cmd
}

func genCommand() *cobra.Command {
	var directory string
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate a frame",
		RunE: func(cmd *cobra.Command, args []string) error {

			baseDir := fileutil.GetFileDir(directory)
			fileName := fileutil.GetFileName(directory)

			fm, err := fileutil.NewFileManager(&baseDir)
			if err != nil {
				return err
			}

			fileExists := fm.FileExists(fileName)
			if !fileExists {
				return fmt.Errorf("Could not find the file under: %v", directory)
			}

			var jsonDSL FrameDSLModel
			fileError := fm.LoadFromFile(fileName, &jsonDSL)
			if fileError != nil {

			}
			output, err := generateFrame(jsonDSL)
			if err != nil {
				return err
			}

			if output.ID == "" {
				return nil
			}

			frameJson, err := json.Marshal(output)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(frameJson))

			return nil
		},
	}

	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Frame working directory")
	cmd.MarkFlagRequired("directory")

	return cmd
}
