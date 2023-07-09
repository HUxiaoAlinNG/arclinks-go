package util

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"arclinks-go/app/vars"
)

func ContainsInt64(arr []int64, i int64) bool {
	for _, v := range arr {
		if v == i {
			return true
		}
	}
	return false
}

func ContainsUInt8(arr []uint8, i uint8) bool {
	for _, v := range arr {
		if v == i {
			return true
		}
	}
	return false
}

func HashSimpleMd5(str string) string {
	hash := md5.New()
	_, err := io.WriteString(hash, str)
	if err != nil {
		return ""
	}
	return string(hash.Sum(nil))
}

//  ParseTime 解析 202005211717 格式时间
func ParseTime(timeStr string) (t time.Time, err error) {
	if len(timeStr) != 12 {
		return t, errors.New("不支持的格式")
	}
	year := timeStr[0:4]
	month := timeStr[4:6]
	day := timeStr[6:8]
	hour := timeStr[8:10]
	minute := timeStr[10:12]

	t, err = time.Parse(time.RFC3339, fmt.Sprintf("%s-%s-%sT%s:%s:00+08:00", year, month, day, hour, minute))
	if err != nil {
		return t, errors.New("解析时间出错")
	}
	return t, nil
}

// GetTimeString 传入：2021-08-19 14:57:31.163988436 +0800 CST m=+11.572627955 , 输出 021-08-19 14:57:31
func GetTimeString(t time.Time) string {
	var ServerTimezone = time.FixedZone("CST", 8*3600) // 东八
	return t.In(ServerTimezone).Format("2006-01-02 15:04:05")
}

func GetVhost() string {
	vhost := vars.RabbitMQSetting.Vhost

	return vhost
}

func Int64SliceToString(is []int64, sep string) string {
	slice := make([]string, len(is))
	for i, v := range is {
		slice[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(slice, sep)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

//RandStringRunes 产生随机字符
func RandStringRunes(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//StructHashCode 获取struct值的hash
func StructHashCode(structObj interface{}) (string, error) {
	bytes, err := json.Marshal(structObj)
	if err != nil {
		return "", nil
	}
	md5Bytes := md5.Sum(bytes)
	return fmt.Sprintf("%x", md5Bytes), nil
}
