/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    bitmap
 * @Date:    2021/11/24 2:02 下午
 * @package: bitmap
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package bitmap

type Bitmap struct {
	words  []uint64
	length int
}

func New() *Bitmap {
	return &Bitmap{}
}

func (bm *Bitmap) Has(num int) bool {
	word, bit := num/64, uint(num%64)
	return word < len(bm.words) && (bm.words[word]&(1<<bit)) != 0
}

func (bm *Bitmap) Add(num int) {
	word, bit := num/64, uint(num%64)
	for word >= len(bm.words) {
		bm.words = append(bm.words, 0)
	}
	// 判断num是否已经存在bitmap中
	if bm.words[word]&(1<<bit) == 0 {
		bm.words[word] |= 1 << bit
		bm.length++
	}
}

func (bm *Bitmap) Len() int {
	return bm.length
}

func (bm *Bitmap) Clear() {
	bm.words = nil
	bm.length = 0
}

//func (bitmap *Bitmap) Len() int {
//	return bitmap.length
//}
//
//func (bitmap *Bitmap) String() string {
//	var buf bytes.Buffer
//	buf.WriteByte('{')
//	for i, v := range bitmap.words {
//		if v == 0 {
//			continue
//		}
//		for j := uint(0); j < 64; j++ {
//			if v&(1<<j) != 0 {
//				if buf.Len() > len("{") {
//					buf.WriteByte(' ')
//				}
//				fmt.Fprintf(&buf, "%d", 64*uint(i)+j)
//			}
//		}
//	}
//	buf.WriteByte('}')
//	fmt.Fprintf(&buf,"\nLength: %d", bitmap.length)
//	return buf.String()
//}
