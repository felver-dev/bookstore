package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Storage interface {
	Sauvegarder(data any) error
	Charger(data any) error
}

type JSONStorage struct {
	filename string
}

func NewJSONStorage(filename string) *JSONStorage {
	return &JSONStorage{filename: filename}
}

func (js *JSONStorage) Sauvegarder(data any) error {
	dir := filepath.Dir(js.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("imppossible de créer le dossier %s : %v", dir, err)
	}

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("erreur lors de la conversion en JSON : %v", err)
	}

	err = os.WriteFile(js.filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de la l'écriture du fichier %s : %v", js.filename, err)
	}

	return nil
}

func (js *JSONStorage) Charger(data any) error {
	if _, err := os.Stat(js.filename); os.IsNotExist(err) {
		return nil
	}

	fileData, err := os.ReadFile(js.filename)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier %s : %v", js.filename, err)
	}

	if len(fileData) == 0 {
		return nil
	}

	err = json.Unmarshal(fileData, data)
	if err != nil {
		return fmt.Errorf("erreur lors de la conversion depuis JSON : %v", err)
	}
	return nil
}

func (js *JSONStorage) Existe() bool {

	_, err := os.Stat(js.filename)
	return !os.IsNotExist(err)
}

func (js *JSONStorage) Supprimer() error {
	if !js.Existe() {
		return nil
	}
	err := os.Remove(js.filename)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du fichier %s : %v", js.filename, err)
	}
	return nil
}

func (js *JSONStorage) TailleEnOctets() (int64, error) {
	if !js.Existe() {
		return 0, nil
	}

	info, err := os.Stat(js.filename)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de la lecture des informations  du fichier : %v", err)
	}

	return info.Size(), nil
}
