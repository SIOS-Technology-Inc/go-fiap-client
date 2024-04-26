package fiap

import (
	"testing"
	"time"
	"errors"
	"net/http"
	"io"
	"strings"


	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/testutil"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/model"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap/tools"
)

func TestFetchOncePointBoundary(t *testing.T){
	// テストケースを定義
	testCases := []struct {
		name            string
		body            string
		expectedPointCount	int
		expectedPoints  map[string][]model.Value
	}{
		{
			name: "0 points",
			body: `
			<body>
			</body>
			`, // 0個のpointを含むレスポンス
			expectedPointCount: 0,
			expectedPoints: map[string][]model.Value{},
		},
		{
			name: "1 point",
			body: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
				</point>
			</body>
			`, // 1個のpointを含むレスポンス
			expectedPointCount: 1,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {},
			},
		},
		{
			name: "2 points",
			body: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
				</point>
				<point id="http://xxxxxxxx/tokyo/building1/Room102/">
				</point>
			</body>
			`, // 2個のpointを含むレスポンス
			expectedPointCount: 2,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {},
				"http://xxxxxxxx/tokyo/building1/Room102/": {},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mockの有効化
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// 下記URLにPOSTしたときの挙動を定義
			responder := testutil.CustomBodyResponder(tc.body)
			httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

			// テスト対象の関数を実行
			_, points, _, _, err := FetchOnce(
					"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
					[]model.UserInputKey{
							{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
					},
					&model.FetchOnceOption{},
			)

			assert.NoError(t, err)
			assert.Len(t, points, tc.expectedPointCount)
		})
	}
}

func TestFetchOncePointValueBoundary(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
		name            string
		body            string
		expectedValueCount       int
		expectedPoints  map[string][]model.Value
	}{
		{
			name: "when response contain no Point.Value",
			body: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/"></point>
			</body>
			`, // pointに0個のvalueを含むレスポンス
			expectedValueCount:       0,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {},
			},
		},
		{
			name: "when response contain 1 Point.Value",
			body: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
					<value time="2012-02-02T16:34:05.000+09:00">30</value>
				</point>
			</body>
			`, // pointに1個のvalueを含むレスポンス
			expectedValueCount: 1,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					{
						Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "30",
					},
				},
			},
		},
		{
			name: "2 points",
			body: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
					<value time="2012-02-02T16:34:05.000+09:00">30</value>
					<value time="2012-02-02T16:35:05.000+09:00">40</value>
				</point>
			</body>
			`, // pointに2個のvalueを含むレスポンス
			expectedValueCount: 2,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					{
						Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "30",
					},
					{
						Time: time.Date(2012, 2, 2, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "40",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// mockの有効化
					httpmock.Activate()
					defer httpmock.DeactivateAndReset()

					// 下記URLにPOSTしたときの挙動を定義
					responder := testutil.CustomBodyResponder(tc.body)
					httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

					// テスト対象の関数を実行
					_, points, _, _, err := FetchOnce(
							"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							[]model.UserInputKey{
									{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
							},
							&model.FetchOnceOption{},
					)

					assert.NoError(t, err)
					assert.Len(t, points["http://xxxxxxxx/tokyo/building1/Room101/"], tc.expectedValueCount)
					assert.Equal(t, tc.expectedPoints, points)
			})
	}
}

func TestFetchOncePointSetBoundary(t *testing.T){
	// テストケースを定義
	testCases := []struct {
		name            string
		body            string
		expectedPointSetCount	int
		expectedPointSets  map[string](model.ProcessedPointSet)
	}{
		{
			name: "0 pointSets",
			body: `
			<body>
			</body>
			`, // 0個のpointSetを含むレスポンス
			expectedPointSetCount:       0,
			expectedPointSets: map[string](model.ProcessedPointSet)(nil),
		},
		{
			name: "1 pointSets",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
				</pointSet>
			</body>
			`, // 1個のpointSetを含むレスポンス
			expectedPointSetCount: 1,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/": {},
			},
		},
		{
			name: "2 pointSets",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/" />
				<pointSet id="http://xxxxxxxx/tokyo/building2/" />
			</body>
			`, // 2個のpointSetを含むレスポンス
			expectedPointSetCount: 2,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/": {},
				"http://xxxxxxxx/tokyo/building2/": {},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
				// mockの有効化
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				// 下記URLにPOSTしたときの挙動を定義
				responder := testutil.CustomBodyResponder(tc.body)
				httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

				// テスト対象の関数を実行
				pointSets, _, _, _, err := FetchOnce(
						"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
						[]model.UserInputKey{
								{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
						},
						&model.FetchOnceOption{},
				)

				assert.NoError(t, err)
				assert.Len(t, pointSets, tc.expectedPointSetCount)
				assert.Equal(t, tc.expectedPointSets, pointSets)
				
		})
	}
}

func TestFetchOncePointSetPointSetIDBoundary(t *testing.T){
	// テストケースを定義
	testCases := []struct {
		name            string
		body            string
		expectedPointSetPointSetIdCount	int
		expectedPointSets  map[string](model.ProcessedPointSet)
	}{
		{
			name: "0 PointSets.PointSetId",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
				</pointSet>
			</body>
			`, // 0個のpointSetIDを含むレスポンス
			expectedPointSetPointSetIdCount:       0,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/": {},
			},
		},
		{
			name: "1 pointSets.PointSetId",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
					<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
				</pointSet>
			</body>
			`, // 1個のpointSetIDを含むレスポンス
			expectedPointSetPointSetIdCount: 1,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/": {
					PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
				},
			},
		},
		{
			name: "2 pointSets.PointSetId",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
					<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
					<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
				</pointSet>
			</body>
			`, // 2個のpointSetIDを含むレスポンス
			expectedPointSetPointSetIdCount: 2,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/": {
					PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/", "http://xxxxxxxx/tokyo/building1/Room102/"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
				// mockの有効化
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				// 下記URLにPOSTしたときの挙動を定義
				responder := testutil.CustomBodyResponder(tc.body)
				httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

				// テスト対象の関数を実行
				pointSets, _, _, _, err := FetchOnce(
						"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
						[]model.UserInputKey{
								{ID: "http://xxxxxxxx/tokyo/building1/"},
						},
						&model.FetchOnceOption{},
				)

				assert.NoError(t, err)
				assert.Len(t, pointSets["http://xxxxxxxx/tokyo/building1/"].PointSetID, tc.expectedPointSetPointSetIdCount)
				assert.Equal(t, tc.expectedPointSets, pointSets)
				
		})
	}
}

