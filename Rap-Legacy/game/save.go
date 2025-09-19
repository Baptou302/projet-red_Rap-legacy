package game // Déclare le package "game", utilisé pour organiser le code

import (
	"encoding/json" // Pour encoder et décoder les structures en JSON
	"errors"        // Pour créer des erreurs personnalisées
	"os"            // Pour lire/écrire et manipuler fichiers
	"path/filepath" // Pour gérer les chemins de fichiers de manière portable
	"time"          // Pour gérer le temps et timestamps
)

// -----------------------------
// Type Save (utilisé par game.go)
// -----------------------------
type Save struct {
	Name      string   `json:"name"`         // Nom de la sauvegarde
	Class     string   `json:"class"`        // Classe du joueur
	Inventory []string `json:"inventory"`    // Liste des objets possédés
	Created   int64    `json:"created_unix"` // Timestamp Unix de création
	// Position du joueur
	PlayerX float64 `json:"player_x"` // Coordonnée X du joueur
	PlayerY float64 `json:"player_y"` // Coordonnée Y du joueur
	// Stats du joueur
	Ego      int `json:"ego"`      // Ego du joueur
	Flow     int `json:"flow"`     // Flow du joueur
	Charisma int `json:"charisma"` // Charisme du joueur
}

// -----------------------------
// Chemins / constantes
// -----------------------------
const savesDir = "saves"       // Répertoire où sont stockées les sauvegardes
const savesFile = "saves.json" // Nom du fichier global contenant toutes les saves

// -----------------------------
// Helpers FS
// -----------------------------
func ensureSavesPath() error {
	// Vérifie si le dossier de sauvegarde existe
	if _, err := os.Stat(savesDir); os.IsNotExist(err) {
		// Si non, le crée avec les permissions 755
		if err := os.MkdirAll(savesDir, 0o755); err != nil {
			return err // Retourne une erreur si échec
		}
	}
	// Vérifie si le fichier JSON existe
	path := filepath.Join(savesDir, savesFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Si non, crée un fichier vide initialisé avec "[]"
		if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
			return err
		}
	}
	return nil // Tout est ok
}

// Retourne le chemin complet du fichier de sauvegarde global
func saveFilePath() string {
	return filepath.Join(savesDir, savesFile)
}

// -----------------------------
// Lecture / écriture (liste unique saves.json)
// -----------------------------
func LoadAllSaves() ([]Save, error) {
	if err := ensureSavesPath(); err != nil { // Assure que le dossier/fichier existent
		return nil, err
	}
	path := saveFilePath()         // Chemin complet du fichier
	data, err := os.ReadFile(path) // Lecture du contenu du fichier
	if err != nil {                // Si erreur lors de la lecture
		return nil, err
	}
	var saves []Save    // Slice pour stocker les saves
	if len(data) == 0 { // Si fichier vide
		return []Save{}, nil // Retourne slice vide
	}
	if err := json.Unmarshal(data, &saves); err != nil { // Décodage JSON
		return []Save{}, err // Erreur si JSON invalide
	}
	return saves, nil // Retourne la liste de saves
}

// Écrit toutes les sauvegardes dans le fichier JSON global
func SaveAll(saves []Save) error {
	if err := ensureSavesPath(); err != nil { // Assure existence du dossier/fichier
		return err
	}
	path := saveFilePath()                        // Chemin du fichier
	d, err := json.MarshalIndent(saves, "", "  ") // Encode en JSON lisible
	if err != nil {                               // Si erreur d'encodage
		return err
	}
	return os.WriteFile(path, d, 0o644) // Écrit le JSON dans le fichier
}

// -----------------------------
// Fonctions CRUD (liste)
// -----------------------------

// GetSave retourne la save correspondant au nom exact
func GetSave(name string) (Save, bool, error) {
	saves, err := LoadAllSaves() // Charge toutes les saves
	if err != nil {              // Si erreur
		return Save{}, false, err
	}
	for _, s := range saves { // Parcourt toutes les saves
		if s.Name == name { // Si nom correspond
			return s, true, nil // Retourne la save trouvée
		}
	}
	return Save{}, false, nil // Non trouvé
}

// Vérifie si une sauvegarde existe
func SaveExists(name string) (bool, error) {
	_, ok, err := GetSave(name) // Appelle GetSave
	return ok, err              // Retourne vrai si trouvé, faux sinon
}

