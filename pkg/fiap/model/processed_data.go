package model

/*
ProcessedPointSet is a type for the processed point set.
ProcessedPointSet は処理されたポイントセットの型です。

この型は、Fetcherメソッドの戻り値としてpointSetを返す際に使用します。
1つのPointSetIDに対応するPointSetIDとPointIDの配列の組を表します。
*/
type ProcessedPointSet struct {
	PointSetID []string `json:"point_set_id"`
	PointID    []string `json:"point_id"`
}
