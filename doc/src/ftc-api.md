## Manage FTC API 

### Apps using FTC API

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