func TestFetchOncePointSetPointIdBoundary(t *testing.T){
	// テストケースを定義
	testCases := []struct {
		name            string
		body            string
		expectedPointSetPointIdCount	int
		expectedPointSets  map[string](model.ProcessedPointSet)
	}{
		{
			name: "0 PointSets.PointId",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/">
				</pointSet>
			</body>
			`, // 0個のpointIDを含むレスポンス
			expectedPointSetPointIdCount:       0,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/Room101/": {},
			},
		},
		{
			name: "1 pointSets.PointId",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/">
					<point id="http://xxxxxxxx/tokyo/building1/Room101/Temperature/" />
				</pointSet>
			</body>
			`, // 1個のpointIDを含むレスポンス
			expectedPointSetPointIdCount: 1,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					PointID: []string{"http://xxxxxxxx/tokyo/building1/Room101/Temperature/"},
				},
			},
		},
		{
			name: "2 pointSets.PointId",
			body: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/">
					<point id="http://xxxxxxxx/tokyo/building1/Room101/Temperature/" />
					<point id="http://xxxxxxxx/tokyo/building1/Room101/Humidity/" />
				</pointSet>
			</body>
			`, // 2個のpointIDを含むレスポンス
			expectedPointSetPointIdCount: 2,
			expectedPointSets: map[string](model.ProcessedPointSet){
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					PointID: []string{"http://xxxxxxxx/tokyo/building1/Room101/Temperature/", "http://xxxxxxxx/tokyo/building1/Room101/Humidity/"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
				// mockの有効化
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				// 下記URLにPOSTしたときの挙動を定義
				responder := testutil.CustomBodyResponder(tc.body)
				httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

				// テスト対象の関数を実行
				pointSets, _, _, _, err := FetchOnce(
						"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
						[]model.UserInputKey{
								{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
						},
						&model.FetchOnceOption{},
				)

				assert.NoError(t, err)
				assert.Len(t, pointSets["http://xxxxxxxx/tokyo/building1/Room101/"].PointID, tc.expectedPointSetPointIdCount)
				assert.Equal(t, tc.expectedPointSets, pointSets)
				
		})
	}
}

func TestFetchOnceBoundaryCombination(t *testing.T){
	expectedPoints := map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time: time.Date(2012, 2, 2, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "30",
			},
			{
				Time: time.Date(2012, 2, 3, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "40",
			},
		},
		"http://xxxxxxxx/tokyo/building1/Room102/": {
			{
				Time: time.Date(2012, 2, 4, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "50",
			},
			{
				Time: time.Date(2012, 2, 5, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "60",
			},
		},
	}
	expectedPointSets := map[string](model.ProcessedPointSet){
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/", "http://xxxxxxxx/tokyo/building1/Room102/"},
			PointID: []string{"http://xxxxxxxx/tokyo/building1/Temperature/", "http://xxxxxxxx/tokyo/building1/Humidity/"},
		},
		"http://xxxxxxxx/tokyo/building2/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building2/Room101/", "http://xxxxxxxx/tokyo/building2/Room102/"},
			PointID: []string{"http://xxxxxxxx/tokyo/building2/Temperature/", "http://xxxxxxxx/tokyo/building2/Humidity/"},
		},
	}

	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomBodyResponder(`
	<body>
		<point id="http://xxxxxxxx/tokyo/building1/Room101/">
			<value time="2012-02-02T16:35:05.000+09:00">30</value>
			<value time="2012-02-03T16:35:05.000+09:00">40</value>
		</point>
		<point id="http://xxxxxxxx/tokyo/building1/Room102/">
			<value time="2012-02-04T16:35:05.000+09:00">50</value>
			<value time="2012-02-05T16:35:05.000+09:00">60</value>
		</point>
		<pointSet id="http://xxxxxxxx/tokyo/building1/">
			<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
			<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
			<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
			<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
		</pointSet>
		<pointSet id="http://xxxxxxxx/tokyo/building2/">
			<point id="http://xxxxxxxx/tokyo/building2/Temperature/" />
			<point id="http://xxxxxxxx/tokyo/building2/Humidity/" />
			<pointSet id="http://xxxxxxxx/tokyo/building2/Room101/" />
			<pointSet id="http://xxxxxxxx/tokyo/building2/Room102/" />
		</pointSet>
	</body>
	`)

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	pointSets, points, _, _, err := FetchOnce(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
			{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
		},
		&model.FetchOnceOption{},
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedPointSets, pointSets)
	assert.Equal(t, expectedPoints, points)
}

func TestFetchOnceBodyIsEmpty(t *testing.T){
	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomBodyResponder(`
	<body>
	</body>
	`)

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	pointSets, points, _, _, err := FetchOnce(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
			{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
		},
		&model.FetchOnceOption{},
	)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(pointSets))
	assert.Equal(t, 0, len(points))
}

func TestFetchOnceWithRepeatedPointSetId(t *testing.T) {
	testCases := []struct {
			name          string
			responderBody string
			expectedPointSets map[string]model.ProcessedPointSet
		}{
					{
						name: "when same pointset id is not repeated",
						responderBody: `
						<body>
								<pointSet id="http://xxxxxxxx/tokyo/building1/">
										<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
										<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
								</pointSet>
								<pointSet id="http://xxxxxxxx/tokyo/building2/">
										<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
										<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
								</pointSet>
						</body>
						`,
						expectedPointSets: map[string]model.ProcessedPointSet{
								"http://xxxxxxxx/tokyo/building1/": {
										PointSetID: []string{
												"http://xxxxxxxx/tokyo/building1/Room101/",
										},
										PointID: []string{
												"http://xxxxxxxx/tokyo/building1/Temperature/",
										},
								},
								"http://xxxxxxxx/tokyo/building2/": {
									PointSetID: []string{
											"http://xxxxxxxx/tokyo/building1/Room102/",
									},
									PointID: []string{
											"http://xxxxxxxx/tokyo/building1/Humidity/",
									},
							},
						},
				},
				{
					name: "when same pointset id repeated 2 times",
					responderBody: `
					<body>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
									<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
							</pointSet>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
									<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
							</pointSet>
					</body>
					`,
					expectedPointSets: map[string]model.ProcessedPointSet{
							"http://xxxxxxxx/tokyo/building1/": {
									PointSetID: []string{
											"http://xxxxxxxx/tokyo/building1/Room101/",
											"http://xxxxxxxx/tokyo/building1/Room102/",
									},
									PointID: []string{
											"http://xxxxxxxx/tokyo/building1/Temperature/",
											"http://xxxxxxxx/tokyo/building1/Humidity/",
									},
							},
					},
			},
			{
					name: "when same pointset id repeated 3 times",
					responderBody: `
					<body>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
									<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
							</pointSet>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
									<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
							</pointSet>
							<pointSet id="http://xxxxxxxx/tokyo/building1/">
									<point id="http://xxxxxxxx/tokyo/building1/Illuminance/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room103/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room104/" />
							</pointSet>		
					</body>
					`,
					expectedPointSets: map[string]model.ProcessedPointSet{
							"http://xxxxxxxx/tokyo/building1/": {
									PointSetID: []string{
											"http://xxxxxxxx/tokyo/building1/Room101/",
											"http://xxxxxxxx/tokyo/building1/Room102/",
											"http://xxxxxxxx/tokyo/building1/Room103/",
											"http://xxxxxxxx/tokyo/building1/Room104/",
									},
									PointID: []string{
											"http://xxxxxxxx/tokyo/building1/Temperature/",
											"http://xxxxxxxx/tokyo/building1/Humidity/",
											"http://xxxxxxxx/tokyo/building1/Illuminance/",
									},
							},
					},
			},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// mockの有効化
					httpmock.Activate()
					defer httpmock.DeactivateAndReset()

					// 下記URLにPOSTしたときの挙動を定義
					responder := testutil.CustomBodyResponder(tc.responderBody)
					httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

					// テスト対象の関数を実行
					pointSets, _, _, _, err := FetchOnce(
							"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							[]model.UserInputKey{
									{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
							},
							&model.FetchOnceOption{},
					)

					assert.NoError(t, err)
					assert.Equal(t, tc.expectedPointSets, pointSets)
			})
	}
}

func TestFetchOnceWithRepeatedPointId(t *testing.T) {
	testCases := []struct {
			name          string
			responderBody string
			expectedPoints map[string][]model.Value
	}{
			{
				name: "when same point id is not repeated",
				responderBody: `
				<body>
						<point id="http://xxxxxxxx/tokyo/building1/Room101/">
								<value time="2012-02-02T16:35:05.000+09:00">30</value>
						</point>
						<point id="http://xxxxxxxx/tokyo/building1/Room102/">
								<value time="2012-02-04T16:35:05.000+09:00">40</value>
								<value time="2012-02-05T16:35:05.000+09:00">50</value>
						</point>
				</body>
				`,
				expectedPoints: map[string][]model.Value{
					"http://xxxxxxxx/tokyo/building1/Room101/": {
							{
									Time: time.Date(2012, 2, 2, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
									Value: "30",
							},
					},
					"http://xxxxxxxx/tokyo/building1/Room102/": {
							{
									Time: time.Date(2012, 2, 4, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
									Value: "40",
							},
							{
									Time: time.Date(2012, 2, 5, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
									Value: "50",
							},
					},
				},
			},
			{
					name: "when same point id repeated 2 times",
					responderBody: `
					<body>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-02T16:35:05.000+09:00">30</value>
									<value time="2012-02-03T16:35:05.000+09:00">40</value>
							</point>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-04T16:35:05.000+09:00">50</value>
									<value time="2012-02-05T16:35:05.000+09:00">60</value>
							</point>
					</body>
					`,
					expectedPoints: map[string][]model.Value{
							"http://xxxxxxxx/tokyo/building1/Room101/": {
									{
											Time: time.Date(2012, 2, 2, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "30",
									},
									{
											Time: time.Date(2012, 2, 3, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "40",
									},
									{
											Time: time.Date(2012, 2, 4, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "50",
									},
									{
											Time: time.Date(2012, 2, 5, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "60",
									},
							},
					},
			},
			{
					name: "when same point id repeated 3 times",
					responderBody: `
					<body>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-02T16:35:05.000+09:00">30</value>
									<value time="2012-02-03T16:35:05.000+09:00">40</value>
							</point>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-04T16:35:05.000+09:00">50</value>
							</point>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-05T16:35:05.000+09:00">60</value>
									<value time="2012-02-06T16:35:05.000+09:00">70</value>
							</point>
					</body>
					`,
					expectedPoints: map[string][]model.Value{
							"http://xxxxxxxx/tokyo/building1/Room101/": {
									{
											Time: time.Date(2012, 2, 2, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "30",
									},
									{
											Time: time.Date(2012, 2, 3, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "40",
									},
									{
											Time: time.Date(2012, 2, 4, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "50",
									},
									{
											Time: time.Date(2012, 2, 5, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "60",
									},
									{
											Time: time.Date(2012, 2, 6, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
											Value: "70",
									},
							},
					},
			},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// mockの有効化
					httpmock.Activate()
					defer httpmock.DeactivateAndReset()

					// 下記URLにPOSTしたときの挙動を定義
					responder := testutil.CustomBodyResponder(tc.responderBody)
					httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

					// テスト対象の関数を実行
					_, points, _, _, err := FetchOnce(
							"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							[]model.UserInputKey{
									{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
							},
							&model.FetchOnceOption{},
					)

					assert.NoError(t, err)
					assert.Equal(t, tc.expectedPoints, points)
			})
	}
}


func TestFetchOnceCursor(t *testing.T) {
	testCases := []struct {
			name          string
			responderBody string
			expectedCursor string
	}{
			{
					name: "when cursor is not empty",
					responderBody: `
					<header>
							<OK/>
							<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="2f9f9cc6-7530-3cc9-faee-e894edeb1566">
									<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
							</query>
					</header>
					<body>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-02T16:34:05.000+09:00">30</value>
							</point>
					</body>
					`,
					expectedCursor: "2f9f9cc6-7530-3cc9-faee-e894edeb1566",
			},
			{
					name: "when cursor is empty",
					responderBody: `
					<header>
							<OK/>
							<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="">
									<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
							</query>
					</header>
					<body>
							<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-02T16:34:05.000+09:00">30</value>
							</point>
					</body>
					`,
					expectedCursor: "",
			},
	}

	for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
					// mockの有効化
					httpmock.Activate()
					defer httpmock.DeactivateAndReset()

					// 下記URLにPOSTしたときの挙動を定義
					responder := testutil.CustomHeaderBodyResponder(tc.responderBody)
					httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

					// テスト対象の関数を実行
					_, _, cursor, _, err := FetchOnce(
							"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
							[]model.UserInputKey{
									{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
							},
							&model.FetchOnceOption{},
					)

					assert.NoError(t, err)
					assert.Equal(t, tc.expectedCursor, cursor)
			})
	}
}

func TestFetchOnceFiapFetchInputError(t *testing.T){
		testcases := []struct{
			name string
			connectionURL string
			keys []model.UserInputKey
			wantError []string
		}{
			{
				name: "when connectionURL is empty",
				connectionURL: "",
				keys: []model.UserInputKey{
					{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
				},
				wantError: []string{
					"fiapFetch error",
					"connectionURL is empty",
				},
			},
			{
				name: "when connection url is invalid",
				connectionURL: "htrp://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
				keys: []model.UserInputKey{
					{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
				},
				wantError: []string{
					"fiapFetch error",
					"invalid connectionURL",
				},
			},
			{
				name: "when keys is empty",
				connectionURL: "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
				keys: []model.UserInputKey{},
				wantError: []string{
					"fiapFetch error",
					"keys is empty",
				},
			},
			{
				name: "when keys id is empty",
				connectionURL: "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
				keys: []model.UserInputKey{
					{ID: ""},
				},
				wantError: []string{
					"fiapFetch error",
					"keys.ID is empty",
				},
			},
		};

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				_, _, _, _, err := FetchOnce(tc.connectionURL, tc.keys, &model.FetchOnceOption{})
				for _, want := range tc.wantError {
					assert.Contains(t, err.Error(), want)
				}
			})
		}
}

func TestFetchOnceFiapFetchRequestError(t *testing.T){
	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		httpmock.NewErrorResponder(errors.New("mocked error")))

	// テスト対象の関数を実行
	_, _, _, _, err := FetchOnce(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
			{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
		},
		&model.FetchOnceOption{},
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fiapFetch error")
	assert.Contains(t, err.Error(), "client.Call error")
}

func TestFetchOnceProcessQueryRSError(t *testing.T){
	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := httpmock.NewStringResponder(200, `
	<?xml version='1.0' encoding='utf-8'?>
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Header/>
			<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
					</ns2:queryRS>
			</soapenv:Body>
	</soapenv:Envelope>
	`)

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	_, _, _, _, err := FetchOnce(
		"http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
		[]model.UserInputKey{
			{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
		},
		&model.FetchOnceOption{},
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "processQueryRS error")
	assert.Contains(t, err.Error(), "transport is nil")
}

func TestFetchFetchOnce1(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomHeaderBodyResponder(`
	<header>
		<OK/>
		<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="">
			<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
		</query>
	</header>
	<body>
		<point id="http://xxxxxxxx/tokyo/building1/Room101/">
			<value time="2012-02-02T16:34:05.000+09:00">30</value>
		</point>
		<pointSet id="http://xxxxxxxx/tokyo/building1/">
			<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
			<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
			<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
		</pointSet>
	</body>
	`)
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	pointSets, points, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/", "http://xxxxxxxx/tokyo/building1/Humidity/"},
		},
	}, pointSets)
	assert.Equal(t, map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time:  time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "30",
			},
		},
	}, points)
}

func TestFetchFetchOnce2(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		func (req *http.Request) (*http.Response, error){
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			if strings.Contains(string(body), "cursor") {
				responseWithoutCursor := httpmock.NewStringResponse(200,`
					<?xml version='1.0' encoding='utf-8'?>
					<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
					<soapenv:Header/>
					<soapenv:Body>
						<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
							<transport xmlns="http://gutp.jp/fiap/2009/11/">
								<header>
								<OK/>
									<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage">
										<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
									</query>
								</header>
								<body>
									<point id="http://xxxxxxxx/tokyo/building1/Room102/">
										<value time="2012-02-03T16:34:05.000+09:00">40</value>
									</point>
									<pointSet id="http://xxxxxxxx/tokyo/building1/">
										<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
										<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
									</pointSet>
								</body>
								</transport>
						</ns2:queryRS>
					</soapenv:Body>
					</soapenv:Envelope>
				`)
				return responseWithoutCursor, nil
			} else {
				responseWithCursor := httpmock.NewStringResponse(200, `
				<?xml version='1.0' encoding='utf-8'?>
				<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
				<soapenv:Header/>
				<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
						<transport xmlns="http://gutp.jp/fiap/2009/11/">
							<header>
							<OK/>
								<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb">
									<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
								</query>
							</header>
							<body>
								<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-02T16:34:05.000+09:00">30</value>
								</point>
								<pointSet id="http://xxxxxxxx/tokyo/building2/">
									<point id="http://xxxxxxxx/tokyo/building2/Temperature/" />
									<pointSet id="http://xxxxxxxx/tokyo/building2/Room101/" />
								</pointSet>
							</body>
							</transport>
					</ns2:queryRS>
				</soapenv:Body>
				</soapenv:Envelope>
				`)
				return responseWithCursor, nil
			}
		},
	)
	// テスト対象の関数を実行
	pointSets, points, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
		},
		"http://xxxxxxxx/tokyo/building2/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building2/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building2/Temperature/"},
		},
	}, pointSets)
	assert.Equal(t, map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time:  time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "30",
			},
		},
		"http://xxxxxxxx/tokyo/building1/Room102/": {
			{
				Time:  time.Date(2012, 2, 3, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "40",
			},
		},
	}, points)
}

func TestFetchFetchOncePointSetsBoundary(t *testing.T){
	testCases := []struct {
		name string
		responseBody string
		expectedPointSets map[string]model.ProcessedPointSet
	}{
		{
			name: "when fecthOncePointSets is empty",
			responseBody: `
			<body>
			</body>
			`,
			expectedPointSets: map[string]model.ProcessedPointSet{},
		},
		{
			name: "when fecthOncePointSets has 1 pointSet",
			responseBody: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
					<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
					<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
				</pointSet>
			</body>
			`,
			expectedPointSets: map[string]model.ProcessedPointSet{
				"http://xxxxxxxx/tokyo/building1/": {
					PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
					PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
				},
			},
		},
		{
			name: "when fecthOncePointSets has 2 pointSets",
			responseBody: `
			<body>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
					<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
					<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
				</pointSet>
				<pointSet id="http://xxxxxxxx/tokyo/building2/">
					<pointSet id="http://xxxxxxxx/tokyo/building2/Room101/" />
					<point id="http://xxxxxxxx/tokyo/building2/Temperature/" />
				</pointSet>
			</body>
			`,
			expectedPointSets: map[string]model.ProcessedPointSet{
				"http://xxxxxxxx/tokyo/building1/": {
					PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
					PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
				},
				"http://xxxxxxxx/tokyo/building2/": {
					PointSetID: []string{"http://xxxxxxxx/tokyo/building2/Room101/"},
					PointID:    []string{"http://xxxxxxxx/tokyo/building2/Temperature/"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mockの有効化
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// 下記URLにPOSTしたときの挙動を定義
			responder := testutil.CustomBodyResponder(tc.responseBody)
			httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

			// テスト対象の関数を実行
			pointSets, _, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
				{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
			}, &model.FetchOption{})
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPointSets, pointSets)
		})
	}
}

