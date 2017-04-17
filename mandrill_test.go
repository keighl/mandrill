package mandrill

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %[1]T) - Got %v (type %[2]T)", b, a)
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %[1]T) - Got %v (type %[2]T)", b, a)
	}
}

func testTools(code int, body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: tr}

	client := &Client{"APIKEY", server.URL + "/", httpClient}
	return server, client
}

// ClientWithKey //////

func Test_ClientWithKey(t *testing.T) {
	c := ClientWithKey("CHEEEEEESE")
	refute(t, c, nil)
}

// MessagesSendTemplate //////////

func Test_MessagesSendTemplate_Success(t *testing.T) {
	server, m := testTools(200, `[{"email":"bob@example.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()
	responses, err := m.MessagesSendTemplate(&Message{}, "cheese", map[string]string{"name": "bob"})

	expect(t, len(responses), 1)
	expect(t, err, nil)

	correctMessagesResponse := &MessagesResponse{
		Email:           "bob@example.com",
		Status:          "sent",
		RejectionReason: "hard-bounce",
		Id:              "1",
	}
	expect(t, reflect.DeepEqual(correctMessagesResponse, responses[0]), true)
}

func Test_MessagesSendTemplate_Fail(t *testing.T) {
	server, m := testTools(400, `{"status":"error","code":12,"name":"Unknown_Subaccount","message":"No subaccount exists with the id 'customer-123'"}`)
	defer server.Close()
	responses, err := m.MessagesSendTemplate(&Message{}, "cheese", map[string]string{"name": "bob"})

	expect(t, len(responses), 0)

	correctMessagesResponse := &Error{
		Status:  "error",
		Code:    12,
		Name:    "Unknown_Subaccount",
		Message: "No subaccount exists with the id 'customer-123'",
	}
	expect(t, reflect.DeepEqual(correctMessagesResponse, err), true)
}

// MessagesSend //////////

func Test_MessageSend_Success(t *testing.T) {
	server, m := testTools(200, `[{"email":"bob@example.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)
	defer server.Close()
	responses, err := m.MessagesSend(&Message{})

	expect(t, len(responses), 1)
	expect(t, err, nil)

	correctMessagesResponse := &MessagesResponse{
		Email:           "bob@example.com",
		Status:          "sent",
		RejectionReason: "hard-bounce",
		Id:              "1",
	}
	expect(t, reflect.DeepEqual(correctMessagesResponse, responses[0]), true)
}

func Test_MessageSend_Fail(t *testing.T) {
	server, m := testTools(400, `{"status":"error","code":12,"name":"Unknown_Subaccount","message":"No subaccount exists with the id 'customer-123'"}`)
	defer server.Close()
	responses, err := m.MessagesSend(&Message{})

	expect(t, len(responses), 0)

	correctMessagesResponse := &Error{
		Status:  "error",
		Code:    12,
		Name:    "Unknown_Subaccount",
		Message: "No subaccount exists with the id 'customer-123'",
	}
	expect(t, reflect.DeepEqual(correctMessagesResponse, err), true)
}

// Ping //////////

func Test_Ping_Success(t *testing.T) {
	server, m := testTools(200, `"PONG!"`)
	defer server.Close()
	response, err := m.Ping()

	expect(t, response, "PONG!")
	expect(t, err, nil)
}

func Test_Ping_Fail(t *testing.T) {
	server, m := testTools(400, `{"status":"error","code":-1,"name":"Invalid_Key","message":"Invalid API key"}`)
	defer server.Close()
	response, err := m.Ping()

	expect(t, response, "")

	correctMessagesResponse := &Error{
		Status:  "error",
		Code:    -1,
		Name:    "Invalid_Key",
		Message: "Invalid API key",
	}
	expect(t, reflect.DeepEqual(correctMessagesResponse, err), true)
}

// AddTemplate //////////

func Test_AddTemplate_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "slug": "example-template",
	    "name": "Example Template",
	    "labels": [
	        "example-label"
	    ],
	    "code": "<div mc:edit=\"editable\">editable content</div>",
	    "subject": "example subject",
	    "from_email": "from.email@example.com",
	    "from_name": "Example Name",
	    "text": "Example text",
	    "publish_name": "Example Template",
	    "publish_code": "<div mc:edit=\"editable\">different than draft content</div>",
	    "publish_subject": "example publish_subject",
	    "publish_from_email": "from.email.published@example.com",
	    "publish_from_name": "Example Published Name",
	    "publish_text": "Example published text",
	    "published_at": "2013-01-01 15:30:40",
	    "created_at": "2013-01-01 15:30:27",
	    "updated_at": "2013-01-01 15:30:49"
	}`)
	defer server.Close()
	response, err := m.AddTemplate(&Template{})

	expect(t, err, nil)

	correctResponse := &Template{
		Name:      "Example Template",
		Slug:      "example-template",
		Subject:   "example subject",
		FromEmail: "from.email@example.com",
		FromName:  "Example Name",
		HTML:      "<div mc:edit=\"editable\">editable content</div>",
		Text:      "Example text",
		Labels:    []string{"example-label"},
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_AddTemplate_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": -1,
	    "name": "Invalid_Key",
	    "message": "Invalid API key"
	}`)
	defer server.Close()
	_, err := m.AddTemplate(&Template{})

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    -1,
		Name:    "Invalid_Key",
		Message: "Invalid API key",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// UpdateTemplate //////////

func Test_UpdateTemplate_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "slug": "example-template",
	    "name": "Example Template",
	    "labels": [
	        "example-label"
	    ],
	    "code": "<div mc:edit=\"editable\">editable content</div>",
	    "subject": "example subject",
	    "from_email": "from.email@example.com",
	    "from_name": "Example Name",
	    "text": "Example text",
	    "publish_name": "Example Template",
	    "publish_code": "<div mc:edit=\"editable\">different than draft content</div>",
	    "publish_subject": "example publish_subject",
	    "publish_from_email": "from.email.published@example.com",
	    "publish_from_name": "Example Published Name",
	    "publish_text": "Example published text",
	    "published_at": "2013-01-01 15:30:40",
	    "created_at": "2013-01-01 15:30:27",
	    "updated_at": "2013-01-01 15:30:49"
	}`)
	defer server.Close()
	response, err := m.UpdateTemplate(&Template{})

	expect(t, err, nil)

	correctResponse := &Template{
		Name:      "Example Template",
		Slug:      "example-template",
		Subject:   "example subject",
		FromEmail: "from.email@example.com",
		FromName:  "Example Name",
		HTML:      "<div mc:edit=\"editable\">editable content</div>",
		Text:      "Example text",
		Labels:    []string{"example-label"},
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_UpdateTemplate_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": -1,
	    "name": "Invalid_Key",
	    "message": "Invalid API key"
	}`)
	defer server.Close()
	_, err := m.UpdateTemplate(&Template{})

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    -1,
		Name:    "Invalid_Key",
		Message: "Invalid API key",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// DeleteTemplate //////////

