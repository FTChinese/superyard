# Login

## Check if a staff exsits

GET `/staff/exists?k={name|email}&v={:value}`

## Login to CMS

POST `/staff/auth`

## Request a password reset letter

POST `/staff/password-reset/letter`

## Verify password reset link

GET `/staff/password-reset/tokens/{token}`

## Allow user to reset password

POST `/staff/password-reset`