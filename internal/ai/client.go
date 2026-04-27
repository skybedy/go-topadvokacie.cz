package ai

import (
	"context"

	"lexdemo/internal/model"
)

const SystemPrompt = "Jsi AI asistent pro českého komerčního právníka. Nejsi advokát. Pomáháš analyzovat dokumenty, strukturovat informace, upozorňovat na rizika a navrhovat otázky. Odpovídej česky, strukturovaně a opatrně. Nepřidávej fakta, která nejsou ve vstupu."

type Client interface {
	Analyze(ctx context.Context, action string, inputA string, inputB string, options Options) (model.Result, error)
}

type Options struct {
	DetailLevel string
	Perspective string
}

type Action struct {
	ID          string
	Label       string
	Description string
	NeedsSecond bool
}

type PromptTemplate struct {
	ID          string
	Label       string
	Version     string
	Category    string
	Description string
	Instruction string
	NeedsSecond bool
}

var Actions = []Action{
	{ID: "contract-analysis", Label: "Analýza smlouvy", Description: "Strukturované rozebrání smlouvy do klíčových právních a obchodních částí."},
	{ID: "client-summary", Label: "Shrnutí pro klienta", Description: "Krátké vysvětlení lidským jazykem bez právnického balastu."},
	{ID: "risk-points", Label: "Rizikové body", Description: "Seznam rizik včetně závažnosti, důvodu a návrhu řešení."},
	{ID: "change-proposal", Label: "Návrh změn", Description: "Doporučené úpravy textu, důvod a návrh formulace."},
	{ID: "client-questions", Label: "Otázky na klienta", Description: "Skutkové, obchodní, právní a důkazní otázky pro další práci."},
	{ID: "consistency-check", Label: "Kontrola konzistence", Description: "Rozpory, chybějící části, duplicity a nejasnosti."},
	{ID: "plain-language", Label: "Převod do srozumitelné řeči", Description: "Zachování významu právního textu v řeči pro laiky."},
	{ID: "compare-versions", Label: "Porovnání dvou verzí", Description: "Věcné, stylistické a rizikové změny mezi dokumenty.", NeedsSecond: true},
}

