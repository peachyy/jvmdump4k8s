package notify

import (
	"fmt"
	"github.com/braumye/grobot"
	"jvmdump4k8s/config"
)

//https://github.com/braumye/grobot
func SendDingtalk(fileurl string) {
	token := config.GlobalConfig.NotifyDingToken
	podName := config.GlobalConfig.PodName
	fmt.Printf("开始推送钉钉 token %s %s \n", token, fileurl)
	robot, _ := grobot.New("dingtalk", token)
	// 发送文本消息
	err := robot.SendTextMessage(fmt.Sprintf("报警 %s应用发生OOM , dump文件%s ", podName, fileurl))
	fmt.Println("推送钉钉完成 err=", err)
}
