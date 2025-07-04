// ==========================================
// internal/cli/menu.go
// INTERFACE CLI PRINCIPALE AVEC TOUS LES MENUS
// ==========================================

package cli

import (
	"fmt"
	"strings"

	"github.com/felver-dev/bookstore/internal/services"
)

// ========================================
// INTERFACE CLI PRINCIPALE
// Coordonne tous les gestionnaires et offre un menu complet
// ========================================

type CLI struct {
	gestionnaireLivres   *services.GestionnaireLivres
	gestionnaireMembres  *services.GestionnaireMembres
	gestionnaireEmprunts *services.GestionnaireEmprunts
}

// NewCLI crée une nouvelle instance de l'interface CLI
func NewCLI(gl *services.GestionnaireLivres, gm *services.GestionnaireMembres, ge *services.GestionnaireEmprunts) *CLI {
	return &CLI{
		gestionnaireLivres:   gl,
		gestionnaireMembres:  gm,
		gestionnaireEmprunts: ge,
	}
}

// Run démarre l'application et affiche le menu principal
func (cli *CLI) Run() error {
	fmt.Println("📚 Bienvenue dans le Système de Gestion de Librairie !")
	fmt.Println("📁 Données sauvegardées dans le dossier 'data/'")

	for {
		cli.afficherMenuPrincipal()
		choix := LireEntreeEntierAvecLimites("Votre choix : ", 0, 4)

		var err error
		switch choix {
		case 1:
			err = cli.menuLivres()
		case 2:
			err = cli.menuMembres()
		case 3:
			err = cli.menuEmprunts()
		case 4:
			cli.afficherStatistiques()
		case 0:
			fmt.Println("\n👋 Au revoir ! Toutes les données ont été sauvegardées.")
			return nil
		}

		if err != nil {
			AfficherErreur(err.Error())
		}

		AttendreEntree("")
	}
}

// ========================================
// MENU PRINCIPAL
// ========================================

func (cli *CLI) afficherMenuPrincipal() {
	AfficherTitre("📚 GESTION DE LIBRAIRIE - MENU PRINCIPAL")
	fmt.Println("1. 📖 Gestion des Livres")
	fmt.Println("2. 👥 Gestion des Membres")
	fmt.Println("3. 📋 Gestion des Emprunts")
	fmt.Println("4. 📊 Statistiques")
	fmt.Println("0. 🚪 Quitter")
	AfficherSeparateur("-", 50)
}

// ========================================
// SOUS-MENU LIVRES
// ========================================

func (cli *CLI) menuLivres() error {
	for {
		AfficherTitre("📖 GESTION DES LIVRES")
		fmt.Println("1. ➕ Ajouter un livre")
		fmt.Println("2. 📋 Lister tous les livres")
		fmt.Println("3. 📗 Lister les livres disponibles")
		fmt.Println("4. 🔍 Rechercher des livres")
		fmt.Println("5. ✏️  Modifier un livre")
		fmt.Println("6. 🗑️  Supprimer un livre")
		fmt.Println("0. ⬅️  Retour au menu principal")
		AfficherSeparateur("-", 50)

		choix := LireEntreeEntierAvecLimites("Votre choix : ", 0, 6)

		var err error
		switch choix {
		case 1:
			err = cli.ajouterLivre()
		case 2:
			cli.listerLivres()
		case 3:
			cli.listerLivresDisponibles()
		case 4:
			cli.rechercherLivres()
		case 5:
			err = cli.modifierLivre()
		case 6:
			err = cli.supprimerLivre()
		case 0:
			return nil
		}

		if err != nil {
			AfficherErreur(err.Error())
		}

		AttendreEntree("")
	}
}

func (cli *CLI) ajouterLivre() error {
	AfficherTitre("➕ AJOUTER UN LIVRE")

	// Saisir les informations du livre
	titre := LireEntreeObligatoire("Titre du livre : ")
	auteur := LireEntreeObligatoire("Auteur : ")
	isbn := LireEntreeObligatoire("ISBN (10 ou 13 caractères) : ")

	// Proposer une liste de genres
	genres := []string{
		"Roman", "Science-fiction", "Fantasy", "Policier", "Thriller",
		"Romance", "Historique", "Biographie", "Essai", "Poésie",
		"Théâtre", "Bande dessinée", "Manga", "Jeunesse", "Documentaire",
		"Guide pratique", "Cuisine", "Art", "Sport", "Autre",
	}
	indexGenre := LireChoixDansListe("Choisissez le genre :", genres)
	genre := genres[indexGenre]

	datePublication := LireEntreeObligatoire("Date de publication (JJ/MM/AAAA) : ")

	// Appeler le service pour ajouter le livre
	err := cli.gestionnaireLivres.AjouterLivre(titre, auteur, isbn, genre, datePublication)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Livre '%s' ajouté avec succès !", titre))
	return nil
}

