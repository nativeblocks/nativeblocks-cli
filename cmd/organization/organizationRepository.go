package organization

import (
	"errors"

	"github.com/nativeblocks/cli/library/fileutil"
	"github.com/nativeblocks/cli/library/graphqlutil"
)

const orgFileName = "organization"

const organizationsQuery = `
  query organizations {
    organizations {
      id
      name
    }
  }
`

func GetOrganizations(fm fileutil.FileManager, regionUrl, accessToken string) ([]OrganizationModel, error) {
	client := graphqlutil.NewClient()

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	apiResponse, err := client.Execute(
		regionUrl,
		headers,
		organizationsQuery,
		nil,
	)
	if err != nil {
		return nil, errors.New("failed to fetch organizations: " + err.Error())
	}

	var orgResp OrganizationsResponse
	err = graphqlutil.Parse(apiResponse, &orgResp)
	if err != nil {
		return nil, err
	}

	if len(orgResp.Organizations) == 0 {
		return nil, errors.New("no organizations found")
	}
	orgs := mapResponseToModel(orgResp)
	return orgs, nil
}

func SelectOrganization(fm *fileutil.FileManager, orgModel *OrganizationModel) error {
	if err := fm.SaveToFile(orgFileName, orgModel); err != nil {
		return errors.New("failed to save organization config: " + err.Error())
	}
	return nil
}

func GetOrganization(fm fileutil.FileManager) (*OrganizationModel, error) {
	var model OrganizationModel
	if err := fm.LoadFromFile(orgFileName, &model); err != nil {
		return nil, errors.New("organization not set. Please select an organization first using 'nativeblocks organization set'")
	}
	return &model, nil
}
