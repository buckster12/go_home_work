package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

var AsteriskValues = false

var IgnoreCase = AsteriskValues
var IgnoreDash = AsteriskValues
var AllowSpecialChars = !AsteriskValues
var AllowDigitsInWords = !AsteriskValues

func getSliceOfWords(inputText string) []string {
	// clear special chars
	if !AllowSpecialChars {
		for _, clearChar := range []string{`.`, `,`, `!`, `?`, `:`, `;`, `"`, `'`, "`"} {
			inputText = strings.ReplaceAll(inputText, clearChar, ` `)
		}
	}

	// Split the text and get a slice with words
	re := regexp.MustCompile(`\s+`)
	return re.Split(inputText, -1)
}

func Top10(inputText string) []string {
	const top = 10

	words := getSliceOfWords(inputText)

	// get words map and count of entries
	resultMap := getWordCountMap(words)

	var uniqueWords = make([]string, 0, len(resultMap))
	for value := range resultMap {
		uniqueWords = append(uniqueWords, value)
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		return resultMap[uniqueWords[i]] > resultMap[uniqueWords[j]]
	})

	if len(uniqueWords) > top {
		return uniqueWords[:10]
	}
	return uniqueWords
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
		resultMap[word1]++
	}

	return resultMap
}