func (cli *CLI) listerLivres() {
	AfficherTitre("📋 LISTE DE TOUS LES LIVRES")

	livres := cli.gestionnaireLivres.ListerLivres()

	if len(livres) == 0 {
		AfficherInfo("Aucun livre enregistré.")
		return
	}

	cli.afficherTableauLivres(livres)
}

func (cli *CLI) listerLivresDisponibles() {
	AfficherTitre("📗 LIVRES DISPONIBLES À L'EMPRUNT")

	livres := cli.gestionnaireLivres.listerLivresDisponibles()

	if len(livres) == 0 {
		AfficherInfo("Aucun livre disponible actuellement.")
		return
	}

	cli.afficherTableauLivres(livres)
}

func (cli *CLI) rechercherLivres() {
	AfficherTitre("🔍 RECHERCHER DES LIVRES")

	terme := LireEntreeObligatoire("Terme de recherche (titre, auteur ou genre) : ")

	resultats := cli.gestionnaireLivres.RechercherLivres(terme)

	fmt.Printf("\n🎯 %d résultat(s) trouvé(s) pour '%s' :\n", len(resultats), terme)

	if len(resultats) == 0 {
		AfficherInfo("Aucun livre correspondant.")
		return
	}

	cli.afficherTableauLivres(resultats)
}

func (cli *CLI) modifierLivre() error {
	AfficherTitre("✏️ MODIFIER UN LIVRE")

	id := LireEntreeEntierObligatoire("ID du livre à modifier : ")

	livre, _ := cli.gestionnaireLivres.TrouverLivreParID(id)
	if livre == nil {
		return fmt.Errorf("aucun livre trouvé avec l'ID %d", id)
	}

	// Afficher les informations actuelles
	fmt.Println("\nInformations actuelles :")
	livre.AfficherDetails()

	AfficherInfo("Laissez vide pour conserver la valeur actuelle.")

	// Saisir les nouvelles valeurs
	fmt.Printf("Nouveau titre (%s) : ", livre.Titre)
	nouveauTitre := LireEntree()

	fmt.Printf("Nouvel auteur (%s) : ", livre.Auteur)
	nouvelAuteur := LireEntree()

	fmt.Printf("Nouvel ISBN (%s) : ", livre.ISBN)
	nouvelISBN := LireEntree()

	fmt.Printf("Nouveau genre (%s) : ", livre.Genre)
	nouveauGenre := LireEntree()

	fmt.Printf("Nouvelle date de publication (%s) : ", livre.DatePublication.Format("02/01/2006"))
	nouvelleDateStr := LireEntree()

	// Appeler le service pour modifier
	err := cli.gestionnaireLivres.ModifierLivre(id, nouveauTitre, nouvelAuteur, nouvelISBN, nouveauGenre, nouvelleDateStr)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Livre ID %d modifié avec succès !", id))
	return nil
}

func (cli *CLI) supprimerLivre() error {
	AfficherTitre("🗑️ SUPPRIMER UN LIVRE")

	id := LireEntreeEntierObligatoire("ID du livre à supprimer : ")

	livre, _ := cli.gestionnaireLivres.TrouverLivreParID(id)
	if livre == nil {
		return fmt.Errorf("aucun livre trouvé avec l'ID %d", id)
	}

	// Afficher le livre à supprimer
	fmt.Println("\nLivre à supprimer :")
	livre.AfficherDetails()

	// Demander confirmation
	if !LireConfirmation("\n⚠️ Êtes-vous sûr de vouloir supprimer ce livre ?") {
		AfficherInfo("Suppression annulée.")
		return nil
	}

	titre := livre.Titre
	err := cli.gestionnaireLivres.SupprimerLivre(id)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Livre '%s' (ID: %d) supprimé avec succès !", titre, id))
	return nil
}

// ========================================
// SOUS-MENU MEMBRES
// ========================================

func (cli *CLI) menuMembres() error {
	for {
		AfficherTitre("👥 GESTION DES MEMBRES")
		fmt.Println("1. ➕ Inscrire un membre")
		fmt.Println("2. 📋 Lister tous les membres")
		fmt.Println("3. ✅ Lister les membres actifs")
		fmt.Println("4. 🔍 Rechercher des membres")
		fmt.Println("5. ✏️  Modifier un membre")
		fmt.Println("6. ⛔ Suspendre un membre")
		fmt.Println("7. ✅ Réactiver un membre")
		fmt.Println("8. 🗑️  Supprimer un membre")
		fmt.Println("0. ⬅️  Retour au menu principal")
		AfficherSeparateur("-", 50)

		choix := LireEntreeEntierAvecLimites("Votre choix : ", 0, 8)

		var err error
		switch choix {
		case 1:
			err = cli.inscrireMembre()
		case 2:
			cli.listerMembres()
		case 3:
			cli.listerMembresActifs()
		case 4:
			cli.rechercherMembres()
		case 5:
			err = cli.modifierMembre()
		case 6:
			err = cli.suspendirMembre()
		case 7:
			err = cli.reactiverMembre()
		case 8:
			err = cli.supprimerMembre()
		case 0:
			return nil
		}

		if err != nil {
			AfficherErreur(err.Error())
		}

		AttendreEntree("")
	}
}

