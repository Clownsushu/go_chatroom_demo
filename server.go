package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type UserStruct struct {
	Nickname string
	Room     string
}

// 连接池
var Conntions []*websocket.Conn

var Users = make(UserType)

type UserType map[net.Addr]UserStruct

var Rooms = make(RoomType)

type RoomType map[string][]*websocket.Conn

// ws地址
var WsPort string = ":8888"

var Up = websocket.Upgrader{
	ReadBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(WsPort, nil)
	if err != nil {
		fmt.Println("Err : ", err)
		return
	}
}

func Handler(resp http.ResponseWriter, req *http.Request) {
	//创建链接
	conn, err := Up.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Println("链接失败: ", err)
		return
	}
	defer conn.Close()
	//追加连接池
	Conntions = append(Conntions, conn)

	go readMessage(conn)
	for {

	}
}

func readMessage(conn *websocket.Conn) {
	for {
		t, data, e := conn.ReadMessage()

		if e != nil {
			delData(conn)
			fmt.Println("errr: ", e, conn.RemoteAddr())
			break
		}
		fmt.Println("接收到: ", t, string(data))
		//判断是否存在指定字符&^|0! 存在就存入昵称等信息
		f := strings.Split(string(data), "&^|0!")
		if len(f) == 2 {
			//设置房间
			setRooms(conn, f[0], f[1])
			//设置用户
			setUser(conn, f[0], f[1])
		} else {
			//广播消息
			broadcastMessage(conn, data, 1)
		}

	}
}

func broadcastMessage(conn *websocket.Conn, msg []byte, t int) {
	room := Users[conn.RemoteAddr()].Room
	message := ""
	switch t {
	case 1:
		message = Users[conn.RemoteAddr()].Nickname + "：" + string(msg)
	default:
		message = string(msg)
	}
	if len(Rooms[room]) > 0 {
		for _, v := range Rooms[room] {
			if v.RemoteAddr() != conn.RemoteAddr() {
				v.WriteMessage(websocket.TextMessage, []byte(message))
			}
		}
	}
}

func setRooms(conn *websocket.Conn, nickname, room string) {
	Rooms[room] = append(Rooms[room], conn)
	msg := "欢迎加入'" + room + "'聊天室, 当前聊天室有" + strconv.Itoa(len(Rooms[room])) + "人在线"
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
	//循环给其他人发生消息
	welcome := "欢迎'" + nickname + "'加入聊天室" + ", 当前聊天室有" + strconv.Itoa(len(Rooms[room])) + "人在线!"
	for _, v := range Rooms[room] {
		if v.RemoteAddr() != conn.RemoteAddr() {
			v.WriteMessage(websocket.TextMessage, []byte(welcome))
		}
	}
}

func setUser(conn *websocket.Conn, nickname, room string) {
	user := make(map[string]string)
	user["nickname"] = nickname
	user["room"] = room
	//type UserType map[net.Addr]UserStruct
	Users[conn.RemoteAddr()] = UserStruct{nickname, room}
	fmt.Println(Users)
}

func delData(conn *websocket.Conn) {
	room := Users[conn.RemoteAddr()].Room
	nickname := Users[conn.RemoteAddr()].Nickname

	//发送退出消息
	msg := "'" + nickname + "'退出了聊天室!"
	broadcastMessage(conn, []byte(msg), 2)

	delete(Users, conn.RemoteAddr())

	roomSli := make([]*websocket.Conn, 0)

	for _, v := range Rooms[room] {
		if v.RemoteAddr() != conn.RemoteAddr() {
			roomSli = append(roomSli, v)
			break
		}
	}
	Rooms[room] = roomSli
	//fmt.Println("rooms: ", Rooms[room])
}
