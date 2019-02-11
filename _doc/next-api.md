# Access FTC API

## App Registration

    POST /ftc-api/apps

### Input
```json
{
    "name": "User Login", // required, max 255 chars
    "slug": "user-login", // required, max 255 chars
    "repoUrl": "https://github.com/user-login", // required, 120 chars
    "description": "UI for user login", // optional, 511 chars
    "homeUrl": "https://www.ftchinese.com/user" // optional, 255 chars
}
```

Owner name will extracted from request header `X-User-Name`.

* `400 Bad Request` if request body cannot be parsed as JSON.
	{
		"message": "Problems parsing JSON"
	}

* `422 Unprocessable Entity` 

if required fields are missing
```json
{
	"message": "Validation failed",
	"field": "name | slug | repoUrl",
	"code": "missing"
}
```

or the length of  any of the fields exceeds max chars
```json
{
	"message": "The length of xxx should not exceed 255 chars",
	"field": "email | slug | repoUrl | description | homeUrl",
	"code": "invalid"
}
```

or the slugified name of the app is taken
```json
{
	"message": "Validation failed",
	"field": "slug",
	"code": "already_exists"
}
```

* `204 No Content` for success.

## List Apps

	GET /ftc-api/apps?page=<number>

`page` defaults to 1 if it is missing, or is not a number.

* `400 Bad Request` if query string cannot be parsed.

* `200 OK` with body
```json
[
	{
		"name": "User Login",
		"slug": "user-login",
		"clientId": "20 hexdecimal numbers",
		"clientSecret": "64 hexdecimal numbers",
		"repoUrl": "https://github.com/user-login",
		"description": "UI for user login",
		"homeUrl": "https://www.ftchinese.com/user",
		"isActive": true,
		"createdAt": "",
		"updatedAt": "",
		"ownedBy": "foo.bar"
	}
]
```

## Show an App Information

	GET /ftc-api/apps/{name}

* `400 Bad Request` if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

* `404 Not Found` if the app does not exist

* `200 OK`. See response for List Apps.

## Update an App

	PATCH /ftc-api/apps/{name}

### Input
```json
{
	"name": "User Login", // max 60 chars, required
	"slug": "user-login", // max 60 chars, required
	"repoUrl": "https://github.com/user-login", // 120 chars, required
	"description": "UI for user login", // 500 chars, optional
	"homeUrl": "https://www.ftchinese.com/user" // 120 chars, optional
}
```

* `400 Bad Request`

if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

or if request body cannot be parsed as JSON
```json
{
	"message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity` is the same as App Registration

* `204 No Content` for success.

## Delete an App

	DELETE /ftc-api/apps/{name}


* `400 Bad Request` if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.

## Transfer an App

	POST /ftc-api/apps/{name}/transfer

### Input
```json
{
	"newOwner": "foo.baz"
}
```

### Response

* `400 Bad Request`

if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

or if request body cannot be parsed as JSON
```json
{
	"message": "Problems parsing JSON"
}
```

* `404 Not Found` if the new owner is not found.

* `204 No Content` for success.

## Create an Access Token

	POST /ftc-api/tokens

### Input
```json
{
	"description": "string", // optional, max 255 chars
	"myftId": "string", // optional
	"ownedByApp": "string" // optional
}
```

The creator of this token will always be recorded.

`myftId` and `ownedByApp` should be mutually exclusive.
If `ownedByApp` is present, it means this access token is created for an app. In such case `myftId` must be empty.
If both `myftId` and `ownedByApp` are empty, it must be a personal access token.

### Response

* `400 Bad Request` if request body cannot be parsed
```json
{
	"message": "Problems parsing JSON"
}
```

* `204 No Content` for success.

## List a User's Personal Access Tokens

	GET /ftc-api/tokens/personal


* 200 OK with body
```json
[
	{
		"id": 1,
		"token": "40 hexdecimal numbers",
		"description": "",
		"myftId": "",
		"createdAt": "",
		"updatedAt": "",
		"lastUsedAt": ""
	}
]
```

## Delete a Personal Access Token

	DELETE /ftc-api/tokens/personal/{tokenId}


* `400 Bad Request` if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.

## Show an App's Access Tokens

	GET /ftc-api/tokens/app/{name}

* `400 Bad Request` if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

* `200 OK` with body
```json
[
	{
		"id": 1,
		"token": "40 hexdecimal numbers",
		"description": "",
		"myftId": "",
		"createdAt": "",
		"updatedAt": "",
		"lastUsedAt": ""
	}
]
```

## Delete an App's Access Token

	DELETE /ftc-api/tokens/app/{name}/{tokenId}

* `400 Bad Request` if request URL does not contain `name` and `tokenId` part, or tokenId < 1.
```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.