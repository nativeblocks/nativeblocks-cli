package projectModule

import (
	"errors"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
)

const ProjectFileName = "project"

const projectsQuery = `
  query projects($organizationId: String!) {
    projects(organizationId: $organizationId) {
      id
      name
      platform
      apiKeys {
        name
        apiKey
      }
    }
  }
`

const installedIntegrationsQuery = `
	query integrationsInstalled($organizationId: String!, $projectId: String!, $kind: String!) {
		integrationsInstalled(organizationId: $organizationId, projectId: $projectId, kind: $kind) {
			integrationKeyType
			integrationVersion
			integrationId
			integrationPlatformSupport
			integrationKind
			integrationProperties {
				key
				value
				type
			}
			integrationData {
				key
				type
			}
			integrationEvents {
				event
			}
			integrationSlots {
				slot
			}
		}
	}
`

func GetProjects(regionUrl string, accessToken string, organizationId string) ([]ProjectModel, error) {
	client := graphqlutil.NewClient()

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	variables := map[string]interface{}{
		"organizationId": organizationId,
	}

	apiResponse, err := client.Execute(
		regionUrl,
		headers,
		projectsQuery,
		variables,
	)
	if err != nil {
		return nil, errors.New("failed to fetch projects: " + err.Error())
	}

	var projResp ProjectsResponse
	err = graphqlutil.Parse(apiResponse, &projResp)
	if err != nil {
		return nil, err
	}
	if len(projResp.Projects) == 0 {
		return nil, errors.New("no projects found")
	}
	return mapProjectsResponseToModel(projResp), nil
}

func SelectProject(fm fileutil.FileManager, projectModel *ProjectModel) error {
	if err := fm.SaveToFile(ProjectFileName, projectModel); err != nil {
		return errors.New("failed to save project config: " + err.Error())
	}
	return nil
}

func GetProject(fm fileutil.FileManager) (*ProjectModel, error) {
	var model ProjectModel
	if err := fm.LoadFromFile(ProjectFileName, &model); err != nil {
		return nil, errors.New("project not set. Please select a project first using 'nativeblocks project set'")
	}
	return &model, nil
}

func GetInstalledIntegration(regionUrl string, accessToken string, organizationId string, projectId string, kind string) ([]IntegrationProjectModel, error) {

	client := graphqlutil.NewClient()

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	variables := map[string]interface{}{
		"organizationId": organizationId,
		"projectId":      projectId,
		"kind":           kind,
	}

	apiResponse, err := client.Execute(
		regionUrl,
		headers,
		installedIntegrationsQuery,
		variables,
	)
	if err != nil {
		return nil, errors.New("failed to fetch installed integrations: " + err.Error())
	}

	var installedIntegrationResponse InstalledIntegrationResponse
	err = graphqlutil.Parse(apiResponse, &installedIntegrationResponse)
	if err != nil {
		return nil, err
	}

	return mapIntegrationsResponseToModel(installedIntegrationResponse), nil
}
