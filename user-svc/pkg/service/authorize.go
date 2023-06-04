package service

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	mt "math/rand"
	"strconv"
	"strings"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
)

type authorize struct {
	secret             string
	saltSize           int
	delimiter          string
	stretchingPassword int
}

func InitAuthorize() *authorize {
	obj := &authorize{
		secret:             econf.GetString("user-svc.oauth.secret"),
		saltSize:           econf.GetInt("user-svc.oauth.salt"),
		delimiter:          econf.GetString("user-svc.oauth.delimiter"),
		stretchingPassword: econf.GetInt("user-svc.oauth.stretchingPassword"),
	}
	return obj
}

func (p *authorize) Hash(pass string) (string, error) {
	saltSecretInfo, err := saltSecret()
	if err != nil {
		return "", err
	}

	salt, err := p.salt(p.secret + saltSecretInfo)
	if err != nil {
		return "", err
	}

	iteration := randInt(1, 3)
	hash, err := p.hash(pass, saltSecretInfo, salt, int64(iteration))
	if err != nil {
		return "", err
	}
	iterationString := strconv.Itoa(iteration)
	password := p.hashJoin(saltSecretInfo, iterationString, hash, salt)
	elog.Debug("gen hash data", zap.String("saltSecret", saltSecretInfo), zap.String("iterationString", iterationString), zap.String("hash", hash), zap.String("salt", salt))
	return password, nil

}

// "saltSecret": "CklLQbaItQSHjuV-0278p2WoODqtB635IE4FPWx5nAi9RJ2kv5D5fFbK5qNL3VFkVrpjCg3JS4v7"
//  "iterationString": "5"
// hash 54ac3643d100c56c7efe9af2295c3371a1c4681feeca63399670798d
// salt 54e6909e6758cc6ad41f791da5e57e380b0fbde395b83de1

// CklLQbaItQSHjuV-0278p2WoODqtB635IE4FPWx5nAi9RJ2kv5D5fFbK5qNL3VFkVrpjCg3JS4v7$5$54ac3643d100c56c7efe9af2295c3371a1c4681feeca63399670798d$54e6909e6758cc6ad41f791da5e57e380b0fbde395b83de1
// CklLQbaItQSHjuV-0278p2WoODqtB635IE4FPWx5nAi9RJ2kv5D5fFbK5qNL3VFkVrpjCg3JS4v7$5$c40691fd4c4d4d436a664ce4112e2c2a4792fa700b9c84d9d8f538ac$54e6909e6758cc6ad41f791da5e57e380b0fbde395b83de1

// 校验密码是否有效
func (p *authorize) Verify(uid int64, hashing string, pass string) error {
	data := p.trimSaltHash(hashing)
	iteration, _ := strconv.ParseInt(data["iteration_string"], 10, 64)

	has, err := p.hash(pass, data["salt_secret"], data["salt"], iteration)
	if err != nil {
		return err
	}

	hashJoin := p.hashJoin(data["salt_secret"], data["iteration_string"], has, data["salt"])
	if hashJoin == hashing {
		return nil
	}
	return errors.New("not equal")
}

func (p *authorize) hash(pass string, salt_secret string, salt string, iteration int64) (string, error) {
	var pass_salt = salt_secret + pass + salt + salt_secret + pass + salt + pass + pass + salt
	var i int

	hash_pass := p.secret
	hash_start := sha512.New()
	hash_center := sha256.New()
	hash_output := sha256.New224()

	i = 0
	for i <= p.stretchingPassword {
		i = i + 1
		hash_start.Write([]byte(pass_salt + hash_pass))
		hash_pass = hex.EncodeToString(hash_start.Sum(nil))
	}

	i = 0
	for int64(i) <= iteration {
		i = i + 1
		hash_pass = hash_pass + hash_pass
	}

	i = 0
	for i <= p.stretchingPassword {
		i = i + 1
		hash_center.Write([]byte(hash_pass + salt_secret))
		hash_pass = hex.EncodeToString(hash_center.Sum(nil))
	}
	hash_output.Write([]byte(hash_pass + p.secret))
	hash_pass = hex.EncodeToString(hash_output.Sum(nil))

	return hash_pass, nil
}

func (p *authorize) trimSaltHash(hash string) map[string]string {
	str := strings.Split(hash, p.delimiter)

	return map[string]string{
		"salt_secret":      str[0],
		"iteration_string": str[1],
		"hash":             str[2],
		"salt":             str[3],
	}
}

func (p *authorize) hashJoin(saltSecret, iteration, hash, salt string) string {
	return strings.Join([]string{saltSecret, iteration, hash, salt}, p.delimiter)
}

func (p *authorize) salt(secret string) (string, error) {
	buf := make([]byte, p.saltSize, p.saltSize+md5.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(buf)
	hash.Write([]byte(secret))
	return hex.EncodeToString(hash.Sum(buf)), nil
}

func saltSecret() (string, error) {
	rb := make([]byte, randInt(10, 100))
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

func randInt(min int, max int) int {
	return min + mt.Intn(max-min)
}
