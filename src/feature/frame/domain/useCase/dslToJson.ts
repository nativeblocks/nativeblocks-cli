import { ActionDSLModel, ActionTriggerDSLModel, BlockDSLModel, FrameDSLModel, generateId } from "../model/dslModel";
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
import { validateJsonWithDynamicSchema } from "./validateSchema";

function processActions(
  frameId: string,
  key: string,
  inputActions: ActionDSLModel[],
  variables: VariableModel[]
): ActionModel[] {
  return inputActions?.map((inputAction) => {
    const actionId = generateId();
    const newAction: ActionModel = {
      id: actionId,
      frameId: frameId,
      key: key,
      event: inputAction.event,
      triggers: processTriggers(actionId, inputAction.triggers, "", variables),
    };
    return newAction;
  });
}

function processTriggers(
  actionId: string,
  triggers: ActionTriggerDSLModel[],
  parentId: string = "",
  variables: VariableModel[]
): ActionTriggerModel[] {
  let flatTriggers: ActionTriggerModel[] = [];

  triggers?.forEach((trigger) => {
    const newTrigger = {
      id: generateId(),
      actionId: actionId,
      parentId: parentId,
      keyType: trigger.keyType,
      then: trigger.then,
      name: trigger.name,
      integrationVersion: trigger.integrationVersion,
      data: [],
      properties: [],
    } as ActionTriggerModel;

    const triggerProperties =
      trigger.properties?.map((property) => {
        return {
          id: generateId(),
          actionTriggerId: newTrigger.id,
          key: property.key,
          type: property.type,
          value: property.value,
          description: property.description ?? "",
          valuePicker: property.valuePicker ?? "",
          valuePickerGroup: property.valuePickerGroup ?? "",
          valuePickerOptions: property.valuePickerOptions ?? "",
        } as TriggerPropertyModel;
      }) ?? [];

    const triggerData =
      trigger.data?.map((dataItem) => {
        return {
          id: generateId(),
          actionTriggerId: newTrigger.id,
          key: dataItem.key,
          value: dataItem.value,
          type: dataItem.type,
          description: dataItem.description ?? "",
        } as TriggerDataModel;
      }) ?? [];

    newTrigger.properties = triggerProperties;
    newTrigger.data = triggerData;

    newTrigger.data?.forEach((dataEntry) => {
      const variable = variables.find((v) => v.key === dataEntry.value);
      if (variable) {
        dataEntry.value = variable.key;
        dataEntry.type = variable.type;
      } else {
        throw new Error(`Variable with key '${dataEntry.value}' not found`);
      }
    });
    flatTriggers.push(newTrigger);

    if (trigger.triggers && trigger.triggers.length > 0) {
      const subTriggerResults = processTriggers(actionId, trigger.triggers, newTrigger.id, variables);
      flatTriggers = flatTriggers.concat(subTriggerResults);
    }
  });
  return flatTriggers;
}

