package controller

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/k8s/stream"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

// PodController 定义控制器结构体
type PodController struct {
	pod service.PodService // 使用服务接口
}

// NewPodController 创建新的 DeploymentsController 实例
func NewPodController(pod service.PodService) *PodController {
	return &PodController{pod: pod}
}

func (c *PodController) List(ctx *gin.Context) {
	var (
		id, _         = strconv.Atoi(ctx.Param("id"))
		ns            = ctx.Param("namespace")
		labelSelector = ctx.Query("labelSelector")
	)
	data, err := c.pod.List(ctx, id, ns, labelSelector)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PodController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.pod.Get(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PodController) GetDetail(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.pod.GetDetail(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PodController) GetLogs(ctx *gin.Context) {
	var (
		id, _           = strconv.Atoi(ctx.Param("id"))
		namespace       = ctx.Param("namespace")
		name            = ctx.Param("name")
		container       = ctx.Param("container")
		sinceSeconds    = time.Now().Unix()
		tailLinesString = ctx.Query("tailLines")
		followString    = ctx.Query("follow")
		token           = ctx.GetHeader("Sec-Websocket-Protocol")
	)
	// 解析 tailLines (int64)
	tailLines, err := strconv.ParseInt(tailLinesString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"msg":     "the ailLines parameter is in the wrong format; it should be a number.",
		})
		return
	}
	// 解析 follow (Bool)
	follow, err := strconv.ParseBool(followString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"msg":     "the follow parameter is in the wrong format; it should be a boolean value.",
		})
		return
	}
	// 调用服务层，获取日志流
	stream, cancel, err := c.pod.GetLogs(ctx, id, namespace, name, container, follow, &tailLines, &sinceSeconds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"msg":     err,
		})
		return
	}
	defer cancel()
	defer stream.Close()
	logger.Infof("Kubernetes logs stream established for pod %s", name)
	if follow {
		var upGrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// 允许跨域（可根据实际情况收紧）
			CheckOrigin: func(r *http.Request) bool { return true },
			Subprotocols: []string{
				token,
			},
		}
		// 如果 follow 值为 true 时，升级为 WebSocket 连接
		conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "WebSocket upgrade failed", "err": err.Error()})
			return
		}
		defer conn.Close()
		logger.Infof("WebSocket connected for logs: cluster=%d, pod=%s, container=%s", id, name, container)

		reader := bufio.NewReader(stream)
		// 启动日志流循环
		for {
			line, err := reader.ReadBytes('\n')
			if len(line) > 0 {
				if writeErr := conn.WriteMessage(websocket.TextMessage, line); writeErr != nil {
					logger.Infof("WebSocket disconnected for logs (write error): cluster=%d, pod=%s, container=%s", id, name, container)
					return
				}
			}
			if err != nil {
				logger.Infof("Log stream ended or an exception occurred, error: %v", err)
				return
			}
		}
	} else {
		// 非实时日志，一次性返回
		content, err := io.ReadAll(stream)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"msg":     fmt.Sprintf("读取日志内容失败: %v", err),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": string(content)})
	}
}

