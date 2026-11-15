import { useEffect, useId, useMemo, useRef, useState, type FormEvent } from "react";
import {
  categories,
  categoryById,
  competitions,
  countries,
  events,
  genderLabel,
  promoters,
  weightChoices,
  type NewAthlete,
} from "../data";
import { formatResult, joinNames } from "../format";
import { eventScores, isRecorded, needsResult } from "../scoring";
import { eventIds, type Athlete, type Category, type EventDef, type EventId, type EventResult, type Gender, type ResultUpdate } from "../types";
import { Segmented } from "./Segmented";

type EntryMode = "event" | "card";
type ResultList = "unentered" | "entered";

type PromoterDeskProps = {
  athletes: Athlete[];
  onSave: (athleteId: string, updates: Partial<Record<EventId, ResultUpdate>>, promoterId: string) => void;
  onAddAthlete: (athlete: NewAthlete) => void;
  onWithdraw: (athleteId: string, withdrawn: boolean) => void;
};

type Notice = {
  text: string;
};

type LiftDraft = {
  weight: string;
  value: string;
};

type Parsed = { ok: true; value: ResultUpdate } | { ok: false; message: string };

function blankLift(): LiftDraft {
  return { weight: "", value: "" };
}

function emptyCard(): Record<EventId, LiftDraft> {
  return Object.fromEntries(eventIds.map((id) => [id, blankLift()])) as Record<EventId, LiftDraft>;
}

function draftFromResult(result: EventResult): LiftDraft {
  return {
    weight: result.weight === null ? "" : String(result.weight),
    value: result.value === null ? "" : String(result.value),
  };
}

function cardFromAthlete(athlete: Athlete): Record<EventId, LiftDraft> {
  const draft = emptyCard();
  for (const event of events) {
    draft[event.id] = draftFromResult(athlete.results[event.id]);
  }
  return draft;
}

function sameLift(result: EventResult, value: ResultUpdate): boolean {
  return result.weight === value.weight && result.value === value.value;
}

function valueLabel(event: EventDef): string {
  if (event.measure === "reps") return "the reps";
  if (event.measure === "time") return "the time in seconds";
  return "the distance in metres";
}

function parseLift(draft: LiftDraft, event: EventDef, gender: Gender): Parsed {
  const weightRaw = draft.weight.trim();
  const valueRaw = draft.value.trim();
  if (weightRaw === "" && valueRaw === "") {
    return { ok: true, value: { weight: null, value: null } };
  }
  if (weightRaw === "" || valueRaw === "") {
    return { ok: false, message: `Choose a weight and enter ${valueLabel(event)}.` };
  }
  const weight = Number(weightRaw);
  const amount = Number(valueRaw);
  if (!event.weights[gender].includes(weight)) {
    return { ok: false, message: "Choose one of the listed weights." };
  }
  if (event.measure === "reps") {
    if (!Number.isInteger(amount) || amount < 1 || amount > 50) {
      return { ok: false, message: "Enter reps from 1 to 50." };
    }
    return { ok: true, value: { weight, value: amount } };
  }
  if (!Number.isFinite(amount) || amount <= 0) {
    return {
      ok: false,
      message: event.measure === "time" ? "Enter a time in seconds." : "Enter a distance in metres.",
    };
  }
  const rounded = Math.round(amount * 10) / 10;
  if (event.measure === "time" && rounded > 600) {
    return { ok: false, message: "Enter a time up to 600 seconds." };
  }
  if (event.measure === "distance" && rounded > 500) {
    return { ok: false, message: "Enter a distance up to 500 metres." };
  }
  return { ok: true, value: { weight, value: rounded } };
}

