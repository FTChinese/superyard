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

* `400 Bad Request` if request body cannot be parsed as JSON.
	{
		"message": "Problems parsing JSON"
	}

- 422 Unprocessable Entity if required fields are missing,
	{
		"message": "Validation failed",
		"field": "name | slug | repoUrl",
		"code": "missing"
	}
or the length of  any of the fields exceeds max chars,
	{
		"message": "The length of xxx should not exceed 255 chars",
		"field": "email | slug | repoUrl | description | homeUrl",
		"code": "invalid"
	}
or the slugified name of the app is taken
	{
		"message": "Validation failed",
		"field": "slug",
		"code": "already_exists"
	}

- `204 No Content` for success.