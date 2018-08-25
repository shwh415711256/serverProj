package cache

import (
	"DBServer/db/redisdb"
	"github.com/garyburd/redigo/redis"
	"github.com/cihub/seelog"
	"Common/model"
	"DBServer/module/account"
	"DBServer/module/sign"
	"strconv"
	"github.com/coreos/etcd/clientv3"
)

const (
	AccountExprieTime = 3600 * 24
)

var (
	wxAccountValueStr []string = []string {
		"openid", "nickname", "avatarurl", "gendar", "city",
		"province", "country", "language",
	}
)

func GetAccountInfo(gameId string, openId string) (*model.AccountInfo, bool){
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	accountInfo := &model.AccountInfo{}
	gold, err := redis.Uint64(pool.Do("hget","accountInfo/" + gameId + openId, "gold"))
	if err != nil {
		if err == redis.ErrNil {
			accountInfo, err = account.LoadAccountInfoFromDB(gameId, openId)
			if err != nil {
				seelog.Errorf("getAccountInfo From DB error[%v]", err)
				return nil, false
			}
			SetAccountInfo(gameId, accountInfo)
			return accountInfo, true
		}
		seelog.Errorf("redis getAccountInfo error[%v]", err)
		return nil, false
	}
	money, _ := redis.Uint64(pool.Do("hget","accountInfo/" + gameId + openId, "money"))
	score, _ := redis.Uint64(pool.Do("hget","accountInfo/" + gameId + openId, "score"))
	accountInfo.WxOpenId = openId
	accountInfo.Gold = gold
	accountInfo.Money = money
	accountInfo.Score =  score
	pool.Do("EXPIRE", "accountInfo/" + gameId + openId, AccountExprieTime)
	seelog.Debugf("%v", accountInfo)
	return accountInfo, true
}

func SetAccountInfo(gameId string, info *model.AccountInfo) {
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	pool.Do("hmset", "accountInfo/" + gameId + info.WxOpenId, "gold", info.Gold,
		"money", info.Money, "score", info.Score)
	pool.Do("EXPIRE", "accountInfo/" + gameId + info.WxOpenId, AccountExprieTime)
}

func UpdateAccountExprie(gameId string, openId string) {
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	pool.Do("EXPIRE", "accountInfo/" + gameId + openId, AccountExprieTime)
}

func DelAccountInfo(gameId string, openId string) {
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	pool.Do("DEL", "accountInfo/" + gameId + openId)
}

func GetSignInfo(gameId string, openId string, signType int, data *model.RpcGetSignInfoResponse) bool {
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	values, err := redis.Values(pool.Do("hmget", "signInfo/"+gameId+strconv.Itoa(signType)+openId, "sign_sumday", "today_sign"))
	if err != nil {
		if err == redis.ErrNil {
			signInfo, err := sign.LoadSignInfoFromDB(gameId, openId, signType)
			if err != nil {
				seelog.Errorf("getSignInfo From DB error[%v]", err)
				return false
			}
			pool.Do("hmset", "signInfo/"+gameId+strconv.Itoa(signType)+openId, "sign_sumday", signInfo.SignSumDay,
				"today_sign", signInfo.IsTodaySign)
			data.SignSumDay = signInfo.SignSumDay
			data.IsTodaySign = signInfo.IsTodaySign
			return true
		}
	}
	signSumDay, err := strconv.ParseInt(string(values[0].([]byte)), 10, 32)
	if err != nil {
		seelog.Errorf("getSignInfo parse signSumDay error[%v]", err)
		return false
	}
	data.SignSumDay = int(signSumDay)
	if string(values[1].([]byte)) == "1" {
		data.IsTodaySign = true
	}else {
		data.IsTodaySign = false
	}
	return true
}

/*
func GetWxAccountInfo(openId string) (map[string]string, bool){
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	doStr := fmt.Sprintf("hmget wxAccountInfo/%s ", openId)
	l := len(wxAccountValueStr)
	for i := 0; i < l; i ++ {
		doStr += wxAccountValueStr[i] + " "
	}
	info, err := redis.StringMap(pool.Do(doStr))
	if err != nil && err == redis.ErrNil{
		seelog.Debugf("redis getAccountInfo error[%v]", err)
		return nil, false
	}
	return info, true
}

func SetWxAccountInfo(info *model.WxUserInfo) {
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	doStr := fmt.Sprintf("hmset wxAccountInfo/%s ", info.OpenId)
	l := len(wxAccountValueStr)
	for i := 0; i < l; i ++ {
		doStr += wxAccountValueStr[i] + " "
	}
	pool.Do(doStr, info.NickName, info.AvatarUrl, info.Gendar, info.City, info.Province, info.Country, info.Lang)
	pool.Do("EXPIRE wxAccountInfo/%s %d", info.OpenId, AccountExprieTime)
}

func UpdateWxAccountExprie(openId string) {
	pool := redisdb.RedisPool.Get()
	defer pool.Close()
	pool.Do("EXPIRE wxAccountInfo/%s %d", openId, AccountExprieTime)
}*/