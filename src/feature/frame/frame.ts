import { input } from "@inquirer/prompts";
import { Command } from "commander";
import fs from "fs";
import path from "path";
import { createDefaultDir } from "../../infrastructure/utility/FileUitl";
import { frameRepository } from "./data/FrameRepositoryImpl";
import { findData, findKeyTypes, findProperties } from "./extractKeyType";
import { processFrame } from "./frameValidator";
import { generateBaseSchema } from "./generateSchema";

export function frame(program: Command) {
  return program.command("frame").description("Manage frame");
}

export function generateFrame(program: Command) {
  return program
    .command("gen")
    .description("Generate a frame")
    .option("-d, --directory", "Frame working directory")
    .argument("<directory>", "Frame working directory")
    .action(async (directory) => {
      const jsonFilePath: string = path.join(directory);
      if (fs.existsSync(jsonFilePath)) {
        try {
          const data: string = fs.readFileSync(jsonFilePath, "utf-8");
          const json = JSON.parse(data);
          const output = await processFrame(json);
          console.log(JSON.stringify(output, null, 2));
        } catch (e) {
          console.log(e);
        }
      } else {
        console.log(`could not retrieve from ${jsonFilePath}`);
      }
    });
}

export function syncFrame(program: Command) {
  return program
    .command("sync")
    .description("Sync a frame")
    .option("-d, --directory", "Frame working directory")
    .option("-apiKey, --apiKey", "Project api key")
    .argument("<directory>", "Frame working directory")
    .argument("<apiKey>", "Project api key")
    .action(async (directory, apiKey) => {
      const jsonFilePath: string = path.join(directory);
      if (fs.existsSync(jsonFilePath)) {
        try {
          const data: string = fs.readFileSync(jsonFilePath, "utf-8");
          const json = JSON.parse(data);
          const output = await processFrame(json);
          if (output.data) {
            const result = await frameRepository.syncFrame(
              apiKey,
              output.data.frameProduction.route,
              JSON.stringify(output.data.frameProduction)
            );
            if (result.onSuccess) {
              console.log(`Frame synced`);
            } else {
              console.log(`Sync faild: ${result.onError}`);
            }
          }
        } catch (e) {
          console.log(e);
        }
      } else {
        console.log(`could not retrieve from ${jsonFilePath}`);
      }
    });
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
          const blockKeyTypes = findKeyTypes(jsonFilePath + "/block");
          const blockPropertyKeys = findProperties(jsonFilePath + "/block");
          const blockDataKeys = findData(jsonFilePath + "/block");

          const actionKeyTypes = findKeyTypes(jsonFilePath + "/action");
          const actionPropertyKeys = findProperties(jsonFilePath + "/action");
          const actionDataKeys = findData(jsonFilePath + "/action");

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
