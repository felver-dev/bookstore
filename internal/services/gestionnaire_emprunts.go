package services

import (
	"fmt"
	"time"

	"github.com/felver-dev/bookstore/internal/models"
	"github.com/felver-dev/bookstore/internal/storage"
)

type GestionnaireEmprunts struct {
	emprunts            []models.Emprunt
	prochainID          int
	stockage            storage.Storage
	gestionnaireLivres  *GestionnaireLivres
	gestionnaireMembres *GestionnaireMembres
}

type statMembre struct {
	nom   string
	count int
}

type statLivre struct {
	titre string
	count int
}

func (ge *GestionnaireEmprunts) sauvegarderEmprunts() error {
	return ge.stockage.Sauvegarder(ge.emprunts)
}

func (ge *GestionnaireEmprunts) ChargerEmprunts() error {
	err := ge.stockage.Charger(&ge.emprunts)
	if err != nil {
		return err
	}

	// Calculer le prochain ID
	for _, emprunt := range ge.emprunts {
		if emprunt.ID >= ge.prochainID {
			ge.prochainID = emprunt.ID + 1
		}
	}

	return nil
}

func NouveauGestionnaireEmprunts(stockage storage.Storage, gl *GestionnaireLivres, gm *GestionnaireMembres) *GestionnaireEmprunts {
	ge := &GestionnaireEmprunts{
		emprunts:            make([]models.Emprunt, 0),
		prochainID:          1,
		stockage:            stockage,
		gestionnaireLivres:  gl,
		gestionnaireMembres: gm,
	}

	ge.ChargerEmprunts()
	ge.mettreAJourStatutsEmprunts() // Vérifier les retards au démarrage
	return ge
}

func (ge *GestionnaireEmprunts) EmprunterLivre(livreID, membreID int) error {
	// 1. VÉRIFICATIONS PRÉALABLES

	// Vérifier que le livre existe et est disponible
	livre, _ := ge.gestionnaireLivres.TrouverLivreParID(livreID)
	if livre == nil {
		return fmt.Errorf("livre ID %d introuvable", livreID)
	}

	if !livre.EstDisponible() {
		return fmt.Errorf("le livre '%s' n'est pas disponible (actuellement emprunté)", livre.Titre)
	}

	// Vérifier que le membre existe et peut emprunter
	membre, _ := ge.gestionnaireMembres.TrouverMembreParID(membreID)
	if membre == nil {
		return fmt.Errorf("membre ID %d introuvable", membreID)
	}

	if !membre.PeutEmprunter() {
		if !membre.Actif {
			return fmt.Errorf("le membre %s est suspendu et ne peut pas emprunter", membre.Nom)
		}
		return fmt.Errorf("le membre %s a atteint la limite de %d emprunts simultanés",
			membre.Nom, models.LIMIT_EMPRUNTS_SIMULTANES)
	}

	// Vérifier que ce membre n'a pas déjà emprunté ce livre et ne l'a pas encore rendu
	for _, emprunt := range ge.emprunts {
		if emprunt.LivreID == livreID && emprunt.MembreID == membreID && emprunt.DateRetourEffectif == nil {
			return fmt.Errorf("le membre %s a déjà emprunté ce livre et ne l'a pas encore rendu", membre.Nom)
		}
	}

	// 2. CRÉER L'EMPRUNT
	maintenant := time.Now()
	dateRetourPrevu := maintenant.AddDate(0, 0, models.DUREE_EMPRUMT_JOURS) // Ajouter 14 jours

	nouvelEmprunt := models.Emprunt{
		ID:                 ge.prochainID,
		LivreID:            livreID,
		MembreID:           membreID,
		DateEmprunt:        maintenant,
		DateRetourPrevu:    dateRetourPrevu,
		DateRetourEffectif: nil, // Pas encore rendu
		Statut:             models.STATUT_EN_COURS,

		// Informations dénormalisées pour faciliter l'affichage
		TitreLivre: livre.Titre,
		NomMembre:  membre.Nom,
	}

	// 3. METTRE À JOUR LES ÉTATS
	// Marquer le livre comme emprunté
	if err := ge.gestionnaireLivres.MarquerCommeEmprunte(livreID); err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du livre : %v", err)
	}

	// Mettre à jour les compteurs du membre
	if err := ge.gestionnaireMembres.AjouterEmpruntAuMembre(membreID); err != nil {
		// Si erreur, annuler la modification du livre
		ge.gestionnaireLivres.MarquerCommeDisponible(livreID)
		return fmt.Errorf("erreur lors de la mise à jour du membre : %v", err)
	}

	// 4. ENREGISTRER L'EMPRUNT
	ge.emprunts = append(ge.emprunts, nouvelEmprunt)
	ge.prochainID++

	return ge.sauvegarderEmprunts()
}

