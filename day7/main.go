package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node struct {
	name     string
	count    int
	children []Node
}

func create_node(name string, count int) Node {
	return Node{name, count, make([]Node, 0)}
}

func (n *Node) add_child(name string, count int) {
	new_child := create_node(name, count)
	n.children = append(n.children, new_child)
}

func (n *Node) build_tree(rules map[string]Node) {
	// Loop through this rule's children
	for i, child := range n.children {
		// Find the appropraite rule for this child and copy it's values in
		rule := rules[child.name]
		for _, rule_child := range rule.children {
			n.children[i].add_child(rule_child.name, rule_child.count)
		}
		n.children[i].build_tree(rules)
	}
}

func (n *Node) contains(name string) bool {
	if len(n.children) == 0 {
		return false
	}
	for _, child := range n.children {
		if child.name == name {
			return true
		}
		if child.contains(name) {
			return true
		}
	}
	return false
}

func (n *Node) n_bags_inside() int {
	if len(n.children) == 0 {
		return 0 // If we don't have children, we contain no bags
	}
	total := 0
	for _, node := range n.children {
		total += node.count
		total += node.count * node.n_bags_inside()
	}
	return total
}

func parse_line(line string) (Node, string) {
	tokens := strings.Split(line, " contain ")
	parent := strings.Replace(tokens[0], " bags", "", 1)
	parent = strings.Trim(parent, " ,.")
	// fmt.Println(parent)
	parent_node := create_node(parent, 1)
	child_tokens := strings.Split(tokens[1], ", ")
	for _, t := range child_tokens {
		t = strings.Replace(t, "bags", "", 1)
		t = strings.Replace(t, "bag", "", 1)
		t = strings.Trim(t, " ,.")
		if t == "no other" {
			continue
		}
		count := int(t[0]) - 48
		child_name := t[2:]
		parent_node.add_child(child_name, count)
	}
	return parent_node, parent
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatalf("Uh oh %e", err)
	}
	defer f.Close()

	rules := make(map[string]Node)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		newrule, name := parse_line(line)
		rules[name] = newrule
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	total_w_gold := 0
	for _, rule := range rules {
		rule.build_tree(rules)
		contains_shiny_gold := rule.contains("shiny gold")
		// fmt.Println(rule.name, contains_shiny_gold)
		if contains_shiny_gold {
			total_w_gold += 1
		}
	}
	fmt.Println("# w shiny gold", total_w_gold)
	// fmt.Println(rules)
	shiny_gold := rules["shiny gold"]
	fmt.Println("# bags in shiny gold", shiny_gold.n_bags_inside())
}
