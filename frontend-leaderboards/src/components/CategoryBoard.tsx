import { useEffect, useId, useMemo, useRef, useState } from "react";
import { categories, events, genderLabel, weightLine } from "../data";
import { formatResult, ordinal, pointsLabel } from "../format";
import { inRegion, regionById, regions, type RegionFilter } from "../regions";
import { standings } from "../scoring";
import type { Athlete, AthleteStanding, EventDef, EventId } from "../types";
import { ScoreTable } from "./ScoreTable";

type CategoryBoardProps = {
  athletes: Athlete[];
  categoryId: string;
  region: RegionFilter;
  selectedId: string | null;
  onSelect: (athleteId: string) => void;
  onCategory: (categoryId: string) => void;
  onRegion: (region: RegionFilter) => void;
};

function placed(rows: AthleteStanding[], eventId: EventId, place: number) {
  return rows.filter((row) => row.events[eventId].place === place);
}

export function CategoryBoard({
  athletes,
  categoryId,
  region,
  selectedId,
  onSelect,
  onCategory,
  onRegion,
}: CategoryBoardProps) {
  const [eventId, setEventId] = useState<EventId | null>(null);
  const regionLabel = region === "all" ? "All regions" : regionById(region).label;
  const boards = useMemo(
    () =>
      categories.map((category) => {
        const group = athletes.filter(
          (athlete) => athlete.categoryId === category.id && inRegion(athlete, region),
        );
        return { category, group, rows: standings(group) };
      }),
    [athletes, region],
  );

  const selected = boards.find((board) => board.category.id === categoryId) ?? boards[0];
  const rows = selected.rows;
  const leaders = rows.slice(0, 3);
  const openEvent = events.find((event) => event.id === eventId) ?? null;

  return (
    <section className="mt-8">
      <h2 className="font-display text-4xl uppercase leading-none tracking-wide">Weight classes</h2>

      <div className="mt-4 grid gap-3 sm:grid-cols-2">
        <label className="block">
          <span className="mb-1 block text-sm text-muted">Class</span>
          <select
            value={selected.category.id}
            onChange={(event) => onCategory(event.target.value)}
            className="w-full rounded-md border border-line bg-panel px-3 py-3 font-display text-lg uppercase tracking-wide text-chalk"
          >
            {(["men", "women"] as const).map((gender) => (
              <optgroup key={gender} label={genderLabel(gender)}>
                {boards
                  .filter((board) => board.category.gender === gender)
                  .map((board) => (
                    <option key={board.category.id} value={board.category.id}>
                      {board.category.short}
                    </option>
                  ))}
              </optgroup>
            ))}
          </select>
        </label>
        <label className="block">
          <span className="mb-1 block text-sm text-muted">Region</span>
          <select
            value={region}
            onChange={(event) => onRegion(event.target.value as RegionFilter)}
            className="w-full rounded-md border border-line bg-panel px-3 py-3 font-display text-lg uppercase tracking-wide text-chalk"
          >
            <option value="all">All regions</option>
            {regions.map((item) => (
              <option key={item.id} value={item.id}>
                {item.label}
              </option>
            ))}
          </select>
        </label>
      </div>

      <div className="mt-8">
        <p className="font-display text-sm uppercase tracking-[0.18em] text-faint">
          {genderLabel(selected.category.gender)}
        </p>
        <h3 className="mt-1 font-display text-4xl uppercase leading-none tracking-wide">
          {selected.category.label}
        </h3>
        <p className="mt-2 text-sm text-muted">
          {selected.group.length} {selected.group.length === 1 ? "athlete" : "athletes"}
          {region !== "all" ? ` · ${regionLabel}` : ""}
        </p>

        {selected.group.length === 0 ? (
          <p className="mt-5 rounded-lg border border-dashed border-line px-4 py-8 text-center text-muted">
            No athletes from {regionLabel} in this class.
          </p>
        ) : (
          <>
        <ol className="mt-5 flex flex-col gap-2">
          {leaders.map((row) => (
            <li key={row.athlete.id}>
              <button
                type="button"
                onClick={() => onSelect(row.athlete.id)}
                className={`flex w-full items-center justify-between gap-3 rounded-md border px-3 py-3 text-left sm:px-4 ${
                  row.place === 1 ? "border-brass bg-brass/5" : "border-line bg-panel"
                }`}
              >
                <span>
                  <span className="font-mono text-xs text-brass sm:text-sm">
                    {row.tied ? `Tied ${ordinal(row.place)}` : ordinal(row.place)}
                  </span>
                  <span className="mt-1 block font-display text-2xl uppercase leading-none tracking-wide">
                    {row.athlete.name}
                  </span>
                </span>
                <span className="shrink-0 font-mono text-sm text-chalk">{pointsLabel(row.total)}</span>
              </button>
            </li>
          ))}
        </ol>

        <h4 className="mt-8 font-display text-sm uppercase tracking-[0.18em] text-faint">Events</h4>
        <ul className="mt-2 grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
          {events.map((event) => (
            <li key={event.id}>
              <button
                type="button"
                onClick={() => setEventId(event.id)}
                className="w-full rounded-md border border-line bg-panel px-3 py-2 text-left hover:border-chalk/30"
              >
                <p className="font-display text-lg uppercase leading-none tracking-wide">{event.short}</p>
                <p className="mt-1 text-xs text-muted">Men {weightLine(event, "men")}</p>
                <p className="text-xs text-muted">Women {weightLine(event, "women")}</p>
                <p className="mt-1 text-[11px] text-faint">{event.note}</p>
                <p className="mt-2 text-sm text-brass">View places</p>
              </button>
            </li>
          ))}
        </ul>

        <div className="mt-6">
          <ScoreTable
            sections={[
              {
                key: selected.category.id,
                title: "",
                caption: "",
                categoryId: null,
                rows,
              },
            ]}
            pointsLabel="Points"
            showClass={false}
            showEventPoints
            selectedId={selectedId}
            onSelect={onSelect}
          />
        </div>
          </>
        )}
      </div>

      {openEvent && selected.group.length > 0 && (
        <EventPlaces
          event={openEvent}
          rows={rows}
          onClose={() => setEventId(null)}
          onSelect={(athleteId) => {
            setEventId(null);
            onSelect(athleteId);
          }}
        />
      )}
    </section>
  );
}

