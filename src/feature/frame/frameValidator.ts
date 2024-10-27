import {
  ActionModel,
  ActionTriggerModel,
  BlockDataModel,
  BlockModel,
  BlockPropertyModel,
  BlockSlotModel,
  generateId,
  TriggerDataModel,
  TriggerPropertyModel,
} from "./frameModel";
import Ajv from "ajv";

export class VariableModel {
  id: string;
  frameId: string;
  key: string;
  value: string;
  type: string;

  constructor(frameId: string, key: string, value: string, type: string) {
    this.id = generateId();
    this.frameId = frameId;
    this.key = key;
    this.value = value;
    this.type = type;
  }
}

interface InputFrame {
  $schema: string;
  name: string;
  route: string;
  type: string;
  isStarter: boolean;
  variables: VariableModel[];
  blocks: BlockModel[];
  actions: ActionModel[];
}

function processActions(frameId: string, key: string, inputActions: ActionModel[], variables: VariableModel[]): ActionModel[] {
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
  triggers: ActionTriggerModel[],
  parentId: string = "",
  variables: VariableModel[]
): ActionTriggerModel[] {
  let flatTriggers: ActionTriggerModel[] = [];

  triggers?.forEach((trigger) => {
    const newTrigger = new ActionTriggerModel(
      actionId,
      parentId,
      trigger.keyType,
      trigger.then,
      trigger.name,
      trigger.integrationVersion,
      [],
      [],
      trigger.triggers ? processTriggers(actionId, trigger.triggers, "", variables) : []
    );

    const triggerProperties =
      trigger.properties?.map((property) => {
        return new TriggerPropertyModel(
          newTrigger.id,
          property.key,
          property.type,
          property.value,
          property.description ?? "",
          property.valuePicker ?? "",
          property.valuePickerGroup ?? "",
          property.valuePickerOptions ?? ""
        );
      }) ?? [];

    const triggerData =
      trigger.data?.map((dataItem) => {
        return new TriggerDataModel(
          newTrigger.id,
          dataItem.key,
          dataItem.value,
          dataItem.type,
          dataItem.description ?? ""
        );
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
  blocks: BlockModel[],
  parentId: string = "",
  variables: VariableModel[],
  onNewAction: (actions: ActionModel[]) => void
): BlockModel[] {
  let flatBlocks: BlockModel[] = [];

  blocks?.forEach((block, index) => {
    const slot = block.slot ?? "content";
    const newBlock = new BlockModel(
      frameId,
      block.keyType,
      block.key,
      block.visibilityKey,
      block.position ?? index,
      slot,
      block.integrationVersion ?? 0,
      parentId,
      [],
      [],
      [],
      [],
      []
    );
    onNewAction(processActions(frameId, block.key, block.actions, variables))

    const blockProperties =
      block.properties?.map((property) => {
        if (property.value) {
          return new BlockPropertyModel(
            newBlock.id,
            property.key,
            property.type,
            property.value,
            property.value,
            property.value,
            property.value,
            property.description ?? "",
            property.valuePicker ?? "",
            property.valuePickerGroup ?? "",
            property.valuePickerOptions ?? ""
          );
        } else {
          return new BlockPropertyModel(
            newBlock.id,
            property.key,
            property.type,
            property.value,
            property.valueMobile,
            property.valueTablet,
            property.valueDesktop,
            property.description ?? "",
            property.valuePicker ?? "",
            property.valuePickerGroup ?? "",
            property.valuePickerOptions ?? ""
          );
        }
      }) ?? [];

    const blockData =
      block.data?.map((dataItem) => {
        return new BlockDataModel(newBlock.id, dataItem.key, dataItem.value, dataItem.type, dataItem.description ?? "");
      }) ?? [];

    const blocSlots =
      block.slots?.map((slot) => {
        return new BlockSlotModel(newBlock.id, slot.slot, slot.description ?? "");
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

async function validateJsonWithDynamicSchema(jsonData: any, schemaUrl: string) {
  const ajv = new Ajv({ allErrors: true });
  try {
    const schemaResponse = await fetch(schemaUrl);
    const schema = await schemaResponse.json();
    const validate = ajv.compile(schema);
    const isValid = validate(jsonData);
    if (!isValid) {
      return { valid: false, errors: validate.errors };
    }
    return { valid: true, errors: null };
  } catch (error: any) {
    return { valid: false, errors: [error.message] };
  }
}

export async function processFrame(inputFrame: InputFrame) {
  if (!inputFrame.$schema) {
    throw Error("Please provide $schema for the json file");
  }
  const validateResult = await validateJsonWithDynamicSchema(inputFrame, inputFrame.$schema);
  if (validateResult.valid) {
    const frameId = generateId();
    const variables = inputFrame.variables.map((variable) => {
      return new VariableModel(frameId, variable.key, variable.value, variable.type);
    });
    const actions: ActionModel[] = [];
    const flatBlocks = processBlocks(frameId, inputFrame.blocks, "", variables ?? [], (blockActions) => {
      actions.push(...blockActions);
    });
    return {
      data: {
        frameProduction: {
          id: frameId,
          name: inputFrame.name,
          route: inputFrame.route,
          type: inputFrame.type ?? "FRAME",
          isStarter: inputFrame.isStarter,
          checksum: "",
          projectId: "",
          routeArguments: getWordsBetweenCurly(inputFrame.route).map((arg) => {
            return { name: arg };
          }),
          variables: variables,
          blocks: flatBlocks,
          actions: actions
        },
      },
    };
  } else {
    validateResult.errors?.forEach((e) => {
      console.log(e);
    });
    return { data: null };
  }
}
