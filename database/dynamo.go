package database

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jamesonwilliams/golang-sso-google/auth"
	"log"
)

type DynamoDatabase struct {
	Region    string
	TableName string
}

func (ddb *DynamoDatabase) RetrieveUser(email string) (auth.User, error) {
	sess := session.New(&aws.Config{Region: aws.String(ddb.Region)})
	svc := dynamodb.New(sess)

	params := &dynamodb.QueryInput{
		TableName: aws.String(ddb.TableName), // Required
		Select:    aws.String("ALL_ATTRIBUTES"),

		ConsistentRead: aws.Bool(true),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": { // Required
				ComparisonOperator: aws.String("EQ"), // Required
				AttributeValueList: []*dynamodb.AttributeValue{
					{ // Required
						S: aws.String(email),
					},
				},
			},
		},
	}

	resp, err := svc.Query(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return auth.User{}, err
	}

	// Pretty-print the response data.
	fmt.Println(resp)

	users := []auth.User{}
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &users)
	if err != nil || len(users) <= 0 {
		return auth.User{}, errors.New("No users")
	}

	fmt.Println(users[0])

	return users[0], err
}

func (ddb *DynamoDatabase) SaveUser(user *auth.User) error {
	fmt.Println("saving user " + user.Email + "...")
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Println("Failed to marshal data", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(ddb.TableName),
		ConditionExpression: aws.String("attribute_not_exists(email)"),
	}
	log.Println(input)

	sess := session.New(&aws.Config{Region: aws.String(ddb.Region)})
	svc := dynamodb.New(sess)
	result, err := svc.PutItem(input)
	if err != nil {
		log.Println("Failed to save item to table", err)
	}

	log.Println("Successfully added item to table.", result)

	return err
}