function EventPlaces({
  event,
  rows,
  onClose,
  onSelect,
}: {
  event: EventDef;
  rows: AthleteStanding[];
  onClose: () => void;
  onSelect: (athleteId: string) => void;
}) {
  const titleId = useId();
  const closeRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    closeRef.current?.focus();
    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    function onKey(keyEvent: KeyboardEvent) {
      if (keyEvent.key === "Escape") onClose();
    }
    window.addEventListener("keydown", onKey);
    return () => {
      document.body.style.overflow = previous;
      window.removeEventListener("keydown", onKey);
    };
  }, [onClose]);

  const places = [1, 2, 3, 4, 5, 6]
    .filter((place) => placed(rows, event.id, place).length > 0)
    .slice(0, 3);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button type="button" aria-label="Close event places" className="absolute inset-0 bg-black/70" onClick={onClose} />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        className="relative z-10 max-h-[calc(100svh-2rem)] w-full overflow-y-auto rounded-lg border border-line bg-panel px-5 py-5 sm:max-w-md"
      >
        <div className="flex items-start justify-between gap-3">
          <div>
            <p className="font-display text-sm uppercase tracking-[0.18em] text-faint">Event</p>
            <h3 id={titleId} className="mt-1 font-display text-3xl uppercase leading-none tracking-wide">
              {event.name}
            </h3>
          </div>
          <button
            ref={closeRef}
            type="button"
            onClick={onClose}
            className="rounded-md border border-line px-2 py-1 text-sm text-muted hover:text-chalk"
          >
            Close
          </button>
        </div>
        <p className="mt-3 text-xs text-muted">Men {weightLine(event, "men")}</p>
        <p className="text-xs text-muted">Women {weightLine(event, "women")}</p>
        <p className="mt-1 text-[11px] text-faint">{event.note}</p>

        {places.length === 0 ? (
          <p className="mt-5 text-sm text-muted">No result yet</p>
        ) : (
          <ol className="mt-5 flex flex-col gap-4">
            {places.map((place) => {
              const group = placed(rows, event.id, place);
              return (
                <li key={place}>
                  <p className="font-mono text-sm text-brass">
                    {group.length > 1 ? `Tied ${ordinal(place)}` : ordinal(place)}
                  </p>
                  <ul className="mt-1">
                    {group.map((row) => (
                      <li key={row.athlete.id}>
                        <button
                          type="button"
                          onClick={() => onSelect(row.athlete.id)}
                          className="flex w-full items-baseline justify-between gap-3 py-1 text-left"
                        >
                          <span className="text-chalk">{row.athlete.name}</span>
                          <span className="shrink-0 font-mono text-sm text-brass">
                            {formatResult(event, row.athlete.results[event.id])}
                          </span>
                        </button>
                      </li>
                    ))}
                  </ul>
                </li>
              );
            })}
          </ol>
        )}
      </div>
    </div>
  );
}
