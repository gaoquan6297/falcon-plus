// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/gaoquan6297/falcon-plus/modules/alarm/g"
	"github.com/gaoquan6297/falcon-plus/modules/alarm/model"
	"github.com/gaoquan6297/falcon-plus/modules/alarm/redi"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func ConsumeSms() {
	for {
		L := redi.PopAllSms()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendSmsList(L)
	}
}

func SendSmsList(L []*model.Sms) {
	for _, sms := range L {
		SmsWorkerChan <- 1
		go SendSms(sms)
	}
}

type CallRl struct {
	To           string `json:"to"`
	MediaName    string `json:"mediaName"`
	MediaTxt     string `json:"mediaTxt"`
	AppId        string `json:"appId"`
	DisplayNum   string `json:"displayNum"`
	PlayTimes    string `json:"playTimes"`
	RespUrl      string `json:"respUrl"`
	UserData     string `json:"userData"`
	MaxCallTime  string `json:"maxCallTime"`
	Speed        string `json:"speed"`
	Volume       string `json:"volume"`
	Pitch        string `json:"pitch"`
	Bgsound      string `json:"bgsound"`
}

func SendSms(sms *model.Sms) {
	defer func() {
		<-SmsWorkerChan
	}()
	accountSid   := g.Config().Call.AccountSid
	accountToken := g.Config().Call.AccountToken
	appId        := g.Config().Call.AppId
	mediaText    := g.Config().Call.MediaText
	nowdate      := time.Now().Format("20060102150405")
	signature := strings.Join([]string{accountSid, accountToken,nowdate}, "")
	h := md5.New()
	h.Write([]byte(signature))
	sig := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	// api.sms: https://app.cloopen.com:8883/2013-12-26/Accounts/
	url := strings.Join([]string{g.Config().Api.Sms,accountSid,"/Calls/LandingCalls?sig=",sig},"")
	call_msg := CallRl{
		sms.Tos,
		"",
		mediaText,
		appId,
		"",
		"2",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}
	auth := base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{accountSid, nowdate},":")))
	jsonCall, err := json.Marshal(call_msg)
	if err != nil {
		log.Errorf("生成json字符串错误")
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonCall))
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Authorization", auth)
	if err != nil {
		log.Errorf("The request error happened: %s",err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Errorf("The response error happened: %s",err)
	}

	msg, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("response body content: %s", msg)
	}
	//url := g.Config().Api.Sms
	//r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	//r.Param("tos", sms.Tos)
	//r.Param("content", sms.Content)
	//resp, err := r.String()
	//if err != nil {
	//	log.Errorf("send sms fail, tos:%s, content:%s, error:%v", sms.Tos, sms.Content, err)
	//}
	//log.Debugf("send sms:%v, resp:%v, url:%s", sms, resp, url)

}
