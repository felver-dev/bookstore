package models

import (
	"fmt"
	"strings"
	"time"
)

type Livre struct {
	ID              int       `json:"id"`
	Titre           string    `json:"titre"`
	Auteur          string    `json:"auteur"`
	ISBN            string    `json:"isbn"`
	Genre           string    `json:"genre"`
	DatePublication time.Time `json:"date_publication"`
	Disponible      bool      `json:"disponible"`
	NombreEmprunts  int       `json:"nombre_emprunts"`
	DateAjout       time.Time `json:"date_ajout"`
}

// Permet d'afficher un livre de manière simple
// Elle est appelée automatiquement quanf on fait fmt.Print(livre)
func (l Livre) String() string {
	statut := "📗 Disponible"
	if !l.Disponible {
		statut = "📕 Emprunté"
	}

	return fmt.Sprintf("ID: %d | %s par %s | %s | %s ", l.ID, l.Titre, l.Auteur, l.Genre, statut)
}

// AfficherDetails() montre toutes les informations d'un livre dans un tableau
func (l Livre) AfficherDetails() {
	fmt.Printf("┌%s┐\n", strings.Repeat("─", 60))
	fmt.Printf("│ Livre #%d%s│\n", l.ID, strings.Repeat(" ", 60-len(fmt.Sprintf(" Livre #%d", l.ID))))
	fmt.Printf("├%s┤\n", strings.Repeat("─", 60))
	fmt.Printf("│ Titre         : %-40s │\n", l.Titre)
	fmt.Printf("│ Auteur        : %-40s │\n", l.Auteur)
	fmt.Printf("│ ISBN          : %-40s │\n", l.ISBN)
	fmt.Printf("│ Genre         : %-40s │\n", l.Genre)
	fmt.Printf("│ Publication   : %-40s │\n", l.DatePublication.Format("02/01/2006"))

	// Afficher le statut avec des couleurs (émojis)
	statut := "📗 Disponible"
	if !l.Disponible {
		statut = "📕 Emprunté"
	}
	fmt.Printf("│ Statut        : %-40s │\n", statut)
	fmt.Printf("│ Emprunts      : %-40d │\n", l.NombreEmprunts)
	fmt.Printf("│ Ajouté le     : %-40s │\n", l.DateAjout.Format("02/01/2006 15:04:05"))
	fmt.Printf("└%s┘\n", strings.Repeat("─", 60))
}

func (l Livre) EstDisponible() bool {
	return l.Disponible
}

func (l *Livre) MarquerCommeEmprunte() {
	l.Disponible = false
	l.NombreEmprunts++
}

func (l *Livre) MarquerCommeDispponible() {
	l.Disponible = true
}
