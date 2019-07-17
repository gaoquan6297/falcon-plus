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
	"fmt"
	"strings"

	"github.com/gaoquan6297/falcon-plus/common/model"
	"github.com/gaoquan6297/falcon-plus/common/utils"
	"github.com/gaoquan6297/falcon-plus/modules/alarm/g"
)

func BuildCommonSMSContent(event *model.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonIMContent(event *model.Event) string {
	status := ""
	if event.Status == "OK"{
		status = strings.Join([]string{"<font color=#008000 size=6>","</font>"},event.Status)
	} else{
		status = strings.Join([]string{"<font color=#FF0000 size=6>","</font>"},event.Status)
	}
	return fmt.Sprintf(
		"告警等级: P%d    \n 告警状态: %s    \n 告警主机: %s    \n 告警说明: %s    \n 告警指标: %s %s %s %s %s %s    \n 告警次数: %d    \n 告警时间: %s",
		event.Priority(),
		status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
	link := g.Link(event)
	return fmt.Sprintf(
		"%s\r\nP%d\r\nEndpoint:%s\r\nMetric:%s\r\nTags:%s\r\n%s: %s%s%s\r\nNote:%s\r\nMax:%d, Current:%d\r\nTimestamp:%s\r\n%s\r\n",
		event.Status,
		event.Priority(),
		event.Endpoint,
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		event.Func(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.Note(),
		event.MaxStep(),
		event.CurrentStep,
		event.FormattedTime(),
		link,
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateIMContent(event *model.Event) string {
	return BuildCommonIMContent(event)
}
