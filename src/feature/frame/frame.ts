import { Command } from "commander";
import fs from "fs";
import path from "path";
import { frameRepository } from "./data/FrameRepositoryImpl";
import { generateFrame } from "./domain/useCase/dslToJson";
import { mapFrameModelToDSL } from "./domain/useCase/jsonToDsl";

export function frame(program: Command) {
  return program.command("frame").description("Manage frame");
}

export function genFrame(program: Command) {
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
          const output = await generateFrame(json);
          console.log(JSON.stringify(output, null, 2));
        } catch (e) {
          console.log(e);
        }
      } else {
        console.log(`could not retrieve from ${jsonFilePath}`);
      }
    });
}

export function pushFrame(program: Command) {
  return program
    .command("push")
    .description("Push a frame")
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
          const output = await generateFrame(json);
          if (output.data) {
            const result = await frameRepository.pushFrame(
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

export function pullFrame(program: Command) {
  return program
    .command("pull")
    .description("Pull a frame")
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
          if (json) {
            const result = await frameRepository.pullFrame(apiKey, json.route);
            if (result.onSuccess) {
              const frame = mapFrameModelToDSL(result.onSuccess);
              fs.writeFileSync(directory, JSON.stringify(frame));
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
