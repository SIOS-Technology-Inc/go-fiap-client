package testutil

import 	(
	"github.com/jarcoal/httpmock"
	"fmt"
)

// CustomBodyResponder returns a FIAP response with the given body content.
// CustomBodyResponderは指定されたbodyの内容を持つFIAPレスポンスを返します。
func CustomBodyResponder(bodyContent string) httpmock.Responder {
	responseTemplate := `<?xml version='1.0' encoding='utf-8'?>
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Header/>
			<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
							<transport xmlns="http://gutp.jp/fiap/2009/11/">
									<header>
											<OK/>
											<query id="e3264a29-b4a6-41dd-a6bb-cbf57b76e571" type="storage" acceptableSize="1000">
													<key id="xxxxxxxx/tokyo/building1/" attrName="time" select="maximum"/>
											</query>
									</header>
									%s
							</transport>
					</ns2:queryRS>
			</soapenv:Body>
	</soapenv:Envelope>`
	return httpmock.NewStringResponder(200, fmt.Sprintf(responseTemplate, bodyContent))
}

// CustomHeaderBodyResponder returns a FIAP response with the given header and body content.
// CustomHeaderBodyResponderは指定されたヘッダーとボディの内容を持つFIAPレスポンスを返します。
func CustomHeaderBodyResponder(bodyContent string) httpmock.Responder {
	responseTemplate := `<?xml version='1.0' encoding='utf-8'?>
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Header/>
			<soapenv:Body>
					<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
							<transport xmlns="http://gutp.jp/fiap/2009/11/">
									%s
							</transport>
					</ns2:queryRS>
			</soapenv:Body>
	</soapenv:Envelope>`
	return httpmock.NewStringResponder(200, fmt.Sprintf(responseTemplate, bodyContent))
}

// CustomTransportResponder returns a FIAP response with the given transport content.
// CustomTransportResponderは指定されたtransportの内容を持つFIAPレスポンスを返します。
func CustomTransportResponder(bodyContent string) httpmock.Responder {
	responseTemplate := `<?xml version='1.0' encoding='utf-8'?>
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Header/>
			<soapenv:Body>
				<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
					%s
				</ns2:queryRS>
			</soapenv:Body>
	</soapenv:Envelope>`
	return httpmock.NewStringResponder(200, fmt.Sprintf(responseTemplate, bodyContent))
}

// CustomTransportStatusCodeResponder returns a FIAP response with the given body content and status code.
// CustomTransportStatusCodeResponderは指定されたbodyの内容とステータスコードを持つFIAPレスポンスを返します。
func CustomTransportStatusCodeResponder(bodyContent string, statusCode int) httpmock.Responder {
	responseTemplate := `<?xml version='1.0' encoding='utf-8'?>
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
			<soapenv:Header/>
			<soapenv:Body>
				<ns2:queryRS xmlns:ns2="http://soap.fiap.org/">
					%s
				</ns2:queryRS>
			</soapenv:Body>
	</soapenv:Envelope>`
	return httpmock.NewStringResponder(statusCode, fmt.Sprintf(responseTemplate, bodyContent))
}