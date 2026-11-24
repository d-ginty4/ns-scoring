import {
  eventIds,
  type Athlete,
  type Category,
  type Competition,
  type EventDef,
  type EventId,
  type EventResult,
  type Gender,
  type Promoter,
} from "./oldTypes";

type AthleteSeed = Omit<Athlete, "competitionId" | "withdrawn">;

export const events: EventDef[] = [
  {
    id: "log",
    name: "Log for reps",
    short: "Log",
    weights: { men: [85, 100, 115, 130], women: [50, 60, 70, 80] },
    weightLabel: "kg",
    measure: "reps",
    note: "For reps",
  },
  {
    id: "deadlift",
    name: "Deadlift for reps",
    short: "Deadlift",
    weights: { men: [210, 240, 270, 300], women: [125, 150, 175, 200] },
    weightLabel: "kg",
    measure: "reps",
    note: "For reps",
  },
  {
    id: "loading",
    name: "Sandbag loading",
    short: "Loading",
    weights: { men: [100, 110, 120, 130], women: [60, 70, 80, 90] },
    weightLabel: "kg",
    measure: "time",
    note: "4 bags from 10, 8, 6 and 4 m",
  },
  {
    id: "farmers",
    name: "Farmers",
    short: "Farmers",
    weights: { men: [110, 130], women: [70, 90] },
    weightLabel: "kg each hand",
    measure: "time",
    note: "For time",
  },
  {
    id: "carry",
    name: "Sandbag carry",
    short: "Carry",
    weights: { men: [110, 130], women: [70, 90] },
    weightLabel: "kg",
    measure: "distance",
    note: "For distance",
  },
];

export const categories: Category[] = [
  { id: "m72", gender: "men", label: "Under 72 kg", short: "U72" },
  { id: "m80", gender: "men", label: "Under 80 kg", short: "U80" },
  { id: "m90", gender: "men", label: "Under 90 kg", short: "U90" },
  { id: "m105", gender: "men", label: "Under 105 kg", short: "U105" },
  { id: "mopen", gender: "men", label: "Open", short: "Open" },
  { id: "mjr", gender: "men", label: "Junior", short: "Junior" },
  { id: "m90m", gender: "men", label: "Under 90 kg master", short: "U90 master" },
  { id: "mopenm", gender: "men", label: "Open masters", short: "Open masters" },
  { id: "w57", gender: "women", label: "Under 57 kg", short: "U57" },
  { id: "w64", gender: "women", label: "Under 64 kg", short: "U64" },
  { id: "w73", gender: "women", label: "Under 73 kg", short: "U73" },
  { id: "w82", gender: "women", label: "Under 82 kg", short: "U82" },
  { id: "wopen", gender: "women", label: "Open", short: "Open" },
  { id: "wjr", gender: "women", label: "Junior", short: "Junior" },
  { id: "w73m", gender: "women", label: "Under 73 kg master", short: "U73 master" },
  { id: "wopenm", gender: "women", label: "Open master", short: "Open master" },
];

export const competitions: Competition[] = [
  { id: "ireland", name: "Ireland Nationals", country: "Ireland" },
  { id: "england", name: "England Nationals", country: "England" },
  { id: "scotland", name: "Scotland Nationals", country: "Scotland" },
  { id: "wales", name: "Wales Nationals", country: "Wales" },
];

const homeCompetition: Record<string, string> = {
  Ireland: "ireland",
  England: "england",
  Scotland: "scotland",
  Wales: "wales",
};

export const promoters: Promoter[] = [
  {
    id: "northshore",
    name: "Northshore Strength",
    short: "Northshore",
    region: "United Kingdom",
  },
  {
    id: "baltic",
    name: "Baltic Iron",
    short: "Baltic Iron",
    region: "Latvia",
  },
  {
    id: "cape",
    name: "Cape Strong",
    short: "Cape Strong",
    region: "South Africa",
  },
  {
    id: "pacific",
    name: "Pacific Barbell",
    short: "Pacific Barbell",
    region: "Australia",
  },
  {
    id: "heartland",
    name: "Heartland Shows",
    short: "Heartland",
    region: "United States",
  },
];

type Lift = [number, number] | null;

