package db

import (
	"errors"
	"fmt"
	"function/logger"
	"function/models"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

var tablename = os.Getenv("TABLE_NAME")

var albums = []models.Album{
	{ID: uuid.NewString(), Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: uuid.NewString(), Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: uuid.NewString(), Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	{ID: uuid.NewString(), Title: "Kind of Blue", Artist: "Miles Davis", Price: 45.99},
	{ID: uuid.NewString(), Title: "A Love Supreme", Artist: "John Coltrane", Price: 59.99},
	{ID: uuid.NewString(), Title: "Time Out", Artist: "The Dave Brubeck Quartet", Price: 29.99},
	{ID: uuid.NewString(), Title: "Giant Steps", Artist: "John Coltrane", Price: 44.99},
	{ID: uuid.NewString(), Title: "The Shape of Jazz to Come", Artist: "Ornette Coleman", Price: 23.99},
	{ID: uuid.NewString(), Title: "Out to Lunch!", Artist: "Eric Dolphy", Price: 37.99},
	{ID: uuid.NewString(), Title: "Mingus Ah Um", Artist: "Charles Mingus", Price: 34.99},
	{ID: uuid.NewString(), Title: "Getz/Gilberto", Artist: "Stan Getz & João Gilberto", Price: 28.99},
	{ID: uuid.NewString(), Title: "Moanin'", Artist: "Art Blakey & The Jazz Messengers", Price: 31.99},
	{ID: uuid.NewString(), Title: "Speak No Evil", Artist: "Wayne Shorter", Price: 39.99},
	{ID: uuid.NewString(), Title: "Somethin' Else", Artist: "Cannonball Adderley", Price: 36.99},
	{ID: uuid.NewString(), Title: "The Black Saint and the Sinner Lady", Artist: "Charles Mingus", Price: 48.99},
}

var log = logger.GetLogger()

func getDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func CreateSampleDataRecords() {
	err := CreateAlbum(albums[0])
	checkAvailable(err)

	// Promise.allSettled([for album of albums {createAlbum(album)}])
	var waitGroup sync.WaitGroup
	for _, album := range albums[1:] {
		waitGroup.Add(1)

		go func(album models.Album) {
			defer waitGroup.Done()
			CreateAlbum(album)
		}(album)
	}

	waitGroup.Wait()
}

func checkAvailable(err error) error {
	retries := 0
	for err != nil {
		if retries > 3 {
			log.Fatal("retired three times...")
		}

		aerr, ok := err.(awserr.Error)
		if !ok {
			log.Fatal("err.(awserr.Error) is not awserr.Error")
		}

		if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			log.Error("Table is unavailable ... waiting 5 seconds")
			time.Sleep(5 * time.Second)
			log.Info("Checking...")
			err = CreateAlbum(albums[0])
		}
	}

	return nil
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
	CreateSampleDataRecords()
}

func GetAlbums() ([]models.Album, error) {
	albums := make([]models.Album, 0)
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
	item, err := dynamodbattribute.MarshalMap(models.AlbumRecord{
		PK:     "album#" + album.ID,
		SK:     album.Artist,
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
		"pk": {
			S: aws.String("album#" + id),
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
	sk, err := GetSingleAlbum(id)
	// TODO: Refactor
	if err != nil {
		log.Errorf("Error deleting album: %s", err)
		return errors.New("error deleting album")
	}

	inputKey := map[string]*dynamodb.AttributeValue{
		"pk": {
			S: aws.String("album#" + id),
		},
		"sk": {
			S: aws.String(sk.ID),
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
