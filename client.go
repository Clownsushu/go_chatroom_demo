package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"strconv"
)

var wsUrl string = "ws://localhost:8888"

var rloadMsg string = "\n请重新打开客户端"

func main() {
	var nickname string
	var chatRoomNum int

	for {
		fmt.Print("请输入您的昵称, 按回车确认：")
		nreader := bufio.NewReader(os.Stdin)
		s, _, _ := nreader.ReadLine()
		nickname = string(s)
		if nickname == "" {
			fmt.Println("昵称不能是空白!!!")
		} else {
			break
		}
	}

	for {
		fmt.Print("请输入聊天室编号, 按回车确认：")
		reader := bufio.NewReader(os.Stdin)
		b, _, _ := reader.ReadLine()
		chatRoomNum, _ = strconv.Atoi(string(b))
		if chatRoomNum > 0 {
			break
		} else {
			fmt.Println("请输入聊天室编号, 编号是数字!!!")
		}
	}
	//开始socket
	dl := websocket.Dialer{}
	conn, _, err := dl.Dial(wsUrl, nil)
	if err != nil {
		fmt.Println("链接websocket失败: ", err, rloadMsg)
		return
	}
	//链接上之后发送昵称和聊天室编号
	data := nickname + "&^|0!" + strconv.Itoa(chatRoomNum)
	conn.WriteMessage(websocket.TextMessage, []byte(data))
	fmt.Println("输入你要发送的消息, 按回车键确认!!!")
	//关闭客户端
	defer conn.Close()
	//创建协程监听消息
	go sendMessage(conn)
	//创建协程读消息
	go getMessage(conn)
	for {
	}
}

func getMessage(conn *websocket.Conn) {
	for {
		_, d, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("接收消息失败: ", err, rloadMsg)
			return
		} else {
			fmt.Println(string(d))
			//fmt.Print("请输入您要发送的消息, 按回车键确认：")
		}
	}
}

func sendMessage(conn *websocket.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		b, _, _ := reader.ReadLine()
		err := conn.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			fmt.Println("消息发送失败: ", err, rloadMsg)
			return
		}
	}
}
