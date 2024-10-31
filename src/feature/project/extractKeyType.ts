import * as fs from "fs";
import * as path from "path";

export function findKeyTypes(dirPath: string): string[] {
  const keyTypes: string[] = [];

  function searchDirectory(currentPath: string) {
    const entries = fs.readdirSync(currentPath, { withFileTypes: true });
    entries.forEach((entry) => {
      const entryPath = path.join(currentPath, entry.name);
      if (entry.isDirectory()) {
        searchDirectory(entryPath);
      } else if (entry.isFile() && entry.name === "integration.json") {
        const fileContent = fs.readFileSync(entryPath, "utf-8");
        try {
          const jsonData = JSON.parse(fileContent);
          keyTypes.push(jsonData.keyType);
        } catch (error) {
          console.error(`Error parsing ${entryPath}:`, error);
        }
      }
    });
  }

  searchDirectory(dirPath);
  return keyTypes.filter((item) => !!item);
}

export function findData(dirPath: string): { key: string; value: string; type: string }[] {
  const data: { key: string; value: string; type: string }[] = [];

  function searchDirectory(currentPath: string) {
    const entries = fs.readdirSync(currentPath, { withFileTypes: true });
    entries.forEach((entry) => {
      const entryPath = path.join(currentPath, entry.name);
      if (entry.isDirectory()) {
        searchDirectory(entryPath);
      } else if (entry.isFile() && entry.name === "data.json") {
        const fileContent = fs.readFileSync(entryPath, "utf-8");
        try {
          const jsonData = JSON.parse(fileContent);
          jsonData.forEach((item: any) => {
            data.push({
              key: item.key,
              value: "",
              type: item.type,
            });
          });
        } catch (error) {
          console.error(`Error parsing ${entryPath}:`, error);
        }
      }
    });
  }

  searchDirectory(dirPath);
  return data.filter((item) => !!item.key);
}

export function findProperties(dirPath: string): { key: string; value: string; type: string }[] {
  const properties: { key: string; value: string; type: string }[] = [];

  function searchDirectory(currentPath: string) {
    const entries = fs.readdirSync(currentPath, { withFileTypes: true });
    entries.forEach((entry) => {
      const entryPath = path.join(currentPath, entry.name);
      if (entry.isDirectory()) {
        searchDirectory(entryPath);
      } else if (entry.isFile() && entry.name === "properties.json") {
        const fileContent = fs.readFileSync(entryPath, "utf-8");
        try {
          const jsonData = JSON.parse(fileContent);
          jsonData.forEach((item: any) => {
            properties.push({
              key: item.key,
              value: item.value,
              type: item.type,
            });
          });
        } catch (error) {
          console.error(`Error parsing ${entryPath}:`, error);
        }
      }
    });
  }

  searchDirectory(dirPath);
  return properties.filter((item) => !!item.key);
}
