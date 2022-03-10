package sterling_go

import (
	"fmt"
	"net/http"
	"testing"
)

var spayT = &SPay{
	key:     "000000010000001000000101000001010000011100001011010011010001000100010010000100010000110100001001000001110000001000000100000010000000000100001100000000110000010100000111000010110000110100011011",
	cypher:  "0000000100000010000000110010010100000111000010110000110100010001",
	baseUrl: baseUrl,
	appId:   16227,
	client:  http.DefaultClient,
}

func TestSPay_ListBanks(t *testing.T) {

	tests := []struct {
		name    string
		ref     *int64
		wantErr bool
	}{
		{
			"Test 1",
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotBanks, err := spayT.ListBanks(tt.ref)
			if err != nil {
				fmt.Printf("banks returned: %+v \n", gotBanks)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBanks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
