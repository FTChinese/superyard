# FTC Users

## Find a User

	GET /search/user?k=<name|email>&v=<value>

* `400 Bad Request` if url query string cannot be parsed:
```json
{
	"message": "Bad request"
}
```

or either `k` or `v` cannot be found in query string:
```json
{
	"message": "Both 'k' and 'v' should be present in query string"
}
```

or if the value of url query parameter `k` is neither `name` nor `email`
```json
{
	"message": "The value of 'k' must be one of 'name' or 'email'"
}
```

* `404 Not Found` if the the user with the specified `name` or `email` is not found.

* 200 OK with body:
```json
{
	"id": "",
	"name": "",
	"email": ""
}
```

## Get a User's Profile

	GET /ftc-user/profile/{userId}

* `400 Bad Request` if request URL does not contain `userId` part
```json
{
	"message": "Invalid request URI"
}
```

* `404 Not Found` if the the user is not found.
```json
{
	"id": "",
	"name": "",
	"email": "",
	"gender": "M | F",
	"familyName": "",
	"givenName": "",
	"mobileNumber": "",
	"birthdate": "",
	"address": "",
	"createdAt": "",
	"membership": {
		"tier": "standard | premium",
		"bilingCycle": "year | month",
		"expireDate": ""
	}
}
```