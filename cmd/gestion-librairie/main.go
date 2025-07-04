package main

import (
	"log"

	"github.com/felver-dev/bookstore/internal/cli"
	"github.com/felver-dev/bookstore/internal/services"
	"github.com/felver-dev/bookstore/internal/storage"
)

func main() {
	// ========================================
	// INITIALISATION DE L'APPLICATION
	// ========================================

	// 1. Créer les systèmes de stockage pour chaque type de données
	// Chaque service aura son propre fichier JSON
	stockageLivres := storage.NewJSONStorage("data/livres.json")
	stockageMembres := storage.NewJSONStorage("data/membres.json")
	stockageEmprunts := storage.NewJSONStorage("data/emprunts.json")

	// 2. Créer les services (la logique métier de notre application)
	// Ces services contiennent toutes les règles de gestion de la librairie
	gestionnaireL := services.NouveauGestionnaireLivres(stockageLivres)
	gestionnaireM := services.NouveauGestionnaireMembres(stockageMembres)
	gestionnaireE := services.NouveauGestionnaireEmprunts(stockageEmprunts, gestionnaireL, gestionnaireM)

	// 3. Créer l'interface utilisateur en ligne de commande
	// Elle va utiliser tous les services pour offrir un menu complet
	cliApp := cli.NewCLI(gestionnaireL, gestionnaireM, gestionnaireE)

	// 4. Démarrer l'application
	// Si une erreur se produit, on arrête le programme
	if err := cliApp.Run(); err != nil {
		log.Fatal("Erreur lors du démarrage :", err)
	}
}
