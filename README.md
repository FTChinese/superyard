## Sitemap

Since restful is resource-oriented, it might be a better idea to group those resources by owner, not by what you, from the standpoint of machine, have.

### Article

### Search

* GET `/search/staff?k={login|email|display}&v={}` Return id
* GET `/search/user?k={email|name}&v={}` Return uuid

### Login

FTC staff self-service section. Consumed by `backyard-user`.

* POST `/staff/auth`
* POST `/staff/password-reset/email`
* GET `/staff/password-reset/tokens/{token}`
* POST `/staff/password-reset`

### Staff admin

Request header must contain `X-User-Name` field and this user's privileges will be checked to see if he/she has the power to perform those actions.



### Personal settings

Request header must contain `X-User-Name` field.

* GET `/user/profile` Show a logged in staff's information.
* PATCH `/user/display-name` Change display name
* PATCH `/user/email` Change email
* POST `/user/password` Change password
* POST `/user/myft` Add a myft account
* DELETE `/user/myft/:id` Delete a myft account

### Admin

Request header must contain `X-User-Name` field and this user's privileges will be checked to see if he/she has the power to perform those actions.

* GET `/admin/staff/exists?k={name|email}&v={:value}`
* POST `/admin/staff/new`
* PUT `/admin/staff/new` Activate a deactivated user.
* GET `/admin/staff/roster` All staff

* GET `/admin/staff/profile` Show a staff's profile
* PATCH `/admin/staff/profile` Update staff's profile
* DELETE `/admin/staff/profile` Delete a staff

* GET `/admin/vip-roster` Show all myft accounts that are granted VIP.
* PUT `/admin/vip/{myftId}` Grant vip to a myft account
* DELETE `/admin/vip/{myftId}` Delete vip status of a myft account

### FTC Apps
* GET `/apps` Show urls
* POST `/apps/ftc` Create a new ftc app
* GET `/apps/ftc` Show all ftc apps. Anyone can see details of an app created by any others.
* GET `/apps/ftc/:appId` Show a ftc app
* POST `/apps/ftc/:appId` Only owner can edit it. So posted data should include owner id.
* DELETE `/apps/ftc/:appId`
* POST `/apps/ftc/:appId/reset-secret`
* POST `/apps/ftc/:appId/transfer`

* GET `/apps/ftc/:appId/tokens` List all tokens owned by this app

### CMS apps:

* POST `/apps/cms` Create a new cms app.
* GET `/apps/cms` List all cms apps.
* GET `/apps/cms/:appId` Show an app info.
* POST `/apps/cms/:appId` Update an app info owned by a user.
* DELETE `/apps/cms/:appId`
* POST `/apps/cms/:appId/transfer`

#### Permissions of CMS app
* GET `/apps/perms/:clientId` Get the unix permission of an app `clientId`. This is used by machines, not human.

### Personl Access Tokens
* GET `/tokens` Show urls
* POST `/tokens/personal` Create a new personal access token.
* GET `/tokens/personal/:ownerId` Show all access tokens a user owns
* DELETE `/tokens/personal/:ownerId` Revoke all access tokens
* GET `/tokens/personal/:ownerId/:tokenId` Show an access token
* POST `/token/personal/:ownerId/:tokenId` Update an access token
* DELETE `/token/personal/:ownerId/:tokenId` Delete an access token
* POST `/token/personal/:ownerId/:tokenId/regenerate`

### Client Access Tokens
Make sure current user owns the app which owns the token.

* DELETE `/tokens/client/:clientId` Revoke all tokens applied by an app
* DELETE `/tokens/client/:clientId/:tokenId` Revoke an access token owned by an app.
