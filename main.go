package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type soapRQ struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soap,attr"`
	XMLNsXSI  string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD  string   `xml:"xmlns:xsd,attr"`
	Body      soapBody
}

type soapBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	GetCursOnDate *GetCursOnDate
}

type GetCursOnDate struct {
	XMLs string `xml:"xmlns,attr"`
	OnDate time.Time `xml:"On_date"`
}

func main() {
	soapAction := "http://web.cbr.ru/GetCursOnDate"
	ws := "http://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"
	payload := &GetCursOnDate{
		XMLs: "http://web.cbr.ru/",
		OnDate:  time.Now(),
	}
	res, err := soapCall(ws, soapAction, payload)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("RESULT:" + string(res))
}

//func soapCallHandleResponse(ws string, action string, payloadInterface interface{}, result interface{}) error {
//	body, err := soapCall(ws, action, payloadInterface)
//	if err != nil {
//		return err
//	}
//
//	err = xml.Unmarshal(body, &result)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func soapCall(ws string, action string, payloadInterface *GetCursOnDate) ([]byte, error) {
	v := soapRQ{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsXSD:  "http://www.w3.org/2001/XMLSchema",
		XMLNsXSI:  "http://www.w3.org/2001/XMLSchema-instance",
		Body: soapBody{
			GetCursOnDate: payloadInterface,
		},
	}
	payload, err := xml.MarshalIndent(v, "", "  ")

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return bodyBytes, nil
}