func Test_DeleteTemplate_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "slug": "example-template",
	    "name": "Example Template",
	    "labels": [
	        "example-label"
	    ],
	    "code": "<div mc:edit=\"editable\">editable content</div>",
	    "subject": "example subject",
	    "from_email": "from.email@example.com",
	    "from_name": "Example Name",
	    "text": "Example text",
	    "publish_name": "Example Template",
	    "publish_code": "<div mc:edit=\"editable\">different than draft content</div>",
	    "publish_subject": "example publish_subject",
	    "publish_from_email": "from.email.published@example.com",
	    "publish_from_name": "Example Published Name",
	    "publish_text": "Example published text",
	    "published_at": "2013-01-01 15:30:40",
	    "created_at": "2013-01-01 15:30:27",
	    "updated_at": "2013-01-01 15:30:49"
	}`)
	defer server.Close()
	response, err := m.DeleteTemplate("")

	expect(t, err, nil)

	correctResponse := &Template{
		Name:      "Example Template",
		Slug:      "example-template",
		Subject:   "example subject",
		FromEmail: "from.email@example.com",
		FromName:  "Example Name",
		HTML:      "<div mc:edit=\"editable\">editable content</div>",
		Text:      "Example text",
		Labels:    []string{"example-label"},
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_DeleteTemplate_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": -1,
	    "name": "Invalid_Key",
	    "message": "Invalid API key"
	}`)
	defer server.Close()
	_, err := m.DeleteTemplate("")

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    -1,
		Name:    "Invalid_Key",
		Message: "Invalid API key",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// TemplateInfo //////////

