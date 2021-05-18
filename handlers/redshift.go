package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	rsa "github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

type RedshiftEvent struct {
	ClusterId  string       `json:"cluster_id"`
	Database   string       `json:"database"`
	DbUser     string       `json:"database_user"`
	Arn        string       `json:"arn"`
	Statements []*Statement `json:"statements"`
}

type Statement struct {
	Sql  string `json:"sql"`
	Name string `json:"name"`
}

var rsSvc *rsa.RedshiftDataAPIService

func init() {
	// init redshift client
	mySession := session.Must(session.NewSession())
	rsSvc = rsa.New(mySession) // may need to add configurations
	if rsSvc == nil {
		log.Panic("did not connect to the redshift api service")
	}
}

func (r *RedshiftEvent) ExecuteStatements(ctx context.Context) error {
	input := &rsa.ExecuteStatementInput{
		ClusterIdentifier: &r.ClusterId,
		Database:          &r.Database,
		DbUser:            &r.DbUser,
		SecretArn:         &r.Arn,
	}

	for _, statement := range r.Statements {
		input.SetStatementName(statement.Name)
		input.SetSql(statement.Sql)

		err := input.Validate()
		if err != nil {
			return fmt.Errorf("the event did not pass validation: %s", err)
		}

		output, err := rsSvc.ExecuteStatement(input)
		if err != nil {
			return fmt.Errorf("execution of the statement failed: %s", err)
		}

		log.Println(output.GoString())
	}

	return nil
}
