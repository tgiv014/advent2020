package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type encryptor struct {
	key      int
	loopsize int
}

func (e *encryptor) transform(subject int) int {
	return transform(e.key, subject, e.loopsize)
}

func (e *encryptor) findLoopSize() int {
	var i int
	res := 1
	for i = 0; res != e.key; i++ {
		res *= 7
		res %= 20201227
	}
	e.loopsize = i
	return i
}

func transform(key, subject, loopSize int) int {
	var res int = 1
	for i := 0; i < loopSize; i++ {
		res *= subject
		res %= 20201227
	}
	return res
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	cardPubKey := encryptor{}
	doorPubKey := encryptor{}
	scanner := bufio.NewScanner(f)
	for lineNum := 0; scanner.Scan(); lineNum++ {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		switch lineNum {
		case 0:
			cardPubKey.key = n
		case 1:
			doorPubKey.key = n
		default:
			fmt.Println("Unprocessed lines")
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	fmt.Println(cardPubKey.findLoopSize())
	fmt.Println(doorPubKey.findLoopSize())
	fmt.Println("Found loop sizes")

	encKey1 := cardPubKey.transform(doorPubKey.key)
	encKey2 := doorPubKey.transform(cardPubKey.key)
	if encKey1 == encKey2 {
		fmt.Println("Verified encryption key", encKey1)
	} else {
		fmt.Println("Issue encountered in calculating encryption key")
	}
}
