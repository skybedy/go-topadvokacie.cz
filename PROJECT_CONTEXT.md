# PROJECT_CONTEXT.md

## Strucny popis projektu

LexPilot Demo je lokalni Go MVP pro komercniho pravnika Filipa. Cilem je ukazat rozdil mezi rucnim pouzivanim ChatGPT a strukturovanym pravnickym workflow nastrojem nad smlouvami a pravnimi texty.

Projekt neprodava predstavu "AI pravnika". Ukazuje opakovatelne pracovni postupy: pravnik vlozi dokument, vybere workflow nebo ulozeny prompt a dostane strukturovany vystup jako pracovni podklad.

## Aktualni stav

- Aplikace bezi jako jednoduchy Go HTTP server.
- Frontend je server-side renderovane HTML s Tailwind CSS pres CDN.
- Drobne chovani ve frontendu je ve Vanilla JavaScriptu.
- Bez `OPENAI_API_KEY` aplikace bezi v mock demo rezimu s pripravenymi odpovedmi.
- S `OPENAI_API_KEY` vola OpenAI Chat Completions API pres rucni klient v Go standardni knihovne.
- Hlavni formular podporuje volbu delky vystupu: strucne, standardne, detailne.
- Hlavni formular podporuje perspektivu vystupu: pro pravnika, pro klienta, pro vyjednavani.
- Vystupni sekce maji tlacitko pro kopirovani textu do schranky.
- Homepage obsahuje kratky demo scenar pro Filipa s rychlymi odkazy na vhodne workflow.
- Upload souboru probiha pred analyzou pres endpoint `/upload-text`; vytazeny text se vlozi do viditelneho textarea pole a az potom uzivatel spousti analyzu.
- `.env.example` nastavuje `OPENAI_MODEL=gpt-5-nano` a `OPENAI_TIMEOUT_SECONDS=180`.
- Go modul se jmenuje `lexdemo`.
- Projekt aktualne nema testovaci soubory, ale `go test ./...` prochazi.
- Adresar momentalne neni inicializovany jako git repozitar, proto `git status` selze s hlaskou `not a git repository`.

## Hlavni workflow

Vestavene workflow:

- Analyza smlouvy
- Shrnuti pro klienta
- Rizikove body
- Navrh zmen
- Otazky na klienta
- Kontrola konzistence
- Prevod pravniho textu do srozumitelne reci
- Porovnani dvou verzi dokumentu

Prompt knihovna:

- Kontrola smlouvy
- Srozumitelne pro klienta
- Navrh e-mailu bez odeslani
- Protiargumentace protistrany
- Checklist pred podpisem
- Revize obchodnich podminek
- Extrakce povinnosti a lhut
- Red flags pred podpisem
- Vyjednavaci pozice
- Priprava hovoru s klientem
- Co ve smlouve chybi
- Komentare do revize
- Executive summary pro jednatele

## Hlavni adresare a soubory

- `cmd/lexpilot/main.go` - vstupni bod aplikace, nacitani `.env`, volba mock/OpenAI klienta, spusteni serveru.
- `internal/ai/client.go` - AI interface, system prompt, seznam workflow, prompt knihovna a volby vystupu.
- `internal/ai/mock.go` - mock AI klient s pripravenymi vystupy pro lokalni demo bez API klice.
- `internal/ai/openai.go` - OpenAI klient pres `net/http`, JSON response format a fallback pri nevalidnim JSON vystupu.
- `internal/model/result.go` - datove struktury `Result`, `Section` a `Example`.
- `internal/web/server.go` - HTTP routing, handlery, upload/parser souboru, renderovani sablon a chybove hlasky.
- `internal/web/examples.go` - fiktivni ukazkove pravni texty.
- `templates/` - HTML sablony pro pracovni stul, prompty, ukazky a info stranku.
- `static/app.js` - prepinani druheho dokumentu, vyber ukazky, nacitani souboru do textarea a kopirovani vystupnich sekci.
- `static/app.css` - drobne CSS upravy, hlavne vetsi pismo kvuli citelnosti.
- `examples/` - fiktivni textove priklady.
- `README.md` - zakladni dokumentace projektu.
- `.env.example` - priklad konfigurace bez tajnych hodnot.

## Jak projekt spustit

```bash
go run ./cmd/lexpilot
```

Vychozi URL:

```text
http://localhost:8080
```

Volitelna konfigurace:

```bash
cp .env.example .env
```

Zakladni promenne:

```env
ADDR=:8080
OPENAI_API_KEY=
OPENAI_MODEL=gpt-5-nano
OPENAI_TIMEOUT_SECONDS=180
```

Kdyz `OPENAI_API_KEY` chybi, aplikace bezi v mock rezimu.

## Jak projekt testovat

```bash
go test ./...
```

Aktualni stav k 2026-04-27:

```text
?   	lexdemo/cmd/lexpilot	[no test files]
?   	lexdemo/internal/ai	[no test files]
?   	lexdemo/internal/model	[no test files]
?   	lexdemo/internal/web	[no test files]
```

Pro rychlou kontrolu spusteni:

```bash
go run ./cmd/lexpilot
```

Pak otevrit `http://localhost:8080`.

## Znama omezeni a problemy

- Upload podporuje `.txt`, `.md`, `.markdown`, `.csv`, `.rst`, `.log`, `.docx` a `.pdf`.
- DOCX import je implementovan primo v Go pres ZIP/XML extrakci textu z `word/document.xml`, headeru a footeru.
- PDF import pouziva lokalni `pdftotext` z balicku `poppler-utils`, pokud je v systemu dostupny.
- Stary binarni `.doc` zatim neni podporovan; uzivatel ma dokument ulozit jako `.docx` nebo PDF.
- Limit uploadu je 5 MB, vytazeny text se pro demo zkracuje na 512 KB.
- Prompty a workflow jsou zapsane primo v Go kodu, ne v databazi ani v externich souborech.
- Historie analyz se neuklada.
- Export do DOCX/PDF neni implementovan.
- Neexistuje prihlaseni, role, opravneni ani auditni logy.
- Neni zde anonymizace citlivych udaju.
- OpenAI klient pouziva Chat Completions API a ocekava JSON objekt ve strukture `model.Result`.
- UI je demo dashboard, ne produkcni pravni aplikace.

## Poznamky pro dalsi navazani

- Zachovej formulaci upozorneni: "Demo nastroj. Vystupy AI slouzi pouze jako pracovni podklad pro pravnika a nenahrazuji odborne pravni posouzeni."
- E-mailove workflow smi pripravovat jen draft; nic se nema odesilat.
- Nejvetsi produktova hodnota dema je v ulozenych, pojmenovanych a verzovanych workflow, ne v samotnem API volani.
- Pro dalsi wow efekt dava smysl zlepsit ukazkove dokumenty, pridat export vystupu a doplnit male testy pro prompt knihovnu.
- Pred pridanim zavislosti vzdy zvaz, jestli se tim demo opravdu vyrazne zlepsi.