func (ge *GestionnaireEmprunts) RetournerLivre(empruntID int) error {
	// 1. TROUVER L'EMPRUNT
	emprunt, index := ge.TrouverEmpruntParID(empruntID)
	if emprunt == nil {
		return fmt.Errorf("emprunt ID %d introuvable", empruntID)
	}

	// Vérifier que l'emprunt n'est pas déjà terminé
	if emprunt.DateRetourEffectif != nil {
		return fmt.Errorf("cet emprunt est déjà terminé (livre rendu le %s)",
			emprunt.DateRetourEffectif.Format("02/01/2006"))
	}

	// 2. METTRE À JOUR L'EMPRUNT
	emprunt.MarquerCommeRendu()

	// 3. METTRE À JOUR LES ÉTATS
	// Marquer le livre comme disponible
	if err := ge.gestionnaireLivres.MarquerCommeDisponible(emprunt.LivreID); err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du livre : %v", err)
	}

	// Mettre à jour les compteurs du membre
	if err := ge.gestionnaireMembres.RetirerEmpruntDuMembre(emprunt.MembreID); err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du membre : %v", err)
	}

	// 4. SAUVEGARDER
	ge.emprunts[index] = *emprunt
	return ge.sauvegarderEmprunts()
}

func (ge *GestionnaireEmprunts) ListerEmprunts() []models.Emprunt {
	// Mettre à jour les statuts avant de retourner la liste
	ge.mettreAJourStatutsEmprunts()
	return ge.emprunts
}

func (ge *GestionnaireEmprunts) ListerEmpruntsEnCours() []models.Emprunt {
	var enCours []models.Emprunt

	// Mettre à jour les statuts d'abord
	ge.mettreAJourStatutsEmprunts()

	for _, emprunt := range ge.emprunts {
		if emprunt.DateRetourEffectif == nil { // Pas encore rendu
			enCours = append(enCours, emprunt)
		}
	}

	return enCours
}

func (ge *GestionnaireEmprunts) ListerEmpruntsEnRetard() []models.Emprunt {
	var enRetard []models.Emprunt

	// Mettre à jour les statuts d'abord
	ge.mettreAJourStatutsEmprunts()

	for _, emprunt := range ge.emprunts {
		if emprunt.EstEnRetard() {
			enRetard = append(enRetard, emprunt)
		}
	}

	return enRetard
}

func (ge *GestionnaireEmprunts) ListerEmpruntsParLivre(livreID int) []models.Emprunt {
	var empruntsLivre []models.Emprunt

	for _, emprunt := range ge.emprunts {
		if emprunt.LivreID == livreID {
			empruntsLivre = append(empruntsLivre, emprunt)
		}
	}

	return empruntsLivre
}

func (ge *GestionnaireEmprunts) ListerEmpruntsParMembre(membreID int) []models.Emprunt {
	var empruntsMemb []models.Emprunt

	for _, emprunt := range ge.emprunts {
		if emprunt.MembreID == membreID {
			empruntsMemb = append(empruntsMemb, emprunt)
		}
	}

	return empruntsMemb
}

func (ge *GestionnaireEmprunts) TrouverEmpruntParID(id int) (*models.Emprunt, int) {
	for i, emprunt := range ge.emprunts {
		if emprunt.ID == id {
			return &ge.emprunts[i], i
		}
	}
	return nil, -1
}

func (ge *GestionnaireEmprunts) TrouverEmpruntActifParLivre(livreID int) (*models.Emprunt, int) {
	for i, emprunt := range ge.emprunts {
		if emprunt.LivreID == livreID && emprunt.DateRetourEffectif == nil {
			return &ge.emprunts[i], i
		}
	}
	return nil, -1
}

func (ge *GestionnaireEmprunts) ObtenirEmpruntsParPeriode(dateDebut, dateFin time.Time) []models.Emprunt {
	var empruntsP []models.Emprunt

	for _, emprunt := range ge.emprunts {
		if emprunt.DateEmprunt.After(dateDebut) && emprunt.DateEmprunt.Before(dateFin) {
			empruntsP = append(empruntsP, emprunt)
		}
	}

	return empruntsP
}

func (ge *GestionnaireEmprunts) ObtenirEmpruntsARendreAujourdhui() []models.Emprunt {
	var aRendreAujourdhui []models.Emprunt
	aujourd_hui := time.Now().Truncate(24 * time.Hour)

	for _, emprunt := range ge.emprunts {
		if emprunt.DateRetourEffectif == nil { // Pas encore rendu
			dateRetourTrunc := emprunt.DateRetourPrevu.Truncate(24 * time.Hour)
			if dateRetourTrunc.Equal(aujourd_hui) {
				aRendreAujourdhui = append(aRendreAujourdhui, emprunt)
			}
		}
	}

	return aRendreAujourdhui
}

