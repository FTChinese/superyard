## Sitemap

Since restful is resource-oriented, it might be a better idea to group those resources by owner, not by what you, from the standpoint of machine, have.

### Article

### Login

FTC staff self-service section. Consumed by `backyard-user`.

* GET `/staff/exists?k={name|email}&v={:value}`
* POST `/staff/auth`
* POST `/staff/password-reset/letter`
* GET `/staff/password-reset/tokens/{token}`
* POST `/staff/password-reset`
Workflow:
1. Users sends token and new password;
2. Use the token to find out the email associated with it;
3. Use the email to identify which row to update.

### Personal settings

Request header must contain `X-User-Name` field.

* GET `/user/profile` Show a logged in staff's information.
* PATCH `/user/display-name` Change display name
* PATCH `/user/email` Change email
* PATCH `/user/password` Change password

* GET `/user/myft` List all myft accounts
* POST `/user/myft` Add a myft account
* DELETE `/user/myft/:id` Delete a myft account

### Admin

Request header must contain `X-User-Name` field and this user's privileges will be checked to see if he/she has the power to perform those actions.

* POST `/admin/staff/new`
* GET `/admin/staff/roster?page=<number>` All staff

* GET `/admin/staff/profile/{name}` Show a staff's profile
* PUT `/admin/staff/profile/{name}` Restore a deleted staff
* PATCH `/admin/staff/profile/{name}` Update staff's profile
* DELETE `/admin/staff/profile/{name}?rmvip=true|false` Delete a staff

* GET `/admin/vip` Show all myft accounts that are granted VIP.
* PUT `/admin/vip/{myftId}` Grant vip to a myft account
* DELETE `/admin/vip/{myftId}` Delete vip status of a myft account

### FTC Apps
* POST `/ftc-api/apps` Create a new ftc app
* GET `/ftc-api/apps` Show all ftc apps. Anyone can see details of an app created by any others.
* GET `/ftc-api/apps/{name}` Show a ftc app
* PATCH `/ftc-api/apps/{name}` Only owner can edit it. So posted data should include owner id.
* DELETE `/ftc-api/apps/{name}`
* POST `/ftc-api/apps/{name}/transfer`

### Personl Access Tokens
* POST `/ftc-api/tokens` Create an access token. It could belong to a person or an app, depending on the data passed in.
<!-- * POST `/ftc-api/tokens/personal` Create a new personal access token. -->
* GET `/ftc-api/tokens/personal` Show all access tokens a user owns
<!-- * DELETE `/ftc-api/tokens/personal/:userName` Revoke all access tokens -->
<!-- * PATCH `/ftc-api/tokens/personal/:userName/:tokenId` Update the description of an access token. -->
* DELETE `/ftc-api/token/personal/{tokenId}` Delete an access token

<!-- * POST `/ftc-api/tokens/app` Create a new token for an app -->
* GET `/ftc-api/tokens/app/{name}` Show all access tokens owned by an app.
<!-- * DELETE `/ftc-api/tokens/app/:slugName` Revoke all tokens owned by an app -->
<!-- * PATCH `/ftc-api/tokens/app/:slugName/:tokenId` Update the description of an app token -->
* DELETE `/ftc-api/tokens/app/{name}/{tokenId}` Revoke an access token owned by an app.

### CMS apps

For CMS apps, there's no owership. Anybody can edit.

* POST `/apps/cms` Create a new cms app.
* GET `/apps/cms` List all cms apps.
* GET `/apps/cms/:name` Show an app info.
* POST `/apps/cms/:name` Update an app info.
* DELETE `/apps/cms/:name`
* GET `/apps/cms/:name/perms` Get the unix permission of an app.

* GET `/tokens/cms-api` Show all access tokens
* POST `/tokens/cms-api` Create a new token to access cms-api
* PATCH `/tokens/cms-api/:tokenId` Update description of an access token
* DELETE `/tokens/cms-api/:tokenId` Delete an access token.

### Search

* GET `/search/user?k=<name|email>&v=:value`
* GET `/search/orders?{start=YYYY-MM-DD&end=YYYY-MM-DD}` Show all orders within the specified time range

### User

* GET `/ftc-user/profile/{userId}` Show a user's profile, vip status, membership
* GET `/ftc-user/profile/{userId}/orders` Show a user's orders
* GET `/ftc-user/order/{orderId}` Show a single order

### Stats

* GET `/stats/new-users?{start=YYYY-MM-DD&end=YYYY-MM-DD}`
* GET `/stats/new-members?{start=YYYY-MM-DD&end=YYYY-MM-DD}`
* GET `/stats/new-orders?{start=YYYY-MM-DD&end=YYYY-MM-DD}`