func (cli *CLI) inscrireMembre() error {
	AfficherTitre("➕ INSCRIRE UN MEMBRE")

	nom := LireEntreeObligatoire("Nom complet : ")
	email := LireEntreeObligatoire("Adresse email : ")
	telephone := LireEntreeObligatoire("Numéro de téléphone : ")

	err := cli.gestionnaireMembres.AjouterMembre(nom, email, telephone)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Membre '%s' inscrit avec succès !", nom))
	return nil
}

func (cli *CLI) listerMembres() {
	AfficherTitre("📋 LISTE DE TOUS LES MEMBRES")

	membres := cli.gestionnaireMembres.ListerMembres()

	if len(membres) == 0 {
		AfficherInfo("Aucun membre inscrit.")
		return
	}

	cli.afficherTableauMembres(membres)
}

func (cli *CLI) listerMembresActifs() {
	AfficherTitre("✅ MEMBRES ACTIFS")

	membres := cli.gestionnaireMembres.ListerMembresActifs()

	if len(membres) == 0 {
		AfficherInfo("Aucun membre actif.")
		return
	}

	cli.afficherTableauMembres(membres)
}

func (cli *CLI) rechercherMembres() {
	AfficherTitre("🔍 RECHERCHER DES MEMBRES")

	terme := LireEntreeObligatoire("Terme de recherche (nom ou email) : ")

	resultats := cli.gestionnaireMembres.RechercherMembres(terme)

	fmt.Printf("\n🎯 %d résultat(s) trouvé(s) pour '%s' :\n", len(resultats), terme)

	if len(resultats) == 0 {
		AfficherInfo("Aucun membre correspondant.")
		return
	}

	cli.afficherTableauMembres(resultats)
}

func (cli *CLI) modifierMembre() error {
	AfficherTitre("✏️ MODIFIER UN MEMBRE")

	id := LireEntreeEntierObligatoire("ID du membre à modifier : ")

	membre, _ := cli.gestionnaireMembres.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	// Afficher les informations actuelles
	fmt.Println("\nInformations actuelles :")
	membre.AfficherDetails()

	AfficherInfo("Laissez vide pour conserver la valeur actuelle.")

	// Saisir les nouvelles valeurs
	fmt.Printf("Nouveau nom (%s) : ", membre.Nom)
	nouveauNom := LireEntree()

	fmt.Printf("Nouvel email (%s) : ", membre.Email)
	nouvelEmail := LireEntree()

	fmt.Printf("Nouveau téléphone (%s) : ", membre.Telephone)
	nouveauTelephone := LireEntree()

	err := cli.gestionnaireMembres.ModifierMembre(id, nouveauNom, nouvelEmail, nouveauTelephone)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Membre ID %d modifié avec succès !", id))
	return nil
}

func (cli *CLI) suspendireMembre() error {
	AfficherTitre("⛔ SUSPENDRE UN MEMBRE")

	id := LireEntreeEntierObligatoire("ID du membre à suspendre : ")

	membre, _ := cli.gestionnaireMembres.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	fmt.Println("\nMembre à suspendre :")
	membre.AfficherDetails()

	if !LireConfirmation("\n⚠️ Êtes-vous sûr de vouloir suspendre ce membre ?") {
		AfficherInfo("Suspension annulée.")
		return nil
	}

	err := cli.gestionnaireMembres.SuspendirMembre(id)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Membre '%s' (ID: %d) suspendu avec succès !", membre.Nom, id))
	return nil
}

func (cli *CLI) reactiverMembre() error {
	AfficherTitre("✅ RÉACTIVER UN MEMBRE")

	id := LireEntreeEntierObligatoire("ID du membre à réactiver : ")

	membre, _ := cli.gestionnaireMembres.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	err := cli.gestionnaireMembres.ReactiverMembre(id)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Membre '%s' (ID: %d) réactivé avec succès !", membre.Nom, id))
	return nil
}

func (cli *CLI) supprimerMembre() error {
	AfficherTitre("🗑️ SUPPRIMER UN MEMBRE")

	id := LireEntreeEntierObligatoire("ID du membre à supprimer : ")

	membre, _ := cli.gestionnaireMembres.TrouverMembreParID(id)
	if membre == nil {
		return fmt.Errorf("aucun membre trouvé avec l'ID %d", id)
	}

	// Afficher le membre à supprimer
	fmt.Println("\nMembre à supprimer :")
	membre.AfficherDetails()

	if !LireConfirmation("\n⚠️ Êtes-vous sûr de vouloir supprimer ce membre ?") {
		AfficherInfo("Suppression annulée.")
		return nil
	}

	nom := membre.Nom
	err := cli.gestionnaireMembres.SupprimerMembre(id)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Membre '%s' (ID: %d) supprimé avec succès !", nom, id))
	return nil
}

