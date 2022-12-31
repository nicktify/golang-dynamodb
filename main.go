package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
)

type Value struct {
	Value     string `dynamodbav:"value" json:"value"`
	Timestamp string `dynamodbav:"timestamp" json:"timestamp"`
}

type Device struct {
	DeviceID string  `dynamodbav:"deviceID" json:"deviceID"`
	Testing  string  `dynamodbav:"testing" json:"testing"`
	Flow     []Value `dynamodbav:"last_flow" json:"last_flow"`
}

type Person struct {
	name string
	age  int
}

type GetDeviceBody struct {
	DeviceID string
}

type CreateItemBody struct {
	ID   string
	Name string
	Age  string
}

func getItem(deviceID string) (Device, error) {
	var device Device

	cfg, err := config.LoadDefaultConfig(context.TODO())
	client := dynamodb.NewFromConfig(cfg)

	env, _ := godotenv.Read(".env")

	response, err := client.GetItem(context.Background(),
		&dynamodb.GetItemInput{
			TableName: aws.String(env["DEVICE_FILTER_TABLENAME"]),
			Key: map[string]types.AttributeValue{
				"deviceID": &types.AttributeValueMemberS{Value: *aws.String(deviceID)},
			},
		})
	if err != nil {
		return device, nil
	}

	if response.Item == nil {
		return device, nil
	}

	err = attributevalue.UnmarshalMap(response.Item, &device)

	if err != nil {
		return device, nil
	}
	return device, nil
}

func getItemHandleFunc(w http.ResponseWriter, req *http.Request) {
	var body GetDeviceBody

	fmt.Println(req.Body)

	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := getItem(body.DeviceID)

	fmt.Println(body)

	json.NewEncoder(w).Encode(device)
}

func main() {
	http.HandleFunc("/get-item", getItemHandleFunc)
	http.ListenAndServe(":4000", nil)
}
