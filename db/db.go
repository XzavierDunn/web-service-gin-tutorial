package db

import (
	"example/web-service-gin/models"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

var tablename = "Web-Service-Gin-Tutorial-Albums"

func getDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func createSampleDataRecords() {
	var albums = []models.Album{
		{ID: uuid.NewString(), Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: uuid.NewString(), Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: uuid.NewString(), Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}

	for _, album := range albums {
		CreateAlbum(album)
	}
}

func createAlbumTable() {
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

	svc := getDynamoSession()
	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created Albums table: ", tablename)
}

func InitTableWithData() {
	createAlbumTable()
	createSampleDataRecords()
}

func GetAlbums() ([]models.Album, error) {
	var albums []models.Album
	svc := getDynamoSession()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: &tablename,
	})
	if err != nil {
		log.Printf("Query API call failed: %s", err)
		return albums, err
	}

	for _, item := range result.Items {
		album := models.Album{}
		err = dynamodbattribute.UnmarshalMap(item, &album)
		if err != nil {
			log.Printf("Got error unmarshalling: %s", err)
			return albums, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func CreateAlbum(album models.Album) error {
	svc := getDynamoSession()
	av, err := dynamodbattribute.MarshalMap(models.Album{
		ID:     album.ID,
		Title:  album.Title,
		Artist: album.Artist,
		Price:  album.Price,
	})
	if err != nil {
		log.Printf("Got error marshalling new album: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %s", err)
		return err
	}

	return nil
}
