package frame

import (
	"encoding/json"
	"fmt"

	"github.com/nativeblocks/cli/cmd/auth"
	"github.com/nativeblocks/cli/cmd/project"
	"github.com/nativeblocks/cli/cmd/region"
	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/spf13/cobra"
)

func FrameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frame",
		Short: "Manage frames",
	}

	cmd.AddCommand(genCommand())
	cmd.AddCommand(pushCommand())
	cmd.AddCommand(pullCommand())
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
				return fmt.Errorf("could not find the file under: %v", directory)
			}

			var jsonDSL FrameDSLModel
			err = fm.LoadFromFile(fileName, &jsonDSL)
			if err != nil {
				return err
			}

			output, err := generateFrame(jsonDSL)
			if err != nil {
				return err
			}

			if output.Data.FrameProduction.Id == "" {
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

func pushCommand() *cobra.Command {
	var directory string
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push a frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseFm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			region, err := region.GetRegion(*baseFm)
			if err != nil {
				return err
			}

			auth, err := auth.AuthGet(*baseFm)
			if err != nil {
				return err
			}

			project, err := project.GetProject(*baseFm)
			if err != nil {
				return err
			}

			baseDir := fileutil.GetFileDir(directory)
			fileName := fileutil.GetFileName(directory)

			inputFm, err := fileutil.NewFileManager(&baseDir)
			if err != nil {
				return err
			}

			fileExists := inputFm.FileExists(fileName)
			if !fileExists {
				return fmt.Errorf("could not find the file under: %v", directory)
			}

			var jsonDSL FrameDSLModel
			err = inputFm.LoadFromFile(fileName, &jsonDSL)
			if err != nil {
				return err
			}

			output, err := generateFrame(jsonDSL)
			if err != nil {
				return err
			}

			err = pushFrame(output, region.Url, auth.AccessToken, project.APIKeys[0].APIKey)
			if err != nil {
				return err
			}

			fmt.Printf("Frame successfully synced \n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Frame working directory")
	cmd.MarkFlagRequired("directory")

	return cmd
}

func pullCommand() *cobra.Command {
	var directory string
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull a frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseFm, err := fileutil.NewFileManager(nil)
			if err != nil {
				return err
			}

			region, err := region.GetRegion(*baseFm)
			if err != nil {
				return err
			}

			auth, err := auth.AuthGet(*baseFm)
			if err != nil {
				return err
			}

			project, err := project.GetProject(*baseFm)
			if err != nil {
				return err
			}

			baseDir := fileutil.GetFileDir(directory)
			fileName := fileutil.GetFileName(directory)

			inputFm, err := fileutil.NewFileManager(&baseDir)
			if err != nil {
				return err
			}

			fileExists := inputFm.FileExists(fileName)
			if !fileExists {
				return fmt.Errorf("could not find the file under: %v", directory)
			}

			var jsonDSL FrameDSLModel
			err = inputFm.LoadFromFile(fileName, &jsonDSL)
			if err != nil {
				return err
			}

			if jsonDSL.Route == "" {
				return fmt.Errorf("could not find frame route")
			}

			err = pullFrame(*inputFm, region.Url, auth.AccessToken, project.APIKeys[0].APIKey, fileName, jsonDSL.Schema, jsonDSL.Route)
			if err != nil {
				return err
			}

			fmt.Printf("Frame successfully synced \n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Frame working directory")
	cmd.MarkFlagRequired("directory")

	return cmd
}
