import type { EventDef, EventResult, EventScore } from "./oldTypes";
import { isRecorded } from "./scoring";

function trimNumber(value: number): string {
  const rounded = Math.round(value * 10) / 10;
  return Number.isInteger(rounded) ? String(rounded) : rounded.toFixed(1);
}

export function resultLines(event: EventDef, result: EventResult): { performance: string; weight: string } | null {
  if (!isRecorded(result)) return null;
  const amount = trimNumber(result.value);
  const weight = `${result.weight}kg`;
  if (event.measure === "reps") return { performance: `${amount} reps`, weight };
  if (event.measure === "time") return { performance: `${amount}s`, weight };
  return { performance: `${amount} m`, weight };
}

export function formatResult(event: EventDef, result: EventResult): string {
  const lines = resultLines(event, result);
  if (!lines) return "—";
  return `${lines.performance} - ${lines.weight}`;
}

export function formatPoints(points: number): string {
  const rounded = Math.round(points * 100) / 100;
  if (Number.isInteger(rounded)) return String(rounded);
  return rounded.toFixed(2).replace(/0$/, "");
}

export function pointsLabel(points: number): string {
  const rendered = formatPoints(points);
  return `${rendered} ${rendered === "1" ? "pt" : "pts"}`;
}

export function ordinal(place: number): string {
  const mod100 = place % 100;
  if (mod100 >= 11 && mod100 <= 13) return `${place}th`;
  switch (place % 10) {
    case 1:
      return `${place}st`;
    case 2:
      return `${place}nd`;
    case 3:
      return `${place}rd`;
    default:
      return `${place}th`;
  }
}

export function joinNames(names: string[]): string {
  if (names.length === 0) return "";
  if (names.length === 1) return names[0];
  if (names.length === 2) return `${names[0]} and ${names[1]}`;
  return `${names.slice(0, -1).join(", ")}, and ${names[names.length - 1]}`;
}

export function resultTitle(event: EventDef, result: EventResult, score: EventScore): string {
  if (!isRecorded(result) || score.place === null) {
    return `${event.name}: not entered`;
  }
  return `${event.name}: ${formatResult(event, result)}. ${ordinal(score.place)} in this ranking, ${formatPoints(score.points)} points.`;
}

export function columnHint(event: EventDef): string {
  if (event.measure === "reps") return "reps";
  if (event.measure === "time") return "time";
  return "distance";
}
