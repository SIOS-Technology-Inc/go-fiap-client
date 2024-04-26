package fiap

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
	"regexp"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/testutil"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
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
								assert.Equal(t, id, QueryRS.Transport.Body.PointSet[0].PointSetId[i])
						}
					}
					if len(QueryRS.Transport.Body.PointSet[0].PointId) != 0 {
						for i, id := range tc.wantPointIds {
								assert.Equal(t, id, QueryRS.Transport.Body.PointSet[0].PointId[i])
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
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Temperature/", QueryRS.Transport.Body.PointSet[i].PointId[0])
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Humidity/", QueryRS.Transport.Body.PointSet[i].PointId[1])
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Room101/", QueryRS.Transport.Body.PointSet[i].PointSetId[0])
							assert.Equal(t,"http://xxxxxxxx/tokyo/building1/Room102/", QueryRS.Transport.Body.PointSet[i].PointSetId[1])
					}
				} else {
					assert.Empty(t, QueryRS.Transport.Body.PointSet)
				}
		})
	}
}

func TestFiapFetchInputErrors(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
			name       string
			connectionURL string
			keys       []model.UserInputKey
			wantError  string
	}{
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

func TestFiapFetchRequestFailure(t *testing.T) {
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


type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		QueryRQ *model.QueryRQ `xml:"queryRQ"`
	} `xml:"Body"`
}

