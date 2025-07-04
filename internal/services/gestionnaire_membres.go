package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/felver-dev/bookstore/internal/models"
	"github.com/felver-dev/bookstore/internal/storage"
	"github.com/felver-dev/bookstore/internal/validators"
)

type GestionnaireMembres struct {
	membres    []models.Membre
	prochainID int
	stockage   storage.Storage
}

func (gm *GestionnaireMembres) SauvegarderMembres() error {
	return gm.stockage.Sauvegarder(&gm.membres)
}

func (gm *GestionnaireMembres) ChargerMembres() error {
	err := gm.stockage.Charger(&gm.membres)
	if err != nil {
		return err
	}

	for _, membre := range gm.membres {
		if membre.ID >= gm.prochainID {
			gm.prochainID = membre.ID
		}
	}

	return nil
}

func NouveauGestionnaireMembres(stokage storage.Storage) *GestionnaireMembres {
	gm := &GestionnaireMembres{
		membres:    make([]models.Membre, 0),
		prochainID: 1,
		stockage:   stokage,
	}

	gm.ChargerMembres()
	return gm
}

func (gm *GestionnaireMembres) AjouterMembre(nom, email, telephone string) error {
	if !validators.ValiderNom(nom) {
		return fmt.Errorf("le nom du membre est invalide")
	}

	if !validators.ValiderEmail(email) {
		return fmt.Errorf("l'adresse email est invalide")
	}

	if !validators.ValiderTelephone(telephone) {
		return fmt.Errorf("le numéro de téléphone est invalide")
	}

	for _, membre := range gm.membres {
		if strings.EqualFold(membre.Email, email) {
			return fmt.Errorf("un membre avec l'email %s existe déjà (ID: %d %s)", email, membre.ID, membre.Nom)
		}
	}

	maintenant := time.Now()
	nouveauMembre := models.Membre{
		ID:              gm.prochainID,
		Nom:             strings.TrimSpace(nom),
		Email:           strings.ToLower(strings.TrimSpace(email)), // Email en minuscules
		Telephone:       strings.TrimSpace(telephone),
		DateInscription: maintenant,
		NombreEmprunts:  0,    // Aucun emprunt au début
		EmpruntsActifs:  0,    // Aucun emprunt actif au début
		Actif:           true, // Membre actif par défaut
	}

	gm.membres = append(gm.membres, nouveauMembre)
	gm.prochainID++

	return gm.SauvegarderMembres()

}

func (gm *GestionnaireMembres) ListerMembres() []models.Membre {
	return gm.membres
}

func (gm *GestionnaireMembres) ListerMembresActifs() []models.Membre {
	var actifs []models.Membre

	for _, membre := range gm.membres {
		if membre.Actif {
			actifs = append(actifs, membre)
		}
	}

	return actifs
}

func (gm *GestionnaireMembres) RechercherMembres(terme string) []models.Membre {
	var resultats []models.Membre
	terme = strings.ToLower(terme)

	for _, membre := range gm.membres {
		if strings.Contains(strings.ToLower(membre.Nom), terme) ||
			strings.Contains(strings.ToLower(membre.Email), terme) {
			resultats = append(resultats, membre)
		}
	}

	return resultats
}

func (gm *GestionnaireMembres) TrouverMembreParID(id int) (*models.Membre, int) {
	for i, membre := range gm.membres {
		if membre.ID == id {
			return &gm.membres[i], i
		}
	}
	return nil, -1
}

func (gm *GestionnaireMembres) TrouverMembreParEmail(email string) (*models.Membre, int) {
	emailNettoye := strings.ToLower(strings.TrimSpace(email))

	for i, membre := range gm.membres {
		if strings.EqualFold(membre.Email, emailNettoye) {
			return &gm.membres[i], i
		}
	}
	return nil, -1
}

