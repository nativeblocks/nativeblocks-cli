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
	// cmd.AddCommand(pushCommand())
	// cmd.AddCommand(pullCommand())
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

// func pushCommand() *cobra.Command {
// 	var directory string
// 	cmd := &cobra.Command{
// 		Use:   "push",
// 		Short: "Push a frame",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			baseFm, err := fileutil.NewFileManager(nil)
// 			if err != nil {
// 				return err
// 			}

// 			var regionConfig RegionConfig
// 			if err := baseFm.LoadFromFile(RegionFileName, &regionConfig); err != nil {
// 				return fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
// 			}

// 			var authConfig AuthConfig
// 			if err := baseFm.LoadFromFile(AuthFileName, &authConfig); err != nil {
// 				return fmt.Errorf("not authenticated. Please login first using 'nativeblocks auth'")
// 			}

// 			var projectConfig ProjectConfig
// 			if err := baseFm.LoadFromFile(ProjectFileName, &projectConfig); err != nil {
// 				return fmt.Errorf("project not set. Please select a project first using 'nativeblocks project set'")
// 			}

// 			baseDir := fileutil.GetFileDir(directory)
// 			fileName := fileutil.GetFileName(directory)

// 			fm, err := fileutil.NewFileManager(&baseDir)
// 			if err != nil {
// 				return err
// 			}

// 			fileExists := fm.FileExists(fileName)
// 			if !fileExists {
// 				return fmt.Errorf("could not find the file under: %v", directory)
// 			}

// 			var jsonDSL FrameDSLModel
// 			fileError := fm.LoadFromFile(fileName, &jsonDSL)
// 			if fileError != nil {

// 			}
// 			output, err := generateFrame(jsonDSL)
// 			if err != nil {
// 				return err
// 			}

// 			if output.ID == "" {
// 				return nil
// 			}

// 			client := graphqlutil.NewClient()

// 			jsonBytes, _ := json.Marshal(output)
// 			input := map[string]interface{}{
// 				"route":     jsonDSL.Route,
// 				"frameJson": string(jsonBytes),
// 			}

// 			variables := map[string]interface{}{
// 				"input": input,
// 			}

// 			headers := map[string]string{
// 				"Authorization": "Bearer " + authConfig.AccessToken,
// 				"Api-Key":       "Bearer " + projectConfig.APIKey,
// 			}

// 			resp, err := client.Execute(
// 				regionConfig.URL,
// 				headers,
// 				syncFrameMutation,
// 				variables,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("sync failed: %v", err)
// 			}

// 			json.Marshal(resp.Data)
// 			fmt.Printf("Frame successfully synced \n")

// 			return nil
// 		},
// 	}

// 	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Frame working directory")
// 	cmd.MarkFlagRequired("directory")

// 	return cmd
// }

// func pullCommand() *cobra.Command {
// 	var directory string
// 	cmd := &cobra.Command{
// 		Use:   "pull",
// 		Short: "Pull a frame",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			baseFm, err := fileutil.NewFileManager(nil)
// 			if err != nil {
// 				return err
// 			}

// 			var regionConfig RegionConfig
// 			if err := baseFm.LoadFromFile(RegionFileName, &regionConfig); err != nil {
// 				return fmt.Errorf("region not set. Please set region first using 'nativeblocks region set <url>'")
// 			}

// 			var authConfig AuthConfig
// 			if err := baseFm.LoadFromFile(AuthFileName, &authConfig); err != nil {
// 				return fmt.Errorf("not authenticated. Please login first using 'nativeblocks auth'")
// 			}

// 			var projectConfig ProjectConfig
// 			if err := baseFm.LoadFromFile(ProjectFileName, &projectConfig); err != nil {
// 				return fmt.Errorf("project not set. Please select a project first using 'nativeblocks project set'")
// 			}

// 			baseDir := fileutil.GetFileDir(directory)
// 			fileName := fileutil.GetFileName(directory)

// 			fm, err := fileutil.NewFileManager(&baseDir)
// 			if err != nil {
// 				return err
// 			}

// 			fileExists := fm.FileExists(fileName)
// 			if !fileExists {
// 				return fmt.Errorf("could not find the file under: %v", directory)
// 			}

// 			var jsonDSL FrameDSLModel
// 			fileError := fm.LoadFromFile(fileName, &jsonDSL)
// 			if fileError != nil {

// 			}

// 			if jsonDSL.Route == "" {
// 				return fmt.Errorf("could not find frame route")
// 			}

// 			client := graphqlutil.NewClient()

// 			variables := map[string]interface{}{
// 				"route": jsonDSL.Route,
// 			}

// 			headers := map[string]string{
// 				"Authorization": "Bearer " + authConfig.AccessToken,
// 				"Api-Key":       "Bearer " + projectConfig.APIKey,
// 			}

// 			resp, err := client.Execute(
// 				regionConfig.URL,
// 				headers,
// 				getFrameQuery,
// 				variables,
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to process response: %v", err)
// 			}
// 			responseData, err := json.Marshal(resp.Data)
// 			if err != nil {
// 				return fmt.Errorf("failed to process response: %v", err)
// 			}

// 			var frameResponse FrameWrapper
// 			if err := json.Unmarshal(responseData, &frameResponse); err != nil {
// 				fmt.Printf("Debug - Raw response: %s\n", string(responseData))
// 				return fmt.Errorf("failed to parse organizations response: %v", err)
// 			}

// 			frame := mapFrameModelToDSL(frameResponse.Frame)
// 			if frame.Route == "" {
// 				return fmt.Errorf("could not find frame route %v", frame.Route)
// 			}
// 			if err := fm.SaveToFile(fileName, frame); err != nil {
// 				return err
// 			}

// 			fmt.Printf("Frame successfully synced \n")

// 			return nil
// 		},
// 	}

// 	cmd.Flags().StringVarP(&directory, "directory", "d", "", "Frame working directory")
// 	cmd.MarkFlagRequired("directory")

// 	return cmd
// }
