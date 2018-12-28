package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"v2ray.com/core/app/proxyman/command"
	statscmd "v2ray.com/core/app/stats/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

const (
	API_ADDRESS = "127.0.0.1"
	API_PORT    = 10085
	INBOUND_TAG = "proxy"
	LEVEL       = 0
	EMAIL       = "123@gmail.com"
	UUID        = "2601070b-ab53-4352-a290-1d44414581ee"
	ALTERID     = 32
)

type handlerServiceClient struct {
	API_ADDRESS string
	API_PORT    int32
	sClient     statscmd.StatsServiceClient
}

func NewStats(addr string, port int32) *handlerServiceClient {
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	sClient := statscmd.NewStatsServiceClient(cmdConn)
	svr := handlerServiceClient{API_ADDRESS: addr, API_PORT: port, sClient: sClient}
	return &svr
}

func (h *handlerServiceClient) QueryStats(pattern string, reset bool) {
	sresp, err := h.sClient.QueryStats(context.Background(), &statscmd.QueryStatsRequest{
		Pattern: pattern,
		Reset_:  reset,
	})

	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		// log.Printf("%v", sresp)
		for _, stat := range sresp.Stat {
			log.Printf("%s:%d\n", stat.Name, stat.Value)
		}
	}
}

func (h *handlerServiceClient) GetStats(name string, reset bool) {
	sresp, err := h.sClient.GetStats(context.Background(), &statscmd.GetStatsRequest{
		Name:   name,
		Reset_: reset,
	})

	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		log.Println("GetStats")
		log.Printf("%s:%d\n", sresp.Stat.Name, sresp.Stat.Value)
	}
}

func addUser(c command.HandlerServiceClient) {
	resp, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: INBOUND_TAG,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: LEVEL,
				Email: EMAIL,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               UUID,
					AlterId:          ALTERID,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		log.Printf("ok: %v", resp)
	}
}
func removeUser(c command.HandlerServiceClient) {
	resp, err := c.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: INBOUND_TAG,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: EMAIL,
		}),
	})
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		log.Printf("ok: %v", resp)
	}
}

func main() {
	// cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", API_ADDRESS, API_PORT), grpc.WithInsecure())
	// if err != nil {
	// 	panic(err)
	// }
	// sClient := statscmd.NewStatsServiceClient(cmdConn)
	// queryStat(sClient, "")
	// hsClient := command.NewHandlerServiceClient(cmdConn)
	// addUser(hsClient)
	// removeUser(hsClient)

	stats := NewStats("127.0.0.1", 10085)
	stats.QueryStats("", false)
	stats.GetStats("inbound>>>api>>>traffic>>>downlink", false)
}
