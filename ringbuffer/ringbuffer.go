package ringbuffer

type RingBuffer struct {
	pr  int    // 开始读的位置
	pw  int    // 开始写的位置
	buf []byte // 实际buffer
}

// RingBuffer的容量为capacity，实际上内部的buffer大小应该为capacity+1
func New(capacity int) *RingBuffer {
	return &RingBuffer{
		pr:  0,
		pw:  0,
		buf: make([]byte, capacity+1),
	}
}

// 向RingBuffer写入一个slice，若剩余空间小于len(p)，则返回0，否则返回len(p)
func (this *RingBuffer) Write(p []byte) int {
	if this.FreeSize() < len(p) {
		return 0
	}

	n := len(this.buf) - this.pw

	if this.pw < this.pr || len(p) <= n {
		this.pw += copy(this.buf[this.pw:], p)
	} else {
		copy(this.buf[this.pw:], p[0:n])
		this.pw = copy(this.buf[0:], p[n:])
	}

	return len(p)
}

// 从RingBuffer中读取n个byte到p中，返回读取的长度，若RingBuffer的Size不足n，则返回0
func (this *RingBuffer) Read(p []byte, n int) int {
	if this.Size() < n {
		return 0
	}

	m := len(this.buf) - this.pr

	if this.pr < this.pw || n <= m {
		this.pr += copy(p[:n], this.buf[this.pr:])
	} else {
		copy(p, this.buf[this.pr:])
		this.pr = copy(p[m:n], this.buf[0:])
	}

	return n
}

// 返回RingBuffer中的已用空间大小
func (this *RingBuffer) Size() int {
	if this.pr <= this.pw {
		return this.pw - this.pr
	} else {
		return len(this.buf) - this.pr + this.pw
	}
}

// 返回RingBuffer的总容量
func (this *RingBuffer) Capacity() int {
	return len(this.buf) - 1
}

// 返回RingBuffer的未用空间大小
func (this *RingBuffer) FreeSize() int {
	return this.Capacity() - this.Size()
}

// 重置RingBuffer
func (this *RingBuffer) Reset() {
	this.pr, this.pw = 0, 0
}

// 删除最新的n个byte，若Size不足n，则重置RingBuffer
func (this *RingBuffer) RemoveNewest(n int) {
	if this.Size() <= n {
		this.Reset()
		return
	}
	if this.pr < this.pw || n <= this.pw {
		this.pw -= n
	} else {
		this.pw = len(this.buf) - n + this.pw
	}
}

// 删除最旧的n个byte，若Size不足n，则重置RingBuffer
func (this *RingBuffer) RemoveOldest(n int) {
	if this.Size() <= n {
		this.Reset()
		return
	}
	if this.pr < this.pw || n <= len(this.buf)-this.pr {
		this.pr += n
	} else {
		this.pr = n - len(this.buf) + this.pr
	}
}
