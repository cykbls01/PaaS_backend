package util

import (
	"Cloud/pkg/util/ws"
	"encoding/json"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
	"unicode/utf8"
)

// ssh流式处理器
type StreamHandler struct {
	WsConn *ws.WsConnection
	ResizeEvent chan remotecommand.TerminalSize
}

// web终端发来的包
type XtermMessage struct {
	MsgType string `json:"type"`	// 类型:resize客户端调整终端, input客户端输入
	Input string `json:"input"`	// msgtype=input情况下使用
	Rows uint16 `json:"rows"`	// msgtype=resize情况下使用
	Cols uint16 `json:"cols"`// msgtype=resize情况下使用
}

// executor回调获取web是否resize
func (handler *StreamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <- handler.ResizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *StreamHandler) Read(p []byte) (size int, err error) {
	var (
		msg *ws.WsMessage
		xtermMsg XtermMessage
	)

	// 读web发来的输入
	if msg, err = handler.WsConn.WsRead(); err != nil {
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return
	}

	//web ssh调整了终端大小
	if xtermMsg.MsgType == "resize" {
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.ResizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" {	// web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

// executor回调向web端输出
func (handler *StreamHandler) Write(p []byte) (size int, err error) {
	var (
		copyData []byte
	)

	// 产生副本
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	if !utf8.Valid(copyData){
		return
	}
	err = handler.WsConn.WsWrite(websocket.TextMessage, copyData)
	return
}
