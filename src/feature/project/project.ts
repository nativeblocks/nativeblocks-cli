import { Command } from "commander";
import fs from "fs";
import path from "path";
import { findData, findKeyTypes, findProperties } from "./extractKeyType";
import { generateBaseSchema } from "./generateSchema";
import { createDefaultDir } from "../../infrastructure/utility/FileUitl";

export function project(program: Command) {
  return program.command("project").description("Manage project");
}

export function generateSchema(program: Command) {
  return program
    .command("gen-schema")
    .description("Generate a frame")
    .option("-d, --directory", "Frame working root directory")
    .argument("<directory>", "Frame working root directory")
    .action((directory) => {
      const jsonFilePath: string = path.join(directory + "/.nativeblocks/integrations");
      if (fs.existsSync(jsonFilePath)) {
        try {
          let blockKeyTypes = [] as any[];
          let blockPropertyKeys = [] as any[];
          let blockDataKeys = [] as any[];
          let actionKeyTypes = [] as any[];
          let actionPropertyKeys = [] as any[];
          let actionDataKeys = [] as any[];
          
          if (fs.existsSync(jsonFilePath + "/block")) {
            blockKeyTypes = findKeyTypes(jsonFilePath + "/block");
            blockPropertyKeys = findProperties(jsonFilePath + "/block");
            blockDataKeys = findData(jsonFilePath + "/block");
          }

          if (fs.existsSync(jsonFilePath + "/action")) {
            actionKeyTypes = findKeyTypes(jsonFilePath + "/action");
            actionPropertyKeys = findProperties(jsonFilePath + "/action");
            actionDataKeys = findData(jsonFilePath + "/action");
          }

          const schema = generateBaseSchema(
            new Set(["ROOT", ...blockKeyTypes]),
            new Set(actionKeyTypes),
            blockPropertyKeys,
            blockDataKeys,
            actionPropertyKeys,
            actionDataKeys
          );

          console.log("=========================================================================================");
          const path = createDefaultDir(directory);
          fs.writeFileSync(`${path}/schema.json`, JSON.stringify(schema));
          console.log(`The result saved into ${path}/schema.json`);
          console.log("=========================================================================================");
        } catch (e) {
          console.log(e);
        }
      } else {
        console.log(`could not retrieve from ${jsonFilePath}`);
      }
    });
}
