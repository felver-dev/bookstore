package services

import (
	"github.com/felver-dev/bookstore/internal/models"
	"github.com/felver-dev/bookstore/internal/storage"
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
