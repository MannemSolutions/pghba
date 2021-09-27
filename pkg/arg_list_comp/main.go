package arg_list_comp

/*
This package is meant to easily define a list as an argument.
There are some examples, like:
* jinja / go template, but there are decided to be too verbose to easily parse as argument
* python list comprehension, but even that is too verbose
* use regular expressions, e.a.:
  * 'item[123]' would be converted in ['item1','item2','item3']
  * 'item[4-6]' would be converted in ['item4','item5','item6']
  * '(item|part)[7-8]' would be converted in ['item7','item8','part7','part8']
* bash glob lists, e.a.:
  * 'item{1,2,3}' would be converted in ['item1','item2','item3']
  * 'item{99..101}' would be converted in ['item99','item100','item101']
  * '{item,part}{7..9}' would be converted in ['item7','item8','part7','part8']

After deep consideration, bash glob has been decided to be the basis.
the main reason is that regular expressions works with characters instead of integers, and as such something like
'user{1..20}' would not be possible with regular expressions.
Do note that we added an extension, so that brackets with character lists would also work (on top of glob options). e.a.
'user[0-9a-f]' would also generate a list of 16 users with last character being a hex value.

For performance reasons, the initial implementation will closely resemble a generator (like in python),
which means there will be a .Next() method which wll return `the next element, and true` unless the loop is done,
in which case it will return `"" and false`.
*/
