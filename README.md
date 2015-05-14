# Virtebi Algorithm meet corpus parser
Implementation Virtebi Algorithm by split joined words, into readable words in golang

It begins, when I was searching for solution for parsing domain names, into readable string. And found this
http://stackoverflow.com/questions/195010/how-can-i-split-multiple-joined-words

# How to Install
If you setup go environment correctly, then you can do this
```
go install https://github.com/raitucarp/go.virtebi.corpus 
```
Please see commandline usage

note: I do not release binary build yet. Maybe soon


# How to use it as Package
first get it
```
go get https://github.com/raitucarp/go.virtebi.corpus 
```
then include it
```
package main

import (
	"fmt"
	"github.com/raitucarp/go.virtebi.corpus/corpus"
)

func main() {
	c := corpus.NewCorpus()
	// you need to load it manually
	// it will takes longer time for its first build
	// corpus.db
	c.Load()

	origin, result, prob := c.Match("facebook is true")
	fmt.Println(origin, result, prob)
}

```

# Usage
```
Usage of ./virtebi-corpus:
  -format="text": print result with formating, json, xml or text
  -raw=false: print result in raw per line, it does not use any format
  -withprob=false: print probes value
example: ./virtebi-corpus -format=json "facebookiscool" "whatdoyouwant"

```

# Example
Here are some example doing it, with all of command line flag

## JSON formatting
```
$ ./virtebi-corpus -format=json "facebookiscool" "whatdoyouwant"
{"items":[{"origin":"facebookiscool","result":"facebook is cool"},{"origin":"whatdoyouwant","result":"whatdoyou want"}],"Length":3}

```

## XML formatting
```
$ ./virtebi-corpus -format=xml "facebookiscool" "whatdoyouwant" 
<Output></items><items><origin>facebookiscool</origin><result>facebook is cool</result></items><items><origin>whatdoyouwant</origin><result>whatdoyou want</result></items><Length>3</Length></Output>

```

## Text formatting
this example is withprob value 
```
$ ./virtebi-corpus -withprob "facebookiscool" "whatdoyouwant"  
Length = 3
Original		Result		Prob
facebookiscool		facebook is cool		0.000000000030183641197357495
whatdoyouwant		whatdoyou want		0.00000000007057345569603739
```

## JSON with probvalue
```
$ ./virtebi-corpus -format=json -withprob=true "facebookiscool" "whatdoyouwant"  
{"items":[{"origin":"facebookiscool","result":"facebook is cool","prob":3.0183641197357495e-11},{"origin":"whatdoyouwant","result":"whatdoyou want","prob":7.057345569603739e-11}],"Length":3}

```

## RAW
Finally, raw
```
$ ./virtebi-corpus -format=json -withprob=true -raw=true "facebookiscool" "whatdoyouwant"
facebook is cool
whatdoyou want
```

You can't do formating(json, xml, text) or pass withprob true, because raw is just kind of raw text. It contains result per line

# TODO
- Write testing
- Collect more complete corpus
- Fix some bugs
- Refactor some algorithm
- create http server
- ~~Doing with whitespace, etc~~

# Bugs
~~Currently, If I want to parse this string "thisisstring" it would output: "thisis string".
"thisisastring" to thisis a string, etc
I want you to help me to find out what is going on.~~

# License
The MIT License (MIT)

Copyright (c) 2015 Ribhararnus Pracutiar

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
