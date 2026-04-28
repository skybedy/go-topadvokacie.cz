# DECISIONS.md

## 2026-04-27 - Aktualni projektova rozhodnuti

### Projekt zustava male lokalni Go demo

FilipAiPilot je vedome male MVP, ne produkcni pravni system. Prioritou je rychle a srozumitelne predvest Filipovi praktickou hodnotu AI workflow nad pravnimi texty.

Duvod: Filip uz zna obecny ChatGPT chat. Demo ma ukazat hlavne rozdil mezi rucnim promptovanim a opakovatelnym pracovnim nastrojem.

### Hlavni stack je Go + server-side HTML

Pouzite technologie:

- Go 1.22
- `net/http`
- `html/template`
- server-side renderovane HTML
- Tailwind CSS pres CDN
- Vanilla JavaScript
- `.env` konfigurace

Duvod: minimum zavislosti, jednoduche lokalni spusteni, dobra citelnost kodu a snadna udrzba maleho MVP.

### Bez API klice musi demo fungovat v mock rezimu

Kdyz neni nastaven `OPENAI_API_KEY`, aplikace pouzije `MockAIClient` a predpripravene vystupy.

Duvod: demo musi jit ukazat okamzite i bez OpenAI billingu, site nebo funkcniho API klice.

### OpenAI model se nastavuje pres `.env`

Model se ridi promennou `OPENAI_MODEL`. Aktualni `.env.example` pouziva `gpt-5-nano`. Timeout se ridi `OPENAI_TIMEOUT_SECONDS`, vychozi hodnota v kodu je 180 sekund.

Duvod: modely a dostupnost API se meni, proto nema byt model natvrdo svazany s kodem ani s API klicem.

### Nepouziva se automaticke odesilani e-mailu

E-mailovy prompt muze pripravit pouze navrh e-mailu bez odeslani.

Duvod: Filip zatim nechce sverovat AI e-mailovou schranku ani automatickou komunikaci. Bezpecnejsi a vhodnejsi pro demo je jen draft, ktery pravnik zkontroluje.

### Vystup AI je strukturovany jako `Result`

AI klient ocekava JSON objekt s poli:

- `title`
- `summary`
- `sections`
- `warnings`
- `raw`

Duvod: strukturovany vystup se lepe zobrazuje v UI a ukazuje hodnotu workflow nastroje oproti volnemu chatu.

### Hlavni jednotkou produktu je ulozeny prompt

Samostatna vrstva vestavenych workflow byla odstranena. Pracovni panel pouziva jen ulozene pravni prompty z `internal/ai/client.go`.

Duvod: kdyz je aplikace napojena na realny model, paralelni seznam "workflow" a "promptu" zbytecne motal UI i produktove vysvetleni. Ulozeny prompt je srozumitelnejsi jednotka: ma nazev, kategorii, verzi, popis a instrukci.

### Vystup ma nastavitelny rozsah a perspektivu

Pracovni panel podporuje delku vystupu `brief`, `standard` a `detailed` a perspektivu `lawyer`, `client` a `negotiation`.

Duvod: stejne vstupni pravni zneni muze pravnik potrebovat jako interni analyzu, klientsky srozumitelny vystup nebo vyjednavaci podklad. Tato volba lepe ukazuje rozdil mezi obecnym chatem a pracovnim nastrojem.

### Demo rozsiruje prompt knihovnu misto pridavani slozite infrastruktury

Byly pridany prakticke prompty: red flags pred podpisem, vyjednavaci pozice, priprava hovoru s klientem, chybejici oblasti smlouvy, komentare do revize a executive summary pro jednatele.

Duvod: pro Filipovo demo ma vetsi hodnotu sirsi sada konkretnich pravnickych workflow nez predcasna databaze, prihlaseni nebo produkcni infrastruktura.

### Kopirovani sekci je preferovana mala UX funkce

Vystupni karty maji tlacitko pro kopirovani jednotlive sekce do schranky.

Duvod: pravnik bude casti vystupu pravdepodobne presouvat do e-mailu, poznamek, revize nebo klientskych podkladu. Kopirovani ukazuje praktickou pouzitelnost bez integrace do e-mailu.

### Upload podporuje text, DOCX a PDF v MVP rezimu

Podporovane jsou `.txt`, `.md`, `.markdown`, `.csv`, `.rst`, `.log`, `.docx` a `.pdf`. DOCX se parsuje primo v Go pres ZIP/XML. PDF se prevadi pres lokalni `pdftotext` z balicku `poppler-utils`. Stary `.doc` zatim podporovany neni.

Duvod: Filip bude pravdepodobne pracovat hlavne s PDF a Word dokumenty. Tahle MVP cesta prida prakticky wow efekt bez tezke infrastruktury a bez velke Go zavislosti.

### Soubor se pred analyzou nacita do viditelneho textoveho pole

Frontend po vyberu souboru vola `/upload-text`, zobrazi stav nacitani a vlozi vytazeny text do textarea. Teprve potom uzivatel spousti analyzu.

Duvod: pravnik vidi, co se z dokumentu opravdu nacetlo, muze text rychle zkontrolovat nebo upravit a demo pusobi kontrolovaneji nez skryty upload az pri odeslani formulare.

### Dlouhe operace maji viditelny stav

UI zobrazuje spinner pri nacitani souboru do textoveho pole a samostatny loading panel pri cekani na odpoved modelu.

Duvod: PDF/DOCX parsing i OpenAI odpoved mohou trvat nekolik sekund. Filip musi videt, ze aplikace pracuje a ne zasekla se.

### Deploy script zustava jednoduchy a nedestruktivni

Projekt ma `deploy.sh`, ktery provede `git pull`, `go test ./...`, `go build -o filipaipilot ./cmd/filipaipilot` a restartuje systemd sluzbu `filipaipilot`, pokud existuje.

Duvod: FilipAiPilot je male Go MVP bez Node buildu. Agresivni mazani zdrojaku a `skip-worktree` uklid z jine aplikace by tady pridaly krehkost bez zasadniho prinosu.

### Bezpecnostni upozorneni je soucast produktu

V aplikaci ma zustat jasna formulace:

> Demo nastroj. Vystupy AI slouzi pouze jako pracovni podklad pro pravnika a nenahrazuji odborne pravni posouzeni.

Duvod: jde o pravni domenu a demo nesmi pusobit jako nahrada odborneho posouzeni.
