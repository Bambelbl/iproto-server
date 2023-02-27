package server

import (
	"context"
	"github.com/Bambelbl/iproto-server/api"
	"github.com/Bambelbl/iproto-server/packet/request_packet"
	"github.com/Bambelbl/iproto-server/packet/response_packet"
	"github.com/Bambelbl/iproto-server/rate_limiter"
	"github.com/Bambelbl/iproto-server/storage"
	"log"
	"net"
	"sync"
	"time"
)

const (
	MAX_PACKET_SIZE = 350
)

type IprotoServer struct {
	listener        net.Listener
	loger           *log.Logger
	quit            chan struct{}
	queueForClients chan struct{}
	wg              sync.WaitGroup
	stor            *storage.Storage
	rateLimiter     *rate_limiter.RateLimiter
}

func NewServer(addr string, loger *log.Logger, maxClients int, scale_rps int64, limit_rps uint32) *IprotoServer {
	s := &IprotoServer{
		loger:           loger,
		quit:            make(chan struct{}),
		queueForClients: make(chan struct{}, maxClients),
		rateLimiter:     rate_limiter.NewRateLimiter(loger, scale_rps, limit_rps),
	}
	stor := storage.NewSimpleStorageRepo()
	s.stor = &stor
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.loger.Fatal("Server: listen err: %s", err.Error())
	}
	s.listener = l
	return s
}

func (s *IprotoServer) Serve() {
	s.loger.Println("Server starts to serve...")
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					s.loger.Println("Server: accept error: %s", err)
				}
			} else {
				s.wg.Add(1)
				s.queueForClients <- struct{}{}
				go func() {
					endOfHandler := make(chan struct{})
					ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
					go s.handleConnection(conn, ctx, endOfHandler)
					select {
					case <-ctx.Done():
						s.loger.Println("Server: timeout for handler")
					case <-endOfHandler:
						s.loger.Println("Server: handler finished")
					}
					close(endOfHandler)
					s.wg.Done()
				}()
			}
		}
	}()
}

func (s *IprotoServer) handleConnection(conn net.Conn, ctx context.Context, endOfHandler chan struct{}) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.loger.Println("Server: connection close error: %s", err.Error())
		}
	}(conn)
	buf := make([]byte, MAX_PACKET_SIZE)
	_, err := conn.Read(buf)
	if err != nil {
		s.loger.Println("Server: read from request error: %s", err.Error())
		<-s.queueForClients
		endOfHandler <- struct{}{}
		return
	}
	var responseBody string
	var returnCode uint32
	var requestPacket request_packet.IprotoPacketRequest
	if s.rateLimiter.ValidRate(conn.RemoteAddr().String()) {
		requestPacket, err = request_packet.Unmarshal(buf)
		if err != nil {
			s.loger.Println("Server: unmarshal error: %s", err.Error())
			responseBody = "Invalid body in request packet"
			returnCode = 1
		} else {
			responseBody, returnCode = api.Handler(requestPacket, s.stor)
		}
	} else {
		responseBody = "Too many requests"
		returnCode = 1
	}
	_, err = conn.Write(response_packet.Marshal(response_packet.IprotoPacketResponse{
		Header: response_packet.IprotoHeader{
			Func_id:     requestPacket.Header.Func_id,
			Body_length: 0,
			Request_id:  requestPacket.Header.Request_id},
		Return_code: returnCode,
		Body:        responseBody,
	}))
	if err != nil {
		s.loger.Println("Server: write response error: %s", err.Error())
		<-s.queueForClients
		endOfHandler <- struct{}{}
		return
	}
	<-s.queueForClients
	endOfHandler <- struct{}{}
}

func (s *IprotoServer) Stop() error {
	close(s.quit)
	close(s.queueForClients)
	s.rateLimiter.Stop()
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}
