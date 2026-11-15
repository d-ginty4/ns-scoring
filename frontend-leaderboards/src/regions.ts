import type { Athlete } from "./types";

export const regions = [
  {
    id: "north-america",
    label: "North America",
    countries: ["United States", "Canada", "Mexico"],
  },
  {
    id: "uk-ireland",
    label: "UK & Ireland",
    countries: ["England", "Scotland", "Wales", "Northern Ireland", "Ireland"],
  },
  {
    id: "mainland-europe",
    label: "Mainland Europe",
    countries: [
      "Austria",
      "Belgium",
      "Bulgaria",
      "Croatia",
      "Czechia",
      "Denmark",
      "Estonia",
      "Finland",
      "France",
      "Germany",
      "Greece",
      "Hungary",
      "Iceland",
      "Italy",
      "Latvia",
      "Lithuania",
      "Netherlands",
      "Norway",
      "Poland",
      "Portugal",
      "Romania",
      "Spain",
      "Sweden",
      "Switzerland",
      "Ukraine",
    ],
  },
  {
    id: "oceania",
    label: "Oceania",
    countries: ["Australia", "New Zealand"],
  },
] as const;

export type RegionId = (typeof regions)[number]["id"];

export type RegionFilter = RegionId | "all";

const countryToRegion = new Map<string, RegionId>(
  regions.flatMap((region) => region.countries.map((country) => [country, region.id] as const)),
);

export function regionById(id: RegionId) {
  const region = regions.find((item) => item.id === id);
  if (!region) return regions[0];
  return region;
}

export function regionForCountry(country: string): RegionId | null {
  return countryToRegion.get(country) ?? null;
}

export function inRegion(athlete: Athlete, region: RegionFilter): boolean {
  if (region === "all") return true;
  return regionForCountry(athlete.country) === region;
}
