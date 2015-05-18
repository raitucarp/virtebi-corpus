# How to Install
If you setup go environment correctly, then you can do this
```
go install https://github.com/raitucarp/go.virtebi.corpus 
```
Please see commandline usage

note: I do not release binary build yet. Maybe soon

# Usage
```
Usage of ./virtebi-corpus:
  -build=false: build corpus database
  -format="text": print result with formating, json, xml or text
  -listen=0: create http server that listen to a port, example -listen=8989
  -raw=false: print result in raw per line, it does not use any format
  -withprob=false: print probes value
  
example: ./virtebi-corpus -format=json "facebookiscool" "whatdoyouwant" "thisisatext"

```

# Example
Here are some example doing it, with all of command line flag

## JSON formatting
```
$ ./virtebi-corpus -format=json "facebookiscool" "whatdoyouwant" "thisisatext"
{"items":[{"origin":"facebookiscool","result":"facebook is cool"},{"origin":"whatdoyouwant","result":"what do you want"},{"origin":"thisisatext","result":"this is a text"}],"Length":3}
```

## XML formatting
```
$ ./virtebi-corpus -format=xml "facebookiscool" "whatdoyouwant"
<results><items><origin>facebookiscool</origin><result>facebook is cool</result></items><items><origin>whatdoyouwant</origin><result>what do you want</result></items><length>2</length></results>
```

## Text formatting
this example is withprob value 
```
$ ./virtebi-corpus -withprob "facebookiscool" "whatdoyouwant" 
Length = 2
Original		Result		Prob
facebookiscool		facebook is cool		0.00000000042221252134040073
whatdoyouwant		what do you want		0.00000000006427053322919637
```

## JSON with probvalue
```
$ ./virtebi-corpus -format=json -withprob=true "facebookiscool" "whatdoyouwant"
{"items":[{"origin":"facebookiscool","result":"facebook is cool","prob":4.2221252134040073e-10},{"origin":"whatdoyouwant","result":"what do you want","prob":6.427053322919637e-11}],"length":2}
```

## RAW
Finally, raw
```
$ ./virtebi-corpus -format=json -withprob=true -raw=true "facebookiscool" "whatdoyouwant"  "thistextislong"
facebook is cool
what do you want
this text is long
```

You can't do formating(json, xml, text) or pass withprob with raw formatting, because raw is just kind of raw text. It contains result per line

# TODO
- ~~Write testing~~
- Collect more complete corpus
- Fix some bugs
- Refactor some algorithm
- create http server
- ~~Doing with whitespace, etc~~

# Bugs
Currently, If I want to parse this string "thisisactuallyalongtext" it would output: "is actually along text".
I will investigate this case more in the future

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
