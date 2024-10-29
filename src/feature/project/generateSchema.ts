export function generateBaseSchema(
  blockKeyTypes: Set<string>,
  actionKeyTypes: Set<string>,
  blockProperties: { key: string; value: string; type: string }[],
  blockData: { key: string; value: string; type: string }[],
  actionProperties: { key: string; value: string; type: string }[],
  actionData: { key: string; value: string; type: string }[]
): any {
  const baseSchema = {
    $schema: "http://json-schema.org/draft-07/schema#",
    type: "object",
    required: ["name", "route", "isStarter", "type", "variables", "blocks"],
    properties: {
      name: {
        type: "string",
      },
      route: {
        type: "string",
      },
      type: {
        type: "string",
        enum: ["FRAME", "BOTTOM_SHEET", "DIALOG"]
      },
      isStarter: {
        type: "boolean",
      },
      variables: {
        type: "array",
        uniqueItems: true,
        items: {
          type: "object",
          required: ["key", "value", "type"],
          properties: {
            key: {
              type: "string",
            },
            value: {
              type: "string",
            },
            type: {
              type: "string",
              enum: ["STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"],
            },
          },
        },
      },
      blocks: {
        type: "array",
        maxItems: 1,
        items: {
          $ref: "#/definitions/block",
        },
      },
    },
    definitions: {
      block: {
        type: "object",
        required: [
          "keyType",
          "key",
          "visibilityKey",
          "slot",
          "slots",
          "integrationVersion",
          "data",
          "properties",
          "actions",
          "blocks",
        ],
        properties: {
          keyType: {
            type: "string",
            enum: Array.from(blockKeyTypes),
          },
          key: {
            type: "string",
          },
          visibilityKey: {
            type: "string",
          },
          slot: {
            type: "string",
          },
          slots: {
            type: "array",
            items: {
              type: "object",
              properties: {
                slot: {
                  type: "string",
                },
              },
            },
          },
          integrationVersion: {
            type: "integer",
          },
          data: {
            type: "array",
            items: {
              type: "object",
              required: ["key", "value", "type"],
              properties: {
                key: {
                  type: "string",
                  enum: Array.from(
                    new Set(
                      blockData.map((item) => {
                        return item.key;
                      })
                    )
                  ),
                },
                value: {
                  type: "string",
                },
                type: {
                  type: "string",
                  enum: ["STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"],
                },
              },
            },
          },
          properties: {
            type: "array",
            items: {
              type: "object",
              required: ["key", "type"],
              properties: {
                key: {
                  type: "string",
                  enum: Array.from(
                    new Set(
                      blockProperties.map((item) => {
                        return item.key;
                      })
                    )
                  ),
                },
                value: {
                  type: "string"
                },
                valueMobile: {
                  type: "string"
                },
                valueTablet: {
                  type: "string"
                },
                valueDesktop: {
                  type: "string"
                },
                type: {
                  type: "string",
                  enum: ["STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"],
                },
              },
            },
          },
          actions: {
            type: "array",
            items: {
              type: "object",
              required: ["event", "triggers"],
              properties: {
                event: {
                  type: "string",
                },
                triggers: {
                  type: "array",
                  items: {
                    $ref: "#/definitions/trigger",
                  },
                },
              },
            },
          },
          blocks: {
            type: "array",
            items: {
              $ref: "#/definitions/block",
            },
          },
        },
      },
      trigger: {
        type: "object",
        required: ["keyType", "then", "name", "integrationVersion", "properties", "data", "triggers"],
        properties: {
          keyType: {
            type: "string",
            enum: Array.from(actionKeyTypes),
          },
          then: {
            type: "string",
            enum: ["NEXT", "END", "SUCCESS", "FAILURE"],
          },
          integrationVersion: {
            type: "integer",
          },
          name: {
            type: "string",
          },
          properties: {
            type: "array",
            items: {
              type: "object",
              required: ["key", "value", "type"],
              properties: {
                key: {
                  type: "string",
                  enum: Array.from(
                    new Set(
                      actionProperties.map((item) => {
                        return item.key;
                      })
                    )
                  ),
                },
                value: {
                  type: "string"
                },
                type: {
                  type: "string",
                  enum: ["STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"],
                },
              },
            },
          },
          data: {
            type: "array",
            items: {
              type: "object",
              required: ["key", "value", "type"],
              properties: {
                key: {
                  type: "string",
                  enum: Array.from(
                    new Set(
                      actionData.map((item) => {
                        return item.key;
                      })
                    )
                  ),
                },
                value: {
                  type: "string",
                },
                type: {
                  type: "string",
                  enum: ["STRING", "INT", "LONG", "DOUBLE", "FLOAT", "BOOLEAN"],
                },
              },
            },
          },
          triggers: {
            type: "array",
            items: {
              $ref: "#/definitions/trigger",
            },
          },
        },
      },
    },
  };

  return baseSchema;
}