export function PromoterDesk({ athletes, onSave, onAddAthlete, onWithdraw }: PromoterDeskProps) {
  const formId = useId();
  const formRef = useRef<HTMLFormElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const [mode, setMode] = useState<EntryMode>("event");
  const promoterId = promoters[0].id;
  const [competitionId, setCompetitionId] = useState<string | null>(null);
  const [categoryId, setCategoryId] = useState(categories[0].id);
  const [hiddenClassIds, setHiddenClassIds] = useState<string[]>([]);
  const [pickingClasses, setPickingClasses] = useState(false);
  const [eventId, setEventId] = useState<EventId>("log");
  const [athleteId, setAthleteId] = useState<string | null>(null);
  const [query, setQuery] = useState("");
  const [resultList, setResultList] = useState<ResultList>("unentered");
  const [draft, setDraft] = useState<LiftDraft>(blankLift);
  const [cardDraft, setCardDraft] = useState<Record<EventId, LiftDraft>>(emptyCard);
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<Notice | null>(null);
  const [adding, setAdding] = useState(false);
  const [confirmingDropOut, setConfirmingDropOut] = useState(false);
  const [pointsSort, setPointsSort] = useState<"asc" | "desc">("desc");
  const [pinTop, setPinTop] = useState(112);

  const event = events.find((item) => item.id === eventId) ?? events[0];
  const athlete = athletes.find((item) => item.id === athleteId) ?? null;
  const competition = competitions.find((item) => item.id === competitionId) ?? null;
  const competitionAthletes = useMemo(
    () => (competitionId ? athletes.filter((item) => item.competitionId === competitionId) : []),
    [athletes, competitionId],
  );
  const category = categories.find((item) => item.id === categoryId) ?? null;
  const shownClasses = useMemo(
    () => classesWithAthletes(competitionAthletes, hiddenClassIds),
    [competitionAthletes, hiddenClassIds],
  );

  const visibleAthletes = useMemo(() => {
    if (!category) return [];
    const needle = query.trim().toLowerCase();
    return competitionAthletes
      .filter((item) => {
        if (item.categoryId !== category.id) return false;
        if (!needle) return true;
        const weightClass = categoryById(item.categoryId);
        return `${item.name} ${item.country} ${weightClass.short} ${weightClass.label}`
          .toLowerCase()
          .includes(needle);
      })
      .slice()
      .sort((a, b) => a.name.localeCompare(b.name));
  }, [category?.id, competitionAthletes, query]);

  const classAthletes = useMemo(
    () => (category ? competitionAthletes.filter((item) => item.categoryId === category.id) : []),
    [category, competitionAthletes],
  );
  const eventComplete =
    mode === "event" &&
    classAthletes.length > 0 &&
    classAthletes.every((item) => isRecorded(item.results[eventId]));
  const pointsById = useMemo(() => {
    if (!eventComplete) return null;
    const points = new Map<string, number>();
    for (const [id, score] of eventScores(classAthletes, eventId)) points.set(id, score.points);
    return points;
  }, [classAthletes, eventComplete, eventId]);

  const unentered = visibleAthletes.filter((item) =>
    mode === "event" ? needsResult(item, eventId) : events.some((entry) => needsResult(item, entry.id)),
  );
  const entered = visibleAthletes.filter((item) =>
    mode === "event" ? !needsResult(item, eventId) : events.every((entry) => !needsResult(item, entry.id)),
  );
  const rankedEntered = pointsById
    ? entered.slice().sort((a, b) => {
        const diff = (pointsById.get(a.id) ?? 0) - (pointsById.get(b.id) ?? 0);
        const ordered = pointsSort === "asc" ? diff : -diff;
        return ordered || a.name.localeCompare(b.name);
      })
    : entered;
  const listed = resultList === "unentered" ? unentered : rankedEntered;
  const listedIndex = athleteId ? listed.findIndex((item) => item.id === athleteId) : -1;

  function fillFromAthlete(nextId: string, nextMode: EntryMode, nextEventId: EventId) {
    const next = athletes.find((item) => item.id === nextId);
    if (!next) return;
    if (nextMode === "event") {
      setDraft(draftFromResult(next.results[nextEventId]));
    } else {
      setCardDraft(cardFromAthlete(next));
    }
  }

  function selectAthlete(nextId: string, keepNotice = false) {
    setAthleteId(nextId);
    setError(null);
    if (!keepNotice) setNotice(null);
    fillFromAthlete(nextId, mode, eventId);
  }

  function moveInList(index: number) {
    const next = listed[index];
    if (next) selectAthlete(next.id);
  }

  useEffect(() => {
    const header = document.querySelector("header");
    if (!header) return;
    const update = () => setPinTop(Math.ceil(header.getBoundingClientRect().height));
    update();
    const observer = new ResizeObserver(update);
    observer.observe(header);
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!athleteId || !window.matchMedia("(max-width: 1023px)").matches) return;
    formRef.current?.scrollIntoView({ behavior: "smooth", block: "start" });
  }, [athleteId]);

  useEffect(() => {
    const list = listRef.current;
    if (!list || !athleteId) return;
    const selected = list.querySelector<HTMLElement>(`[data-athlete-id="${CSS.escape(athleteId)}"]`);
    if (!selected) return;
    const listRect = list.getBoundingClientRect();
    const itemRect = selected.getBoundingClientRect();
    if (itemRect.top < listRect.top) list.scrollTop -= listRect.top - itemRect.top;
    else if (itemRect.bottom > listRect.bottom) list.scrollTop += itemRect.bottom - listRect.bottom;
  }, [athleteId, resultList, query, categoryId, mode, eventId, pointsSort]);

  function changeMode(nextMode: EntryMode) {
    setMode(nextMode);
    setError(null);
    setNotice(null);
    if (athleteId) fillFromAthlete(athleteId, nextMode, eventId);
  }

  function changeEvent(nextEventId: EventId) {
    setEventId(nextEventId);
    setError(null);
    setNotice(null);
    if (athleteId) fillFromAthlete(athleteId, "event", nextEventId);
  }

  function openCompetition(nextId: string) {
    setCompetitionId(nextId);
    setAthleteId(null);
    setQuery("");
    setError(null);
    setNotice(null);
  }

  function leaveCompetition() {
    setCompetitionId(null);
    setAthleteId(null);
    setQuery("");
    setError(null);
    setNotice(null);
  }

  function toggleClass(id: string) {
    setHiddenClassIds((current) =>
      current.includes(id) ? current.filter((item) => item !== id) : [...current, id],
    );
  }

  function showAllClasses() {
    setHiddenClassIds([]);
    const occupied = classesWithAthletes(competitionAthletes, []);
    if (occupied.some((item) => item.id === categoryId)) return;
    setCategoryId(occupied[0]?.id ?? "");
    setAthleteId(null);
    setQuery("");
    setError(null);
    setNotice(null);
  }

  function hideAllClasses() {
    setHiddenClassIds(classesWithAthletes(competitionAthletes, []).map((item) => item.id));
  }

  function changeCategory(nextId: string) {
    setCategoryId(nextId);
    setAthleteId(null);
    setQuery("");
    setError(null);
    setNotice(null);
  }

  function saveEvent(submitEvent: FormEvent) {
    submitEvent.preventDefault();
    if (!athlete) {
      setError("Choose an athlete first.");
      return;
    }
    const parsed = parseLift(draft, event, athlete.gender);
    if (!parsed.ok) {
      setError(parsed.message);
      return;
    }
    if (parsed.value.weight === null || parsed.value.value === null) {
      setError(`Choose a weight and enter ${valueLabel(event)}, or use Clear to remove a result.`);
      return;
    }
    if (sameLift(athlete.results[event.id], parsed.value)) {
      setError("That result is already on the board.");
      return;
    }
    const wasEmpty = !isRecorded(athlete.results[event.id]);
    onSave(athlete.id, { [event.id]: parsed.value }, promoterId);
    setError(null);
    setNotice({
      text: `Saved ${athlete.name} · ${event.name} · ${formatResult(event, { ...parsed.value, enteredBy: promoterId })}`,
    });
    setDraft({ weight: String(parsed.value.weight), value: String(parsed.value.value) });
    if (wasEmpty) {
      const next = unentered.find((item) => item.id !== athlete.id);
      if (next) selectAthlete(next.id, true);
    }
  }

  function clearEvent() {
    if (!athlete) return;
    if (!isRecorded(athlete.results[event.id])) {
      setError("There is no result to clear.");
      return;
    }
    onSave(athlete.id, { [event.id]: { weight: null, value: null, skipped: false } }, promoterId);
    setDraft(blankLift());
    setError(null);
    setNotice({
      text: `Cleared ${athlete.name} · ${event.name}`,
    });
  }

  function markSkip(eventId: EventId, skipped: boolean) {
    if (!athlete) return;
    const target = events.find((item) => item.id === eventId) ?? event;
    const order = events.findIndex((item) => item.id === target.id) + 1;
    onSave(athlete.id, { [eventId]: { weight: null, value: null, skipped } }, promoterId);
    if (eventId === event.id) setDraft(blankLift());
    setError(null);
    setNotice({
      text: skipped
        ? `${athlete.name} is not doing event ${order} - ${target.name}`
        : `${athlete.name} is doing event ${order} - ${target.name}`,
    });
    if (skipped) {
      const next = unentered.find((item) => item.id !== athlete.id);
      if (next) selectAthlete(next.id, true);
    }
  }

  function dropOut(withdrawn: boolean) {
    if (!athlete) return;
    const stillNeeded = needsResult(athlete, event.id);
    onWithdraw(athlete.id, withdrawn);
    setError(null);
    setNotice({
      text: withdrawn ? `${athlete.name} dropped out` : `${athlete.name} is back in the competition`,
    });
    if (withdrawn && stillNeeded) {
      const next = unentered.find((item) => item.id !== athlete.id);
      if (next) selectAthlete(next.id, true);
    }
  }

  function requestDropOut() {
    if (!athlete) return;
    if (athlete.withdrawn) {
      dropOut(false);
      return;
    }
    setConfirmingDropOut(true);
  }

  function confirmDropOut() {
    if (!athlete) return;
    const start = events.findIndex((item) => item.id === event.id);
    const following = events.slice(start < 0 ? 0 : start);
    const updates: Partial<Record<EventId, ResultUpdate>> = {};
    for (const item of following) {
      updates[item.id] = { weight: null, value: null, skipped: true };
    }
    const stillNeeded = needsResult(athlete, event.id);
    onSave(athlete.id, updates, promoterId);
    setDraft(blankLift());
    setCardDraft((current) => {
      const next = { ...current };
      for (const item of following) next[item.id] = blankLift();
      return next;
    });
    setError(null);
    setNotice({
      text: `${athlete.name} dropped out. Results for the following events are did not participate.`,
    });
    setConfirmingDropOut(false);
    if (stillNeeded) {
      const next = unentered.find((item) => item.id !== athlete.id);
      if (next) selectAthlete(next.id, true);
    }
  }

  function saveCard(submitEvent: FormEvent) {
    submitEvent.preventDefault();
    if (!athlete) {
      setError("Choose an athlete first.");
      return;
    }
    const updates: Partial<Record<EventId, ResultUpdate>> = {};
    const problems: string[] = [];
    for (const item of events) {
      const parsed = parseLift(cardDraft[item.id], item, athlete.gender);
      if (!parsed.ok) {
        problems.push(`${item.name}: ${parsed.message}`);
        continue;
      }
      if (!sameLift(athlete.results[item.id], parsed.value)) {
        updates[item.id] = parsed.value;
      }
    }
    if (problems.length > 0) {
      setError(problems.join(" "));
      return;
    }
    if (Object.keys(updates).length === 0) {
      setError("Nothing has changed on this card.");
      return;
    }
    onSave(athlete.id, updates, promoterId);
    setError(null);
    const changed = events.filter((item) => Object.prototype.hasOwnProperty.call(updates, item.id));
    setNotice({
      text: `Saved ${athlete.name} · ${joinNames(changed.map((item) => item.short))}`,
    });
  }

  if (!competition) {
    return (
      <section>
        <h2 className="font-display text-3xl uppercase leading-none tracking-wide">Your competitions</h2>
        <ul className="mt-4 grid gap-3 sm:grid-cols-2">
          {competitions.map((item) => {
            const field = athletes.filter((athlete) => athlete.competitionId === item.id);
            return (
              <li key={item.id}>
                <button
                  type="button"
                  onClick={() => openCompetition(item.id)}
                  className="flex h-full w-full items-center gap-4 rounded-lg border border-line bg-panel p-4 text-left hover:border-brass"
                >
                  <img src="/logo.png" alt="" className="h-24 w-24 shrink-0 object-contain" />
                  <span className="flex min-w-0 flex-col">
                    <span className="font-display text-3xl uppercase leading-none tracking-wide text-chalk">
                      {item.name}
                    </span>
                    <span className="mt-2 text-sm text-muted">{item.country}</span>
                    <span className="mt-4 font-mono text-sm text-brass">{field.length} athletes</span>
                  </span>
                </button>
              </li>
            );
          })}
        </ul>
      </section>
    );
  }

  return (
    <section>
      <div className="grid items-center gap-3 sm:grid-cols-[auto_minmax(0,1fr)_auto]">
        <div>
          <button
            type="button"
            onClick={leaveCompetition}
            className="inline-flex items-center gap-1.5 rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
          >
            <svg viewBox="0 0 20 20" className="size-4" aria-hidden="true">
              <path fill="currentColor" d="M12.8 4.2 7 10l5.8 5.8-1.2 1.2L4.6 10l7-7z" />
            </svg>
            Back
          </button>
        </div>
        <h2 className="text-center font-display text-3xl uppercase leading-none tracking-wide">{competition.name}</h2>
        <div className="flex flex-wrap items-center gap-3 sm:justify-end">
          <button
            type="button"
            onClick={() => setAdding(true)}
            className="rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
          >
            Add athletes
          </button>
          <Segmented
            label="Entry mode"
            value={mode}
            onChange={changeMode}
            options={[
              { value: "event", label: "One event" },
              { value: "card", label: "All events" },
            ]}
          />
        </div>
      </div>

      <div className="mt-5">
        <button
          type="button"
          aria-expanded={pickingClasses}
          onClick={() => setPickingClasses((open) => !open)}
          className={`rounded-md border px-3 py-2 text-sm ${
            pickingClasses ? "border-brass bg-raised text-chalk" : "border-line text-muted hover:text-chalk"
          }`}
        >
          Filter classes
        </button>
        {pickingClasses && (
          <ClassFilter
            athletes={competitionAthletes}
            hidden={hiddenClassIds}
            onToggle={toggleClass}
            onShowAll={showAllClasses}
            onHideAll={hideAllClasses}
          />
        )}
      </div>

      <WeightClasses
        athletes={competitionAthletes}
        classes={shownClasses}
        mode={mode}
        eventId={event.id}
        value={category?.id ?? ""}
        onChange={changeCategory}
      />

      {mode === "event" && <EventPicker value={event.id} onChange={changeEvent} />}

      <div className="mt-6 grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px] lg:items-start">
        <div className={athlete ? "max-lg:order-2" : undefined}>
          {category && (
            <h3 className="font-display text-2xl uppercase leading-none tracking-wide">
              {genderLabel(category.gender)} {category.short}
            </h3>
          )}

          <div className={`${category ? "mt-3" : ""} flex flex-wrap items-center gap-3`}>
            <Segmented
              label="Result list"
              value={resultList}
              onChange={setResultList}
              options={[
                { value: "unentered", label: `No results ${unentered.length}` },
                { value: "entered", label: `Results ${entered.length}` },
              ]}
            />
            {eventComplete && resultList === "entered" && (
              <div className="ml-auto">
                <Segmented
                  label="Points order"
                  value={pointsSort}
                  onChange={setPointsSort}
                  options={[
                    { value: "desc", label: "High to low" },
                    { value: "asc", label: "Low to high" },
                  ]}
                />
              </div>
            )}
          </div>

          <label className="mt-3 block">
            <span className="sr-only">Search athletes</span>
            <input
              value={query}
              onChange={(change) => setQuery(change.target.value)}
              placeholder={category ? `Search ${category.short}` : "Search athletes"}
              className="w-full rounded-md border border-line bg-ink px-3 py-2 text-sm text-chalk outline-none placeholder:text-faint focus:border-brass"
            />
          </label>

          <div ref={listRef} className="mt-3 max-h-[32rem] space-y-4 overflow-y-auto pr-1">
            <AthleteGroup
              athletes={listed}
              event={event}
              mode={mode}
              pointsById={resultList === "entered" ? pointsById : null}
              selectedId={athleteId}
              onSelect={selectAthlete}
            />
            {listed.length === 0 && (
              <p className="rounded-md border border-dashed border-line px-3 py-6 text-center text-sm text-muted">
                {category ? "No athletes match." : "Choose a weight class to show."}
              </p>
            )}
          </div>
        </div>

        <form
          ref={formRef}
          onSubmit={mode === "event" ? saveEvent : saveCard}
          style={{ scrollMarginTop: pinTop, top: pinTop }}
          className={`rounded-lg border border-line bg-panel p-4 lg:sticky ${
            athlete ? "max-lg:sticky max-lg:z-20 max-lg:order-1 max-lg:shadow-[0_12px_28px_rgba(0,0,0,0.45)]" : ""
          }`}
        >
          {notice && (
            <div className="mb-4 rounded-md border border-brass/40 bg-brass/10 px-3 py-2 text-sm" aria-live="polite">
              <p>{notice.text}</p>
            </div>
          )}
          {!athlete ? (
            <p className="text-sm text-muted">Select an athlete to enter a result.</p>
          ) : (
            <>
              <div className="mb-3 grid grid-cols-4 gap-2" role="group" aria-label="Move through athletes">
                <button
                  type="button"
                  onClick={() => moveInList(0)}
                  disabled={listed.length === 0 || listedIndex === 0}
                  className="rounded-md border border-line px-2 py-1.5 text-sm text-muted hover:text-chalk disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:text-muted"
                >
                  First
                </button>
                <button
                  type="button"
                  onClick={() => moveInList(listedIndex - 1)}
                  disabled={listedIndex <= 0}
                  className="rounded-md border border-line px-2 py-1.5 text-sm text-muted hover:text-chalk disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:text-muted"
                >
                  Previous
                </button>
                <button
                  type="button"
                  onClick={() => moveInList(listedIndex + 1)}
                  disabled={listedIndex < 0 || listedIndex >= listed.length - 1}
                  className="rounded-md border border-line px-2 py-1.5 text-sm text-muted hover:text-chalk disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:text-muted"
                >
                  Next
                </button>
                <button
                  type="button"
                  onClick={() => moveInList(listed.length - 1)}
                  disabled={listed.length === 0 || listedIndex === listed.length - 1}
                  className="rounded-md border border-line px-2 py-1.5 text-sm text-muted hover:text-chalk disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:text-muted"
                >
                  Last
                </button>
              </div>
              <p className="font-display text-2xl uppercase leading-none tracking-wide">{athlete.name}</p>
              <p className="mt-1 text-sm text-muted">
                {genderLabel(athlete.gender)} · {categoryById(athlete.categoryId).label}
              </p>

              {mode === "event" ? (
                <LiftFields
                  idPrefix={formId}
                  event={event}
                  gender={athlete.gender}
                  draft={draft}
                  onChange={setDraft}
                />
              ) : (
                <div className="mt-4 space-y-4">
                  {events.map((item) => (
                    <LiftFields
                      key={item.id}
                      idPrefix={`${formId}-${item.id}`}
                      event={item}
                      gender={athlete.gender}
                      draft={cardDraft[item.id]}
                      skipped={!athlete.withdrawn && athlete.results[item.id].skipped}
                      onToggleSkip={
                        athlete.withdrawn
                          ? undefined
                          : () => markSkip(item.id, !athlete.results[item.id].skipped)
                      }
                      onChange={(next) => setCardDraft((current) => ({ ...current, [item.id]: next }))}
                    />
                  ))}
                </div>
              )}

              {mode === "event" && <CurrentLine athlete={athlete} event={event} />}

              {error && (
                <p role="alert" className="mt-3 text-sm text-ember">
                  {error}
                </p>
              )}

              <div className="mt-4 flex flex-wrap gap-2">
                <button type="submit" className="rounded-md border-b-[3px] border-black/30 bg-brass px-3 py-2 text-sm font-bold tracking-wide text-white uppercase hover:bg-ember">
                  {mode === "event"
                    ? isRecorded(athlete.results[event.id])
                      ? "Update result"
                      : "Enter result"
                    : events.some((item) => isRecorded(athlete.results[item.id]))
                      ? "Update results"
                      : "Enter results"}
                </button>
                {mode === "event" && (
                  <button
                    type="button"
                    onClick={clearEvent}
                    className="rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
                  >
                    Clear
                  </button>
                )}
              </div>
              <div className="mt-2 flex flex-wrap gap-2">
                {mode === "event" && !athlete.withdrawn && (
                  <button
                    type="button"
                    onClick={() => markSkip(event.id, !athlete.results[event.id].skipped)}
                    className="rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
                  >
                    {athlete.results[event.id].skipped ? "Doing this event" : "Did not participate"}
                  </button>
                )}
                <button
                  type="button"
                  onClick={requestDropOut}
                  className="rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
                >
                  {athlete.withdrawn ? "Back in" : "Dropped out"}
                </button>
              </div>
            </>
          )}
        </form>
      </div>
      {confirmingDropOut && athlete && (
        <DropOutWarning
          athleteName={athlete.name}
          events={events.slice(Math.max(0, events.findIndex((item) => item.id === event.id)))}
          onCancel={() => setConfirmingDropOut(false)}
          onConfirm={confirmDropOut}
        />
      )}
      {adding && (
        <AddAthleteDialog
          competitionId={competition.id}
          categoryId={category?.id || shownClasses[0]?.id || categories[0].id}
          onClose={() => setAdding(false)}
          onAdd={(draft) => {
            onAddAthlete(draft);
            setHiddenClassIds((current) => current.filter((id) => id !== draft.categoryId));
            setCategoryId(draft.categoryId);
            setResultList("unentered");
            setAthleteId(null);
            setAdding(false);
          }}
        />
      )}
    </section>
  );
}

