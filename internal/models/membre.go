package models

import (
	"fmt"
	"strings"
	"time"
)

type Membre struct {
	ID              int       `json:"id"`
	Nom             string    `json:"nom"`
	Email           string    `json:"email"`
	Telephone       string    `json:"telephone"`
	DateInscription time.Time `json:"date_inscription"`
	NombreEmprunts  int       `json:"nombre_emprunts"`
	EmpruntsActifs  int       `json:"emprunts_actifs"`
	Actif           bool      `json:"actif"`
}

const (
	LIMIT_EMPRUNTS_SIMULTANES = 3
)

// Affiche un membre simplement
func (m Membre) String() string {
	statut := "✅ Actif"

	if !m.Actif {
		statut = "❌ Suspendu"
	}

	return fmt.Sprintf("ID: %d | %s | %s | Emprunts: %d/%d | %s  ", m.ID, m.Nom, m.Email, m.EmpruntsActifs, LIMIT_EMPRUNTS_SIMULTANES, statut)
}

// AfficherDetails() montre toutes les informations d'un membre
func (m Membre) AfficherDetails() {
	fmt.Printf("┌%s┐\n", strings.Repeat("─", 60))
	fmt.Printf("│ Membre #%d%s│\n", m.ID, strings.Repeat(" ", 60-len(fmt.Sprintf(" Membre #%d", m.ID))))
	fmt.Printf("├%s┤\n", strings.Repeat("─", 60))
	fmt.Printf("│ Nom           : %-40s │\n", m.Nom)
	fmt.Printf("│ Email         : %-40s │\n", m.Email)
	fmt.Printf("│ Téléphone     : %-40s │\n", m.Telephone)
	fmt.Printf("│ Inscrit le    : %-40s │\n", m.DateInscription.Format("02/01/2006"))

	statut := "✅ Actif"
	if !m.Actif {
		statut = "❌ Suspendu"
	}
	fmt.Printf("│ Statut        : %-40s │\n", statut)
	fmt.Printf("│ Emprunts totaux : %-37d │\n", m.NombreEmprunts)
	fmt.Printf("│ Emprunts actifs : %-37d │\n", m.EmpruntsActifs)
	fmt.Printf("└%s┘\n", strings.Repeat("─", 60))
}

func (m Membre) PeutEmprunter() bool {
	return m.Actif && m.EmpruntsActifs < LIMIT_EMPRUNTS_SIMULTANES
}

func (m *Membre) AjouterEmprunt() {
	m.EmpruntsActifs++
	m.NombreEmprunts++
}

func (m *Membre) RetirerEmprunt() {
	if m.EmpruntsActifs > 0 {
		m.EmpruntsActifs--
	}
}

func (m *Membre) Suspendre() {
	m.Actif = false
}

func (m *Membre) Reactiver() {
	m.Actif = true
}