function result(lift: Lift, promoterId: string): EventResult {
  if (!lift) return { weight: null, value: null, enteredBy: null, skipped: false };
  return { weight: lift[0], value: lift[1], enteredBy: promoterId, skipped: false };
}

function results(
  promoterId: string,
  log: Lift,
  deadlift: Lift,
  loading: Lift,
  farmers: Lift,
  carry: Lift,
): Record<EventId, EventResult> {
  return {
    log: result(log, promoterId),
    deadlift: result(deadlift, promoterId),
    loading: result(loading, promoterId),
    farmers: result(farmers, promoterId),
    carry: result(carry, promoterId),
  };
}

export function weightLine(event: EventDef, gender: Gender): string {
  return `${event.weights[gender].join(" · ")} ${event.weightLabel}`;
}

export function weightChoices(event: EventDef, gender?: Gender): string {
  if (gender) return weightLine(event, gender);
  return `Men ${weightLine(event, "men")}. Women ${weightLine(event, "women")}`;
}

const namedAthletes = [
  {
    id: "callum",
    name: "Callum Reid",
    country: "Scotland",
    gender: "men",
    categoryId: "m72",
    bodyweight: 70,
    results: results("northshore", [100, 8], [240, 6], [110, 28.4], [110, 24.6], [110, 40]),
  },
  {
    id: "mateusz",
    name: "Mateusz Kwiatkowski",
    country: "Poland",
    gender: "men",
    categoryId: "m72",
    bodyweight: 71,
    results: results("baltic", [100, 10], [240, 8], [110, 26.1], [110, 22.1], [110, 36]),
  },
  {
    id: "kenji",
    name: "Kenji Sato",
    country: "Japan",
    gender: "men",
    categoryId: "m80",
    bodyweight: 78,
    results: results("pacific", [85, 12], [210, 10], [100, 22], [110, 26.8], null),
  },
  {
    id: "jonas",
    name: "Jonas Berg",
    country: "Sweden",
    gender: "men",
    categoryId: "m80",
    bodyweight: 79,
    results: results("baltic", [115, 6], [270, 5], [120, 31.5], [130, 23.4], [130, 22]),
  },
  {
    id: "andre",
    name: "André Silva",
    country: "Brazil",
    gender: "men",
    categoryId: "m90",
    bodyweight: 87,
    results: results("heartland", [115, 8], [270, 4], [120, 29.2], [130, 21], [130, 26]),
  },
  {
    id: "thomas",
    name: "Thomas Wright",
    country: "England",
    gender: "men",
    categoryId: "m90",
    bodyweight: 89,
    results: results("northshore", [100, 9], [240, 9], [110, 24.8], [110, 19.5], [110, 48]),
  },
  {
    id: "ivo",
    name: "Ivo Petrov",
    country: "Bulgaria",
    gender: "men",
    categoryId: "m105",
    bodyweight: 99,
    results: results("baltic", [130, 5], [300, 4], [130, 34.6], [130, 20.2], [130, 30]),
  },
  {
    id: "malik",
    name: "Malik Johnson",
    country: "United States",
    gender: "men",
    categoryId: "m105",
    bodyweight: 104,
    results: results("heartland", [130, 7], [300, 3], [130, 33.1], [130, 18.8], [130, 34]),
  },
  {
    id: "eoin",
    name: "Eoin Murphy",
    country: "Ireland",
    gender: "men",
    categoryId: "mopen",
    bodyweight: 118,
    results: results("northshore", [115, 9], [270, 6], [120, 27.4], [110, 17.4], [130, 20]),
  },
  {
    id: "henrik",
    name: "Henrik Dahl",
    country: "Norway",
    gender: "men",
    categoryId: "mjr",
    bodyweight: 84,
    results: results("baltic", [130, 9], [300, 6], [130, 30.2], [130, 16.9], [130, 38]),
  },
  {
    id: "samuel",
    name: "Samuel Okonkwo",
    country: "Nigeria",
    gender: "men",
    categoryId: "m90m",
    bodyweight: 86,
    results: results("cape", [130, 8], [300, 5], [130, 32], [130, 18.1], [130, 36]),
  },
  {
    id: "fraser",
    name: "Fraser MacLeod",
    country: "Scotland",
    gender: "men",
    categoryId: "mopenm",
    bodyweight: 132,
    results: results("northshore", [115, 10], [270, 7], [120, 26.5], [130, 19.6], null),
  },
  {
    id: "aoife",
    name: "Aoife Nolan",
    country: "Ireland",
    gender: "women",
    categoryId: "w57",
    bodyweight: 54,
    results: results("northshore", [60, 6], [150, 5], [70, 30.4], [70, 26.2], [70, 34]),
  },
  {
    id: "yui",
    name: "Yui Tanaka",
    country: "Japan",
    gender: "women",
    categoryId: "w57",
    bodyweight: 56,
    results: results("pacific", [70, 8], [175, 4], [80, 33.8], [90, 24.8], [90, 18]),
  },
  {
    id: "lara",
    name: "Lara Costa",
    country: "Portugal",
    gender: "women",
    categoryId: "w64",
    bodyweight: 62,
    results: results("cape", [60, 5], [150, 4], [70, 32.1], [70, 28.4], null),
  },
  {
    id: "freya",
    name: "Freya Lind",
    country: "Iceland",
    gender: "women",
    categoryId: "w64",
    bodyweight: 63,
    results: results("baltic", [70, 6], [175, 5], [80, 31], [90, 22.6], [90, 24]),
  },
  {
    id: "naomi",
    name: "Naomi Okello",
    country: "Kenya",
    gender: "women",
    categoryId: "w73",
    bodyweight: 69,
    results: results("cape", [70, 5], [175, 6], [80, 29.6], [70, 21], [70, 40]),
  },
  {
    id: "hannah",
    name: "Hannah Brooks",
    country: "Canada",
    gender: "women",
    categoryId: "w73",
    bodyweight: 72,
    results: results("heartland", [80, 4], [150, 8], [70, 25.2], [90, 23.5], null),
  },
  {
    id: "zuzanna",
    name: "Zuzanna Wójcik",
    country: "Poland",
    gender: "women",
    categoryId: "w82",
    bodyweight: 78,
    results: results("baltic", [80, 6], [200, 4], [90, 36.4], [90, 20.4], [90, 28]),
  },
  {
    id: "priya",
    name: "Priya Raman",
    country: "India",
    gender: "women",
    categoryId: "w82",
    bodyweight: 80,
    results: results("pacific", [70, 9], [175, 6], [80, 28.8], [90, 21.7], [90, 26]),
  },
  {
    id: "elise",
    name: "Elise Moreau",
    country: "France",
    gender: "women",
    categoryId: "wopen",
    bodyweight: 90,
    results: results("northshore", [80, 5], [200, 3], [90, 34.2], [70, 19.2], [90, 22]),
  },
  {
    id: "amara",
    name: "Amara Diallo",
    country: "Senegal",
    gender: "women",
    categoryId: "w73m",
    bodyweight: 70,
    results: results("cape", [70, 7], [175, 7], [80, 27.5], [90, 22], null),
  },
  {
    id: "britt",
    name: "Britt Solberg",
    country: "Norway",
    gender: "women",
    categoryId: "wopenm",
    bodyweight: 96,
    results: results("baltic", [80, 4], [200, 4], [90, 35.1], [90, 25.1], [90, 16]),
  },
] satisfies AthleteSeed[];