function DropOutWarning({
  athleteName,
  events: affected,
  onCancel,
  onConfirm,
}: {
  athleteName: string;
  events: EventDef[];
  onCancel: () => void;
  onConfirm: () => void;
}) {
  const titleId = useId();

  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") onCancel();
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [onCancel]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4" onClick={onCancel}>
      <div
        role="alertdialog"
        aria-modal="true"
        aria-labelledby={titleId}
        onClick={(event) => event.stopPropagation()}
        className="w-full max-w-md rounded-lg border border-line bg-panel p-4"
      >
        <h3 id={titleId} className="font-display text-2xl uppercase tracking-wide">
          Dropped out
        </h3>
        <p role="alert" className="mt-3 text-sm text-chalk">
          This sets {athleteName}&apos;s results for the following events to did not participate.
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-5 text-sm text-muted">
          {affected.map((item) => (
            <li key={item.id}>{item.name}</li>
          ))}
        </ul>
        <div className="mt-4 flex flex-wrap gap-2">
          <button
            type="button"
            onClick={onConfirm}
            className="rounded-md border-b-[3px] border-black/30 bg-ember px-3 py-2 text-sm font-bold tracking-wide text-white uppercase hover:bg-brass"
          >
            Dropped out
          </button>
          <button
            type="button"
            autoFocus
            onClick={onCancel}
            className="rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}

function AddAthleteDialog({
  competitionId,
  categoryId,
  onClose,
  onAdd,
}: {
  competitionId: string;
  categoryId: string;
  onClose: () => void;
  onAdd: (athlete: NewAthlete) => void;
}) {
  const titleId = useId();
  const competition = competitions.find((item) => item.id === competitionId);
  const [name, setName] = useState("");
  const [country, setCountry] = useState(competition?.country ?? countries[0]);
  const [classId, setClassId] = useState(categoryId);
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [onClose]);

  function submit(event: FormEvent) {
    event.preventDefault();
    const trimmedName = name.trim();
    if (!trimmedName) {
      setFormError("Enter a name.");
      return;
    }
    if (!country) {
      setFormError("Choose a country.");
      return;
    }
    onAdd({
      name: trimmedName,
      country,
      categoryId: classId,
      competitionId,
    });
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4" onClick={onClose}>
      <form
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        onClick={(event) => event.stopPropagation()}
        onSubmit={submit}
        className="w-full max-w-md rounded-lg border border-line bg-panel p-4"
      >
        <h3 id={titleId} className="font-display text-2xl uppercase tracking-wide">
          Add athlete
        </h3>
        <label className="mt-4 block text-sm">
          <span className="mb-1 block text-xs text-muted">Name</span>
          <input
            autoFocus
            value={name}
            onChange={(change) => setName(change.target.value)}
            className="w-full rounded-md border border-line bg-ink px-3 py-2 text-chalk outline-none focus:border-brass"
          />
        </label>
        <label className="mt-3 block text-sm">
          <span className="mb-1 block text-xs text-muted">Country</span>
          <select
            value={country}
            onChange={(change) => setCountry(change.target.value)}
            className="w-full rounded-md border border-line bg-ink px-3 py-2 text-chalk outline-none focus:border-brass"
          >
            {countries.map((item) => (
              <option key={item} value={item}>
                {item}
              </option>
            ))}
          </select>
        </label>
        <label className="mt-3 block text-sm">
          <span className="mb-1 block text-xs text-muted">Weight class</span>
          <select
            value={classId}
            onChange={(change) => setClassId(change.target.value)}
            className="w-full rounded-md border border-line bg-ink px-3 py-2 text-chalk outline-none focus:border-brass"
          >
            <optgroup label="Women">
              {categories
                .filter((item) => item.gender === "women")
                .map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.label}
                  </option>
                ))}
            </optgroup>
            <optgroup label="Men">
              {categories
                .filter((item) => item.gender === "men")
                .map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.label}
                  </option>
                ))}
            </optgroup>
          </select>
        </label>
        {formError && (
          <p role="alert" className="mt-3 text-sm text-ember">
            {formError}
          </p>
        )}
        <div className="mt-4 flex flex-wrap gap-2">
          <button
            type="submit"
            className="rounded-md border-b-[3px] border-black/30 bg-brass px-3 py-2 text-sm font-bold tracking-wide text-white uppercase hover:bg-ember"
          >
            Add athlete
          </button>
          <button
            type="button"
            onClick={onClose}
            className="rounded-md border border-line px-3 py-2 text-sm text-muted hover:text-chalk"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

