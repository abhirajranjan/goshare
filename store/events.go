package store

import (
	"archive/zip"
	"io"
	"time"

	"goshare/resources"
)

type broadcast struct {
	senderCode string
	validTill  time.Time
	files      []resources.Media

	// internal use
	senderNotifyCh chan *senderNotifyEvent
}

func (b *broadcast) ReceiverEvents() chan *senderNotifyEvent {
	return b.senderNotifyCh
}

func newBroadcast(senderCode string, validTill time.Time, files []resources.Media) broadcast {
	return broadcast{
		senderCode: senderCode,
		validTill:  validTill,
		files:      files,
	}
}

//

type transfer struct {
	currentFileName  string
	currentZipWriter io.Writer
	zipWriter        *zip.Writer
}

//

type ackErr struct {
	ack chan struct{}
	err *error
}

func (ae *ackErr) WaitForAck() {
	<-ae.ack
}

func (ae *ackErr) Ack() {
	ae.ack <- struct{}{}
}

func (ae *ackErr) SetError(err error) {
	ae.err = &err
}

func (ae *ackErr) Error() error {
	if ae.err == nil {
		return nil
	}

	return *ae.err
}

func newAckErr() ackErr {
	return ackErr{
		ack: make(chan struct{}, 1),
	}
}

// Events

type broadcastEvent struct {
	ackErr
	broadcast
	otp string
}

func newBroadcastEvent(otp string, b broadcast) *broadcastEvent {
	b.senderNotifyCh = make(chan *senderNotifyEvent, 10)

	return &broadcastEvent{
		ackErr:    newAckErr(),
		otp:       otp,
		broadcast: b,
	}
}

//

type getMetaDataEvent struct {
	ackErr
	// params
	otp string

	// returns
	media      []resources.Media
	senderCode string
}

func (md *getMetaDataEvent) Set(bc *broadcast) {
	md.media = bc.files
	md.senderCode = bc.senderCode
}

func newGetMetadataEvent(otp string) *getMetaDataEvent {
	return &getMetaDataEvent{
		ackErr: newAckErr(),
		otp:    otp,
	}
}

//

type receiverEvent struct {
	ackErr

	// params
	otp          string
	id           string
	receiverCode string
	w            io.Writer
}

func newReceiverEvent(otp, id, receiverCode string, w io.Writer) *receiverEvent {
	return &receiverEvent{
		otp:          otp,
		id:           id,
		receiverCode: receiverCode,
		w:            w,
	}
}

//

type transferEvent struct {
	ackErr

	// params
	id           string
	recieverCode string
	w            io.Writer
}

func newTransferEvent(id, recieverCode string, w io.Writer) *transferEvent {
	return &transferEvent{
		ackErr:       newAckErr(),
		id:           id,
		recieverCode: recieverCode,
		w:            w,
	}
}

//

type senderNotifyEvent struct {
	ackErr

	// params
	id           string
	receiverCode string
}

func newSenderNotifyEvent(id, receiverCode string) *senderNotifyEvent {
	return &senderNotifyEvent{
		ackErr:       newAckErr(),
		id:           id,
		receiverCode: receiverCode,
	}
}

//

type dataEvent struct {
	ackErr

	// param
	id       string
	filename string

	// return
	w io.Writer
}

func (sd *dataEvent) Set(tv *transfer) {
	sd.w = tv.currentZipWriter
}

func newDataEvent(id, filename string) *dataEvent {
	return &dataEvent{
		ackErr:   newAckErr(),
		id:       id,
		filename: filename,
	}
}
