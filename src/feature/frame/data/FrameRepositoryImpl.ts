import { gql } from "graphql-request";
import { getGraphqlClient, handleNetworkError } from "../../../infrastructure/network/NetworkComponent";
import { ResultModel } from "../../../infrastructure/result/model/ResultModel";
import { FrameModel } from "../domain/model/model";
import { FrameRepository } from "./frameRepository";

const SYNC_FRAME_MUTATION = gql`
  mutation syncFrame($input: SyncFrameInput!) {
    syncFrame(input: $input) {
      id
    }
  }
`;

const GET_FRAME_QUERY = gql`
  query frame($route: String!) {
    frame(route: $route) {
      id
      name
      route
      isStarter
      type
      variables {
        key
        value
        type
      }
      blocks {
        id
        parentId
        slot
        keyType
        key
        visibilityKey
        position
        properties {
          key
          valueDesktop
          valueMobile
          valueTablet
          type
        }
        data {
          key
          value
          type
        }
        slots {
          slot
        }
      }
      actions {
        key
        event
        triggers {
          id
          parentId
          keyType
          then
          name
          properties {
            key
            value
            type
          }
          data {
            key
            value
            type
          }
        }
      }
    }
  }
`;

class FrameRepositoryImpl implements FrameRepository {
  private readonly graphqlClient: any;

  constructor(graphqlClient: any) {
    this.graphqlClient = graphqlClient;
  }

  async pushFrame(apiKey: string, route: string, frameJson: string): Promise<ResultModel<any>> {
    try {
      const result = await this.graphqlClient.request(
        SYNC_FRAME_MUTATION,
        {
          input: {
            route: route,
            frameJson: frameJson,
          },
        },
        {
          "Api-Key": `Bearer ${apiKey}`,
        }
      );
      return {
        onSuccess: result.syncFrame?.id,
      };
    } catch (error: any) {
      return {
        onError: handleNetworkError(error).errorMessage,
      };
    }
  }

  async pullFrame(apiKey: string, route: string): Promise<ResultModel<FrameModel>> {
    try {
      const result = await this.graphqlClient.request(
        GET_FRAME_QUERY,
        {
          route: route,
        },
        {
          "Api-Key": `Bearer ${apiKey}`,
        }
      );
      return <ResultModel<FrameModel>>{
        onSuccess: result?.frame,
      };
    } catch (error: any) {
      return <ResultModel<FrameModel>>{
        onError: handleNetworkError(error).errorMessage,
      };
    }
  }
}

export const frameRepository: FrameRepository = new FrameRepositoryImpl(getGraphqlClient());
