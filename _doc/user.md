# Users

## Show a User's FTC Account

    GET /users/ftc/account/{id}

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
    "id": "c14e81a6-2967-424d-afae-7f67e8924727",
    "unionId": null,
    "email": "niweiguo@outlook.com",
    "userName": null,
    "isVip": false,
    "mobile": null,
    "createdAt": "2019-03-24T13:20:16Z",
    "updatedAt": "2019-03-24T13:38:55Z",
    "nickname": null,
    "membership": {
        "tier": null,
        "cycle": null,
        "expireDate": null
    }
}
```

## Show a User's Orders

    GET /users/ftc/orders/{id}?page=<number>&per_page=<number>

It lists an FTC account's order. If the account is bound to a wechat account, orders purchased by wechat account is also listed.

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

## Show FTC Account's Login History

    GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>

### Response

```json
[
    {
        "userId": "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae",
        "loginMethod": "email",
        "clientType": "web",
        "clientVersion": "0.5.10",
        "userIp": "::1",
        "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
        "createdAt": "2019-04-21T02:52:16Z"
    },
    {
        "userId": "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae",
        "loginMethod": "email",
        "clientType": "web",
        "clientVersion": "0.5.10",
        "userIp": "::1",
        "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
        "createdAt": "2019-04-21T02:49:40Z"
    }
]
```

## Show Wechat Account

    GET /users/wx/account/{id}

### Response

```json

```

## Show Wechat User's Orders

    GET /users/wx/orders/{id}?page=<number>&per_page=<number>

### Response

```json
{
    "id": "",
    "unionId": "ogfvwjn5kmva3hRz4_SvRujh4mJM",
    "email": "",
    "userName": null,
    "isVip": false,
    "mobile": null,
    "nickname": "倪卫国",
    "membership": {
        "tier": null,
        "cycle": null,
        "expireDate": null
    },
    "createdAt": "2019-03-25T04:06:04Z",
    "updatedAt": "2019-03-25T04:06:04Z"
}
```
## Show Wechat Account's Login History

    GET /usrs/wx/login-history/{id}?page=<number>&per_page=<number>

### Response

```json
[
    {
        "unionId": "ogfvwjn5kmva3hRz4_SvRujh4mJM",
        "openid": "ob7fA0tTAcNnfB7rt9z3eKUe4EAM",
        "appId": "wxacddf1c20516eb69",
        "clientType": "android",
        "clientVersion": "2.0.2",
        "userIp": "192.168.10.29",
        "userAgent": "okhttp/3.12.0",
        "CreatedAt": "2019-03-25T04:06:04Z",
        "UpdatedAt": "2019-03-25T04:06:04Z"
    }
]
```