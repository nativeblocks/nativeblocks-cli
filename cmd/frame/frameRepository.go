package frame

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
)

const syncFrameMutation = `
  mutation syncFrame($input: SyncFrameInput!) {
    syncFrame(input: $input) {
      id
    }
  }
`

const getFrameQuery = `
  query frame($route: String!) {
    frame(route: $route) {
      id
      name
      route
      isStarter
      type
      variables {
        key
        value
        type
      }
      blocks {
        id
        parentId
        slot
        keyType
        key
        visibilityKey
        position
        properties {
          key
          valueDesktop
          valueMobile
          valueTablet
          type
        }
        data {
          key
          value
          type
        }
        slots {
          slot
        }
      }
      actions {
        key
        event
        triggers {
          id
          parentId
          keyType
          then
          name
          properties {
            key
            value
            type
          }
          data {
            key
            value
            type
          }
        }
      }
    }
  }
`

func pushFrame(output FrameProductionDataWrapper, regionUrl string, accessToken string, apiKey string) error {

	if output.Data.FrameProduction.Id == "" {
		return errors.New("could not genereate frame, please check your input")
	}

	client := graphqlutil.NewClient()

	jsonBytes, _ := json.Marshal(output.Data.FrameProduction)
	input := map[string]interface{}{
		"route":     output.Data.FrameProduction.Route,
		"frameJson": string(jsonBytes),
	}

	variables := map[string]interface{}{
		"input": input,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Api-Key":       "Bearer " + apiKey,
	}

	_, err := client.Execute(
		regionUrl,
		headers,
		syncFrameMutation,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	return nil
}

func pullFrame(fm fileutil.FileManager, regionUrl string, accessToken string, apiKey string, fileName string, schema string, route string) error {
	client := graphqlutil.NewClient()

	variables := map[string]interface{}{
		"route": route,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Api-Key":       "Bearer " + apiKey,
	}

	apiResponse, err := client.Execute(
		regionUrl,
		headers,
		getFrameQuery,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	var frameResponse FrameWrapper
	err = graphqlutil.Parse(apiResponse, &frameResponse)
	if err != nil {
		return err
	}

	frame := mapFrameModelToDSL(frameResponse.Frame, schema)
	if frame.Route == "" {
		return fmt.Errorf("could not find frame route %v", frame.Route)
	}
	if err := fm.SaveToFile(fileName, frame); err != nil {
		return err
	}

	return nil
}
