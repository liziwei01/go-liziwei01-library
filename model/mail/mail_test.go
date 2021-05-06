package mail

import (
	"context"
	"fmt"
	"testing"
)

var ctx = context.Background()

func TestSendByAws(t *testing.T) {
	for i := 0; i < 1; i++ {
		var subject = fmt.Sprintf("[Check Code]")
		var mail = &AwsMail{
			From:        "alssylk@gmail.com",
			Product:     "aws",
			To:          []string{"ziweili1@link.cuhk.edu.cn"},
			Subject:     subject,
			TextContent: "123456",
			// Cc:          []string{"liziwei01@baidu.com"},
			// HTMLContent: "This message body contains HTML formatting. It can, for example, contain links like this one: <a class=\"ulink\" href=\"http://docs.aws.amazon.com/ses/latest/DeveloperGuide\" target=\"_blank\">Amazon SES Developer Guide</a>.",
		}
		err1 := mail.Send(ctx)
		fmt.Print(err1)
	}
}