func TestFetchFetchOncePointsBoundary(t *testing.T){
	testCases := []struct {
		name string
		responseBody string
		expectedPoints map[string][]model.Value
	}{
		{
			name: "when fecthOncePoints is empty",
			responseBody: `
			<body>
			</body>
			`,
			expectedPoints: map[string][]model.Value{},
		},
		{
			name: "when fecthOncePoints has 1 points",
			responseBody: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
					<value time="2012-02-02T16:34:05.000+09:00">30</value>
				</point>
			</body>
			`,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					{
						Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "30",
					},
				},
			},
		},
		{
			name: "when fecthOncePoints has 2 points",
			responseBody: `
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
					<value time="2012-02-02T16:34:05.000+09:00">30</value>
				</point>
				<point id="http://xxxxxxxx/tokyo/building1/Room102/">
					<value time="2012-02-02T16:34:05.000+09:00">20</value>
					<value time="2012-02-02T17:34:05.000+09:00">25</value>
				</point>
			</body>
			`,
			expectedPoints: map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					{
						Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "30",
					},
				},
				"http://xxxxxxxx/tokyo/building1/Room102/": {
					{
						Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "20",
					},
					{
						Time: time.Date(2012, 2, 2, 17, 34, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "25",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mockの有効化
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// 下記URLにPOSTしたときの挙動を定義
			responder := testutil.CustomBodyResponder(tc.responseBody)
			httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

			// テスト対象の関数を実行
			_, points, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
				{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
			}, &model.FetchOption{})
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPoints, points)
		})
	}
}

func TestFetchNotRepeatedPointSetId(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomHeaderBodyResponder(`
	<header>
		<OK/>
		<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="">
			<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
		</query>
	</header>
	<body>
		<pointSet id="http://xxxxxxxx/tokyo/building1/">
			<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
			<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
		</pointSet>
	</body>
	`)
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	pointSets, _, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Humidity/"},
		},
	}, pointSets)
}

func TestFetchRepeatedPointSetId(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		func (req *http.Request) (*http.Response, error){
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			if strings.Contains(string(body), `cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb"`) {
				responseWithoutCursor := httpmock.NewStringResponse(200,`
					<?xml version='1.0' encoding='utf-8'?>
					<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
					<soapenv:Header/>
					<soapenv:Body>
						<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
							<transport xmlns="http://gutp.jp/fiap/2009/11/">
								<header>
								<OK/>
									<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage">
										<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
									</query>
								</header>
								<body>
									<pointSet id="http://xxxxxxxx/tokyo/building1/">
										<point id="http://xxxxxxxx/tokyo/building1/Humidity/" />
										<pointSet id="http://xxxxxxxx/tokyo/building1/Room102/" />
									</pointSet>
								</body>
								</transport>
						</ns2:queryRS>
					</soapenv:Body>
					</soapenv:Envelope>
				`)
				return responseWithoutCursor, nil
			} else {
				responseWithCursor := httpmock.NewStringResponse(200, `
				<?xml version='1.0' encoding='utf-8'?>
				<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
				<soapenv:Header/>
				<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
						<transport xmlns="http://gutp.jp/fiap/2009/11/">
							<header>
							<OK/>
								<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb">
									<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
								</query>
							</header>
							<body>
								<pointSet id="http://xxxxxxxx/tokyo/building1/">
									<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
									<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
								</pointSet>
							</body>
							</transport>
					</ns2:queryRS>
				</soapenv:Body>
				</soapenv:Envelope>
				`)
				return responseWithCursor, nil
			}
		},
	)
	// テスト対象の関数を実行
	pointSets, _, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/", "http://xxxxxxxx/tokyo/building1/Room102/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/", "http://xxxxxxxx/tokyo/building1/Humidity/"},
		},
	}, pointSets)
}

