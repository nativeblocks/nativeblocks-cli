package integration

import (
	"errors"
	"fmt"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
)

const integrationsQuery = `
  query integrations($organizationId: String!, $kind: String!, $platformSupport: String!, $page: Int!, $limit: Int!) {
    integrations(
      organizationId: $organizationId
      kind: $kind
      platformSupport: $platformSupport
      page: $page
      limit: $limit
    ) {
      id
      name
      keyType
      imageIcon
      version
      deprecated
      deprecatedReason
      description
      documentation
      platformSupport
      kind
      public
      properties {
        id
        key
        value
        type
        description
        valuePicker
        valuePickerGroup
        valuePickerOptions
        deprecated
        deprecatedReason
      }
      events {
        id
        event
        description
        deprecated
        deprecatedReason
      }
      data {
        id
        key
        type
        description
        deprecated
        deprecatedReason
      }
      slots {
        id
        slot
        description
        deprecated
        deprecatedReason
      }
    }
  }
`

const syncIntegrationMutation = `
	mutation syncIntegration($input: SyncIntegrationInput!) {
		syncIntegration(input: $input) {
			id
			keyType
			name
			imageIcon
			price
			description
			kind
			documentation
		}
	}
`

func GetIntegrations(fm fileutil.FileManager, regionUrl string, accessToken string, organizationId string, kind string, platformSupport string) ([]IntegrationModel, error) {
	client := graphqlutil.NewClient()

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	variables := map[string]interface{}{
		"organizationId":  organizationId,
		"kind":            kind,
		"platformSupport": platformSupport,
		"page":            0,
		"limit":           1000,
	}

	apiResponse, err := client.Execute(
		regionUrl,
		headers,
		integrationsQuery,
		variables,
	)
	if err != nil {
		return nil, errors.New("failed to fetch projects: " + err.Error())
	}

	var integrationResponse IntegrationsResponse
	err = graphqlutil.Parse(apiResponse, &integrationResponse)
	if err != nil {
		return nil, err
	}
	if len(integrationResponse.Integrations) == 0 {
		return nil, errors.New("no projects found")
	}
	return mapIntegrationsResponseToModel(integrationResponse), nil
}

func SyncIntegration(fm fileutil.FileManager, regionUrl string, accessToken string, organizationId string, jsonInput IntegrationModel) error {
	client := graphqlutil.NewClient()

	input := map[string]interface{}{
		"organizationId":   organizationId,
		"name":             jsonInput.Name,
		"description":      jsonInput.Description,
		"documentation":    "",
		"imageIcon":        "",
		"keyType":          jsonInput.KeyType,
		"kind":             jsonInput.Kind,
		"platformSupport":  jsonInput.PlatformSupport,
		"price":            0,
		"version":          jsonInput.Version,
		"deprecated":       jsonInput.Deprecated,
		"deprecatedReason": jsonInput.DeprecatedReason,
		"public":           false,
	}

	variables := map[string]interface{}{
		"input": input,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	_, err := client.Execute(
		regionUrl,
		headers,
		syncIntegrationMutation,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	return nil
}

func mapIntegrationsResponseToModel(response IntegrationsResponse) []IntegrationModel {
	var integrationModels []IntegrationModel
	for _, integration := range response.Integrations {
		var properties []IntegrationPropertyModel
		for _, prop := range integration.Properties {
			properties = append(properties, IntegrationPropertyModel{
				Key:              prop.Key,
				Value:            prop.Value,
				Type:             prop.Type,
				Description:      prop.Description,
				Deprecated:       prop.Deprecated,
				DeprecatedReason: prop.DeprecatedReason,
			})
		}

		var data []IntegrationDataModel
		for _, dataItem := range integration.Data {
			data = append(data, IntegrationDataModel{
				Key:              dataItem.Key,
				Type:             dataItem.Type,
				Description:      dataItem.Description,
				Deprecated:       dataItem.Deprecated,
				DeprecatedReason: dataItem.DeprecatedReason,
			})
		}

		var events []IntegrationEventModel
		for _, event := range integration.Events {
			events = append(events, IntegrationEventModel{
				Event:            event.Event,
				Description:      event.Description,
				Deprecated:       event.Deprecated,
				DeprecatedReason: event.DeprecatedReason,
			})
		}

		var slots []IntegrationSlotModel
		for _, slot := range integration.Slots {
			slots = append(slots, IntegrationSlotModel{
				Slot:             slot.Slot,
				Description:      slot.Description,
				Deprecated:       slot.Deprecated,
				DeprecatedReason: slot.DeprecatedReason,
			})
		}

		integrationModel := IntegrationModel{
			Name:             integration.Name,
			KeyType:          integration.KeyType,
			Version:          integration.Version,
			Id:               integration.Id,
			PlatformSupport:  integration.PlatformSupport,
			Kind:             integration.Kind,
			Public:           integration.Public,
			Deprecated:       integration.Deprecated,
			DeprecatedReason: integration.DeprecatedReason,
			Properties:       properties,
			Data:             data,
			Events:           events,
			Slots:            slots,
		}
		integrationModels = append(integrationModels, integrationModel)
	}
	return integrationModels
}
