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

func (gl *GestionnaireLivres) TrouverLivreParID(id int) (*models.Livre, int) {

	for i, livre := range gl.livres {
		if livre.ID == id {
			return &gl.livres[i], i
		}
	}
	return nil, -1
}

func (gl *GestionnaireLivres) TrouverLivreParISBN(isbn string) (*models.Livre, int) {
	isbnNettoye := strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", "")

	for i, livre := range gl.livres {
		if strings.EqualFold(livre.ISBN, isbnNettoye) {
			return &gl.livres[i], i
		}
	}
	return nil, -1
}

func (gl *GestionnaireLivres) ModifierLivre(id int, nouveauTitre, nouvelAuteur, nouvelISBN, nouveauGenre, nouvelleDateStr string) error {
	livre, index := gl.TrouverLivreParID(id)

	if livre == nil {
		return fmt.Errorf("aucun livre trouvé avec l'ID %d", id)
	}

	if nouveauTitre != "" {
		if !validators.ValiderTitre(nouveauTitre) {
			return fmt.Errorf("le nouveau titre est invalide")
		}
		livre.Titre = strings.TrimSpace(nouveauTitre)
	}

	if nouvelAuteur != "" {
		if !validators.ValiderNom(nouvelAuteur) {
			return fmt.Errorf("le nouveau nom d'auteur est invalide")
		}
		livre.Auteur = strings.TrimSpace(nouvelAuteur)
	}

	if nouvelISBN != "" {
		if !validators.ValiderISBN(nouvelISBN) {
			return fmt.Errorf("le nouvel ISBN est invalide")
		}

		// Vérifier l'unicité du nouvel ISBN
		isbnNettoye := strings.ReplaceAll(strings.ReplaceAll(nouvelISBN, "-", ""), " ", "")
		for _, l := range gl.livres {
			if l.ID != livre.ID && strings.EqualFold(l.ISBN, isbnNettoye) {
				return fmt.Errorf("l'ISBN %s est déjà utilisé par le livre ID %d", nouvelISBN, l.ID)
			}
		}
		livre.ISBN = isbnNettoye
	}

	if nouveauGenre != "" {
		if !validators.ValiderGenre(nouveauGenre) {
			return fmt.Errorf("le nouveau genre '%s' n'est pas reconnu", nouveauGenre)
		}
		livre.Genre = nouveauGenre
	}

	if nouvelleDateStr != "" {
		nouvelleDate, err := validators.ValiderDatePublication(nouvelleDateStr)
		if err != nil {
			return fmt.Errorf("nouvelle date de publication invalide : %v", err)
		}
		livre.DatePublication = nouvelleDate
	}

	gl.livres[index] = *livre
	return gl.sauvegarderLivres()

}
