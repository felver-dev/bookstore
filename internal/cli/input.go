package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ========================================
// FONCTIONS UTILITAIRES POUR LES SAISIES UTILISATEUR
// Ces fonctions gèrent toutes les interactions avec l'utilisateur
// ========================================

// LireEntree lit une ligne d'entrée utilisateur et supprime les espaces
func LireEntree() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// LireEntreeObligatoire lit une entrée qui ne peut pas être vide
// Continue à demander tant que l'utilisateur n'entre rien
func LireEntreeObligatoire(message string) string {
	for {
		fmt.Print(message)
		entree := LireEntree()
		if entree != "" {
			return entree
		}
		fmt.Println("❌ Cette information est obligatoire.")
	}
}

// LireEntreeEntier lit un nombre entier avec validation
func LireEntreeEntier(message string) (int, error) {
	fmt.Print(message)
	entreeStr := LireEntree()

	if entreeStr == "" {
		return 0, fmt.Errorf("aucune valeur saisie")
	}

	// Convertir la chaîne en entier
	valeur, err := strconv.Atoi(entreeStr)
	if err != nil {
		return 0, fmt.Errorf("'%s' n'est pas un nombre valide", entreeStr)
	}

	return valeur, nil
}

// LireEntreeEntierObligatoire lit un entier obligatoire avec validation
func LireEntreeEntierObligatoire(message string) int {
	for {
		valeur, err := LireEntreeEntier(message)
		if err == nil {
			return valeur
		}
		fmt.Printf("❌ Erreur : %s\n", err.Error())
	}
}

// LireEntreeEntierAvecLimites lit un entier dans une plage donnée
func LireEntreeEntierAvecLimites(message string, min, max int) int {
	for {
		valeur := LireEntreeEntierObligatoire(message)
		if valeur >= min && valeur <= max {
			return valeur
		}
		fmt.Printf("❌ La valeur doit être entre %d et %d.\n", min, max)
	}
}

// LireConfirmation demande une confirmation oui/non
func LireConfirmation(message string) bool {
	fmt.Print(message + " (oui/non) : ")
	reponse := strings.ToLower(LireEntree())

	// Accepter plusieurs variantes de "oui"
	return reponse == "oui" || reponse == "o" || reponse == "yes" || reponse == "y"
}

// LireChoixDansListe affiche une liste d'options et demande à l'utilisateur de choisir
func LireChoixDansListe(message string, options []string) int {
	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	return LireEntreeEntierAvecLimites("Votre choix : ", 1, len(options)) - 1 // Retourner l'index (0-based)
}

// AfficherSeparateur affiche une ligne de séparation visuelle
func AfficherSeparateur(caractere string, longueur int) {
	fmt.Println(strings.Repeat(caractere, longueur))
}

// AfficherTitre affiche un titre encadré
func AfficherTitre(titre string) {
	longueur := len(titre) + 4
	if longueur < 40 {
		longueur = 40
	}

	fmt.Println("\n" + strings.Repeat("=", longueur))
	fmt.Printf("  %s\n", titre)
	fmt.Println(strings.Repeat("=", longueur))
}

// AfficherSucces affiche un message de succès avec un emoji
func AfficherSucces(message string) {
	fmt.Printf("\n✅ %s\n", message)
}

// AfficherErreur affiche un message d'erreur avec un emoji
func AfficherErreur(message string) {
	fmt.Printf("\n❌ %s\n", message)
}

// AfficherInfo affiche un message d'information avec un emoji
func AfficherInfo(message string) {
	fmt.Printf("\nℹ️  %s\n", message)
}

// AfficherAvertissement affiche un avertissement avec un emoji
func AfficherAvertissement(message string) {
	fmt.Printf("\n⚠️  %s\n", message)
}

// AttendreEntree affiche un message et attend que l'utilisateur appuie sur Entrée
func AttendreEntree(message string) {
	if message == "" {
		message = "Appuyez sur Entrée pour continuer..."
	}
	fmt.Println(message)
	LireEntree()
}
