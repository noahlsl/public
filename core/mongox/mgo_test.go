package mongox

import "testing"

func TestMgo(t *testing.T) {
	cfg := &Cfg{
		Host:     "127.0.0.1",
		Port:     27017,
		Username: "msg",
		Password: "3QPEtoht7yXYRavp",
		Database: "msg",
	}
	_ = cfg.NewDatabase()
}
