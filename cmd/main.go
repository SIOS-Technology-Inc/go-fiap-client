package main

import (
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func main() {
	
	// FetchRawOnce関数の呼び出し
	key := []model.UserInputKey{
		{
			ID:              "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
			MinMaxIndicator: "maximum",
		},
	}

	acceptableSize := 1000
	
	fiap.FetchOnce(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		key,
		&model.FetchOnceOption{
			AcceptableSize: &acceptableSize,
		})
}
