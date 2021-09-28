/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    consts
 * @Date:    2021/8/26 5:51 下午
 * @package: consts
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package consts

const (
	MaxInt64  = int64(^uint64(0) >> 1)
	MaxUint64 = ^uint64(0)
	MaxInt    = int(^uint(0) >> 1)
	MaxUint   = ^uint(0)
	MaxInt32  = int32(^uint32(0) >> 1)
	MaxUint32 = ^uint32(0)
	MaxInt16  = int16(^uint16(0) >> 1)
	MaxUint16 = ^uint16(0)
	MaxInt8   = int8(^uint8(0) >> 1)
	MaxUint8  = ^uint8(0)
)
