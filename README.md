# Overview

## Development

完全同[API](https://github.com/FTChinese/subscription-api/blob/master/_doc/development.md)项目。

差异之处如下：

1. 本程序是web app，因此包含了HTML模版文件，放在`web`文件夹下，使用的模版引擎是[https://github.com/flosch/pongo2](https://github.com/flosch/pongo2)，语法和Django、Nunjucks类似。这些模版文件在build时嵌入到二进制文件中，见`web/render.go`文件。
2. 模版文件并非手写，而是用[superyard-react](https://github.com/FTChinese/superyard-react)中的脚本生成。
3. 跟目录下的`client_version_next`和`client_version_ng`也由前端项目的脚本生成。


## See documentation

`godoc -http=:6060`

## Articles

## Login

* POST `/staff/auth`
* POST `/staff/password-reset/letter`
* GET `/staff/password-reset/tokens/{token}`
* POST `/staff/password-reset`

## Personal settings

Request header must contain `X-User-Name` field.

* GET `/user/profile` Show a logged in staff's information.
* PATCH `/user/display-name` Change display name
* PATCH `/user/email` Change email
* PATCH `/user/password` Change password

* GET `/user/myft` List all myft accounts
* POST `/user/myft` Add a myft account
* DELETE `/user/myft/:id` Delete a myft account

## Admin

Request header must contain `X-User-Name` field and this user's privileges will be checked to see if he/she has the power to perform those actions.

* GET `/admin/staff/exists?k={name|email}&v={:value}` Checks if a staff exists
* POST `/admin/staff/new`
* GET `/admin/staff/roster?page=<number>` All staff

* GET `/admin/staff/profile/{name}` Show a staff's profile
* PUT `/admin/staff/profile/{name}` Restore a deleted staff
* PATCH `/admin/staff/profile/{name}` Update staff's profile
* DELETE `/admin/staff/profile/{name}?rmvip=true|false` Delete a staff

* GET `/admin/vip` Show all myft accounts that are granted VIP.
* PUT `/admin/vip/{myftId}` Grant vip to a myft account
* DELETE `/admin/vip/{myftId}` Delete vip status of a myft account

## FTC API 

### Apps using FTC API

* POST `/ftc-api/apps` Create a new ftc app
* GET `/ftc-api/apps?page=<number>` Show all ftc apps. Anyone can see details of an app created by any others.
* GET `/ftc-api/apps/{name}` Show a ftc app
* PATCH `/ftc-api/apps/{name}` Only owner can edit it. So posted data should include owner id.
* DELETE `/ftc-api/apps/{name}`
* POST `/ftc-api/apps/{name}/transfer`

### Access Tokens
* POST `/ftc-api/tokens` Create an access token. It could belong to a person or an app, depending on the data passed in.

* GET `/ftc-api/tokens/personal` Show all access tokens a user owns

* DELETE `/ftc-api/token/personal/{tokenId}` Delete an access token

* GET `/ftc-api/tokens/app/{name}` Show all access tokens owned by an app.

* DELETE `/ftc-api/tokens/app/{name}/{tokenId}` Revoke an access token owned by an app.

## CMS API

### CMS API Apps
For CMS apps, there's no owership. Anybody can edit.

* POST `/apps/cms` Create a new cms app.
* GET `/apps/cms` List all cms apps.
* GET `/apps/cms/:name` Show an app info.
* POST `/apps/cms/:name` Update an app info.
* DELETE `/apps/cms/:name`
* GET `/apps/cms/:name/perms` Get the unix permission of an app.

### CMS API Tokens
* GET `/tokens/cms-api` Show all access tokens
* POST `/tokens/cms-api` Create a new token to access cms-api
* PATCH `/tokens/cms-api/:tokenId` Update description of an access token
* DELETE `/tokens/cms-api/:tokenId` Delete an access token.

## Subscription

### Promotion Schedules

* `GET /subscription/promotion?page=<int>` List all promotion schedules

* `POST /subscription/promotion/schedule` Create a new promotion schedule.

* `GET /subscription/promotion/schedule/:id`
* `DELETE /subscription/promotion/schedule/:id`

* `POST /subscription/promotion/pricing` Create procing for this promotion schedule.

* `POST /subscription/promotion/banner` Create promotion content for this schedule.
* `PATCH /subscription/promotion/banner/:id`

* GET `/subscription/plans` Show all plans
* POST `/subscripiton/plans/new` Create a new group of pricing plans.  
* DELETE `/subscription/plans/delete/:id` Delete a set of plans.

## FTC User

* GET `/ftc-user/profile/{userId}` Show a user's profile, vip status, membership
* GET `/ftc-user/profile/{userId}/orders` Show a user's orders
* GET `/ftc-user/profile/{userId}/login?page=<number>` Show a user's login history

## Search

* GET `/search/user?k=<name|email>&v=:value`
* GET `/search/orders?{start=YYYY-MM-DD&end=YYYY-MM-DD}` Show all orders within the specified time range

## Data

* GET `/stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD`
* GET `/stats/subscription/daily?{start=YYYY-MM-DD&end=YYYY-MM-DD}`
* GET `/stats/orders/daily?{start=YYYY-MM-DD&end=YYYY-MM-DD}`