func TestFetchNotRepeatedPointId(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomHeaderBodyResponder(`
	<header>
		<OK/>
		<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="">
			<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
		</query>
	</header>
	<body>
		<point id="http://xxxxxxxx/tokyo/building1/Room101/">
			<value time="2012-02-02T16:34:05.000+09:00">30</value>
		</point>
	</body>
	`)
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	_, points, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "30",
			},
		},
	}, points)
}

func TestFetchRepeatedPointId(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		func (req *http.Request) (*http.Response, error){
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			if strings.Contains(string(body), `cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb"`) {
				responseWithoutCursor := httpmock.NewStringResponse(200,`
					<?xml version='1.0' encoding='utf-8'?>
					<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
					<soapenv:Header/>
					<soapenv:Body>
						<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
							<transport xmlns="http://gutp.jp/fiap/2009/11/">
								<header>
								<OK/>
									<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage">
										<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
									</query>
								</header>
								<body>
									<point id="http://xxxxxxxx/tokyo/building1/Room101/">
										<value time="2012-02-02T16:34:05.000+09:00">2</value>
										<value time="2012-02-03T16:34:05.000+09:00">3</value>
									</point>
								</body>
								</transport>
						</ns2:queryRS>
					</soapenv:Body>
					</soapenv:Envelope>
				`)
				return responseWithoutCursor, nil
			} else {
				responseWithCursor := httpmock.NewStringResponse(200, `
				<?xml version='1.0' encoding='utf-8'?>
				<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
				<soapenv:Header/>
				<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
						<transport xmlns="http://gutp.jp/fiap/2009/11/">
							<header>
							<OK/>
								<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb">
									<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
								</query>
							</header>
							<body>
								<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-01T16:34:05.000+09:00">1</value>
								</point>
							</body>
							</transport>
					</ns2:queryRS>
				</soapenv:Body>
				</soapenv:Envelope>
				`)
				return responseWithCursor, nil
			}
		},
	)
	// テスト対象の関数を実行
	_, points, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time: time.Date(2012, 2, 1, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "1",
			},
			{
				Time: time.Date(2012, 2, 2, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "2",
			},
			{
				Time: time.Date(2012, 2, 3, 16, 34, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "3",
			},
		},
	}, points)
}

