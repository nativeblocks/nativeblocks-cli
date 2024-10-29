export interface FrameModel {
  id: string;
  name: string;
  route: string;
  routeArguments: { name: string }[];
  type: string;
  isStarter: boolean;
  projectId: string;
  checksum: string;
  variables: VariableModel[];
  blocks: BlockModel[];
  actions: ActionModel[];
}

export interface VariableModel {
  id: string;
  frameId: string;
  key: string;
  value: string;
  type: string;
}

export interface BlockModel {
  id: string;
  frameId: string;
  keyType: string;
  key: string;
  visibilityKey: string;
  position: number;
  slot: string | null;
  integrationVersion: number;
  parentId?: string | null;
  data: BlockDataModel[];
  properties: BlockPropertyModel[];
  slots: BlockSlotModel[];
}

export interface BlockPropertyModel {
  id: string;
  blockId: string;
  key: string;
  valueMobile: string;
  valueTablet: string;
  valueDesktop: string;
  type: string;
  description?: string;
  valuePicker: string;
  valuePickerGroup: string;
  valuePickerOptions: string;
}

export interface BlockDataModel {
  id: string;
  blockId: string;
  key: string;
  value: string;
  type: string;
  description?: string;
}

export interface BlockSlotModel {
  id: string;
  blockId: string;
  slot: string;
  description?: string;
}

export interface ActionModel {
  id: string;
  frameId: string;
  key: string;
  event: string;
  triggers: ActionTriggerModel[];
}

export interface ActionTriggerModel {
  id: string;
  actionId: string;
  parentId: string | null;
  keyType: string;
  then: string;
  name: string;
  integrationVersion: number;
  properties: TriggerPropertyModel[];
  data: TriggerDataModel[];
}

export interface TriggerPropertyModel {
  id: string;
  actionTriggerId: string;
  key: string;
  value: string;
  type: string;
  description?: string;
  valuePicker?: string;
  valuePickerGroup?: string;
  valuePickerOptions?: string;
}

export interface TriggerDataModel {
  id: string;
  actionTriggerId: string;
  key: string;
  value: string;
  type: string;
  description?: string;
}