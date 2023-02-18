package database

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	path = "./database/db.json"
)

type User struct {
	ChatId   int64    `json:"id"`
	Time     string   `json:"time"`
	Crypto   []string `json:"crypto"`
	Currency string   `json:"currency"`
	Ok       bool     `json:"ok"`
	IsActive bool     `json:"isActive"`
}

func CreateUser(chatid int64) {
	if isInDatabase(chatid) {
		return
	}
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)
	u := User{chatid, "", nil, "🇺🇸", false, false}
	users = append(users, u)
	data, _ := json.MarshalIndent(users, "", " ")
	os.WriteFile(path, data, 0755)
}

func SetActive(chatid int64, flag bool) error {
	db, err := os.ReadFile(path)
	if err != nil {
		return errors.New("file read error")
	}
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			users[i].IsActive = flag
			break
		}
	}
	data, _ := json.MarshalIndent(users, "", " ")
	os.WriteFile(path, data, 0755)
	return nil
}

func GetActive(chatid int64) bool {
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			return users[i].IsActive
		}
	}
	return false
}

func SetTime(chatid int64, tm string) error {
	if !isRigtTimeFormat(tm) {
		return errors.New("Wrong time format")
	}
	db, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			users[i].Time = tm
			break
		}
	}

	data, _ := json.MarshalIndent(users, "", " ")
	os.WriteFile(path, data, 0755)
	return nil
}

func AddCrypto(chatid int64, cryptocurrency string) error {
	cryptocurrency = strings.Replace(cryptocurrency, "https://coinmarketcap.com/currencies/", "", 1)
	cryptocurrency = strings.Trim(cryptocurrency, "/")
	if !cryptoValidate(cryptocurrency) {
		return errors.New("Wrong cryptocurrency name")
	}
	db, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			users[i].Crypto = append(users[i].Crypto, cryptocurrency)
			break
		}
	}

	data, _ := json.MarshalIndent(users, "", " ")
	os.WriteFile(path, data, 0755)

	return nil
}

func SetCurrency(chatid int64, currency string) error {
	if !currencyValidate(currency) {
		return errors.New("WRONG CURRENCY")
	}
	db, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			users[i].Currency = currency
			break
		}
	}

	data, _ := json.MarshalIndent(users, "", " ")
	os.WriteFile(path, data, 0755)
	return nil
}

func SetOk(chatid int64, ok bool) {
	db, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			users[i].Ok = ok
			break
		}
	}

	data, _ := json.MarshalIndent(users, "", " ")
	os.WriteFile(path, data, 0755)
}

func GetOk(chatid int64) bool {
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			return users[i].Ok
		}
	}
	return false
}

func GetCrypto(chatid int64) []string {
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			return users[i].Crypto
		}
	}
	return nil
}

func GetCurrency(chatid int64) string {
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			return users[i].Currency
		}
	}
	return ""
}

func GetTime(chatid int64) string {
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)

	for i := 0; i < len(users); i++ {
		if users[i].ChatId == chatid {
			return users[i].Time
		}
	}
	return ""
}

func isInDatabase(chatid int64) bool {
	db, _ := os.ReadFile(path)
	var users []User
	json.Unmarshal(db, &users)
	for _, user := range users {
		if user.ChatId == chatid {
			return true
		}
	}
	return false
}

func isRigtTimeFormat(tm string) bool {
	arr := strings.Split(tm, ":")
	if len(arr) != 2 {
		return false
	}
	hours, err := strconv.Atoi(arr[0])
	if err != nil || hours < 0 || hours > 23 {
		return false
	}
	minutes, err := strconv.Atoi(arr[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return false
	}
	return true
}

func cryptoValidate(crypto string) bool {
	crypto = "https://coinmarketcap.com/currencies/" + crypto
	if resp, err := http.Get(crypto); err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func currencyValidate(cur string) bool {
	if cur != "🇷🇺" && cur != "🇺🇸" {
		return false
	}
	return true
}
