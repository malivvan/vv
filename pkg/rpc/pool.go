package rpc

import (
	"github.com/malivvan/vv/pkg/rpc/bytesutil"
	"sync"
	"time"
)

type Message struct {
	conn *Conn
	seq  uint32
	buf  []byte
}

func (msg *Message) Conn() *Conn            { return msg.conn }
func (msg *Message) Body() []byte           { return msg.buf }
func (msg *Message) Reply(buf []byte) error { return msg.conn.send(msg.seq, buf) }

var contextPool sync.Pool

func acquireMessage(conn *Conn, seq uint32, buf []byte) *Message {
	v := contextPool.Get()
	if v == nil {
		v = &Message{}
	}
	msg := v.(*Message)
	msg.conn = conn
	msg.seq = seq
	msg.buf = buf
	return msg
}

func releaseMessage(msg *Message) { contextPool.Put(msg) }

type pendingWrite struct {
	buf  *bytesutil.ByteBuffer // payload
	wait bool                  // signal to caller if they're waiting
	err  error                 // keeps track of any socket errors on write
	wg   sync.WaitGroup        // signals the caller that this write is complete
}

var pendingWritePool sync.Pool

func acquirePendingWrite(buf *bytesutil.ByteBuffer, wait bool) *pendingWrite {
	v := pendingWritePool.Get()
	if v == nil {
		v = &pendingWrite{}
	}
	pw := v.(*pendingWrite)
	pw.buf = buf
	pw.wait = wait
	return pw
}

func releasePendingWrite(pw *pendingWrite) { pw.err = nil; pendingWritePool.Put(pw) }

type pendingRequest struct {
	dst []byte         // dst to copy response to
	err error          // error while waiting for response
	wg  sync.WaitGroup // signals the caller that the response has been received
}

var pendingRequestPool sync.Pool

func acquirePendingRequest(dst []byte) *pendingRequest {
	v := pendingRequestPool.Get()
	if v == nil {
		v = &pendingRequest{}
	}
	pr := v.(*pendingRequest)
	pr.dst = dst
	return pr
}

func releasePendingRequest(pr *pendingRequest) {
	pr.dst = nil
	pr.err = nil
	pendingRequestPool.Put(pr)
}

var zeroTime time.Time

var timerPool sync.Pool

func AcquireTimer(timeout time.Duration) *time.Timer {
	v := timerPool.Get()
	if v == nil {
		return time.NewTimer(timeout)
	}
	t := v.(*time.Timer)
	t.Reset(timeout)
	return t
}

func ReleaseTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	timerPool.Put(t)
}
