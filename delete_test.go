package boltrepo

import "testing"

func TestDelete(t *testing.T) {
	br := getTestRepo()
	defer br.DB.Close()

	t.Run(``, func(t *testing.T) {
	})
}
