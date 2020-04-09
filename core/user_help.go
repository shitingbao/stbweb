package core

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/scrypt"
)

type User struct {
	Name       string
	Pwd        []byte //对应mysql的varbinary,末尾不会填充，不能使用binary，因为不足会使用ox0填充导致取出的时候多18位的0
	Avatar     string
	Email      string
	Phone      string
	Salt       string
	UpdateTime string
}

//GetUser 获取一个user
func GetUser(username string) *User {
	u := User{}
	if err := Ddb.QueryRow("SELECT name,password,avatar,email,phone,salt,updatetime FROM user where name=?", username).Scan(&u.Name, &u.Pwd, &u.Avatar, &u.Email, &u.Phone, &u.Salt, &u.UpdateTime); err != nil {
		LOG.WithFields(logrus.Fields{"get user": err}).Error("user")
	}
	return &u
}

//IsExistUser 判断用户是否存在，存在为true
func IsExistUser(username string) bool {
	num := 0
	if err := Ddb.QueryRow("SELECT count(*) FROM user where name=?", username).Scan(&num); err != nil {
		LOG.WithFields(logrus.Fields{"get user": err}).Error("user")
	}
	if num > 0 {
		return true
	}
	return false
}

//前端的hex字符串
func huexEncode(md5Pwd string) []byte {
	decoded, err := hex.DecodeString(md5Pwd)
	if err != nil {
		LOG.WithFields(logrus.Fields{"decode": err}).Error("hex")
	}
	return decoded
}

//BuildIserSalt 随机获取用户中一段+uuid生成随机盐，防止代码泄密密码生成过程被破解
func BuildIserSalt(user string) string {
	rand.Seed(time.Now().UnixNano())
	sl := rand.Intn(len(user))
	return user[sl:] + base64.RawURLEncoding.EncodeToString(uuid.NewUUID())
}

//buildUserPassword 根据密码文本和盐生成密文
func buildUserPassword(pwdMd5, salt []byte) ([]byte, error) {
	return scrypt.Key(pwdMd5, salt, 16384, 8, 1, 32)
}

//Equal 密文验证
func (u *User) Equal(pwd string) bool {
	bPwd := BuildPas(pwd, u.Salt)
	return bytes.Equal(bPwd, u.Pwd)
}

//BuildPas 解析前端的hex密码文本，并调用密文生成函数
func BuildPas(pwd, salt string) []byte {
	bPwd, err := buildUserPassword(huexEncode(pwd), []byte(salt))
	if err != nil {
		LOG.WithFields(logrus.Fields{"pwd": err}).Error("validPwdMd5")
	}
	return bPwd
}
