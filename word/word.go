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
	var jin, mu, shui, huo, tu int
	for _, d := range ws {
		switch d.Radicals {
		case "金", "钅", "几", "刀", "戈", "匕", "刂", "玉", "石", "皿", "西", "贝", "兑", "辛":
			fmt.Printf(" [%s] ", d.Word)
			jin++
		case "木", "|", "乙", "卩", "三", "弓", "东", "禾", "户", "门", "竹", "瓜", "舟":
			mu++
		case "水", "冫", "氵", "辶", "月", "子", "耳", "鱼", "黑", "雨", "川", "癸", "亥":
			shui++
		case "火", "丿", "乄", "忄", "心", "丙", "赤":
			huo++
		case "土", "阝", "一", "夕", "幺", "乎", "尸", "辰", "丑", "田", "良":
			tu++
		}
		if _, ok := wm[d.Word]; ok {
			wm[d.Word] = append(wm[d.Word], d)
		} else {
			wm[d.Word] = []IDict{d}
		}
	}
	fmt.Println(jin, mu, shui, huo, tu)
	return nil
}

/*
男金： 靖、铭、琛、川、承、司、斯、宗、骁、聪、在、钩、锦、铎、楚、铮、钦、则
女金： 真、心、新、悦、西、兮、楚、初、千、锐、素、锦、静、镜、斯、舒、瑜、童
男木： 楠、景、茗、聿、启、尧、言、嘉、桉、桐、筒、竹、林、乔、栋、家、翊、松
女木： 楠、景、茗、聿、启、尧、言、嘉、桉、桐、筒、竹、林、乔、栋、家、翊、松
男水： 清、澈、泫、浚、润、泽、向、凡、文、浦、洲、珩、玄、洋、淮、雨、子、云
女水： 妍、澜、淇、沐、潆、盈、雨、文、冰、雯、溪、子、云、汐、潞、淇、妙、涵
男火： 卓、昱、南、晨、知、宁、年、易、晗、炎、焕、哲、煦、旭、明、阳、朗、典
女火： 灿、夏、珞、煊、晴、彤、诺、宁、恬、钧、灵、昭、琉、晨、曦、南、毓、冉
男土： 辰、宸、野、安、为、亦，围、岚、也、以、延、允、容、恩、衡、宇、硕、已
女土： 意、也、坤、辰、伊、米、安、恩、以、容、宛、岚、又、衣、亚、悠、允、画

*/

func CheckWord(w string) []IDict {
	if d, ok := wm[w]; ok {
		return d
	}
	return nil
}
