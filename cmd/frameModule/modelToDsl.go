package frameModule

func findActionTriggerChildren(triggers []ActionTriggerModel, parentId string) []ActionTriggerDSLModel {
	var children []ActionTriggerDSLModel

	for _, trigger := range triggers {
		if trigger.ParentId == parentId {
			child := ActionTriggerDSLModel{
				KeyType:            trigger.KeyType,
				Then:               trigger.Then,
				Name:               trigger.Name,
				IntegrationVersion: trigger.IntegrationVersion,
				Properties:         make([]TriggerPropertyDSLModel, len(trigger.Properties)),
				Data:               make([]TriggerDataDSLModel, len(trigger.Data)),
			}

			for i, prop := range trigger.Properties {
				child.Properties[i] = mapTriggerPropertyModelToDSL(prop)
			}

			for i, data := range trigger.Data {
				child.Data[i] = mapTriggerDataModelToDSL(data)
			}

			child.Triggers = findActionTriggerChildren(triggers, trigger.Id)
			children = append(children, child)
		}
	}

	if children == nil {
		return make([]ActionTriggerDSLModel, 0)
	} else {
		return children
	}
}

func buildActionTriggerTree(triggers []ActionTriggerModel) []ActionTriggerDSLModel {
	var roots []ActionTriggerDSLModel

	for _, trigger := range triggers {
		if trigger.ParentId == "" {
			root := ActionTriggerDSLModel{
				KeyType:            trigger.KeyType,
				Then:               trigger.Then,
				Name:               trigger.Name,
				IntegrationVersion: trigger.IntegrationVersion,
				Properties:         make([]TriggerPropertyDSLModel, len(trigger.Properties)),
				Data:               make([]TriggerDataDSLModel, len(trigger.Data)),
			}

			for i, prop := range trigger.Properties {
				root.Properties[i] = mapTriggerPropertyModelToDSL(prop)
			}

			for i, data := range trigger.Data {
				root.Data[i] = mapTriggerDataModelToDSL(data)
			}

			root.Triggers = findActionTriggerChildren(triggers, trigger.Id)
			roots = append(roots, root)
		}
	}

	if roots == nil {
		return make([]ActionTriggerDSLModel, 0)
	} else {
		return roots
	}
}

func findBlockChildren(blocks []BlockModel, parentId string, actions []ActionModel) []BlockDSLModel {
	var children []BlockDSLModel

	for _, block := range blocks {
		if block.ParentId == parentId {
			child := BlockDSLModel{
				KeyType:            block.KeyType,
				Key:                block.Key,
				VisibilityKey:      block.VisibilityKey,
				Slot:               block.Slot,
				IntegrationVersion: block.IntegrationVersion,
				Data:               make([]BlockDataDSLModel, len(block.Data)),
				Properties:         make([]BlockPropertyDSLModel, len(block.Properties)),
				Slots:              make([]BlockSlotDSLModel, len(block.Slots)),
			}

			for i, data := range block.Data {
				child.Data[i] = mapBlockDataModelToDSL(data)
			}

			for i, prop := range block.Properties {
				child.Properties[i] = mapBlockPropertyModelToDSL(prop)
			}

			for i, slot := range block.Slots {
				child.Slots[i] = mapBlockSlotModelToDSL(slot)
			}

			for _, action := range actions {
				if action.Key == block.Key {
					child.Actions = append(child.Actions, mapActionModelToDSL(action))
				}
			}

			if child.Actions == nil {
				child.Actions = make([]ActionDSLModel, 0)
			}

			child.Blocks = findBlockChildren(blocks, block.Id, actions)
			children = append(children, child)
		}
	}

	if children == nil {
		return make([]BlockDSLModel, 0)
	} else {
		return children
	}
}

func buildBlockTreeWithActions(blocks []BlockModel, actions []ActionModel) []BlockDSLModel {
	var dslBlocks []BlockDSLModel
	for _, block := range blocks {
		if block.ParentId == "" {
			root := BlockDSLModel{
				KeyType:            block.KeyType,
				Key:                block.Key,
				VisibilityKey:      block.VisibilityKey,
				Slot:               "null",
				IntegrationVersion: block.IntegrationVersion,
				Data:               make([]BlockDataDSLModel, len(block.Data)),
				Properties:         make([]BlockPropertyDSLModel, len(block.Properties)),
				Slots:              make([]BlockSlotDSLModel, len(block.Slots)),
			}

			for i, data := range block.Data {
				root.Data[i] = mapBlockDataModelToDSL(data)
			}

			for i, prop := range block.Properties {
				root.Properties[i] = mapBlockPropertyModelToDSL(prop)
			}

			for i, slot := range block.Slots {
				root.Slots[i] = mapBlockSlotModelToDSL(slot)
			}

			for _, action := range actions {
				if action.Key == block.Key {
					root.Actions = append(root.Actions, mapActionModelToDSL(action))
				}
			}

			if root.Actions == nil {
				root.Actions = make([]ActionDSLModel, 0)
			}

			root.Blocks = findBlockChildren(blocks, block.Id, actions)
			dslBlocks = append(dslBlocks, root)
		}
	}

	if dslBlocks == nil {
		return make([]BlockDSLModel, 0)
	} else {
		return dslBlocks
	}
}

func mapVariableModelToDSL(variable VariableModel) VariableDSLModel {
	return VariableDSLModel{
		FrameId: variable.FrameId,
		Key:     variable.Key,
		Value:   variable.Value,
		Type:    variable.Type,
	}
}

func mapBlockDataModelToDSL(data BlockDataModel) BlockDataDSLModel {
	return BlockDataDSLModel{
		Key:   data.Key,
		Value: data.Value,
		Type:  data.Type,
	}
}

func mapBlockPropertyModelToDSL(property BlockPropertyModel) BlockPropertyDSLModel {
	return BlockPropertyDSLModel{
		Key:          property.Key,
		ValueMobile:  property.ValueMobile,
		ValueTablet:  property.ValueTablet,
		ValueDesktop: property.ValueDesktop,
		Type:         property.Type,
	}
}

func mapBlockSlotModelToDSL(slot BlockSlotModel) BlockSlotDSLModel {
	return BlockSlotDSLModel{
		Slot: slot.Slot,
	}
}

func mapActionModelToDSL(action ActionModel) ActionDSLModel {
	return ActionDSLModel{
		Key:      action.Key,
		Event:    action.Event,
		Triggers: buildActionTriggerTree(action.Triggers),
	}
}

func mapTriggerPropertyModelToDSL(property TriggerPropertyModel) TriggerPropertyDSLModel {
	return TriggerPropertyDSLModel{
		Key:   property.Key,
		Value: property.Value,
		Type:  property.Type,
	}
}

func mapTriggerDataModelToDSL(data TriggerDataModel) TriggerDataDSLModel {
	return TriggerDataDSLModel{
		Key:   data.Key,
		Value: data.Value,
		Type:  data.Type,
	}
}

func mapFrameModelToDSL(frame FrameModel, schema string) FrameDSLModel {
	variables := make([]VariableDSLModel, len(frame.Variables))
	for i, variable := range frame.Variables {
		variables[i] = mapVariableModelToDSL(variable)
	}

	return FrameDSLModel{
		Schema:    schema,
		Name:      frame.Name,
		Route:     frame.Route,
		Type:      frame.Type,
		IsStarter: frame.IsStarter,
		Variables: variables,
		Blocks:    buildBlockTreeWithActions(frame.Blocks, frame.Actions),
	}
}
