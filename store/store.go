package store

import (
	"archive/zip"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var ErrNotExist error = errors.New("does not exist")
var ErrAlreadyExist error = errors.New("already exist")

type OtpGenerator interface {
	GenerateOTP(secret string) (token string, origin time.Time, interval time.Duration)
}

type Store struct {
	otpGen     OtpGenerator
	broadcasts map[string]broadcast
	transfers  map[string]transfer

	broadcastEventCh   chan *broadcastEvent
	getMetadataEventCh chan *getMetaDataEvent
	receiverEventCh    chan *receiverEvent

	transferEventCh chan *transferEvent
	dataEventCh     chan *dataEvent
}

func NewStore(otpGen OtpGenerator) *Store {
	s := &Store{
		otpGen:     otpGen,
		broadcasts: map[string]broadcast{},
		transfers:  map[string]transfer{},

		broadcastEventCh:   make(chan *broadcastEvent),
		getMetadataEventCh: make(chan *getMetaDataEvent),
		receiverEventCh:    make(chan *receiverEvent),
		transferEventCh:    make(chan *transferEvent),
		dataEventCh:        make(chan *dataEvent),
	}

	s.runEventLoop()
	return s
}

func (s *Store) runEventLoop() {
	go s.runDataEventLoop()
	go s.runTransferEventLoop()
}

func (s *Store) runDataEventLoop() {
	for {
		select {
		case req := <-s.broadcastEventCh:
			if _, ok := s.broadcasts[req.otp]; ok {
				req.SetError(ErrAlreadyExist)
				req.Ack()
				continue
			}

			s.broadcasts[req.otp] = req.broadcast
			req.Ack()

		case req := <-s.getMetadataEventCh:
			dv, ok := s.broadcasts[req.otp]
			if !ok {
				req.SetError(ErrNotExist)
				req.Ack()
				continue
			}

			req.Set(&dv)
			req.Ack()

		case req := <-s.receiverEventCh:
			if _, ok := s.broadcasts[req.otp]; !ok {
				req.SetError(ErrNotExist)
				req.Ack()
				continue
			}

			senderCh := s.broadcasts[req.otp].senderNotifyCh

			go func() {
				transferEvent := newTransferEvent(req.id, req.receiverCode, req.w)
				s.transferEventCh <- transferEvent
				transferEvent.WaitForAck()

				senderNotifyEvent := newSenderNotifyEvent(req.id, req.receiverCode)
				senderCh <- senderNotifyEvent
				senderNotifyEvent.WaitForAck()

				req.Ack()
			}()
		}
	}
}

func (s *Store) runTransferEventLoop() {
	for {
		select {
		case req := <-s.transferEventCh:
			s.transfers[req.id] = transfer{
				zipWriter: zip.NewWriter(req.w),
			}

			req.Ack()

		case req := <-s.dataEventCh:
			tv, ok := s.transfers[req.id]
			if !ok {
				req.SetError(ErrNotExist)
				req.Ack()
				continue
			}

			if tv.currentFileName != req.filename {
				w, err := tv.zipWriter.CreateHeader(&zip.FileHeader{
					Name: req.filename,
				})

				if err != nil {
					req.SetError(err)
					req.Ack()
					slog.Warn("tv.zipWriter.CreateHeader", slog.String("error", err.Error()))
					continue
				}

				s.transfers[req.id] = transfer{
					currentFileName:  req.filename,
					zipWriter:        tv.zipWriter,
					currentZipWriter: w,
				}

				tv = s.transfers[req.id]
			}

			req.Set(&tv)
			req.Ack()
		}
	}
}

func transferId() string {
	return uuid.NewString()
}
