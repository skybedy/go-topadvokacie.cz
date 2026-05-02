# PROJECT_CONTEXT.md

## Strucny popis projektu

FilipAiPilot je lokalni Go MVP pro komercniho pravnika Filipa. Cilem je ukazat rozdil mezi rucnim pouzivanim ChatGPT a strukturovanym pravnickym workflow nastrojem nad smlouvami a pravnimi texty.

Projekt neprodava predstavu "AI pravnika". Ukazuje opakovatelne pracovni postupy: pravnik vlozi dokument, vybere ulozeny pravni prompt a dostane strukturovany vystup jako pracovni podklad.

## Aktualni stav

- Aplikace bezi jako jednoduchy Go HTTP server.
- Frontend je server-side renderovane HTML s Tailwind CSS pres CDN.
- Drobne chovani ve frontendu je ve Vanilla JavaScriptu.
- Bez odpovidajiciho API klice aplikace bezi v mock demo rezimu s pripravenymi odpovedmi.
- S `AI_PROVIDER=openai` a `OPENAI_API_KEY` vola OpenAI Chat Completions API pres rucni klient v Go standardni knihovne.
- S `AI_PROVIDER=gemini` a `GEMINI_API_KEY` vola Gemini GenerateContent API pres rucni klient v Go standardni knihovne.
- Hlavni formular podporuje volbu delky vystupu: strucne, standardne, detailne.
- Hlavni formular podporuje perspektivu vystupu: pro pravnika, pro klienta, pro vyjednavani.
- Vystupni sekce maji tlacitko pro kopirovani textu do schranky.
- UI zobrazuje preloader pri nacitani uploadovaneho souboru i pri cekani na odpoved modelu.
- Cekaci stav analyzy se renderuje primo v karte "Vystup se zobrazi tady", aby uzivatel videl spinner ve stejnem miste, kde se po dokonceni ukaze vysledek.
- Homepage obsahuje kratky demo scenar pro Filipa s rychlymi odkazy na vhodne ulozene prompty.
- Upload souboru probiha pred analyzou pres endpoint `/upload-text`; vytazeny text se vlozi do viditelneho textarea pole a az potom uzivatel spousti analyzu.
- `.env.example` obsahuje `AI_PROVIDER`, `OPENAI_MODEL`, `GEMINI_MODEL` a sdileny timeout `AI_TIMEOUT_SECONDS=180`.
- Go modul se jmenuje `filipaipilot`.
- Projekt ma zakladni testy pro upload/parser v `internal/web/server_test.go`.
- Adresar je git repozitar na vetvi `main`.

## Prompt knihovna

- Kontrola smlouvy
- Srozumitelne pro klienta
- Otazky na klienta
- Navrh e-mailu bez odeslani
- Protiargumentace protistrany
- Checklist pred podpisem
- Revize obchodnich podminek
- Kontrola konzistence
- Extrakce povinnosti a lhut
- Red flags pred podpisem
- Vyjednavaci pozice
- Priprava hovoru s klientem
- Co ve smlouve chybi
- Komentare do revize
- Executive summary pro jednatele
- Porovnani dvou verzi

## Hlavni adresare a soubory

- `cmd/filipaipilot/main.go` - vstupni bod aplikace, nacitani `.env`, volba mock/OpenAI/Gemini klienta, spusteni serveru.
- `internal/ai/client.go` - AI interface, system prompt, prompt knihovna a volby vystupu.
- `internal/ai/mock.go` - mock AI klient s pripravenymi vystupy pro lokalni demo bez API klice.
- `internal/ai/openai.go` - OpenAI klient pres `net/http`, JSON response format a fallback pri nevalidnim JSON vystupu.
- `internal/ai/gemini.go` - Gemini klient pres `net/http`, GenerateContent API a fallback pri nevalidnim JSON vystupu.
- `internal/model/result.go` - datove struktury `Result`, `Section` a `Example`.
- `internal/web/server.go` - HTTP routing, handlery, upload/parser souboru, renderovani sablon a chybove hlasky.
- `internal/web/examples.go` - fiktivni ukazkove pravni texty.
- `templates/` - HTML sablony pro pracovni stul, prompty, ukazky a info stranku.
- `static/app.js` - prepinani druheho dokumentu, vyber ukazky, nacitani souboru do textarea, loading stav analyzy a kopirovani vystupnich sekci.
- `static/app.css` - drobne CSS upravy, hlavne vetsi pismo kvuli citelnosti a animace spinneru.
- `examples/` - fiktivni textove priklady.
- `README.md` - zakladni dokumentace projektu.
- `deploy.sh` - jednoduchy deploy script pro Ubuntu VPS: git pull, testy, build a restart systemd sluzby.
- `.env.example` - priklad konfigurace bez tajnych hodnot.

## Jak projekt spustit

```bash
go run ./cmd/filipaipilot
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
GEMINI_API_KEY=
GEMINI_MODEL=gemini-2.5-flash
AI_TIMEOUT_SECONDS=180
```

Kdyz chybi API klic odpovidajici zvolenemu provideru, aplikace bezi v mock rezimu.

## Jak projekt testovat

```bash
go test ./...
```

Aktualni stav k 2026-04-27:

```text
?   	filipaipilot/cmd/filipaipilot	[no test files]
?   	filipaipilot/internal/ai	[no test files]
?   	filipaipilot/internal/model	[no test files]
?   	filipaipilot/internal/web	[no test files]
```

Pro rychlou kontrolu spusteni:

```bash
go run ./cmd/filipaipilot
```

Pak otevrit `http://localhost:8080`.

## Znama omezeni a problemy

- Upload podporuje `.txt`, `.md`, `.markdown`, `.csv`, `.rst`, `.log`, `.docx` a `.pdf`.
- DOCX import je implementovan primo v Go pres ZIP/XML extrakci textu z `word/document.xml`, headeru a footeru.
- PDF import pouziva lokalni `pdftotext` z balicku `poppler-utils`, pokud je v systemu dostupny.
- Stary binarni `.doc` zatim neni podporovan; uzivatel ma dokument ulozit jako `.docx` nebo PDF.
- Limit uploadu je 5 MB, vytazeny text se pro demo zkracuje na 512 KB.
- Prompty jsou zapsane primo v Go kodu, ne v databazi ani v externich souborech.
- Historie analyz se neuklada.
- Export do DOCX/PDF neni implementovan.
- Neexistuje prihlaseni, role, opravneni ani auditni logy.
- Neni zde anonymizace citlivych udaju.
- OpenAI klient pouziva Chat Completions API a ocekava JSON objekt ve strukture `model.Result`.
- Gemini klient pouziva GenerateContent API a ocekava JSON objekt ve strukture `model.Result`.
- UI je demo dashboard, ne produkcni pravni aplikace.
- Pri zmene statickych assetu je vhodne menit query string u `app.js` a `app.css`, aby se v prohlizeci neudrzela stara cache a nebil se novy loading UI se starymi styly.

## Poznamky pro dalsi navazani

- Zachovej formulaci upozorneni: "Demo nastroj. Vystupy AI slouzi pouze jako pracovni podklad pro pravnika a nenahrazuji odborne pravni posouzeni."
- E-mailovy prompt smi pripravovat jen draft; nic se nema odesilat.
- Nejvetsi produktova hodnota dema je v ulozenych, pojmenovanych a verzovanych promptech, ne v samotnem API volani.
- Pro dalsi wow efekt dava smysl zlepsit ukazkove dokumenty, pridat export vystupu a doplnit male testy pro prompt knihovnu.
- Pred pridanim zavislosti vzdy zvaz, jestli se tim demo opravdu vyrazne zlepsi.
