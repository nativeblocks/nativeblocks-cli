package frame

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/xeipuuv/gojsonschema"
)

func processActions(frameID, key string, inputActions []ActionDSLModel, variables []VariableModel) ([]ActionModel, error) {
	var actions []ActionModel

	for _, inputAction := range inputActions {
		actionID := generateId()
		subTriggers, err := processTriggers(actionID, inputAction.Triggers, nil, variables)
		if err != nil {
			return nil, err
		}

		newAction := ActionModel{
			ID:       actionID,
			FrameID:  frameID,
			Key:      key,
			Event:    inputAction.Event,
			Triggers: subTriggers,
		}
		actions = append(actions, newAction)
	}

	return actions, nil
}

func processTriggers(actionID string, triggers []ActionTriggerDSLModel, parentID *string, variables []VariableModel) ([]ActionTriggerModel, error) {
	var flatTriggers []ActionTriggerModel

	for _, trigger := range triggers {
		newTrigger := ActionTriggerModel{
			ID:                 generateId(),
			ActionID:           actionID,
			ParentID:           parentID,
			KeyType:            trigger.KeyType,
			Then:               trigger.Then,
			Name:               trigger.Name,
			IntegrationVersion: trigger.IntegrationVersion,
			Properties:         []TriggerPropertyModel{},
			Data:               []TriggerDataModel{},
		}

		for _, property := range trigger.Properties {
			newProperty := TriggerPropertyModel{
				ID:                 generateId(),
				ActionTriggerID:    newTrigger.ID,
				Key:                property.Key,
				Type:               property.Type,
				Value:              property.Value,
				Description:        property.Description,
				ValuePicker:        property.ValuePicker,
				ValuePickerGroup:   property.ValuePickerGroup,
				ValuePickerOptions: property.ValuePickerOptions,
			}
			newTrigger.Properties = append(newTrigger.Properties, newProperty)
		}

		for _, dataItem := range trigger.Data {
			newData := TriggerDataModel{
				ID:              generateId(),
				ActionTriggerID: newTrigger.ID,
				Key:             dataItem.Key,
				Value:           dataItem.Value,
				Type:            dataItem.Type,
				Description:     dataItem.Description,
			}
			newTrigger.Data = append(newTrigger.Data, newData)
		}

		for i, dataEntry := range newTrigger.Data {
			for _, variable := range variables {
				if variable.Key == dataEntry.Value {
					newTrigger.Data[i].Value = variable.Key
					newTrigger.Data[i].Type = variable.Type
					break
				} else {
					return nil, errors.New("Variable with key '" + dataEntry.Value + "' not found")
					break
				}
			}
		}

		flatTriggers = append(flatTriggers, newTrigger)

		if len(trigger.Triggers) > 0 {
			subTriggers, err := processTriggers(actionID, trigger.Triggers, &newTrigger.ID, variables)
			if err != nil {
				return nil, err
			}
			flatTriggers = append(flatTriggers, subTriggers...)
		}
	}

	if flatTriggers == nil {
		flatTriggers = []ActionTriggerModel{}
	}

	return flatTriggers, nil
}

