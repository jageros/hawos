/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    tbl
 * @Date:    2021/9/2 4:00 下午
 * @package: gamedata
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package xlsx

import "github.com/jager/hawox/evq"

var (
	XLSX_DATA_RELOAD_EVENT = evq.CreateEventID()
)

const (
	TblRoomConfig = "room_config"
)
