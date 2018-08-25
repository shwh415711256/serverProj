package model

type OneRoomInfo struct {
	RoomType int
	RoomNeedPlayerNum int
	RoomMaxNum int
	RoomWaitToCloseTime int64
	RoomPlayGameTime int64
	NeedPlayerReady int
}

type CreateRoomRequest struct {
	GameId string
	OpenId string
	RoomType int
}

type CreateRoomResponse struct {
	Result int
	RoomServerAddr string
	RoomId uint64
}

type RandMatchEnemy struct{
	Score int
}

