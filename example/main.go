package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/hashsearch"
)

func main() {
	arr := hashsearch.NewIntArr()

	arr.Sort()
	arr.OrderedAppend(1)
	spew.Dump(arr)

	arr.WarningUnorderedAppend(10)
	arr.WarningUnorderedAppend(9)
	arr.WarningUnorderedAppend(8)
	spew.Dump(arr)

	arr.Sort()
	spew.Dump(arr)

	arr.WarningUnorderedAppend(0)
	spew.Dump(arr)

	arr.Sort()
	spew.Dump(arr)

	arr.OrderedAppend(1)
	spew.Dump(arr)

	arr.OrderedAppend(7)
	spew.Dump(arr)

	arr.OrderedAppend(7)
	spew.Dump(arr)

	arr.OrderedAppend(6)
	spew.Dump(arr)

	fmt.Println(arr.Has(11))
	fmt.Println(arr.Has(99))
	fmt.Println(arr.Has(0))
	fmt.Println(arr.Has(1))
	fmt.Println(arr.Has(7))
}
