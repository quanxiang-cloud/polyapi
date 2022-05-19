package arrange

// large deque cap threshold
const dqTooLarge = 4096

// show size when grow
const dqDebugShowGrow = false

// nodeDeque impliments a deque of *Node
type nodeDeque struct {
	// real data is [head,tail)
	// buffer d is cycle, that is to say, next(len(d)-1)=0, prev(0)=len(d)-1
	// so if tail<head, data range is [head, end, 0, tail)
	// head points to the first elem  available for read
	// tail points to the first space available for write
	head, tail int
	d          []*Node
}

// newNodeDeque new deque object
func newNodeDeque(bufSize int) *nodeDeque {
	r := &nodeDeque{}
	r.Init(bufSize)
	return r
}

// Init with buffer
func (dq *nodeDeque) Init(bufSize int) {
	if nil == dq.d {
		if bufSize <= 0 {
			bufSize = 8 //default buffer size
		}
		dq.newBuf(bufSize)
	}
	dq.Clear()
	return
}

// newBuf create new buffer
func (dq *nodeDeque) newBuf(bufSize int) {
	if bufSize > 0 {
		if dqDebugShowGrow {
			println("newBuf", dq.Cap(), "=>", bufSize)
		}
		dq.d = make([]*Node, bufSize, bufSize) //with the same cap and len
	}
}

// Clear all deque data
func (dq *nodeDeque) Clear() {
	dq.head, dq.tail = 0, 0
}

// PushFront push to front of deque
func (dq *nodeDeque) PushFront(v *Node) (ok bool) {
	if ok = true; ok {
		if nil == dq.d { //init if needed
			dq.Init(-1)
		}

		dq.head = dq.prev(dq.head) //move head to prev empty space
		dq.d[dq.head] = v
		if dq.head == dq.tail { //head reaches tail, buffer full
			dq.grow()
		}
	}
	return
}

// PushBack push to back of deque
func (dq *nodeDeque) PushBack(v *Node) (ok bool) {
	if ok = true; ok {
		if nil == dq.d { //init if needed
			dq.Init(-1)
		}
		dq.d[dq.tail] = v
		dq.tail = dq.next(dq.tail)
		if dq.tail == dq.head { // tail catches up with head, buffer full
			dq.grow()
		}
	}
	return
}

// nextCap get the next buffer size when grow
func (dq *nodeDeque) nextCap() int {
	oldCap := dq.Cap()
	if oldCap < dqTooLarge { // little size, 2*oldCap=>2^(n+1)
		newCap := oldCap * 2
		// if newCap!=2^n, then newCap=>2^(n+1), eg: 3=>6=>8
		// loop to remove the lowest binary digit 1, eg: 10110=>10100=>10000
		for t := 2 * (newCap & (newCap - 1)); t != 0; t &= (t - 1) {
			newCap = t
		}
		return newCap
	}

	// large size, grow by dqTooLarge, at least +50%*dqTooLarge
	const m = dqTooLarge // 4096
	return ((oldCap+m/2)/m + 1) * m
}

// grow when buffer is full
func (dq *nodeDeque) grow() {
	if dq.tail == dq.head { // tail catches up with head, buffer full
		d := dq.d
		dq.newBuf(dq.nextCap())
		h := copy(dq.d, d[dq.head:])
		t := copy(dq.d[h:], d[:dq.tail])
		dq.head, dq.tail = 0, h+t
	}
}

// PopFront pop front of deque
func (dq *nodeDeque) PopFront() (front *Node, ok bool) {
	if ok = dq.head != dq.tail; ok {
		front = dq.d[dq.head]
		dq.head = dq.next(dq.head)
	}
	return
}

// PopBack pop back of deque
func (dq *nodeDeque) PopBack() (back *Node, ok bool) {
	if ok = dq.head != dq.tail; ok {
		dq.tail = dq.prev(dq.tail)
		back = dq.d[dq.tail]
	}
	return
}

// Front get front data
func (dq *nodeDeque) Front() (front *Node, ok bool) {
	if ok = dq.head != dq.tail; ok {
		front = dq.d[dq.head]
	}
	return
}

// Back get back data
func (dq *nodeDeque) Back() (back *Node, ok bool) {
	if ok = dq.head != dq.tail; ok {
		t := dq.prev(dq.tail)
		back = dq.d[t]
	}
	return
}

// Cap get data buffer size
func (dq *nodeDeque) Cap() int {
	return len(dq.d)
}

// Size get size of deque
func (dq *nodeDeque) Size() (size int) {
	if dq.tail >= dq.head {
		size = dq.tail - dq.head
	} else {
		size = dq.Cap() - (dq.head - dq.tail)
	}
	return
}

// Empty check if deque is empty
func (dq *nodeDeque) Empty() bool {
	return dq.head == dq.tail
}

// next buff index
func (dq *nodeDeque) next(idx int) (r int) {
	if r = idx + 1; r >= dq.Cap() {
		r = 0
	}
	return
}

// prev buff index
func (dq *nodeDeque) prev(idx int) (r int) {
	if r = idx - 1; r < 0 {
		r = dq.Cap() - 1
	}
	return
}
