package main

import "fmt"

func main() {
	s:="MCMXCIV"
	fmt.Println(romanToInt(s))

}

func romanToInt(s string) int {
	sum := 0
	num := map[string]int{
		"I": 1,
		"V": 5,
		"X": 10,
		"L": 50,
		"C": 100,
		"D": 500,
		"M": 1000,
	}
	for i, v := range s {
		if i+1> len(s){
			break
		} else if int(s[i+1]) > num[string(v)] {
			sum += num[string(v)] - int(s[i+1])
		} else {
			sum += num[string(v)]
		}
	}
	return sum

}
