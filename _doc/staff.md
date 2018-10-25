## Login

    POST /staff/auth


Input
```json
{
    "userName": "foo.bar",
    "password": "abcedfg",
    "userIp": "127.0.0.1"
}
```

* `400 Bad Request` if body content cannot be parsed as JSON
```json
{
    "message": "Problems parsing JSON"
}
```

* `404 Not Found` if `userName` does not exist or `password` is wrong.

* `200 OK` with body:
```json
{
    "id": 1,
    "email": "foo.bar@ftchinese.com",
    "userName": "foo.bar",
    "displayName": "Foo Bar",
    "department": "tech",
    "groupMembers": 3
}
```

## Forgot Password

    POST /staff/password-reset/letter

### Input
```json
{
    "email": "foo.bar@ftchinese.com"
}
```

### Response

* `400 Bad Request` if request body cannot be parsed as JSON.
```json
{
    "message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity` if `email` is missing or invalid.
```json
{
    "message": "Validation failed"
    "field": "email",
    "code": "missing_field" | "invalid"
}
```

* `404 Not Found` if the `email` is not found.

* `500 Internal Server Error` if token cannot be generated, or token cannot be saved, or email cannot be sent.
```json
{
    "message": "xxxxxxx"
}
```

* `204 No Content` if password reset letter is sent.

## Verify Password Reset Link

    GET /staff/password-reset/tokens/{token}


* `400 Bad Request` if request URL does not contain `token` part
```json
{
    "message": "Invalid request URI"
}
```

* `404 Not Found` if the token does not exist

* `200 OK` with body
```json
{
    "email": "foo.bar@ftchinese.com"
}
```

## Reset Password

    POST /staff/password-reset


### Input
```json
{
    "token": "reset token client extracted from url",
    "password": "8 to 128 chars"
}
```

### Response

`400 Bad Request` if request body cannot be parsed as JSON.
```json
{
    "message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity` if validation failed.
```json
{
    "message": "Validation failed | The length of password should not exceed 128 chars",
    "field": "password",
    "code": "missing_field | invalid"
}
```

* `404 Not Found` if the token is expired or not found.

* `204 No Content` if password is reset succesfully.

## Show Your Personal Data

    GET /user/profile


* `404 Not Found` if this user does not exist.

* `200 OK` with body:
```json
{
    "id": "",
    "userName": "",
    "email": "",
    "isActive": "",
    "displayName": "",
    "department": "",
    "groupMembers": "",
    "createdAt": "",
    "deactivatedAt": "",
    "updatedAt": "",
    "lastLoginAt": "",
    "lastLoginIp": ""
}
```

## Update Display Name

    PATCH /user/display-name


### Input
```json
{
    "displayName": "max 20 chars"
}
```

### Response

* `400 Bad Request` if request body cannot be parsed as JSON.
```json
{
    "message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity` 

if validation failed:
```json
{
    "message": "Validation failed | The length of displayName should not exceed 20 chars",
    "field": "displayName",
    "code": "missing_field | invalid"
}
```

if this `displayName` already exists
```json
{
    "message": "Validation failed",
    "field": "displayName",
    "code": "already_exists"
}
```

* `204 No Content` for success

## Update Email

    PATCH /user/email


### Input
```json
{
    "email": "max 20 chars"
}
```

### Response

* `400 Bad Request` if request body cannot be parsed as JSON.
```json
{
    "message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity` for validation failure:
```json
{
    "message": "Validation failed | The length of email should not exceed 20 chars",
    "field": "email",
    "code": "missing_field | invalid"
}
```

if the email to use already exists
```json
{
		"message": "Validation failed",
		"field": "email",
		"code": "already_exists"
}
```

* `204 No Content` for success.

## Update Password

    PATCH /user/password


### Input
```json
{
    "old": "max 128 chars",
    "new": "max 128 chars"
}
```

The max length limit is random.
Password actually should not have length limit.
But hashing extremely long strings takes time.

### Response

* `400 Bad Request` if request body cannot be parsed as JSON.
```json
{
    "message": "Problems parsing JSON"
}
```

* `422 Unprocessable Entity` if either `old` or `new` is missing in request body, or password is too long.
```json
{
    "message": "Validation failed | Password should not execeed 128 chars",
    "field": "password",
    "code": "missing_field | invalid"
}
```

* `403 Forbidden` if old password is wrong
```json
{
    "message": "wrong password"
}
```

* `204 No Content` for success.

## List My FTC Accounts

    GET /user/myft

* `200 OK`
```json
[
    {
        "myftId": "",
        "myftEmail": "",
        "isVip": "boolean"
    }
]
```

## Link My FTC Account

    POST /user/myft

### Input
```json
{
    "email": "string",
    "password": "string"
}
```

* `400 Bad Request` if request body cannot be parsed as JSON.
```json
{
    "message": "Problems parsing JSON"
}
```

* `404 Not Found` if `email` + `password` verification failed.

* `422 Unprocessable Entity` if the ftc account to add already exist.
```json
{
    "message": "Validation failed",
    "field": "email",
    "code": "already_exists"
}
```

* `204 No Content`

## Delete a My FTC Account

    DELETE /user/myft/{id}


* `400 Bad Request` if request URL does not contain `id` part
```json
{
    "message": "Invalid request URI"
}
```

* `204 No Content` for success