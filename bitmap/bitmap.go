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
