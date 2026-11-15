import { useState } from "react";
import { PromoterDesk } from "./components/PromoterDesk";
import { createAthlete, events, seedAthletes, type NewAthlete } from "./data";
import type { Athlete, EventId, ResultUpdate } from "./types";

export default function App() {
  const [athletes, setAthletes] = useState<Athlete[]>(() => structuredClone(seedAthletes));
  const [dark, setDark] = useState(() => document.documentElement.classList.contains("dark"));

  function toggleTheme() {
    const next = !document.documentElement.classList.contains("dark");
    document.documentElement.classList.toggle("dark", next);
    localStorage.setItem("theme", next ? "dark" : "light");
    setDark(next);
  }

  function setWithdrawn(athleteId: string, withdrawn: boolean) {
    setAthletes((current) =>
      current.map((athlete) => (athlete.id === athleteId ? { ...athlete, withdrawn } : athlete)),
    );
  }

  function addAthlete(draft: NewAthlete) {
    setAthletes((current) => [...current, createAthlete(draft, current)]);
  }

  function saveResults(
    athleteId: string,
    updates: Partial<Record<EventId, ResultUpdate>>,
    promoterId: string,
  ) {
    setAthletes((current) =>
      current.map((athlete) => {
        if (athlete.id !== athleteId) return athlete;
        const nextResults = { ...athlete.results };
        for (const event of events) {
          if (!Object.prototype.hasOwnProperty.call(updates, event.id)) continue;
          const update = updates[event.id];
          if (!update) continue;
          const previous = athlete.results[event.id];
          const skipped = update.skipped === true;
          const recorded = !skipped && update.weight !== null && update.value !== null;
          const nextResult = {
            weight: recorded ? update.weight : null,
            value: recorded ? update.value : null,
            enteredBy: recorded || skipped ? promoterId : null,
            skipped,
          };
          if (
            previous.weight === nextResult.weight &&
            previous.value === nextResult.value &&
            previous.skipped === nextResult.skipped
          ) {
            continue;
          }
          nextResults[event.id] = nextResult;
        }
        return { ...athlete, results: nextResults };
      }),
    );
  }

  return (
    <div className="min-h-svh">
      <header className="sticky top-0 z-30 border-b-4 border-brass bg-brand-black text-white shadow-[0_4px_20px_rgba(0,0,0,0.35)]">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-3 px-4 py-4 sm:px-6">
          <div className="flex min-w-0 items-center gap-3">
            <img
              src="/logo.png"
              alt="Natural Strongman"
              className="h-12 w-12 shrink-0 object-contain md:h-24 md:w-24"
            />
            <div className="min-w-0">
              <h1 className="font-display text-[1.7rem] uppercase leading-[0.95] tracking-wide text-white sm:text-[2rem]">
                Natural Strongman Open 2027
              </h1>
              <h2 className="mt-2 font-display text-sm font-semibold uppercase tracking-[0.18em] text-brass">
                Promoter desk
              </h2>
            </div>
          </div>
          <div className="flex shrink-0 items-center gap-3">
            <button
              type="button"
              role="switch"
              aria-checked={dark}
              aria-label="Dark mode"
              onClick={toggleTheme}
              className="relative grid h-9 w-[4.75rem] shrink-0 grid-cols-2 overflow-hidden rounded-full border border-white/25 bg-white/10 p-0.5"
            >
            <span
              aria-hidden="true"
              className={`absolute top-0.5 bottom-0.5 left-0.5 w-[calc(50%-2px)] rounded-full bg-brass transition-transform duration-300 ease-out motion-reduce:transition-none ${
                dark ? "translate-x-full" : "translate-x-0"
              }`}
            />
            <span className={`relative z-10 flex items-center justify-center ${dark ? "text-white/55" : "text-white"}`}>
              <SunIcon />
            </span>
            <span className={`relative z-10 flex items-center justify-center ${dark ? "text-white" : "text-white/55"}`}>
              <MoonIcon />
            </span>
            </button>
            <button type="button" className="rounded-md border border-white/25 px-3 py-2 text-sm font-medium text-white">
              Log out
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6 sm:px-6">
        <PromoterDesk
          athletes={athletes}
          onSave={saveResults}
          onAddAthlete={addAthlete}
          onWithdraw={setWithdrawn}
        />
      </main>
    </div>
  );
}

function SunIcon() {
  return (
    <svg viewBox="0 0 24 24" className="size-4" aria-hidden="true">
      <circle cx="12" cy="12" r="4" fill="none" stroke="currentColor" strokeWidth="2" />
      <path
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        d="M12 2.5v2.5M12 19v2.5M4.9 4.9l1.8 1.8M17.3 17.3l1.8 1.8M2.5 12H5M19 12h2.5M4.9 19.1l1.8-1.8M17.3 6.7l1.8-1.8"
      />
    </svg>
  );
}

function MoonIcon() {
  return (
    <svg viewBox="0 0 24 24" className="size-4" aria-hidden="true">
      <path
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinejoin="round"
        d="M20.5 14.5A8.2 8.2 0 0 1 9.5 3.2 7 7 0 1 0 20.5 14.5z"
      />
    </svg>
  );
}
