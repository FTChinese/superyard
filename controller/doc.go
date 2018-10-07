/*
Package controller manages backyard-api's endpoints.

Herein ftc account refers to an account registered at www.ftchinese.com; staff account or CMS account refers to account used by ftc staff to access internal systems.

CMS Login

- Authenticate user's password:
	POST /staff/auth

- User forgot password and request a password-reset letter:
	POST /staff/password-reset/letter

- Verify a password reset link:
	GET /staff/password-reset/tokens/{token}

- User is allowed to reset password after password reset link is verified:
	POST /staff/password-reset

Personal settings

To access endpoints in this section, the staff must already logged in to CMS and the header in each request must contain X-User-Name field.

- User wants to see personal settings:
	GET /user/profile

- User wants to update his/her displayed name:
	PATCH /user/display-name

- User wants to change email address (shall we allow this?)
	PATCH /user/email

- User can change password:
	PATCH /user/password

- User wants to see the ftc accounts he/she owns:
	GET /user/myft

- User can link an ftc account to current CMS account:
	POST /user/myft

- User can unlink an ftc account from current CMS account:
	DELETE /user/myft/{id}

Admin

This section could only be accessed by a superuser (administrator). Request header must contain X-User-Name field and this user's privileges will be checked to see if he/she has the power to perform those actions.

- Check if a user with the specified user name or email exists:
	GET /staff/exists?k={name|email}&v={:value}

- Create a new CMS user:
	POST /admin/staff/new

- List all staff:
	GET /admin/staff/roster?page=<number>

- Show a staff's profile
	GET /admin/staff/profile/{name}

- Reinstate a previsouly deactivated CMS account:
	PUT /admin/staff/profile/{name}

- Update staff's profile (not including password)
	PATCH /admin/staff/profile/{name}

- Deactivate a staff's account:
	DELETE /admin/staff/profile/{name}?rmvip=true|false

- Show all ftc accounts that are granted VIP:
	GET /admin/vip

- Grant vip to an ftc account
	PUT /admin/vip/{myftId}

- Revoke VIP of an ftc account
	DELETE /admin/vip/{myftId}

Access FTC API

- Create a new ftc app
	POST /ftc-api/apps

- Show all ftc apps. There's previlege restrictions. Anyone logged in can view them:
	GET /ftc-api/apps?page=<number>

- Show the details of a ftc app:
	GET /ftc-api/apps/{name}

- Edit an app. Only owner can edit it. So posted data should include owner id:
	PATCH /ftc-api/apps/{name}

- Delete an app. Only owner can perform this action:
	DELETE /ftc-api/apps/{name}

- Transfer ownership to another CMS user:
	POST /ftc-api/apps/{name}/transfer

-  Create an access token. It could belong to a person or an app, depending on the data passed in:
	POST /ftc-api/tokens

- Show all access tokens a user owns:
	GET /ftc-api/tokens/personal

- Delete an access token:
	DELETE /ftc-api/token/personal/{tokenId}

- Show all access tokens owned by an app:
	GET /ftc-api/tokens/app/{name}

- Revoke an access token owned by an app:
	DELETE /ftc-api/tokens/app/{name}/{tokenId}

Customer Service

- Find an ftc user:
	GET /search/user?k=<name|email>&v=:value

- Show an ftc user's profile
	GET /ftc-user/profile/{userId}

- Show orders a user placed (may not acutally paid):
	GET /ftc-user/profile/{userId}/orders

- Show a user's login history:
	GET /ftc-user/profile/{userId}/login?page=<number>

Statistics

- Show the number of new user signup on a daily basis within the specifed date range:
	GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD
*/
package controller
