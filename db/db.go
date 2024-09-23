package db

import (
	"example/web-service-gin/models"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tablename = "Web-Service-Gin-Tutorial-Albums"

func GetDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func CreateAlbumTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("title"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tablename),
		Tags: []*dynamodb.Tag{
			{
				Key:   aws.String("Important"),
				Value: aws.String("NOT AT ALL"),
			},
			{
				Key:   aws.String("DELETE"),
				Value: aws.String("YES"),
			},
			{
				Key:   aws.String("Purpose"),
				Value: aws.String("Golang tutorial"),
			},
		},
	}

	svc := GetDynamoSession()
	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created Albums table: ", tablename)
}

func GetAlbums() []models.Album {
	var albums []models.Album
	svc := GetDynamoSession()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: &tablename,
	})
	if err != nil {
		log.Printf("Query API call failed: %s", err)
	}

	for _, item := range result.Items {
		album := models.Album{}

		// TODO: Fix unmarshalling issues
		err = dynamodbattribute.UnmarshalMap(item, &album)
		if err != nil {
			log.Printf("Got error unmarshalling: %s", err)
			// return false, err.Error()
		}

		albums = append(albums, album)
	}

	return albums
}

func CreateAlbum(album models.Album) (bool, string) {
	svc := GetDynamoSession()
	av, err := dynamodbattribute.MarshalMap(models.MarshalledAlbum{
		ID:     album.ID.String(),
		Title:  album.Title,
		Artist: album.Artist,
		Price:  album.Price,
	})
	if err != nil {
		log.Printf("Got error marshalling new album: %s", err)
		return false, err.Error()
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %s", err)
		return false, err.Error()
	}

	return true, ""
}
