# TODO.md

## Ted

- [ ] Udrzet demo spustitelne jednim prikazem `go run ./cmd/lexpilot`.
- [ ] Pri kazde vetsi zmene aktualizovat `PROJECT_CONTEXT.md`, `TODO.md` a `DECISIONS.md`.
- [ ] Overovat zmeny pres `go test ./...`.
- [ ] Pripravit kratky demo scenar pro Filipa: analyza smlouvy, otazky na klienta, srozumitelne shrnuti, porovnani dvou verzi.
- [ ] Zvazit tlacitko pro kopirovani vystupu nebo jednotlivych sekci.

## Dalsi kroky

- [ ] Pridat volbu delky vystupu: strucne, standardne, detailne.
- [ ] Pridat volbu perspektivy vystupu: pro pravnika, pro klienta, pro vyjednavani.
- [ ] Pridat prompt "Vyjednavaci pozice" s mirnou, standardni a tvrdsi variantou pozadavku.
- [ ] Pridat prompt "Red flags pro klienta" s maximalne peti zasadnimi body.
- [ ] Pridat prompt "Priprav call s klientem" s agendou, riziky a otazkami.
- [ ] Pridat prompt "Co ve smlouve chybi" pro identifikaci neupravenych nebo nejasnych oblasti.
- [ ] Pridat prompt "Komentar do revize" pro navrhy komentaru k ustanovenim.
- [ ] Pridat prompt "Executive summary pro jednatele".
- [ ] Zlepsit ukazkove dokumenty tak, aby lepe demonstrovaly rozdil mezi workflow.
- [ ] Doplnit zakladni unit testy pro `ActionByID`, `PromptTemplateByID`, upload validaci a `.env` parsing.

## Pozdeji

- [ ] PDF/DOCX import pomoci vhodnych parseru.
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
