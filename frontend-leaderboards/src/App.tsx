import { useEffect, useRef, useState } from "react";
import { AthleteSheet } from "./components/AthleteSheet";
import { AllResults } from "./components/AllResults";
import { CategoryBoard } from "./components/CategoryBoard";
import { RegionalBoard } from "./components/RegionalBoard";
import { seedAthletes } from "./data";
import { regions, type RegionFilter, type RegionId } from "./regions";
import type { Athlete, Gender } from "./types";

type ViewId = "board" | "regions" | "classes";

const views: { id: ViewId; label: string }[] = [
  { id: "board", label: "Overall leaderboard" },
  { id: "regions", label: "Regional leaderboards" },
  { id: "classes", label: "Weight classes" },
];

export default function App() {
  const [athletes, setAthletes] = useState<Athlete[]>(() => structuredClone(seedAthletes));
  const [view, setView] = useState<ViewId>("board");
  const [boardGender, setBoardGender] = useState<Gender>("men");
  const [regionId, setRegionId] = useState<RegionId>(regions[0].id);
  const [classRegion, setClassRegion] = useState<RegionFilter>("all");
  const [categoryId, setCategoryId] = useState("m90");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [dark, setDark] = useState(() => document.documentElement.classList.contains("dark"));
  const headerRef = useRef<HTMLElement>(null);

  useEffect(() => {
    const header = headerRef.current;
    if (!header) return;
    function syncHeaderHeight() {
      document.documentElement.style.setProperty("--site-header", `${header?.offsetHeight ?? 0}px`);
    }
    syncHeaderHeight();
    const observer = new ResizeObserver(syncHeaderHeight);
    observer.observe(header);
    return () => observer.disconnect();
  }, []);

  const selected = athletes.find((athlete) => athlete.id === selectedId) ?? null;

  function toggleTheme() {
    const next = !document.documentElement.classList.contains("dark");
    document.documentElement.classList.toggle("dark", next);
    localStorage.setItem("theme", next ? "dark" : "light");
    setDark(next);
  }

  function restore() {
    setAthletes(structuredClone(seedAthletes));
  }

  return (
    <div className="min-h-svh">
      <header
        ref={headerRef}
        className="sticky top-0 z-30 border-b-4 border-brass bg-brand-black text-white shadow-[0_4px_20px_rgba(0,0,0,0.35)]"
      >
        <div className="mx-auto flex max-w-6xl flex-col gap-4 px-4 py-4 sm:px-6">
          <div className="flex items-center justify-between gap-3">
            <div className="flex min-w-0 items-center gap-3">
              <img
                src="/logo.png"
                alt="Natural Strongman"
                className="h-12 w-12 shrink-0 object-contain md:h-24 md:w-24"
              />
              <h1 className="font-display text-[1.7rem] uppercase leading-[0.95] tracking-wide text-white sm:text-[2rem]">
                Natural Strongman Open 2027
              </h1>
            </div>
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
          </div>
          <nav aria-label="Views" className="flex flex-wrap gap-1">
            {views.map((item) => {
              const active = view === item.id;
              return (
                <button
                  key={item.id}
                  type="button"
                  aria-current={active ? "page" : undefined}
                  onClick={() => setView(item.id)}
                  className={`border-b-[3px] px-3 py-2 text-sm font-semibold tracking-wide uppercase ${
                    active
                      ? "border-brass text-white"
                      : "border-transparent text-white/60 hover:border-white/30 hover:text-white"
                  }`}
                >
                  {item.label}
                </button>
              );
            })}
          </nav>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6 sm:px-6">
        {view === "board" && (
          <AllResults
            athletes={athletes}
            gender={boardGender}
            onGender={setBoardGender}
            selectedId={selectedId}
            onSelect={setSelectedId}
          />
        )}
        {view === "regions" && (
          <RegionalBoard
            athletes={athletes}
            regionId={regionId}
            onRegion={setRegionId}
            gender={boardGender}
            onGender={setBoardGender}
            selectedId={selectedId}
            onSelect={setSelectedId}
          />
        )}
        {view === "classes" && (
          <CategoryBoard
            athletes={athletes}
            categoryId={categoryId}
            region={classRegion}
            selectedId={selectedId}
            onSelect={setSelectedId}
            onCategory={setCategoryId}
            onRegion={setClassRegion}
          />
        )}
        <footer className="mt-12 flex flex-col gap-2 border-t border-line py-6 text-sm text-faint sm:flex-row sm:items-center sm:justify-between">
          <p>Mock board for a single meet. Results stay in this browser session.</p>
          <button type="button" onClick={restore} className="text-left text-muted hover:text-chalk hover:underline">
            Restore original results
          </button>
        </footer>
      </main>

      <BackToTop />

      {selected && (
        <AthleteSheet
          athlete={selected}
          athletes={athletes}
          onClose={() => setSelectedId(null)}
          showClassLink={view !== "classes"}
          onShowClass={() => {
            setCategoryId(selected.categoryId);
            setSelectedId(null);
            setView("classes");
          }}
          onShowOverall={() => {
            setBoardGender(selected.gender);
            setSelectedId(null);
            setView("board");
          }}
        />
      )}
    </div>
  );
}

function BackToTop() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    function onScroll() {
      setVisible(window.scrollY > 480);
    }
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  if (!visible) return null;

  return (
    <button
      type="button"
      onClick={() => {
        const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        window.scrollTo({ top: 0, behavior: reduceMotion ? "auto" : "smooth" });
      }}
      className="fixed right-4 bottom-4 z-30 inline-flex items-center gap-2 rounded-full border border-white/20 bg-brand-black px-4 py-2 text-sm font-semibold tracking-wide text-white uppercase shadow-[0_8px_24px_rgba(0,0,0,0.35)]"
    >
      <UpIcon />
      Back to top
    </button>
  );
}

function UpIcon() {
  return (
    <svg viewBox="0 0 24 24" className="size-4" aria-hidden="true">
      <path
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M6 14l6-6 6 6"
      />
    </svg>
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