// ========================================
// SOUS-MENU EMPRUNTS COMPLET AVEC TOUTES LES NOUVELLES FONCTIONNALITÉS
// ========================================

func (cli *CLI) menuEmprunts() error {
	for {
		AfficherTitre("📋 GESTION DES EMPRUNTS")
		fmt.Println("1. 📚 Emprunter un livre")
		fmt.Println("2. 📤 Retourner un livre")
		fmt.Println("3. 📋 Lister tous les emprunts")
		fmt.Println("4. 📘 Lister les emprunts en cours")
		fmt.Println("5. ⚠️  Lister les emprunts en retard")
		fmt.Println("6. 👤 Emprunts d'un membre")
		fmt.Println("7. 📖 Historique d'un livre")
		fmt.Println("8. 📅 Prolonger un emprunt")
		fmt.Println("9. 📅 Emprunts à rendre aujourd'hui")
		fmt.Println("10. ❌ Annuler un emprunt")
		fmt.Println("11. 📊 Rapport détaillé des emprunts")
		fmt.Println("0. ⬅️  Retour au menu principal")
		AfficherSeparateur("-", 50)

		choix := LireEntreeEntierAvecLimites("Votre choix : ", 0, 11)

		var err error
		switch choix {
		case 1:
			err = cli.emprunterLivre()
		case 2:
			err = cli.retournerLivre()
		case 3:
			cli.listerEmprunts()
		case 4:
			cli.listerEmpruntsEnCours()
		case 5:
			cli.listerEmpruntsEnRetard()
		case 6:
			cli.listerEmpruntsParMembre()
		case 7:
			cli.listerHistoriqueLivre()
		case 8:
			err = cli.prolongerEmprunt()
		case 9:
			cli.listerEmpruntsARendreAujourdhui()
		case 10:
			err = cli.annulerEmprunt()
		case 11:
			cli.genererRapportEmprunts()
		case 0:
			return nil
		}

		if err != nil {
			AfficherErreur(err.Error())
		}

		AttendreEntree("")
	}
}

func (cli *CLI) emprunterLivre() error {
	AfficherTitre("📚 EMPRUNTER UN LIVRE")

	// Afficher les livres disponibles
	livresDisponibles := cli.gestionnaireLivres.ListerLivresDisponibles()
	if len(livresDisponibles) == 0 {
		AfficherInfo("Aucun livre disponible actuellement.")
		return nil
	}

	fmt.Println("\nLivres disponibles :")
	cli.afficherTableauLivres(livresDisponibles)

	livreID := LireEntreeEntierObligatoire("\nID du livre à emprunter : ")

	// Afficher les membres actifs
	membresActifs := cli.gestionnaireMembres.ListerMembresActifs()
	if len(membresActifs) == 0 {
		AfficherInfo("Aucun membre actif.")
		return nil
	}

	fmt.Println("\nMembres actifs :")
	cli.afficherTableauMembres(membresActifs)

	membreID := LireEntreeEntierObligatoire("\nID du membre : ")

	err := cli.gestionnaireEmprunts.EmprunterLivre(livreID, membreID)
	if err != nil {
		return err
	}

	AfficherSucces("Emprunt enregistré avec succès ! 📚")
	AfficherInfo(fmt.Sprintf("Le livre doit être rendu dans %d jours.", models.DUREE_EMPRUNT_JOURS))
	return nil
}

func (cli *CLI) retournerLivre() error {
	AfficherTitre("📤 RETOURNER UN LIVRE")

	// Afficher les emprunts en cours
	empruntsEnCours := cli.gestionnaireEmprunts.ListerEmpruntsEnCours()
	if len(empruntsEnCours) == 0 {
		AfficherInfo("Aucun emprunt en cours.")
		return nil
	}

	fmt.Println("\nEmprunts en cours :")
	cli.afficherTableauEmprunts(empruntsEnCours)

	empruntID := LireEntreeEntierObligatoire("\nID de l'emprunt à retourner : ")

	err := cli.gestionnaireEmprunts.RetournerLivre(empruntID)
	if err != nil {
		return err
	}

	AfficherSucces("Retour enregistré avec succès ! 📤")
	return nil
}

func (cli *CLI) listerEmprunts() {
	AfficherTitre("📋 LISTE DE TOUS LES EMPRUNTS")

	emprunts := cli.gestionnaireEmprunts.ListerEmprunts()

	if len(emprunts) == 0 {
		AfficherInfo("Aucun emprunt enregistré.")
		return
	}

	cli.afficherTableauEmprunts(emprunts)
}

func (cli *CLI) listerEmpruntsEnCours() {
	AfficherTitre("📘 EMPRUNTS EN COURS")

	emprunts := cli.gestionnaireEmprunts.ListerEmpruntsEnCours()

	if len(emprunts) == 0 {
		AfficherInfo("Aucun emprunt en cours.")
		return
	}

	cli.afficherTableauEmprunts(emprunts)
}

