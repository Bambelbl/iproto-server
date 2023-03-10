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

const (
	CLIENT_INVALID_BODY      = 401
	CLIENT_TOO_MANY_REQUESTS = 402
)

type IprotoServer struct {
	listener        net.Listener
	logger          *log.Logger
	quit            chan struct{}
	queueForClients chan struct{}
	wg              sync.WaitGroup
	stor            *storage.Storage
	rateLimiter     *rate_limiter.RateLimiter
}

// NewIprotoServer initializes IprotoServer and starts it to listen
func NewIprotoServer(addr string, logger *log.Logger, maxClients int, scale_rps int64, limit_rps uint32) *IprotoServer {
	s := &IprotoServer{
		logger:          logger,
		quit:            make(chan struct{}),
		queueForClients: make(chan struct{}, maxClients),
		rateLimiter:     rate_limiter.NewRateLimiter(logger, scale_rps, limit_rps),
	}
	stor := storage.NewSimpleStorageRepo()
	s.stor = &stor
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatal("Server: listen err: %s", err.Error())
	}
	s.listener = l
	return s
}

// Serve listen and serve for IprotoServer
func (s *IprotoServer) Serve() {
	s.logger.Println("Server starts to serve...")
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
					s.logger.Println("Server: accept error: %s", err)
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
						s.logger.Println("Server: timeout for handler")
					case <-endOfHandler:
						s.logger.Println("Server: handler finished")
					}
					close(endOfHandler)
					s.wg.Done()
				}()
			}
		}
	}()
}

// handleConnection handler for incoming requests to IprotoServer
func (s *IprotoServer) handleConnection(conn net.Conn, ctx context.Context, endOfHandler chan struct{}) {
	defer func() {
		<-s.queueForClients
		endOfHandler <- struct{}{}
		err := conn.Close()
		if err != nil {
			s.logger.Println("Server: connection close error: %s", err.Error())
		}
	}()
	buf := make([]byte, MAX_PACKET_SIZE)
	_, err := conn.Read(buf)
	if err != nil {
		s.logger.Println("Server: read from request error: %s", err.Error())
		return
	}
	var responseBody string
	var returnCode uint32
	var requestPacket request_packet.IprotoPacketRequest
	if s.rateLimiter.ValidRate(conn.RemoteAddr().String()) {
		requestPacket, err = request_packet.Unmarshal(buf)
		if err != nil {
			s.logger.Println("Server: unmarshal error: %s", err.Error())
			responseBody = "Invalid body in request packet"
			returnCode = CLIENT_INVALID_BODY
		} else {
			responseBody, returnCode = api.Handler(requestPacket, s.stor)
		}
	} else {
		responseBody = "Too many requests"
		returnCode = CLIENT_TOO_MANY_REQUESTS
	}
	response, err := response_packet.Marshal(response_packet.IprotoPacketResponse{
		Header: response_packet.IprotoHeader{
			Func_id:     requestPacket.Header.Func_id,
			Body_length: 0,
			Request_id:  requestPacket.Header.Request_id},
		Return_code: returnCode,
		Body:        responseBody,
	})
	if err != nil {
		s.logger.Println("Server: marshal response error: %s", err.Error())
		return
	}
	_, err = conn.Write(response)
	if err != nil {
		s.logger.Println("Server: write response error: %s", err.Error())
		return
	}
}

// Stop shutdown to IprotoServer
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
