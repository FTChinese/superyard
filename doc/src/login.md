# Login

## Check if a staff exsits

GET `/staff/exists?k={name|email}&v={:value}`

Response:

`204 No Content` if exists.

`400 Bad Request` if any of `k` or `v` is missing from query string, or `k` is neither `name` nore `email`.

`404 Not Found` if not exists.

`k` specifies the criteria to filter data, either by `name` or by `email`.

## Login to CMS

POST `/staff/auth`

Input `{userName: string, password: string, userIp: string}`

Response:

`400 Bad Request` if body content cannot be parsed as JSON
```json
{ "message": "Problems parsing JSON" }
```

`404 Not Found` if `userName` does not exist or `password` is wrong.


## Request a password reset letter

POST `/staff/password-reset/letter`

## Verify password reset link

GET `/staff/password-reset/tokens/{token}`

## Allow user to reset password

POST `/staff/password-reset`