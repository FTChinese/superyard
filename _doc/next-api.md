# Access Next API

## App Registration

    POST /next/apps

### Input
```json
{
    "name": "Olive Graves",
    "slug": "oliver-graves",
    "repoUrl": "http://mofvacfel.iq/hebbupzaf",
    "description": "Hovovvod mi sogat haski faat za ikomur ti jo utnomov.",
    "homeUrl": "http://po.ie/selo"
}
```

* `name: string` required, max 256 chars
* `slug: string` required, max 256 chars
* `repoUrl: string` required, 256 chars
* `description: string` optional, 512 chars
* `homeUrl: string` optional, 256 chars. The url where this app is run.

Owner name will be extracted from request header `X-User-Name`.

### Response

* `400 Bad Request` if request body cannot be parsed as JSON.
	{
		"message": "Problems parsing JSON"
	}

* `422 Unprocessable Entity` 

if required fields are missing
```json
{
	"message": "Validation failed",
	"error": {
	    "field": "name | slug | repoUrl",
    	"code": "missing"
	}
}
```

or the length of  any of the fields exceeds max chars
```json
{
	"message": "The length of xxx should not exceed 255 chars",
	"error": {
	    "field": "email | slug | repoUrl | description | homeUrl",
    	"code": "invalid"
	}
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

## List All Apps

	GET /next/apps?page=<number>

`page` defaults to 1 if it is missing, or is not a number.

### Response

* `400 Bad Request` if query string cannot be parsed.

* `200 OK` with body
```json
[
	{
        "id": 23,
        "name": "Olive Graves",
        "slug": "oliver-graves",
        "clientId": "a2332f5250fb5d591ef0",
        "clientSecret": "90444ef69d74ced39a9a2d59102e8976d2295bf9c051f35074e99695306df718",
        "repoUrl": "http://mofvacfel.iq/hebbupzaf",
        "description": "Hovovvod mi sogat haski faat za ikomur ti jo utnomov.",
        "homeUrl": "http://po.ie/selo",
        "isActive": true,
        "createdAt": "2019-02-10T12:30:10Z",
        "updatedAt": "2019-02-10T12:30:10Z",
        "ownedBy": "foo.bar"
    },
    {
        "id": 22,
        "name": "Electric Kit",
        "slug": "electric-kit",
        "clientId": "50251719118d85ed816f",
        "clientSecret": "a82bc5b42cfad21d315a1136e40271bc99efe6c1089b24fd80db7a7f141cd753",
        "repoUrl": "https://githu.com/FTChinese/nisi",
        "description": "odit hic doloribus alias impedit nam deleniti doloremque cum.",
        "homeUrl": "http://www.ftchinese.com/odit",
        "isActive": true,
        "createdAt": "2019-02-10T10:01:28Z",
        "updatedAt": "2019-02-10T10:01:28Z",
        "ownedBy": "velit"
    }
]
```

## Show an App Information

	GET /next/apps/{name}

### Response

* `400 Bad Request` if request URL does not contain `name` part

```json
{
	"message": "Invalid request URI"
}
```

* `404 Not Found` if the app does not exist

* `200 OK`
```json
{
    "id": 13,
    "name": "Side Tag Mount",
    "slug": "side-tag-mount",
    "clientId": "bfc465fde898eb8450fa",
    "clientSecret": "698ebe8a61da1fe49994e8f7a412a955a0a5bc7390dafa90179cbd78045fa007",
    "repoUrl": "https://githu.com/FTChinese/alias",
    "description": "nihil quis iste dolorem ipsa minima.",
    "homeUrl": "http://www.ftchinese.com/enim",
    "isActive": true,
    "createdAt": "2019-02-10T09:04:19Z",
    "updatedAt": "2019-02-10T09:04:19Z",
    "ownedBy": "dicta_eligendi_vero"
}
```

## Update an App

	PATCH /next/apps/{name}

### Input

```json
{
    "name": "Isaiah Bradley",
    "slug": "isaiah-bradley",
    "repoUrl": "http://jicun.ir/ede",
    "description": "Va lawjum wuamotod ji ju iwibtac az cidaje osapu to.",
    "homeUrl": "http://of.sx/ifeko"
}
```

* `name: string` optional, max 256 chars
* `slug: string` optional, max 256 chars
* `repoUrl: string` optional, 256 chars
* `description: string` optional, 512 chars
* `homeUrl: string` optional, 256 chars. The url where this app is run.

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

* `422 Unprocessable Entity` is the same as App Registration

* `204 No Content` for success.

## Delete an App

	DELETE /next/apps/{name}

### Response

* `400 Bad Request` if request URL does not contain `name` part
```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.

## Transfer an App

	POST /next/apps/{name}/transfer

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

## Create an Access Token for an App

	POST /next/apps/{name}/tokens

### Input

None.

### Response

* `400 Bad Request` if the app name cannot be extracted from url.

* `204 No Content` for success.

## List Tokens of an App

    GET /next/apps/{name}/tokens
    
### Response

```json
[
    {
        "id": 36,
        "token": "61280219e84b32968e08e8f4ec897012bd5909fe",
        "createdAt": "2019-02-10T14:25:04Z",
        "updatedAt": "2019-02-10T14:25:04Z",
        "lastUsedAt": null
    }
]
```

## Delete a Token of an App

    DELETE /next/apps/{name}/tokens/{id}
 
 ### Response
 
 * `204 No Content`

## Create a Personal Access Token

    POST /next/keys
 
### Input
```json
{
    "description": "This is a personal access token",
    "myftEmail": "LisaGarcia@Jabberstorm.gov"
}
```

* `description: string` Optional. Max 256 chars.
* `myftEmail: string` Optional. The email used to login on ftchinese.com. If set, this token will be able to access personal data of this email. Otherwise the token could only access public data.

### Response

* `204 No Content`

## List a User's Personal Access Tokens

	GET /next/keys

* `200 OK` with body

### Response

```json
[
    {
        "id": 37,
        "token": "8ac699fb19e2aeb274763afe9b7df075145624a6",
        "createdAt": "2019-02-10T14:41:27Z",
        "updatedAt": "2019-02-10T14:41:27Z",
        "lastUsedAt": null,
        "description": "This is a personal access token",
        "myftEmail": "LisaGarcia@Jabberstorm.gov",
        "createdBy": "foo.bar"
    }
]
```

## Delete a Personal Access Token

	DELETE /next/keys/{tokenId}

### Response

* `400 Bad Request` if request URL does not contain `name` part

```json
{
	"message": "Invalid request URI"
}
```

* `204 No Content` for success.