function fieldCopy(event: EventDef): { label: string; placeholder: string; mode: "numeric" | "decimal" } {
  if (event.measure === "reps") return { label: "Reps", placeholder: "1", mode: "numeric" };
  if (event.measure === "time") return { label: "Time (s)", placeholder: "30", mode: "decimal" };
  return { label: "Distance (m)", placeholder: "40", mode: "decimal" };
}

function LiftFields({
  idPrefix,
  event,
  gender,
  draft,
  skipped = false,
  onToggleSkip,
  onChange,
}: {
  idPrefix: string;
  event: EventDef;
  gender: Gender;
  draft: LiftDraft;
  skipped?: boolean;
  onToggleSkip?: () => void;
  onChange: (draft: LiftDraft) => void;
}) {
  const field = fieldCopy(event);
  return (
    <fieldset className="mt-4">
      <legend className="text-sm text-chalk">{event.name}</legend>
      <p className="mb-2 text-xs text-faint">
        {weightChoices(event, gender)}
        {event.note ? ` · ${event.note}` : ""}
      </p>
      <div className="grid grid-cols-2 gap-2">
        <label className="block text-sm">
          <span className="mb-1 block text-xs text-muted">Weight</span>
          <select
            id={`${idPrefix}-weight`}
            value={draft.weight}
            onChange={(change) => onChange({ ...draft, weight: change.target.value })}
            className="w-full rounded-md border border-line bg-ink px-2 py-2 text-chalk outline-none focus:border-brass"
          >
            <option value="">Choose</option>
            {event.weights[gender].map((weight) => (
              <option key={weight} value={weight}>
                {weight}
              </option>
            ))}
          </select>
        </label>
        <label className="block text-sm">
          <span className="mb-1 block text-xs text-muted">{field.label}</span>
          <input
            id={`${idPrefix}-value`}
            value={draft.value}
            onChange={(change) => onChange({ ...draft, value: change.target.value })}
            inputMode={field.mode}
            placeholder={field.placeholder}
            className="w-full rounded-md border border-line bg-ink px-3 py-2 font-mono text-chalk outline-none focus:border-brass"
          />
        </label>
      </div>
      {onToggleSkip && (
        <button
          type="button"
          onClick={onToggleSkip}
          className="mt-2 rounded-md border border-line px-2 py-1 text-xs text-muted hover:text-chalk"
        >
          {skipped ? "Doing this event" : "Did not participate"}
        </button>
      )}
    </fieldset>
  );
}

