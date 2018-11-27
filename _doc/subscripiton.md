# Subscription

## List Promotion Schedules by Pagination

    GET /subscription/promos?page=<int>

Output 5 items per page. If page parameter is omitted, default to 1.

### Response

* `400 Bad Request` if query parameter cannot be parsed.

* `500 Internal Server Error` for any database error.

* `200 OK`:

```json
[
    {
        "id": 1,
        "name": "Compaign name",
        "description": "Notes on this campaign",
        "startAt": "2018-11-10T16:00:00Z",
        "endAt": "2018-11-11T16:00:00Z",
        "plans": {
            "standard_year": {
                "tier": "standard",
                "cycle": "year",
                "price": 198.00,
                "id": 10,
                "descripiton": "FT中文网 - 标准会员",
                "message": ""
            },
            "standard_month": {
                "tier": "standard",
                "cycle": "month",
                "price": 28.00,
                "id": 5,
                "descripiton": "FT中文网 - 标准会员",
                "message": ""
            },
            "premium_year": {
                "tier": "standard",
                "cycle": "year",
                "price": 1998.00,
                "id": 100,
                "descripiton": "FT中文网 - 高端会员",
                "message": ""
            }
        },
        "banner": {
            "heading": "FT中文网会员订阅服务",
            "subHeading": "欢迎您",
            "content": [
                "希望全球视野的FT中文网，能够带您站在高海拔的地方俯瞰世界，引发您的思考，从不同的角度看到不一样的事物，见他人之未见！"
            ]
        },
        "isEnabled": true,
        "createdAt": "2006-01-02T15:04:05Z",
        "updatedAt": "2006-01-02T15:04:05Z",
        "createdBy": "weiguo.ni"
    }
]
```

`plans` and `banner` could be `null` if user failed to complete the form.

## Create a New Promotion Schedule

    POST /subscription/promos

### Input

```json
{
    "name": "Double Eleven Promotion", // Required. Max 256 chars
    "description": "Whatever you like in case you forgot why you created this promotion", // Optional. Max 256 chars
    "startAt": "2006-01-02T15:04:05Z", // Required. ISO 8601 date time string. When will this promotion start.
    "endAt": "2006-01-02T15:04:05Z" // Requried. ISO8601 date time string. When will this promotion end.
}
```

Client should make sure the `startAt` and `endAt` is properly formatted, and must be a time not later than now (what's the purpose if you prepare to scheudle a compaign in the passed?).

### Response

* `400 Bad Request` if request body cannot be parsed as valid JSON.

* `422 Unprocessable Entity`

if `name` is empty:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "name",
        "code": "missing_field"
    }
}
```

if `name` exceeds 256 chars:
```json
{
    "message": "The length of name should not exceed 256 chars",
    "error": {
        "field": "name",
        "code": "invalid"
    }
}
```

if `description` exceeds 256 chars:
```json
{
    "message": "The length of description should not exceed 256 chars",
    "error": {
        "field": "description",
        "code": "invalid"
    }
}
```

if `startAt` is empty:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "startAt",
        "code": "missing_field"
    }
}
```

if `endAt` is empty:
```json
{
    "message": "Validation faild",
    "error": {
        "field": "endAt",
        "code": "missing_field"
    }
}
```

* `200 OK`returns the inserted row's id.
```json
{
    "id": 2
}
```

## Get One Promotion Compaign

    GET /subscription/promos/{id}

### Response

* `400 Bad Request` if URL parameter `id` is not a number or cannot be extracted.

* `404 Not Found` if the promotion with this `id` does not exist.

* `200 OK`. See `List Promotion Schedules by Pagination` for the structure of response. The data is one item of the array.

## Activate a Promotion Campaign

    PUT /subscription/promos/{id}

After an compaign is created, it won't be used until you explicitly activate it.

Not input body.

### Response

* `400 Bad Request` if `id` cannot be parsed to an integer.

* `422 Unprocessable Entity` if

`plans` or `banner` column if `null`:
```json
{
    "message": "Pleans complete the pricing plans | Please complete the promotion banner content",
    "error": {
        "field": "plans | banner",
        "code": "missing_field"
    }
}
```

* `204 No Content` if performed successfully.

## Delete a Promotion Campaign

    DELETE /subscription/promos/{id}

### Response

* `400 Bad Request` if `id` cannot be parsed to an integer.

* `204 No Content` if performed successfully.

## Set/Update Pricing Plans for a Promotion

    PATCH /subscription/promos/{id}/pricing

### Input

```json
{
    "standard_year": {
        "tier": "standard",
        "cycle": "year",
        "price": 198.00,
        "id": 10,
        "descripiton": "FT中文网 - 标准会员",
        "message": ""
    },
    "standard_month": {
        "tier": "standard",
        "cycle": "month",
        "price": 28.00,
        "id": 5,
        "descripiton": "FT中文网 - 标准会员",
        "message": ""
    },
    "premium_year": {
        "tier": "standard",
        "cycle": "year",
        "price": 1998.00,
        "id": 100,
        "descripiton": "FT中文网 - 高端会员",
        "message": ""
    }
}
```

### Response

* `400 Bad Request` if `id` cannot be parsed to an integer.

* `422 Unprocessable Entity`

if `standard_year` does not exist:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "standard_year",
        "code": "missing_field"
    }
}
```

if `premium_year` does not exist:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "premium_year",
        "code": "missing_field"
    }
}
```

if `standard_year`, `standard_month` and `premium_year` exists but their fields are invalid:

if `tier` is empty:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.tier",
        "code": "missing_field"
    }
}
```

if `cycle` is empty:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.cycle",
        "code": "missing_field"
    }
}
```

if `tier` is not one of `standard` or `premium`:
```json
{
    "message": "Tier must be one of standard or premium",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.tier",
        "code": "invalid"
    }
}
```

if `cyle` is not one of `year` or `month`:
```json
{
    "message": "Tier must be one of standard or premium",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.cycle",
        "code": "invalid"
    }
}
```

if `price` is less than or equal to 0:
```json
{
    "message": "Price must be greated than 0",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.price",
        "code": "invalid"
    }
}
```

if `descripiton` is empty:
```json
{
    "message": "Validation failed",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.description",
        "code": "missing_field"
    }
}
```

If `description` is greater than 128 chars:
```json
{
    "message": "The length of description should not exceed 128 chars",
    "error": {
        "field": "<standard_year|standard_month|premium_year>.price",
        "code": "invalid"
    }
}
```

* `204 No Content` if performed successfully.

## Set/Update Barrier Banner for a Promotion

    PATCH /subscrption/promos/{id}/banner

### Input

```json
{
    "heading": "FT中文网会员订阅服务",
    "subHeading": "欢迎您",
    "content": [
        "希望全球视野的FT中文网，能够带您站在高海拔的地方俯瞰世界，引发您的思考，从不同的角度看到不一样的事物，见他人之未见！"
    ]
}
```

### Response

* `400 Bad Request` if `id` cannot be parsed to an integer.

* `204 No Content` if performed successfully.