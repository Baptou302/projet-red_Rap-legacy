package game

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// -----------------------------
// Type Save (utilisé par game.go)
// -----------------------------
type Save struct {
	Name      string   `json:"name"`
	Class     string   `json:"class"`
	Inventory []string `json:"inventory"`
	Created   int64    `json:"created_unix"`
	// Position (float64 pour être compatible avec NewPlayer)
	PlayerX float64 `json:"player_x"`
	PlayerY float64 `json:"player_y"`
	// Stats
	Ego      int `json:"ego"`
	Flow     int `json:"flow"`
	Charisma int `json:"charisma"`
}

// -----------------------------
// Chemins / constantes
// -----------------------------
const savesDir = "saves"
const savesFile = "saves.json"

// -----------------------------
// Helpers FS
// -----------------------------
func ensureSavesPath() error {
	// Crée dossier si nécessaire
	if _, err := os.Stat(savesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(savesDir, 0o755); err != nil {
			return err
		}
	}
	// Crée le fichier JSON s'il n'existe pas (contenu initial "[]")
	path := filepath.Join(savesDir, savesFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func saveFilePath() string {
	return filepath.Join(savesDir, savesFile)
}

// -----------------------------
// Lecture / écriture (liste unique saves.json)
// -----------------------------
func LoadAllSaves() ([]Save, error) {
	if err := ensureSavesPath(); err != nil {
		return nil, err
	}
	path := saveFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var saves []Save
	if len(data) == 0 {
		return []Save{}, nil
	}
	if err := json.Unmarshal(data, &saves); err != nil {
		return []Save{}, err
	}
	return saves, nil
}

func SaveAll(saves []Save) error {
	if err := ensureSavesPath(); err != nil {
		return err
	}
	path := saveFilePath()
	d, err := json.MarshalIndent(saves, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, d, 0o644)
}

// -----------------------------
// Fonctions CRUD (liste)
// -----------------------------

// GetSave retourne la save par nom (match exact)
func GetSave(name string) (Save, bool, error) {
	saves, err := LoadAllSaves()
	if err != nil {
		return Save{}, false, err
	}
	for _, s := range saves {
		if s.Name == name {
			return s, true, nil
		}
	}
	return Save{}, false, nil
}

func SaveExists(name string) (bool, error) {
	_, ok, err := GetSave(name)
	return ok, err
}

// CreateSave : crée une nouvelle save et l'ajoute à saves.json
// renvoie la Save créée.
func CreateSave(name, class string) (Save, error) {
	if name == "" {
		return Save{}, errors.New("nom de sauvegarde vide")
	}

	// vérifie existence
	exists, err := SaveExists(name)
	if err != nil {
		return Save{}, err
	}
	if exists {
		return Save{}, errors.New("une sauvegarde avec ce nom existe déjà")
	}

	// inventaire par classe
	var inv []string
	switch class {
	case "Lyricistes", "lyricistes", "lyriciste":
		inv = []string{"Micro", "Cristalline - mystérieuse", "Cigarette électronique"}
	case "Performeurs", "performeurs", "performer":
		inv = []string{"Micro", "Cristalline - tonic", "Téléphone"}
	case "Hitmakers", "hitmakers", "hitmaker":
		inv = []string{"Micro", "Cristalline - suspicieuse", "Téléphone"}
	default:
		inv = []string{"Micro"}
	}

	now := time.Now().Unix()
	newSave := Save{
		Name:      name,
		Class:     class,
		Inventory: inv,
		Created:   now,
		PlayerX:   100,
		PlayerY:   100,
		Ego:       100,
		Flow:      10,
		Charisma:  5,
	}

	saves, err := LoadAllSaves()
	if err != nil {
		return Save{}, err
	}
	saves = append(saves, newSave)
	if err := SaveAll(saves); err != nil {
		return Save{}, err
	}
	return newSave, nil
}

// OverwriteSave remplace une save existante (ou l'ajoute si inexistante)
func OverwriteSave(updated Save) error {
	saves, err := LoadAllSaves()
	if err != nil {
		return err
	}
	found := false
	for i := range saves {
		if saves[i].Name == updated.Name {
			saves[i] = updated
			found = true
			break
		}
	}
	if !found {
		saves = append(saves, updated)
	}
	return SaveAll(saves)
}

// DeleteSave supprime une save du fichier
func DeleteSave(name string) error {
	saves, err := LoadAllSaves()
	if err != nil {
		return err
	}
	newSaves := make([]Save, 0, len(saves))
	found := false
	for _, s := range saves {
		if s.Name == name {
			found = true
			continue
		}
		newSaves = append(newSaves, s)
	}
	if !found {
		return errors.New("sauvegarde introuvable")
	}
	return SaveAll(newSaves)
}

// ListSaves : alias pratique
func ListSaves() ([]Save, error) {
	return LoadAllSaves()
}

// -----------------------------
// Fonctions optionnelles (fichier individuel par save)
// -----------------------------
func SaveToFileSingle(name string, s Save) error {
	if err := ensureSavesPath(); err != nil {
		return err
	}
	path := filepath.Join(savesDir, name+".json")
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func LoadFromFileSingle(name string) (Save, error) {
	if err := ensureSavesPath(); err != nil {
		return Save{}, err
	}
	path := filepath.Join(savesDir, name+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		return Save{}, err
	}
	var s Save
	if err := json.Unmarshal(b, &s); err != nil {
		return Save{}, err
	}
	return s, nil
}
