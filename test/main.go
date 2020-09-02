package main

import "fmt"

var stocks []string

func main() {
	stocks = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
	size := len(stocks) / 5
	var j int
	for i := 0; i < len(stocks); i += size {
		j += size
		if j > len(stocks) {
			j = len(stocks)
		}
		// do what do you want to with the sub-slice, here just printing the sub-slices
		fmt.Println(stocks[i:j])
	}
}
