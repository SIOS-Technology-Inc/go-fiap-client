package fiap

import (
	"testing"
	"time"
	"errors"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/testutil"
)

func TestFiapFetchResponseOnlyPoint(t *testing.T) {
	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomBodyResponder(`
		<body>
			<point id="http://xxxxxxxx/tokyo/building1/Room101/">
				<value time="2012-02-02T16:34:05.000+09:00">30</value>
			</point>
			<point id="http://xxxxxxxx/tokyo/building1/Room102/">
				<value time="2012-02-02T16:34:05.000+09:00">20</value>
				<value time="2012-02-02T17:34:05.000+09:00">25</value>
			</point>
		</body>
	`)

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	httpResponse, QueryRS, err := fiapFetch(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
			{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
		},
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, "200", httpResponse.Status)
	assert.Equal(t, "http://xxxxxxxx/tokyo/building1/Room101/", QueryRS.Transport.Body.Point[0].Id)
	assert.Equal(t, time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)), QueryRS.Transport.Body.Point[0].Value[0].Time)
	assert.Equal(t, "30", QueryRS.Transport.Body.Point[0].Value[0].Value)
	assert.Equal(t, "http://xxxxxxxx/tokyo/building1/Room102/", QueryRS.Transport.Body.Point[1].Id)
	assert.Equal(t, time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)), QueryRS.Transport.Body.Point[1].Value[0].Time)
	assert.Equal(t, "20", QueryRS.Transport.Body.Point[1].Value[0].Value)
	assert.Equal(t, time.Date(2012, 2, 2, 17, 34, 5, 0, time.FixedZone("", 9*60*60)), QueryRS.Transport.Body.Point[1].Value[1].Time)
	assert.Equal(t, "25", QueryRS.Transport.Body.Point[1].Value[1].Value)
	assert.Empty(t, QueryRS.Transport.Body.PointSet)
}

func TestFiapFetchResponseOnlyOnePointSet(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
			name       string
			responder  httpmock.Responder
			wantStatus string
			wantRootId string
			wantPointIds    []string
			wantPoinSetIds    []string
	}{
			{
					name: "when pointSet contain 0 pointId and 0 pointSetId",
					responder: testutil.CustomBodyResponder(
						`<body>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
							</pointSet>
						</body>`,
					), // 0個のpointIdを返すレスポンス
					wantStatus: "200",
					wantRootId: "http://xxxxxxxx/tokyo/building1/",
					wantPointIds:    []string{},
					wantPoinSetIds:    []string{},
			},
			{
					name: "when pointSet contain 1 pointId and 0 pointSetId",
					responder: testutil.CustomBodyResponder(
						`<body>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
								<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
							</pointSet>
						</body>`,
					), // 1個のpointIdを返すレスポンス
					wantStatus: "200",
					wantRootId: "http://xxxxxxxx/tokyo/building1/",
					wantPointIds:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
					wantPoinSetIds:    []string{},
			},
			{
				name: "when pointSet contain 2 pointId and 0 pointSetId",
				responder: testutil.CustomBodyResponder(
					`<body>
						<pointSet id="http://xxxxxxxx/tokyo/building1/">
							<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
							<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
						</pointSet>
					</body>`,
				), // 2個のpointIdを返すレスポンス
				wantStatus: "200",
				wantRootId: "http://xxxxxxxx/tokyo/building1/",
				wantPointIds:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/", "http://xxxxxxxx/tokyo/building1/Humidity/"},
				wantPoinSetIds:    []string{},
			},
			{
				name: "when pointSet contain 0 pointId and 1 pointSetId",
				responder: testutil.CustomBodyResponder(
				`<body>
					<pointSet id="http://xxxxxxxx/tokyo/building1/">
						<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
					</pointSet>
				</body>`,
				), // 1個のpointSetIdを返すレスポンス
				wantStatus: "200",
				wantRootId: "http://xxxxxxxx/tokyo/building1/",
				wantPointIds:    []string{},
				wantPoinSetIds:    []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			},
			{
				name: "when pointSet contain 0 pointId and 2 pointSetId",
				responder: testutil.CustomBodyResponder(
					`<body>
						<pointSet id="http://xxxxxxxx/tokyo/building1/">
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
						</pointSet>
					</body>`,
				), // 2個のpointSetIdを返すレスポンス
				wantStatus: "200",
				wantRootId: "http://xxxxxxxx/tokyo/building1/",
				wantPointIds:    []string{},
				wantPoinSetIds:    []string{"http://xxxxxxxx/tokyo/building1/Room101/", "http://xxxxxxxx/tokyo/building1/Room102/"},
			},
			{
				name: "when pointSet contain 2 pointId and 2 pointSetId",
				responder: testutil.CustomBodyResponder(
					`<body>
						<pointSet id="http://xxxxxxxx/tokyo/building1/">
							<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
							<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
						</pointSet>
					</body>`,
				), // 2個のpointIdと2個のpointSetIdを返すレスポンス
				wantStatus: "200",
				wantRootId: "http://xxxxxxxx/tokyo/building1/",
				wantPointIds:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/", "http://xxxxxxxx/tokyo/building1/Humidity/"},
				wantPoinSetIds:    []string{"http://xxxxxxxx/tokyo/building1/Room101/", "http://xxxxxxxx/tokyo/building1/Room102/"},
			},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// mockの有効化
					httpmock.Activate()
					defer httpmock.DeactivateAndReset()

					// 下記URLにPOSTしたときの挙動を定義
					httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", tc.responder)

					// テスト対象の関数を実行
					httpResponse, QueryRS, err := fiapFetch(
							"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							[]model.UserInputKey{
									{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
							},
							nil,
					)

					assert.NoError(t, err)
					assert.Equal(t, tc.wantStatus, httpResponse.Status)
					assert.Equal(t, tc.wantRootId, QueryRS.Transport.Body.PointSet[0].Id)
					if len(QueryRS.Transport.Body.PointSet[0].PointSetId) != 0 {
						for i, id := range tc.wantPoinSetIds {
								assert.Equal(t, id, *QueryRS.Transport.Body.PointSet[0].PointSetId[i])
						}
					}
					if len(QueryRS.Transport.Body.PointSet[0].PointId) != 0 {
						for i, id := range tc.wantPointIds {
								assert.Equal(t, id, *QueryRS.Transport.Body.PointSet[0].PointId[i])
						}
					}
			})
	}
}