var PromptLibrary = []PromptTemplate{
	{
		ID:          "prompt-contract-review",
		Label:       "Kontrola smlouvy",
		Version:     "v1.0",
		Category:    "Smlouvy",
		Description: "Rychlá právní a obchodní revize smlouvy před dalším čtením.",
		Instruction: "Proveď kontrolu smlouvy. Zaměř se na předmět, strany, lhůty, cenu, odpovědnost, sankce, ukončení, nejasnosti a chybějící části. U každého bodu odliš, co plyne přímo ze vstupu a co je doporučené prověřit.",
	},
	{
		ID:          "prompt-plain-client",
		Label:       "Srozumitelně pro klienta",
		Version:     "v1.0",
		Category:    "Klientská komunikace",
		Description: "Převod právního textu do krátkého klientského vysvětlení.",
		Instruction: "Převeď právní text do srozumitelné řeči pro klienta. Zachovej význam, vynech právnický balast a upozorni na praktické dopady. Nepiš jako advokátní stanovisko.",
	},
	{
		ID:          "prompt-email-draft",
		Label:       "Návrh e-mailu bez odeslání",
		Version:     "v0.9",
		Category:    "Drafty",
		Description: "Návrh pracovního e-mailu, který se nikam neodesílá.",
		Instruction: "Připrav návrh e-mailu vycházející z vloženého textu. E-mail musí být věcný, profesionální a opatrný. Neodesílej jej, pouze navrhni text a přidej poznámky, co má právník před odesláním ověřit.",
	},
	{
		ID:          "prompt-counterparty",
		Label:       "Protiargumentace protistrany",
		Version:     "v1.0",
		Category:    "Vyjednávání",
		Description: "Simulace námitek, které může vznést druhá strana.",
		Instruction: "Vžij se do role protistrany a sepiš realistické protiargumenty k textu nebo návrhu. Poté navrhni věcné reakce, které může právník použít při vyjednávání. Nezveličuj a nepřidávej fakta mimo vstup.",
	},
	{
		ID:          "prompt-signing-checklist",
		Label:       "Checklist před podpisem",
		Version:     "v1.0",
		Category:    "Checklisty",
		Description: "Praktický seznam bodů, které je potřeba zkontrolovat před podpisem.",
		Instruction: "Vytvoř checklist před podpisem. Rozděl jej na právní, obchodní, procesní a důkazní body. U každého bodu napiš, proč je důležitý a co má právník nebo klient ověřit.",
	},
	{
		ID:          "prompt-terms-review",
		Label:       "Revize obchodních podmínek",
		Version:     "v1.0",
		Category:    "Obchodní podmínky",
		Description: "Kontrola obchodních podmínek z pohledu B2B/B2C rizik.",
		Instruction: "Zkontroluj obchodní podmínky. Zaměř se na identifikaci poskytovatele, objednávku, platby, reklamace, odpovědnost, odstoupení, změny podmínek, ochranu spotřebitele a nejasné formulace.",
	},
	{
		ID:          "prompt-obligations-deadlines",
		Label:       "Extrakce povinností a lhůt",
		Version:     "v1.0",
		Category:    "Extrakce",
		Description: "Vytahuje kdo má co udělat, do kdy a co hrozí při nesplnění.",
		Instruction: "Extrahuj z dokumentu povinnosti a lhůty. U každé položky uveď povinnou stranu, obsah povinnosti, lhůtu nebo spouštěcí událost, sankci nebo následek a míru nejasnosti.",
	},
	{
		ID:          "prompt-red-flags",
		Label:       "Red flags před podpisem",
		Version:     "v1.0",
		Category:    "Rizika",
		Description: "Krátký seznam nejzásadnějších bodů, které má klient vidět před rozhodnutím.",
		Instruction: "Najdi maximálně 5 nejdůležitějších red flags před podpisem. U každého bodu uveď, proč je problém prakticky důležitý, jaký může mít obchodní dopad a jaký další krok doporučuješ právníkovi ověřit.",
	},
	{
		ID:          "prompt-negotiation-position",
		Label:       "Vyjednávací pozice",
		Version:     "v1.0",
		Category:    "Vyjednávání",
		Description: "Připravuje mírnou, standardní a tvrdší variantu požadavku ke klíčovým ustanovením.",
		Instruction: "Vyber ustanovení vhodná k vyjednávání. U každého bodu navrhni mírnou variantu požadavku, standardní kompromis a tvrdší vyjednávací pozici. Přidej stručný argument pro klienta a realistickou reakci protistrany.",
	},
	{
		ID:          "prompt-client-call",
		Label:       "Příprava hovoru s klientem",
		Version:     "v1.0",
		Category:    "Klientská práce",
		Description: "Agenda hovoru, otázky, rozhodnutí a podklady, které má klient dodat.",
		Instruction: "Připrav 15minutový hovor s klientem. Vytvoř agendu, otázky na klienta, body k obchodnímu rozhodnutí, podklady k doplnění a doporučený další krok po hovoru.",
	},
	{
		ID:          "prompt-missing-clauses",
		Label:       "Co ve smlouvě chybí",
		Version:     "v1.0",
		Category:    "Kontrola",
		Description: "Upozorní na oblasti, které ve vstupu nejsou řešené nebo jsou nejasné.",
		Instruction: "Zkontroluj, které obvyklé oblasti nejsou ve vstupu řešené nebo jsou nejasné. Nevyvozuj, že jde vždy o chybu. U každé položky napiš, proč může být důležitá a jakou otázku má právník položit klientovi.",
	},
	{
		ID:          "prompt-review-comments",
		Label:       "Komentáře do revize",
		Version:     "v1.0",
		Category:    "Revize",
		Description: "Navrhuje stručné komentáře k ustanovením pro práci v revizním režimu.",
		Instruction: "Navrhni komentáře do revize smlouvy. Každý komentář napiš věcně a použitelně pro právníka: co upravit, proč, jaké rozhodnutí je potřeba a co případně ověřit u klienta.",
	},
	{
		ID:          "prompt-executive-summary",
		Label:       "Executive summary pro jednatele",
		Version:     "v1.0",
		Category:    "Rozhodování",
		Description: "Manažerské shrnutí pro rychlé obchodní rozhodnutí.",
		Instruction: "Připrav executive summary pro jednatele nebo vedení. Stručně uveď, zda dokument podepsat, nepodepsat, nebo podepsat po úpravách. Přidej tři hlavní rizika, obchodní dopad a rozhodnutí, která musí udělat klient.",
	},
}

func ActionByID(id string) Action {
	for _, action := range Actions {
		if action.ID == id {
			return action
		}
	}
	for _, prompt := range PromptLibrary {
		if prompt.ID == id {
			return Action{
				ID:          prompt.ID,
				Label:       prompt.Label,
				Description: prompt.Description,
				NeedsSecond: prompt.NeedsSecond,
			}
		}
	}
	return Actions[0]
}

func PromptTemplateByID(id string) (PromptTemplate, bool) {
	for _, prompt := range PromptLibrary {
		if prompt.ID == id {
			return prompt, true
		}
	}
	return PromptTemplate{}, false
}
