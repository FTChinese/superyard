# Search

## Find a User

	GET /search/user?k=<name|email>&v=<value>

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

* `200 OK` with body:
```json
{
    "id": "a7ae6327-d294-42b1-9485-9e93ac39ac54",
    "unionId": null,
    "email": "xPeterson@Meejo.net",
    "userName": "in_maxime",
    "isVip": true
}
```

## Find an Order

    GET /search/order?id=<order id>

### Response

* `400 Bad Request` if query parameter cannot be parsed, or the value of `id` is empty.

* `404 Not Found` if the order is not found.

* `200 OK`:
```json
{
    "orderId": "FT378780AC514DA666",
    "userId": "ff03c3b9-f8c9-4228-936b-cd1099a2113e",
    "loginMethod": 1,
    "tier": "standard",
    "cycle": "year",
    "listPrice": 258,
    "netPrice": 258,
    "payMethod": "tenpay",
    "createdAt": "2019-02-09T10:17:24Z",
    "confirmedAt": "2019-02-09T10:17:24Z",
    "startDate": "2019-02-09",
    "endDate": "2020-02-10",
    "clientType": 3,
    "clientVersion": "1.1.1",
    "userIp": "78.163.95.6",
    "userAgent": "Mozilla/5.0 (iPod; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) CriOS/44.0.2403.67 Mobile/12H143 Safari/600.1.4"
}
```