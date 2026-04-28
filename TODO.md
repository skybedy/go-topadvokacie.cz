# TODO.md

## Ted

- [ ] Udrzet demo spustitelne jednim prikazem `go run ./cmd/filipaipilot`.
- [ ] Pri kazde vetsi zmene aktualizovat `PROJECT_CONTEXT.md`, `TODO.md` a `DECISIONS.md`.
- [ ] Overovat zmeny pres `go test ./...`.
- [x] Zobrazit cekani na odpoved modelu primo v karte vysledku jako viditelny spinner, ne jen zmenu textu.
- [x] Pripravit kratky demo scenar pro Filipa: analyza smlouvy, red flags, priprava hovoru, porovnani dvou verzi.
- [x] Pridat tlacitko pro kopirovani jednotlivych vystupnich sekci.

## Dalsi kroky

- [x] Pridat volbu delky vystupu: strucne, standardne, detailne.
- [x] Pridat volbu perspektivy vystupu: pro pravnika, pro klienta, pro vyjednavani.
- [x] Pridat prompt "Vyjednavaci pozice" s mirnou, standardni a tvrdsi variantou pozadavku.
- [x] Pridat prompt "Red flags pro klienta" s maximalne peti zasadnimi body.
- [x] Pridat prompt "Priprav call s klientem" s agendou, riziky a otazkami.
- [x] Pridat prompt "Co ve smlouve chybi" pro identifikaci neupravenych nebo nejasnych oblasti.
- [x] Pridat prompt "Komentar do revize" pro navrhy komentaru k ustanovenim.
- [x] Pridat prompt "Executive summary pro jednatele".
- [x] Odstranit oddelenou umelou vrstvu workflow a pouzivat jen ulozene pravni prompty.
- [ ] Zlepsit ukazkove dokumenty tak, aby lepe demonstrovaly rozdil mezi pravnimi prompty.
- [ ] Doplnit zakladni unit testy pro `PromptTemplateByID`, upload validaci a `.env` parsing.
- [ ] Doladit demo scenar pro porovnani dvou verzi tak, aby umel predvyplnit i dokument B.

## Pozdeji

- [x] MVP PDF/DOCX import: DOCX pres ZIP/XML, PDF pres lokalni `pdftotext`.
- [x] Nacitat uploadovany soubor nejdriv do viditelneho textarea pole pred spustenim analyzy.
- [x] Pridat viditelny preloader pro upload souboru a cekani na odpoved modelu.
- [x] Pridat jednoduchy `deploy.sh` pro Ubuntu VPS.
- [ ] Robustnejsi PDF/DOCX import pro slozite dokumenty, tabulky a skenovana PDF.
- [ ] Export vystupu do DOCX/PDF.
- [ ] Historie analyz.
- [ ] Sprava promptu pres UI.
- [ ] Ulozeni promptu mimo Go kod, napr. verzovane soubory nebo databaze.
- [ ] Testovaci sada dokumentu pro ladeni kvality promptu.
- [ ] Mereni spotreby tokenu a ceny.
- [ ] Anonymizace nebo redakce citlivych udaju pred odeslanim do AI.
- [ ] Vlastni vzory a knihovna klauzuli.
- [ ] RAG nad internimi znalostmi nebo vzory.
- [ ] Role, opravneni a auditni logy pro produkcni nasazeni.
- [ ] Retencni politika pro vstupy a vystupy.
