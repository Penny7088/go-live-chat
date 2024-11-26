package ws

import (
	"context"
	"errors"
	"strconv"
	"time"

	"lingua_exchange/pkg/timewheel"
)

const (
	heartbeatInterval = 30 // 心跳检测的间隔时间（秒）
	heartbeatTimeout  = 75 // 心跳检测的超时时间（超时时间为间隔的2.5倍以上）
)

var health *heartbeat

// heartbeat 结构体用于管理客户端的心跳检测
type heartbeat struct {
	timeWheel *timewheel.SimpleTimeWheel[*Client] // 时间轮定时器，用于心跳管理
}

// 初始化心跳管理
// 使用一个 1 秒精度的时间轮来管理心跳检测任务
func init() {
	health = &heartbeat{
		timeWheel: timewheel.NewSimpleTimeWheel[*Client](1*time.Second, 100, health.handle),
	}
}

// Start 启动心跳检测的后台协程
// ctx: 用于控制心跳检测生命周期的上下文
func (h *heartbeat) Start(ctx context.Context) error {
	// 启动时间轮
	go h.timeWheel.Start()

	// 等待上下文结束信号
	<-ctx.Done()

	// 停止时间轮
	h.timeWheel.Stop()

	// 返回错误以通知退出
	return errors.New("heartbeat exit")
}

// insert 添加客户端到心跳检测队列
// c: 客户端对象
func (h *heartbeat) insert(c *Client) {
	// 将客户端添加到时间轮队列，等待 heartbeatInterval 秒后进行下一次检测
	h.timeWheel.Add(strconv.FormatInt(c.cid, 10), c, time.Duration(heartbeatInterval)*time.Second)
}

// delete 从心跳检测队列中移除客户端
// c: 客户端对象
func (h *heartbeat) delete(c *Client) {
	// 根据客户端的唯一 ID 从时间轮中移除
	h.timeWheel.Remove(strconv.FormatInt(c.cid, 10))
}

// handle 处理心跳检测的回调函数
// timeWheel: 当前时间轮对象
// key: 客户端标识符
// c: 客户端对象
func (h *heartbeat) handle(timeWheel *timewheel.SimpleTimeWheel[*Client], key string, c *Client) {
	// 如果客户端已关闭，直接返回
	if c.Closed() {
		return
	}

	// 计算上次心跳时间与当前时间的差值
	interval := int(time.Now().Unix() - c.lastTime)

	// 如果超时则关闭连接
	if interval > heartbeatTimeout {
		c.Close(2000, "心跳检测超时，连接已关闭")
		return
	}

	// 如果超过心跳间隔时间，则推送一次 "ping" 消息作为心跳
	if interval > heartbeatInterval {
		_ = c.Write(&ClientResponse{Event: "ping"})
	}

	// 重新将客户端添加到时间轮，等待下一次心跳检测
	timeWheel.Add(key, c, time.Duration(heartbeatInterval)*time.Second)
}