const athletesPerCategoryPerCompetition = 10;
const athletesPerCategory = athletesPerCategoryPerCompetition * competitions.length;

const bodyweightRange: Record<string, [number, number]> = {
  m72: [62, 71],
  m80: [72, 79],
  m90: [80, 89],
  m105: [90, 104],
  mopen: [105, 140],
  mjr: [67, 96],
  m90m: [74, 89],
  mopenm: [96, 138],
  w57: [48, 56],
  w64: [57, 63],
  w73: [64, 72],
  w82: [73, 81],
  wopen: [82, 108],
  wjr: [52, 74],
  w73m: [58, 72],
  wopenm: [76, 105],
};

const origins: { country: string; men: string[]; women: string[]; lasts: string[] }[] = [
  {
    country: "Scotland",
    men: ["Graham", "Duncan", "Hamish", "Finlay", "Rory", "Angus"],
    women: ["Isla", "Fiona", "Mhairi", "Kirsty", "Eilidh", "Ailsa"],
    lasts: ["Campbell", "Murray", "Wallace", "Gordon", "Sinclair", "Burns"],
  },
  {
    country: "England",
    men: ["Harry", "George", "Oliver", "Arthur", "Edward", "Jack"],
    women: ["Amelia", "Florence", "Beatrice", "Clara", "Edith", "Rose"],
    lasts: ["Bennett", "Clarke", "Hughes", "Palmer", "Shaw", "Turner"],
  },
  {
    country: "Ireland",
    men: ["Cillian", "Oisin", "Padraig", "Sean", "Niall", "Darragh"],
    women: ["Niamh", "Siobhan", "Orla", "Ciara", "Maeve", "Saoirse"],
    lasts: ["Byrne", "Kelly", "Murphy", "O'Brien", "Walsh", "Doyle"],
  },
  {
    country: "Wales",
    men: ["Rhys", "Gareth", "Owain", "Dylan", "Ioan", "Sion"],
    women: ["Carys", "Bethan", "Seren", "Lowri", "Anwen", "Ffion"],
    lasts: ["Evans", "Davies", "Price", "Griffiths", "Morgan", "Howells"],
  },
  {
    country: "Norway",
    men: ["Sondre", "Magnus", "Erlend", "Jonas", "Leif", "Sverre"],
    women: ["Ingrid", "Astrid", "Liv", "Sigrid", "Maren", "Thea"],
    lasts: ["Hansen", "Nilsen", "Berg", "Hauge", "Solberg", "Dahl"],
  },
  {
    country: "Sweden",
    men: ["Erik", "Axel", "Nils", "Viktor", "Lars", "Oskar"],
    women: ["Linnea", "Elsa", "Maja", "Astrid", "Klara", "Ebba"],
    lasts: ["Lindqvist", "Andersson", "Bergstrom", "Holm", "Ek", "Nyström"],
  },
  {
    country: "Finland",
    men: ["Eero", "Mikko", "Onni", "Veikko", "Aaro", "Tapio"],
    women: ["Aino", "Helmi", "Siiri", "Venla", "Elli", "Pilvi"],
    lasts: ["Korhonen", "Virtanen", "Mäkinen", "Laine", "Heikkinen", "Nieminen"],
  },
  {
    country: "Estonia",
    men: ["Mart", "Ragnar", "Karl", "Toomas", "Andres", "Jaan"],
    women: ["Kadri", "Liis", "Mari", "Kertu", "Anu", "Piret"],
    lasts: ["Tamm", "Saar", "Kask", "Mägi", "Sepp", "Ilves"],
  },
  {
    country: "Poland",
    men: ["Piotr", "Tomasz", "Krzysztof", "Adam", "Bartek", "Wojciech"],
    women: ["Anna", "Kasia", "Magda", "Ola", "Natalia", "Agnieszka"],
    lasts: ["Nowak", "Wisniewski", "Kaminski", "Lewandowski", "Zielinski", "Kowalski"],
  },
  {
    country: "Czechia",
    men: ["Jan", "Tomas", "Lukas", "Ondrej", "Petr", "Marek"],
    women: ["Tereza", "Karolina", "Barbora", "Eliska", "Adela", "Nikola"],
    lasts: ["Novak", "Svoboda", "Dvorak", "Horak", "Kucera", "Vesely"],
  },
  {
    country: "Germany",
    men: ["Lukas", "Felix", "Jonas", "Max", "Leon", "Paul"],
    women: ["Lena", "Mia", "Greta", "Klara", "Anja", "Heike"],
    lasts: ["Schneider", "Fischer", "Weber", "Becker", "Hoffmann", "Keller"],
  },
  {
    country: "Netherlands",
    men: ["Daan", "Bram", "Luuk", "Sven", "Thijs", "Joris"],
    women: ["Sanne", "Fleur", "Lotte", "Noor", "Anouk", "Eva"],
    lasts: ["de Vries", "Jansen", "Bakker", "Visser", "Smit", "Meijer"],
  },
  {
    country: "France",
    men: ["Hugo", "Louis", "Mathis", "Julien", "Antoine", "Nicolas"],
    women: ["Camille", "Manon", "Chloe", "Lea", "Ines", "Margot"],
    lasts: ["Bernard", "Dubois", "Laurent", "Petit", "Roux", "Morel"],
  },
  {
    country: "Spain",
    men: ["Alejandro", "Pablo", "Diego", "Javier", "Sergio", "Alvaro"],
    women: ["Lucia", "Marta", "Carmen", "Elena", "Nuria", "Raquel"],
    lasts: ["Garcia", "Martinez", "Lopez", "Sanchez", "Romero", "Navarro"],
  },
  {
    country: "Portugal",
    men: ["Tiago", "Diogo", "Rui", "Goncalo", "Miguel", "Nuno"],
    women: ["Ines", "Beatriz", "Catarina", "Mariana", "Rita", "Joana"],
    lasts: ["Santos", "Ferreira", "Oliveira", "Costa", "Pereira", "Almeida"],
  },
  {
    country: "Italy",
    men: ["Marco", "Luca", "Matteo", "Andrea", "Giovanni", "Alessio"],
    women: ["Giulia", "Chiara", "Francesca", "Elena", "Sofia", "Alessia"],
    lasts: ["Rossi", "Russo", "Ferrari", "Esposito", "Bianchi", "Romano"],
  },
  {
    country: "Bulgaria",
    men: ["Georgi", "Nikolai", "Dimitar", "Ivan", "Stefan", "Petar"],
    women: ["Elena", "Maria", "Desislava", "Iva", "Nadezhda", "Viktoria"],
    lasts: ["Ivanov", "Dimitrov", "Georgiev", "Nikolov", "Petrov", "Stoyanov"],
  },
  {
    country: "Ukraine",
    men: ["Oleksandr", "Dmytro", "Andriy", "Taras", "Bohdan", "Maksym"],
    women: ["Olena", "Yulia", "Sofiia", "Kateryna", "Dariia", "Nataliia"],
    lasts: ["Shevchenko", "Kovalenko", "Bondarenko", "Melnyk", "Kravchenko", "Savchenko"],
  },
  {
    country: "United States",
    men: ["Ethan", "Noah", "Caleb", "Marcus", "Cole", "Owen"],
    women: ["Harper", "Maya", "Claire", "Avery", "Nora", "Quinn"],
    lasts: ["Carter", "Brooks", "Hayes", "Reed", "Foster", "Bennett"],
  },
  {
    country: "Canada",
    men: ["Logan", "Evan", "Nathan", "Cole", "Ryan", "Blake"],
    women: ["Megan", "Lauren", "Abby", "Jenna", "Paige", "Tessa"],
    lasts: ["MacDonald", "Tremblay", "Gagnon", "Campbell", "Roy", "Wilson"],
  },
  {
    country: "Brazil",
    men: ["Lucas", "Rafael", "Bruno", "Gustavo", "Felipe", "Caio"],
    women: ["Ana", "Juliana", "Camila", "Fernanda", "Bianca", "Leticia"],
    lasts: ["Oliveira", "Souza", "Almeida", "Ribeiro", "Carvalho", "Gomes"],
  },
  {
    country: "Argentina",
    men: ["Mateo", "Santiago", "Facundo", "Agustin", "Nicolas", "Lautaro"],
    women: ["Valentina", "Martina", "Camila", "Lucia", "Sofia", "Catalina"],
    lasts: ["Gonzalez", "Fernandez", "Diaz", "Acosta", "Perez", "Suarez"],
  },
  {
    country: "Mexico",
    men: ["Santiago", "Diego", "Emiliano", "Mateo", "Javier", "Hector"],
    women: ["Renata", "Ximena", "Daniela", "Regina", "Fernanda", "Paola"],
    lasts: ["Hernandez", "Ramirez", "Cruz", "Flores", "Reyes", "Morales"],
  },
  {
    country: "Nigeria",
    men: ["Chinedu", "Emeka", "Tunde", "Ifeanyi", "Kunle", "Segun"],
    women: ["Adaeze", "Chioma", "Ngozi", "Funke", "Zainab", "Halima"],
    lasts: ["Adeyemi", "Okafor", "Balogun", "Eze", "Nwosu", "Adebayo"],
  },
  {
    country: "Kenya",
    men: ["Brian", "Kevin", "Daniel", "Peter", "Joseph", "David"],
    women: ["Wanjiku", "Akinyi", "Nyambura", "Chebet", "Aisha", "Faith"],
    lasts: ["Otieno", "Mwangi", "Kamau", "Wanjala", "Kipchoge", "Omondi"],
  },
  {
    country: "South Africa",
    men: ["Johan", "Pieter", "Thabo", "Sipho", "Ruan", "Liam"],
    women: ["Annelie", "Lerato", "Zanele", "Marieke", "Naledi", "Sarah"],
    lasts: ["Botha", "van der Merwe", "Nkosi", "Dlamini", "Pretorius", "Mokoena"],
  },
  {
    country: "Senegal",
    men: ["Moussa", "Ibrahima", "Ousmane", "Cheikh", "Mamadou", "Abdou"],
    women: ["Awa", "Fatou", "Aminata", "Mariam", "Khady", "Astou"],
    lasts: ["Diop", "Ndiaye", "Ba", "Sow", "Fall", "Sarr"],
  },
  {
    country: "India",
    men: ["Arjun", "Rohan", "Vikram", "Kabir", "Aman", "Dev"],
    women: ["Ananya", "Meera", "Diya", "Kavya", "Sana", "Ishita"],
    lasts: ["Sharma", "Patel", "Singh", "Iyer", "Reddy", "Nair"],
  },
  {
    country: "Japan",
    men: ["Haruto", "Sota", "Ren", "Kaito", "Yuki", "Daiki"],
    women: ["Hana", "Aoi", "Mei", "Sakura", "Rin", "Yuna"],
    lasts: ["Sato", "Suzuki", "Takahashi", "Watanabe", "Ito", "Nakamura"],
  },
  {
    country: "Australia",
    men: ["Lachlan", "Cooper", "Mitchell", "Riley", "Harrison", "Nate"],
    women: ["Matilda", "Ruby", "Chloe", "Georgia", "Ellie", "Zoe"],
    lasts: ["Walker", "Thompson", "Harris", "Young", "King", "Scott"],
  },
  {
    country: "New Zealand",
    men: ["Wiremu", "James", "Callan", "Tane", "Hamish", "Beau"],
    women: ["Aroha", "Maia", "Charlotte", "Holly", "Grace", "Paige"],
    lasts: ["Williams", "Wilson", "Taylor", "Brown", "Ngata", "Rangi"],
  },
];

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value));
}

