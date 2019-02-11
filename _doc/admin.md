# Administration

All request must contain header `X-User-Name`.

## Staff Exists

    GET /admin/staff/exists?k={name|email}&v={value}

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

## Search an Account

    GET /admin/account/search?k={name|email}&v={value}

### Response

* `400 Bad Request` either `k` or `v` are empty, or `k` is not one of `email` or `name`

* `404 Not Found` if the search found no result.

* `200 OK`:

```json
{
    "id": 4,
    "email": "mGordon@Realcube.net",
    "userName": "illo_rem",
    "displayName": "Mr. Dr. Arthur Woods",
    "department": "tech",
    "groupMembers": 3
}
```

## Create an Account for a Staff

    POST /admin/accounts

### Input

```json
{
    "email": "foo.bar@ftchinese.com", // required, unique, max 256 chars
    "userName": "foo.bar", // required, unique, max 256 chars
    "displayName": "Foo Bar", // optional, unique, max 256 chars
    "department": "tech", // optinal, max 256 chars
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
    "message": "The length of xxx should not exceed 256 chars",
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

    GET /admin/accounts?page=<number>

`page` defaults to 1 if omitted or is not a number. Returns 20 entires per page.

### Response

200 OK with an array:
```json
[
    {
        "id": 1,
        "email": "officiis@Camido.mil",
        "userName": "RalphAdams",
        "displayName": "Christopher Wright",
        "department": "tech",
        "groupMembers": 3
    },
    {
        "id": 2,
        "email": "SeanRose@Rhyloo.biz",
        "userName": "in_optio",
        "displayName": "Albert Gilbert",
        "department": "tech",
        "groupMembers": 3
    }
]
```

## Show a Staff's Profile

    GET /admin/accounts/{name}

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
    "id": 48,
    "email": "facere_pariatur_sed@Rhybox.org",
    "userName": "aliquid",
    "displayName": "Theresa Diaz",
    "department": "tech",
    "groupMembers": 3,
    "isActive": true,
    "createdAt": "2019-02-09T23:29:53Z",
    "deactivatedAt": null,
    "updatedAt": "2019-02-09T23:29:53Z",
    "lastLoginAt": null,
    "lastLoginIp": null
}
```
## Update a Staff's Profile

    PATCH /admin/accounts/profile/{name}

 Input and response are identical to creating a new staff `POST /admin/account`.

### Input
```json
{
    "userName": "foo.bar",
    "email": "foo.bar@ftchinese.com",
    "displayName": "Foo Bar",
    "department": "tech",
    "groupMembers": 3
}
```

All fields are optional.

### Response

* `400 Bad Request` if the `{name}` part in URL cannot be extracted, or input body cannot be parsed into JSON.

* `422 Unprocessable Entity` if validation failed, or `userName` is duplicate.

* `204 No Content` for success.
 
 ## Delete Staff

    DELETE /admin/staff/profile/{name}

This endpoint performs those operations:

1. Turn the `is_active` column in the `backyard.staff` table to `false` (not actually delete any data), and record the time of this operation;

2. If `revokeVip` is `true`, find all the FTC accounts associated with this staff, and turn the `is_vip` column of `cmstmp01.userinfo` to `false`;

3. Deletes all the FTC accounts associated with this staff (deletes data in the `backyard.staf_myft` table);

4. Deletes all personal access tokens in the `oauth.access` table.

### Input
```json
{
  "revokeVip": true
}
```

`revokeVip` indicates whether the vip granted to the ftc accounts associated with this staff should also be removed. It defaults to `true`.

### Response

* `400 Bad Request` if request URL does not contain `name`.
```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.

## Restore a Removed Staff

    PUT /admin/accounts/{name}

### Input

None.

### Response

* `400 Bad Request` if url does not contain the `name` part.
```json
{
    "message": "Invalid request URI"
}
```

* `204 No Content`

## List All FTC VIPs

    GET /admin/vip?page=<number>

Pagination defaults to 20 items.

### Response

* `200 OK`:

```json
[
    {
        "id": "6f39f2c9-40ce-4a0a-82b4-8a9bdcb9d906",
        "unionId": null,
        "email": "DonaldGutierrez@Eabox.mil",
        "userName": "dJackson",
        "isVip": true
    },
    {
        "id": "a7ae6327-d294-42b1-9485-9e93ac39ac54",
        "unionId": null,
        "email": "xPeterson@Meejo.net",
        "userName": "in_maxime",
        "isVip": true
    }
]
```

## Grant VIP to an FTC Account

    PUT /admin/vip/{email}

### Response

* `400 Bad Request` if `email` is not present in URL.
```json
{
    "message": "Invalid request URI"
}
```

* `204 No Content` if granted.

## Revoke VIP of an FTC Account

    DELETE /admin/vip/{email}

### Input

None

### Response

* `400 Bad Request` if `email` is not present in URL.

```json
{
    "message": "Invalid request URI"
}
```

* `204 No Content` for success.