func TestFiapFetchRequestOptionAndKeys(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
		name string
		option *model.FetchOnceOption
		keys []model.UserInputKey
		expectedRequestOption *model.FetchOnceOption
		expectedRequestKeys []model.Key
	}{
		{
				name: "when all inputs are valid",
				option: &model.FetchOnceOption{
					AcceptableSize: 1000,
					Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
				},
				keys: []model.UserInputKey{
						{
							ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
							Eq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
							Neq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
							Lt: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
							Gt: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
							Lteq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
							Gteq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
							MinMaxIndicator: model.SelectTypeMaximum,
						},
						{
							ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
						},
				},
				expectedRequestOption: &model.FetchOnceOption{
					AcceptableSize: 1000,
					Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
				},
				expectedRequestKeys: []model.Key{
					{
						Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
						Eq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
						Neq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
						Lt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
						Gt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
						Lteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
						Gteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
						Select: "maximum",
					},
					{
						Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					},
				},
		},
		{
			name: "when option.AcceptableSize is set but option.Cursor is not set",
			keys: []model.UserInputKey{
					{
						ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					},
			},
			option: &model.FetchOnceOption{
					AcceptableSize: 1000,
			},
			expectedRequestOption: &model.FetchOnceOption{
				AcceptableSize: 1000,
			},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
				},
			},
		},
		{
			name: "when option.Cursor is set but option.AcceptableSize is not set",
			keys: []model.UserInputKey{
					{
						ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					},
			},
			option: &model.FetchOnceOption{
				Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
			},
			expectedRequestOption: &model.FetchOnceOption{
				Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
			},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
				},
			},
		},
		{
			name: "when option.AcceptableSize and option.Cursor are set",
			keys: []model.UserInputKey{
					{
						ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					},
			},
			option: &model.FetchOnceOption{
				AcceptableSize: 1000,
				Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
			},
			expectedRequestOption: &model.FetchOnceOption{
				AcceptableSize: 1000,
				Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
			},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
				},
			},
		},
		{
			name : "when UserInputKey.MinMaxIndicator is SelectTypeMaximum",
			keys: []model.UserInputKey{
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					MinMaxIndicator: model.SelectTypeMaximum,
				},
			},
			option: nil,
			expectedRequestOption: &model.FetchOnceOption{},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					Select: "maximum",
				},
			},
		},
		{
			name : "when UserInputKey.MinMaxIndicator is SelectTypeMinimum",
			keys: []model.UserInputKey{
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					MinMaxIndicator: model.SelectTypeMinimum,
				},
			},
			option: nil,
			expectedRequestOption: &model.FetchOnceOption{},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					Select: "minimum",
				},
			},
		},
		{
			name : "when keys contain one UserInputKey",
			keys: []model.UserInputKey{
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					Eq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Neq: testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Lt: testutil.TimeToTimep(time.Date(2021, 1, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Gt: testutil.TimeToTimep(time.Date(2021, 1, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Lteq: testutil.TimeToTimep(time.Date(2021, 1, 5, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Gteq: testutil.TimeToTimep(time.Date(2021, 1, 6, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					MinMaxIndicator: model.SelectTypeMaximum,
				},
			},
			option: nil,
			expectedRequestOption: &model.FetchOnceOption{},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					Eq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Neq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Lt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Gt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Lteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 5, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Gteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 6, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Select: "maximum",
				},
			},
		},
		{
			name : "when keys contain two UserInputKey",
			keys: []model.UserInputKey{
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					Eq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Neq: testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Lt: testutil.TimeToTimep(time.Date(2021, 1, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Gt: testutil.TimeToTimep(time.Date(2021, 1, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Lteq: testutil.TimeToTimep(time.Date(2021, 1, 5, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Gteq: testutil.TimeToTimep(time.Date(2021, 1, 6, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					MinMaxIndicator: model.SelectTypeMaximum,
				},
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Humidity_TD",
					Eq: testutil.TimeToTimep(time.Date(2021, 1, 6, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Neq: testutil.TimeToTimep(time.Date(2021, 1, 5, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Lt: testutil.TimeToTimep(time.Date(2021, 1, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Gt: testutil.TimeToTimep(time.Date(2021, 1, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Lteq: testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					Gteq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
					MinMaxIndicator: model.SelectTypeMinimum,
				},
			},
			option: nil,
			expectedRequestOption: &model.FetchOnceOption{},
			expectedRequestKeys: []model.Key{
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					Eq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Neq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Lt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Gt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Lteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 5, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Gteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 6, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Select: "maximum",
				},
				{
					Id: "http://kurimoto/nukaya/vaisala/B-2/Humidity_TD",
					Eq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 6, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Neq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 5, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Lt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 4, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Gt: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 3, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Lteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Gteq: tools.TimeToString(testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))),
					Select: "minimum",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// httpmockの有効化
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Matcher function
			matcher := httpmock.NewMatcher("",func(req *http.Request) bool {
				envelope := &Envelope{}
				err := xml.NewDecoder(req.Body).Decode(envelope)
				if err != nil {
						return false
				}
				if envelope.Body.QueryRQ.Transport.Header.Query.AcceptableSize == tc.expectedRequestOption.AcceptableSize &&
					envelope.Body.QueryRQ.Transport.Header.Query.Cursor == tc.expectedRequestOption.Cursor &&
					len(envelope.Body.QueryRQ.Transport.Header.Query.Key) == len(tc.expectedRequestKeys) {
						for i := range envelope.Body.QueryRQ.Transport.Header.Query.Key {
									if envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Id == tc.expectedRequestKeys[i].Id && 
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].AttrName == "time" &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Eq == tc.expectedRequestKeys[i].Eq &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Neq == tc.expectedRequestKeys[i].Neq &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Lt == tc.expectedRequestKeys[i].Lt &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Gt == tc.expectedRequestKeys[i].Gt &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Lteq == tc.expectedRequestKeys[i].Lteq &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Gteq == tc.expectedRequestKeys[i].Gteq &&
								 envelope.Body.QueryRQ.Transport.Header.Query.Key[i].Select == tc.expectedRequestKeys[i].Select {
								return true
							}
						}
					return false
				}
				return false
			})

			// Responder function
			responder := testutil.CustomBodyResponder("")

			// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
			httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
				matcher,
				responder,
			)

			// テスト対象の関数を実行
			httpResponse, _, err := fiapFetch(
				"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
				tc.keys,
				tc.option,
			)
			if err != nil {
				t.Fatalf("fiapFetch failed: %v", err)
			}
			if httpResponse == nil {
				t.Fatal("httpResponse is nil")
			}
			assert.Equal(t, 200, httpResponse.StatusCode)
		})
	}
}

func TestFiapFetchRequestWithoutSpecificString(t *testing.T) {
	testCases := []struct {
			name string
			keys []model.UserInputKey
			option *model.FetchOnceOption
			notExpectedString []string  // 含まれてほしくない文字列を追加
	}{
			{
					name: "when option is nil",
					keys: []model.UserInputKey{
							{
								ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
							},
					},
					option: nil,
					notExpectedString: []string{"acceptableSize", "cursor"},
			},
			{
					name: "when option.AcceptableSize is set but option.Cursor is not set",
					keys: []model.UserInputKey{
							{
								ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
							},
					},
					option: &model.FetchOnceOption{
							AcceptableSize: 1000,
					},
					notExpectedString: []string{"cursor"},
			},
			{
				name: "when option.Cursor is set but option.AcceptableSize is not set",
				keys: []model.UserInputKey{
						{
							ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
						},
				},
				option: &model.FetchOnceOption{
					Cursor: "9c383255-d5ae-4e86-9ce6-b72800ee65f7",
				},
				notExpectedString: []string{"acceptableSize"},
		},
		{
			name: "when UserInputKey.MinMaxIndicator is SelectTypeNone",
			keys: []model.UserInputKey{
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
					MinMaxIndicator: model.SelectTypeNone,
				},
			},
			option: nil,
			notExpectedString: []string{"select"},
		},
		{
			name: "when UserInputKey.MinMaxIndicator is not set",
			keys: []model.UserInputKey{
				{
					ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD",
				},
			},
			option: nil,
			notExpectedString: []string{"select"},
		},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// httpmockの有効化
					httpmock.Activate()
					defer httpmock.DeactivateAndReset()

					// Matcher function
					matcher := httpmock.NewMatcher("", func(req *http.Request) bool {
						bodyBytes, err := io.ReadAll(req.Body)
						if err != nil {
							return false
						}
						bodyString := string(bodyBytes)

						for _, tag := range tc.notExpectedString {
							if strings.Contains(bodyString, tag) {
								return false
							}
						}
						return true
					})

					// Responder function
					responder := testutil.CustomBodyResponder("")

					// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
					httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							matcher,
							responder,
					)

					// テスト対象の関数を実行
					httpResponse, _, err := fiapFetch(
							"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							tc.keys,
							tc.option,
					)
					if err != nil {
							t.Fatalf("fiapFetch failed: %v", err)
					}
					if httpResponse == nil {
							t.Fatal("httpResponse is nil")
					}
					assert.Equal(t, 200, httpResponse.StatusCode)
			})
	}
}

func TestFiapFetchRequestQueryType(t *testing.T) {
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Matcher function
	matcher := httpmock.NewMatcher("",func(req *http.Request) bool {
		envelope := &Envelope{}
		err := xml.NewDecoder(req.Body).Decode(envelope)
		if err != nil {
				return false
		}
		if envelope.Body.QueryRQ.Transport.Header.Query.Type == "storage" {	
			return true
		}
		return false
	})

	// Responder function
	responder := testutil.CustomBodyResponder("")

	// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
	httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		matcher,
		responder,
	)

	// テスト対象の関数を実行
	httpResponse, _, err := fiapFetch(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
				{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
		},
		nil,
)
if err != nil {
	t.Fatalf("fiapFetch failed: %v", err)
}
if httpResponse == nil {
	t.Fatal("httpResponse is nil")
}
	assert.Equal(t, 200, httpResponse.StatusCode)
}

func TestFiapFetchRequestUuid(t *testing.T) {
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Matcher function
	matcher := httpmock.NewMatcher("",func(req *http.Request) bool {
		envelope := &Envelope{}
		err := xml.NewDecoder(req.Body).Decode(envelope)
		if err != nil {
				return false
		}

		uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
		return uuidRegex.MatchString(envelope.Body.QueryRQ.Transport.Header.Query.Id)
	})

	// Responder function
	responder := testutil.CustomBodyResponder("")

	// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
	httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		matcher,
		responder,
	)

	// テスト対象の関数を実行
	httpResponse, _, err := fiapFetch(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
				{ID: "http://kurimoto/nukaya/vaisala/B-2/Temperature_TD"},
		},
		nil,
)
if err != nil {
	t.Fatalf("fiapFetch failed: %v", err)
}
if httpResponse == nil {
	t.Fatal("httpResponse is nil")
}
	assert.Equal(t, 200, httpResponse.StatusCode)
}