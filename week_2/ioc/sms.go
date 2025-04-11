package ioc

import (
	"github.com/Kirby980/study/week_2/internal/service/sms"
	"github.com/Kirby980/study/week_2/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 换内存，还是换别的
	return memory.NewService()
}
