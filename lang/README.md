# Extended Parser

Extended parser from [_Writing An Interpreter In Go_](https://interpreterbook.com/), **_highly_** recommend the book.

This is a lexer and parser lib for use in [Blitzbase](https://github.com/orangeseeds/blitzbase). It is used to parse the language for assigning access permission for collections. Extra features in this particular parser include, property accessors parser for syntax like, `collection.name` or `$request.auth.id`. 

The parser parses expression in a right-to-left manner, no operator precedence in this parser.

A statement like this 
`$request.data.is_valid_date & collection.exists | collection.id != 10;` 
will be parsed as such;
`($request.data.count != (collection.exists | (collection.id = 10)))`

This language is supposed to have only one statement in the entire program, we chain a bunch of expressions to evaluate if the request to access the collection or any kind of data is valid or not.

This is inspired by pocketbase's access permission language. To check permission use `$request` and the request is supposed to have, these attributes for now,
- **auth**: includes record data of the auth-collection which the auth_token points to.
- **data**: contains the data coming with the request, in case of POST/PUT/PATH requests.
- **method**: is the request method, GET/POST/PUT/PATCH/DELETE...

So, an example check could be,
`$request.auth.name ~= $request.data.name`
This is supposed to check, the tokens, auth-record name to the name field in the incoming data and run a `LIKE` check on them, and if it passes then request will be able to access the collection.