func (ge *GestionnaireEmprunts) CalculerDureeEmpruntsTermines() float64 {
	var totalJours int
	var count int

	for _, emprunt := range ge.emprunts {
		if emprunt.DateRetourEffectif != nil { // Emprunt terminé
			duree := emprunt.DateRetourEffectif.Sub(emprunt.DateEmprunt)
			totalJours += int(duree.Hours() / 24)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return float64(totalJours) / float64(count)
}

func (ge *GestionnaireEmprunts) PrologerEmprunt(empruntID int, joursSupplementaires int) error {
	emprunt, index := ge.TrouverEmpruntParID(empruntID)
	if emprunt == nil {
		return fmt.Errorf("emprunt ID %d introuvable", empruntID)
	}

	// Vérifier que l'emprunt est en cours
	if emprunt.DateRetourEffectif != nil {
		return fmt.Errorf("impossible de prolonger un emprunt déjà terminé")
	}

	// Prolonger la date de retour
	emprunt.DateRetourPrevu = emprunt.DateRetourPrevu.AddDate(0, 0, joursSupplementaires)

	// Mettre à jour le statut si nécessaire
	emprunt.MettreAjourStatut()

	ge.emprunts[index] = *emprunt
	return ge.sauvegarderEmprunts()
}

func (ge *GestionnaireEmprunts) AnnulerEmprunt(empruntID int) error {
	emprunt, index := ge.TrouverEmpruntParID(empruntID)
	if emprunt == nil {
		return fmt.Errorf("emprunt ID %d introuvable", empruntID)
	}

	// Vérifier que l'emprunt peut être annulé (pas encore rendu)
	if emprunt.DateRetourEffectif != nil {
		return fmt.Errorf("impossible d'annuler un emprunt déjà terminé")
	}

	// Remettre le livre disponible
	if err := ge.gestionnaireLivres.MarquerCommeDisponible(emprunt.LivreID); err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du livre : %v", err)
	}

	// Retirer l'emprunt du membre
	if err := ge.gestionnaireMembres.RetirerEmpruntDuMembre(emprunt.MembreID); err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du membre : %v", err)
	}

	// Supprimer l'emprunt de la liste
	ge.emprunts = append(ge.emprunts[:index], ge.emprunts[index+1:]...)

	return ge.sauvegarderEmprunts()
}

func (ge *GestionnaireEmprunts) ObtenirStatistiques() map[string]interface{} {
	stats := make(map[string]interface{})

	// Mettre à jour les statuts avant de calculer les stats
	ge.mettreAJourStatutsEmprunts()

	total := len(ge.emprunts)
	stats["total"] = total

	if total == 0 {
		return stats
	}

	// Compter par statut
	enCours := 0
	rendus := 0
	enRetard := 0

	for _, emprunt := range ge.emprunts {
		switch emprunt.Statut {
		case models.STATUT_EN_COURS:
			enCours++
		case models.STATUT_RENDU:
			rendus++
		case models.STATUT_EN_RETARD:
			enRetard++
		}
	}

	stats["en_cours"] = enCours
	stats["rendus"] = rendus
	stats["en_retard"] = enRetard

	// Durée moyenne des emprunts terminés
	dureeM := ge.CalculerDureeEmpruntsTermines()
	stats["duree_moyenne_jours"] = dureeM

	// Emprunts par mois (12 derniers mois)
	empruntsParMois := ge.calculerEmpruntsParMois()
	stats["emprunts_par_mois"] = empruntsParMois

	// Membre le plus actif (le plus d'emprunts)
	membrePlusActif := ge.trouverMembrePlusActif()
	if membrePlusActif != nil {
		stats["membre_plus_actif"] = map[string]interface{}{
			"nom":      membrePlusActif.nom,
			"emprunts": membrePlusActif.count,
		}
	}

	// Livre le plus emprunté
	livrePlusEmprunte := ge.trouverLivrePlusEmprunte()
	if livrePlusEmprunte != nil {
		stats["livre_plus_emprunte"] = map[string]interface{}{
			"titre":    livrePlusEmprunte.titre,
			"emprunts": livrePlusEmprunte.count,
		}
	}

	return stats
}
func (ge *GestionnaireEmprunts) NettoierEmpruntsAnciens(ageMaxAnnees int) error {
	dateLimit := time.Now().AddDate(-ageMaxAnnees, 0, 0)
	var empruntsAGarder []models.Emprunt
	supprimesCount := 0

	for _, emprunt := range ge.emprunts {
		// Garder l'emprunt s'il est récent OU s'il n'est pas encore terminé
		if emprunt.DateEmprunt.After(dateLimit) || emprunt.DateRetourEffectif == nil {
			empruntsAGarder = append(empruntsAGarder, emprunt)
		} else {
			supprimesCount++
		}
	}

	if supprimesCount > 0 {
		ge.emprunts = empruntsAGarder
		err := ge.sauvegarderEmprunts()
		if err != nil {
			return fmt.Errorf("erreur lors de la sauvegarde après nettoyage : %v", err)
		}
	}

	return nil
}

