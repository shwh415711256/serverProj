package main

import (
	"Common/model"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
	"golang.org/x/net/websocket"
	"RoomServer/src/server/msg"
	"strconv"
	"sync"
)

func TestCreateRoom(gameid string, openid string) uint64{
	url := "https://cc.h5youxi.fun/createRoom"
	data := &model.CreateRoomRequest{
		GameId:gameid,
		OpenId:openid,
		RoomType:1,
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return 0
	}
	body := bytes.NewBuffer(bts)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-type", "application/json")
	response, _ := client.Do(req)
	defer response.Body.Close()
	if response == nil {
		return 0
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("TestCreateRoom error, response:%v, gameid:%s, openid:%s", response, gameid, openid)
		return 0
	}
	resdata, err := ioutil.ReadAll(response.Body)
//	fmt.Println("CreateRoomResp:%s, gameid:%s, openid:%s", string(resdata), gameid, openid)
	resp := &model.CreateRoomResponse{}
	err = json.Unmarshal(resdata, resp)
	if err != nil {
		return 0
	}
	return resp.RoomId
}

func TestConnectWss() *websocket.Conn {
	origin := "wss://cc.h5youxi.fun"
	url := "wss://cc.h5youxi.fun/wss"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println("TestConnectWss error, can not connet websocket,err:", err)
		return nil
	}
	return ws
}

func TestLogin(conn *websocket.Conn, gameid string, openid string) bool{
	req := &msg.LoginReq{
		GameId:gameid,
		OpenId:openid,
	}
	v := map[string]msg.LoginReq{
		"LoginReq": *req,
	}
	sendbts, err := json.Marshal(v)
	if err != nil {
		return false
	}
	_, err = conn.Write(sendbts)
	if err != nil {
		return false
	}
//	fmt.Println("send LoginReq success, data:", string(reqbts))
	/*
	respdata := make([]byte,512)
	n, err := conn.Read(respdata)
	if err != nil {
		return false
	}
	var m map[string]json.RawMessage
	err = json.Unmarshal(respdata[:n], &m)
	if err != nil {
		fmt.Println("login   ",err)
		return false
	}
	resp := &msg.LoginResp{}
	err = json.Unmarshal(m["LoginResp"], resp)
	if err != nil {
		fmt.Println("login   ",err)
		return false
	}
	if resp.Result != "ok" {
		fmt.Errorf("Login error, gameid:%s, openid:%s", gameid, openid)
		return false
	}
	fmt.Println("Login resp:,", string(respdata))
	*/
	return true
}

func TestEnterRoom(conn *websocket.Conn, gameid string, openid string, roomid uint64) bool{
	enter := &msg.EnterRoomRequest{
		GameId:gameid,
		RoomId:roomid,
		OpenId:openid,
	}
	v := map[string]msg.EnterRoomRequest{
		"EnterRoomRequest": *enter,
	}
	reqbts, err := json.Marshal(v)
	if err != nil {
		return false
	}
	_, err = conn.Write(reqbts)
	if err != nil {
		return false
	}
//	fmt.Println("send EnterRoomRequest success, data:", string(reqbts))
	/*
	respdata := make([]byte, 4096)
	n, err := conn.Read(respdata)
	if err != nil {
		fmt.Println("enterRoom   ",err)
		return false
	}
	var m map[string]json.RawMessage
	err = json.Unmarshal(respdata[:n], &m)
	if err != nil {
		fmt.Println("enterRoom   ",err)
		return false
	}
	resp := &msg.EnterRoomResponse{}
	err = json.Unmarshal(m["EnterRoomResponse"], resp)
	if err != nil {
		fmt.Println("enterRoom   ",err)
		return false
	}
	fmt.Println("EnterRoomRequest resp:", string(respdata))
	*/
	return true
}

func TestStartGame(conn *websocket.Conn, roomid uint64, openid string) bool{
	start := &msg.StartGameRequest{
		RoomId:roomid,
		OpenId:openid,
	}
	v := map[string]msg.StartGameRequest{
		"StartGameRequest": *start,
	}
	reqbts, err := json.Marshal(v)
	if err != nil {
		fmt.Println("1111111111  ",err)
		return false
	}
	_, err = conn.Write(reqbts)
	if err != nil {
		fmt.Println("2222222222  ",err)
		return false
	}
//	fmt.Println("send StartGameRequest success, data:", string(reqbts))
	/*
	respdata := make([]byte, 4096)
	n, err := conn.Read(respdata)
	if err != nil {
		return false
	}
	var m map[string]json.RawMessage
	err = json.Unmarshal(respdata[:n], &m)
	if err != nil {
		return false
	}
	resp := &msg.StartGameResponse{}
	err = json.Unmarshal(m["StartGameResponse"], resp)
	if err != nil {
		fmt.Println("start   ",err, "   ", string(respdata))
		return false
	}
	fmt.Println("StartGameRequest resp:", string(respdata))
	return resp.Result == 1
	*/
	return true
}

