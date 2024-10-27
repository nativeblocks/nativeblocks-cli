import { GraphQLClient } from "graphql-request";
import { regionRepository } from "../../feature/region/data/RegionRepositoryImpl";
import fs from "fs";
import os from "os";
import path from "path";
import { execSync } from "child_process";

let instance: GraphQLClient | null = null;
export const getGraphqlClient = () => {
  if (instance == null) {
    instance = new GraphQLClient("", {
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
    });
    const endpoint = regionRepository.get();
    if (endpoint.onSuccess) {
      instance.setEndpoint(endpoint.onSuccess);
    }
  }
  return instance;
};

function getAuthToken(): string {
  const userHomeDir: string = os.homedir();
  const credentialPath: string = path.join(userHomeDir, ".nativeblocks/cli/credential.json");
  const directory = path.dirname(credentialPath);
  if (fs.existsSync(directory)) {
    try {
      const data: string = fs.readFileSync(credentialPath, "utf-8");
      const json = JSON.parse(data);
      return json.accessToken;
    } catch (e) {
      return "";
    }
  } else {
    return "";
  }
}

interface ErrorResult {
  errorClass: string;
  errorMessage: string;
}

function isOnline(customUrl: string): boolean {
  try {
    execSync(`ping -c 1 ${customUrl}`, { stdio: "ignore" });
    return true;
  } catch (error) {
    return false;
  }
}

export function handleNetworkError(error: any): ErrorResult {
  const endpoint = regionRepository.get();

  if (!endpoint.onSuccess)
    return {
      errorClass: "CONNECTION",
      errorMessage: "Please first set the region endpoint url",
    };

  if (!isOnline("8.8.8.8"))
    return {
      errorClass: "CONNECTION",
      errorMessage: "Please check the connection and try again",
    };

  try {
    const response = error.response;
    if (response && response.errors && response.errors.length > 0) {
      const firstError = response.errors[0];
      const errorClass = firstError.extensions ? firstError.extensions.classification : "";
      const errorMessage = firstError.message || "Unknown error";
      return { errorClass, errorMessage };
    } else {
      console.error("Unexpected error format:", error.message);
    }
  } catch (_) {}
  return { errorClass: "UnknownClass", errorMessage: "Unknown error" };
}
