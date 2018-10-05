# CMS User Settings

Request header must contain `X-User-Name` field.

* GET `/user/profile` Show a logged in staff's information.
* PATCH `/user/display-name` Change display name
* PATCH `/user/email` Change email
* PATCH `/user/password` Change password

* GET `/user/myft` List all myft accounts
* POST `/user/myft` Add a myft account
* DELETE `/user/myft/:id` Delete a myft account