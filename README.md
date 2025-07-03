# 📚 Système de Gestion de Librairie

Un système complet de gestion de librairie développé en Go avec une architecture modulaire et claire.

## 🚀 Fonctionnalités

### 📖 Gestion des Livres
- ➕ Ajouter de nouveaux livres avec validation ISBN
- 📋 Lister tous les livres ou seulement les disponibles
- 🔍 Rechercher par titre, auteur ou genre
- ✏️ Modifier les informations d'un livre
- 🗑️ Supprimer des livres (si non empruntés)

### 👥 Gestion des Membres
- ➕ Inscrire de nouveaux membres
- 📋 Gérer les membres actifs et suspendus
- ✏️ Modifier les informations des membres
- ⛔ Suspendre/réactiver des membres
- 📊 Limite de 3 emprunts simultanés par membre

### 📋 Gestion des Emprunts
- 📚 Emprunter des livres (durée : 14 jours)
- 📤 Retourner des livres
- ⚠️ Détection automatique des retards
- 📊 Historique complet des emprunts
- 👤 Consulter les emprunts par membre

### 📊 Statistiques
- Livres les plus empruntés
- Membres les plus actifs  
- Emprunts en retard
- Données globales de la librairie

## 🏗️ Architecture