func Test_TemplateInfo_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "slug": "example-template",
	    "name": "Example Template",
	    "labels": [
	        "example-label"
	    ],
	    "code": "<div mc:edit=\"editable\">editable content</div>",
	    "subject": "example subject",
	    "from_email": "from.email@example.com",
	    "from_name": "Example Name",
	    "text": "Example text",
	    "publish_name": "Example Template",
	    "publish_code": "<div mc:edit=\"editable\">different than draft content</div>",
	    "publish_subject": "example publish_subject",
	    "publish_from_email": "from.email.published@example.com",
	    "publish_from_name": "Example Published Name",
	    "publish_text": "Example published text",
	    "published_at": "2013-01-01 15:30:40",
	    "created_at": "2013-01-01 15:30:27",
	    "updated_at": "2013-01-01 15:30:49"
	}`)
	defer server.Close()
	response, err := m.TemplateInfo("")

	expect(t, err, nil)

	correctResponse := &Template{
		Name:      "Example Template",
		Slug:      "example-template",
		Subject:   "example subject",
		FromEmail: "from.email@example.com",
		FromName:  "Example Name",
		HTML:      "<div mc:edit=\"editable\">editable content</div>",
		Text:      "Example text",
		Labels:    []string{"example-label"},
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_TemplateInfo_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": 5,
	    "name": "Unknown_Template",
	    "message": "No such template \"Example Template\""
	}`)
	defer server.Close()

	_, err := m.TemplateInfo("")

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    5,
		Name:    "Unknown_Template",
		Message: "No such template \"Example Template\"",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// AddSubaccount //////////

func Test_AddSubaccount_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "id": "cust-123",
	    "name": "ABD Widgets, Inc.",
	    "custom_quota": 42,
	    "status": "active",
	    "reputation": 42,
	    "created_at": "2013-01-01 15:30:27",
	    "first_sent_at": "2013-01-01 15:30:29",
	    "sent_weekly": 42,
	    "sent_monthly": 42,
	    "sent_total": 42
	}`)
	defer server.Close()

	response, err := m.AddSubaccount(&Subaccount{})

	expect(t, err, nil)

	correctResponse := &Subaccount{
		Id:           "cust-123",
		Name:         "ABD Widgets, Inc.",
		Quota:        42,
		Reputation:   42,
		Status:       "active",
		Sent_hourly:  0,
		Sent_weekly:  42,
		Sent_monthly: 42,
		Sent_total:   42,
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_AddSubaccount_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": -1,
	    "name": "Invalid_Key",
	    "message": "Invalid API key"
	}`)
	defer server.Close()

	_, err := m.AddSubaccount(&Subaccount{})

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    -1,
		Name:    "Invalid_Key",
		Message: "Invalid API key",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// UpdateSubaccount //////////

func Test_UpdateSubaccount_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "id": "cust-123",
	    "name": "ABD Widgets, Inc.",
	    "custom_quota": 43,
	    "status": "active",
	    "reputation": 43,
	    "created_at": "2013-01-01 15:30:27",
	    "first_sent_at": "2013-01-01 15:30:29",
	    "sent_weekly": 43,
	    "sent_monthly": 43,
	    "sent_total": 43
	}`)
	defer server.Close()

	response, err := m.UpdateSubaccount(&Subaccount{})

	expect(t, err, nil)

	correctResponse := &Subaccount{
		Id:           "cust-123",
		Name:         "ABD Widgets, Inc.",
		Quota:        43,
		Reputation:   43,
		Status:       "active",
		Sent_hourly:  0,
		Sent_weekly:  43,
		Sent_monthly: 43,
		Sent_total:   43,
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_UpdateSubaccount_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": -1,
	    "name": "Invalid_Key",
	    "message": "Invalid API key"
	}`)
	defer server.Close()

	_, err := m.UpdateSubaccount(&Subaccount{})

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    -1,
		Name:    "Invalid_Key",
		Message: "Invalid API key",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// DeleteSubaccount //////////

func Test_DeleteSubaccount_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "id": "cust-123",
	    "name": "ABD Widgets, Inc.",
	    "custom_quota": 43,
	    "status": "active",
	    "reputation": 43,
	    "created_at": "2013-01-01 15:30:27",
	    "first_sent_at": "2013-01-01 15:30:29",
	    "sent_weekly": 43,
	    "sent_monthly": 43,
	    "sent_total": 43
	}`)
	defer server.Close()

	response, err := m.DeleteSubaccount("")

	expect(t, err, nil)

	correctResponse := &Subaccount{
		Id:           "cust-123",
		Name:         "ABD Widgets, Inc.",
		Quota:        43,
		Reputation:   43,
		Status:       "active",
		Sent_hourly:  0,
		Sent_weekly:  43,
		Sent_monthly: 43,
		Sent_total:   43,
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_DeleteSubaccount_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": 12,
	    "name": "Unknown_Subaccount",
	    "message": "No subaccount exists with the id 'customer-123'"
	}`)
	defer server.Close()

	_, err := m.DeleteSubaccount("")

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    12,
		Name:    "Unknown_Subaccount",
		Message: "No subaccount exists with the id 'customer-123'",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// SubaccountInfo //////////

