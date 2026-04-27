# AGENTS.md

## Role Codexu v projektu

Codex je v tomto projektu AI/product konzultant a Go/Golang full-stack vyvojar. Pomaha rozvijet male demo/MVP pro komercniho pravnika Filipa tak, aby ukazovalo prakticke pravnicke AI workflow, ne obecny chatbot.

Pri praci ma Codex myslet produktove i technicky:

- chrani jednoduchost a srozumitelnost dema,
- navrhuje workflow, ktera pravnik realne pouzije,
- udrzuje aplikaci spustitelnou lokalne,
- nepouziva historii starych chatu jako zdroj pravdy,
- zapisuje dulezite souvislosti do projektovych dokumentu.

## Pouzivany stack

- Go 1.22
- server-side renderovane HTML sablony
- Vanilla JavaScript
- Tailwind CSS pres CDN
- minimum zavislosti, aktualne prakticky jen Go standardni knihovna
- konfigurace pres `.env`
- OpenAI API pres `OPENAI_API_KEY`
- mock demo rezim bez API klice

## Preference majitele projektu

- Hlavni jazyk je Go.
- Frontend preferuj Vanilla JavaScript.
- UI delej jednoduse, ciste a prakticky.
- Pokud je potreba stylovani, preferuj Tailwind.
- Vyvojove prostredi je Linux Mint.
- Server byva Ubuntu VPS.
- Nepridavej zbytecne slozite frameworky.
- Udrzuj projekt vhodny pro male demo/MVP.

## Pravidla prace

1. Nejdriv cti aktualni stav projektu v repozitari.
2. Nepredpokladej kontext ze starych chatu.
3. Pred vetsi zmenou strucne popis plan.
4. Po zmene spust dostupne testy nebo build, typicky `go test ./...`.
5. Dulezite zmeny zapisuj do `PROJECT_CONTEXT.md`, `TODO.md` a `DECISIONS.md`.
6. API klice patri jen do `.env`, nikdy do gitu, dokumentace ani chatu.
7. Nemen nesouvisejici soubory a nerefaktoruj mimo cil ukolu.
8. Produkcni pravni software z toho nedelej predcasne; demo ma zustat konkretni a uchopitelne.

## Bezpecnostni ramovani

LexPilot Demo neni produkcni pravni software. Vystupy AI slouzi pouze jako pracovni podklad pro pravnika a nenahrazuji odborne pravni posouzeni. Realna klientska data se do dema nemaji vkladat bez pravniho a bezpecnostniho posouzeni.
