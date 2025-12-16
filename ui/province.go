package ui

import (
	"fmt"
	"strings"
	"tb_helper/constant"
)

// 根据省份名获取省份
func GetProvinceByName(name string) *constant.Province {
	shortName := name
	if len([]rune(name)) >= 2 {
		shortName = string([]rune(name)[:2])
	}
	for _, p := range constant.ProvinceAndCityArr {
		if strings.Contains(p.Name, shortName) {
			return &p
		}
	}
	return nil
}

// 根据省份名获取城市列表（可指定默认选中的城市名）
func GetCities(provinceName, cityName string) []string {
	p := GetProvinceByName(provinceName)
	if p == nil {
		return []string{}
	}

	var cities []string
	cities = append(cities, "请选择")
	cities = append(cities, p.Cities...)

	// 查找默认城市
	if cityName != "" {
		var selected string
		shortCity := cityName
		if len([]rune(cityName)) >= 2 {
			shortCity = string([]rune(cityName)[:2])
		}
		for _, c := range p.Cities {
			if c == cityName { // 完全匹配
				selected = c
				break
			} else if selected == "" && strings.Contains(c, shortCity) {
				selected = c
			}
		}
		if selected != "" {
			fmt.Printf("默认选中城市: %s\n", selected)
		}
	}
	return cities
}
