package pubsub

import (
	"goshare/internal/models"
	"strings"
	"sync"
	_ "unsafe" // for go linkname

	"github.com/pkg/errors"
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

type ProtoFunc func(Id uint64, data []byte) (resid uint64, resData []byte)

type Pubsub struct {
	mapper      map[uint64][]byte
	mapperMu    sync.RWMutex
	protoMapper map[string]ProtoFunc
}

func NewPubsub() *Pubsub {
	p := &Pubsub{
		mapper:      map[uint64][]byte{},
		mapperMu:    sync.RWMutex{},
		protoMapper: map[string]ProtoFunc{},
	}
	p.initProto()
	return p
}

func (p *Pubsub) initProto() {
	p.protoMapper[DirectProtoTag] = directProto
}

func (p *Pubsub) GenerateRoom() uint64 {
	return uint64(fastrand())
}

func (p *Pubsub) Broadcast(model *models.BroadcastModel) error {
	proto := model.Proto
	if proto == "" {
		proto = DirectProtoTag
	}

	proto = sanitizeProto(proto)
	protofn, ok := p.protoMapper[proto]
	if !ok {
		return errors.Wrapf(models.DomainErr{}, "undefined %s protocol", model.Proto)
	}

	id, data := protofn(model.Id, model.Data)
	p.mapperMu.Lock()
	p.mapper[id] = data
	p.mapperMu.Unlock()

	return nil
}

func (p *Pubsub) FetchRoom(roomID uint64) ([]byte, error) {
	p.mapperMu.RLock()
	defer p.mapperMu.RUnlock()

	return p.mapper[roomID], nil
}

func sanitizeProto(protoName string) string {
	return strings.ToLower(protoName)
}
