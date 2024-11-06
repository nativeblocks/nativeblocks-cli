package frame

import (
	"encoding/json"
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
	"github.com/spf13/cobra"
)

const (
	ProjectFileName = "project"
	RegionFileName  = "region"
	AuthFileName    = "auth"
)

type RegionConfig struct {
	URL string `json:"url"`
}

type AuthConfig struct {
	AccessToken string `json:"accessToken"`
}

type ProjectConfig struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	APIKey  string `json:"apiKey"`
	KeyName string `json:"keyName"`
}

const syncFrameMutation = `
  mutation syncFrame($input: SyncFrameInput!) {
    syncFrame(input: $input) {
      id
    }
  }
`

func FrameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frame",
		Short: "Manage frames",
	}

	cmd.AddCommand(genCommand())
	cmd.AddCommand(pushCommand())
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

			var regionConfig RegionConfig
			if err := baseFm.LoadFromFile(RegionFileName, &regionConfig); err != nil {
				return fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
			}

			var authConfig AuthConfig
			if err := baseFm.LoadFromFile(AuthFileName, &authConfig); err != nil {
				return fmt.Errorf("not authenticated. Please login first using 'nativeblocks auth'")
			}

			var projectConfig ProjectConfig
			if err := baseFm.LoadFromFile(ProjectFileName, &projectConfig); err != nil {
				return fmt.Errorf("project not set. Please select a project first using 'nativeblocks project list'")
			}

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

			client := graphqlutil.NewClient()

			jsonBytes, _ := json.Marshal(output)
			input := map[string]interface{}{
				"route":     jsonDSL.Route,
				"frameJson": string(jsonBytes),
			}

			variables := map[string]interface{}{
				"input": input,
			}

			headers := map[string]string{
				"Authorization": "Bearer " + authConfig.AccessToken,
				"Api-Key":       "Bearer " + projectConfig.APIKey,
			}

			resp, err := client.Execute(
				regionConfig.URL,
				headers,
				syncFrameMutation,
				variables,
			)
			if err != nil {
				return fmt.Errorf("sync failed: %v", err)
			}

			json.Marshal(resp.Data)
			fmt.Printf("Frame successfully synced \n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Frame working directory")
	cmd.MarkFlagRequired("directory")

	return cmd
}
