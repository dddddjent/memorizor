# Word Service

# API

All the apis below require an `Authorization` field in the header

## Get /page

Get the number of total pages

#### Request

- None

#### Response:

```json
{
	"pages": 5
}
```

## Get /list/:page/?method=time

Get the list of all the words under this account.

#### Request

- page is the index of a specific page.
  - A single page contains 10 words.
  - An invalid page index should return an error.
  - However, a page index too large will return an empty list
- method: time/alphabetic
  - Sort by time in ascending (old to new) order
  - Sort by alphabetic order of the words in ascending order

#### Response

```json
{
    "list": a list of words
}
```

## Get /today

#### Request

- None

#### Response

- Return a list of words that should be reviewed in the past days in Fibonacci number

```json
{
    "today": [
        [{word}],
        []
    ]
}
```

## Post /word

#### Request

- add
  - Only `word`, `explanation`, `url` are allowed in the request
  - `word` shouldn't be empty, or contains things other than English characters
  - It will be converted to a word with uppercase only at the front
  - Beware that adding an existing word will update it

```json
{
	"method": "add",
	"parameters": {
		"word": "nice",
		"explanation": "something",
		"url": "www.google.com"
	}
}
```

- update
  - Only `word`, `explanation`, `url` are allowed in the request
  - `word` shouldn't be empty, or contains things other than English characters
  - It will be converted to a word with uppercase only at the front
  - Beware that updating a non-existing word will create one

```json
{
	"method": "add",
	"parameters": {
		"word": "nice",
		"explanation": "something",
		"url": "www.google.com"
	}
}
```

- delete
  - `id` is the word's id. It needs to be an existing word

```json
{
	"method": "delete",
	"parameters": {
		"id": "126adf2d-e583-4497-9d5c-1c4c9822ced1"
	}
}
```

- click
  - update the last clicked time of a word to the current calling time
  - `id` is the word's id. It needs to be an existing word

```json
{
	"method": "click",
	"parameters": {
		"id": "98bf1a38-2ce1-41a5-b643-ab41a35ce2b3"
	}
}
```

#### Response

- All the methods will return if it processes successfully

```json
{
	"message": "ok"
}
```
