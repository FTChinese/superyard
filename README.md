## Sitemap

Since restful is resource-oriented, it might be a better idea to group those resources by owner, not by what you, from the standpoint of machine, have.

### Article

### Search

* GET `/search/staff?k={login|email|display}&v={}` Return id
* GET `/search/user?k={email|name}&v={}` Return uuid

### Login

FTC staff self-service section. Consumed by `backyard-user`.

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
* POST `/user/password` Change password

* GET `/user/myft` List all myft accounts
* POST `/user/myft` Add a myft account
* DELETE `/user/myft/:id` Delete a myft account

### Admin

Request header must contain `X-User-Name` field and this user's privileges will be checked to see if he/she has the power to perform those actions.

* GET `/admin/staff/exists?k={name|email}&v={:value}`
* POST `/admin/staff/new`
* PUT `/admin/staff/new` Activate a deactivated user.
* GET `/admin/staff/roster` All staff

* GET `/admin/staff/profile/{name}` Show a staff's profile
* PATCH `/admin/staff/profile/{name}` Update staff's profile
* DELETE `/admin/staff/profile/{name}` Delete a staff

* GET `/admin/vip-roster` Show all myft accounts that are granted VIP.
* PUT `/admin/vip/{myftId}` Grant vip to a myft account
* DELETE `/admin/vip/{myftId}` Delete vip status of a myft account

### FTC Apps
* POST `/apps/ftc` Create a new ftc app
* GET `/apps/ftc` Show all ftc apps. Anyone can see details of an app created by any others.
* GET `/apps/ftc/:name` Show a ftc app
* POST `/apps/ftc/:name` Only owner can edit it. So posted data should include owner id.
* DELETE `/apps/ftc/:name`
* POST `/apps/ftc/:name/transfer`

### Personl Access Tokens
* GET `/tokens/next-api` Show urls
* POST `/tokens/next-api/personal` Create a new personal access token.
* GET `/tokens/next-api/personal/:userName` Show all access tokens a user owns
* DELETE `/tokens/next-api/personal/:userName` Revoke all access tokens
* PATCH `/tokens/next-api/personal/:userName/:tokenId` Update the description of an access token.
* DELETE `/token/next-api/personal/:userName/:tokenId` Delete an access token

* POST `/tokens/next-api/app` Create a new token for an app
* GET `/tokens/next-api/app/:slugName` Show all access tokens owned by an app.
* DELETE `/tokens/next-api/app/:slugName` Revoke all tokens owned by an app
* PATCH `/tokens/next-api/app/:slugName/:tokenId` Update the description of an app token
* DELETE `/tokens/next-api/app/:slugName/:tokenId` Revoke an access token owned by an app.

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

* GET `/ftc-user/:userId` Show a user's profile, vip status, membership, orders placed, articles favoured, reading history.

### Order

* GET `/order/:orderId` Show an order

### Stats

* GET `/stats/new-users?{start=YYYY-MM-DD&end=YYYY-MM-DD}`
* GET `/stats/new-members?{start=YYYY-MM-DD&end=YYYY-MM-DD}`
* GET `/stats/new-orders?{start=YYYY-MM-DD&end=YYYY-MM-DD}`