func TestFetchEmpty(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	responder := testutil.CustomHeaderBodyResponder(`
	<header>
		<OK/>
		<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="">
			<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
		</query>
	</header>
	<body>
	</body>
	`)
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	pointSets, points, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.NoError(t, err)
	assert.Equal(t, 0, len(pointSets))
	assert.Equal(t, 0, len(points))
}

func TestFetchFetchOnceError1(t *testing.T){
	// mockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		httpmock.NewErrorResponder(errors.New("mocked error")))

	// テスト対象の関数を実行
	_, _, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FetchOnce error on loop iteration 1")
	assert.Contains(t, err.Error(), "fiapFetch error")
	assert.Contains(t, err.Error(), "client.Call error")
}

func TestFetchFetchOnceError2(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		func (req *http.Request) (*http.Response, error){
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			if strings.Contains(string(body), `cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb"`) {
				return nil, errors.New("mocked error")
			} else {
				responseWithCursor := httpmock.NewStringResponse(200, `
				<?xml version='1.0' encoding='utf-8'?>
				<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
				<soapenv:Header/>
				<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
						<transport xmlns="http://gutp.jp/fiap/2009/11/">
							<header>
							<OK/>
								<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" cursor="a93f7094-4fd1-8e9a-749c-08e222bb0afb">
									<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
								</query>
							</header>
							<body>
								<point id="http://xxxxxxxx/tokyo/building1/Room101/">
									<value time="2012-02-01T16:34:05.000+09:00">1</value>
								</point>
							</body>
							</transport>
					</ns2:queryRS>
				</soapenv:Body>
				</soapenv:Envelope>
				`)
				return responseWithCursor, nil
			}
		},
	)
	// テスト対象の関数を実行
	_, _, _, err := Fetch("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", []model.UserInputKey{
		{ID: "http://xxxxxxxx/tokyo/building1/Room101/"},
	}, &model.FetchOption{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FetchOnce error on loop iteration 2")
	assert.Contains(t, err.Error(), "fiapFetch error")
	assert.Contains(t, err.Error(), "client.Call error")
}

func TestFetchByIdsWithKeyIdBoundary(t *testing.T){
	testcases := []struct {
		name string
		ids string
		expectedIds []string
	}{
		{
			name: "when one id is specified",
			ids: "http://xxxxxxxx/tokyo/building1/",
			expectedIds: []string{"http://xxxxxxxx/tokyo/building1/"},
		},
		{
			name: "when two ids are specified",
			ids: `"http://xxxxxxxx/tokyo/building1/", "http://xxxxxxxx/tokyo/building2/"`,
			expectedIds: []string{"http://xxxxxxxx/tokyo/building1/", "http://xxxxxxxx/tokyo/building2/"},
		},
	}

	for _, tc := range testcases {
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

				for _, id := range tc.expectedIds {
					if strings.Contains(bodyString, id) {
						return true
					}
				}
				return false
			})

			// Responder function
			responder := testutil.CustomBodyResponder(`
			<body>
				<point id="http://xxxxxxxx/tokyo/building1/Room101/">
					<value time="2012-02-04T16:35:05.000+09:00">50</value>
				</point>
				<pointSet id="http://xxxxxxxx/tokyo/building1/">
					<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
					<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
				</pointSet>
			</body>
			`)

			expectedPoints := map[string][]model.Value{
				"http://xxxxxxxx/tokyo/building1/Room101/": {
					{
						Time:  time.Date(2012, 2, 4, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
						Value: "50",
					},
				},
			}
			expectedPointSets := map[string]model.ProcessedPointSet{
				"http://xxxxxxxx/tokyo/building1/": {
					PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
					PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
				},
			}

			// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
			httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
					matcher,
					responder,
			)
			// テスト対象の関数を実行
			pointSets, points, _, _ := FetchByIdsWithKey("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", model.UserInputKeyNoID{}, tc.ids)
			assert.Equal(t, expectedPointSets, pointSets)
			assert.Equal(t, expectedPoints, points)
		})
	}
}

