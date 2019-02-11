package staff

const SignupLetter = `
Dear {{.NormalizeName}},

Welcome to join FTC.

The following is your credentials to sign in to FTC Content Management System.

Login name: {{.UserName}}
Password: {{.GetPassword}}

The password is an automatically generated random string. You're suggested to sign in the Content Management System and change it as soon as possible.

You can login via: http://superyard.ftchinese.com.

This email contains sensitive data. Do not leak it to anyone else.

Thanks,
FTC Dev Team`

const PasswordResetLetter = `
{{.NormalizeName}}

We heard that you lost your FTC CMS password. Sorry about that!

But don’t worry! You can use the following link to reset your password:

http://superyard.ftchinese.com/password-reset/{{.Token}}

If you don’t use this link within 3 hours, it will expire. To get a new password reset link, visit http://superyard.ftchinese.com.

Thanks,
FTC Dev Team`