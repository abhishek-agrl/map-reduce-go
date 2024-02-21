package workers

import (
	"log"
	"slices"
	"strings"
)

func startMap(inputStr string) [number_of_reducers][]WordMap {
	dataArr := strings.Fields(strings.Trim(strings.ReplaceAll(inputStr, ".", " "), " "))
	var wordMapArr []WordMap

	for _, dataWord := range dataArr {
		idx := slices.IndexFunc(wordMapArr, func(wordObj WordMap) bool { return wordObj.Word == dataWord })
		if idx != -1 {
			wordMapArr[idx].Count += 1
		} else {
			wordMapArr = append(wordMapArr, WordMap{dataWord, 1})
		}
	}

	log.Printf("%v Unique words found", len(wordMapArr))

	buckets := [number_of_reducers][]WordMap{}

	for _, wordMap := range wordMapArr {
		bucketIdx := int(rune(wordMap.Word[0]) % number_of_reducers)
		buckets[bucketIdx] = append(buckets[bucketIdx], wordMap)
	}

	return buckets
}
