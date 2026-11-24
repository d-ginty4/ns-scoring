export const eventIds = ["log", "deadlift", "loading", "farmers", "carry"] as const;

export type EventId = (typeof eventIds)[number];

export type Gender = "men" | "women";

export type Measure = "reps" | "time" | "distance";

export type EventDef = {
  id: EventId;
  name: string;
  short: string;
  weights: Record<Gender, number[]>;
  weightLabel: string;
  measure: Measure;
  note: string;
};

export type Category = {
  id: string;
  gender: Gender;
  label: string;
  short: string;
};

export type Promoter = {
  id: string;
  name: string;
  short: string;
  region: string;
};

export type Competition = {
  id: string;
  name: string;
  country: string;
};

export type EventResult = {
  weight: number | null;
  value: number | null;
  enteredBy: string | null;
  skipped: boolean;
};

export type ResultUpdate = {
  weight: number | null;
  value: number | null;
  skipped?: boolean;
};

export type Athlete = {
  id: string;
  name: string;
  country: string;
  gender: Gender;
  categoryId: string;
  competitionId: string;
  bodyweight: number;
  withdrawn: boolean;
  results: Record<EventId, EventResult>;
};

export type EventScore = {
  place: number | null;
  points: number;
};

export type AthleteStanding = {
  athlete: Athlete;
  events: Record<EventId, EventScore>;
  total: number;
  wins: number;
  place: number;
  tied: boolean;
};