func processBlocks(frameID string, blocks []BlockDSLModel, parentID *string, parentSlots []BlockSlotModel, variables []VariableModel, onNewAction func([]ActionModel, error)) ([]BlockModel, error) {
	var flatBlocks []BlockModel

	for index, block := range blocks {
		newBlock := BlockModel{
			ID:                 generateId(),
			FrameID:            frameID,
			KeyType:            block.KeyType,
			Key:                block.Key,
			VisibilityKey:      block.VisibilityKey,
			Position:           index,
			Slot:               block.Slot,
			IntegrationVersion: block.IntegrationVersion,
			ParentID:           parentID,
			Data:               []BlockDataModel{},
			Properties:         []BlockPropertyModel{},
			Slots:              []BlockSlotModel{},
		}

		if newBlock.Slot == nil {
			contentSlot := "content"
			newBlock.Slot = &contentSlot
		}

		onNewAction(processActions(frameID, block.Key, block.Actions, variables))

		for _, property := range block.Properties {
			newProperty := BlockPropertyModel{
				ID:                 generateId(),
				BlockID:            newBlock.ID,
				Key:                property.Key,
				Type:               property.Type,
				ValueMobile:        property.ValueMobile,
				ValueTablet:        property.ValueTablet,
				ValueDesktop:       property.ValueDesktop,
				Description:        property.Description,
				ValuePicker:        property.ValuePicker,
				ValuePickerGroup:   property.ValuePickerGroup,
				ValuePickerOptions: property.ValuePickerOptions,
			}
			newBlock.Properties = append(newBlock.Properties, newProperty)
		}

		for _, dataItem := range block.Data {
			newData := BlockDataModel{
				ID:          generateId(),
				BlockID:     newBlock.ID,
				Key:         dataItem.Key,
				Value:       dataItem.Value,
				Type:        dataItem.Type,
				Description: dataItem.Description,
			}
			newBlock.Data = append(newBlock.Data, newData)
		}

		for _, slotItem := range block.Slots {
			newSlot := BlockSlotModel{
				ID:          generateId(),
				BlockID:     newBlock.ID,
				Slot:        slotItem.Slot,
				Description: slotItem.Description,
			}
			newBlock.Slots = append(newBlock.Slots, newSlot)
		}

		for _, dataEntry := range newBlock.Data {
			for _, variable := range variables {
				if variable.Key == dataEntry.Value {
					dataEntry.Value = variable.Key
					dataEntry.Type = variable.Type
					break
				} else {
					return nil, errors.New("Variable with key '" + dataEntry.Value + "' not found")
					break
				}
			}
		}

		flatBlocks = append(flatBlocks, newBlock)

		if len(block.Blocks) > 0 {
			subBlocks, err := processBlocks(frameID, block.Blocks, &newBlock.ID, newBlock.Slots, variables, onNewAction)
			if err != nil {
				return nil, err
			}
			flatBlocks = append(flatBlocks, subBlocks...)
		}
	}

	return flatBlocks, nil
}

func getWordsBetweenCurly(text string) []string {
	re := regexp.MustCompile(`\{(.*?)\}`)
	matches := re.FindAllStringSubmatch(text, -1)

	var result []string
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1])
		}
	}
	return result
}

func convertRouteArguments(route string) []RouteArgument {
	args := getWordsBetweenCurly(route)
	routeArguments := make([]RouteArgument, len(args))

	for i, arg := range args {
		routeArguments[i] = RouteArgument{Name: arg}
	}
	return routeArguments
}

func generateFrame(frameDSL FrameDSLModel) (FrameModel, error) {
	if frameDSL.Schema == "" {
		return FrameModel{}, errors.New("Please provide $schema for the json file")
	}

	schemaLoader := gojsonschema.NewReferenceLoader(frameDSL.Schema)
	documentLoader := gojsonschema.NewGoLoader(frameDSL)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return FrameModel{}, err
	}

	if !result.Valid() {
		for _, errz := range result.Errors() {
			fmt.Printf("- %s\n", errz)
		}
		return FrameModel{}, nil
	}

	frameID := generateId()

	var variables []VariableModel
	for _, variable := range frameDSL.Variables {
		variables = append(variables, VariableModel{
			ID:      generateId(),
			FrameID: frameID,
			Key:     variable.Key,
			Value:   variable.Value,
			Type:    variable.Type,
		})
	}

	var actions []ActionModel
	blocks, err := processBlocks(frameID, frameDSL.Blocks, nil, []BlockSlotModel{}, variables, func(blockActions []ActionModel, err error) {
		actions = append(actions, blockActions...)
	})

	frame := FrameModel{
		ID:             frameID,
		Name:           frameDSL.Name,
		Route:          frameDSL.Route,
		RouteArguments: convertRouteArguments(frameDSL.Route),
		Type:           frameDSL.Type,
		IsStarter:      frameDSL.IsStarter,
		Variables:      variables,
		Blocks:         blocks,
		Actions:        actions,
	}

	if frame.Actions == nil {
		frame.Actions = []ActionModel{}
	}
	if frame.Blocks == nil {
		frame.Blocks = []BlockModel{}
	}
	return frame, err
}