function EventLines({ event }: { event: EventDef }) {
  const order = events.findIndex((item) => item.id === event.id) + 1;
  return (
    <span className="block min-w-0">
      <span className="flex items-baseline gap-2 font-display text-lg uppercase leading-none tracking-wide">
        <span className="font-mono text-sm text-faint">{order}</span>
        <span>{event.short}</span>
      </span>
      <span className="mt-2 block text-sm normal-case tracking-normal text-muted">
        Women {event.weights.women.join(" · ")}
      </span>
      <span className="mt-0.5 block text-sm normal-case tracking-normal text-muted">
        Men {event.weights.men.join(" · ")}
      </span>
    </span>
  );
}

function EventPicker({ value, onChange }: { value: EventId; onChange: (eventId: EventId) => void }) {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const selected = events.find((item) => item.id === value) ?? events[0];

  useEffect(() => {
    if (!open) return;
    function onPointer(event: PointerEvent) {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false);
    }
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") setOpen(false);
    }
    document.addEventListener("pointerdown", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("pointerdown", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <div ref={rootRef} className="relative mt-5">
      <button
        type="button"
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-label="Event"
        onClick={() => setOpen((current) => !current)}
        className="flex w-full items-start justify-between gap-3 rounded-md border border-line bg-panel px-3 py-3 text-left text-chalk"
      >
        <EventLines event={selected} />
        <svg viewBox="0 0 20 20" className="mt-1 size-4 shrink-0 text-muted" aria-hidden="true">
          <path fill="currentColor" d="M5.2 7.5 10 12.3l4.8-4.8 1.2 1.2L10 14.7 4 8.7z" />
        </svg>
      </button>
      {open && (
        <ul
          role="listbox"
          aria-label="Event"
          className="absolute z-20 mt-1 max-h-96 w-full overflow-y-auto rounded-md border border-line bg-panel shadow-[0_8px_24px_rgba(0,0,0,0.28)]"
        >
          {events
            .filter((item) => item.id !== value)
            .map((item) => (
              <li key={item.id}>
                <button
                  type="button"
                  role="option"
                  aria-selected={false}
                  onClick={() => {
                    onChange(item.id);
                    setOpen(false);
                  }}
                  className="flex w-full px-3 py-3 text-left hover:bg-raised"
                >
                  <EventLines event={item} />
                </button>
              </li>
            ))}
        </ul>
      )}
    </div>
  );
}

function CurrentLine({ athlete, event }: { athlete: Athlete; event: EventDef }) {
  const result = athlete.results[event.id];
  if (isRecorded(result)) {
    return <p className="mt-2 text-xs text-faint">Current result {formatResult(event, result)}</p>;
  }
  if (athlete.withdrawn) return <p className="mt-2 text-xs text-faint">Dropped out.</p>;
  if (result.skipped) return <p className="mt-2 text-xs text-faint">Did not participate.</p>;
  return <p className="mt-2 text-xs text-faint">No result yet.</p>;
}

function classesWithAthletes(athletes: Athlete[], hiddenIds: string[]): Category[] {
  const hidden = new Set(hiddenIds);
  return (["women", "men"] as const).flatMap((gender) =>
    categories.filter(
      (item) =>
        item.gender === gender &&
        !hidden.has(item.id) &&
        athletes.some((athlete) => athlete.categoryId === item.id),
    ),
  );
}

function WeightClasses({
  athletes,
  classes,
  mode,
  eventId,
  value,
  onChange,
}: {
  athletes: Athlete[];
  classes: Category[];
  mode: EntryMode;
  eventId: EventId;
  value: string;
  onChange: (categoryId: string) => void;
}) {
  return (
    <div className="mt-3 space-y-3">
      {(["women", "men"] as const).map((gender) => {
        const group = classes.filter((item) => item.gender === gender);
        if (group.length === 0) return null;
        return (
          <div key={gender}>
            <p className="mb-2 font-display text-sm uppercase tracking-[0.16em] text-faint">{genderLabel(gender)}</p>
            <div className="flex flex-wrap gap-2">
              {group.map((item) => {
                const field = athletes.filter((athlete) => athlete.categoryId === item.id);
                const missing = field.filter((athlete) => needsResult(athlete, eventId)).length;
                const active = item.id === value;
                return (
                  <button
                    key={item.id}
                    type="button"
                    aria-pressed={active}
                    onClick={() => onChange(item.id)}
                    className={`rounded-md border px-2.5 py-1.5 text-sm ${
                      active ? "border-brass bg-raised text-chalk" : "border-line text-muted hover:text-chalk"
                    }`}
                  >
                    <span className="font-medium">{item.short}</span>
                    <span className="ml-2 font-mono text-xs text-faint">{field.length}</span>
                    {mode === "event" && missing > 0 && (
                      <span className="ml-2 font-mono text-xs text-brass">{missing} left</span>
                    )}
                  </button>
                );
              })}
            </div>
          </div>
        );
      })}
    </div>
  );
}

function ClassFilter({
  athletes,
  hidden,
  onToggle,
  onShowAll,
  onHideAll,
}: {
  athletes: Athlete[];
  hidden: string[];
  onToggle: (categoryId: string) => void;
  onShowAll: () => void;
  onHideAll: () => void;
}) {
  const occupied = classesWithAthletes(athletes, []);
  const allShown = occupied.every((item) => !hidden.includes(item.id));
  const noneShown = occupied.length > 0 && occupied.every((item) => hidden.includes(item.id));
  return (
    <div className="mt-3 rounded-lg border border-line bg-panel p-4">
      <div className="mb-3 flex flex-wrap gap-2">
        <button
          type="button"
          onClick={onShowAll}
          disabled={allShown}
          className="rounded-md border border-line px-3 py-1.5 text-sm text-muted hover:text-chalk disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:text-muted"
        >
          Select all
        </button>
        <button
          type="button"
          onClick={onHideAll}
          disabled={noneShown}
          className="rounded-md border border-line px-3 py-1.5 text-sm text-muted hover:text-chalk disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:text-muted"
        >
          Deselect all
        </button>
      </div>
      {(["women", "men"] as const).map((gender) => {
        const group = occupied.filter((item) => item.gender === gender);
        if (group.length === 0) return null;
        return (
          <fieldset key={gender} className="mt-3 first:mt-0">
            <legend className="mb-2 font-display text-sm uppercase tracking-[0.16em] text-faint">
              {genderLabel(gender)}
            </legend>
            <div className="flex flex-wrap gap-x-4 gap-y-2">
              {group.map((item) => (
                <label key={item.id} className="flex items-center gap-2 text-sm text-chalk">
                  <input
                    type="checkbox"
                    checked={!hidden.includes(item.id)}
                    onChange={() => onToggle(item.id)}
                    className="accent-brass"
                  />
                  {item.short}
                </label>
              ))}
            </div>
          </fieldset>
        );
      })}
    </div>
  );
}

function formatPoints(points: number): string {
  const rounded = Math.round(points * 10) / 10;
  const text = Number.isInteger(rounded) ? String(rounded) : rounded.toFixed(1);
  return `${text} ${rounded === 1 ? "pt" : "pts"}`;
}

function AthleteGroup({
  athletes,
  event,
  mode,
  pointsById,
  selectedId,
  onSelect,
}: {
  athletes: Athlete[];
  event: EventDef;
  mode: EntryMode;
  pointsById: Map<string, number> | null;
  selectedId: string | null;
  onSelect: (athleteId: string, keepNotice?: boolean) => void;
}) {
  if (athletes.length === 0) return null;
  return (
    <ul className="space-y-2">
        {athletes.map((athlete) => {
          const selected = athlete.id === selectedId;
          const filled = events.filter((item) => isRecorded(athlete.results[item.id])).length;
          const eventResult = athlete.results[event.id];
          const eventText = athlete.withdrawn && !isRecorded(eventResult)
            ? "Dropped out"
            : eventResult.skipped
              ? "Did not participate"
              : formatResult(event, eventResult);
          const cardText = athlete.withdrawn ? "Dropped out" : `${filled}/5`;
          const points = pointsById?.get(athlete.id);
          const resultText = mode === "event" ? eventText : cardText;
          return (
            <li key={athlete.id}>
              <button
                type="button"
                data-athlete-id={athlete.id}
                onClick={() => onSelect(athlete.id)}
                className={`flex w-full items-center justify-between gap-3 rounded-md border px-3 py-2 text-left ${
                  selected ? "border-brass bg-raised" : "border-line hover:border-chalk/30"
                }`}
              >
                <span>
                  <span className="block font-medium text-chalk">{athlete.name}</span>
                  <span className="block text-xs text-muted">
                    {athlete.country}
                  </span>
                </span>
                <span className="shrink-0 whitespace-nowrap font-mono text-sm text-brass">
                  {points === undefined ? resultText : `${formatPoints(points)} · ${resultText}`}
                </span>
              </button>
            </li>
          );
        })}
    </ul>
  );
}
