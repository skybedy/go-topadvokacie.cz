# FilipAiPilot

FilipAiPilot je lokální MVP aplikace pro běžného komerčního právníka. Ukazuje, jak lze obecné AI použití zabalit do konkrétních právnických workflow nad dokumenty, smlouvami a právními texty.

Aplikace běží jako jednoduchý Go server se server-side renderovaným HTML frontendem a Tailwind CSS přes CDN. Umí volat OpenAI i Gemini. Bez odpovídajícího API klíče funguje v mock režimu s připravenými odpověďmi.

## Spuštění

```bash
go run ./cmd/filipaipilot
```

Otevřete:

```text
http://localhost:8080
```

Volitelně vytvořte `.env` podle `.env.example`:

```bash
cp .env.example .env
```

Pro reálné OpenAI volání nastavte:

```env
AI_PROVIDER=openai
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
AI_TIMEOUT_SECONDS=180
```

Pro reálné Gemini volání nastavte:

```env
AI_PROVIDER=gemini
GEMINI_API_KEY=...
GEMINI_MODEL=gemini-2.5-flash
AI_TIMEOUT_SECONDS=180
```

## Produkční build a deploy

Na Ubuntu VPS je nejjednodušší buildnout Go binárku a spouštět ji přes systemd za Nginx reverse proxy:

```bash
go test ./...
go build -o filipaipilot ./cmd/filipaipilot
```

Pro PDF upload nainstalujte systémový parser:

```bash
sudo apt install poppler-utils
```

Součástí projektu je jednoduchý deploy script:

```bash
./deploy.sh
```

Script provede `git pull`, `go test ./...`, build binárky `filipaipilot` a pokusí se restartovat systemd službu `filipaipilot.service`, pokud na serveru existuje. `.env` zůstává mimo git a musí být připravený na serveru.

## Struktura

```text
cmd/filipaipilot/main.go     vstupní bod aplikace
internal/ai/             AI rozhraní, mock klient a OpenAI/Gemini klient
internal/model/          sdílené datové struktury
internal/web/            HTTP handlery a ukázková data
templates/               server-side HTML šablony
static/                  drobný CSS a Vanilla JavaScript
examples/                fiktivní právní texty
```

## Právní režimy

Dokument lze vložit ručně nebo nahrát jako `.txt`, `.md`, `.csv`, `.rst`, `.log`, `.docx` nebo `.pdf`. Po výběru souboru se jeho text nejdřív načte do textového pole, takže právník vidí obsah před spuštěním analýzy. DOCX import používá jednoduché ZIP/XML vytěžení textu přímo v Go. PDF import používá lokální nástroj `pdftotext` z balíčku `poppler-utils`, pokud je v systému dostupný. Starší binární `.doc` zatím podporovaný není; dokument je vhodné uložit jako `.docx` nebo PDF.

Pracovní panel používá jen uložené právní prompty. Neexistuje oddělená umělá vrstva workflow; každý režim je pojmenovaný, verzovaný prompt s vlastní instrukcí. Panel navíc umožňuje nastavit délku výstupu (`stručně`, `standardně`, `detailně`) a perspektivu (`pro právníka`, `pro klienta`, `pro vyjednávání`). Výstupní sekce lze kopírovat samostatně do schránky.

## Prompt knihovna

Součástí dema je i osobní právnická prompt knihovna. Nejde o další chat, ale o uložené a verzované šablony pro opakovanou práci:

- kontrola smlouvy,
- převod do srozumitelné řeči pro klienta,
- otázky na klienta,
- návrh e-mailu bez odeslání,
- protiargumentace protistrany,
- checklist před podpisem,
- revize obchodních podmínek,
- kontrola konzistence,
- extrakce povinností a lhůt,
- red flags před podpisem,
- vyjednávací pozice,
- příprava hovoru s klientem,
- kontrola chybějících oblastí ve smlouvě,
- komentáře do revize,
- executive summary pro jednatele,
- porovnání dvou verzí dokumentu.

V testovací verzi jsou prompty uložené přímo v Go kódu v `internal/ai/client.go`. V produkci by dávalo smysl přesunout je do databáze nebo verzovaných souborů, přidat vlastní prompty pro konkrétního právníka a měřit kvalitu výstupů na ukázkových dokumentech.

Výstup má strukturu `Result`: `Title`, `Summary`, `Sections`, `Warnings` a `Raw` fallback.

## Náklady a AI API

ChatGPT Plus paušál se pro tuto aplikaci typicky nepoužívá. Lokální aplikace volá externí AI API (OpenAI nebo Gemini), takže je potřeba samostatný API klíč a usage se účtuje zvlášť podle tokenů.

Bez API klíče aplikace běží v mock režimu zdarma a je použitelná pro prezentaci právních režimů.

## Bezpečnost a důvěrnost dat

Demo není produkční právní software. Při použití AI API se vložený text odesílá externí službě podle aktuálních podmínek a nastavení daného API účtu. Do demo prostředí nevkládejte skutečná důvěrná klientská data bez odpovídajícího právního a bezpečnostního posouzení.

Pro produkční nasazení doporučuji řešit minimálně:

- anonymizaci nebo pseudonymizaci dokumentů,
- řízení přístupů a auditní logy,
- šifrování dat v klidu i při přenosu,
- retenční politiku pro vstupy a výstupy,
- kontrolu promptů a verzování workflow,
- jasné UI upozornění, že AI výstup je pouze pracovní podklad.

## Produkční roadmapa

- robustnější PDF/DOCX import včetně složitějšího formátování a skenovaných dokumentů
- vlastní vzory a knihovna klauzulí
- RAG nad interní znalostní bází
- historie analýz
- exporty do DOCX/PDF
- anonymizace a redakce citlivých údajů

## Upozornění

Demo nástroj. Výstupy AI slouží pouze jako pracovní podklad pro právníka a nenahrazují odborné právní posouzení.
