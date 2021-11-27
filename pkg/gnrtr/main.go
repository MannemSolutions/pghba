package gnrtr

import (
	"fmt"
	"regexp"
)

/*
This package is meant to easily define a list within an argument passed to hba.
hba has an object Rules{}, which can use arguments in the gnrtr format and expand into multiple Rules.
Gnrtr formatted arguments can be used on connType, database, user, and address, and as such many rules can be defined
with one call of pghba.

There are some existing implementations that come close to this implementation, like:
* jinja / go template, but there are decided to be too verbose to easily parse as argument
* python list comprehension, but even that is too verbose
* use regular expressions, e.a.:
  * 'item[123]' would be converted in ['item1','item2','item3']
  * 'item[4-6]' would be converted in ['item4','item5','item6']
  * '(item|part)[7-8]' would be converted in ['item7','item8','part7','part8']
* bash glob lists, e.a.:
  * 'item{1,2,3}' would be converted in ['item1','item2','item3']
  * 'item{99..101}' would be converted in ['item99','item100','item101']
  * 'item_{a..z}' would be converted in ['item_a','item_b','item_c']
  * '{item,part}{7..9}' would be converted in ['item7','item8','part7','part8']

They have been used as an inspiration, but this module is built ground up.

After deep consideration, regexp has been decided to be the basis, with the addition of the bash loops defined with
`{1..10}`, or `{a..z}`. Regexp is designed as a filter language, and as such has many options that has not been
implemented, like:
* `{1,3}` will not insert the character 1 to 3 times
* `*` would mean 'any number of characters hereafter' and is not really deterministic by nature
* `.` would mean 'any character' and could be implemented, but would multiply the number of results with 128
  If this would still be intended, then use character loops (like [a-zA-Z0-9] or similar)
* And many more, for similar reasons.

The following is implemented:
* array. Example: (item1|item2|item3) (like in regexp)
* charlist. Example: [a-z_,|] (like in regexp)
* charloop. Example: {a..z} (like in bash)
* intloop. Example: {1..99} (like in bash)

*/

var (
	reIntLoops     = regexp.MustCompile(`{(\d+)..(\d+)}`)
	reIntLoop      = regexp.MustCompile(fmt.Sprintf("^%s$", reIntLoops.String()))
	reCharLoops    = regexp.MustCompile(`{(\S)..(\S)}`)
	reCharLoop     = regexp.MustCompile(fmt.Sprintf("^%s$", reCharLoops.String()))
	reCharLists    = regexp.MustCompile(`\[([^]]+)]`)
	reCharList     = regexp.MustCompile(fmt.Sprintf("^%s$", reCharLists.String()))
	reArrays       = regexp.MustCompile(`\(([^)]+)\)`)
	reArray        = regexp.MustCompile(fmt.Sprintf("^%s$", reArrays.String()))
	rePlaceholders = regexp.MustCompile(`\$\{([0-9]+)\}`)
)