func TestFetchByIdsWithKeyMissingId(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 下記URLにPOSTしたときの挙動を定義
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", 
		httpmock.NewErrorResponder(errors.New("mocked error")))

	// テスト対象の関数を実行
	_, _, _, err := FetchByIdsWithKey("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", model.UserInputKeyNoID{
		Lteq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
		Gteq: testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))),
	}, "http://xxxxxxxx/tokyo/building1/")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Fetch error")
}

func TestFetchLatestCheckHttpReqAndSuccess(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
	toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
	ids := "http://xxxxxxxx/tokyo/building1/"
	expectedReqStrings := []string{"http://xxxxxxxx/tokyo/building1", tools.TimeToString(fromDate),tools.TimeToString(toDate), "maximum"}

	expectedPoints := map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time:  time.Date(2021, 1, 1, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "50",
			},
		},
	}
	expectedPointSets := map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
		},
	}

	// Matcher function
	matcher := httpmock.NewMatcher("", func(req *http.Request) bool {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return false
		}
		bodyString := string(bodyBytes)

		for _, expectedString := range expectedReqStrings {
			if !strings.Contains(bodyString, expectedString) {
				return false
			}
		}
		return true
	})

	// Responder function
	responder := testutil.CustomBodyResponder(`
		<body>
			<point id="http://xxxxxxxx/tokyo/building1/Room101/">
				<value time="2021-01-01T16:35:05.000+09:00">50</value>
			</point>
			<pointSet id="http://xxxxxxxx/tokyo/building1/">
				<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
				<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
			</pointSet>
		</body>
	`)
	// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
	httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
			matcher,
			responder,
	)
	// テスト対象の関数を実行
	pointSets, points,_, _ := FetchLatest("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
	fromDate,
	toDate,
	ids)

	assert.Equal(t, expectedPointSets, pointSets)
	assert.Equal(t, expectedPoints, points)
}

