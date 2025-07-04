package validators

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ValiderEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValiderTelephone(telephone string) bool {
	nettoye := regexp.MustCompile(`[\s\-s.]`).ReplaceAllString(telephone, "")
	re := regexp.MustCompile(`^\+?[0-9]{8,15}$`)
	return re.MatchString(nettoye)
}

func ValiderISBN(isbn string) bool {
	nettoye := strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", "")

	if len(nettoye) == 10 {
		return validerISBN10(nettoye)
	} else if len(nettoye) == 13 {
		return validerISBN13(nettoye)
	}

	return false
}

func validerISBN10(isbn string) bool {

	if len(isbn) != 10 {
		return false
	}

	for i := 0; i < 9; i++ {
		if isbn[i] < '0' || isbn[i] > '9' {
			return false
		}
	}

	dernier := isbn[9]
	if dernier != 'X' && (dernier < '0' || dernier > '9') {
		return false
	}

	return true
}

func validerISBN13(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}

	for _, char := range isbn {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func ValiderDatePublication(dateStr string) (time.Time, error) {

	date, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	if date.After(time.Now()) {
		return time.Time{}, fmt.Errorf("la date de publication ne peut pas être dans le future")
	}

	anneeMin := time.Date(1440, 1, 1, 0, 0, 0, 0, time.UTC)
	if date.Before(anneeMin) {
		return time.Time{}, fmt.Errorf("la date de publication semble trop ancienne")
	}

	return date, nil
}

func ValiderGenre(genre string) bool {
	genresValides := []string{
		"Roman", "Science-fiction", "Fantasy", "Policier", "Thriller",
		"Romance", "Historique", "Biographie", "Essai", "Poésie",
		"Théâtre", "Bande dessinée", "Manga", "Jeunesse", "Documentaire",
		"Guide pratique", "Cuisine", "Art", "Sport", "Autre",
	}

	for _, genreValide := range genresValides {
		if strings.EqualFold(genre, genreValide) {
			return true
		}
	}

	return false
}

func ValiderNom(nom string) bool {
	if len(strings.TrimSpace(nom)) < 2 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s\-']+$`)
	return re.MatchString(nom)
}

func ValiderTitre(titre string) bool {
	return len(strings.TrimSpace(titre)) >= 1
}
