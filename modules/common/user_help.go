package common

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"stbweb/core"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/scrypt"
)

type user struct {
	Name       string
	Pwd        string
	Avatar     string
	Email      string
	Phone      string
	Salt       string
	UpdateTime string
}

func getUser(username string) *user {
	u := user{}
	if err := core.Ddb.QueryRow("SELECT name,password,avatar,email,phone,salt,updatetime FROM user where name=?", username).Scan(&u.Name, &u.Pwd, &u.Avatar, &u.Email, &u.Phone, &u.Salt, &u.UpdateTime); err != nil {
		core.LOG.WithFields(logrus.Fields{"get user": err}).Error("user")
	}
	return &u
}

//isExistUser 判断用户是否存在，存在为true
func isExistUser(username string) bool {
	num := 0
	if err := core.Ddb.QueryRow("SELECT count(*) FROM user where name=?", username).Scan(&num); err != nil {
		core.LOG.WithFields(logrus.Fields{"get user": err}).Error("user")
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
		core.LOG.WithFields(logrus.Fields{"decode": err}).Error("hex")
	}
	return decoded
}

//随机获取用户中一段+uuid生成随机盐，防止代码泄密密码生成过程被破解
func buildIserSalt(user string) string {
	rand.Seed(time.Now().UnixNano())
	sl := rand.Intn(len(user))
	return user[sl:] + base64.RawURLEncoding.EncodeToString(uuid.NewUUID())
}

//buildUserPassword 根据密码文本和盐生成密文
func buildUserPassword(pwdMd5, salt []byte) ([]byte, error) {
	return scrypt.Key(pwdMd5, salt, 16384, 8, 1, 32)
}

//md5加密
// func md5Encode(salt string) {
// 	data := []byte(salt)
// 	has := md5.Sum(data)
// 	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
// 	fmt.Println(md5str1)
// }

//密文验证
func (u *user) equal(pwd string) bool {
	bPwd := buildPas(pwd, u.Salt)
	return bytes.Equal(bPwd, []byte(u.Pwd))
}

//解析前端的hex密码文本，并调用密文生成函数
func buildPas(pwd, salt string) []byte {
	bPwd, err := buildUserPassword(huexEncode(pwd), []byte(salt))
	if err != nil {
		core.LOG.WithFields(logrus.Fields{"pwd": err}).Error("validPwdMd5")
	}
	return bPwd
}
