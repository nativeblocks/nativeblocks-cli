package integration

import (
	"errors"
	"fmt"
	"sync"

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

const syncIntegrationPropertyMutation = `
	mutation syncIntegrationProperties($input: SyncIntegrationPropertiesInput!) {
		syncIntegrationProperties(input: $input) {
			key
			type
			value
			description
			valuePicker
			valuePickerGroup
			valuePickerOptions
		}
	}
`

const syncIntegrationDataMutation = `
mutation syncIntegrationData($input: SyncIntegrationDataInput!) {
    syncIntegrationData(input: $input) {
        key
        type
    }
}`

const syncIntegrationEventsMutation = `
mutation syncIntegrationEvents($input: SyncIntegrationEventsInput!) {
    syncIntegrationEvents(input: $input) {
        event
    }
}`

const syncIntegrationSlotsMutation = `
mutation syncIntegrationSlots($input: SyncIntegrationSlotsInput!) {
    syncIntegrationSlots(input: $input) {
        slot
    }
}`

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

func SyncIntegration(regionUrl string, accessToken string, organizationId string, directory string) error {
	integrationFileName := "integration.json"
	propertiesFileName := "properties.json"
	dataFileName := "data.json"
	eventsFileName := "events.json"
	slotsFileName := "slots.json"

	inputFm, err := fileutil.NewFileManager(&directory)
	if err != nil {
		return err
	}

	fileExists := inputFm.FileExists(integrationFileName)
	if !fileExists {
		return fmt.Errorf("could not find the file under: %v", directory)
	}

	var jsonInput IntegrationModel
	err = inputFm.LoadFromFile(integrationFileName, &jsonInput)
	if err != nil {
		return err
	}

	if jsonInput.KeyType == "" {
		return fmt.Errorf("could not find integration keyType")
	}

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

	// Get initial integration sync response
	apiResponse, err := client.Execute(
		regionUrl,
		headers,
		syncIntegrationMutation,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	var syncIntegrationsResponse SyncIntegrationsResponse
	err = graphqlutil.Parse(apiResponse, &syncIntegrationsResponse)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	if inputFm.FileExists(propertiesFileName) {
		var jsonInput []IntegrationPropertyModel
		err = inputFm.LoadFromFile(propertiesFileName, &jsonInput)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := SyncIntegrationProperties(*inputFm, regionUrl, accessToken, organizationId, syncIntegrationsResponse.Integration.Id, jsonInput); err != nil {
				fmt.Printf("Error syncing properties: %v\n", err)
			}
		}()
	}

	if inputFm.FileExists(eventsFileName) {
		var jsonInput []IntegrationEventModel
		err = inputFm.LoadFromFile(eventsFileName, &jsonInput)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := SyncIntegrationEvents(*inputFm, regionUrl, accessToken, organizationId, syncIntegrationsResponse.Integration.Id, jsonInput); err != nil {
				fmt.Printf("Error syncing events: %v\n", err)
			}
		}()
	}

	if inputFm.FileExists(dataFileName) {
		var jsonInput []IntegrationDataModel
		err = inputFm.LoadFromFile(dataFileName, &jsonInput)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := SyncIntegrationData(*inputFm, regionUrl, accessToken, organizationId, syncIntegrationsResponse.Integration.Id, jsonInput); err != nil {
				fmt.Printf("Error syncing data: %v\n", err)
			}
		}()
	}

	if inputFm.FileExists(slotsFileName) {
		var jsonInput []IntegrationSlotModel
		err = inputFm.LoadFromFile(slotsFileName, &jsonInput)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := SyncIntegrationSlots(*inputFm, regionUrl, accessToken, organizationId, syncIntegrationsResponse.Integration.Id, jsonInput); err != nil {
				fmt.Printf("Error syncing slots: %v\n", err)
			}
		}()
	}

	wg.Wait()
	return nil
}

func SyncIntegrationProperties(fm fileutil.FileManager, regionUrl string, accessToken string, organizationId string, integrationId string, jsonInput []IntegrationPropertyModel) error {
	if len(jsonInput) == 0 {
		return nil
	}

	client := graphqlutil.NewClient()

	var properties []map[string]interface{}
	for _, prop := range jsonInput {
		properties = append(properties, map[string]interface{}{
			"key":                prop.Key,
			"value":              prop.Value,
			"type":               prop.Type,
			"description":        prop.Description,
			"valuePicker":        prop.ValuePicker,
			"valuePickerGroup":   prop.ValuePickerGroup,
			"valuePickerOptions": prop.ValuePickerOptions,
			"deprecated":         prop.Deprecated,
			"deprecatedReason":   prop.DeprecatedReason,
		})
	}

	input := map[string]interface{}{
		"organizationId": organizationId,
		"integrationId":  integrationId,
		"properties":     properties,
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
		syncIntegrationPropertyMutation,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	return nil
}

func SyncIntegrationEvents(fm fileutil.FileManager, regionUrl string, accessToken string, organizationId string, integrationId string, jsonInput []IntegrationEventModel) error {
	if len(jsonInput) == 0 {
		return nil
	}

	client := graphqlutil.NewClient()

	var events []map[string]interface{}
	for _, prop := range jsonInput {
		events = append(events, map[string]interface{}{
			"event":            prop.Event,
			"description":      prop.Description,
			"deprecated":       prop.Deprecated,
			"deprecatedReason": prop.DeprecatedReason,
		})
	}

	input := map[string]interface{}{
		"organizationId": organizationId,
		"integrationId":  integrationId,
		"events":         events,
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
		syncIntegrationEventsMutation,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	return nil
}

func SyncIntegrationData(fm fileutil.FileManager, regionUrl string, accessToken string, organizationId string, integrationId string, jsonInput []IntegrationDataModel) error {
	if len(jsonInput) == 0 {
		return nil
	}

	client := graphqlutil.NewClient()

	var data []map[string]interface{}
	for _, prop := range jsonInput {
		data = append(data, map[string]interface{}{
			"key":              prop.Key,
			"type":             prop.Type,
			"description":      prop.Description,
			"deprecated":       prop.Deprecated,
			"deprecatedReason": prop.DeprecatedReason,
		})
	}

	input := map[string]interface{}{
		"organizationId": organizationId,
		"integrationId":  integrationId,
		"data":           data,
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
		syncIntegrationDataMutation,
		variables,
	)
	if err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}

	return nil
}

func SyncIntegrationSlots(fm fileutil.FileManager, regionUrl string, accessToken string, organizationId string, integrationId string, jsonInput []IntegrationSlotModel) error {
	if len(jsonInput) == 0 {
		return nil
	}

	client := graphqlutil.NewClient()

	var slots []map[string]interface{}
	for _, prop := range jsonInput {
		slots = append(slots, map[string]interface{}{
			"slot":             prop.Slot,
			"description":      prop.Description,
			"deprecated":       prop.Deprecated,
			"deprecatedReason": prop.DeprecatedReason,
		})
	}

	input := map[string]interface{}{
		"organizationId": organizationId,
		"integrationId":  integrationId,
		"slots":          slots,
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
		syncIntegrationSlotsMutation,
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
