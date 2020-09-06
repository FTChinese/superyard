package letter

// Template keys
const (
	keySignUp             = "signUp"
	keyPwReset            = "passwordReset"
	keyOrderConfirmed     = "orderConfirmed"
	keyManualUpsertMember = "manualUpsertMember"
)

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
	keyOrderConfirmed: `
FT中文网用户 {{.Name}},

感谢您订阅FT中文网会员服务。

您于{{.OrderCreatedAt}}购买FT中文网会员的订单，由于某些故障我方未能及时获取支付完成信息，现已由FT中文网客服确认，会员信息已经更新。

本次订单信息

订单号 {{.OrderID}}
支付金额 {{.OrderAmount}}
支付方式 {{.PayMethod}}
订阅周期 {{.OrderStartDate}} 至 {{.OrderEndDate}}

最新会员状态：

会员类型 {{.Tier}}
到期时间 {{.ExpirationDate}}

如有疑问，请联系客服：subscriber.service@ftchinese.com。

再次感谢您对FT中文网的支持。

FT中文网`,
	keyManualUpsertMember: `
FT中文网用户 {{.Name}},

您的会员订阅状态已由客服更新。

最新会员状态：

会员类型 {{.Tier}}
到期时间 {{.ExpirationDate}}

如有疑问，请联系客服：subscriber.service@ftchinese.com。

感谢您对FT中文网的支持。

FT中文网`,
}
