import { events } from "./data";
import type { Athlete, AthleteStanding, EventId, EventResult, EventScore, Gender } from "./oldTypes";

export function isRecorded(result: EventResult): result is EventResult & { weight: number; value: number } {
  return result.weight !== null && result.value !== null;
}

export function needsResult(athlete: Athlete, eventId: EventId): boolean {
  if (athlete.withdrawn || athlete.results[eventId].skipped) return false;
  return !isRecorded(athlete.results[eventId]);
}

function optionRank(eventId: EventId, gender: Gender, weight: number): number {
  const event = events.find((item) => item.id === eventId);
  return event?.weights[gender].indexOf(weight) ?? -1;
}

function scoreEvent(group: Athlete[], eventId: EventId): Map<string, EventScore> {
  const event = events.find((item) => item.id === eventId);
  const lowerIsBetter = event?.measure === "time";
  const ranked = group
    .filter((athlete) => isRecorded(athlete.results[eventId]))
    .map((athlete) => {
      const result = athlete.results[eventId];
      return {
        id: athlete.id,
        rank: optionRank(eventId, athlete.gender, result.weight as number),
        value: result.value as number,
      };
    })
    .sort((a, b) => {
      if (a.rank !== b.rank) return b.rank - a.rank;
      return lowerIsBetter ? a.value - b.value : b.value - a.value;
    });

  const maxPoints = group.length;
  const scores = new Map<string, EventScore>();
  let index = 0;

  while (index < ranked.length) {
    let end = index;
    while (
      end + 1 < ranked.length &&
      ranked[end + 1].rank === ranked[index].rank &&
      ranked[end + 1].value === ranked[index].value
    ) {
      end += 1;
    }
    let pointsSum = 0;
    for (let placeIndex = index; placeIndex <= end; placeIndex += 1) {
      pointsSum += maxPoints - (placeIndex + 1) + 1;
    }
    const points = pointsSum / (end - index + 1);
    for (let placeIndex = index; placeIndex <= end; placeIndex += 1) {
      scores.set(ranked[placeIndex].id, { place: index + 1, points });
    }
    index = end + 1;
  }

  for (const athlete of group) {
    if (!scores.has(athlete.id)) {
      scores.set(athlete.id, { place: null, points: 0 });
    }
  }

  return scores;
}

export function eventScores(group: Athlete[], eventId: EventId): Map<string, EventScore> {
  return scoreEvent(group, eventId);
}

export function standings(group: Athlete[]): AthleteStanding[] {
  const byEvent = Object.fromEntries(
    events.map((event) => [event.id, scoreEvent(group, event.id)]),
  ) as Record<EventId, Map<string, EventScore>>;

  const rows: AthleteStanding[] = group.map((athlete) => {
    const eventScores = {} as Record<EventId, EventScore>;
    let total = 0;
    let wins = 0;
    for (const event of events) {
      const score = byEvent[event.id].get(athlete.id) ?? { place: null, points: 0 };
      eventScores[event.id] = score;
      total += score.points;
      if (score.place === 1) wins += 1;
    }
    return { athlete, events: eventScores, total, wins, place: 0, tied: false };
  });

  rows.sort(
    (a, b) =>
      b.total - a.total || b.wins - a.wins || a.athlete.name.localeCompare(b.athlete.name),
  );

  for (let index = 0; index < rows.length; ) {
    let end = index;
    while (
      end + 1 < rows.length &&
      rows[end + 1].total === rows[index].total &&
      rows[end + 1].wins === rows[index].wins
    ) {
      end += 1;
    }
    const tied = end > index;
    for (let cursor = index; cursor <= end; cursor += 1) {
      rows[cursor] = { ...rows[cursor], place: index + 1, tied };
    }
    index = end + 1;
  }

  return rows;
}
