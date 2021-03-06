package tcp_server

import (
  "bufio"
  "log"
  "net"
)

// Client holds info about connection
type Client struct {
  conn     net.Conn
  Server   *server
  incoming chan string // Channel for incoming data from client
}

// TCP server
type server struct {
  clients                  []*Client
  address                  string        // Address to open connection: localhost:9999
  joins                    chan net.Conn // Channel for new connections
  onNewClientCallback      func(c *Client)
  onClientConnectionClosed func(c *Client, err error)
  onNewMessage             func(c *Client, message string)
  listener                 net.Listener
}

// Read client data from channel
func (c *Client) listen() {
  reader := bufio.NewReader(c.conn)
  for {
    message, err := reader.ReadString('\n')
    if err != nil {
      c.conn.Close()
      c.Server.onClientConnectionClosed(c, err)
      return
    }
    c.Server.onNewMessage(c, message)
  }
}

func (c *Client) Send(message string) error {
  _, err := c.conn.Write([]byte(message))
  return err
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
  s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
  s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, message string)) {
  s.onNewMessage = callback
}

// Creates new Client instance and starts listening
func (s *server) newClient(conn net.Conn) {
  client := &Client{
    conn:   conn,
    Server: s,
  }
  go client.listen()
  s.onNewClientCallback(client)
}

// Listens new connections channel and creating new client
func (s *server) listenChannels() {
  for {
    select {
    case conn := <-s.joins:
      s.newClient(conn)
    }
  }
}

func (s *server) Close() {
  s.listener.Close()
}

// Start network server
func (s *server) Listen() {
  go s.listenChannels()

  listener, err := net.Listen("tcp", s.address)
  if err != nil {
    log.Fatal("Error starting TCP server.")
  }

  s.listener = listener

  for {
    conn, _ := listener.Accept()
    //if err != nil {
    //  break
    //}
    s.joins <- conn
  }
}

// Creates new tcp server instance
func New(address string) *server {
  log.Println("Creating server with address", address)
  server := &server{
    address: address,
    joins:   make(chan net.Conn),
  }

  server.OnNewClient(func(c *Client) {})
  server.OnNewMessage(func(c *Client, message string) {})
  server.OnClientConnectionClosed(func(c *Client, err error) {})

  return server
}