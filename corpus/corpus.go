package corpus

import (
	"github.com/cheggaaa/pb"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"fmt"
)

// create word pattern
// Don't know whether a-z+ or \w is the fastest
var wordPattern = regexp.MustCompile("[a-z]+")

// create corpus structure
type Corpus struct {
	// dictionary is just
	// word and the value
	Dictionary map[string]float64

	// get maximum longest word
	// it's important to use this
	Max float64

	// total number of dictionary
	Total float64
}

// load the corpus into existence,
// let the war begin
func (c *Corpus) Load() {


	// if it's not found
	// then builddictionary from text
	// maybe this is take longer time
	// because it's first build database
	if c.IsAlreadyBuilt() {
		// build from database cache
		c.buildDictionaryFromCache()
	} else {
		// build from corpus.txt
		c.buildDictionaryFromText()
	}

	// set max value
	c.getMax()
	// set total value
	c.getTotal()
}

// Create dictionary from text,
// that is corpus.txt. It's important to keep
// corpus.txt and dont delete it.
// Especially from very early corpus build.
func (c *Corpus) buildDictionaryFromText() {
	// open corpus.db
	// corpus.db is standard db name
	db, err := leveldb.OpenFile("corpus.db", nil)

	// defer closing database
	defer db.Close()

	// if error is not null
	// then panic
	if err != nil {
		panic(err)
	}

	// load content from corpus.txt
	contents, err := ioutil.ReadFile("corpus.txt")

	// if error then panic
	if err != nil {
		panic(err)
	}

	// convert byte to string
	text := string(contents)

	c.saveToDB(text, db)
}

func (c *Corpus) IsAlreadyBuilt() bool {
	// open corpus.db
	// corpus.db is standard db name
	db, err := leveldb.OpenFile("corpus.db", nil)

	// defer closing database
	defer db.Close()

	// if error is not null
	// then panic
	if err != nil {
		panic(err)
	}

	// testing to get word "a"
	_, err = db.Get([]byte("a"), nil)

	if err != nil {
		return false
	} else {
		return true
	}
}

func (c *Corpus) Build() {
	c.buildDictionaryFromText()
}

func (c *Corpus) saveToDB(text string, db *leveldb.DB) {
	dictionary := make(map[string]float64)
	words := strings.Split(text, "\n")
	bar := pb.StartNew(len(words))
	fmt.Println("Please wait for build corpus.db")
	for i := 0; i < len(words); i++ {
		word_score := strings.Fields(words[i])
		if len(word_score) > 1 {
			word := word_score[0]
			score := word_score[1]
			dictionary[word], _ = strconv.ParseFloat(score, 64)
			db.Put([]byte(word), []byte(score), nil)
		}
		bar.Increment()
	}
	fmt.Println("Build corpus.db done.")
	c.Dictionary = dictionary
}

// Build dictionary from database.
// As long as, the database is exist
func (c *Corpus) buildDictionaryFromCache() {
	// open corpus.db
	// corpus.db is standard db name
	db, err := leveldb.OpenFile("corpus.db", nil)

	// defer closing database
	defer db.Close()

	// if error is not null
	// then panic
	if err != nil {
		panic(err)
	}

	// make dictionary
	dict := make(map[string]float64)

	// do iterator
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		word := iter.Key()
		weight := iter.Value()
		w, err := strconv.ParseFloat(string(weight), 64)
		if err != nil {
			panic(err)
		}
		dict[string(word)] = w
	}
	// release iterator
	iter.Release()
	err = iter.Error()
	if err != nil {
		panic(err)
	}

	// apply dict to corpus.Dictionary
	c.Dictionary = dict
}

// build total value
func (c *Corpus) getTotal() {
	// sum is float
	var sum float64

	// dp sum weight value
	for _, weight := range c.Dictionary {
		sum += weight
	}
	// apply sum to total
	c.Total = sum
}

// build maximum value
// from dictionary
func (c *Corpus) getMax() {
	// simple longest
	// that is empty string
	longest := ""

	// do compare the longest word
	// from dictionary
	for word, _ := range c.Dictionary {
		if len(word) > len(longest) {
			longest = word
		}
	}

	// apply the longest to
	// corpus max
	c.Max = float64(len(longest))
}

// Match is the rule here,
// take string and search in dictionary.
// Virtebi algorithm in action
func (c *Corpus) Match(text string) (string, string, float64) {
	fields := strings.Fields(text)
	text = strings.Join(fields, "")

	// create probs
	probs := []float64{1.0}
	lasts := []float64{0}

	// algorithm doing the job
	for i := 1; i < len(text)+1; i++ {
		fi := float64(i) - c.Max
		m := math.Max(0, fi)

		var maxprob, maxk float64
		for j := int(m); j < i; j++ {
			prob := probs[j] * c.wordProb(text[j:i])
			if prob > maxprob {
				maxprob = prob
				jfloat := float64(j)
				if jfloat > maxk {
					maxk = jfloat
				}
			}

		}

		// append maxprob to probs
		probs = append(probs, maxprob)
		lasts = append(lasts, maxk)
	}
	// create array of words
	words := []string{}

	// get text length
	textLength := len(text)

	// finalize the result
	for a := 0; a < textLength; a++ {
		start := int(lasts[textLength])
		words = append(words, text[start:textLength])
		textLength = int(lasts[textLength])
	}

	// reverse the array of string
	words = reverseString(words)

	// return origin, result, and probs
	return text, strings.Join(words, " "), probs[len(probs)-1]
}

// Corpus constructor
func NewCorpus() *Corpus {
	corpus := new(Corpus)
	return corpus
}

// Search entire text, assign value to the word.
// Based on their distribution cmiiw
func groupCorpus(text string) map[string]float64 {
	a := make(map[string]float64)

	words := strings.Split(text, "\n")
	for i := 0; i < len(words); i++ {
		word_score := strings.Fields(words[i])
		if len(word_score) > 1 {
			word := word_score[0]
			score := word_score[1]
			a[word], _ = strconv.ParseFloat(score, 64)
		}
		fmt.Println(word_score)
	}
	/*for i := 0; i < len(words); i++ {
		word := words[i]
		a[word] = a[word] + 1
	}*/
	return a
}

// Search word probability
func (c *Corpus) wordProb(s string) float64 {
	// if there is such word
	// in dictionary
	if val, ok := c.Dictionary[s]; ok {
		// return word weight value
		// divided by total of corpus
		return val / c.Total
	} else {
		return 0
	}
}

// Since, sort.Reverse is not actually works.
// This function is replace it. This one is
// awesome
func reverseString(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
