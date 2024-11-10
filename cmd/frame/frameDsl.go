package frame

import (
	"github.com/google/uuid"
)

func generateId() string {
	return uuid.New().String()
}

type FrameDSLModel struct {
	Schema    string             `json:"$schema"`
	Name      string             `json:"name"`
	Route     string             `json:"route"`
	Type      string             `json:"type"`
	IsStarter bool               `json:"isStarter"`
	Variables []VariableDSLModel `json:"variables"`
	Blocks    []BlockDSLModel    `json:"blocks"`
}

type VariableDSLModel struct {
	FrameId string `json:"frameId"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type BlockDSLModel struct {
	KeyType            string                  `json:"keyType"`
	Key                string                  `json:"key"`
	VisibilityKey      string                  `json:"visibilityKey"`
	Slot               string                  `json:"slot,omitempty"`
	IntegrationVersion int                     `json:"integrationVersion"`
	Blocks             []BlockDSLModel         `json:"blocks"`
	Data               []BlockDataDSLModel     `json:"data"`
	Properties         []BlockPropertyDSLModel `json:"properties"`
	Slots              []BlockSlotDSLModel     `json:"slots"`
	Actions            []ActionDSLModel        `json:"actions"`
}

type BlockPropertyDSLModel struct {
	Key                string `json:"key"`
	ValueMobile        string `json:"valueMobile"`
	ValueTablet        string `json:"valueTablet"`
	ValueDesktop       string `json:"valueDesktop"`
	Type               string `json:"type"`
	Description        string `json:"description"`
	ValuePicker        string `json:"valuePicker"`
	ValuePickerGroup   string `json:"valuePickerGroup"`
	ValuePickerOptions string `json:"valuePickerOptions"`
}

type BlockDataDSLModel struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type BlockEventDSLModel struct {
	EventName string           `json:"eventName"`
	Actions   []ActionDSLModel `json:"actions"`
}

type BlockSlotDSLModel struct {
	Slot        string `json:"slot"`
	Description string `json:"description,omitempty"`
}

type ActionDSLModel struct {
	Key      string                  `json:"key"`
	Event    string                  `json:"event"`
	Triggers []ActionTriggerDSLModel `json:"triggers"`
}

type ActionTriggerDSLModel struct {
	KeyType            string                    `json:"keyType"`
	Then               string                    `json:"then"`
	Name               string                    `json:"name"`
	IntegrationVersion int                       `json:"integrationVersion"`
	Properties         []TriggerPropertyDSLModel `json:"properties"`
	Data               []TriggerDataDSLModel     `json:"data"`
	Triggers           []ActionTriggerDSLModel   `json:"triggers"`
}

type TriggerPropertyDSLModel struct {
	Key                string `json:"key"`
	Value              string `json:"value"`
	Type               string `json:"type"`
	Description        string `json:"description,omitempty"`
	ValuePicker        string `json:"valuePicker,omitempty"`
	ValuePickerGroup   string `json:"valuePickerGroup,omitempty"`
	ValuePickerOptions string `json:"valuePickerOptions,omitempty"`
}

type TriggerDataDSLModel struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}
