import fs from "fs";
import { gql } from "graphql-request";
import os from "os";
import path from "path";
import { getGraphqlClient, handleNetworkError } from "../../../infrastructure/network/NetworkComponent";
import { ResultModel } from "../../../infrastructure/result/model/ResultModel";
import { ProjectRepository } from "./projectRepository";

export const GET_PROJECTS_QUERY = gql`
  query projects($organizationId: String!) {
    projects(organizationId: $organizationId) {
      id
      name
      platform
      apiKeys {
        name
        apiKey
        expireAt
      }
    }
  }
`;

export const GET_PROJECT_QUERY = gql`
  query project($id: String!) {
    project(id: $id) {
      id
      name
      platform
      apiKeys {
        name
        apiKey
        expireAt
      }
      organization {
        id
      }
    }
  }
`;

export type ProjectModel = {
  id: string;
  name: string;
  platform: string;
  apiKeys: ProjectApiKeyModel[];
};

export type ProjectApiKeyModel = {
  name: string;
  apiKey: string;
};

class ProjectRepositoryImpl implements ProjectRepository {
  private readonly graphqlClient: any;

  private userHomeDir: string = os.homedir();
  private projectPath: string = path.join(this.userHomeDir, ".nativeblocks/cli/project.json");
  private directory = path.dirname(this.projectPath);

  constructor(graphqlClient: any) {
    this.graphqlClient = graphqlClient;
  }

  async getInstalledIntegrations(projectId: string): Promise<ResultModel<any>> {
    throw Error();
  }

  async projects(organizationId: string): Promise<ResultModel<ProjectModel[]>> {
    try {
      const result = await this.graphqlClient.request(GET_PROJECTS_QUERY, {
        organizationId: organizationId,
      });
      return {
        onSuccess:
          result.projects?.map((item: any) => {
            return {
              id: item?.id ?? "",
              name: item?.name ?? "",
              platform: item?.platform ?? "",
            } as ProjectModel;
          }) ?? [],
      };
    } catch (error: any) {
      return {
        onError: handleNetworkError(error).errorMessage,
      };
    }
  }

  async set(id: string): Promise<ResultModel<string>> {
    if (!fs.existsSync(this.directory)) {
      fs.mkdirSync(this.directory, { recursive: true });
    }
    try {
      const result = await this.graphqlClient.request(GET_PROJECT_QUERY, {
        id: id,
      });
      const json = {
        id: result.project?.id ?? "",
        name: result.project?.name ?? "",
        platform: result.project?.platform ?? "",
        apiKeys:
          result.project?.apiKeys?.map((key: any) => {
            return {
              name: key.name ?? "",
              apiKey: key.apiKey ?? "",
            } as ProjectApiKeyModel;
          }) ?? [],
      } as ProjectModel;
      fs.writeFileSync(this.projectPath, JSON.stringify(json));
      return {
        onSuccess: `project saved to file successfully at ${this.projectPath}`,
      };
    } catch (e) {
      return {
        onError: `project could not save ${e}`,
      };
    }
  }

  async get(): Promise<ResultModel<ProjectModel>> {
    if (fs.existsSync(this.projectPath)) {
      try {
        const data: string = fs.readFileSync(this.projectPath, "utf-8");
        const json = JSON.parse(data);
        return {
          onSuccess: json,
        };
      } catch (e) {
        return {
          onError: `project could not retrieve ${e}`,
        };
      }
    } else {
      return {
        onError: `project could not retrieve from ${this.projectPath}, please set the organization id`,
      };
    }
  }
}

export const projectRepository: ProjectRepository = new ProjectRepositoryImpl(getGraphqlClient());
