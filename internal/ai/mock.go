package ai

import (
	"context"
	"strings"

	"filipaipilot/internal/model"
)

type MockAIClient struct{}

func NewMockAIClient() *MockAIClient {
	return &MockAIClient{}
}

func (c *MockAIClient) Analyze(ctx context.Context, action string, inputA string, inputB string, options Options) (model.Result, error) {
	select {
	case <-ctx.Done():
		return model.Result{}, ctx.Err()
	default:
	}

	switch action {
	case "prompt-contract-review":
		return result("Prompt knihovna: Kontrola smlouvy", "Tento prompt simuluje Filipovu uloženou šablonu pro rychlou kontrolu smlouvy. Výstup odděluje fakta z dokumentu od bodů k ověření.", []model.Section{
			{"Co prompt kontroluje", []string{"Určitost předmětu plnění a identifikaci stran.", "Lhůty, cenu, odpovědnost, sankce a ukončení.", "Chybějící přílohy, definice a proces předání."}},
			{"Ukázkový výstup", []string{"Doporučuji ověřit, zda dokument obsahuje přesný akceptační proces.", "Pokud odpovědnost není limitovaná, jde o bod k obchodnímu rozhodnutí.", "Sankce je vhodné navázat na konkrétní porušení."}},
		}), nil
	case "prompt-plain-client":
		return result("Prompt knihovna: Srozumitelně pro klienta", "Prompt převede právní text do krátkého klientského vysvětlení bez zbytečného právnického stylu.", []model.Section{
			{"Klientské vysvětlení", []string{"Text říká, kdo má co dodat, kdy a co se stane při porušení.", "Nejdůležitější je pohlídat si sankce, ukončení a praktické důsledky pro běžný provoz."}},
			{"Poznámka pro právníka", []string{"Před odesláním klientovi zkontrolujte, zda zjednodušení nezměnilo právní význam."}},
		}), nil
	case "prompt-client-questions":
		return result("Prompt knihovna: Otázky na klienta", "Tyto otázky pomohou doplnit skutkový a obchodní kontext před úpravou dokumentu.", []model.Section{
			{"Skutkové", []string{"Kdo bude za klienta přebírat plnění?", "Existují přílohy nebo technická specifikace, na které text odkazuje?"}},
			{"Obchodní", []string{"Jaká je maximální přijatelná výše smluvní pokuty?", "Je pro klienta důležitější rychlé ukončení, nebo stabilita vztahu?"}},
			{"Právní", []string{"Má být odpovědnost limitována?", "Má smlouva řešit mlčenlivost i po skončení vztahu?"}},
			{"Důkazní", []string{"Jak se bude prokazovat předání, vady a reklamace?", "Kde budou ukládány objednávky a akceptace?"}},
		}), nil
	case "prompt-email-draft":
		return result("Prompt knihovna: Návrh e-mailu bez odeslání", "Prompt připraví pouze pracovní návrh e-mailu. Nic se automaticky neodesílá.", []model.Section{
			{"Draft e-mailu", []string{"Dobrý den, zasílám několik bodů k doplnění smlouvy. Prosím zejména o potvrzení lhůt, odpovědnosti a sankcí.", "Po doplnění podkladů připravím upravené znění."}},
			{"Ověřit před použitím", []string{"Adresáta, tón komunikace, obchodní kontext a případné citlivé informace."}},
		}), nil
	case "prompt-counterparty":
		return result("Prompt knihovna: Protiargumentace protistrany", "Prompt pomáhá připravit se na vyjednávání tím, že předvídá realistické námitky druhé strany.", []model.Section{
			{"Možné námitky", []string{"Limit odpovědnosti je pro protistranu příliš nízký.", "Smluvní pokuta může být vnímána jako nepřiměřená.", "Výpovědní doba nemusí odpovídat investici do spolupráce."}},
			{"Možné reakce", []string{"Navrhnout odstupňování podle závažnosti porušení.", "Spojit vyšší odpovědnost s konkrétními typy škody.", "Doplnit přechodné období při ukončení."}},
		}), nil
	case "prompt-signing-checklist":
		return result("Prompt knihovna: Checklist před podpisem", "Prompt vytvoří praktický seznam věcí, které má právník nebo klient projít před podpisem.", []model.Section{
			{"Právní", []string{"Jsou správně označeny smluvní strany?", "Je jasné, kdy lze odstoupit nebo vypovědět smlouvu?", "Je vyřešen vztah smluvní pokuty a náhrady škody?"}},
			{"Obchodní", []string{"Odpovídají platby a lhůty domluvenému dealu?", "Je klient schopen povinnosti reálně splnit?"}},
			{"Důkazní", []string{"Je jasné, jak se bude prokazovat předání, vady a komunikace?"}},
		}), nil
	case "prompt-terms-review":
		return result("Prompt knihovna: Revize obchodních podmínek", "Prompt prochází obchodní podmínky pohledem typických B2B/B2C rizik.", []model.Section{
			{"Kontrolní body", []string{"Identifikace poskytovatele, objednávkový proces a platební podmínky.", "Reklamace, odpovědnost, odstoupení a jednostranné změny podmínek.", "Spotřebitelská ustanovení, pokud jde o B2C vztah."}},
			{"Typická rizika", []string{"Nejasné změny podmínek bez oznámení.", "Chybějící reklamační proces.", "Příliš široké omezení odpovědnosti."}},
		}), nil
	case "prompt-consistency-check":
		return result("Prompt knihovna: Kontrola konzistence", "Text působí jako dokument, u kterého je vhodné ověřit návaznost definic, lhůt a sankcí.", []model.Section{
			{"Rozpory", []string{"Zkontrolujte, zda se stejná strana neoznačuje různými názvy.", "Ověřte, zda výpovědní doba odpovídá ustanovením o ukončení."}},
			{"Chybějící části", []string{"Může chybět proces předání plnění.", "U smluvní pokuty často chybí vztah k náhradě škody."}},
			{"Duplicity", []string{"Prověřte, zda se povinnost mlčenlivosti neopakuje s rozdílným rozsahem."}},
		}), nil
	case "prompt-obligations-deadlines":
		return result("Prompt knihovna: Extrakce povinností a lhůt", "Prompt vytáhne z dokumentu praktickou tabulku kdo má co udělat, do kdy a co hrozí při nesplnění.", []model.Section{
			{"Povinnosti", []string{"Dodavatel: dodat plnění podle smlouvy do sjednaného termínu.", "Objednatel: zaplatit cenu ve lhůtě splatnosti.", "Obě strany: zachovat mlčenlivost, pokud je sjednána."}},
			{"Lhůty a následky", []string{"Prodlení s dodáním může spustit smluvní pokutu.", "Prodlení s platbou může vést k úroku, výpovědi nebo pozastavení plnění."}},
		}), nil
	case "prompt-red-flags":
		return resultWithOptions("Prompt knihovna: Red flags před podpisem", "Krátký výstup pro rychlé rozhodnutí klienta. V reálném režimu by AI vybrala jen nejzásadnější rizika z vloženého textu.", options, []model.Section{
			{"Top red flags", []string{"Neomezená odpovědnost může vytvořit nepřiměřené ekonomické riziko.", "Smluvní pokuta bez jasného stropu může být obchodně tvrdá.", "Chybějící akceptační proces komplikuje dokazování, zda bylo plnění řádně předáno.", "Jednostranná změna podmínek bez možnosti ukončení je bod k vyjednávání.", "Nejasné ukončení smlouvy může klienta svázat déle, než čeká."}},
			{"Doporučený další krok", []string{"Ověřit s klientem, která rizika jsou obchodně přijatelná a která je nutné upravit před podpisem."}},
		}), nil
	case "prompt-negotiation-position":
		return resultWithOptions("Prompt knihovna: Vyjednávací pozice", "Prompt převádí právní rizika do konkrétních vyjednávacích požadavků.", options, []model.Section{
			{"Odpovědnost", []string{"Mírná varianta: doplnit rozumný limit odpovědnosti.", "Standardní kompromis: limit na výši odměny za posledních 12 měsíců.", "Tvrdší pozice: vyloučit nepřímé škody a limitovat odpovědnost kromě úmyslu."}},
			{"Smluvní pokuta", []string{"Mírná varianta: snížit sazbu nebo zavést celkový strop.", "Standardní kompromis: navázat pokutu na závažná porušení.", "Tvrdší pozice: požadovat odstranění pokuty a ponechat jen náhradu škody."}},
			{"Argumentace", []string{"Pro klienta: cílem není oslabit smlouvu, ale nastavit předvídatelné riziko.", "Možná reakce protistrany: bez sankce nebude mít závazek dostatečnou váhu."}},
		}), nil
	case "prompt-client-call":
		return resultWithOptions("Prompt knihovna: Příprava hovoru s klientem", "Prompt připraví právníka na rychlý a věcný call nad dokumentem.", options, []model.Section{
			{"Agenda 15 minut", []string{"1. Potvrdit obchodní cíl smlouvy.", "2. Projít tři hlavní rizika.", "3. Rozhodnout, které body vyjednávat.", "4. Domluvit další podklady a termín revize."}},
			{"Otázky na klienta", []string{"Jaká je maximální přijatelná odpovědnost?", "Je důležitější rychlé ukončení, nebo stabilita vztahu?", "Kdo bude prakticky potvrzovat předání plnění?"}},
			{"Po hovoru", []string{"Připravit revizní komentáře a krátké klientské shrnutí bodů k rozhodnutí."}},
		}), nil
	case "prompt-missing-clauses":
		return resultWithOptions("Prompt knihovna: Co ve smlouvě chybí", "Prompt nehledá jen chyby, ale hlavně oblasti, které je potřeba vědomě potvrdit.", options, []model.Section{
			{"Chybějící nebo nejasné oblasti", []string{"Akceptační proces: není jasné, jak se potvrdí předání plnění.", "Limit odpovědnosti: pokud chybí, klient potřebuje znát maximální expozici.", "Rozhodné právo a řešení sporů: vhodné ověřit, zda odpovídá dohodě.", "Kontaktní osoby a doručování: prakticky důležité pro výzvy a reklamace."}},
			{"Otázky pro klienta", []string{"Bylo vynechání těchto oblastí záměrné?", "Existuje příloha nebo objednávka, která tyto body řeší mimo hlavní smlouvu?"}},
		}), nil
	case "prompt-review-comments":
		return resultWithOptions("Prompt knihovna: Komentáře do revize", "Prompt připraví pracovní komentáře, které právník může upravit a vložit do revizního režimu.", options, []model.Section{
			{"Komentáře", []string{"Doporučuji doplnit přesný akceptační proces, protože bez něj může být sporné, kdy bylo plnění převzato.", "Prosím potvrdit, zda je výše smluvní pokuty obchodně přijatelná; aktuální formulace může být pro klienta tvrdá.", "Navrhuji doplnit limit odpovědnosti, aby bylo riziko předvídatelné."}},
			{"Ověřit u klienta", []string{"Zda má klient vyjednávací prostor u odpovědnosti, sankcí a ukončení."}},
		}), nil
	case "prompt-executive-summary":
		return resultWithOptions("Prompt knihovna: Executive summary pro jednatele", "Manažerský výstup převádí právní rozbor do rozhodnutí, ne do dlouhého stanoviska.", options, []model.Section{
			{"Doporučení", []string{"Podepsat až po úpravách odpovědnosti, sankcí a akceptačního procesu."}},
			{"Tři hlavní rizika", []string{"Neomezená odpovědnost může vytvořit nepřiměřené finanční riziko.", "Smluvní pokuta může být vysoká vzhledem k významu porušení.", "Nejasné předání plnění může vést ke sporu o splnění."}},
			{"Rozhodnutí vedení", []string{"Stanovit maximální akceptovatelný limit odpovědnosti.", "Rozhodnout, zda je obchodně nutné trvat na snížení sankcí.", "Potvrdit, kdo bude za klienta přebírat plnění."}},
		}), nil
	case "prompt-compare-versions":
		return result("Prompt knihovna: Porovnání dvou verzí", "Druhá verze podle demo výstupu mění zejména rozložení rizika a zpřesňuje některé procesní kroky.", []model.Section{
			{"Věcné změny", []string{"Zkontrolujte, zda druhá verze nemění rozsah odpovědnosti nebo sankcí.", "Ověřte dopady nových lhůt na provoz klienta."}},
			{"Stylistické změny", []string{"Některé formulace mohou být kratší, ale právně méně přesné.", "Sjednoťte terminologii mezi oběma verzemi."}},
			{"Rizikové změny", []string{"Pozor na vypuštění limitu odpovědnosti nebo práva odstoupit.", "Pozor na nenápadné rozšíření mlčenlivosti či zákazu konkurence."}},
		}), nil
	default:
		return result("Analýza smlouvy", "Dokument je potřeba číst jako pracovní podklad. Demo výstup zvýrazňuje hlavní části, které by právník typicky kontroloval.", []model.Section{
			{"Shrnutí", []string{"Text upravuje obchodní vztah a stanoví základní povinnosti stran.", "Pro další práci je vhodné ověřit ekonomické parametry a praktickou vymahatelnost."}},
			{"Smluvní strany", []string{"Identifikace stran by měla obsahovat název, IČO, sídlo a oprávněné osoby.", "Zkontrolujte, zda jsou strany označeny konzistentně v celém dokumentu."}},
			{"Předmět smlouvy", []string{"Předmět musí být dostatečně určitý.", "Pokud dokument odkazuje na přílohy, měly by být přiložené a verzované."}},
			{"Lhůty", []string{"Lhůty je vhodné navázat na konkrétní události.", "Pozor na neurčité výrazy typu bez zbytečného odkladu bez dalšího kontextu."}},
			{"Platby", []string{"Ověřte cenu, splatnost, DPH, fakturaci a následky prodlení."}},
			{"Povinnosti", []string{"Povinnosti by měly být měřitelné a přiřazené konkrétní straně."}},
			{"Sankce", []string{"Sankce by měly odpovídat závažnosti porušení.", "Prověřte vztah smluvní pokuty a náhrady škody."}},
			{"Ukončení", []string{"Zkontrolujte výpověď, odstoupení, následky ukončení a vypořádání."}},
			{"Nejasnosti", []string{"Neurčité pojmy je vhodné definovat nebo nahradit přesnější formulací."}},
			{"Doporučení", []string{"Doplnit chybějící definice, limity odpovědnosti a jasný akceptační postup."}},
		}), nil
	}
}

func result(title, summary string, sections []model.Section) model.Result {
	return model.Result{
		Title:    title,
		Summary:  summary,
		Sections: sections,
		Warnings: []string{
			"Demo režim používá předpřipravené odpovědi a nehodnotí plný právní kontext.",
			"Výstup je pracovní podklad pro právníka, nikoli právní stanovisko.",
		},
		Raw: rawPreview(summary),
	}
}

func resultWithOptions(title, summary string, options Options, sections []model.Section) model.Result {
	sections = append([]model.Section{
		{"Nastavení výstupu", []string{
			"Délka: " + mockDetailLabel(options.DetailLevel),
			"Perspektiva: " + mockPerspectiveLabel(options.Perspective),
		}},
	}, sections...)
	return result(title, summary, sections)
}

func mockDetailLabel(value string) string {
	switch value {
	case "brief":
		return "stručně"
	case "detailed":
		return "detailně"
	default:
		return "standardně"
	}
}

func mockPerspectiveLabel(value string) string {
	switch value {
	case "client":
		return "pro klienta"
	case "negotiation":
		return "pro vyjednávání"
	default:
		return "pro právníka"
	}
}

func rawPreview(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	return "Mock fallback: " + s
}
