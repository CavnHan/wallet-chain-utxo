package bitcoin

import (
	"fmt"
	"testing"
)

func TestGetLatestBlick(t *testing.T) {
	client, err := NewBtcClient("nd-314-634-477.p2pify.com", "lucid-swanson", "salon-ahead-vanish-dial-curdle-arise")
	if err != nil {
		t.Error(err)
	}
	block := client.GetLatestBlock()
	fmt.Println(block)

}
