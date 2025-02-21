package common

import (
	"fmt"

	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"google.golang.org/protobuf/proto"
)

func EncodeFrames(msgs []ehp.EsphomeMessageTyper) (retFrames []frameshakers.Frame, err error) {
	for _, msg := range msgs {
		id := ehp.MessageType(0)
		et, ok := msg.(ehp.EsphomeMessageTyper)
		if ok {
			id = et.EsphomeMessageType()
		}
		if !ok || id == 0 {
			return nil, fmt.Errorf("internal error: implement ID for type %T", msg)
		}
		raw, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		retFrames = append(retFrames, frameshakers.Frame{
			Type: int(id), Data: raw,
		})
	}
	return retFrames, nil
}

func DecodeFrame(frame frameshakers.Frame) (ehp.MessageType, ehp.EsphomeMessageTyper, error) {
	mt := ehp.MessageType(frame.Type)
	msg := ehp.MessageByType(mt)
	if msg == nil {
		return 0, nil, fmt.Errorf("unknown message: %d", frame.Type)
	}
	if err := proto.Unmarshal(frame.Data, msg); err != nil {
		return 0, nil, err
	}
	return mt, msg, nil
}

func Enum[To ~int32, From ~int32](from From) To {
	return To(int(from))
}

func Enums[To ~int32, From ~int32](from []From) []To {
	ret := make([]To, len(from))
	for i, e := range from {
		ret[i] = To(int(e))
	}
	return ret
}
