package tencent_test

import (
	"context"
	"os"
	"testing"

	"github.com/Kirby980/study/week_2/internal/service/sms/tencent"
	"github.com/go-playground/assert/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func TestSender(t *testing.T) {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		t.Fatal()
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}

	s := tencent.NewService(c, "xxxx", "xxxx")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:   "发送验证码",
			tplId:  "1877556",
			params: []string{"123456"},
			// 改成你的手机号码
			numbers: []string{"10086"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}
}
