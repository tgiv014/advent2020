package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type ingredient struct {
	name     string
	nFoodsIn int
	all      *allergen
}

type allergen struct {
	name               string
	ingredientMentions map[string]int
	ing                *ingredient
}

type food struct {
	ingredients []*ingredient
	allergens   []*allergen
}

var reIngredients *regexp.Regexp

func init() {
	reIngredients = regexp.MustCompile(`^([a-z ]+) \(contains ([a-z ,]+)\)$`)
}
func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	foods := make([]*food, 0)
	allIngredients := make(map[string]*ingredient)
	allAllergens := make(map[string]*allergen)

	// Read in the seed
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := reIngredients.FindSubmatch(scanner.Bytes())
		primaryIngredients := strings.Split(string(m[1]), " ")
		allergens := strings.Split(string(m[2]), ", ")
		myIngredients := []*ingredient{}
		myAllergens := []*allergen{}
		for _, ing := range primaryIngredients {
			// Only insert new ingredient to the `all` list if needed
			if _, ok := allIngredients[ing]; !ok {
				allIngredients[ing] = &ingredient{ing, 0, nil}
			}
			myIngredients = append(myIngredients, allIngredients[ing])
			allIngredients[ing].nFoodsIn++
		}

		for _, all := range allergens {
			// Only insert new ingredient to the `all` list if needed
			if _, ok := allAllergens[all]; !ok {
				allAllergens[all] = &allergen{all, make(map[string]int), nil}
			}
			myAllergens = append(myAllergens, allAllergens[all])
		}
		newFood := food{myIngredients, myAllergens}
		foods = append(foods, &newFood)
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	for _, fd := range foods {
		for _, allergen := range fd.allergens {
			for _, ing := range fd.ingredients {
				allergen.ingredientMentions[ing.name]++
			}
		}
	}

	// Loop until all allergens have been assigned ingredients
	nFound := 0
	for nFound < len(allAllergens) {
		for _, all := range allAllergens {
			// Skip if already assigned
			if all.ing != nil {
				continue
			}
			// Get max number of mentions for this allergen
			max := 0
			for _, v := range all.ingredientMentions {
				if v > max {
					max = v
				}
			}
			// See if there is only one maximum
			nMaximums := 0
			var maxKey string
			for k, v := range all.ingredientMentions {
				if v == max {
					nMaximums++
					maxKey = k
				}
			}
			// If there is more than one maximum, skip for now
			if nMaximums > 1 {
				continue
			}
			// Success!
			nFound++
			all.ing = allIngredients[maxKey]
			allIngredients[maxKey].all = all
			// Delete the newly found ingredient from other allergens' mention lists
			for ak, av := range allAllergens {
				if ak == all.name {
					continue
				}
				delete(av.ingredientMentions, maxKey)
			}
			break
		}
	}

	for _, all := range allAllergens {
		fmt.Println(all.name, "->", all.ing.name)
	}

	nSafeMentions := 0
	unsafeIngredients := []*ingredient{}
	for _, ing := range allIngredients {
		if ing.all == nil {
			nSafeMentions += ing.nFoodsIn
		} else {
			unsafeIngredients = append(unsafeIngredients, ing)
		}
	}
	fmt.Println("Safe Mentions:", nSafeMentions)
	sort.Slice(unsafeIngredients, func(i, j int) bool {
		return unsafeIngredients[i].all.name < unsafeIngredients[j].all.name
	})
	unsafeIngredientNames := []string{}
	for _, ing := range unsafeIngredients {
		unsafeIngredientNames = append(unsafeIngredientNames, ing.name)
	}
	fmt.Println("Canonical Unsafe Ingedients:", strings.Join(unsafeIngredientNames, ","))
}
