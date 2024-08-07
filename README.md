# Blitzbase

Blitzbase is a realtime database API inspired by pocketbase & supabase with realtime subscriptions using SSEs(which I plan to upgrade to websockets) and a REST API written in go.

#### Running the project
(make sure that sqlite is installed in your device.)
```
$   go run ./cmd/
```

<!--
 To see it in action open index.html from inside the [ui](https://github.com/orangeseeds/Blitzbase/tree/main/ui) folder.
-->

### Features

- Collection Based Data: Tables are represented as collections.
- Subscribable Events: Events with specific topics can be subscribed to, like Edit events on certain collections, or even specific records.
- Permission Parser: It is used to parse the language for assigning access permission for collections. Extra features in this particular parser include, property accessors parser for syntax like, `collection.name` or `$request.auth.id`(similar to pocketbase). 

The parser parses expression in a right-to-left manner, no operator precedence in this parser(which will be added soon).

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
