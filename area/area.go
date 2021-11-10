/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    area
 * @Date:    2021/11/9 5:35 下午
 * @package: area
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package area

import (
	"encoding/json"
)

type ICounty interface {
	GetCode() string
	GetName() string
}

type ICity interface {
	ICounty
	GetCounties() []ICounty
}

type IProvince interface {
	ICounty
	GetCities() []ICity
}

type County struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (a *County) GetCode() string {
	return a.Code
}

func (a *County) GetName() string {
	return a.Name
}

type City struct {
	County
	Counties  []*County `json:"children"`
	iCounties []ICounty
}

func (c *City) GetCounties() []ICounty {
	return c.iCounties
}

type Province struct {
	County
	Cities  []*City `json:"children"`
	iCities []ICity
}

func (p *Province) GetCities() []ICity {
	return p.iCities
}

var provinces []*Province
var iProvinces []IProvince

func init() {
	if err := json.Unmarshal([]byte(citystr), &provinces); err != nil {
		panic(err)
	}
	for _, p := range provinces {
		for _, city := range p.Cities {
			for _, county := range city.Counties {
				city.iCounties = append(city.iCounties, county)
			}
			p.iCities = append(p.iCities, city)
		}
		iProvinces = append(iProvinces, p)
	}
}

func GetProvinces() []IProvince {
	return iProvinces
}

func GetCities(province string) []ICity {
	for _, p := range iProvinces {
		if p.GetName() == province {
			return p.GetCities()
		}
	}
	return nil
}

func GetCounties(province, city string) []ICounty {
	cs := GetCities(province)
	for _, p := range cs {
		if p.GetName() == city {
			return p.GetCounties()
		}
	}
	return nil
}
