package main

import (
	"fmt"
	"github.com/mofax/iso8583"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func sendIso(writer http.ResponseWriter, request *http.Request) {

	var response Response
	var iso Iso8583

	reqBody, _ := ioutil.ReadAll(request.Body)
	req := string(reqBody)
	log.Println("New Request")
	log.Printf("ISO Message: %v\n", req)

	err := doProducer(broker, topic1, req)

	if err != nil {
		errDesc := fmt.Sprintf("Failed sent to Kafka\nError: %v", err)
		response.ResponseCode, response.ResponseDescription = 500, errDesc
		log.Println(err)
		jsonFormatter(writer, response, 500)
	} else {
		msg, err := consumeResponse(broker, group, []string{topic2})
		if err != nil {
			errDesc := fmt.Sprintf("Failed to get response from Kafka\nError: %v", err)
			response.ResponseCode, response.ResponseDescription = 500, errDesc
			log.Println(err)
			jsonFormatter(writer, response, 500)
		} else {

			if msg == "" {
				errDesc := "Got empty response"
				response.ResponseCode, response.ResponseDescription = 500, errDesc
				log.Println(errDesc)
				jsonFormatter(writer, response, 500)
			} else {

				header := msg[0:4]
				data := msg[4:]

				isoStruct := iso8583.NewISOStruct("spec1987.yml", false)

				isoParsed, err := isoStruct.Parse(data)
				if err != nil {
					log.Printf("Error parsing iso message\nError: %v", err)
				}

				iso.Header, _ = strconv.Atoi(header)
				iso.MTI = isoParsed.Mti.String()
				iso.Hex, _ = iso8583.BitMapArrayToHex(isoParsed.Bitmap)

				iso.Message, err = isoParsed.ToString()
				if err != nil {
					log.Printf("Iso Parsed failed convert to string.\nError: %v", err)
				}

				// create file from response
				event := header + iso.Message
				filename := "Response_to_" + isoParsed.Elements.GetElements()[3] + "@" + fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
				file := CreateFile(filename, event)
				log.Println("File created: ", file)

				desc := "Success"
				iso.ResponseStatus.ResponseCode, iso.ResponseStatus.ResponseDescription = 200, desc
				log.Println(desc)
				jsonFormatter(writer, iso, 200)

			}

		}
	}

}

func sendFile(writer http.ResponseWriter, request *http.Request) {

	var response Response
	var iso Iso8583

	reqBody, _ := ioutil.ReadAll(request.Body)
	filename := string(reqBody)
	filename = "storage/request/" + filename

	if !strings.Contains(filename, ".txt") {
		filename += ".txt"
	}

	log.Println("New Request")

	if CheckExist(filename) {
		req := ReadFile(filename)

		log.Printf("ISO Message: %v\n", req)

		err := doProducer(broker, topic1, req)

		if err != nil {
			errDesc := fmt.Sprintf("Failed sent to Kafka\nError: %v", err)
			response.ResponseCode, response.ResponseDescription = 500, errDesc
			log.Println(err)
			jsonFormatter(writer, response, 500)
		} else {
			msg, err := consumeResponse(broker, group, []string{topic2})
			if err != nil {
				errDesc := fmt.Sprintf("Failed to get response from Kafka\nError: %v", err)
				response.ResponseCode, response.ResponseDescription = 500, errDesc
				log.Println(err)
				jsonFormatter(writer, response, 500)
			} else {

				if msg == "" {
					errDesc := "Got empty response"
					response.ResponseCode, response.ResponseDescription = 500, errDesc
					log.Println(errDesc)
					jsonFormatter(writer, response, 500)
				} else {

					header := msg[0:4]
					data := msg[4:]

					isoStruct := iso8583.NewISOStruct("spec1987.yml", false)

					isoParsed, err := isoStruct.Parse(data)
					if err != nil {
						log.Printf("Error parsing iso message\nError: %v", err)
					}

					iso.Header, _ = strconv.Atoi(header)
					iso.MTI = isoParsed.Mti.String()
					iso.Hex, _ = iso8583.BitMapArrayToHex(isoParsed.Bitmap)

					iso.Message, err = isoParsed.ToString()
					if err != nil {
						log.Printf("Iso Parsed failed convert to string.\nError: %v", err)
					}

					// create file from response
					event := header + iso.Message
					filename := "Response_to_" + isoParsed.Elements.GetElements()[3] + "@" + fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
					file := CreateFile(filename, event)
					log.Println("File created: ", file)

					desc := "Success"
					iso.ResponseStatus.ResponseCode, iso.ResponseStatus.ResponseDescription = 200, desc
					log.Println(desc)
					jsonFormatter(writer, iso, 200)

				}

			}
		}

	} else {

		errDesc := "File not found"
		response.ResponseCode, response.ResponseDescription = 404, errDesc
		log.Println(errDesc)
		log.Println("Process failed")
		jsonFormatter(writer, response, 404)

	}

}