func (ge *GestionnaireEmprunts) ExporterRapportEmprunts() string {
	stats := ge.ObtenirStatistiques()
	rapport := "=== RAPPORT DES EMPRUNTS ===\n\n"

	rapport += fmt.Sprintf("Total des emprunts : %d\n", stats["total"])
	rapport += fmt.Sprintf("En cours : %d\n", stats["en_cours"])
	rapport += fmt.Sprintf("Rendus : %d\n", stats["rendus"])
	rapport += fmt.Sprintf("En retard : %d\n", stats["en_retard"])

	if duree, ok := stats["duree_moyenne_jours"].(float64); ok && duree > 0 {
		rapport += fmt.Sprintf("Durée moyenne : %.1f jours\n", duree)
	}

	rapport += "\n=== EMPRUNTS EN RETARD ===\n"
	empruntsEnRetard := ge.ListerEmpruntsEnRetard()
	if len(empruntsEnRetard) == 0 {
		rapport += "Aucun emprunt en retard\n"
	} else {
		for _, emprunt := range empruntsEnRetard {
			joursRetard := emprunt.CalculerJoursRetard()
			rapport += fmt.Sprintf("- %s (%s) - %d jour(s) de retard\n",
				emprunt.TitreLivre, emprunt.NomMembre, joursRetard)
		}
	}

	return rapport
}

func (ge *GestionnaireEmprunts) calculerEmpruntsParMois() map[string]int {
	empruntsParMois := make(map[string]int)
	maintenant := time.Now()

	// Initialiser les 12 derniers mois à 0
	for i := 11; i >= 0; i-- {
		mois := maintenant.AddDate(0, -i, 0).Format("2006-01")
		empruntsParMois[mois] = 0
	}

	// Compter les emprunts par mois
	for _, emprunt := range ge.emprunts {
		mois := emprunt.DateEmprunt.Format("2006-01")
		if _, existe := empruntsParMois[mois]; existe {
			empruntsParMois[mois]++
		}
	}

	return empruntsParMois
}

func (ge *GestionnaireEmprunts) trouverMembrePlusActif() *statMembre {
	compteurMembres := make(map[int]int)
	nomsMembres := make(map[int]string)

	// Compter les emprunts par membre
	for _, emprunt := range ge.emprunts {
		compteurMembres[emprunt.MembreID]++
		nomsMembres[emprunt.MembreID] = emprunt.NomMembre
	}

	// Trouver le maximum
	var membrePlusActif *statMembre
	maxEmprunts := 0

	for membreID, count := range compteurMembres {
		if count > maxEmprunts {
			maxEmprunts = count
			membrePlusActif = &statMembre{
				nom:   nomsMembres[membreID],
				count: count,
			}
		}
	}

	return membrePlusActif
}

func (ge *GestionnaireEmprunts) trouverLivrePlusEmprunte() *statLivre {
	compteurLivres := make(map[int]int)
	titresLivres := make(map[int]string)

	// Compter les emprunts par livre
	for _, emprunt := range ge.emprunts {
		compteurLivres[emprunt.LivreID]++
		titresLivres[emprunt.LivreID] = emprunt.TitreLivre
	}

	// Trouver le maximum
	var livrePlusEmprunte *statLivre
	maxEmprunts := 0

	for livreID, count := range compteurLivres {
		if count > maxEmprunts {
			maxEmprunts = count
			livrePlusEmprunte = &statLivre{
				titre: titresLivres[livreID],
				count: count,
			}
		}
	}

	return livrePlusEmprunte
}

func (ge *GestionnaireEmprunts) mettreAJourStatutsEmprunts() {
	modifie := false

	for i := range ge.emprunts {
		ancienStatut := ge.emprunts[i].Statut
		ge.emprunts[i].MettreAjourStatut()

		if ge.emprunts[i].Statut != ancienStatut {
			modifie = true
		}
	}

	// Sauvegarder si des modifications ont été apportées
	if modifie {
		ge.sauvegarderEmprunts()
	}
}
