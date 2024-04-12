package testutil

import 	(
	"github.com/jarcoal/httpmock"
	"fmt"
)


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
