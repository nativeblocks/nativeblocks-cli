package frameModule

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func processActions(frameId, key string, inputActions []ActionDSLModel, variables []VariableModel) ([]ActionModel, error) {
	var actions []ActionModel

	for _, inputAction := range inputActions {
		actionId := generateId()
		subTriggers, err := processTriggers(actionId, inputAction.Triggers, "", variables)
		if err != nil {
			return nil, err
		}

		newAction := ActionModel{
			Id:       actionId,
			FrameId:  frameId,
			Key:      key,
			Event:    inputAction.Event,
			Triggers: subTriggers,
		}
		actions = append(actions, newAction)
	}

	return actions, nil
}

func processTriggers(actionId string, triggers []ActionTriggerDSLModel, parentId string, variables []VariableModel) ([]ActionTriggerModel, error) {
	var flatTriggers []ActionTriggerModel

	for _, trigger := range triggers {
		newTrigger := ActionTriggerModel{
			Id:                 generateId(),
			ActionId:           actionId,
			ParentId:           parentId,
			KeyType:            trigger.KeyType,
			Then:               trigger.Then,
			Name:               trigger.Name,
			IntegrationVersion: trigger.IntegrationVersion,
			Properties:         []TriggerPropertyModel{},
			Data:               []TriggerDataModel{},
		}

		if newTrigger.Then == "END" && len(trigger.Triggers) > 0 {
			return nil, errors.New("The " + newTrigger.Name + " can not have a subTrigger because it defines with \"END\" then ")
		}

		for _, property := range trigger.Properties {
			newProperty := TriggerPropertyModel{
				Id:                 generateId(),
				ActionTriggerId:    newTrigger.Id,
				Key:                property.Key,
				Type:               property.Type,
				Value:              property.Value,
				Description:        "",
				ValuePicker:        "",
				ValuePickerGroup:   "",
				ValuePickerOptions: "",
			}

			newTrigger.Properties = append(newTrigger.Properties, newProperty)
		}

		for _, dataItem := range trigger.Data {
			newData := TriggerDataModel{
				Id:              generateId(),
				ActionTriggerId: newTrigger.Id,
				Key:             dataItem.Key,
				Value:           dataItem.Value,
				Type:            dataItem.Type,
				Description:     "",
			}

			newTrigger.Data = append(newTrigger.Data, newData)
		}

		err := findTriggerVariable(variables, newTrigger.Data, newTrigger.Name)
		if err != nil {
			return nil, err
		}

		flatTriggers = append(flatTriggers, newTrigger)

		if len(trigger.Triggers) > 0 {
			subTriggers, err := processTriggers(actionId, trigger.Triggers, newTrigger.Id, variables)
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

func processBlocks(frameId string, blocks []BlockDSLModel, parentId string, parentSlots []BlockSlotModel, variables []VariableModel, onNewAction func([]ActionModel)) ([]BlockModel, error) {
	var flatBlocks []BlockModel

	for index, block := range blocks {
		newBlock := BlockModel{
			Id:                 generateId(),
			FrameId:            frameId,
			KeyType:            block.KeyType,
			Key:                block.Key,
			VisibilityKey:      block.VisibilityKey,
			Position:           index,
			Slot:               block.Slot,
			IntegrationVersion: block.IntegrationVersion,
			ParentId:           parentId,
			Data:               []BlockDataModel{},
			Properties:         []BlockPropertyModel{},
			Slots:              []BlockSlotModel{},
		}

		if newBlock.Slot == "null" {
			emptySlot := ""
			newBlock.Slot = emptySlot
		}

		if len(parentSlots) > 0 {
			contain := containsSlot(parentSlots, newBlock.Slot)
			if !contain {
				return nil, errors.New("The " + newBlock.Key + " used in a wrong slot")
			}
		}

		processedActions, err := processActions(frameId, block.Key, block.Actions, variables)
		if err != nil {
			return nil, err
		}
		onNewAction(processedActions)

		for _, property := range block.Properties {
			newProperty := BlockPropertyModel{
				Id:                 generateId(),
				BlockId:            newBlock.Id,
				Key:                property.Key,
				Type:               property.Type,
				ValueMobile:        property.ValueMobile,
				ValueTablet:        property.ValueTablet,
				ValueDesktop:       property.ValueDesktop,
				Description:        "",
				ValuePicker:        "",
				ValuePickerGroup:   "",
				ValuePickerOptions: "",
			}

			newBlock.Properties = append(newBlock.Properties, newProperty)
		}

		for _, dataItem := range block.Data {
			newData := BlockDataModel{
				Id:          generateId(),
				BlockId:     newBlock.Id,
				Key:         dataItem.Key,
				Value:       dataItem.Value,
				Type:        dataItem.Type,
				Description: "",
			}

			newBlock.Data = append(newBlock.Data, newData)
		}

		for _, slotItem := range block.Slots {
			newSlot := BlockSlotModel{
				Id:          generateId(),
				BlockId:     newBlock.Id,
				Slot:        slotItem.Slot,
				Description: "",
			}

			newBlock.Slots = append(newBlock.Slots, newSlot)
		}

		err = findBlockVariable(variables, newBlock.Data, newBlock.Key)
		if err != nil {
			return nil, err
		}

		flatBlocks = append(flatBlocks, newBlock)

		if len(block.Blocks) > 0 {
			subBlocks, err := processBlocks(frameId, block.Blocks, newBlock.Id, newBlock.Slots, variables, onNewAction)
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

func generateFrame(frameDSL FrameDSLModel) (FrameProductionDataWrapper, error) {
	if frameDSL.Schema == "" {
		return FrameProductionDataWrapper{}, errors.New("please provide $schema for the json file")
	}

	schemaLoader := gojsonschema.NewReferenceLoader(frameDSL.Schema)
	documentLoader := gojsonschema.NewGoLoader(frameDSL)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return FrameProductionDataWrapper{}, err
	}

	if !result.Valid() {
		for _, errz := range result.Errors() {
			fmt.Printf("- %s\n", errz)
		}
		return FrameProductionDataWrapper{}, nil
	}

	frameId := generateId()

	var variables []VariableModel
	for _, variable := range frameDSL.Variables {
		variables = append(variables, VariableModel{
			Id:      generateId(),
			FrameId: frameId,
			Key:     variable.Key,
			Value:   variable.Value,
			Type:    variable.Type,
		})
	}

	var actions []ActionModel

	if len(frameDSL.Blocks) > 0 && frameDSL.Blocks[0].KeyType != "ROOT" {
		return FrameProductionDataWrapper{}, errors.New("first block's keyType must be 'ROOT'")
	}

	blocks, err := processBlocks(frameId, frameDSL.Blocks, "", []BlockSlotModel{}, variables, func(blockActions []ActionModel) {
		actions = append(actions, blockActions...)
	})

	if err != nil {
		return FrameProductionDataWrapper{}, err
	}

	hasDuplicateBlockKey := findDuplicateKeys(blocks)
	if len(hasDuplicateBlockKey) > 0 {
		return FrameProductionDataWrapper{}, errors.New("duplicate block keys found: " + strings.Join(hasDuplicateBlockKey, ","))
	}

	frame := FrameModel{
		Id:             frameId,
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
	production := FrameProductionWrapper{
		FrameProduction: frame,
	}
	wrapper := FrameProductionDataWrapper{
		Data: production,
	}
	return wrapper, err
}

func findBlockVariable(variables []VariableModel, data []BlockDataModel, blockKey string) error {
	for _, dataEntry := range data {
		found := false
		for _, variable := range variables {
			if variable.Key == dataEntry.Value {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no matching variable found for %s block in data entry with key: %s", blockKey, dataEntry.Key)
		}
	}
	return nil
}

func findTriggerVariable(variables []VariableModel, data []TriggerDataModel, triggerName string) error {
	for _, dataEntry := range data {
		found := false
		for _, variable := range variables {
			if variable.Key == dataEntry.Value {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no matching variable found for %s trigger in data entry with key: %s", triggerName, dataEntry.Key)
		}
	}
	return nil
}

func containsSlot(slots []BlockSlotModel, key string) bool {
	for _, slot := range slots {
		if slot.Slot == key {
			return true
		}
	}
	return false
}

func findDuplicateKeys(blocks []BlockModel) []string {
	keyCount := make(map[string]int)
	var duplicates []string

	for _, block := range blocks {
		keyCount[block.Key]++
	}

	for key, count := range keyCount {
		if count > 1 {
			duplicates = append(duplicates, key)
		}
	}
	return duplicates
}
