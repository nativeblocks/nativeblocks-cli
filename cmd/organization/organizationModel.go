package organization

type OrganizationModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationItemResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationsResponse struct {
	Organizations []OrganizationItemResponse `json:"organizations"`
}

func mapResponseToModel(response OrganizationsResponse) []OrganizationModel {
	var models []OrganizationModel
	for _, orgItem := range response.Organizations {
		model := OrganizationModel{
			Id:   orgItem.Id,
			Name: orgItem.Name,
		}
		models = append(models, model)
	}
	return models
}
