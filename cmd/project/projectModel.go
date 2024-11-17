package project

type ProjectModel struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Platform string        `json:"platform"`
	APIKeys  []APIKeyModel `json:"apiKeys"`
}

type APIKeyModel struct {
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

type ProjectsResponse struct {
	Projects []ProjectModel `json:"projects"`
}

type IntegrationPropertyModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type IntegrationDataModel struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type IntegrationEventModel struct {
	Event string `json:"event"`
}

type IntegrationSlotModel struct {
	Slot string `json:"slot"`
}

type IntegrationProjectModel struct {
	IntegrationKeyType         string                     `json:"integrationKeyType"`
	IntegrationVersion         int8                       `json:"integrationVersion"`
	IntegrationID              string                     `json:"integrationId"`
	IntegrationPlatformSupport string                     `json:"integrationPlatformSupport"`
	IntegrationKind            string                     `json:"integrationKind"`
	IntegrationProperties      []IntegrationPropertyModel `json:"integrationProperties"`
	IntegrationData            []IntegrationDataModel     `json:"integrationData"`
	IntegrationEvents          []IntegrationEventModel    `json:"integrationEvents"`
	IntegrationSlots           []IntegrationSlotModel     `json:"integrationSlots"`
}

type InstalledIntegrationResponse struct {
	IntegrationsInstalled []IntegrationProjectModel `json:"integrationsInstalled"`
}

func mapProjectsResponseToModel(response ProjectsResponse) []ProjectModel {
	var projectModels []ProjectModel
	for _, project := range response.Projects {
		var apiKeys []APIKeyModel
		for _, apiKey := range project.APIKeys {
			apiKeys = append(apiKeys, APIKeyModel{
				Name:   apiKey.Name,
				APIKey: apiKey.APIKey,
			})
		}

		projectModel := ProjectModel{
			Id:       project.Id,
			Name:     project.Name,
			Platform: project.Platform,
			APIKeys:  apiKeys,
		}
		projectModels = append(projectModels, projectModel)
	}
	return projectModels
}

func mapIntegrationsResponseToModel(response InstalledIntegrationResponse) []IntegrationProjectModel {
	var integrationModels []IntegrationProjectModel
	for _, integration := range response.IntegrationsInstalled {
		var properties []IntegrationPropertyModel
		for _, prop := range integration.IntegrationProperties {
			properties = append(properties, IntegrationPropertyModel{
				Key:   prop.Key,
				Value: prop.Value,
				Type:  prop.Type,
			})
		}

		var data []IntegrationDataModel
		for _, dataItem := range integration.IntegrationData {
			data = append(data, IntegrationDataModel{
				Key:  dataItem.Key,
				Type: dataItem.Type,
			})
		}

		var events []IntegrationEventModel
		for _, event := range integration.IntegrationEvents {
			events = append(events, IntegrationEventModel{
				Event: event.Event,
			})
		}

		var slots []IntegrationSlotModel
		for _, slot := range integration.IntegrationSlots {
			slots = append(slots, IntegrationSlotModel{
				Slot: slot.Slot,
			})
		}

		integrationModel := IntegrationProjectModel{
			IntegrationKeyType:         integration.IntegrationKeyType,
			IntegrationVersion:         integration.IntegrationVersion,
			IntegrationID:              integration.IntegrationID,
			IntegrationPlatformSupport: integration.IntegrationPlatformSupport,
			IntegrationKind:            integration.IntegrationKind,
			IntegrationProperties:      properties,
			IntegrationData:            data,
			IntegrationEvents:          events,
			IntegrationSlots:           slots,
		}
		integrationModels = append(integrationModels, integrationModel)
	}
	return integrationModels
}
