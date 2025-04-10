package tencent

import (
	"context"
	"fmt"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tplId
	req.PhoneNumberSet = s.toStringPtrString(number)
	req.TemplateParamSet = s.toStringPtrString(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrString(src []string) []*string {
	ptrStr := make([]*string, 0, len(src))
	for _, v := range src {
		ptrStr = append(ptrStr, &v)
	}
	return ptrStr
}
