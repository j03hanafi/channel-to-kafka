package main

import (
	"fmt"
	"github.com/mofax/iso8583"
	"io/ioutil"
	"net/http"
	"strconv"
)

func sendIso(writer http.ResponseWriter, request *http.Request) {

	var response Response
	var iso Iso8583

	reqBody, _ := ioutil.ReadAll(request.Body)
	req := string(reqBody)
	fmt.Println("New Request")
	fmt.Printf("ISO Message: %v\n", req)

	err := doProducer(broker, topic1, req)

	if err != nil {
		errDesc := fmt.Sprintf("Failed sent to Kafka\nError: %v", err)
		response.ResponseCode, response.ResponseDescription = 500, errDesc
		jsonFormatter(writer, response, 500)
	} else {
		msg, err := consumeResponse(broker, group, []string{topic2})
		if err != nil {
			errDesc := fmt.Sprintf("Failed to get response from Kafka\nError: %v", err)
			response.ResponseCode, response.ResponseDescription = 500, errDesc
			jsonFormatter(writer, response, 500)
		} else {

			if msg == "" {
				response.ResponseCode, response.ResponseDescription = 500, "Got empty response"
				jsonFormatter(writer, response, 500)
			} else {

				header := msg[0:4]
				data := msg[4:]

				isoStruct := iso8583.NewISOStruct("spec1987.yml", false)

				isoParsed, err := isoStruct.Parse(data)
				if err != nil {
					fmt.Printf("Error parsing iso message\nError: %v", err)
				}

				iso.Header, _ = strconv.Atoi(header)
				iso.MTI = isoParsed.Mti.String()
				iso.Hex, _ = iso8583.BitMapArrayToHex(isoParsed.Bitmap)

				iso.Message, err = isoParsed.ToString()
				if err != nil {
					fmt.Printf("Iso Parsed failed convert to string.\nError: %v", err)
				}

				//event := header + iso.Message

				iso.ResponseStatus.ResponseCode, iso.ResponseStatus.ResponseDescription = 200, "Success"
				jsonFormatter(writer, iso, 200)

			}

		}
	}

}
