import { gql } from "graphql-request";
import { getGraphqlClient, handleNetworkError } from "../../../infrastructure/network/NetworkComponent";
import { ResultModel } from "../../../infrastructure/result/model/ResultModel";
import { FrameRepository } from "./frameRepository";

const SYNC_FRAME_MUTATION = gql`
  mutation syncFrame($input: SyncFrameInput!) {
    syncFrame(input: $input) {
      id
    }
  }
`;

class FrameRepositoryImpl implements FrameRepository {
  private readonly graphqlClient: any;

  constructor(graphqlClient: any) {
    this.graphqlClient = graphqlClient;
  }

  async syncFrame(apiKey: string, route: string, frameJson: string): Promise<ResultModel<any>> {
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
}

export const frameRepository: FrameRepository = new FrameRepositoryImpl(getGraphqlClient());
