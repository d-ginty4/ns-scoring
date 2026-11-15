import { useMemo, useState } from "react";
import { categoryById, genderLabel, promoterNames } from "../data";
import { standings } from "../scoring";
import type { Athlete, Gender } from "../types";
import { ScoreTable } from "./ScoreTable";
import { Segmented } from "./Segmented";

type AllResultsProps = {
  athletes: Athlete[];
  gender: Gender;
  onGender: (gender: Gender) => void;
  selectedId: string | null;
  onSelect: (athleteId: string) => void;
};

function matches(athlete: Athlete, query: string): boolean {
  const category = categoryById(athlete.categoryId);
  const haystack = [
    athlete.name,
    athlete.country,
    category.label,
    category.short,
    genderLabel(athlete.gender),
    ...promoterNames(athlete),
  ]
    .join(" ")
    .toLowerCase();
  return haystack.includes(query.trim().toLowerCase());
}

export function AllResults({ athletes, gender, onGender, selectedId, onSelect }: AllResultsProps) {
  const [query, setQuery] = useState("");

  const rows = useMemo(() => {
    return standings(athletes.filter((athlete) => athlete.gender === gender)).filter((row) =>
      matches(row.athlete, query),
    );
  }, [athletes, gender, query]);

  return (
    <section className="mt-8">
      <h2 className="font-display text-4xl uppercase leading-none tracking-wide">Overall leaderboard</h2>

      <div className="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <Segmented
          label="Gender"
          value={gender}
          onChange={onGender}
          options={[
            { value: "women", label: "Women" },
            { value: "men", label: "Men" },
          ]}
        />
        <p className="text-sm text-muted" aria-live="polite">
          {rows.length} {rows.length === 1 ? "athlete" : "athletes"}
        </p>
      </div>

      <label className="mt-3 block">
        <span className="sr-only">Search athletes</span>
        <input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search athlete, country, class, or promoter"
          className="w-full rounded-md border border-line bg-panel px-3 py-2 text-sm text-chalk outline-none placeholder:text-faint focus:border-brass"
        />
      </label>

      <div className="mt-4">
        {rows.length === 0 ? (
          <p className="rounded-lg border border-dashed border-line px-4 py-8 text-center text-muted">
            No athletes match that search.
          </p>
        ) : (
          <ScoreTable
            sections={[{ key: "overall", title: "", caption: "", categoryId: null, rows }]}
            pointsLabel="Points"
            showEventPoints
            selectedId={selectedId}
            onSelect={onSelect}
          />
        )}
      </div>

      <details className="mt-4 text-sm text-muted">
        <summary className="text-chalk">How points work</summary>
        <p className="mt-2 max-w-3xl">
          Men and women each have their own leaderboard, because they lift different weights. On each
          event a heavier weight beats any result at a lighter weight. At the same weight, more reps
          or more metres places higher, and a faster time places higher on the loading race and the
          farmers. First place earns one point for every athlete on that board, and the rest count
          down. A tie splits those points. Event wins break a tie on total points, then the
          athlete&apos;s name. The weight-class boards use the same rule inside each class.
        </p>
      </details>
    </section>
  );
}
