package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Inventory struct {
	Items map[string]int
}

func NewInventory() *Inventory {
	return &Inventory{
		Items: make(map[string]int),
	}
}

func (inv *Inventory) AddItem(item string, qty int) {
	inv.Items[item] += qty
}

func (inv *Inventory) RemoveItem(item string, qty int) bool {
	if inv.Items[item] >= qty {
		inv.Items[item] -= qty
		if inv.Items[item] == 0 {
			delete(inv.Items, item)
		}
		return true
	}
	return false
}

func (inv *Inventory) HasItems(requirements map[string]int) bool {
	for item, qty := range requirements {
		if inv.Items[item] < qty {
			return false
		}
	}
	return true
}

// Recette de craft
type Recipe struct {
	Name         string
	Requirements map[string]int
	Result       string
}

// Craft d’une recette
func Craft(inv *Inventory, recipe Recipe) bool {
	if !inv.HasItems(recipe.Requirements) {
		fmt.Println("❌ Tu n'as pas les objets nécessaires pour crafter:", recipe.Name)
		return false
	}

	for item, qty := range recipe.Requirements {
		inv.RemoveItem(item, qty)
	}

	inv.AddItem(recipe.Result, 1)
	fmt.Println("✅ Tu as crafté:", recipe.Result)
	return true
}

// Menu de craft
func CraftMenu(inv *Inventory, recipes []Recipe) {
	fmt.Println("\n===== ⚒️  MENU DE CRAFT ⚒️  =====")
	for i, recipe := range recipes {
		fmt.Printf("[%d] %s (nécessite:", i+1, recipe.Name)
		for item, qty := range recipe.Requirements {
			fmt.Printf(" %dx %s", qty, item)
		}
		fmt.Println(")")
	}
	fmt.Println("[0] Quitter le menu")

	fmt.Print("Choisis un craft: ")
	reader := bufio.NewReader(os.Stdin)
	choiceRaw, _ := reader.ReadString('\n')
	choiceRaw = strings.TrimSpace(choiceRaw)

	if choiceRaw == "0" {
		fmt.Println("Retour au jeu...")
		return
	}

	if choiceRaw == "1" && len(recipes) >= 1 {
		Craft(inv, recipes[0])
	}
}

func main() {
	inv := NewInventory()
	inv.AddItem("Cristalline", 1)
	inv.AddItem("Cigarette Électronique", 1)

	recipes := []Recipe{
		{
			Name: "Méga Cristalline",
			Requirements: map[string]int{
				"Cristalline":            1,
				"Cigarette Électronique": 1,
			},
			Result: "Méga Cristalline",
		},
	}

	fmt.Println("Appuie sur 'f' pour ouvrir le menu de craft.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "f" {
			CraftMenu(inv, recipes)
			fmt.Println("Inventaire actuel:", inv.Items)
		}
	}
}
