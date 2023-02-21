package omailer

import (
	"context"

	"gopkg.in/gomail.v2"
)

// A Dialer is a dialer to an SMTP server.
type Dialer struct {
	// Host represents the host of the SMTP server.
	Host string
	// Port represents the port of the SMTP server.
	Port int
	// Username is the username to use to authenticate to the SMTP server.
	Username string
	// Password is the password to use to authenticate to the SMTP server.
	Password string

	statsC chan interface{}
	stopCh chan struct{}
}

// NewDialer returns a new SMTP Dialer. The given parameters are used to connect
// to the SMTP server.
func NewDialer(host string, port int, username, password string) *Dialer {
	return &Dialer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		statsC:   make(chan interface{}),
		stopCh:   make(chan struct{}),
	}
}

func (d *Dialer) StatsC() <-chan interface{} {
	return d.statsC
}

func (d *Dialer) Stop() {
	close(d.stopCh)
}

type EmailMessage struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Send opens an async connection to the SMTP server, sends the given emails and
// closes the connection.
func (d *Dialer) SendAsync(ctx context.Context, msg *EmailMessage) {
	go func() {
		err := d.Send(ctx, msg)
		if err == nil {
			return
		}
		select {
		case d.statsC <- err:
			return
		default:
		}
	}()
}

// Send opens a connection to the SMTP server, sends the given emails and
// closes the connection.
func (d *Dialer) Send(ctx context.Context, msg *EmailMessage) *EmailError {
	if ctx == nil {
		ctx = context.Background()
	}
	res := make(chan *EmailError, 1)

	go func() {
		defer close(res)
		res <- d.send(msg)
	}()

	select {
	case <-ctx.Done():
		return newEmailError(ctx.Err())
	case r := <-res:
		if r != nil {
			return r
		}
	}
	return nil
}

func (d *Dialer) send(msg *EmailMessage) *EmailError {
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)
	n := gomail.NewDialer(d.Host, d.Port, d.Username, d.Password)
	return newEmailError(n.DialAndSend(m))
}

type EmailError struct {
	err error
}

func (e *EmailError) Error() string {
	return e.err.Error()
}

func newEmailError(err error) *EmailError {
	return &EmailError{err: err}
}
