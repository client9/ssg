# funcs - additional functions for Go templates


## Goals and Alternatives

- not pipeline based.  The template pipeline operations look good, but work for the simplest cases only. After that it gets ugly quick.
- independant and exportable - to serve as a base, or for use in different templating systems
- stdlib only - keep it simple.  Add more advanced functions go in a different module.
- minimal - what a descending sort? Then use `reverse (sort $list)`.  
- perfer naming over additional arguments. `sort`, sortNum` versus some additional argument.
- immutable data structures

## Alternatives

[masterminds/sprig](https://github.com/masterminds/sprig) -- appear semi-abandoned, pipeline based, has a number of unusual functions.

[Hugo](https://gohugo.io/) -- the static site generator has many functions, but has inconsistent design, and implementation is tied to Hugo (hard to rip out).

## Not Included

- **Internationalization Support** -- gets complicated fast.  Good for a separate module.  Note: all string operations are rune-based.
- **Checksum and Hashes (regular or cryptographic)** -- limited uses, many variations. Good for separate module.
- **Cryptography functions**-- limited use, many variations.
- **OS and Environment** -- consider passing these items instead of making them part of the template.
- **Math trig** -- limited utility.

## Recipes

**descending sorts** -- use `reverse (sort $val)`

**sequences**  while we do have `seq 10` (1..10), etc, in Go 1.24 you can do `{{ range 10 }}{{ . }}{{ end }}` will print 0-9.

**chomp** classic function to remove ending "\r\n".  Probably solved with whitespace control in the template, or `trim`, or `trimRight "\r\n"`

**capitalize**  as found in other templating languages is: `firstUpper (lower $str)`

**join with trailing**   `concat (join $list ", ") ","

**join with special case for last**   `concat (join (drop $list -1) ", ") " and " (take $list -1)` --> "a, b, c and d"

**toString**  use `printf "%v" $val` or `{{ $val }}`

**Date Formatting**  see `$date.Format`

**UTC time** see `now.UTC`