func (cli *CLI) listerEmpruntsEnRetard() {
	AfficherTitre("⚠️ EMPRUNTS EN RETARD")

	emprunts := cli.gestionnaireEmprunts.ListerEmpruntsEnRetard()

	if len(emprunts) == 0 {
		AfficherSucces("Aucun emprunt en retard ! 🎉")
		return
	}

	AfficherAvertissement(fmt.Sprintf("%d emprunt(s) en retard détecté(s) !", len(emprunts)))
	cli.afficherTableauEmprunts(emprunts)

	// Afficher les détails des retards
	fmt.Println("\nDétails des retards :")
	for _, emprunt := range emprunts {
		joursRetard := emprunt.CalculerJoursRetard()
		fmt.Printf("• %s (%s) - %d jour(s) de retard\n",
			emprunt.TitreLivre, emprunt.NomMembre, joursRetard)
	}
}

func (cli *CLI) listerEmpruntsParMembre() {
	AfficherTitre("👤 EMPRUNTS D'UN MEMBRE")

	membreID := LireEntreeEntierObligatoire("ID du membre : ")

	membre, _ := cli.gestionnaireMembres.TrouverMembreParID(membreID)
	if membre == nil {
		AfficherErreur(fmt.Sprintf("Aucun membre trouvé avec l'ID %d", membreID))
		return
	}

	emprunts := cli.gestionnaireEmprunts.ListerEmpruntsParMembre(membreID)

	fmt.Printf("\nEmprunts de %s :\n", membre.Nom)

	if len(emprunts) == 0 {
		AfficherInfo("Aucun emprunt pour ce membre.")
		return
	}

	cli.afficherTableauEmprunts(emprunts)

	// Afficher un résumé
	enCours := 0
	rendus := 0
	enRetard := 0

	for _, emprunt := range emprunts {
		switch emprunt.Statut {
		case models.STATUT_EN_COURS:
			enCours++
		case models.STATUT_RENDU:
			rendus++
		case models.STATUT_EN_RETARD:
			enRetard++
		}
	}

	fmt.Printf("\nRésumé : %d total | %d en cours | %d rendus | %d en retard\n",
		len(emprunts), enCours, rendus, enRetard)
}

// ========================================
// NOUVELLES FONCTIONNALITÉS CLI
// ========================================

func (cli *CLI) listerHistoriqueLivre() {
	AfficherTitre("📖 HISTORIQUE D'UN LIVRE")

	livreID := LireEntreeEntierObligatoire("ID du livre : ")

	livre, _ := cli.gestionnaireLivres.TrouverLivreParID(livreID)
	if livre == nil {
		AfficherErreur(fmt.Sprintf("Aucun livre trouvé avec l'ID %d", livreID))
		return
	}

	emprunts := cli.gestionnaireEmprunts.ListerEmpruntsParLivre(livreID)

	fmt.Printf("\nHistorique des emprunts de '%s' :\n", livre.Titre)

	if len(emprunts) == 0 {
		AfficherInfo("Ce livre n'a jamais été emprunté.")
		return
	}

	cli.afficherTableauEmprunts(emprunts)

	// Afficher des statistiques sur ce livre
	fmt.Printf("\nStatistiques : %d emprunt(s) au total\n", len(emprunts))

	// Calculer la durée moyenne d'emprunt pour ce livre
	var totalJours int
	var count int

	for _, emprunt := range emprunts {
		if emprunt.DateRetourEffectif != nil {
			duree := emprunt.DateRetourEffectif.Sub(emprunt.DateEmprunt)
			totalJours += int(duree.Hours() / 24)
			count++
		}
	}

	if count > 0 {
		dureeM := float64(totalJours) / float64(count)
		fmt.Printf("Durée moyenne d'emprunt : %.1f jours\n", dureeM)
	}
}

func (cli *CLI) prolongerEmprunt() error {
	AfficherTitre("📅 PROLONGER UN EMPRUNT")

	empruntsEnCours := cli.gestionnaireEmprunts.ListerEmpruntsEnCours()
	if len(empruntsEnCours) == 0 {
		AfficherInfo("Aucun emprunt en cours à prolonger.")
		return nil
	}

	fmt.Println("\nEmprunts en cours :")
	cli.afficherTableauEmprunts(empruntsEnCours)

	empruntID := LireEntreeEntierObligatoire("\nID de l'emprunt à prolonger : ")

	// Vérifier que l'emprunt existe
	emprunt, _ := cli.gestionnaireEmprunts.TrouverEmpruntParID(empruntID)
	if emprunt == nil {
		return fmt.Errorf("aucun emprunt trouvé avec l'ID %d", empruntID)
	}

	// Afficher les détails de l'emprunt
	fmt.Println("\nEmprunt à prolonger :")
	emprunt.AfficherDetails()

	jours := LireEntreeEntierAvecLimites("\nNombre de jours supplémentaires (1-30) : ", 1, 30)

	// Demander confirmation
	if !LireConfirmation(fmt.Sprintf("Confirmer la prolongation de %d jour(s) ?", jours)) {
		AfficherInfo("Prolongation annulée.")
		return nil
	}

	err := cli.gestionnaireEmprunts.PrologerEmprunt(empruntID, jours)
	if err != nil {
		return err
	}

	AfficherSucces(fmt.Sprintf("Emprunt prolongé de %d jour(s) ! 📅", jours))

	// Afficher la nouvelle date limite
	empruntMisAJour, _ := cli.gestionnaireEmprunts.TrouverEmpruntParID(empruntID)
	if empruntMisAJour != nil {
		AfficherInfo(fmt.Sprintf("Nouvelle date limite : %s",
			empruntMisAJour.DateRetourPrevu.Format("02/01/2006")))
	}

	return nil
}

