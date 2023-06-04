package mysql

import (
	"encoding/json"
	"fmt"

	communityv1 "ecodepost/pb/community/v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetSystemTheme(db *gorm.DB) (res SystemTheme, err error) {
	info, err := getSystemValueJson(db, systemKeyTheme)
	if err != nil {
		err = fmt.Errorf("GetSystemTheme fail, err: %w", err)
		return
	}
	err = json.Unmarshal([]byte(info.ValueJson), &res)
	if err != nil {
		err = fmt.Errorf("GetSystemTheme fail2, err: %w", err)
		return
	}
	return
}

func PutSystemTheme(db *gorm.DB, req *communityv1.SetThemeReq) (err error) {
	var res SystemTheme
	info, err := getSystemValueJson(db.Clauses(clause.Locking{Strength: "UPDATE"}), systemKeyTheme)
	if err != nil {
		err = fmt.Errorf("PutSystemTheme fail, err: %w", err)
		return
	}
	err = json.Unmarshal([]byte(info.ValueJson), &res)
	if err != nil {
		err = fmt.Errorf("PutSystemTheme fail2, err: %w", err)
		return
	}
	res.ThemeName = req.ThemeName
	res.DefaultAppearance = req.DefaultAppearance
	res.IsCustom = req.IsCustom

	customColor := CustomColor{}
	err = json.Unmarshal([]byte(req.GetCustomColor()), &customColor)
	if err != nil {
		err = fmt.Errorf("PutSystemTheme fail3, err: %w", err)
		return
	}
	res.CustomColor = customColor
	jsonByte, err := json.Marshal(res)
	if err != nil {
		err = fmt.Errorf("PutSystemHomeOption json marshal fail, err: %w", err)
		return
	}
	err = putSystemValueJson(db, systemKeyTheme, string(jsonByte))
	if err != nil {
		err = fmt.Errorf("PutSystemHomeOption putSystemValueJson fail, err: %w", err)
		return
	}
	return
}