func TestFetchLatestFetchByIdsWithKeyError(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder := testutil.CustomBodyResponder("")
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	_, _, _, err := FetchLatest("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",nil,nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FetchByIdsWithKey error")
}

func TestFetchOldestCheckHttpReqAndSuccess(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
	toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
	ids := "http://xxxxxxxx/tokyo/building1/"
	expectedReqStrings := []string{"http://xxxxxxxx/tokyo/building1", tools.TimeToString(fromDate), tools.TimeToString(toDate), "minimum"}

	expectedPoints := map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time:  time.Date(2021, 1, 1, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "50",
			},
		},
	}
	expectedPointSets := map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
		},
	}

	// Matcher function
	matcher := httpmock.NewMatcher("", func(req *http.Request) bool {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return false
		}
		bodyString := string(bodyBytes)

		for _, expectedString := range expectedReqStrings {
			if !strings.Contains(bodyString, expectedString) {
				return false
			}
		}
		return true
	})

	// Responder function
	responder := testutil.CustomBodyResponder(`
		<body>
			<point id="http://xxxxxxxx/tokyo/building1/Room101/">
				<value time="2021-01-01T16:35:05.000+09:00">50</value>
			</point>
			<pointSet id="http://xxxxxxxx/tokyo/building1/">
				<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
				<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
			</pointSet>
		</body>
	`)
	// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
	httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
			matcher,
			responder,
	)
	// テスト対象の関数を実行
	pointSets, points, _, _ := FetchOldest("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
	fromDate,
	toDate,
	ids)

	assert.Equal(t, expectedPointSets, pointSets)
	assert.Equal(t, expectedPoints, points)
}

