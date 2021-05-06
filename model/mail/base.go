/*
 * @Author: 		liziwei01
 * @Date: 			2021-04-19 17:00:00
 * @LastEditTime: 	2021-04-19 17:00:00
 * @LastEditors: 	liziwei01
 * @Description: 	main
 * @FilePath: 		github.com/go-liziwei01-library/main.go
 */

package mail

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/go-liziwei01-library/library/conf"
)

const (
	awsConfigFileName = "%s.toml"
	charset           = "UTF-8"
	configPath        = "./conf"
)

// WorkMail 邮件配置
type WorkMail struct {
	Server string `toml:"server"`
	Port   int    `toml:"port"`
	Pwd    string `toml:"Pwd"`
	User   string `toml:"user"`
}

// AwsConfig aws配置信息结构体
type AwsConfig struct {
	Ak       string   `toml:"ak"`
	Sk       string   `toml:"sk"`
	Region   string   `toml:"region"`
	Bucket   string   `toml:"bucket"`
	ACL      string   `toml:"acl"`
	SesMail  string   `toml:"sesMail"`
	WorkMail WorkMail `toml:"workMail"`
}

// getConfig 获取aws配置信息
func getAwsConfig(ctx context.Context, product string) (*AwsConfig, error) {
	// 配置文件跟路径
	awsConfPath := filepath.Join(configPath, fmt.Sprintf(awsConfigFileName, product))
	var cfg *AwsConfig
	if err := conf.Parse(awsConfPath, &cfg); err != nil {
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email getAwsConfig error: %s", err.Error()))
		return nil, err
	}
	return cfg, nil
}

// AwsMail 邮件
type AwsMail struct {
	Product     string   `json:"product"`
	From        string   `json:"from"`
	To          []string `json:"to"`
	Cc          []string `json:"cc,omitempty"`
	HTMLContent string   `json:"htmlContent"`
	TextContent string   `json:"textContent"`
	Attach      string   `json:"attach"`
	Subject     string   `json:"subject"`
}

// AwsSesMail 邮件
type AwsSesMail struct {
	From        *string   `json:"from"`
	To          []*string `json:"to"`
	Cc          []*string `json:"cc,omitempty"`
	Subject     *string   `json:"subject"`
	HTMLContent *string   `json:"htmlContent"`
	TextContent *string   `json:"textContent"`
	Attach      *string   `json:"attach"`
}

// Send 发送邮件
func (m *AwsMail) Send(ctx context.Context) error {
	awsConfig, err := getAwsConfig(ctx, m.Product)
	if err != nil {
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email Send getAwsConfig error: %s", err.Error()))
		return err
	}
	sm, err := praseSes(ctx, m)
	if err != nil {
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email Send praseSes error: %s", err.Error()))
		return err
	}
	creds := credentials.NewStaticCredentials(awsConfig.Ak, awsConfig.Sk, "")
	_, err = creds.Get()
	if err != nil {
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email Send creds error: %s", err.Error()))
		return err
	}
	cfg := aws.NewConfig().WithRegion(awsConfig.Region).WithCredentials(creds)
	svc := ses.New(session.New(), cfg)
	sesBody := &ses.Body{}
	if sm.HTMLContent != nil && *(sm.HTMLContent) != "" {
		sesBody.Html = &ses.Content{
			Charset: aws.String(charset),
			Data:    sm.HTMLContent,
		}
	}
	if sm.TextContent != nil && *(sm.TextContent) != "" {
		sesBody.Text = &ses.Content{
			Charset: aws.String(charset),
			Data:    sm.TextContent,
		}
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: sm.Cc,
			ToAddresses: sm.To,
		},
		Message: &ses.Message{
			Body: sesBody,
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    sm.Subject,
			},
		},
		Source: sm.From,
	}
	// result, err := svc.SendEmail(input)
	_, err = svc.SendEmail(input)
	if err != nil {
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email Send SendEmail error: %s", err.Error()))
		return err
	}
	// logs.SvrLogger.Notice(ctx, fmt.Sprintf("email Send SendEmail success: %s", result))
	return nil
}

// praseSes 转换参数
func praseSes(ctx context.Context, m *AwsMail) (*AwsSesMail, error) {
	awsConfig, err := getAwsConfig(ctx, m.Product)
	if err != nil {
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email praseSes getAwsConfig error: %s", err.Error()))
		return nil, err
	}
	if len(m.To) == 0 {
		err = fmt.Errorf("")
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email praseSes getAwsConfig error: %s", err.Error()))
		return nil, err
	}
	if m.Subject == "" {
		err = fmt.Errorf("email subject is empty")
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email praseSes getAwsConfig error: %s", err.Error()))
		return nil, err
	}
	if m.HTMLContent == "" && m.TextContent == "" {
		err = fmt.Errorf("email Content is empty")
		// logs.SvrLogger.Warning(ctx, fmt.Sprintf("email praseSes getAwsConfig error: %s", err.Error()))
		return nil, err
	}
	sm := &AwsSesMail{}
	sm.Subject = aws.String(m.Subject)
	sm.To = make([]*string, 0)
	for _, to := range m.To {
		sm.To = append(sm.To, aws.String(to))
	}
	sm.Cc = make([]*string, 0)
	for _, cc := range m.Cc {
		sm.Cc = append(sm.To, aws.String(cc))
	}
	sm.HTMLContent = aws.String(m.HTMLContent)
	sm.TextContent = aws.String(m.TextContent)
	sm.Attach = aws.String(m.Attach)
	sm.From = aws.String(m.From)
	if m.From == "" {
		sm.From = aws.String(awsConfig.SesMail)
	}
	return sm, nil
}
