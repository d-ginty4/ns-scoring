import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { eventsUrl, regionsUrl } from "./constants.ts";

const dataDir = path.join(path.dirname(fileURLToPath(import.meta.url)), "data");

const post = async (url: string, record: unknown) => {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(record),
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`${response.status} ${body}`);
  }
};

const main = async () => {
  const regions = JSON.parse(await readFile(path.join(dataDir, "regions.json"), "utf8")) as unknown[];
  for (const region of regions) {
    await post(regionsUrl, region);
  }

  const events = JSON.parse(await readFile(path.join(dataDir, "events.json"), "utf8")) as unknown[];
  for (const event of events) {
    await post(eventsUrl, event);
  }
};

main();
