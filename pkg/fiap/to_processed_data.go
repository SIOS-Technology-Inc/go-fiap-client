package fiap

import (
	"fmt"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func QueryRSToProcessedDatas(data *QueryRS) (pointSets map[string](model.ProcessedPointSet), points map[string](model.ProcessedPoint), cursor string, err error){
	if data == nil {
		return nil, nil, "", fmt.Errorf("QueryRS is nil")
	}	
	// BodyにPointSetが返っていれば、それを処理する
	if data.Transport.Body.PointSet != nil {
		// PointSetIdが同じものがあれば、データを結合する
	}
	// BodyにPointが返っていれば、それを処理する
	if data.Transport.Body.Point != nil {
		// PointIdが同じものがあれば、データを結合する
	}
		
	return
}