// CreateSave : crée une nouvelle save et l'ajoute à saves.json
func CreateSave(name, class string) (Save, error) {
	if name == "" { // Vérifie que le nom n'est pas vide
		return Save{}, errors.New("nom de sauvegarde vide")
	}

	// Vérifie si la save existe déjà
	exists, err := SaveExists(name)
	if err != nil {
		return Save{}, err
	}
	if exists { // Si existe déjà
		return Save{}, errors.New("une sauvegarde avec ce nom existe déjà")
	}

	// Détermine l'inventaire initial selon la classe
	var inv []string
	switch class {
	case "Lyricistes", "lyricistes", "lyriciste":
		inv = []string{"Micro", "Cristalline - mystérieuse", "Cigarette électronique"}
	case "Performeurs", "performeurs", "performer":
		inv = []string{"Micro", "Cristalline - tonic", "Téléphone"}
	case "Hitmakers", "hitmakers", "hitmaker":
		inv = []string{"Micro", "Cristalline - suspicieuse", "Téléphone"}
	default:
		inv = []string{"Micro"} // Inventaire par défaut
	}

	now := time.Now().Unix() // Timestamp actuel
	newSave := Save{
		Name:      name,  // Nom
		Class:     class, // Classe
		Inventory: inv,   // Inventaire
		Created:   now,   // Date de création
		PlayerX:   100,   // Position X initiale
		PlayerY:   100,   // Position Y initiale
		Ego:       100,   // Stat Ego initial
		Flow:      10,    // Stat Flow initial
		Charisma:  5,     // Stat Charisma initial
	}

	saves, err := LoadAllSaves() // Charge toutes les saves existantes
	if err != nil {
		return Save{}, err
	}
	saves = append(saves, newSave)         // Ajoute la nouvelle save
	if err := SaveAll(saves); err != nil { // Sauvegarde le tout dans le fichier
		return Save{}, err
	}
	return newSave, nil // Retourne la save créée
}

// OverwriteSave remplace une save existante ou l'ajoute si inexistante
func OverwriteSave(updated Save) error {
	saves, err := LoadAllSaves() // Charge toutes les saves
	if err != nil {
		return err
	}
	found := false
	for i := range saves { // Parcourt toutes les saves
		if saves[i].Name == updated.Name { // Si match par nom
			saves[i] = updated // Remplace
			found = true
			break
		}
	}
	if !found { // Si non trouvée
		saves = append(saves, updated) // Ajoute la save
	}
	return SaveAll(saves) // Écrit dans le fichier JSON
}

// DeleteSave supprime une save du fichier global
func DeleteSave(name string) error {
	saves, err := LoadAllSaves() // Charge toutes les saves
	if err != nil {
		return err
	}
	newSaves := make([]Save, 0, len(saves)) // Nouvelle slice pour sauvegardes restantes
	found := false
	for _, s := range saves { // Parcourt toutes les saves
		if s.Name == name { // Ignore celle à supprimer
			found = true
			continue
		}
		newSaves = append(newSaves, s) // Ajoute les autres
	}
	if !found { // Si non trouvée
		return errors.New("sauvegarde introuvable")
	}
	return SaveAll(newSaves) // Sauvegarde les autres
}

// ListSaves : alias pratique pour LoadAllSaves
func ListSaves() ([]Save, error) {
	return LoadAllSaves()
}

// -----------------------------
// Fonctions optionnelles (fichier individuel par save)
// -----------------------------

// Sauvegarde une save dans un fichier séparé
func SaveToFileSingle(name string, s Save) error {
	if err := ensureSavesPath(); err != nil {
		return err
	}
	path := filepath.Join(savesDir, name+".json") // Chemin du fichier individuel
	b, err := json.MarshalIndent(s, "", "  ")     // Encode la save en JSON lisible
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644) // Écrit le fichier
}

// Charge une save depuis un fichier individuel
func LoadFromFileSingle(name string) (Save, error) {
	if err := ensureSavesPath(); err != nil {
		return Save{}, err
	}
	path := filepath.Join(savesDir, name+".json") // Chemin du fichier
	b, err := os.ReadFile(path)                   // Lit le fichier
	if err != nil {
		return Save{}, err
	}
	var s Save
	if err := json.Unmarshal(b, &s); err != nil { // Décode JSON
		return Save{}, err
	}
	return s, nil // Retourne la save
}
