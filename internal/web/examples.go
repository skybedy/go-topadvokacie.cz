package web

import "lexdemo/internal/model"

var Examples = []model.Example{
	{
		ID:      "nda",
		Title:   "NDA",
		Summary: "Fiktivní dohoda o mlčenlivosti mezi technologickou firmou a externím konzultantem.",
		Content: `DOHODA O MLČENLIVOSTI

Společnost AlfaTech s.r.o. a konzultant Jan Novák se dohodli, že konzultant bude zachovávat mlčenlivost o obchodních, technických a finančních informacích, které získá při jednání o možné spolupráci.

Mlčenlivost trvá po dobu 5 let od podpisu dohody. Konzultant nesmí informace zpřístupnit třetím osobám bez předchozího písemného souhlasu AlfaTech s.r.o.

Za porušení mlčenlivosti je sjednána smluvní pokuta 500 000 Kč za každé jednotlivé porušení. Tím není dotčeno právo na náhradu škody v plné výši.`,
	},
	{
		ID:      "najem",
		Title:   "Nájemní smlouva",
		Summary: "Krátká ukázka nájmu kancelářských prostor pro obchodní účely.",
		Content: `NÁJEMNÍ SMLOUVA

Pronajímatel přenechává nájemci do užívání kancelář č. 304 v budově Business Park Praha. Nájemné činí 28 000 Kč měsíčně bez DPH a je splatné vždy do 10. dne příslušného měsíce.

Nájem se sjednává na dobu neurčitou s výpovědní dobou 3 měsíce. Nájemce je povinen užívat prostor pouze k administrativním účelům a nesmí jej přenechat třetí osobě bez souhlasu pronajímatele.

Při prodlení s úhradou nájemného delším než 15 dnů je pronajímatel oprávněn smlouvu vypovědět s okamžitou účinností.`,
	},
	{
		ID:      "pokuta",
		Title:   "Smluvní pokuta",
		Summary: "Samostatné ustanovení o smluvní pokutě pro kontrolu přiměřenosti a jasnosti.",
		Content: `SMLUVNÍ POKUTA

V případě, že dodavatel nedodá dílo v termínu uvedeném ve smlouvě, je objednatel oprávněn požadovat smluvní pokutu ve výši 0,5 % z celkové ceny díla za každý den prodlení.

Smluvní pokuta je splatná do 7 dnů od doručení písemné výzvy objednatele. Zaplacením smluvní pokuty není dotčeno právo objednatele na náhradu škody.`,
	},
	{
		ID:      "dodatek",
		Title:   "Dodatek",
		Summary: "Fiktivní dodatek měnící termín dodání a cenu projektu.",
		Content: `DODATEK Č. 1

Smluvní strany se dohodly, že termín předání první verze díla se prodlužuje z 30. 6. 2026 na 31. 7. 2026.

Cena díla se zvyšuje o 120 000 Kč bez DPH z důvodu rozšíření rozsahu prací o integrační modul. Ostatní ustanovení smlouvy zůstávají beze změny.`,
	},
}

func ExampleByID(id string) (model.Example, bool) {
	for _, example := range Examples {
		if example.ID == id {
			return example, true
		}
	}
	return model.Example{}, false
}
