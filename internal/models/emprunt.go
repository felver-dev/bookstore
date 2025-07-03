package models

import (
	"fmt"
	"strings"
	"time"
)

type Emprunt struct {
	ID                 int        `json:"id"`
	LivreID            int        `json:"livre_id"`
	MembreID           int        `json:"membre_id"`
	DateEmprunt        time.Time  `json:"date_emprunt"`
	DateRetourPrevu    time.Time  `json:"date_retour_prevu"`
	DateRetourEffectif *time.Time `json:"date_retour_effectif"`
	Statut             string     `json:"statut"`

	TitreLivre string `json:"titre_livre"`
	NomMembre  string `json:"nom_membre"`
}

const (
	DUREE_EMPRUMT_JOURS = 14
	STATUT_EN_COURS     = "en-cours"
	STATUT_RENDU        = "rendu"
	STATUT_EN_RETARD    = "en-retard"
)

func (e Emprunt) String() string {

	statutEmoji := "📘"

	switch e.Statut {
	case STATUT_EN_COURS:
		statutEmoji = "📘 En cours"
	case STATUT_RENDU:
		statutEmoji = "📘 Rendu"
	case STATUT_EN_RETARD:
		statutEmoji = "⚠️ En retard"
	}

	return fmt.Sprintf("ID: %d | %s par %s | Emprunté le %s | %s",
		e.ID, e.TitreLivre, e.NomMembre,
		e.DateEmprunt.Format("02/01/2006"), statutEmoji)

}

// AfficherDetails() montre toutes les informations d'un emprunt
func (e Emprunt) AfficherDetails() {
	fmt.Printf("┌%s┐\n", strings.Repeat("─", 70))
	fmt.Printf("│ Emprunt #%d%s│\n", e.ID, strings.Repeat(" ", 70-len(fmt.Sprintf(" Emprunt #%d", e.ID))))
	fmt.Printf("├%s┤\n", strings.Repeat("─", 70))
	fmt.Printf("│ Livre         : %-50s │\n", e.TitreLivre)
	fmt.Printf("│ Membre        : %-50s │\n", e.NomMembre)
	fmt.Printf("│ Emprunté le   : %-50s │\n", e.DateEmprunt.Format("02/01/2006 15:04:05"))
	fmt.Printf("│ À rendre le   : %-50s │\n", e.DateRetourPrevu.Format("02/01/2006"))

	// Affichage conditionnel de la date de retour effectif
	if e.DateRetourEffectif != nil {
		fmt.Printf("│ Rendu le      : %-50s │\n", e.DateRetourEffectif.Format("02/01/2006 15:04:05"))
	} else {
		fmt.Printf("│ Rendu le      : %-50s │\n", "Pas encore rendu")
	}

	// Afficher le statut avec des emojis
	statutAffichage := ""
	switch e.Statut {
	case STATUT_EN_COURS:
		statutAffichage = "📘 En cours"
	case STATUT_RENDU:
		statutAffichage = "✅ Rendu"
	case STATUT_EN_RETARD:
		statutAffichage = "⚠️ En retard"
	}
	fmt.Printf("│ Statut        : %-50s │\n", statutAffichage)

	// Calculer et afficher les jours de retard s'il y en a
	if e.Statut == STATUT_EN_RETARD {
		joursRetard := int(time.Since(e.DateRetourPrevu).Hours() / 24)
		fmt.Printf("│ Retard        : %-50s │\n", fmt.Sprintf("%d jour(s)", joursRetard))
	}

	fmt.Printf("└%s┘\n", strings.Repeat("─", 70))
}

func (e Emprunt) EstEnRetard() bool {
	return e.DateRetourEffectif == nil && time.Now().After(e.DateRetourPrevu)
}

func (e Emprunt) CalculerJoursRetard() int {
	if !e.EstEnRetard() {
		return 0
	}

	duree := time.Since(e.DateRetourPrevu)
	return int(duree.Hours() / 24)
}

func (e *Emprunt) MarquerCommeRendu() {

	maintenant := time.Now()
	e.DateRetourEffectif = &maintenant
	e.Statut = STATUT_RENDU

}
func (e *Emprunt) MettreAjourStatut() {
	if e.DateRetourEffectif == nil {
		if time.Now().After(e.DateRetourPrevu) {
			e.Statut = STATUT_EN_RETARD
		} else {
			e.Statut = STATUT_EN_COURS
		}
	} else {
		e.Statut = STATUT_RENDU
	}
}
