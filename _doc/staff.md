# CMS User

## Login

    POST /login

Header: `X-User-Ip: <forwarded user ip>`

### Input

```json
{
    "userName": "foo.bar",
    "password": "abcedfg"
}
```

### Response

* `400 Bad Request` if body content cannot be parsed as JSON
```json
{
    "message": "Problems parsing JSON"
}
```

* `404 Not Found` is `userName` does not exist.

* `403 Forbidden` if `password` is wrong.

* `204 No Content` if user name and password are correct.

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
    "message": "Validation failed",
    "field": "email",
    "code": "missing_field | invalid"
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
	"token": "cecb195a8ee363f6c8f2881a76f2346f84dcce981f113b8cfcfb2c58f00468b5",
	"password": "12345678"
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

* `204 No Content` for success.

## Get User Account
    
    GET /staff/account

### Response

* `404 Not Found`

* `200 OK`
```json
{
    "id": 70,
    "email": "neefrankie@163.com",
    "userName": "foo.bar",
    "displayName": "Foo Bar Updated | null",
    "department": "tech | null",
    "groupMembers": 3
}
```

## Get User Profile

    GET /staff/profile

### Response

* `404 Not Found` if this user does not exist.

* `200 OK` with body:
```json
{
    "id": 70,
    "email": "neefrankie@163.com",
    "userName": "foo.bar",
    "displayName": "Foo Bar Updated | null",
    "department": "tech | null",
    "groupMembers": 3,
    "isActive": true,
    "createdAt": "2019-02-10T06:47:20Z",
    "deactivatedAt": "2019-02-10T10:01:26Z",
    "updatedAt": "2019-02-10T11:52:40Z",
    "lastLoginAt": "2019-02-10T11:52:40Z",
    "lastLoginIp": "127.0.0.1"
}
```

## Update Display Name

    PATCH /staff/display-name

### Input
```json
{
	"displayName": "Victor Nee"
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
    "error": {
        "field": "displayName",
        "code": "missing_field | invalid"
    }
}
```

if this `displayName` already exists
```json
{
    "message": "Validation failed",
    "error": {
        "field": "displayName",
        "code": "already_exists"
    }
}
```

* `204 No Content` for success

## Update Email

    PATCH /staff/email

### Input

```json
{
	"email": "neefrankie@163.com"
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
    "error": {
        "field": "email",
        "code": "missing_field | invalid"
    }
}
```

if the email to use already exists
```json
{
    "message": "Validation failed",
    "error": {
        "field": "email",
        "code": "already_exists"
    }	
}
```

* `204 No Content` for success.

## Update Password

    PATCH /staff/password

### Input
```json
{
	"newPassword": "12345678",
	"oldPassword": "12345678"
}
```

Max 128 chars.

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
    "error": {
        "field": "password",
        "code": "missing_field | invalid"
    }
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

    GET /staff/myft

* `200 OK`

```json
[
    {
        "id": "3622eaf4-eeba-4114-a086-c4a7a6ddce8a",
        "unionId": null,
        "email": "dignissimos_tempora@Youfeed.com",
        "userName": "6Franklin",
        "isVip": false
    },
    {
        "id": "ff020312-a0ec-4440-a737-eb1e947ade10",
        "unionId": null,
        "email": "oBarnes@Centidel.name",
        "userName": "9Hawkins",
        "isVip": false
    }
]
```

## Add My FTC Account

    POST /staff/myft

### Input

```json
{
  "email": "molestias_officiis@Wordify.info",
  "password": "12345678"
}
```

### Response

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
    "error": {
        "field": "email",
        "code": "already_exists"
    }
}
```

* `204 No Content`

## Delete a My FTC Account

    DELETE /staff/myft

### Input

```json
{
    "email": "molestias_officiis@Wordify.info"
}
```

### Response

* `400 Bad Request` if request URL does not contain `id` part

```json
{
    "message": "Invalid request URI"
}
```

* `204 No Content` for success