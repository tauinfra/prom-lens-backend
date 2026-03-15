package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"valyria-backend/internal/core/logger"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

// StreamHandler SSH 流式处理器
type StreamHandler struct {
	WsConn      *WsConnection
	ResizeEvent chan remotecommand.TerminalSize
}

// xtermMessage ssh 端发的数据
type xtermMessage struct {
	Type  string `json:"type"`  // type=resize 调整终端, type=stdin 终端输入
	Stdin string `json:"stdin"` // type=stdin 情况下使用
	Rows  uint16 `json:"rows"`  // type=resize 情况下使用
	Cols  uint16 `json:"cols"`  // type=resize 情况下使用
}

// Next executor 回调获取 Xterm WebSSH 端大小
func (s *StreamHandler) Next() *remotecommand.TerminalSize {
	size, ok := <-s.ResizeEvent
	if !ok {
		return nil
	}
	return &size
}

// Read executor 回调读取 Xterm WebSSH 端输入
func (s *StreamHandler) Read(p []byte) (size int, err error) {
	var (
		msg      *WsMessage
		xtermMsg xtermMessage
	)

	// 使用非阻塞读取或带超时的读取
	msg, err = s.WsConn.WsRead()
	if err != nil {
		if err.Error() == "websocket closed" {
			return 0, io.EOF // 连接关闭返回 EOF
		}
		logger.Error(fmt.Sprintf("WebSocket读取失败: %v", err))
		return 0, err
	}

	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		logger.Error(fmt.Sprintf("JSON解析失败: %v, 原始数据: %s", err, string(msg.Data)))
		return 0, err
	}

	logger.Debug(fmt.Sprintf("收到消息类型: %s, 数据长度: %d", xtermMsg.Type, len(xtermMsg.Stdin)))

	switch xtermMsg.Type {
	case "resize":
		if xtermMsg.Cols == 0 || xtermMsg.Rows == 0 {
			logger.Warn("终端调整大小参数无效，已忽略")
			return 0, nil
		}
		logger.Info(fmt.Sprintf("终端调整大小: %dx%d", xtermMsg.Cols, xtermMsg.Rows))
		select {
		case s.ResizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}:
			// 成功发送调整事件
		default:
			logger.Warn("调整大小通道已满，丢弃事件")
		}
		return 0, nil

	case "stdin":
		if xtermMsg.Stdin == "" {
			return 0, nil
		}

		data := []byte(xtermMsg.Stdin)
		size = len(data)
		if size > len(p) {
			// 输入数据超过缓冲区大小，截断处理
			logger.Warn(fmt.Sprintf("输入数据过大: %d > %d，进行截断", size, len(p)))
			size = len(p)
		}
		copy(p, data[:size])
		logger.Debug(fmt.Sprintf("发送 %d 字节到容器: %q", size, string(data[:size])))
		return size, nil

	default:
		logger.Warn(fmt.Sprintf("未知消息类型: %s", xtermMsg.Type))
		return 0, nil
	}
}

// Write executor 回调向 Xterm WebSSH 端输出
func (s *StreamHandler) Write(p []byte) (size int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	// 优化：直接发送，避免不必要的内存拷贝
	// WebSocket 库内部会处理数据拷贝
	err = s.WsConn.WsWrite(websocket.TextMessage, p)
	if err != nil {
		logger.Error(fmt.Sprintf("WebSocket写入失败: %v", err))
		return 0, err
	}

	logger.Debug(fmt.Sprintf("成功发送 %d 字节到前端", len(p)))
	return len(p), nil
}
