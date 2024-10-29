import { v4 as uuidv4 } from "uuid";

export function generateId(): string {
  return uuidv4();
}

export interface FrameDSLModel {
  $schema: string;
  name: string;
  route: string;
  type: string;
  isStarter: boolean;
  variables: VariableDSLModel[];
  blocks: BlockDSLModel[];
  actions: ActionDSLModel[];
}

export class VariableDSLModel {
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

export interface BlockDSLModel {
  keyType: string;
  key: string;
  visibilityKey: string;
  slot: string | null;
  integrationVersion: number;
  blocks?: BlockDSLModel[] | null;
  data: BlockDataDSLModel[];
  properties: BlockPropertyDSLModel[];
  slots: BlockSlotDSLModel[];
  actions: ActionDSLModel[];
}

export interface BlockPropertyDSLModel {
  key: string;
  value: string;
  valueMobile: string;
  valueTablet: string;
  valueDesktop: string;
  type: string;
  description: string;
  valuePicker: string;
  valuePickerGroup: string;
  valuePickerOptions: string;
}

export interface BlockDataDSLModel {
  key: string;
  value: string;
  type: string;
  description?: string;
}

export interface BlockEventDSLModel {
  eventName: string;
  actions: ActionDSLModel[];
}

export interface BlockSlotDSLModel {
  slot: string;
  description?: string;
}

export interface ActionDSLModel {
  key: string;
  event: string;
  triggers: ActionTriggerDSLModel[];
}

export interface ActionTriggerDSLModel {
  keyType: string;
  then: string;
  name: string;
  integrationVersion: number;
  properties: TriggerPropertyDSLModel[];
  data: TriggerDataDSLModel[];
  triggers: ActionTriggerDSLModel[];
}

export interface TriggerPropertyDSLModel {
  key: string;
  value: string;
  type: string;
  description?: string;
  valuePicker?: string;
  valuePickerGroup?: string;
  valuePickerOptions?: string;
}

export interface TriggerDataDSLModel {
  key: string;
  value: string;
  type: string;
  description?: string;
}
