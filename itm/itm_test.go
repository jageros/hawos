package itm

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type TPlayerMission struct {
	Id          uint64    `xorm:"'id' not null pk autoincr comment('自增ID') UNSIGNED BIGINT" json:"id"`
	PlayerId    uint64    `xorm:"'player_id' not null comment('玩家id') unique(player_id) UNSIGNED BIGINT" json:"player_id"`
	MissionType int       `xorm:"'mission_type' not null default 0 comment('任务类型') unique(player_id) INT" json:"mission_type"`
	Value       string    `xorm:"'value' not null comment('任务数据') TEXT" json:"value"`
	UpdateTime  int64     `xorm:"'update_time' not null default 0 comment('更新时间') INT" json:"update_time"`
	CreateTime  time.Time `xorm:"'create_time' default CURRENT_TIMESTAMP comment('创建时间') TIMESTAMP" json:"create_time"`
}

func TestMarshal(t *testing.T) {
	v := &TPlayerMission{
		Id:          1,
		PlayerId:    2,
		MissionType: 3,
		Value:       "{v:3}",
		UpdateTime:  time.Now().Unix(),
		CreateTime:  time.Now(),
	}
	str, err := json.Marshal(v)
	fmt.Println(string(str), err)
	vv := TPlayerMission{}
	json.Unmarshal(str, &vv)
	fmt.Println(vv)
}
