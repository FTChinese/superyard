# Administration

All request must contain header `X-User-Name`.

## Staff exists

    GET /admin/staff/exists?k={name|email}&v={:value}

### Response

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

* `204 No Content` if the user exists.

## Create an account for a staff

    POST /admin/staff/new

### Input
```json
{
    "email": "foo.bar@ftchinese.com", // required, unique, max 80 chars
    "userName": "foo.bar", // required, unique, max 255 chars
    "displayName": "Foo Bar", // optional, unique, max 255 chars
    "department": "tech", // optinal, max 255 chars
    "groupMembers": 3  // required, > 0
}
```

### Response

* `400 Bad Request` if request body cannot be parsed:
```json
{
    "message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity`:

if any of the required fields is missing
```json
{
    "message": "Validation failed",
    "field": "email | userName | groupMembers",
    "code": "missing"
}
```

if email is not a valid email address
```json
{
    "message": "Validation failed",
    "field": "email",
    "code": "invalid"
}
```

if the length of any string fields is over 255:
```json
{
    "message": "The length of xxx should not exceed 255 chars",
    "field": "email | userName | displayName | department",
    "code": "invalid"
}
```
if any of unique fields is already taken by others:
```json
{
    "message": "Validation failed",
    "field": "email | userName | displayName",
    "code": "already_exists"
}
```

* `204 No Content` if a new staff is created.

## List all staff

    GET /admin/staff/roster?page=<number>

`page` defaults to 1 if omitted or is not a number. Returns 20 entires per page.

### Response

200 OK with an array:
```json
[
    {
		"id": 1,
		"email": "foo.bar@ftchinese.com",
		"userName": "foo.bar",
		"displayName": "Foo Bar",
		"department": "tech",
		"groupMembers": 3
    }
]
```

## Show a Staff's Profile

    GET /admin/staff/profile/{name}

* `400 Bad Request` if url does not contain the `name` part.
```json
{
    "message": "Invalid request URI"
}
```

* `404 Not Found` if the requested user is not found

* `200 OK`
```json
{
    "id": "",
    "userName": "",
    "email": "",
    "isActive": true,
    "displayName": "",
    "department": "",
    "groupMembers": 3,
    "createdAt": "",
    "deactivatedAt": "",
    "updatedAt": "",
    "lastLoginAt": "",
    "lastLoginIp": ""
}
```

## Restore a Removed Staff

    PUT /admin/staff/profile/{name}

No input.

Response:

* `400 Bad Request` if url does not contain the `name` part.
	{
		"message": "Invalid request URI"
	}

* `204 No Content`

## Update a Staff's Profile

    PATCH /admin/staff/profile/{name}

 Input and response are identical to creating a new staff `POST /admin/staff/new`.

 ## Delete Staff

    DELETE /admin/staff/profile/{name}?rmvip=<true|false>

`rmvip` defaults to true if omitted, or cannot be converted to a boolean value.

`name` is a staff's login name.

* `400 Bad Request` if request URL does not contain `name`.
```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.

## List All FTC VIPs

    GET /admin/vip

* `200 OK` with body:
```json
[
    {
		"myftId": "string",
		"myftEmail": "string"
	}
]
```

## Grant VIP to an FTC Account

    PUT /admin/vip/{myftId}


* `400 Bad Request` if `myftId` is not present in URL.
```json
{
    "message": "Invalid request URI"
}
```

* 204 No Content if granted.

## Revoke VIP of an FTC Account

    DELETE /admin/vip/{myftId}


* `400 Bad Request` if `myftId` is not present in URL.
```json
{
    "message": "Invalid request URI"
}
```

* `204 No Content` if revoked successuflly.