package messager

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

type Messager interface {
	Get()
	Put(*pb.NetworkMessage) error
}

type MessageBuffer struct {
	// messages in our heap
	inHeapLock sync.Mutex
	inHeap     map[uint64]bool

	// messages we've made use of
	crossListLock sync.Mutex
	crossList     map[uint64]bool

	// sorted list
	mHeap *messageHeap
}

type messageHeap []*pb.NetworkMessage

func (mh messageHeap) Len() int      { return len(mh) }
func (mh messageHeap) Swap(i, j int) { mh[i], mh[j] = mh[j], mh[i] }
func (mh messageHeap) Less(i, j int) bool {
	mhiNonce := binary.LittleEndian.Uint64(mh[i].Nonce)
	mhjNonce := binary.LittleEndian.Uint64(mh[j].Nonce)
	return mhiNonce < mhjNonce
}

func (mh *messageHeap) Push(x interface{}) {
	*mh = append(*mh, x.(*pb.NetworkMessage))
}

func (mh *messageHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	x := old[n-1]
	*mh = old[0 : n-1]

	return x
}

func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		inHeap:    make(map[uint64]bool),
		crossList: make(map[uint64]bool),
		mHeap:     &messageHeap{},
	}
}

func (mb *MessageBuffer) Get() (*pb.NetworkMessage, error) {
	var nonce uint64
	mb.crossListLock.Lock()
	defer mb.crossListLock.Unlock()

	nextValue := mb.mHeap.Pop()

	switch nextValue {
	case nextValue.(*pb.NetworkMessage):
		nonce = binary.LittleEndian.Uint64(nextValue.(*pb.NetworkMessage).Nonce)
		mb.crossList[nonce] = true
	default:
		return nil, fmt.Errorf("")

	}

	mb.inHeapLock.Lock()
	delete(mb.inHeap, nonce)
	mb.inHeapLock.Unlock()

	return nextValue.(*pb.NetworkMessage), nil
}

func (mb *MessageBuffer) Put(message *pb.NetworkMessage) error {
	mb.inHeapLock.Lock()
	defer mb.inHeapLock.Unlock()

	nonce := binary.LittleEndian.Uint64(message.Nonce)

	if _, ok := mb.inHeap[nonce]; ok {
		// we've already seen this value
		return fmt.Errorf("")
	}

	mb.crossListLock.Lock()
	defer mb.crossListLock.Unlock()

	if _, ok := mb.crossList[nonce]; ok {
		// we've already seen this value
		return fmt.Errorf("")
	}

	mb.mHeap.Push(message)
	mb.inHeap[nonce] = true

	return nil
}
