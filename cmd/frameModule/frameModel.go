package frameModule

type FrameWrapper struct {
	Frame FrameModel `json:"frame"`
}
type FrameProductionWrapper struct {
	FrameProduction FrameModel `json:"frameProduction"`
}
type FrameProductionDataWrapper struct {
	Data FrameProductionWrapper `json:"data"`
}

type FrameModel struct {
	Id             string          `json:"id"`
	Name           string          `json:"name"`
	Route          string          `json:"route"`
	RouteArguments []RouteArgument `json:"routeArguments"`
	Type           string          `json:"type"`
	IsStarter      bool            `json:"isStarter"`
	ProjectId      string          `json:"projectId"`
	Checksum       string          `json:"checksum"`
	Variables      []VariableModel `json:"variables"`
	Blocks         []BlockModel    `json:"blocks"`
	Actions        []ActionModel   `json:"actions"`
}

type RouteArgument struct {
	Name string `json:"name"`
}

type VariableModel struct {
	Id      string `json:"id"`
	FrameId string `json:"frameId"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type BlockModel struct {
	Id                          string               `json:"id"`
	FrameId                     string               `json:"frameId"`
	KeyType                     string               `json:"keyType"`
	Key                         string               `json:"key"`
	VisibilityKey               string               `json:"visibilityKey"`
	Position                    int                  `json:"position"`
	Slot                        string               `json:"slot"`
	IntegrationVersion          int                  `json:"integrationVersion"`
	ParentId                    string               `json:"parentId"`
	Data                        []BlockDataModel     `json:"data"`
	Properties                  []BlockPropertyModel `json:"properties"`
	Slots                       []BlockSlotModel     `json:"slots"`
	IntegrationDeprecated       bool                 `json:"integrationDeprecated"`
	IntegrationDeprecatedReason string               `json:"integrationDeprecatedReason"`
}

type BlockPropertyModel struct {
	Id                 string `json:"id"`
	BlockId            string `json:"blockId"`
	Key                string `json:"key"`
	ValueMobile        string `json:"valueMobile"`
	ValueTablet        string `json:"valueTablet"`
	ValueDesktop       string `json:"valueDesktop"`
	Type               string `json:"type"`
	Description        string `json:"description"`
	ValuePicker        string `json:"valuePicker"`
	ValuePickerGroup   string `json:"valuePickerGroup"`
	ValuePickerOptions string `json:"valuePickerOptions"`
	Deprecated         bool   `json:"deprecated"`
	DeprecatedReason   string `json:"deprecatedReason"`
}

type BlockDataModel struct {
	Id               string `json:"id"`
	BlockId          string `json:"blockId"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

type BlockSlotModel struct {
	Id               string `json:"id"`
	BlockId          string `json:"blockId"`
	Slot             string `json:"slot"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

type ActionModel struct {
	Id       string               `json:"id"`
	FrameId  string               `json:"frameId"`
	Key      string               `json:"key"`
	Event    string               `json:"event"`
	Triggers []ActionTriggerModel `json:"triggers"`
}

type ActionTriggerModel struct {
	Id                          string                 `json:"id"`
	ActionId                    string                 `json:"actionId"`
	ParentId                    string                 `json:"parentId"`
	KeyType                     string                 `json:"keyType"`
	Then                        string                 `json:"then"`
	Name                        string                 `json:"name"`
	IntegrationVersion          int                    `json:"integrationVersion"`
	Properties                  []TriggerPropertyModel `json:"properties"`
	Data                        []TriggerDataModel     `json:"data"`
	IntegrationDeprecated       bool                   `json:"integrationDeprecated"`
	IntegrationDeprecatedReason string                 `json:"integrationDeprecatedReason"`
}

type TriggerPropertyModel struct {
	Id                 string `json:"id"`
	ActionTriggerId    string `json:"actionTriggerId"`
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

type TriggerDataModel struct {
	Id               string `json:"id"`
	ActionTriggerId  string `json:"actionTriggerId"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}