func (cli *CLI) listerEmpruntsARendreAujourdhui() {
	AfficherTitre("📅 EMPRUNTS À RENDRE AUJOURD'HUI")

	emprunts := cli.gestionnaireEmprunts.ObtenirEmpruntsARendreAujourdhui()

	if len(emprunts) == 0 {
		AfficherSucces("Aucun emprunt à rendre aujourd'hui ! 🎉")
		return
	}

	AfficherAvertissement(fmt.Sprintf("%d emprunt(s) à rendre aujourd'hui :", len(emprunts)))
	cli.afficherTableauEmprunts(emprunts)

	fmt.Println("\nRappel : Ces emprunts deviennent en retard à partir de demain !")
}

func (cli *CLI) annulerEmprunt() error {
	AfficherTitre("❌ ANNULER UN EMPRUNT")

	empruntsEnCours := cli.gestionnaireEmprunts.ListerEmpruntsEnCours()
	if len(empruntsEnCours) == 0 {
		AfficherInfo("Aucun emprunt en cours à annuler.")
		return nil
	}

	fmt.Println("\nEmprunts en cours :")
	cli.afficherTableauEmprunts(empruntsEnCours)

	empruntID := LireEntreeEntierObligatoire("\nID de l'emprunt à annuler : ")

	// Vérifier que l'emprunt existe
	emprunt, _ := cli.gestionnaireEmprunts.TrouverEmpruntParID(empruntID)
	if emprunt == nil {
		return fmt.Errorf("aucun emprunt trouvé avec l'ID %d", empruntID)
	}

	// Afficher les détails de l'emprunt
	fmt.Println("\nEmprunt à annuler :")
	emprunt.AfficherDetails()

	AfficherAvertissement("⚠️ ATTENTION : L'annulation d'un emprunt est une action administrative exceptionnelle.")
	AfficherInfo("Le livre redeviendra disponible et les compteurs du membre seront mis à jour.")

	// Demander confirmation avec double vérification
	if !LireConfirmation("Êtes-vous sûr de vouloir annuler cet emprunt ?") {
		AfficherInfo("Annulation de l'emprunt abandonnée.")
		return nil
	}

	// Deuxième confirmation
	fmt.Print("Pour confirmer, tapez 'ANNULER' en majuscules : ")
	confirmation := LireEntree()
	if confirmation != "ANNULER" {
		AfficherInfo("Annulation de l'emprunt abandonnée.")
		return nil
	}

	err := cli.gestionnaireEmprunts.AnnulerEmprunt(empruntID)
	if err != nil {
		return err
	}

	AfficherSucces("Emprunt annulé avec succès ! ❌")
	AfficherInfo("Le livre est maintenant disponible pour un nouvel emprunt.")

	return nil
}

func (cli *CLI) genererRapportEmprunts() {
	AfficherTitre("📊 RAPPORT DÉTAILLÉ DES EMPRUNTS")

	rapport := cli.gestionnaireEmprunts.ExporterRapportEmprunts()
	fmt.Println(rapport)

	// Afficher des statistiques supplémentaires
	stats := cli.gestionnaireEmprunts.ObtenirStatistiques()

	fmt.Println("\n=== STATISTIQUES AVANCÉES ===")

	if duree, ok := stats["duree_moyenne_jours"].(float64); ok && duree > 0 {
		fmt.Printf("Durée moyenne des emprunts : %.1f jours\n", duree)
	}

	if membrePlusActifMap, ok := stats["membre_plus_actif"].(map[string]interface{}); ok {
		nom := membrePlusActifMap["nom"].(string)
		count := membrePlusActifMap["emprunts"].(int)
		fmt.Printf("Membre le plus actif : %s (%d emprunts)\n", nom, count)
	}

	if livrePlusEmprunteMap, ok := stats["livre_plus_emprunte"].(map[string]interface{}); ok {
		titre := livrePlusEmprunteMap["titre"].(string)
		count := livrePlusEmprunteMap["emprunts"].(int)
		fmt.Printf("Livre le plus emprunté : %s (%d emprunts)\n", titre, count)
	}

	// Afficher les emprunts par mois si disponible
	if empruntsParMois, ok := stats["emprunts_par_mois"].(map[string]int); ok {
		fmt.Println("\n=== EMPRUNTS PAR MOIS (12 derniers mois) ===")
		for mois, count := range empruntsParMois {
			if count > 0 {
				fmt.Printf("%s : %d emprunt(s)\n", mois, count)
			}
		}
	}
}

