# Mandrill API via Golang

[![Build Status](https://travis-ci.org/keighl/mandrill.png?branch=master)](https://travis-ci.org/keighl/mandrill) [![Coverage Status](https://coveralls.io/repos/keighl/mandrill/badge.svg)](https://coveralls.io/r/keighl/mandrill)

Stripped down package for sending emails through the Mandrill API. Inspired by [@mostafah's implementation](https://github.com/mostafah/mandrill).

### Installation

    go get -u github.com/kakysha/mandrill


### Documentation

http://godoc.org/github.com/kakysha/mandrill

### Regular Message

https://mandrillapp.com/api/docs/messages.JSON.html#method=send

```go
import (
    m "github.com/keighl/mandrill"
)

client := m.ClientWithKey("XXXXXXXXX")

message := &m.Message{}
message.AddRecipient("bob@example.com", "Bob Johnson", "to")
message.FromEmail = "kyle@example.com"
message.FromName = "Kyle Truscott"
message.Subject = "You won the prize!"
message.HTML = "<h1>You won!!</h1>"
message.Text = "You won!!"

responses, err := client.MessagesSend(message)
```

### Message Using Template (template_content version)

https://mandrillapp.com/api/docs/messages.JSON.html#method=send-template

```go
templateContent := map[string]string{"header": "Bob! You won the prize!"}
responses, err := client.MessagesSendTemplate(message, "you-won", templateContent)
```

### Message Using Template (per recipient merge_vars version)

https://mandrillapp.com/api/docs/messages.JSON.html#method=send-template

```go
// build per-recipient merge vars
merge_values := make(map[string]interface{})
merge_values["discount"] = site.Discount
merge_values["promocode"] = site.Promocode
merge_values["url_tag"] = site.Url_Tag
// bind vars to recipient's email
rcpt_merge_vars := mandrill.MapToRecipientVars(email, merge_values)
// append this recipient vars to all merge_vars
all_merge_vars = append(all_merge_vars, rcpt_merge_vars)
// fill template struct and send
message.MergeVars = all_merge_vars
responses, err := client.MessagesSendTemplate(message, "template-slug", nil)
```

### Templates API

https://mandrillapp.com/api/docs/templates.JSON.html

- AddTemplate()
- UpdateTemplate()
- DeleteTemplate()
- TemplateInfo()

```go
template := &m.Template{
	Name:      "Example Template",
	Subject:   "Account Activation",
	HTML:      "<h1>You've signed up for *|SITE_NAME|*</h1>",
	Text: "You've signed up for *|SITE_NAME|*",
	FromEmail: "noreply@example.com",
	FromName:  "Site Account System",
	Labels: []string{"account", "account-activation"}
}

res, err := client.UpdateTemplate(template)
if err != nil {
	res, err = client.AddTemplate(template)
}
```

### Subaccounts API

https://mandrillapp.com/api/docs/subaccounts.JSON.html

- AddSubaccount()
- UpdateSubaccount()
- DeleteSubaccount()
- SubaccountInfo()

```go
subaccount := &m.Subaccount{
	Id:           "cust-123",
	Name:         "ABC Widgets, Inc.",
	Notes:        "Free plan user, signed up on 2013-01-01 12:00:00",
	Quota:        42,
}

res, err := client.UpdateSubaccount(subaccount)
if err != nil {
	res, err = client.AddSubaccount(subaccount)
}
```

### Including Merge Tags

http://help.mandrill.com/entries/21678522-How-do-I-use-merge-tags-to-add-dynamic-content-

```go
// Global vars
message.GlobalMergeVars = m.MapToVars(map[string]interface{}{"name": "Bob"})

// Recipient vars
bobVars := m.MapToRecipientVars("bob@example.com", map[string]interface{}{"name": "Bob"})
jillVars := m.MapToRecipientVars("jill@example.com", map[string]interface{}{"name": "Jill"})
message.MergeVars = []*m.RcptMergeVars{bobVars, jillVars}
```

### Integration Testing Keys

You can pass special API keys to the client to mock success/err responses from `MessagesSend` or `MessagesSendTemplate`.

```go
// Sending messages will be successful, but without a real API request
c := ClientWithKey("SANDBOX_SUCCESS")

// Sending messages will error, but without a real API request
c := ClientWithKey("SANDBOX_ERROR")
```


