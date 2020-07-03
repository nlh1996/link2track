package ws

// Request .
type Request struct {
	conn     *Connection
	ID			 string
	ByteData []byte
}

// GetConn .
func (r *Request) GetConn() *Connection {
	return r.conn
}

// Send .
func (r *Request) Send(msg []byte) (err error) {
	return r.conn.Send(msg)
}

