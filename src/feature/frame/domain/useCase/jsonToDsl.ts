import { ActionTriggerModel, BlockModel } from "../model/model";

function findActionTriggerChildren(triggers: ActionTriggerModel[], parentId: string): ActionTriggerModel[] {
  return triggers
    .filter((trigger) => trigger.parentId === parentId)
    .map((trigger) => {
      return {
        ...trigger,
        triggers: findActionTriggerChildren(triggers, trigger.id ?? ""),
      };
    });
}

export function buildActionTriggerTree(triggers: ActionTriggerModel[]): ActionTriggerModel[] {
  const roots = triggers.filter((trigger) => trigger.parentId === "");
  return roots.map((root) => {
    return {
      ...root,
      triggers: findActionTriggerChildren(triggers, root.id ?? ""),
    };
  });
}

function findBlockChildren(blocks: BlockModel[], parentId: string): BlockModel[] {
  return blocks
    .filter((block) => block.parentId === parentId)
    .map((block) => {
      const result = {
        ...block,
        blocks: findBlockChildren(blocks, block.id ?? ""),
      };
      return result;
    });
}

export function buildBlockTree(blocks: BlockModel[]): BlockModel[] {
  const roots = blocks.filter((block) => block.parentId === "");
  return roots.map((root) => {
    return {
      ...root,
      blocks: findBlockChildren(blocks, root.id ?? ""),
    };
  });
}