func (gm *GestionnaireMembres) ModifierMembre(id int, nouveauNom, nouvelEmail, nouveauTelephone string) error {
	// 1. TROUVER LE MEMBRE
	membre, index := gm.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	// 2. METTRE À JOUR LES CHAMPS NON VIDES
	if nouveauNom != "" {
		if !validators.ValiderNom(nouveauNom) {
			return fmt.Errorf("le nouveau nom est invalide")
		}
		membre.Nom = strings.TrimSpace(nouveauNom)
	}

	if nouvelEmail != "" {
		if !validators.ValiderEmail(nouvelEmail) {
			return fmt.Errorf("le nouvel email est invalide")
		}

		// Vérifier l'unicité du nouvel email
		emailNettoye := strings.ToLower(strings.TrimSpace(nouvelEmail))
		for _, m := range gm.membres {
			if m.ID != membre.ID && strings.EqualFold(m.Email, emailNettoye) {
				return fmt.Errorf("l'email %s est déjà utilisé par le membre ID %d", nouvelEmail, m.ID)
			}
		}
		membre.Email = emailNettoye
	}

	if nouveauTelephone != "" {
		if !validators.ValiderTelephone(nouveauTelephone) {
			return fmt.Errorf("le nouveau téléphone est invalide")
		}
		membre.Telephone = strings.TrimSpace(nouveauTelephone)
	}

	// 3. SAUVEGARDER LES MODIFICATIONS
	gm.membres[index] = *membre
	return gm.SauvegarderMembres()
}

func (gm *GestionnaireMembres) SuspendirMembre(id int) error {
	membre, index := gm.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	if !membre.Actif {
		return fmt.Errorf("le membre %s est déjà suspendu", membre.Nom)
	}

	membre.Suspendre()
	gm.membres[index] = *membre

	return gm.SauvegarderMembres()
}

func (gm *GestionnaireMembres) ReactiverMembre(id int) error {
	membre, index := gm.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	if membre.Actif {
		return fmt.Errorf("le membre %s est déjà actif", membre.Nom)
	}

	membre.Reactiver()
	gm.membres[index] = *membre

	return gm.SauvegarderMembres()
}

func (gm *GestionnaireMembres) SupprimerMembre(id int) error {
	membre, index := gm.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	// RÈGLE MÉTIER : On ne peut pas supprimer un membre qui a des emprunts en cours
	if membre.EmpruntsActifs > 0 {
		return fmt.Errorf("impossible de supprimer le membre '%s' car il a %d emprunt(s) en cours",
			membre.Nom, membre.EmpruntsActifs)
	}

	// Supprimer le membre de la liste
	gm.membres = append(gm.membres[:index], gm.membres[index+1:]...)

	return gm.SauvegarderMembres()
}

func (gm *GestionnaireMembres) AjouterEmpruntAuMembre(id int) error {
	membre, index := gm.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("membre ID %d introuvable", id)
	}

	if !membre.PeutEmprunter() {
		if !membre.Actif {
			return fmt.Errorf("le membre %s est suspendu", membre.Nom)
		}
		return fmt.Errorf("le membre %s a atteint la limite de %d emprunts simultanés",
			membre.Nom, models.LIMIT_EMPRUNTS_SIMULTANES)
	}

	membre.AjouterEmprunt()
	gm.membres[index] = *membre

	return gm.SauvegarderMembres()
}

func (gm *GestionnaireMembres) RetirerEmpruntDuMembre(id int) error {
	membre, index := gm.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("membre ID %d introuvable", id)
	}

	membre.RetirerEmprunt()
	gm.membres[index] = *membre

	return gm.SauvegarderMembres()
}

func (gm *GestionnaireMembres) ObtenirStatistiques() map[string]interface{} {
	stats := make(map[string]interface{})

	total := len(gm.membres)
	stats["total"] = total

	if total == 0 {
		return stats
	}

	// Compter les membres actifs et suspendus
	actifs := 0
	suspendus := 0
	for _, membre := range gm.membres {
		if membre.Actif {
			actifs++
		} else {
			suspendus++
		}
	}
	stats["actifs"] = actifs
	stats["suspendus"] = suspendus

	// Membre le plus actif (celui qui a emprunté le plus de livres)
	var plusActif models.Membre
	maxEmprunts := -1
	for _, membre := range gm.membres {
		if membre.NombreEmprunts > maxEmprunts {
			maxEmprunts = membre.NombreEmprunts
			plusActif = membre
		}
	}
	if maxEmprunts > 0 {
		stats["plus_actif"] = plusActif
	}

	return stats
}
