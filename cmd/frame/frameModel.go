package frame

type FrameModel struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Route          string          `json:"route"`
	RouteArguments []RouteArgument `json:"routeArguments"`
	Type           string          `json:"type"`
	IsStarter      bool            `json:"isStarter"`
	ProjectID      string          `json:"projectId"`
	Checksum       string          `json:"checksum"`
	Variables      []VariableModel `json:"variables"`
	Blocks         []BlockModel    `json:"blocks"`
	Actions        []ActionModel   `json:"actions"`
}

type RouteArgument struct {
	Name string `json:"name"`
}

type VariableModel struct {
	ID      string `json:"id"`
	FrameID string `json:"frameId"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type BlockModel struct {
	ID                 string               `json:"id"`
	FrameID            string               `json:"frameId"`
	KeyType            string               `json:"keyType"`
	Key                string               `json:"key"`
	VisibilityKey      string               `json:"visibilityKey"`
	Position           int                  `json:"position"`
	Slot               *string              `json:"slot,omitempty"`
	IntegrationVersion int                  `json:"integrationVersion"`
	ParentID           *string              `json:"parentId,omitempty"`
	Data               []BlockDataModel     `json:"data"`
	Properties         []BlockPropertyModel `json:"properties"`
	Slots              []BlockSlotModel     `json:"slots"`
}

type BlockPropertyModel struct {
	ID                 string  `json:"id"`
	BlockID            string  `json:"blockId"`
	Key                string  `json:"key"`
	ValueMobile        string  `json:"valueMobile"`
	ValueTablet        string  `json:"valueTablet"`
	ValueDesktop       string  `json:"valueDesktop"`
	Type               string  `json:"type"`
	Description        *string `json:"description,omitempty"`
	ValuePicker        string  `json:"valuePicker"`
	ValuePickerGroup   string  `json:"valuePickerGroup"`
	ValuePickerOptions string  `json:"valuePickerOptions"`
}

type BlockDataModel struct {
	ID          string  `json:"id"`
	BlockID     string  `json:"blockId"`
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	Type        string  `json:"type"`
	Description *string `json:"description,omitempty"`
}

type BlockSlotModel struct {
	ID          string  `json:"id"`
	BlockID     string  `json:"blockId"`
	Slot        string  `json:"slot"`
	Description *string `json:"description,omitempty"`
}

type ActionModel struct {
	ID       string               `json:"id"`
	FrameID  string               `json:"frameId"`
	Key      string               `json:"key"`
	Event    string               `json:"event"`
	Triggers []ActionTriggerModel `json:"triggers"`
}

type ActionTriggerModel struct {
	ID                 string                 `json:"id"`
	ActionID           string                 `json:"actionId"`
	ParentID           *string                `json:"parentId,omitempty"`
	KeyType            string                 `json:"keyType"`
	Then               string                 `json:"then"`
	Name               string                 `json:"name"`
	IntegrationVersion int                    `json:"integrationVersion"`
	Properties         []TriggerPropertyModel `json:"properties"`
	Data               []TriggerDataModel     `json:"data"`
}

type TriggerPropertyModel struct {
	ID                 string  `json:"id"`
	ActionTriggerID    string  `json:"actionTriggerId"`
	Key                string  `json:"key"`
	Value              string  `json:"value"`
	Type               string  `json:"type"`
	Description        *string `json:"description,omitempty"`
	ValuePicker        *string `json:"valuePicker,omitempty"`
	ValuePickerGroup   *string `json:"valuePickerGroup,omitempty"`
	ValuePickerOptions *string `json:"valuePickerOptions,omitempty"`
}

type TriggerDataModel struct {
	ID              string  `json:"id"`
	ActionTriggerID string  `json:"actionTriggerId"`
	Key             string  `json:"key"`
	Value           string  `json:"value"`
	Type            string  `json:"type"`
	Description     *string `json:"description,omitempty"`
}