function processBlocks(
  frameId: string,
  blocks: BlockDSLModel[],
  parentId: string = "",
  variables: VariableModel[],
  onNewAction: (actions: ActionModel[]) => void
): BlockModel[] {
  let flatBlocks: BlockModel[] = [];

  blocks?.forEach((block, index) => {
    const slot = block.slot ?? "content";
    const newBlock = {
      id: generateId(),
      frameId: frameId,
      keyType: block.keyType,
      key: block.key,
      visibilityKey: block.visibilityKey,
      position: index,
      slot: slot,
      integrationVersion: block.integrationVersion,
      parentId: parentId,
      data: [],
      properties: [],
      slots: [],
    } as BlockModel;
    onNewAction(processActions(frameId, block.key, block.actions, variables));

    const blockProperties =
      block.properties?.map((property) => {
        if (property.value) {
          return {
            id: generateId(),
            blockId: newBlock.id,
            key: property.key,
            type: property.type,
            valueMobile: property.value,
            valueTablet: property.value,
            valueDesktop: property.value,
            description: property.description ?? "",
            valuePicker: property.valuePicker ?? "",
            valuePickerGroup: property.valuePickerGroup ?? "",
            valuePickerOptions: property.valuePickerOptions ?? "",
          } as BlockPropertyModel;
        } else {
          return {
            id: generateId(),
            blockId: newBlock.id,
            key: property.key,
            type: property.type,
            valueMobile: property.valueMobile,
            valueTablet: property.valueTablet,
            valueDesktop: property.valueDesktop,
            description: property.description ?? "",
            valuePicker: property.valuePicker ?? "",
            valuePickerGroup: property.valuePickerGroup ?? "",
            valuePickerOptions: property.valuePickerOptions ?? "",
          } as BlockPropertyModel;
        }
      }) ?? [];

    const blockData =
      block.data?.map((dataItem) => {
        return {
          id: generateId(),
          blockId: newBlock.id,
          key: dataItem.key,
          value: dataItem.value,
          type: dataItem.type,
          description: dataItem.description ?? "",
        } as BlockDataModel;
      }) ?? [];

    const blocSlots =
      block.slots?.map((slot) => {
        return {
          id: generateId(),
          blockId: newBlock.id,
          slot: slot.slot,
          description: slot.description ?? "",
        } as BlockSlotModel;
      }) ?? [];

    newBlock.properties = blockProperties;
    newBlock.data = blockData;
    newBlock.slots = blocSlots;

    const variable = variables.find((v) => v.key === newBlock.visibilityKey);
    if (!variable) {
      throw new Error(`Variable with key '${newBlock.visibilityKey}' not found`);
    }

    newBlock.data?.forEach((dataEntry) => {
      const variable = variables.find((v) => v.key === dataEntry.value);
      if (variable) {
        dataEntry.value = variable.key;
        dataEntry.type = variable.type;
      } else {
        throw new Error(`Variable with key '${dataEntry.value}' not found`);
      }
    });
    flatBlocks.push(newBlock);

    if (block.blocks && block.blocks.length > 0) {
      const subBlockResults = processBlocks(frameId, block?.blocks ?? [], newBlock.id, variables, onNewAction);
      flatBlocks = flatBlocks.concat(subBlockResults);
    }
  });
  return flatBlocks;
}

function getWordsBetweenCurly(text: string): string[] {
  const regex = /\{(.*?)\}/g;
  const matches = Array.from(text.matchAll(regex));

  const result: string[] = [];
  for (const match of matches) {
    if (match[1]) {
      result.push(match[1]);
    }
  }
  return result;
}

export async function generateFrame(frameDSL: FrameDSLModel) {
  if (!frameDSL.$schema) {
    throw Error("Please provide $schema for the json file");
  }
  const validateResult = await validateJsonWithDynamicSchema(frameDSL, frameDSL.$schema);
  if (validateResult.valid) {
    const frameId = generateId();
    const variables = frameDSL.variables.map((variable) => {
      return {
        id: generateId(),
        frameId: frameId,
        key: variable.key,
        value: variable.value,
        type: variable.type,
      } as VariableModel;
    });
    const actions: ActionModel[] = [];
    const blocks = processBlocks(frameId, frameDSL.blocks, "", variables ?? [], (blockActions) => {
      actions.push(...blockActions);
    });
    return {
      data: {
        frameProduction: {
          id: frameId,
          name: frameDSL.name,
          route: frameDSL.route,
          routeArguments: getWordsBetweenCurly(frameDSL.route).map((arg) => {
            return { name: arg };
          }),
          type: frameDSL.type,
          isStarter: frameDSL.isStarter,
          projectId: "",
          checksum: "",
          variables: variables,
          blocks: blocks,
          actions: actions,
        } as FrameModel,
      },
    };
  } else {
    validateResult.errors?.forEach((e) => {
      console.log(e);
    });
    return { data: null };
  }
}
