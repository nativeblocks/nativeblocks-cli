package integrationModule

type IntegrationModel struct {
	Id               string                     `json:"id"`
	Name             string                     `json:"name"`
	Description      string                     `json:"description"`
	KeyType          string                     `json:"keyType"`
	PlatformSupport  string                     `json:"platformSupport"`
	Kind             string                     `json:"kind"`
	ImageIcon        string                     `json:"imageIcon"`
	Version          int                        `json:"version"`
	Documentation    string                     `json:"documentation"`
	Public           bool                       `json:"public"`
	Manageable       bool                       `json:"manageable"`
	Deprecated       bool                       `json:"deprecated"`
	DeprecatedReason string                     `json:"deprecatedReason"`
	Properties       []IntegrationPropertyModel `json:"properties,omitempty"`
	Events           []IntegrationEventModel    `json:"events,omitempty"`
	Data             []IntegrationDataModel     `json:"data,omitempty"`
	Slots            []IntegrationSlotModel     `json:"slots,omitempty"`
}

type IntegrationResponse struct {
	Integration IntegrationModel `json:"integration"`
}
type IntegrationsResponse struct {
	Integrations []IntegrationModel `json:"integrations"`
}

type SyncIntegrationResponse struct {
	Integration IntegrationModel `json:"syncIntegration"`
}

type IntegrationDataModel struct {
	Key              string `json:"key"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

type IntegrationPropertyModel struct {
	Key                string `json:"key"`
	Value              string `json:"value"`
	Type               string `json:"type"`
	Description        string `json:"description"`
	ValuePicker        string `json:"valuePicker"`
	ValuePickerGroup   string `json:"valuePickerGroup"`
	ValuePickerOptions string `json:"valuePickerOptions"`
	Deprecated         bool   `json:"deprecated"`
	DeprecatedReason   string `json:"deprecatedReason"`
}

type IntegrationEventModel struct {
	Event            string `json:"event"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

type IntegrationSlotModel struct {
	Slot             string `json:"slot"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}
