package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func MaxInt(a []int) int {
	var max = 0
	for _, i := range a {
		if i > max {
			max = i
		}
	}
	return max
}

func GetRandomComplex(rmin, rmax, imin, imax float64) (float64, float64) {
	// math/rand uses Fibonacci method to generate randoms, in our case it is
	// good enough
	r := rand.Float64()*(rmax-rmin) + rmin
	i := rand.Float64()*(imax-imin) + imin
	return r, i
}

func SplitColors(colors string) []int {
	cArr := strings.Split(colors, ",")
	if len(cArr) != 3 && len(cArr) != 1 {
		log.Fatalf("Color palette must have either 1 or 3 components, got: %d", len(cArr))
	}
	if len(cArr) == 1 {
		c, err := strconv.Atoi(cArr[0])
		if err != nil {
			log.Fatalf("Invalid color: %s", err)
		}
		return []int{c, c, c}
	}

	r, err := strconv.Atoi(cArr[0])
	if err != nil {
		log.Fatalf("Invalid red color component: %s", err)
	}
	g, err := strconv.Atoi(cArr[1])
	if err != nil {
		log.Fatalf("Invalid green color component: %s", err)
	}
	b, err := strconv.Atoi(cArr[2])
	if err != nil {
		log.Fatalf("Invalid blue color component: %s", err)
	}
	return []int{r, g, b}
}
