package db

import (
	"errors"
	"fmt"
	"time"
	"web-service-gin/logger"
	"web-service-gin/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

var tablename = "Web-Service-Gin-Tutorial-Albums"

var log = logger.GetLogger()

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
		err := CreateAlbum(album)
		for err != nil {
			// type assertion https://go.dev/tour/methods/15
			aerr, ok := err.(awserr.Error)
			if !ok {
				log.Fatal("err.(awserr.Error) is not awserr.Error")
			}

			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				log.Error("Table is unavailable ... waiting 5 seconds")
				time.Sleep(5 * time.Second)
				log.Info("Retrying sample data creation")
				err = CreateAlbum(album)
			}
		}
	}
}

func createAlbumTable() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
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
		log.Errorf("Got error calling CreateTable: %s", err)
		return err
	}

	fmt.Println("Created Albums table: ", tablename)
	return nil
}

// checkIfTableExists
func _() error {
	svc := getDynamoSession()
	response, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	})
	if err != nil {
		log.Errorf("Error describing table: %s", err)
		return err
	}

	log.Infof("Found: %v", response.Table.TableName)
	return nil
}

func InitTableWithData() {
	err := createAlbumTable()
	if err != nil {
		log.Fatalf(err.Error())
	}
	createSampleDataRecords()
}

func GetAlbums() ([]models.Album, error) {
	var albums []models.Album
	svc := getDynamoSession()

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: &tablename,
	})
	if err != nil {
		log.Errorf("Query API call failed: %s", err)
		return albums, err
	}

	for _, item := range result.Items {
		album := models.Album{}
		err = dynamodbattribute.UnmarshalMap(item, &album)
		if err != nil {
			log.Errorf("Got error unmarshalling: %s", err)
			return albums, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func CreateAlbum(album models.Album) error {
	svc := getDynamoSession()
	item, err := dynamodbattribute.MarshalMap(models.Album{
		ID:     album.ID,
		Title:  album.Title,
		Artist: album.Artist,
		Price:  album.Price,
	})
	if err != nil {
		log.Errorf("Got error marshalling new album: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tablename),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Errorf("Got error calling PutItem: %s", err)
		return err
	}

	log.Infof("Saved album: %s", album.ID)
	return nil
}

func GetSingleAlbum(id string) (models.Album, error) {
	album := models.Album{}
	inputKey := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	svc := getDynamoSession()
	input := &dynamodb.GetItemInput{
		Key:       inputKey,
		TableName: &tablename,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		log.Errorf("Error fetching item: %s", err)
		return album, err
	}

	if result.Item == nil {
		return album, errors.New("album not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &album)
	if err != nil {
		log.Errorf("Error unmarshalling item: %s", err)
		return album, err
	}

	return album, nil
}

func DeleteAlbum(id string) error {
	svc := getDynamoSession()
	inputKey := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	input := &dynamodb.DeleteItemInput{
		Key:          inputKey,
		ReturnValues: aws.String("ALL_OLD"),
		TableName:    &tablename,
	}

	result, err := svc.DeleteItem(input)
	if err != nil {
		log.Errorf("Error deleting album: %s", err)
		return errors.New("error deleting album")
	}

	if result.Attributes == nil {
		return errors.New("album does not exist")
	}

	return nil
}