func (c *PodController) GetTailLogs(ctx *gin.Context) {
	var (
		id, _           = strconv.Atoi(ctx.Param("id"))
		namespace       = ctx.Param("namespace")
		name            = ctx.Param("name")
		container       = ctx.Param("container")
		sinceSeconds    = time.Now().Unix()
		tailLinesString = ctx.Query("tailLines")
		followString    = ctx.Query("follow")
	)
	// 解析 tailLines (int64)
	tailLines, err := strconv.ParseInt(tailLinesString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"msg":     "the ailLines parameter is in the wrong format; it should be a number.",
		})
		return
	}
	// 解析 follow (Bool)
	follow, err := strconv.ParseBool(followString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    400,
			"msg":     "the follow parameter is in the wrong format; it should be a boolean value.",
		})
		return
	}
	goCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	stream, cancel, err := c.pod.GetLogs(goCtx, id, namespace, name, container, follow, &tailLines, &sinceSeconds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"msg":     err,
		})
		return
	}
	defer stream.Close()

	// 设置响应头
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.Header("X-Pod-Name", name)
	ctx.Header("X-Namespace", namespace)
	ctx.Header("X-Container", container)

	// 如果是实时日志，设置流式响应
	if follow {
		// 设置流式响应头
		ctx.Header("Transfer-Encoding", "chunked")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("X-Accel-Buffering", "no") // 禁用 Nginx 缓冲
		ctx.Status(http.StatusOK)

		// 流式传输日志
		reader := bufio.NewReader(stream)
		flusher, ok := ctx.Writer.(http.Flusher)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    500,
				"msg":     "Streaming not supported",
			})
			return
		}
		logger.Infof("Starting real-time logs stream: cluster=%d, pod=%s, container=%s", id, name, container)
		// 流式传输日志
		for {
			select {
			case <-ctx.Request.Context().Done():
				// 客户端断开连接
				logger.Infof("Client disconnected from logs stream: cluster=%d, pod=%s, duration=%v, lines=%d", id, name)
				return
			default:
				line, err := reader.ReadBytes('\n')
				if len(line) > 0 {
					// 写入数据并立即刷新
					if _, err := ctx.Writer.Write(line); err != nil {
						logger.Errorf("Write to client failed: %v", err)
						return
					}
					flusher.Flush() // 刷新到前端
				}
				if err != nil {
					if err == io.EOF {
						logger.Infof("Logs stream ended normally: cluster=%d, pod=%s", id, name)
					} else {
						logger.Errorf("Read from logs stream failed: %v", err)
						// 发送错误信息到客户端
						ctx.Writer.Write([]byte(fmt.Sprintf("\n=== Log stream error: %v ===\n", err)))
						flusher.Flush()
					}
					return
				}
			}
		}
	} else {
		// 非实时日志，一次性返回
		content, err := io.ReadAll(stream)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"msg":     fmt.Sprintf("读取日志内容失败: %v", err),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": string(content)})
	}
}

func (c *PodController) Executor(ctx *gin.Context) {
	var (
		wsConn    *stream.WsConnection
		executor  remotecommand.Executor
		handler   *stream.StreamHandler
		id, _     = strconv.Atoi(ctx.Param("id"))
		namespace = ctx.Param("namespace")
		name      = ctx.Param("name")
		container = ctx.Param("container")
		cols, _   = strconv.Atoi(ctx.Query("cols"))
		rows, _   = strconv.Atoi(ctx.Query("rows"))
		token     = ctx.GetHeader("Sec-Websocket-Protocol")
	)
	// 创建执行器
	executor, err := c.pod.Executor(ctx, id, namespace, name, container)
	if err != nil {
		logger.Errorf("kubernetes pod '%v' container '%v' executor creation failed, err: %v", name, container, err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	logger.Infof("kubernetes pod '%v' container '%v' executor created has been successfully.", name, container)
	// 建立 WebSocket
	wsConn, err = stream.WsConn(ctx.Writer, ctx.Request, token)
	if err != nil {
		logger.Errorf("kubernetes pod '%v' container '%v' websocket connection failed, err: %v", name, container, err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	logger.Infof("kubernetes pod '%v' container '%v' websocket connection has been successful.", name, container)
	defer wsConn.WsClose()

	// 创建流式 handler
	handler = &stream.StreamHandler{
		WsConn:      wsConn,
		ResizeEvent: make(chan remotecommand.TerminalSize, 1),
	}

	// 发送初始终端大小
	if cols > 0 && rows > 0 {
		select {
		case handler.ResizeEvent <- remotecommand.TerminalSize{Width: uint16(cols), Height: uint16(rows)}:
		default:
		}
	}

	// 执行远程命令
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	})
	if err != nil {
		logger.Errorf("kubernetes pod '%v' container '%v' executor stream failed, err: %v", name, container, err)
	}
	logger.Infof("kubernetes pod '%v' container '%v' executor stream has been successful.", name, container)
}

func (c *PodController) DebugExecutor(ctx *gin.Context) {
	var (
		id, _     = strconv.Atoi(ctx.Param("id"))
		namespace = ctx.Param("namespace")
		name      = ctx.Param("name")
		container = ctx.Param("container")
	)

	// 测试命令
	if err := c.pod.DebugExecutor(ctx, id, namespace, name, container); err != nil {
		logger.Errorf("DebugExecutor connection failed, error: %v\n", err)
	}
}