func TestFiapFetchResponseOnlyMultiplePointSet(t *testing.T){
	// テストケースを定義
	testCases := []struct {
			name       string
			responder  httpmock.Responder
			wantStatus string
			wantRootId []string
	}{
			{
					name: "when body contain no pointSet",
					responder: testutil.CustomBodyResponder(
						`<body>
						</body>`,
					), // 0個のpointSetを返すレスポンス
					wantStatus: "200",
					wantRootId: []string{},
			},
			{
					name: "when body contain 1 pointSet",
					responder: testutil.CustomBodyResponder(
						`<body>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
								<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
								<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
								<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
								<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
							</pointSet>
						</body>`,
					), // 1個のpointSetを返すレスポンス
					wantStatus: "200",
					wantRootId: []string{"http://xxxxxxxx/tokyo/building1/"},
			},
			{
				name: "when body contain 2 pointSet",
				responder: testutil.CustomBodyResponder(
					`<body>
						<pointSet id="http://xxxxxxxx/tokyo/building1/">
							<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
							<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
						</pointSet>
						<pointSet id="http://xxxxxxxx/tokyo/building2/">
							<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
							<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
							<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
						</pointSet>
					</body>`,
				), // 2個のpointSetを返すレスポンス
				wantStatus: "200",
				wantRootId: []string{"http://xxxxxxxx/tokyo/building1/", "http://xxxxxxxx/tokyo/building2/"},
			},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
				// mockの有効化
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				// 下記URLにPOSTしたときの挙動を定義
				httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", tc.responder)

				// テスト対象の関数を実行
				httpResponse, QueryRS, err := fiapFetch(
						"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
						[]model.UserInputKey{
								{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
						},
						nil,
				)

				assert.NoError(t, err)
				assert.Equal(t, tc.wantStatus, httpResponse.Status)
				if len(QueryRS.Transport.Body.PointSet) != 0 {
					for i, id := range tc.wantRootId {
							assert.Equal(t, id, QueryRS.Transport.Body.PointSet[i].Id)
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Temperature/", *QueryRS.Transport.Body.PointSet[i].PointId[0])
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Humidity/", *QueryRS.Transport.Body.PointSet[i].PointId[1])
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Room101/", *QueryRS.Transport.Body.PointSet[i].PointSetId[0])
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Room102/", *QueryRS.Transport.Body.PointSet[i].PointSetId[1])
					}
				} else {
					assert.Empty(t, QueryRS.Transport.Body.PointSet)
				}
		})
	}
}

func TestInputErrors(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
			name       string
			connectionURL string
			keys       []model.UserInputKey
			wantError  string
	}{
			{
					name: "when connectionURL is empty",
					connectionURL: "",
					keys: []model.UserInputKey{
							{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
					},
					wantError: "connectionURL is empty",
			},
			{
					name: "when connectionURL is invalid",
					connectionURL: "htt://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
					keys: []model.UserInputKey{
							{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
					},
					wantError: "invalid connectionURL",
			},
			{
					name: "when keys is empty",
					connectionURL: "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
					keys: []model.UserInputKey{},
					wantError: "keys is empty",
			},
			{
					name: "when keys.ID is empty",
					connectionURL: "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
					keys: []model.UserInputKey{
							{ID: ""},
					},
					wantError: "keys.ID is empty",
			},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// テスト対象の関数を実行
					httpResponse, QueryRS, err := fiapFetch(
							tc.connectionURL,
							tc.keys,
							nil,
					)

					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.wantError)
					assert.Nil(t, httpResponse)
					assert.Nil(t, QueryRS)
			})
	}
}

func TestRequestFailure(t *testing.T) {
	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		httpmock.NewErrorResponder(errors.New("mocked error")))

	// テスト対象の関数を実行
	httpResponse, QueryRS, err := fiapFetch(
			"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
			[]model.UserInputKey{
					{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
			},
			nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client.Call error")
	assert.Nil(t, httpResponse)
	assert.Nil(t, QueryRS)
}
