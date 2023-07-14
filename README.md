# Overview

## Development

完全同[API](https://github.com/FTChinese/subscription-api/blob/master/_doc/development.md)项目。

差异之处如下：

1. 本程序是前后端分离的single page web app，后端需要输出HTML，模版文件放在`web`文件夹下，使用的模版引擎是[https://github.com/flosch/pongo2](https://github.com/flosch/pongo2)，语法和Django、Nunjucks类似。这些模版文件在build时嵌入到二进制文件中，见`web/render.go`文件。
2. 后端输出的网页只是一个空白页，所有UI都交给前端用react渲染，这个空白网页对应的模版是`web/template/next.html`。
3. 另外一个模版文件`web/template/ng.html`供此前用Angular写的前端app使用，但Angular版已不再开发，仅保留此前已有的功能。
4. 模版文件并非手写，而是用[superyard-react](https://github.com/FTChinese/superyard-react)中的脚本生成。
5. 根目录下的`client_version_next`和`client_version_ng`也由前端项目的脚本生成。


## View documentation

`godoc -http=:6060`

## API Endpoints

The following is an overview all api endpoints. For more details, follow these links:

### Login

Login by FTC staff only.

* POST `/auth/login` Login
* POST `/auth/password-reset/letter` Send user a letter to reset password.
* GET `/auth/password-reset/tokens/{token}` Verify password reset token.
* POST `/auth/password-reset` Reset password.

### Personal settings

Request header must contain `X-User-Name` field.

* GET `/settings/account`
* PATCH `/settings/account/email` Change email
* PATCH `/settings/account/display-name` Change display name
* PATCH `/settings/account/password` Change password
* GET `/settings/account/profile` Show a logged in staff's information.

### OAuth Access Token

* GET `/oauth/apps?page<int>&per_page=<int>` Get a list of apps.
* POST `/oauth/apps` Create a new app for which you can generate access tokens.
* GET `/oauth/apps/:id` Get a specific app
* PATCH `/oauth/apps/:id` Update an app
* DELETE `/oauth/apps/:id` Deactivate an app.

* GET `/oauth/keys?client_id=<string>&page=<number>&per_page=<number>` Get a list of access tokens.
* POST `/oauth/keys` Craete a new access token.
* DELETE `/oauth/keys/:id` Delete an access token.

### FTC Users

* GET `/readers/search?q=<email|username|phone>&kind=<ftc|wechat>`
* GET `/readers/ftc/:id` Get an ftc account
* GET `/readers/ftc/:id/profile` Get the profile of an ftc account
* GET `/readers/wx/:id` Load a wechat account
* GET `/readers/wx/:id/profile` Load more details a wechat account.

### Sandbox Users

FTC user in sandbox mode. Used for testing only.

* POST `/sandbox` Create a test user
* GET `/sandbox` List aoo test users
* GET `/sandbox/:id` Load a test account
* DELETE `/sandbox/:id` Delete a test account
* PATCH `/sandbox/:id/password` Change the password of a test account

### Memberships

* POST `/memberships` Update or create a membership
* DELETE `/memberships` Deelete a membership.

### Snapshots

Show memberships change history

* GET `/snapshots?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>` Show a user's membership change history.

### Apple In-App Purchase

* GET `/iap?page=<int>&per_page=<int>` List a user's IAP record
* GET `/iap/:id` Load a single IAP subscription
* PATCH `/iap/:id` Refresh an existing IAP
* POST `/iap/:id/link` Link iap to an ftc account
* POST `/iap/:id/unlink` Sever links between an iap and ftc account

### Paywall

Please refer to the [Paywall](https://github.com/FTChinese/subscription-api/blob/master/_doc/paywall.md) section on how the paywall system is designed. This project simply transfer requests to [API](https://github.com/FTChinese/subscription-api).

* GET `/paywall` Load paywall data
* POST `/paywall/banner` Create a banner
* POST `/paywall/banner/promo` Create promo banner
* DELETE `/paywall/banner/promo` Delete promo banner
* POST `/paywall/products` Create a product
* GET `/paywall/products` List products
* GET `/paywall/products/:productId` Load product by id
* POST `/paywall/products/:productId/activate` Put a product on paywall
* PATCH `/paywall/products/:productId` Update a product
* POST `/paywall/prices` Create a new price under a specific product.
* GET `/paywall/prices?product_id=<string>&live=<true|false>` List all prices under a product.
* GET `/paywall/prices/:priceId` Load a price
* POST `/paywall/prices/:priceId/activate` Turn a price into active state under a product.
* POST `/paywall/prices/:priceId/deactivate` Turn a price into deactivatie state under a product.
* PATCH `/paywall/prices/:priceId` Update a price
* PATCH `/paywall/prices/:priceId/discounts` Refersh discounts attached to a price.
* DELETE `/paywall/prices/:priceId` Archive a price.
* POST `/paywall/discounts` Create discount
* DELETE `/paywall/discounts/:id` Delete a discount

### Stripe

Link Stripe price and coupon to ftc price.

* GET `/stripe/prices?page=<int>&per_page=<int>&live=<bool>` Get a list of stripe prices
* GET `/stripe/prices/:id?live=<bool>&refresh=<bool>` Load a stripe price
* PATCH `/stripe/prices/:id?live=<bool>` Update a stripe price's metadata
* PATCH `/stripe/prices/:id/activate?live=<bool>` Activate a stripe price.
* PATCH `/stripe/prices/:id/deactivate?live=<bool>` Deactivate a stripe price.
* GET `/stripe/prices/:id/coupons?live<bool>` List coupons attached to a price.
* GET `/stripe/coupons/:id?live=<bool>&refresh=<bool>` Load a stripe coupon
* POST `/stripe/coupons/:id?live=<bool>` Update a coupon
* PATCH `/stripe/coupons/:id/activate?live<bool>` Activate a coupon
* DELETE `/stripe/coupons/:id?live=<bool>`

### Android Release

This is used to publish latest version data to API. You also need to upload Android APK to minio storage for user to download.

* POST `/android/releases` Create a new release.
* GET `/android/releases` List android releases.
* GET `/android/release/:versionName` Load the details of a specific version.
* PATCH `/android/release/:versionName` Update a release.
* DELETE `/android/release/:versionName` Delete a release.

### Wiki

Used to store development documentations.

* GET `/wiki` List wiki articles
* POST `/wiki` Create a new wiki article
* GET `/wiki/:id` Load one wiki article
* PATCH `/wiki/:id` Updatea a wiki article

### Legal Documents

Legal documents as listed here: https://next.ftacademy.cn/terms.

* GET `/legal?page=<int>&per_page=<int>` List legal documents
* POST `/legal` Create a new legal document.
* GET `/legal/:id` Load a legal document.
* PATCH `/legal/:id` Updatea a legal document.
* POST `/legal/:id/publish` Make a legal document publicly visible.
