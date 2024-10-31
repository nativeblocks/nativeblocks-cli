import { Command } from "commander";
import { regionRepository } from "./data/RegionRepositoryImpl";

export function region(program: Command) {
  return program
    .command("region")
    .description("Set or get the active region")
}

export function setRegion(program: Command) {
  return program
    .command("set")
    .description("Set the region")
    .option("-u, --url", "Api region url")
    .argument("<url>", "Region url")
    .action((url) => {
      const result = regionRepository.set(url);
      if (result.onSuccess) {
        console.log(result.onSuccess);
      } else {
        console.error(result.onError);
      }
    });
}

export function getRegion(program: Command) {
  return program
    .command("get")
    .description("Get the active region")
    .action(() => {
      const result = regionRepository.get();
      if (result.onSuccess) {
        console.log(result.onSuccess)
      } else {
        console.error(result.onError)
      }
    });
}