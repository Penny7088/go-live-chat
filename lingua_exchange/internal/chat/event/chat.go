package event

import (
	"log"

	"github.com/redis/go-redis/v9"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/config"
	"lingua_exchange/internal/dao"
	"lingua_exchange/pkg/socket"
)

type ChatEvent struct {
	Redis           *redis.Client
	Config          *config.Config
	RoomStorage     cache.ChatRoomCache
	GroupMemberRepo dao.GroupMemberDao
}

func (c *ChatEvent) OnOpen(client socket.IClient) {
	log.Println("OnOpen client:", client)

}

func (c *ChatEvent) OnMessage(client socket.IClient, message []byte) {
	log.Println("OnMessage client:", client)
}

func (c *ChatEvent) OnClose(client socket.IClient, code int, text string) {
	log.Println("OnClose client:", client)
}
