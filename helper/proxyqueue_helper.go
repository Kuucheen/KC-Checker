package helper

// ProxyQueue This is for displaying the latest proxies in the threadPhase
type ProxyQueue struct {
	data []*Proxy
}

// Enqueue adds a proxy value to the queue.
func (q *ProxyQueue) Enqueue(p *Proxy) {
	if len(q.data) < 5 {
		q.data = append(q.data, p)
	} else {
		q.Dequeue()
		q.Enqueue(p)
	}
}

// Dequeue removes and returns the front proxy value from the queue.
func (q *ProxyQueue) Dequeue() {
	if len(q.data) > 0 {
		q.data = append(q.data[:0], q.data[1:]...)
	}
}
func (q *ProxyQueue) Data() []*Proxy {
	return q.data
}