func MsgParse(conn *websocket.Conn, openid string, gameid string, roomid uint64, data []byte,index int, count *int) bool{
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("msgparse   ",err, "1111   ", string(data))
		return false
	}
	for k, v := range m {
		switch k{
		case "LoginResp":
			{
				resp := &msg.LoginResp{}
				err = json.Unmarshal(v, resp)
				if err != nil {
					fmt.Println("login   ", err)
					return false
				}
				if resp.Result != "ok" {
					fmt.Errorf("Login error, gameid:%s, openid:%s", gameid, openid)
					return false
				}
		//		fmt.Println("login resp ok,openid:", openid)
				if (index%2 == 0 && !resp.IsReConnect) {
					TestEnterRoom(conn, gameid, openid, roomid)
				}
			}
		case "EnterRoomResponse":
			{
				resp := &msg.EnterRoomResponse{}
				err = json.Unmarshal(m["EnterRoomResponse"], resp)
				if err != nil {
					fmt.Println("enterRoom   ",err)
					return false
				}
				if resp.Result != 1 {
					fmt.Println("enterroom error, result:", resp.Result)
					return false
				}
	//			fmt.Println("enterroom ok")
			}
		case "UpdateRoomUserData":
			{
				resp := &msg.UpdateRoomUserData{}
				err = json.Unmarshal(m["UpdateRoomUserData"], resp)
				if err != nil {
					fmt.Println("UpdateRoom   ",err)
					return false
				}
				// 等待开始
				if resp.RoomState != 2{
					fmt.Println("UpdateRoom error")
					return false
				}
				TestStartGame(conn, roomid, openid)
			}
		case "StartGameResponse":
			{
				resp := &msg.StartGameResponse{}
				err = json.Unmarshal(m["StartGameResponse"], resp)
				if err != nil {
					fmt.Println("StartGame   ",err)
					return false
				}
				// 等待开始
				if resp.Result != 1 {
					fmt.Println("StartGame error")
					return false
				}
				*count ++
		//		fmt.Println("StartGameResponse ok, openid:%s", openid)
			}
		case "GameResultMsg":
			{
				resp := &msg.GameResultMsg{}
				err = json.Unmarshal(m["GameResultMsg"], resp)
				if err != nil {
					fmt.Println("GameResul   ",err)
					return false
				}
	//			fmt.Println("openid:",openid, " index:",index, " count:",*count)
				if *count >= 10 {
					return false
				}
				if index % 2 != 0 {
					TestStartGame(conn, roomid, openid)
				}
		//		fmt.Println("GameResultMsg ok")
			}
		default:
			{
				fmt.Println(k, "          ", v)
			}
		}

	}
	return true
}

func Run(conn *websocket.Conn, openid string, gameid string, roomid uint64, index int) {
	TestLogin(conn, gameid, openid)
	playcount := 0
	for {
		data := make([]byte, 4096)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("conn read error,", err)
			break
		}
		if !MsgParse(conn, openid, gameid, roomid, data[:n], index, &playcount) {
			break
		}
	}
	conn.Close()
	wgcount++
	fmt.Println(openid,"    " ,wgcount)
	wg.Done()
}

var wg sync.WaitGroup
var wgcount int

func main() {
	wg = sync.WaitGroup{}
	wgcount = 0
	gameid := "5"
	for i := 1; i < 5000; i += 2 {
		openid1 := "testOpenid_" + strconv.Itoa(i)
		openid2 := "testOpenid_" + strconv.Itoa(i+1)
		//		time.Sleep(500 * time.Millisecond)
		roomid := TestCreateRoom(gameid, openid1)
		conn1 := TestConnectWss()
		conn2 := TestConnectWss()
		wg.Add(2)
		if conn1 != nil && conn2 != nil && roomid != 0 {
			go Run(conn1, openid1, gameid, roomid, i)
			go Run(conn2, openid2, gameid, roomid, i+1)
		}
	}
	wg.Wait()

	/*
	wg := sync.WaitGroup{}
	for i := 1; i < 20; i += 2 {
		openid1 := "testOpenid" + strconv.Itoa(i)
		openid2 := "testOpenid" + strconv.Itoa(i+1)
		wg.Add(1)
		time.Sleep(500 * time.Millisecond)
		roomid := TestCreateRoom(gameid, openid1)
		conn1 := TestConnectWss()
		conn2 := TestConnectWss()
		if conn1 != nil && conn2 != nil && roomid != 0 {
			if !TestLogin(conn1, gameid, openid1) {
				fmt.Println("111111111111")
				wg.Done()
				return
			}
			if !TestLogin(conn2, gameid, openid2) {
				wg.Done()
				return
			}
			if !TestEnterRoom(conn2, gameid, openid2, roomid) {
				wg.Done()
				return
			}
			datas := make([]byte, 4096)
			_, _ = conn1.Read(datas)
			fmt.Println(string(datas))
			if !TestStartGame(conn1, roomid, openid1) {
				wg.Done()
				return
			}
			_, _ = conn2.Read(datas)
			fmt.Println("1111111111111111111111  ", string(datas))
			for {
				TestGameOver(conn1)
				TestGameOver(conn2)
				if !TestStartGame(conn1, roomid, openid1) {
					break
				}
				_, _ = conn2.Read(datas)
				fmt.Println("222222222222222   ", string(datas))
				TestGameOver(conn1)
				TestGameOver(conn2)
				conn1.Close()
				conn2.Close()
				break
			}
		}
	}*/
}
