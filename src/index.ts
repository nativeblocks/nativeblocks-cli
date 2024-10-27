import { Command } from "commander";
import { addIntegration, integration, integrations, syncIntegration } from "./feature/integration/integration";
import { integrationData, syncIntegrationData } from "./feature/integration/integrationData";
import { integrationEvent, syncIntegrationEvent } from "./feature/integration/integrationEvent";
import { integrationProperty, syncIntegrationProperty } from "./feature/integration/integrationProperty";
import { integrationSlot, syncIntegrationSlot } from "./feature/integration/integrationSlot";
import { auth } from "./feature/login/auth";
import { getOrganization, organization, organizationList, setOrganization } from "./feature/organization/organization";
import { getRegion, region, setRegion } from "./feature/region/region";
import { frame, generateFrame, generateSchema, syncFrame } from "./feature/frame/frame";

const program = new Command();

program
  .name("nativeblocks")
  .description("Nativeblocks CLI for integration management")
  .version(`${process.env.npm_package_version}`);

const regionCommand = region(program);
setRegion(regionCommand);
getRegion(regionCommand);

auth(program);

const organizationCommand = organization(program);
organizationList(organizationCommand);
setOrganization(organizationCommand);
getOrganization(organizationCommand);

const integrationCommand = integration(program);
integrations(integrationCommand);
addIntegration(integrationCommand);
syncIntegration(integrationCommand);

const eventCommand = integrationEvent(integrationCommand);
syncIntegrationEvent(eventCommand);

const propertyCommand = integrationProperty(integrationCommand);
syncIntegrationProperty(propertyCommand);

const dataCommand = integrationData(integrationCommand);
syncIntegrationData(dataCommand);

const slotCommand = integrationSlot(integrationCommand);
syncIntegrationSlot(slotCommand);

const frameCommand = frame(program);
generateFrame(frameCommand);
syncFrame(frameCommand);
generateSchema(frameCommand);

program.parse();
