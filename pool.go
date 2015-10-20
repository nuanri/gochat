package gochat

type ClientPool struct {
	pool_map map[string]*Conn
}

func NewClientPool(pool_map map[string]*Conn) *ClientPool {

	p := &ClientPool{
		pool_map: pool_map,
	}

	return p
}

func (p *ClientPool) Add(username string, conn *Conn) {
	p.pool_map[username] = conn
}

func (p *ClientPool) GetByConn(conn *Conn) string {
	username := ""
	for k, v := range p.pool_map {
		if conn == v {
			username = k
			break
		}
	}
	return username
}

func (p *ClientPool) GetByUsername(username string) (*Conn, bool) {
	client_v, ok := p.pool_map[username]
	return client_v, ok
}

func (p *ClientPool) DelByConn(conn *Conn) {

}

func (p *ClientPool) DelByUsername(username string) {

}

func (p *ClientPool) SendToAll(himsg *HiMsg) {
	for _, client_v := range p.pool_map {
		r_msg_b := GetJson(himsg)

		client_v.Send(r_msg_b)

	}
}
