package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/azzz/psychic-octo-pancake/internal/omap"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"io"
	"log"
)

type Server struct {
	queue amqp091.Queue
	conn  *amqp091.Connection
	ch    *amqp091.Channel

	logger *log.Logger

	store *omap.OMap[string, string]
	log   io.Writer

	stop chan struct{}
}

func (s *Server) SetLogger(logger *log.Logger) {
	s.logger = logger
}

func New(url, queue string, elog io.Writer) (*Server, error) {
	var (
		err error
		srv = &Server{
			stop:   make(chan struct{}),
			logger: log.Default(),
			log:    elog,
			store:  omap.New[string, string](),
		}
	)

	srv.conn, err = amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	srv.ch, err = srv.conn.Channel()
	if err != nil {
		_ = srv.Stop()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	srv.queue, err = srv.ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		_ = srv.Stop()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return srv, nil
}

func (s Server) Stop() error {
	s.logger.Print("INFO: stop the server")

	close(s.stop)
	var errs []error

	if s.ch != nil {
		s.logger.Print("INFO: close AMQP channel")
		if err := s.ch.Close(); err != nil {
			errs = append(errs, fmt.Errorf("cloe channel: %w", err))
		}
	}

	if s.conn != nil {
		s.logger.Print("INFO: close AMQP connection")
		if err := s.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (s Server) Start(ctx context.Context) error {
	s.logger.Print("INFO: start the server")
	updates, err := s.ch.Consume(s.queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	for {
		select {
		case <-s.stop:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case upd := <-updates:
			traceId, _ := uuid.NewUUID()
			err := s.handle(upd.Body)
			s.logger.Printf("INFO: traceid=%s: received message", traceId)

			if err == nil {
				if err := upd.Ack(false); err != nil {
					s.logger.Fatalf("FATAL: traceid=%s: ack message: %s", traceId, err)
				}
				break
			}

			if !errors.Is(err, ValueNotFoundErr) {
				s.logger.Printf("WARN: traceid=%s: handle message: %s", traceId, err)
			}

			if err := upd.Reject(false); err != nil {
				s.logger.Fatalf("FATAL: traceid=%s: reject message: %s", traceId, err)
			}
		}
	}
}

func (s Server) handle(data []byte) error {
	var msg Message

	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	cmd, err := s.buildCommand(msg)
	if err != nil {
		return fmt.Errorf("build command: %w", err)
	}

	if err := cmd.Exec(msg); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}

	return nil
}

func (s Server) buildCommand(msg Message) (Command, error) {
	switch msg.Command {
	case addItemCommand:
		return AddItemCommand{s.store}, nil
	case getItemCommand:
		return GetItemCommand{s.log, s.store}, nil
	case getAllItemsCommand:
		return GetAllCommand{s.log, s.store}, nil
	case removeItemCommand:
		return RemoveItemCommand{s.store}, nil
	}

	return nil, UnsupportedCommandErr
}
