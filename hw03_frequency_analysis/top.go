package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

var AsteriskValues = true

var IgnoreCase = AsteriskValues
var IgnoreDash = AsteriskValues
var AllowSpecialChars = !AsteriskValues
var AllowDigitsInWords = !AsteriskValues

func Top10(inputText string) []string {
	type wordEntry struct {
		word  string
		count int
	}

	const top = 10

	// clear special chars
	if !AllowSpecialChars {
		for _, clearChar := range []string{`.`, `,`, `!`, `?`, `:`, `;`, `"`, `'`, "`"} {
			inputText = strings.ReplaceAll(inputText, clearChar, ` `)
		}
	}

	// Split the text and get a slice with words
	re := regexp.MustCompile(`\s+`)
	words := re.Split(inputText, -1)

	// get words map and count of entries
	resultMap := getWordCountMap(words)

	// insert map values to my struct to get data sorted by key,
	// filter values and sort it
	wordEntries := make([]wordEntry, 0)
	for word, count := range resultMap {
		wordEntries = append(wordEntries, wordEntry{word, count})
	}
	sort.Slice(wordEntries, func(i, j int) bool {
		return wordEntries[i].count > wordEntries[j].count
	})

	//fmt.Println(wordEntries)

	outputSlice := make([]string, 0)
	for i, entry := range wordEntries {
		if i > top-1 {
			break
		}
		outputSlice = append(outputSlice, entry.word)
	}

	return outputSlice
}

// create a map, where key is a word
// and value is a number of entries the word in given text
func getWordCountMap(words []string) map[string]int {
	resultMap := make(map[string]int)

	digitChecker := regexp.MustCompile(`\d`)

	for _, word1 := range words {
		if len(word1) == 0 {
			continue
		}
		if IgnoreCase {
			word1 = strings.ToLower(word1)
		}
		if IgnoreDash && word1 == `-` {
			continue
		}
		// skip words with digits
		if !AllowDigitsInWords && digitChecker.MatchString(word1) {
			continue
		}

		// increase counter or create if not exist
		if _, ok := resultMap[word1]; !ok {
			resultMap[word1] = 1
		} else {
			resultMap[word1]++
		}
	}

	return resultMap
}
