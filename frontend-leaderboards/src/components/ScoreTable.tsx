import { useEffect, useRef, useState } from "react";
import { categoryById, events, weightChoices } from "../data";
import { columnHint, formatPoints, pointsLabel as formatPointsLabel, resultLines, resultTitle } from "../format";
import type { AthleteStanding } from "../types";

export type ResultSection = {
  key: string;
  title: string;
  caption: string;
  categoryId: string | null;
  rows: AthleteStanding[];
};

type ScoreTableProps = {
  sections: ResultSection[];
  pointsLabel: string;
  showClass?: boolean;
  showEventPoints: boolean;
  selectedId: string | null;
  onSelect: (athleteId: string) => void;
  onOpenCategory?: (categoryId: string) => void;
};

export function ScoreTable({
  sections,
  pointsLabel,
  showClass = true,
  showEventPoints,
  selectedId,
  onSelect,
  onOpenCategory,
}: ScoreTableProps) {
  const columnCount = 8 + (showClass ? 1 : 0);
  const scrollerRef = useRef<HTMLDivElement>(null);
  const tableRef = useRef<HTMLTableElement>(null);
  const theadRef = useRef<HTMLTableSectionElement>(null);
  const [pin, setPin] = useState<PinnedHeader | null>(null);

  useEffect(() => {
    const scroller = scrollerRef.current;
    const table = tableRef.current;
    const thead = theadRef.current;
    if (!scroller || !table || !thead) return;
    const scrollEl = scroller;
    const tableEl = table;
    const headEl = thead;

    const scheduled = { id: 0 };
    let current: PinnedHeader | null = null;

    function measure() {
      cancelAnimationFrame(scheduled.id);
      scheduled.id = requestAnimationFrame(() => {
        const headerHeight =
          Number.parseFloat(
            getComputedStyle(document.documentElement).getPropertyValue("--site-header"),
          ) || 0;
        const headRect = headEl.getBoundingClientRect();
        const tableRect = tableEl.getBoundingClientRect();
        const box = scrollEl.getBoundingClientRect();
        const stuck =
          headRect.top < headerHeight - 1 && tableRect.bottom > headerHeight + headRect.height;
        const next = stuck
          ? {
              top: headerHeight,
              left: box.left,
              width: box.width,
              scrollLeft: scrollEl.scrollLeft - scrollEl.clientLeft,
              tableWidth: tableEl.offsetWidth,
              colWidths: [...headEl.querySelectorAll("th")].map(
                (cell) => cell.getBoundingClientRect().width,
              ),
            }
          : null;
        if (samePin(current, next)) return;
        current = next;
        setPin(next);
      });
    }

    measure();
    window.addEventListener("scroll", measure, { passive: true });
    window.addEventListener("resize", measure);
    scrollEl.addEventListener("scroll", measure, { passive: true });
    const observer = new ResizeObserver(measure);
    observer.observe(scrollEl);
    return () => {
      cancelAnimationFrame(scheduled.id);
      window.removeEventListener("scroll", measure);
      window.removeEventListener("resize", measure);
      scrollEl.removeEventListener("scroll", measure);
      observer.disconnect();
    };
  }, [sections, showClass, pointsLabel]);

  const columnHeads = (widths?: number[]) => (
    <tr className="border-b border-line text-faint">
      <th
        scope="col"
        style={widths ? { width: widths[0] } : undefined}
        className="w-14 bg-panel px-3 py-3 font-display text-sm uppercase tracking-wide"
      >
        Place
      </th>
      <th
        scope="col"
        style={widths ? { width: widths[1] } : undefined}
        className="bg-panel px-3 py-3 font-display text-sm uppercase tracking-wide"
      >
        Athlete
      </th>
      {showClass && (
        <th
          scope="col"
          style={widths ? { width: widths[2] } : undefined}
          className="bg-panel px-3 py-3 font-display text-sm uppercase tracking-wide"
        >
          Class
        </th>
      )}
      <th
        scope="col"
        style={widths ? { width: widths[showClass ? 3 : 2] } : undefined}
        className="bg-panel px-3 py-3 text-right font-display text-sm uppercase tracking-wide"
      >
        {pointsLabel}
      </th>
      {events.map((event, index) => (
        <th
          key={event.id}
          scope="col"
          style={widths ? { width: widths[(showClass ? 4 : 3) + index] } : undefined}
          className="bg-panel px-3 py-3 text-right"
          title={weightChoices(event)}
        >
                <span className="block font-display text-sm uppercase tracking-wide">{event.short}</span>
                <span className="block font-mono text-[10px] font-normal tracking-normal">{columnHint(event)}</span>
              </th>
            ))}
    </tr>
  );

  return (
    <>
    <div ref={scrollerRef} className="overflow-x-auto rounded-lg border border-line bg-panel">
      <table ref={tableRef} className="w-full border-collapse text-left">
        <caption className="sr-only">Competition results</caption>
        <thead ref={theadRef}>{columnHeads()}</thead>
        <tbody>
          {sections.flatMap((section) => {
            const header = (
              <tr key={`${section.key}-heading`} className="bg-raised">
                <td colSpan={columnCount} className="px-3 py-3">
                  <div className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
                    <div className="flex flex-wrap items-baseline gap-3">
                      <p className="font-display text-lg uppercase tracking-wide text-chalk">
                        {section.title}
                      </p>
                      {section.categoryId && onOpenCategory && (
                        <button
                          type="button"
                          onClick={() => onOpenCategory(section.categoryId as string)}
                          className="text-sm text-brass hover:underline"
                        >
                          Open class
                        </button>
                      )}
                    </div>
                    <p className="max-w-xl text-sm text-muted">{section.caption}</p>
                  </div>
                </td>
              </tr>
            );

            const body = section.rows.map((row) => {
              const category = showClass ? categoryById(row.athlete.categoryId) : null;
              const selected = row.athlete.id === selectedId;
              return (
                <tr
                  key={row.athlete.id}
                  onClick={() => onSelect(row.athlete.id)}
                  className={`border-b border-line align-top ${
                    selected ? "bg-brass/10" : "hover:bg-ink"
                  }`}
                >
                  <td className="px-3 py-3 font-mono text-sm text-brass">
                    {row.tied ? `=${row.place}` : row.place}
                  </td>
                  <td className="px-3 py-3">
                    <button
                      type="button"
                      onClick={() => onSelect(row.athlete.id)}
                      className="text-left"
                    >
                      <span className="block font-display text-xl uppercase leading-none tracking-wide text-chalk">
                        {row.athlete.name}
                      </span>
                      <span className="mt-1 block text-xs text-muted">{row.athlete.country}</span>
                    </button>
                  </td>
                  {category && (
                    <td className="px-3 py-3">
                      <span className="rounded border border-line px-1.5 py-0.5 font-mono text-[11px] text-muted">
                        {category.short}
                      </span>
                    </td>
                  )}
                  <td className="px-3 py-3 text-right font-mono text-base tabular-nums text-chalk">
                    {formatPoints(row.total)}
                  </td>
                  {events.map((event) => {
                    const result = row.athlete.results[event.id];
                    const score = row.events[event.id];
                    const lines = resultLines(event, result);
                    return (
                      <td
                        key={event.id}
                        className="px-3 py-3 text-right"
                        title={resultTitle(event, result, score)}
                      >
                        {lines ? (
                          <>
                            <span className="block whitespace-nowrap font-mono text-sm tabular-nums text-chalk">
                              {lines.performance}
                            </span>
                            <span className="block whitespace-nowrap font-mono text-sm tabular-nums text-muted">
                              {lines.weight}
                            </span>
                          </>
                        ) : (
                          <span className="block font-mono text-sm text-faint">—</span>
                        )}
                        {showEventPoints && (
                          <span className="mt-0.5 block font-mono text-[11px] text-brass">
                            {score.place === null ? "—" : formatPointsLabel(score.points)}
                          </span>
                        )}
                      </td>
                    );
                  })}
                </tr>
              );
            });

            return section.title ? [header, ...body] : body;
          })}
        </tbody>
      </table>
    </div>
    {pin && (
      <div
        aria-hidden="true"
        className="pointer-events-none fixed z-20 overflow-hidden bg-panel shadow-[0_1px_0_var(--color-line)]"
        style={{ top: pin.top, left: pin.left, width: pin.width }}
      >
        <table
          className="border-collapse text-left"
          style={{ width: pin.tableWidth, transform: `translateX(-${pin.scrollLeft}px)` }}
        >
          <thead>{columnHeads(pin.colWidths)}</thead>
        </table>
      </div>
    )}
    </>
  );
}

type PinnedHeader = {
  top: number;
  left: number;
  width: number;
  scrollLeft: number;
  tableWidth: number;
  colWidths: number[];
};

function samePin(current: PinnedHeader | null, next: PinnedHeader | null): boolean {
  if (current === next) return true;
  if (!current || !next) return false;
  return (
    Math.abs(current.top - next.top) < 0.5 &&
    Math.abs(current.left - next.left) < 0.5 &&
    Math.abs(current.width - next.width) < 0.5 &&
    Math.abs(current.scrollLeft - next.scrollLeft) < 0.5 &&
    Math.abs(current.tableWidth - next.tableWidth) < 0.5 &&
    current.colWidths.length === next.colWidths.length &&
    current.colWidths.every((width, index) => Math.abs(width - next.colWidths[index]) < 0.5)
  );
}
