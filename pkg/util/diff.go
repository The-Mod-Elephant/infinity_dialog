package util

import (
	"cmp"
	"slices"
	"strconv"
)

func SortedDifference(slice1, slice2 *[]string) *[]string {
	diff := []string{}
	m := map[string]int{}

	for _, s := range *slice1 {
		m[s] = 1
	}
	for _, s := range *slice2 {
		m[s] += 1
	}
	for k, v := range m {
		if v == 1 {
			diff = append(diff, k)
		}
	}
	slices.SortFunc(diff, func(a, b string) int {
		v1, _ := strconv.Atoi(a)
		v2, _ := strconv.Atoi(b)
		return cmp.Compare(v1, v2)
	})
	return &diff
}
