package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/raitucarp/go.virtebi.corpus/corpus"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"strconv"
)

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

func isDatabaseExist() bool {
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
		return false
	}
	return true
}

// main job
func main() {
	// initialize new corpus
	c := corpus.NewCorpus()

	// flag formatting
	formatting := flag.String("format", "text", "print result with formating, json, xml or text")

	// withprob formating
	withProbe := flag.Bool("withprob", false, "print probes value")

	// build an http server
	listen := flag.Int("listen", 0, "create http server that listen to a port, example -listen=8989")

	// raw formatting
	raw := flag.Bool("raw", false, "print result in raw per line, it does not use any format")

	build := flag.Bool("build", false, "build corpus database")

	// parse flag
	flag.Parse()

	// tail is words here
	tail := flag.Args()

	if *build {
		if !c.IsAlreadyBuilt() {
			c.Build()
		} else {
			fmt.Println("corpus.db is already built")
		}
		os.Exit(0)
	}

	if !c.IsAlreadyBuilt() {
		fmt.Println("You have no db. Please build the database, by using -build=true")
		flag.PrintDefaults()
		os.Exit(0)
	} else {
		c.Load()
	}


	// if tail is empty then print error
	if len(tail) < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("example:", os.Args[0], `-format=json "facebookiscool" "whatdoyouwant"`)
		os.Exit(0)
	}

	// iterate tail, and apply corpus.match.
	// Create Item for each result
	for i := 0; i < len(tail); i++ {
		// match in corpus
		origin, result, prob := c.Match(tail[i])
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

	if *listen > 0 {
		fmt.Println("use http server")
	}

	if *raw {
		// print result per line
		for i := 0; i < len(o.Items); i++ {
			item := o.Items[i]
			fmt.Println(item.Result)
		}
		os.Exit(0)
	}

	// if raw,
	// just prints result perline.
	// No matters other flag is

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
