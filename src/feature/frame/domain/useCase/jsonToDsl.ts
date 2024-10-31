import {
  ActionDSLModel,
  ActionTriggerDSLModel,
  BlockDataDSLModel,
  BlockDSLModel,
  BlockPropertyDSLModel,
  BlockSlotDSLModel,
  FrameDSLModel,
  TriggerDataDSLModel,
  TriggerPropertyDSLModel,
  VariableDSLModel,
} from "../model/dslModel";
import {
  ActionModel,
  ActionTriggerModel,
  BlockDataModel,
  BlockModel,
  BlockPropertyModel,
  BlockSlotModel,
  FrameModel,
  TriggerDataModel,
  TriggerPropertyModel,
  VariableModel,
} from "../model/model";

function findActionTriggerChildren(triggers: ActionTriggerModel[], parentId: string): ActionTriggerDSLModel[] {
  return triggers
    .filter((trigger) => trigger.parentId === parentId)
    .map((trigger) => {
      return {
        keyType: trigger.keyType,
        then: trigger.then,
        name: trigger.name,
        integrationVersion: trigger.integrationVersion,
        properties: trigger.properties.map((prop) => mapTriggerPropertyModelToDSL(prop)),
        data: trigger.data.map((data) => mapTriggerDataModelToDSL(data)),
        triggers: findActionTriggerChildren(triggers, trigger.id ?? ""),
      } as ActionTriggerDSLModel;
    });
}

function buildActionTriggerTree(triggers: ActionTriggerModel[]): ActionTriggerDSLModel[] {
  const roots = triggers.filter((trigger) => trigger.parentId === "");
  return roots.map((root) => {
    return {
      keyType: root.keyType,
      then: root.then,
      name: root.name,
      integrationVersion: root.integrationVersion,
      properties: root.properties.map((prop) => mapTriggerPropertyModelToDSL(prop)),
      data: root.data.map((data) => mapTriggerDataModelToDSL(data)),
      triggers: findActionTriggerChildren(triggers, root.id ?? ""),
    } as ActionTriggerDSLModel;
  });
}

function findBlockChildren(blocks: BlockModel[], parentId: string, actions: ActionModel[]): BlockDSLModel[] {
  return blocks
    .filter((block) => block.parentId === parentId)
    .map((block) => {
      const result = {
        keyType: block.keyType,
        key: block.key,
        visibilityKey: block.visibilityKey,
        slot: block.slot,
        integrationVersion: block.integrationVersion,
        data: block.data.map((dataItem) => mapBlockDataModelToDSL(dataItem)),
        properties: block.properties.map((property) => mapBlockPropertyModelToDSL(property)),
        slots: block.slots.map((slot) => mapBlockSlotModelToDSL(slot)),
        actions:
        actions
          .filter((action) => {
            return action.key === block.key;
          })
          ?.map((act) => {
            return mapActionModelToDSL(act);
          }) ?? [],
        blocks: findBlockChildren(blocks, block.id ?? "", actions),
      };
      return result;
    });
}

function buildBlockTreeWithActions(blocks: BlockModel[], actions: ActionModel[]): BlockDSLModel[] {
  const roots = blocks.filter((block) => block.parentId === "");
  return roots.map((root) => {
    return {
      keyType: root.keyType,
      key: root.key,
      visibilityKey: root.visibilityKey,
      slot: root.slot,
      integrationVersion: root.integrationVersion,
      data: root.data.map((dataItem) => mapBlockDataModelToDSL(dataItem)),
      properties: root.properties.map((property) => mapBlockPropertyModelToDSL(property)),
      slots: root.slots.map((slot) => mapBlockSlotModelToDSL(slot)),
      actions:
        actions
          .filter((action) => {
            return action.key === root.key;
          })
          ?.map((act) => {
            return mapActionModelToDSL(act);
          }) ?? [],
      blocks: findBlockChildren(blocks, root.id ?? "", actions),
    } as BlockDSLModel;
  });
}

function mapVariableModelToDSL(variable: VariableModel): VariableDSLModel {
  return {
    frameId: variable.frameId,
    key: variable.key,
    value: variable.value,
    type: variable.type,
  };
}

function mapBlockDataModelToDSL(data: BlockDataModel): BlockDataDSLModel {
  return {
    key: data.key,
    value: data.value,
    type: data.type,
    description: data.description,
  };
}

function mapBlockPropertyModelToDSL(property: BlockPropertyModel): BlockPropertyDSLModel {
  return {
    key: property.key,
    valueMobile: property.valueMobile,
    valueTablet: property.valueTablet,
    valueDesktop: property.valueDesktop,
    type: property.type,
    description: property.description || "",
    valuePicker: property.valuePicker,
    valuePickerGroup: property.valuePickerGroup,
    valuePickerOptions: property.valuePickerOptions,
  };
}

function mapBlockSlotModelToDSL(slot: BlockSlotModel): BlockSlotDSLModel {
  return {
    slot: slot.slot,
    description: slot.description,
  };
}

function mapActionModelToDSL(action: ActionModel): ActionDSLModel {
  return {
    key: action.key,
    event: action.event,
    triggers: buildActionTriggerTree(action.triggers)
  };
}

function mapTriggerPropertyModelToDSL(property: TriggerPropertyModel): TriggerPropertyDSLModel {
  return {
    key: property.key,
    value: property.value,
    type: property.type,
    description: property.description || "",
    valuePicker: property.valuePicker,
    valuePickerGroup: property.valuePickerGroup,
    valuePickerOptions: property.valuePickerOptions,
  };
}

function mapTriggerDataModelToDSL(data: TriggerDataModel): TriggerDataDSLModel {
  return {
    key: data.key,
    value: data.value,
    type: data.type,
    description: data.description,
  };
}

export function mapFrameModelToDSL(frame: FrameModel): FrameDSLModel {
  return {
    $schema: "https://your-schema-url.com",
    name: frame.name,
    route: frame.route,
    type: frame.type,
    isStarter: frame.isStarter,
    variables: frame.variables.map((variable) => mapVariableModelToDSL(variable)),
    blocks: buildBlockTreeWithActions(frame.blocks, frame.actions),
  };
}