import { ResultModel } from "../../../infrastructure/result/model/ResultModel";

export interface FrameRepository {
  syncFrame(apiKey: string, route: string, frameJson: string): Promise<ResultModel<any>>;
}
