export interface IntegrationsInstalledModel {
  integrationKeyType: string;
  integrationVersion: number;
  integrationId: string;
  integrationPlatformSupport: string;
  integrationKind: string;
  integrationProperties: IntegrationPropertyModel[];
  integrationData: IntegrationDataModel[];
  integrationEvents: IntegrationEventModel[];
  integrationSlots: IntegrationSlotModel[];
}

export interface IntegrationPropertyModel {
  key: string;
  value: string;
  type: string;
}

export interface IntegrationDataModel {
  key: string;
  type: string;
}

export interface IntegrationEventModel {
  event: string;
}

export interface IntegrationSlotModel {
  slot: string;
}