// ========================================
// STATISTIQUES COMPLÈTES
// ========================================

func (cli *CLI) afficherStatistiques() {
	AfficherTitre("📊 STATISTIQUES COMPLÈTES DE LA LIBRAIRIE")

	// Statistiques des livres
	statsLivres := cli.gestionnaireLivres.ObtenirStatistiques()
	fmt.Printf("📖 LIVRES :\n")
	fmt.Printf("   Total : %d livre(s)\n", statsLivres["total"])
	if statsLivres["total"].(int) > 0 {
		fmt.Printf("   Disponibles : %d\n", statsLivres["disponibles"])
		fmt.Printf("   Empruntés : %d\n", statsLivres["empruntes"])

		if plusEmprunte, existe := statsLivres["plus_emprunte"]; existe {
			livre := plusEmprunte.(models.Livre)
			fmt.Printf("   Plus emprunté : %s (%d emprunt(s))\n", livre.Titre, livre.NombreEmprunts)
		}

		// Afficher les genres les plus populaires
		if parGenre, existe := statsLivres["par_genre"].(map[string]int); existe {
			fmt.Println("\n   Répartition par genre :")
			for genre, count := range parGenre {
				fmt.Printf("     %s : %d livre(s)\n", genre, count)
			}
		}
	}

	fmt.Println()

	// Statistiques des membres
	statsMembres := cli.gestionnaireMembres.ObtenirStatistiques()
	fmt.Printf("👥 MEMBRES :\n")
	fmt.Printf("   Total : %d membre(s)\n", statsMembres["total"])
	if statsMembres["total"].(int) > 0 {
		fmt.Printf("   Actifs : %d\n", statsMembres["actifs"])
		fmt.Printf("   Suspendus : %d\n", statsMembres["suspendus"])

		if plusActif, existe := statsMembres["plus_actif"]; existe {
			membre := plusActif.(models.Membre)
			fmt.Printf("   Plus actif : %s (%d emprunt(s))\n", membre.Nom, membre.NombreEmprunts)
		}
	}

	fmt.Println()

	// Statistiques des emprunts
	statsEmprunts := cli.gestionnaireEmprunts.ObtenirStatistiques()
	fmt.Printf("📋 EMPRUNTS :\n")
	fmt.Printf("   Total : %d emprunt(s)\n", statsEmprunts["total"])
	if statsEmprunts["total"].(int) > 0 {
		fmt.Printf("   En cours : %d\n", statsEmprunts["en_cours"])
		fmt.Printf("   Rendus : %d\n", statsEmprunts["rendus"])
		fmt.Printf("   En retard : %d\n", statsEmprunts["en_retard"])

		if duree, existe := statsEmprunts["duree_moyenne_jours"].(float64); existe && duree > 0 {
			fmt.Printf("   Durée moyenne : %.1f jours\n", duree)
		}
	}

	// Alertes et recommandations
	fmt.Println("\n=== ALERTES ET RECOMMANDATIONS ===")

	// Vérifier les emprunts en retard
	empruntsEnRetard := cli.gestionnaireEmprunts.ListerEmpruntsEnRetard()
	if len(empruntsEnRetard) > 0 {
		AfficherAvertissement(fmt.Sprintf("%d emprunt(s) en retard nécessitent un suivi", len(empruntsEnRetard)))
	}

	// Vérifier les emprunts à rendre aujourd'hui
	empruntsAujourdhui := cli.gestionnaireEmprunts.ObtenirEmpruntsARendreAujourdhui()
	if len(empruntsAujourdhui) > 0 {
		AfficherInfo(fmt.Sprintf("%d emprunt(s) à rendre aujourd'hui", len(empruntsAujourdhui)))
	}

	// Taux d'occupation de la librairie
	if statsLivres["total"].(int) > 0 {
		tauxOccupation := float64(statsLivres["empruntes"].(int)) / float64(statsLivres["total"].(int)) * 100
		fmt.Printf("\n📈 Taux d'occupation : %.1f%% des livres sont actuellement empruntés\n", tauxOccupation)

		if tauxOccupation > 80 {
			AfficherInfo("Excellente fréquentation ! Considérez l'ajout de nouveaux livres.")
		} else if tauxOccupation < 20 {
			AfficherInfo("Faible taux d'emprunt. Envisagez des actions de promotion.")
		}
	}
}

// ========================================
// FONCTIONS HELPER POUR L'AFFICHAGE DES TABLEAUX
// ========================================

