package notify

import (
	"fmt"
	"github.com/braumye/grobot"
	"jvmdump4k8s/config"
)

//https://github.com/braumye/grobot
func SendWechat(fileurl string) {
	token := config.GlobalConfig.NotifyWxToken
	podName := config.GlobalConfig.PodName
	fmt.Printf("开始推送微信 token %s %s \n", token, fileurl)
	robot, _ := grobot.New("wechatwork", token)
	// 发送文本消息
	err := robot.SendTextMessage(fmt.Sprintf("报警 %s应用发生OOM , dump文件%s ", podName, fileurl))
	fmt.Println("推送微信完成 err=", err)
}
