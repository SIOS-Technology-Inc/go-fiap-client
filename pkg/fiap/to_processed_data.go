package fiap

import (
	"fmt"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
)

func RawQueryRSToProcessedDatas(data *QueryRS) (pointSets map[string]([]model.ProcessedPointSet), points map[string]([]model.ProcessedPoint), cursor string, err error){
	if data == nil {
		return nil, nil, "", fmt.Errorf("QueryRS is nil")
	}
	
	// BodyにPointSetが返っていれば、それを処理する
	if data.Transport.Body.PointSet != nil {
		// PointSetの数だけ処理を繰り返し、
		// PointSetのIDをキーにしたmapを作成する
		for _, ps := range data.Transport.Body.PointSet {
			// ここでPointSetのIDを取得
			pointSetID := string(ps.Id)
			// ここでPointSetのIDをもとに、出力するデータの構造体を作成
			pointSets[pointSetID] = make([]model.ProcessedPointSet, 0)
			// PointSetの中にPointSetがあれば1階層分処理する
			if ps.PointSet != nil {
				for _, p := range ps.PointSet {
					// ここでPointSetのIDを取得
					pointSetID := string(p.Id)
					// ここでPointSetsの配列内にPointSetのIDを追加
					pointSets[pointSetID] = append(pointSets[pointSetID], model.ProcessedPointSet{PointSetID: []string{pointSetID}})
				}
			}
			// PointSetの中にPointがあれば1階層分処理する
			if ps.Point != nil {
				for _, p := range ps.Point {
					// ここでPointのIDを取得
					pointID := string(p.Id)
					// ここでmap[string]([]model.ProcessedPointSet)の配列内にPointのIDを追加
					pointSets[pointSetID] = append(pointSets[pointSetID], model.ProcessedPointSet{PointID: []string{pointID}})
				}
			}
		}
	}
	// BodyにPointが返っていれば、それを処理する
	if data.Transport.Body.Point != nil {
		
	return
}