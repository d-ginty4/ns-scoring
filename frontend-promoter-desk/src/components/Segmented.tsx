type Option<T extends string> = {
  value: T;
  label: string;
};

type SegmentedProps<T extends string> = {
  label: string;
  value: T;
  onChange: (value: T) => void;
  options: readonly Option<T>[];
};

export function Segmented<T extends string>({
  label,
  value,
  onChange,
  options,
}: SegmentedProps<T>) {
  return (
    <div
      role="group"
      aria-label={label}
      className="inline-flex flex-wrap rounded-md border border-line bg-panel p-1"
    >
      {options.map((option) => {
        const active = option.value === value;
        return (
          <button
            key={option.value}
            type="button"
            aria-pressed={active}
            onClick={() => onChange(option.value)}
            className={`rounded px-3 py-1.5 text-sm font-medium ${
              active ? "bg-chalk text-ink" : "text-muted hover:text-chalk"
            }`}
          >
            {option.label}
          </button>
        );
      })}
    </div>
  );
}
