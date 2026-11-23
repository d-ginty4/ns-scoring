import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { Faker, en, en_CA, en_GB, en_IE, en_US } from "@faker-js/faker";
import { competitionsUrl, divisionsUrl, eventsUrl, regionsUrl } from "./constants.ts";

const athletesPerDivision = 10;

const localesByCountry: Record<string, typeof en> = {
  Ireland: en_IE,
  England: en_GB,
  Scotland: en_GB,
  Wales: en_GB,
  "United States": en_US,
  Canada: en_CA,
};

const fakers = new Map<string, Faker>();

const fakerFor = (country: string) => {
  const existing = fakers.get(country);
  if (existing) {
    return existing;
  }

  const created = new Faker({ locale: [localesByCountry[country] ?? en, en] });
  fakers.set(country, created);
  return created;
};

const uniqueName = (country: string, gender: string, used: Set<string>) => {
  const sex = gender === "female" ? "female" : "male";
  const faker = fakerFor(country);

  for (let attempt = 0; attempt < 20; attempt++) {
    const name = `${faker.person.firstName(sex)} ${faker.person.lastName(sex)}`;
    if (!used.has(name)) {
      used.add(name);
      return name;
    }
  }

  throw new Error(`could not generate a unique ${gender} name for ${country}`);
};

const dataDir = path.join(path.dirname(fileURLToPath(import.meta.url)), "data");

type Envelope<T> = { data: T };

const post = async <T>(url: string, record: unknown) => {
  const response = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(record),
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`${response.status} ${body}`);
  }

  return (await response.json()) as Envelope<T>;
};

const main = async () => {
  const regions = JSON.parse(await readFile(path.join(dataDir, "regions.json"), "utf8")) as { name: string }[];
  const regionIds = new Map<string, number>();
  for (const region of regions) {
    const created = await post<{ id: number; name: string }>(regionsUrl, region);
    regionIds.set(created.data.name, created.data.id);
  }

  const events = JSON.parse(await readFile(path.join(dataDir, "events.json"), "utf8")) as unknown[];
  for (const event of events) {
    await post(eventsUrl, event);
  }

  const divisions = JSON.parse(await readFile(path.join(dataDir, "divisions.json"), "utf8")) as {
    name: string;
    gender: string;
  }[];
  const divisionIds = new Map<string, number>();
  for (const division of divisions) {
    const created = await post<{ id: number }>(divisionsUrl, division);
    divisionIds.set(`${division.gender}:${division.name}`, created.data.id);
  }

  const competitions = JSON.parse(await readFile(path.join(dataDir, "competitions.json"), "utf8")) as {
    name: string;
    country: string;
    region: string;
    iconName?: string;
  }[];
  for (const competition of competitions) {
    const regionId = regionIds.get(competition.region);
    if (regionId === undefined) {
      throw new Error(`region ${competition.region} was not seeded`);
    }

    const created = await post<{ id: number }>(competitionsUrl, {
      name: competition.name,
      country: competition.country,
      regionId,
      ...(competition.iconName ? { iconName: competition.iconName } : {}),
    });

    const usedNames = new Set<string>();
    const athletes = divisions.flatMap((division) => {
      const divisionId = divisionIds.get(`${division.gender}:${division.name}`);
      if (divisionId === undefined) {
        throw new Error(`division ${division.gender} ${division.name} was not seeded`);
      }

      return Array.from({ length: athletesPerDivision }, () => ({
        name: uniqueName(competition.country, division.gender, usedNames),
        country: competition.country,
        divisionId,
      }));
    });

    await post(`${competitionsUrl}/${created.data.id}/athletes/bulk`, { athletes });
  }
};

main();
