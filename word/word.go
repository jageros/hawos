/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    word
 * @Date:    2021/12/17 6:07 下午
 * @package: word
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package word

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	ws []*dict
	wm map[string][]IDict
)

func init() {
	InitFromJsonFile("word.json")
}

type IDict interface {
	GetWord() string
	GetOldWord() string
	GetStrokes() string
	GetPinYin() string
	GetRadicals() string
	GetExplanation() string
	GetMore() string
}

type dict struct {
	Word        string `json:"word"`
	Oldword     string `json:"oldword"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"`
	More        string `json:"more"`
}

func (d *dict) GetWord() string        { return d.Word }
func (d *dict) GetOldWord() string     { return d.Oldword }
func (d *dict) GetStrokes() string     { return d.Strokes }
func (d *dict) GetPinYin() string      { return d.Pinyin }
func (d *dict) GetRadicals() string    { return d.Radicals }
func (d *dict) GetExplanation() string { return d.Explanation }
func (d *dict) GetMore() string        { return d.More }

func InitFromJsonFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	r := io.Reader(f)
	err = json.NewDecoder(r).Decode(&ws)

	if err != nil {
		return err
	}
	wm = map[string][]IDict{}
	index := 0
	for _, d := range ws {
		if d2, ok := wm[d.Word]; ok {
			index++
			if len(d2) >= 2 {
				fmt.Println("-----------------------")
			}
			fmt.Printf("%d: [%s %s/%s]\n", index, d.Word, d.Pinyin, d2[0].GetPinYin())
			wm[d.Word] = append(wm[d.Word], d)
		} else {
			wm[d.Word] = []IDict{d}
		}
	}
	return nil
}

func CheckWord(w string) []IDict {
	if d, ok := wm[w]; ok {
		return d
	}
	return nil
}
