import { v4 as uuidv4 } from "uuid";

export function generateId(): string {
  return uuidv4();
}

export class BlockModel {
  id: string;
  frameId: string;
  keyType: string;
  key: string;
  visibilityKey: string;
  position: number;
  slot: string | null;
  integrationVersion: number;
  parentId: string | null;
  blocks?: BlockModel[] | null;
  data: BlockDataModel[] | null;
  properties: BlockPropertyModel[];
  slots: BlockSlotModel[];
  actions: ActionModel[];

  constructor(
    frameId: string,
    keyType: string,
    key: string,
    visibilityKey: string,
    position: number,
    slot: string | null,
    integrationVersion: number,
    parentId: string | null,
    blocks: BlockModel[] | null,
    data: BlockDataModel[],
    properties: BlockPropertyModel[],
    slots: BlockSlotModel[],
    actions: ActionModel[]
  ) {
    this.id = generateId();
    this.frameId = frameId;
    this.keyType = keyType;
    this.key = key;
    this.visibilityKey = visibilityKey;
    this.position = position;
    this.slot = slot;
    this.integrationVersion = integrationVersion;
    this.parentId = parentId;
    this.blocks = blocks;
    this.data = data;
    this.properties = properties;
    this.slots = slots;
    this.actions = actions;
  }
}

export class BlockPropertyModel {
  id: string;
  key: string;
  value?: string;
  valueMobile?: string;
  valueTablet?: string;
  valueDesktop?: string;
  type: string;
  description?: string;
  valuePicker?: string;
  valuePickerGroup?: string;
  valuePickerOptions?: string;
  blockId: string;

  constructor(
    blockId: string,
    key: string,
    type: string,
    value?: string,
    valueMobile?: string,
    valueTablet?: string,
    valueDesktop?: string,
    description?: string,
    valuePicker?: string,
    valuePickerGroup?: string,
    valuePickerOptions?: string
  ) {
    this.id = generateId();
    this.key = key;
    this.value = value;
    this.valueMobile = valueMobile;
    this.valueTablet = valueTablet;
    this.valueDesktop = valueDesktop;
    this.type = type;
    this.description = description;
    this.valuePicker = valuePicker;
    this.valuePickerGroup = valuePickerGroup;
    this.valuePickerOptions = valuePickerOptions ?? "[]";
    this.blockId = blockId;
  }
}

export class BlockDataModel {
  id: string;
  blockId: string;
  key: string;
  value: string;
  type: string;
  description?: string;

  constructor(blockId: string, key: string, value: string, type: string, description?: string) {
    this.id = generateId();
    this.blockId = blockId;
    this.key = key;
    this.value = value;
    this.type = type;
    this.description = description;
  }
}

export class BlockEventModel {
  eventName: string;
  actions: any[];

  constructor(eventName: string, actions: any[] = []) {
    this.eventName = eventName;
    this.actions = actions;
  }
}

export class BlockSlotModel {
  id: string;
  blockId: string;
  slot: string;
  description?: string;

  constructor(blockId: string, slot: string, description?: string) {
    this.id = generateId();
    this.blockId = blockId;
    this.slot = slot;
    this.description = description;
  }
}

export interface ActionModel {
  id: string;
  frameId: string;
  key: string;
  event: string;
  triggers: ActionTriggerModel[];
}

export class ActionTriggerModel {
  id: string;
  actionId: string;
  parentId: string | null;
  keyType: string;
  then: string;
  name: string;
  integrationVersion: number;
  properties: TriggerPropertyModel[];
  data: TriggerDataModel[];
  triggers: ActionTriggerModel[];

  constructor(
    actionId: string,
    parentId: string | null,
    keyType: string,
    then: string,
    name: string,
    integrationVersion: number,
    properties: TriggerPropertyModel[],
    data: TriggerDataModel[],
    triggers: ActionTriggerModel[]
  ) {
    this.id = generateId();
    this.actionId = actionId;
    this.parentId = parentId;
    this.keyType = keyType;
    this.then = then;
    this.name = name;
    this.integrationVersion = integrationVersion;
    this.properties = properties;
    this.data = data;
    this.triggers = triggers;
  }
}

export class TriggerPropertyModel {
  id: string;
  actionTriggerId: string;
  key: string;
  value?: string;
  type: string;
  description?: string;
  valuePicker?: string;
  valuePickerGroup?: string;
  valuePickerOptions?: string;

  constructor(
    actionTriggerId: string,
    key: string,
    type: string,
    value?: string,
    description?: string,
    valuePicker?: string,
    valuePickerGroup?: string,
    valuePickerOptions?: string
  ) {
    this.id = generateId();
    this.actionTriggerId = actionTriggerId;
    this.key = key;
    this.value = value;
    this.type = type;
    this.description = description;
    this.valuePicker = valuePicker;
    this.valuePickerGroup = valuePickerGroup;
    this.valuePickerOptions = valuePickerOptions ?? "[]";
  }
}

export class TriggerDataModel {
  id: string;
  actionTriggerId: string;
  key: string;
  value: string;
  type: string;
  description?: string;

  constructor(actionTriggerId: string, key: string, value: string, type: string, description?: string) {
    this.id = generateId();
    this.actionTriggerId = actionTriggerId;
    this.key = key;
    this.value = value;
    this.type = type;
    this.description = description;
  }
}
