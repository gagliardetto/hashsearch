package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/hashsearch"
)

func main() {
	arr := hashsearch.NewInt()

	arr.Add(1)
	spew.Dump(arr)

	arr.Add(1)
	spew.Dump(arr)

	arr.Add(7)
	spew.Dump(arr)

	arr.Add(7)
	spew.Dump(arr)

	arr.Add(6)
	spew.Dump(arr)

	fmt.Println(arr.Has(11))
	fmt.Println(arr.Has(99))
	fmt.Println(arr.Has(0))
	fmt.Println(arr.Has(1))
	fmt.Println(arr.Has(7))
}
