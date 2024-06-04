package store

import (
	"io"
	"mime/multipart"

	"goshare/resources"

	"github.com/pkg/errors"
)

func (s *Store) NotifyEvent(senderCode string, media []resources.Media) <-chan resources.EventType {
	eventCh := make(chan resources.EventType)
	go s.notifyEvent(senderCode, media, eventCh)
	return eventCh
}

func (s *Store) notifyEvent(senderCode string, media []resources.Media, eventCh chan resources.EventType) {
	token, origin, validDuration := s.otpGen.GenerateOTP(senderCode)
	validTill := origin.Add(validDuration)

	event := newBroadcastEvent(token, newBroadcast(senderCode, validTill, media))
	s.broadcastEventCh <- event
	event.WaitForAck()

	if err := event.Error(); err != nil {
		panic(errors.Wrap(err, "store.NotifyEvent"))
	}

	eventCh <- resources.OTPEvent(token, validTill)

	for e := range event.ReceiverEvents() {
		e.Ack()
		eventCh <- resources.RecieverEvent(e.receiverCode, e.id)
	}
}

func (s *Store) GetMetaData(otp string) (media []resources.Media, senderCode string) {
	event := newGetMetadataEvent(otp)

	s.getMetadataEventCh <- event
	event.WaitForAck()

	if err := event.Error(); err != nil {
		return nil, ""
	}

	return event.media, event.senderCode
}

func (s *Store) GetData(otp string, recieverCode string) (io.Reader, error) {
	r, w := io.Pipe()
	id := transferId()

	receiverEvent := newReceiverEvent(otp, id, recieverCode, w)
	s.receiverEventCh <- receiverEvent
	receiverEvent.WaitForAck()

	if err := receiverEvent.Error(); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Store) SetData(id string, part *multipart.Part) error {
	event := newDataEvent(id, part.FileName())
	s.dataEventCh <- event
	event.WaitForAck()

	if err := event.Error(); err != nil {
		return err
	}

	_, err := io.Copy(event.w, part)
	return err
}
