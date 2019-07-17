package util

//非线程安全，使用前要自行上锁

type RingBuffer struct {
	buffer []byte
	head int
	tail int
	max_size int
	cur_size int
}

func NewRingBuffer(_buffer_size int) RingBuffer {
	return RingBuffer{
		buffer : make([]byte, _buffer_size),
		head : 0,
		tail : 0,
		max_size : _buffer_size,
		cur_size : 0,
	}
}

func (rb *RingBuffer) GetCapacity() int {
	return rb.max_size
}

func (rb *RingBuffer) GetLength() int {
	return rb.cur_size
}

func (rb *RingBuffer) GetHeadPos() int {
	return rb.head
}

func (rb *RingBuffer) GetTailPos() int {
	return rb.tail
}

// 当cur_size等于0时，head和tail应相等
func (rb *RingBuffer) IsEmpty() bool {
	return 0 == rb.cur_size
}

func (rb *RingBuffer) IsFull() bool {
	return rb.cur_size == rb.max_size
}

func (rb *RingBuffer) Push(data []byte) bool {
	data_len := len(data)

	if rb.cur_size + data_len > rb.max_size {
		return false
	}

	if rb.tail + data_len >= rb.max_size {
		size := rb.max_size - rb.tail
		copy(rb.buffer[rb.tail : rb.max_size], data[0 : size])
		rb.tail = 0
		if data_len > size {
			remain_len := data_len - size
			copy(rb.buffer[rb.tail : rb.tail + remain_len], data[size : data_len])
			rb.tail += remain_len
		}	
	} else {
		copy(rb.buffer[rb.tail : rb.tail + data_len], data)
		rb.tail += data_len
	}

	rb.cur_size += data_len
	return true
}

func (rb *RingBuffer) Pop(pop_size int) (bool, []byte, []byte) {
	if rb.cur_size < pop_size {
		return false, nil, nil
	} else {
		if rb.tail < rb.head {
			head := rb.head
			size := rb.max_size - head
			if size >= pop_size {
				rb.head = (rb.head + pop_size) % rb.max_size // 预防size == pop_size的情况
				rb.cur_size -= pop_size
				return true, rb.buffer[head : head + pop_size], nil
			} else {
				remain_size := pop_size - size
				rb.head = remain_size
				rb.cur_size -= pop_size
				return true, rb.buffer[head : head + size], rb.buffer[0 : remain_size]
			}
			

		} else {
			head := rb.head
			rb.head += pop_size
			rb.cur_size -= pop_size
			return true, rb.buffer[head : head + pop_size], nil
		}
	}
}
