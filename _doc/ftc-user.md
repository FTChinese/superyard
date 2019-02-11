# FTC Users

## Show a User's Account

    GET /user/account

### Response

* `400 Bad Request` if request URL does not contain `userId` part
```json
{
	"message": "Invalid request URI"
}
```

* `404 Not Found` if the the user is not found.

* `200 OK`

```json
{
    "id": "ff03c3b9-f8c9-4228-936b-cd1099a2113e",
    "unionId": null,
    "email": "LisaGarcia@Jabberstorm.gov",
    "userName": "dolor",
    "isVip": false,
    "mobile": null,
    "createdAt": "2019-02-08T12:40:33Z",
    "nickname": null,
    "membership": {
        "tier": null,
        "cycle": null,
        "expireDate": null
    }
}
```

## Show a User's Orders

### Response
```json
[
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
    },
    {
        "orderId": "FTD143C8F7BDAB774D",
        "userId": "ff03c3b9-f8c9-4228-936b-cd1099a2113e",
        "loginMethod": 1,
        "tier": "standard",
        "cycle": "year",
        "listPrice": 258,
        "netPrice": 258,
        "payMethod": "tenpay",
        "createdAt": "2019-02-09T10:18:30Z",
        "confirmedAt": "2019-02-09T10:18:30Z",
        "startDate": "2019-02-09",
        "endDate": "2020-02-10",
        "clientType": 3,
        "clientVersion": "1.1.1",
        "userIp": "57.170.41.5",
        "userAgent": "Opera/9.80 (X11; Linux i686; U; en) Presto/2.2.15 Version/10.10"
    }
]
```