function mulberry32(seed: number): () => number {
  let state = seed >>> 0;
  return () => {
    state = (state + 0x6d2b79f5) >>> 0;
    let t = state;
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

function slug(value: string): string {
  return value
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

function generatedAthletes(existing: AthleteSeed[]): AthleteSeed[] {
  const usedNames = new Set(existing.map((athlete) => athlete.name));
  const usedIds = new Set(existing.map((athlete) => athlete.id));
  const people = origins.flatMap((origin) => {
    const pairs: { name: string; country: string; gender: Gender }[] = [];
    for (const first of origin.men) {
      for (const last of origin.lasts) {
        pairs.push({ name: `${first} ${last}`, country: origin.country, gender: "men" });
      }
    }
    for (const first of origin.women) {
      for (const last of origin.lasts) {
        pairs.push({ name: `${first} ${last}`, country: origin.country, gender: "women" });
      }
    }
    return pairs;
  });
  const rand = mulberry32(0x51a7);
  for (let index = people.length - 1; index > 0; index -= 1) {
    const swap = Math.floor(rand() * (index + 1));
    const current = people[index];
    people[index] = people[swap];
    people[swap] = current;
  }

  let cursor = 0;
  function nextPerson(gender: Gender): { name: string; country: string } {
    while (cursor < people.length) {
      const person = people[cursor];
      cursor += 1;
      if (person.gender === gender && !usedNames.has(person.name)) {
        usedNames.add(person.name);
        return person;
      }
    }
    throw new Error(`Ran out of ${gender} names`);
  }

  const added: AthleteSeed[] = [];
  for (const category of categories) {
    if (category.id === "wjr") continue;
    const have = existing.filter((athlete) => athlete.categoryId === category.id).length;
    const [minWeight, maxWeight] = bodyweightRange[category.id];
    for (let slot = 0; slot < athletesPerCategory - have; slot += 1) {
      const person = nextPerson(category.gender);
      const skill = 1 - slot / Math.max(1, athletesPerCategory - have - 1);
      const local = mulberry32((category.id.charCodeAt(0) << 16) ^ (slot + 1) * 997);
      const span = maxWeight - minWeight;
      const bodyweight = minWeight + Math.round((1 - skill) * span * 0.85 + local() * span * 0.15);
      const promoterId = promoters[Math.floor(local() * promoters.length)].id;
      const blank: EventId | null = null;
      let id = slug(`${person.name}-${category.id}`);
      while (usedIds.has(id)) id = `${id}-${slot}`;
      usedIds.add(id);
      added.push({
        id,
        name: person.name,
        country: person.country,
        gender: category.gender,
        categoryId: category.id,
        bodyweight: clamp(bodyweight, minWeight, maxWeight),
        results: buildResults(category.gender, skill, local, promoterId, blank),
      });
    }
  }
  return added;
}

function buildResults(
  gender: Gender,
  skill: number,
  rand: () => number,
  promoterId: string,
  blank: EventId | null,
): Record<EventId, EventResult> {
  const lifts = {} as Record<EventId, EventResult>;
  for (const event of events) {
    if (event.id === blank) {
      lifts[event.id] = result(null, promoterId);
      continue;
    }
    const weights = event.weights[gender];
    const noisy = clamp(skill + (rand() - 0.5) * 0.4, 0, 1);
    const weightIndex = Math.min(weights.length - 1, Math.floor(noisy * weights.length));
    const weight = weights[weightIndex];
    const weightBias = weights.length === 1 ? 0 : weightIndex / (weights.length - 1);
    let value: number;
    if (event.measure === "reps") {
      value = clamp(Math.round(3 + noisy * 9 - weightBias * 2 + (rand() - 0.5) * 2), 1, 14);
    } else if (event.id === "loading") {
      const seconds = 39 - noisy * 16 + weightBias * 3 + (rand() - 0.5) * 1.8;
      value = Math.round(clamp(seconds, 18, 48) * 10) / 10;
    } else if (event.id === "farmers") {
      const seconds = 30 - noisy * 12 + weightBias * 2.5 + (rand() - 0.5) * 1.4;
      value = Math.round(clamp(seconds, 14, 36) * 10) / 10;
    } else {
      value = clamp(Math.round(12 + noisy * 34 - weightBias * 8 + (rand() - 0.5) * 6), 8, 60);
    }
    lifts[event.id] = result([weight, value], promoterId);
  }
  return lifts;
}

function leastFilled(counts: Record<string, number>): string {
  let best = competitions[0].id;
  for (const competition of competitions) {
    if (counts[competition.id] < counts[best]) best = competition.id;
  }
  return best;
}

function assignCompetitions(athletes: AthleteSeed[]): Athlete[] {
  const byCategory = new Map<string, AthleteSeed[]>();
  for (const athlete of athletes) {
    const group = byCategory.get(athlete.categoryId) ?? [];
    group.push(athlete);
    byCategory.set(athlete.categoryId, group);
  }

  const assigned = new Map<string, string>();
  for (const group of byCategory.values()) {
    const counts = Object.fromEntries(competitions.map((competition) => [competition.id, 0]));
    const share = Math.ceil(group.length / competitions.length);
    const pending: AthleteSeed[] = [];
    for (const athlete of group) {
      const home = homeCompetition[athlete.country];
      if (home && counts[home] < share) {
        assigned.set(athlete.id, home);
        counts[home] += 1;
      } else {
        pending.push(athlete);
      }
    }
    for (const athlete of pending) {
      const competitionId = leastFilled(counts);
      assigned.set(athlete.id, competitionId);
      counts[competitionId] += 1;
    }
  }

  return athletes.map((athlete) => ({
    ...athlete,
    withdrawn: false,
    competitionId: assigned.get(athlete.id) ?? competitions[0].id,
  }));
}

const missingByCategory: Record<string, Partial<Record<EventId, number>>> = {
  w57: {},
  w64: { log: 2 },
  w73: { deadlift: 1, loading: 3 },
  w82: { farmers: 4, carry: 1 },
  wopen: { log: 2, carry: 2 },
  wjr: {},
  w73m: { deadlift: 3 },
  wopenm: { log: 1, deadlift: 1, loading: 1, farmers: 1, carry: 1 },
  m72: { log: 4, deadlift: 1, farmers: 2 },
  m80: {},
  m90: { loading: 3 },
  m105: { log: 5, carry: 2 },
  mopen: {},
  mjr: { farmers: 1 },
  m90m: { deadlift: 2, carry: 2 },
  mopenm: { log: 1, loading: 4 },
};

function placeholderResult(gender: Gender, eventId: EventId): EventResult {
  const event = events.find((item) => item.id === eventId) ?? events[0];
  const weight = event.weights[gender][0];
  const value = event.measure === "reps" ? 5 : event.measure === "time" ? 30 : 20;
  return { weight, value, enteredBy: "northshore", skipped: false };
}

function withFullCards(athletes: Athlete[]): Athlete[] {
  return athletes.map((athlete) => {
    const next = { ...athlete.results };
    for (const eventId of eventIds) {
      const result = next[eventId];
      if (result.weight === null || result.value === null) next[eventId] = placeholderResult(athlete.gender, eventId);
    }
    return { ...athlete, results: next };
  });
}

function applyMissingResults(athletes: Athlete[]): Athlete[] {
  const next = athletes.map((athlete) => ({ ...athlete, results: { ...athlete.results } }));
  for (const competition of competitions) {
    for (const category of categories) {
      const plan = missingByCategory[category.id];
      if (!plan) continue;
      const field = next
        .filter((athlete) => athlete.competitionId === competition.id && athlete.categoryId === category.id)
        .sort((a, b) => a.name.localeCompare(b.name));
      let cursor = 0;
      for (const eventId of eventIds) {
        const count = plan[eventId] ?? 0;
        for (let index = 0; index < count && cursor < field.length; index += 1, cursor += 1) {
          field[cursor].results = {
            ...field[cursor].results,
            [eventId]: { weight: null, value: null, enteredBy: null, skipped: false },
          };
        }
      }
    }
  }
  return next;
}

export const countries: string[] = [
  ...new Set([...origins.map((origin) => origin.country), ...namedAthletes.map((athlete) => athlete.country)]),
].sort((a, b) => a.localeCompare(b));

export type NewAthlete = {
  name: string;
  country: string;
  categoryId: string;
  competitionId: string;
};

function defaultBodyweight(categoryId: string): number {
  const [minWeight, maxWeight] = bodyweightRange[categoryId] ?? [40, 160];
  return Math.round((minWeight + maxWeight) / 2);
}

export function createAthlete(draft: NewAthlete, existing: Athlete[]): Athlete {
  const category = categoryById(draft.categoryId);
  const usedIds = new Set(existing.map((athlete) => athlete.id));
  const base = slug(`${draft.name}-${draft.categoryId}`);
  let id = base;
  let spare = 2;
  while (usedIds.has(id)) {
    id = `${base}-${spare}`;
    spare += 1;
  }
  return {
    id,
    name: draft.name.trim(),
    country: draft.country.trim(),
    gender: category.gender,
    categoryId: category.id,
    competitionId: draft.competitionId,
    bodyweight: defaultBodyweight(category.id),
    withdrawn: false,
    results: results("", null, null, null, null, null),
  };
}

const scoredAthletes = assignCompetitions([...namedAthletes, ...generatedAthletes(namedAthletes)]);

export const seedAthletes: Athlete[] = applyMissingResults(withFullCards(scoredAthletes));

export function genderLabel(gender: Gender): string {
  return gender === "men" ? "Men" : "Women";
}

export function categoryById(id: string): Category {
  const category = categories.find((item) => item.id === id);
  if (!category) {
    throw new Error(`Unknown category ${id}`);
  }
  return category;
}

export function promoterById(id: string | null): Promoter | null {
  if (!id) return null;
  return promoters.find((promoter) => promoter.id === id) ?? null;
}

export function promoterNames(athlete: Athlete): string[] {
  const names = new Set<string>();
  for (const event of events) {
    const promoter = promoterById(athlete.results[event.id].enteredBy);
    if (promoter) names.add(promoter.short);
  }
  return [...names];
}
