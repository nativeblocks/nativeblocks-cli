import { ResultModel } from "../../../infrastructure/result/model/ResultModel";
import { ProjectModel } from "./projectRepositoryImpl";

export interface ProjectRepository {
  getInstalledIntegrations(apiKey: string): Promise<ResultModel<any>>;
  projects(organizationId: string): Promise<ResultModel<ProjectModel[]>>;
  set(id: string): Promise<ResultModel<string>>;
  get(): Promise<ResultModel<ProjectModel>>;
}
