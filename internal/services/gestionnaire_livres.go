package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/felver-dev/bookstore/internal/models"
	"github.com/felver-dev/bookstore/internal/storage"
	"github.com/felver-dev/bookstore/internal/validators"
)

type GestionnaireLivres struct {
	livres     []models.Livre
	prochainID int
	stockage   storage.Storage
}

func (gl *GestionnaireLivres) ChargerLivres() error {
	err := gl.stockage.Charger(&gl.livres)
	if err != nil {
		return err
	}

	for _, livre := range gl.livres {
		if livre.ID >= gl.prochainID {
			gl.prochainID = livre.ID
		}
	}

	return nil
}

func (gl *GestionnaireLivres) sauvegarderLivres() error {
	return gl.stockage.Sauvegarder(gl.livres)
}

func NouveauGestionnaireLivres(stockage storage.Storage) *GestionnaireLivres {
	gl := &GestionnaireLivres{
		livres:     make([]models.Livre, 0),
		prochainID: 1,
		stockage:   stockage,
	}

	gl.ChargerLivres()
	return gl
}

// Methodes publiques

func (gl *GestionnaireLivres) AjouterLivre(titre, auteur, isbn, genre, datePublicationStr string) error {
	if !validators.ValiderTitre(titre) {
		return fmt.Errorf("le titre du livre est invalide")
	}

	if !validators.ValiderNom(auteur) {
		return fmt.Errorf("le nom de l'auteur est invalide")
	}

	if !validators.ValiderISBN(isbn) {
		return fmt.Errorf("l'ISBN est invalide (doit faire 10 ou 13 caractères)")
	}

	if !validators.ValiderGenre(genre) {
		return fmt.Errorf("le genre '%s' n'est pas reconnu", genre)
	}

	datePublication, err := validators.ValiderDatePublication(datePublicationStr)
	if err != nil {
		return fmt.Errorf("date de publication invalide : %v", err)
	}

	for _, livre := range gl.livres {
		if strings.EqualFold(livre.ISBN, isbn) {
			return fmt.Errorf("un livre avec l'ISBN %s existe déjà (ID : %d - %s)", isbn, livre.ID, livre.Titre)
		}
	}

	maintenant := time.Now()
	nouveauLivre := models.Livre{
		ID:              gl.prochainID,
		Titre:           strings.TrimSpace(titre),
		Auteur:          strings.TrimSpace(auteur),
		ISBN:            strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", ""),
		Genre:           genre,
		DatePublication: datePublication,
		Disponible:      true,
		NombreEmprunts:  0,
		DateAjout:       maintenant,
	}

	gl.livres = append(gl.livres, nouveauLivre)
	gl.prochainID++
	return gl.sauvegarderLivres()
}

func (gl *GestionnaireLivres) ListerLivres() []models.Livre {
	var disponibles []models.Livre

	for _, livre := range gl.livres {
		if livre.EstDisponible() {
			disponibles = append(disponibles, livre)
		}
	}

	return disponibles
}

func (gl *GestionnaireLivres) RechercherLivres(terme string) []models.Livre {
	var resultats []models.Livre
	terme = strings.TrimSpace(terme)

	for _, livre := range gl.livres {
		if strings.Contains(strings.ToLower(livre.Titre), terme) || strings.Contains(strings.ToLower(livre.Auteur), terme) || strings.Contains(strings.ToLower(livre.Genre), terme) {

			resultats = append(resultats, livre)
		}
	}

	return resultats
}
