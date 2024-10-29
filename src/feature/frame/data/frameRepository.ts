import { ResultModel } from "../../../infrastructure/result/model/ResultModel";
import { FrameModel } from "../domain/model/model";

export interface FrameRepository {
  pushFrame(apiKey: string, route: string, frameJson: string): Promise<ResultModel<any>>;
  pullFrame(apiKey: string, route: string): Promise<ResultModel<FrameModel>>;
}
