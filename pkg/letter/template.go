package letter

var templates = map[string]string{
	keySignUp: `
Dear {{.DisplayName}},

Welcome to join FTC.

The following is your credentials to sign in to FTC Content Management System.

Login name: {{.LoginName}}
Password: {{.Password}}

The password is an automatically generated random string. You're suggested to sign in the Content Management System and change it as soon as possible.

You can login via: {{.LoginURL}}.

This email contains sensitive data. Do not leak it to anyone else.

Thanks,
FTC Dev Team`,
	keyPwReset: `
{{.DisplayName}}

We heard that you lost your FTC CMS password. Sorry about that!

But don’t worry! You can use the following link to reset your password:

{{.URL}}

If you don’t use this link within 3 hours, it will expire. To get a new password reset link, visit http://superyard.ftchinese.com.

Thanks,
FTC Dev Team`,
}
