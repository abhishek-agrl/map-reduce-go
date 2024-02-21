package workers

import (
	"bytes"
	"encoding/gob"
	"log"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func getWordMapList(bucketGobs [][]byte) [][]string {
	log.Printf("Decompressing Gobs for further processing")
	var wordMapList [][]string
	for _, bucketGob := range bucketGobs {
		buf := bytes.NewBuffer(bucketGob)
		var bucketStr string
		gob.NewDecoder(buf).Decode(&bucketStr)
		wordMapList = append(wordMapList, strings.Split(strings.Trim(bucketStr, "\n"), "\n"))
	}
	return wordMapList
}

func startReducer(bucketGobs [][]byte) string {
	var reducedWordMap []WordMap
	wordMapList := getWordMapList(bucketGobs)
	log.Printf("Decompression Complete, Starting to Merge Count")

	for _, dataArr := range wordMapList {
		for _, line := range dataArr {
			lineArr := strings.Fields(line)
			word := lineArr[0]
			count, _ := strconv.Atoi(lineArr[1])
			idx := slices.IndexFunc(reducedWordMap, func(globalWord WordMap) bool { return globalWord.Word == word })
			if idx > -1 {
				reducedWordMap[idx].Count += uint(count)
			} else {
				reducedWordMap = append(reducedWordMap, WordMap{word, uint(count)})
			}
		}
	}

	log.Printf("Merge Count Complete, Final Count Ready")

	sort.Slice(reducedWordMap, func(i, j int) bool {
		return reducedWordMap[i].Word < reducedWordMap[j].Word
	})

	return convertMaptoText(reducedWordMap)
}