func TestFetchOldestFetchByIdsWithKeyError(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder := testutil.CustomBodyResponder("")
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	_, _, _, err := FetchOldest("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",nil,nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FetchByIdsWithKey error")
}

func TestFetchDateRangeCheckHttpReqAndSuccess(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fromDate := testutil.TimeToTimep(time.Date(2021, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
	toDate := testutil.TimeToTimep(time.Date(2021, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)))
	ids := "http://xxxxxxxx/tokyo/building1/"
	expectedReqStrings := []string{"http://xxxxxxxx/tokyo/building1", tools.TimeToString(fromDate), tools.TimeToString(toDate)}

	expectedPoints := map[string][]model.Value{
		"http://xxxxxxxx/tokyo/building1/Room101/": {
			{
				Time:  time.Date(2021, 1, 1, 16, 35, 5, 0, time.FixedZone("", 9*60*60)),
				Value: "50",
			},
		},
	}
	expectedPointSets := map[string]model.ProcessedPointSet{
		"http://xxxxxxxx/tokyo/building1/": {
			PointSetID: []string{"http://xxxxxxxx/tokyo/building1/Room101/"},
			PointID:    []string{"http://xxxxxxxx/tokyo/building1/Temperature/"},
		},
	}

	// Matcher function
	matcher := httpmock.NewMatcher("", func(req *http.Request) bool {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return false
		}
		bodyString := string(bodyBytes)

		for _, expectedString := range expectedReqStrings {
			if !strings.Contains(bodyString, expectedString) {
				return false
			}
		}
		return true
	})

	// Responder function
	responder := testutil.CustomBodyResponder(`
		<body>
			<point id="http://xxxxxxxx/tokyo/building1/Room101/">
				<value time="2021-01-01T16:35:05.000+09:00">50</value>
			</point>
			<pointSet id="http://xxxxxxxx/tokyo/building1/">
				<pointSet id="http://xxxxxxxx/tokyo/building1/Room101/" />
				<point id="http://xxxxxxxx/tokyo/building1/Temperature/" />
			</pointSet>
		</body>
	`)
	// 下記URLにPOSTし、かつ特定のrequest bodyを送信したときの挙動を定義
	httpmock.RegisterMatcherResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
			matcher,
			responder,
	)
	// テスト対象の関数を実行
	pointSets, points,_, _ := FetchOldest("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",
	fromDate,
	toDate,
	ids)

	assert.Equal(t, expectedPointSets, pointSets)
	assert.Equal(t, expectedPoints, points)
}

func TestFetchDateRangeFetchByIdsWithKeyError(t *testing.T){
	// httpmockの有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder := testutil.CustomBodyResponder("")
	httpmock.RegisterResponder("POST", "http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage", responder)

	// テスト対象の関数を実行
	_, _, _, err := FetchDateRange("http://iot.info.nara-k.ac.jp/axis2/services/FIAPStorage",nil,nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FetchByIdsWithKey error")
}