package requests

const (
	SaveURLReq   = "INSERT INTO url (url, alias) VALUES($1, $2)"
	DeleteURLReq = "DELETE * FROM url WHERE alias = $1"
	GetURLReq    = "SELECT url FROM url WHERE alias = $1"
)
