package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

	// if it's not found
	// then builddictionary from text
	// maybe this is take longer time
	// because it's first build database
	if err != nil {
		// build from corpus.txt
		c.buildDictionaryFromText(db)
	} else {
		// build from database cache
		c.buildDictionaryFromCache(db)
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
func (c *Corpus) buildDictionaryFromText(db *leveldb.DB) {
	// load content from corpus.txt
	contents, err := ioutil.ReadFile("corpus.txt")

	// if error then panic
	if err != nil {
		panic(err)
	}

	// convert byte to string
	text := string(contents)

	// build wordlist from text
	wordList := toWords(text)

	// sort wordlist
	sort.Strings(wordList)

	// apply dictionary from word grouping
	c.Dictionary = groupCorpus(wordList)

	// insert to database
	for word, weight := range c.Dictionary {
		// put every word
		// and it's value to the leveldb database
		err = db.Put([]byte(word), []byte(strconv.FormatFloat(weight, 'f', 6, 64)), nil)

		// if error then continue,
		// keep saving another word
		if err != nil {
			continue
		}
	}

}

// Build dictionary from database.
// As long as, the database is exist
func (c *Corpus) buildDictionaryFromCache(db *leveldb.DB) {
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
	err := iter.Error()
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
	corpus.Load()
	return corpus
}

// Search entire text, assign value to the word.
// Based on their distribution cmiiw
func groupCorpus(words []string) map[string]float64 {
	a := make(map[string]float64)

	for i := 0; i < len(words); i++ {
		word := words[i]
		a[word] = a[word] + 1
	}
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

// Split entire corpus text
// to array of words.
func toWords(text string) []string {
	text = strings.ToLower(text)
	return wordPattern.FindAllString(text, -1)
}

// Items structure
type Item struct {
	// Origin is the original text before parsing
	Origin string `json:"origin" xml:"origin"`
	// Result is the result
	Result string `json:"result" xml:"result"`
	// Prob value
	Prob float64 `json:"prob,omitempty" xml:"prob,omitempty"`
}

// Output structure
type Output struct {
	// Items array, is collection of items
	Items []Item `json:"items" xml:"items"`
	// how many text being parsed?
	Length int `json:="length" xml:"items"`
}

var results []Item

// main job
func main() {
	// initialize new corpus
	corpus := NewCorpus()

	// flag formatting
	formatting := flag.String("format", "text", "print result with formating, json, xml or text")
	// withprob formating
	withProbe := flag.Bool("withprob", false, "print probes value")
	// raw formatting
	raw := flag.Bool("raw", false, "print result in raw per line, it does not use any format")

	// parse flag
	flag.Parse()

	// tail is words here
	tail := flag.Args()

	// if tail is empty then print error
	if len(tail) < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("example:", os.Args[0], `-format=json "facebookiscool" "whatdoyouwant"`)
	}

	// iterate tail, and apply corpus.match.
	// Create Item for each result
	for i := 0; i < len(tail); i++ {
		// match in corpus
		origin, result, prob := corpus.Match(tail[i])
		// create item
		r := Item{
			Origin: origin,
			Result: result,
		}

		// if withprob flag is present,
		// apply value to it
		if *withProbe {
			r.Prob = prob
		}

		// insert to results
		results = append(results, r)
	}

	// create output
	o := Output{}
	o.Items = results
	o.Length = len(results)

	// if raw,
	// just prints result perline.
	// No matters other flag is
	if *raw {
		// print result per line
		for i := 0; i < len(o.Items); i++ {
			item := o.Items[i]
			fmt.Println(item.Result)
		}
	} else {
		// get formatting
		switch *formatting {
		// if json then do marshaling
		case "json":
			output, _ := json.Marshal(o)
			fmt.Println(string(output))
		// xml format same as json
		case "xml":
			output, _ := xml.Marshal(o)
			fmt.Println(string(output))
		// text is a bit complicated,
		// because it's separated by tab and line
		case "text":
			header := "Length = " + strconv.Itoa(o.Length) + "\n"
			header += "Original\t\tResult"
			if *withProbe {
				header += "\t\tProb"
			}

			// print header
			fmt.Println(header)
			// print the text
			for i := 0; i < len(o.Items); i++ {
				item := o.Items[i]
				output := item.Origin + "\t\t" + item.Result
				if *withProbe {
					output += "\t\t" + strconv.FormatFloat(item.Prob, 'f', -1, 64)
				}
				fmt.Println(output)
			}
		}
	}

}
