import { ResultModel } from "../../../infrastructure/result/model/ResultModel";
import { IntegrationsInstalledModel } from "./model/installedIntegrationModel";
import { ProjectModel } from "./projectRepositoryImpl";

export interface ProjectRepository {
  getInstalledIntegrations(organizationId: string, projectId: string, kind: string): Promise<ResultModel<IntegrationsInstalledModel[]>>;
  projects(organizationId: string): Promise<ResultModel<ProjectModel[]>>;
  set(id: string): Promise<ResultModel<string>>;
  get(): Promise<ResultModel<ProjectModel>>;
}
