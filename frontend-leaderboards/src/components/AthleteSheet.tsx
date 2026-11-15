import { useEffect, useId, useRef } from "react";
import { categoryById, events, genderLabel, weightChoices } from "../data";
import { formatResult, ordinal } from "../format";
import { standings } from "../scoring";
import type { Athlete } from "../types";

type AthleteSheetProps = {
  athlete: Athlete;
  athletes: Athlete[];
  onClose: () => void;
  onShowClass: () => void;
  onShowOverall: () => void;
  showClassLink?: boolean;
};

export function AthleteSheet({
  athlete,
  athletes,
  onClose,
  onShowClass,
  onShowOverall,
  showClassLink = true,
}: AthleteSheetProps) {
  const titleId = useId();
  const closeRef = useRef<HTMLButtonElement>(null);
  const category = categoryById(athlete.categoryId);
  const classRows = standings(athletes.filter((item) => item.categoryId === athlete.categoryId));
  const fieldRows = standings(athletes.filter((item) => item.gender === athlete.gender));
  const classRow = classRows.find((row) => row.athlete.id === athlete.id);
  const fieldRow = fieldRows.find((row) => row.athlete.id === athlete.id);

  useEffect(() => {
    closeRef.current?.focus();
    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") onClose();
    }
    window.addEventListener("keydown", onKey);
    return () => {
      document.body.style.overflow = previous;
      window.removeEventListener("keydown", onKey);
    };
  }, [onClose]);

  return (
    <div className="fixed inset-0 z-40">
      <button type="button" aria-label="Close athlete details" className="absolute inset-0 bg-black/70" onClick={onClose} />
      <aside
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        className="absolute inset-y-0 right-0 z-10 flex w-full max-w-md flex-col overflow-y-auto border-l border-line bg-panel px-5 py-5"
      >
        <div className="flex items-start justify-between gap-3">
          <p className="font-display text-sm uppercase tracking-[0.18em] text-faint">
            {genderLabel(athlete.gender)}
          </p>
          <button
            ref={closeRef}
            type="button"
            onClick={onClose}
            className="rounded-md border border-line px-2 py-1 text-sm text-muted hover:text-chalk"
          >
            Close
          </button>
        </div>

        <h2 id={titleId} className="mt-3 font-display text-4xl uppercase leading-none tracking-wide">
          {athlete.name}
        </h2>
        <p className="mt-2 text-muted">{athlete.country}</p>

        {classRow && fieldRow && (
          <div className="mt-5 grid grid-cols-2 gap-3">
            <div className="rounded-md border border-line bg-ink px-3 py-3">
              <p className="font-mono text-3xl text-brass">
                {classRow.tied ? `Tied ${ordinal(classRow.place)}` : ordinal(classRow.place)}
              </p>
              <p className="mt-1 text-sm text-muted">{category.label}</p>
              <p className="text-sm text-muted">
                {classRows.length} {classRows.length === 1 ? "athlete" : "athletes"}
              </p>
            </div>
            <div className="rounded-md border border-line bg-ink px-3 py-3">
              <p className="font-mono text-3xl text-chalk">
                {fieldRow.tied ? `Tied ${ordinal(fieldRow.place)}` : ordinal(fieldRow.place)}
              </p>
              <p className="mt-1 text-sm text-muted">{genderLabel(athlete.gender)}</p>
              <p className="text-sm text-muted">
                {fieldRows.length} {fieldRows.length === 1 ? "athlete" : "athletes"}
              </p>
            </div>
          </div>
        )}

        <div className="mt-4 flex flex-col items-start gap-2">
          {showClassLink && (
            <button type="button" onClick={onShowClass} className="text-left text-sm text-brass hover:underline">
              {category.short} leaderboard
            </button>
          )}
          <button type="button" onClick={onShowOverall} className="text-left text-sm text-brass hover:underline">
            Overall {athlete.gender === "men" ? "men's" : "women's"} leaderboard
          </button>
        </div>

        <ul className="mt-6">
          {events.map((event) => {
            const result = athlete.results[event.id];
            const score = classRow?.events[event.id];
            return (
              <li key={event.id} className="border-t border-line py-3">
                <div className="flex items-baseline justify-between gap-4">
                  <div>
                    <p className="font-medium text-chalk">{event.name}</p>
                    <p className="text-xs text-muted">{weightChoices(event, athlete.gender)}</p>
                  </div>
                  <div className="text-right">
                    <p className="whitespace-nowrap font-mono text-sm tabular-nums">{formatResult(event, result)}</p>
                    <p className="font-mono text-[11px] text-brass">
                      {score?.place ? `${ordinal(score.place)} in class` : "—"}
                    </p>
                  </div>
                </div>
              </li>
            );
          })}
        </ul>
      </aside>
    </div>
  );
}