func (cli *CLI) afficherTableauLivres(livres []models.Livre) {
	fmt.Printf("\n")
	fmt.Printf("│ %-3s │ %-25s │ %-20s │ %-15s │ %-12s │\n", "ID", "Titre", "Auteur", "Genre", "Statut")
	fmt.Printf("├%s┼%s┼%s┼%s┼%s┤\n",
		strings.Repeat("─", 5),
		strings.Repeat("─", 27),
		strings.Repeat("─", 22),
		strings.Repeat("─", 17),
		strings.Repeat("─", 14))

	for _, livre := range livres {
		titre := livre.Titre
		if len(titre) > 25 {
			titre = titre[:22] + "..."
		}

		auteur := livre.Auteur
		if len(auteur) > 20 {
			auteur = auteur[:17] + "..."
		}

		genre := livre.Genre
		if len(genre) > 15 {
			genre = genre[:12] + "..."
		}

		statut := "📗 Disponible"
		if !livre.EstDisponible() {
			statut = "📕 Emprunté"
		}

		fmt.Printf("│ %-3d │ %-25s │ %-20s │ %-15s │ %-12s │\n",
			livre.ID, titre, auteur, genre, statut)
	}

	fmt.Printf("└%s┴%s┴%s┴%s┴%s┘\n",
		strings.Repeat("─", 5),
		strings.Repeat("─", 27),
		strings.Repeat("─", 22),
		strings.Repeat("─", 17),
		strings.Repeat("─", 14))

	fmt.Printf("\nTotal : %d livre(s)\n", len(livres))
}

func (cli *CLI) afficherTableauMembres(membres []models.Membre) {
	fmt.Printf("\n")
	fmt.Printf("│ %-3s │ %-25s │ %-25s │ %-9s │ %-12s │\n", "ID", "Nom", "Email", "Emprunts", "Statut")
	fmt.Printf("├%s┼%s┼%s┼%s┼%s┤\n",
		strings.Repeat("─", 5),
		strings.Repeat("─", 27),
		strings.Repeat("─", 27),
		strings.Repeat("─", 11),
		strings.Repeat("─", 14))

	for _, membre := range membres {
		nom := membre.Nom
		if len(nom) > 25 {
			nom = nom[:22] + "..."
		}

		email := membre.Email
		if len(email) > 25 {
			email = email[:22] + "..."
		}

		emprunts := fmt.Sprintf("%d/%d", membre.EmpruntsActifs, models.LIMITE_EMPRUNTS_SIMULTANES)

		statut := "✅ Actif"
		if !membre.Actif {
			statut = "❌ Suspendu"
		}

		fmt.Printf("│ %-3d │ %-25s │ %-25s │ %-9s │ %-12s │\n",
			membre.ID, nom, email, emprunts, statut)
	}

	fmt.Printf("└%s┴%s┴%s┴%s┴%s┘\n",
		strings.Repeat("─", 5),
		strings.Repeat("─", 27),
		strings.Repeat("─", 27),
		strings.Repeat("─", 11),
		strings.Repeat("─", 14))

	fmt.Printf("\nTotal : %d membre(s)\n", len(membres))
}

func (cli *CLI) afficherTableauEmprunts(emprunts []models.Emprunt) {
	fmt.Printf("\n")
	fmt.Printf("│ %-3s │ %-25s │ %-20s │ %-10s │ %-12s │\n", "ID", "Livre", "Membre", "Emprunté", "Statut")
	fmt.Printf("├%s┼%s┼%s┼%s┼%s┤\n",
		strings.Repeat("─", 5),
		strings.Repeat("─", 27),
		strings.Repeat("─", 22),
		strings.Repeat("─", 12),
		strings.Repeat("─", 14))

	for _, emprunt := range emprunts {
		titre := emprunt.TitreLivre
		if len(titre) > 25 {
			titre = titre[:22] + "..."
		}

		nom := emprunt.NomMembre
		if len(nom) > 20 {
			nom = nom[:17] + "..."
		}

		dateEmprunt := emprunt.DateEmprunt.Format("02/01/2006")

		var statut string
		switch emprunt.Statut {
		case models.STATUT_EN_COURS:
			statut = "📘 En cours"
		case models.STATUT_RENDU:
			statut = "✅ Rendu"
		case models.STATUT_EN_RETARD:
			joursRetard := emprunt.CalculerJoursRetard()
			statut = fmt.Sprintf("⚠️ %d j retard", joursRetard)
		default:
			statut = emprunt.Statut
		}

		fmt.Printf("│ %-3d │ %-25s │ %-20s │ %-10s │ %-12s │\n",
			emprunt.ID, titre, nom, dateEmprunt, statut)
	}

	fmt.Printf("└%s┴%s┴%s┴%s┴%s┘\n",
		strings.Repeat("─", 5),
		strings.Repeat("─", 27),
		strings.Repeat("─", 22),
		strings.Repeat("─", 12),
		strings.Repeat("─", 14))

	fmt.Printf("\nTotal : %d emprunt(s)\n", len(emprunts))
}
