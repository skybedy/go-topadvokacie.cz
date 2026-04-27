# DECISIONS.md

## 2026-04-27 - Aktualni projektova rozhodnuti

### Projekt zustava male lokalni Go demo

LexPilot Demo je vedome male MVP, ne produkcni pravni system. Prioritou je rychle a srozumitelne predvest Filipovi praktickou hodnotu AI workflow nad pravnimi texty.

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

Workflow muze pripravit pouze navrh e-mailu bez odeslani.

Duvod: Filip zatim nechce sverovat AI e-mailovou schranku ani automatickou komunikaci. Bezpecnejsi a vhodnejsi pro demo je jen draft, ktery pravnik zkontroluje.

### Vystup AI je strukturovany jako `Result`

AI klient ocekava JSON objekt s poli:

- `title`
- `summary`
- `sections`
- `warnings`
- `raw`

Duvod: strukturovany vystup se lepe zobrazuje v UI a ukazuje hodnotu workflow nastroje oproti volnemu chatu.

### Prompty jsou zatim v Go kodu

Workflow a prompt knihovna jsou aktualne v `internal/ai/client.go`.

Duvod: pro MVP je to nejjednodussi a nejprehlednejsi. Pozdeji muze davat smysl presun do databaze nebo verzovanych souboru.

### Upload je omezen na textove soubory

Podporovane jsou `.txt`, `.md`, `.markdown`, `.csv`, `.rst` a `.log`. PDF/DOCX je roadmapa.

Duvod: PDF/DOCX vyzaduji specializovane parsovani a zvysily by slozitost mimo hlavni cil dema.

### Bezpecnostni upozorneni je soucast produktu

V aplikaci ma zustat jasna formulace:

> Demo nastroj. Vystupy AI slouzi pouze jako pracovni podklad pro pravnika a nenahrazuji odborne pravni posouzeni.

Duvod: jde o pravni domenu a demo nesmi pusobit jako nahrada odborneho posouzeni.
