# Matching rules

I have a suggestion that folks don’t use and don’t know the full power of Chuck URL matching. So I decided to write this small cheat sheet.

## 1.Equality

Two URLs are matched if they are equal. Chuck ignores the URL schema and port.

So `https://api.github.com/search/repositories` is matched with:

- `https://api.github.com/search/repositories` - identical
- `http://api.github.com/search/repositories` - different schema but still matched
- `https://api.github.com:443/search/repositories` - port is ignored

## 2.Wildcard

`https://api.github.com/users/:username/events` means anything can be instead of **username**.
`https://api.github.com/users/john-doe/events` has matched with URL above.

### Important

You are able to add wildcard URL and concrete URLs simultaneously. In this case, the concrete URL always has a bigger priority.
For example, you added two patterns:
- `https://api.github.com/users/:username/events` 
- `https://api.github.com/users/peter/events`

- url `https://api.github.com/users/peter/events` will be matched with second patterns 
- `https://api.github.com/users/john/events` with first.

## 3. Catch-all wildcard

If you care only about prefix in the URL you are able to use a catch-all wildcard.
Pattern `https://api.github.com/authorizations/events/*` will match with all URLs bellow

- `https://api.github.com/authorizations/events/1/2/2`
- `https://api.github.com/authorizations/events/john`
- `https://api.github.com/authorizations/events/user/update/john`

## 4. Query string matching

If pathes are matched router is checking query string matching according to the following rules.

- Equality match
- Key with wildcard
- Catch-all wildcard

Router ignores the query string keys order (!!!).

### Examples

We added two patterns:
- `https://api.github.com/repos?format=json&token=*&id=:id`
- `https://api.github.com/repos?token=*&format=json`

Instead of path matching in this case `id=*` and `id=:id` means almost the same.

- `https://api.github.com/repos?format=json&token=123456&id=12` is matched with the first pattern.
- `https://api.github.com/repos?format=json&token=8797` with the second.

Url `https://api.github.com/repos?format=xml&token=123` doesn't match with any pattern (format=json != format = xml)

#### Catch-all wildcard example

`https://test.net:443/secure.google.com/v1/authinit?format=json&*` means the `format=json` is mandatory and I don't care about another query params.

URLs 
`https://test.net:443/secure.google.com/v1/authinit?format=json&token=12` 
and 
`https://test.net:443/secure.google.com/v1/authinit?format=json&token=12&code=9876`
are matched with the patterns.

But `https://test.net:443/secure.google.com/v1/authinit?format=xml&token=12` isn't matched.

## Examples and unit tests

The big set of examples you can find in the file with unit tests:
[https://github.com/gavrilaf/grouter/blob/master/grouter_test.go]