func Test_SubaccountInfo_Success(t *testing.T) {
	server, m := testTools(200, `{
	    "id": "cust-123",
	    "name": "ABC Widgets, Inc.",
	    "notes": "Free plan user, signed up on 2013-01-01 12:00:00",
	    "custom_quota": 42,
	    "status": "active",
	    "reputation": 42,
	    "created_at": "2013-01-01 15:30:27",
	    "first_sent_at": "2013-01-01 15:30:29",
	    "sent_weekly": 42,
	    "sent_monthly": 42,
	    "sent_total": 42,
	    "sent_hourly": 42,
	    "hourly_quota": 42,
	    "last_30_days": {
	        "sent": 42,
	        "hard_bounces": 42,
	        "soft_bounces": 42,
	        "rejects": 42,
	        "complaints": 42,
	        "unsubs": 42,
	        "opens": 42,
	        "unique_opens": 42,
	        "clicks": 42,
	        "unique_clicks": 42
	    }
	}`)
	defer server.Close()
	response, err := m.SubaccountInfo("")

	expect(t, err, nil)

	correctResponse := &Subaccount{
		Id:           "cust-123",
		Name:         "ABC Widgets, Inc.",
		Notes:        "Free plan user, signed up on 2013-01-01 12:00:00",
		Quota:        42,
		Reputation:   42,
		Status:       "active",
		Sent_hourly:  42,
		Sent_weekly:  42,
		Sent_monthly: 42,
		Sent_total:   42,
	}
	expect(t, reflect.DeepEqual(correctResponse, response), true)
}

func Test_SubaccountInfo_Fail(t *testing.T) {
	server, m := testTools(400, `{
	    "status": "error",
	    "code": 12,
	    "name": "Unknown_Subaccount",
	    "message": "No subaccount exists with the id 'customer-123'"
	}`)
	defer server.Close()
	_, err := m.SubaccountInfo("")

	refute(t, err, nil)

	correctError := &Error{
		Status:  "error",
		Code:    12,
		Name:    "Unknown_Subaccount",
		Message: "No subaccount exists with the id 'customer-123'",
	}
	expect(t, reflect.DeepEqual(correctError, err), true)
}

// TEST Keys //////////

func Test_SANDBOX_SUCCESS(t *testing.T) {
	client := ClientWithKey("SANDBOX_SUCCESS")
	_, err := client.MessagesSend(&Message{})
	expect(t, err, nil)
}

func Test_SANDBOX_ERROR(t *testing.T) {
	client := ClientWithKey("SANDBOX_ERROR")
	_, err := client.MessagesSend(&Message{})
	refute(t, err, nil)
}

// AddRecipient //////////

func Test_AddRecipient(t *testing.T) {
	m := &Message{}
	m.AddRecipient("bob@example.com", "Bob Johnson", "to")
	tos := []*To{&To{"bob@example.com", "Bob Johnson", "to"}}
	expect(t, reflect.DeepEqual(m.To, tos), true)
}

// ConvertMapToVariables /////

func Test_ConvertMapToVariables(t *testing.T) {
	m := map[string]interface{}{"name": "bob"}
	target := ConvertMapToVariables(m)
	hand := []*Variable{&Variable{"name", "bob"}}
	expect(t, reflect.DeepEqual(target, hand), true)
}

func Test_ConvertMapToVariables_WithString(t *testing.T) {
	m := map[string]string{"name": "bob"}
	target := ConvertMapToVariables(m)
	hand := []*Variable{&Variable{"name", "bob"}}
	expect(t, reflect.DeepEqual(target, hand), true)
}

func Test_ConvertMapToVariables_BadType(t *testing.T) {
	target := ConvertMapToVariables("CHEESE")
	expect(t, len(target), 0)
}

func Test_MapToVars(t *testing.T) {
	m := map[string]interface{}{"name": "bob"}
	target := MapToVars(m)
	hand := []*Variable{&Variable{"name", "bob"}}
	expect(t, reflect.DeepEqual(target, hand), true)
}

// ConvertMapToVariablesForRecipient ////

func Test_ConvertMapToVariablesForRecipient(t *testing.T) {
	m := map[string]string{"name": "bob"}
	target := ConvertMapToVariablesForRecipient("bob@example.com", m)
	hand := &RcptMergeVars{"bob@example.com", ConvertMapToVariables(m)}
	expect(t, reflect.DeepEqual(target, hand), true)
}

func Test_MapToRecipientVars(t *testing.T) {
	m := map[string]string{"name": "bob"}
	target := MapToRecipientVars("bob@example.com", m)
	hand := &RcptMergeVars{"bob@example.com", ConvertMapToVariables(m)}
	expect(t, reflect.DeepEqual(target, hand), true)
}

// Error Interface ////

func Test_ErrorError(t *testing.T) {
	e := Error{Message: "CHEEEEEESE"}
	expect(t, e.Error(), "CHEEEEEESE")
}
