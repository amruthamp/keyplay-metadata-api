package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/bigtable"
	"github.com/fubotv/fubotv-logging/v3/logging"
	"github.com/fubotv/keyplay-metadata-api/app/config"
	"github.com/fubotv/keyplay-metadata-api/app/model"
	"google.golang.org/api/option"
)

const (
	VidaiBtTable     = "VIDAI_DATA_JOBS"
	ColumnFamilyName = "event_family"
)

type DBHandler struct {
	Client *bigtable.Client
	Table  *bigtable.Table
}

// function to connect database
func CreateDBHandler() (*DBHandler, error) {
	cfg := config.GetConfig()
	ctx := context.Background()

	logging.Info(context.Background(), fmt.Sprintf("creating BigTable client for BT_INSTANCE_ID: %s, BT_PROJECT_ID:%s", cfg.DatabaseCfg.BTInstanceId, cfg.DatabaseCfg.BTProjectId))

	clientOpts := []option.ClientOption{
		option.WithCredentialsFile(cfg.DatabaseCfg.BTConnectionCredentials),
	}

	bigTableClient, err := bigtable.NewClient(
		ctx,
		cfg.DatabaseCfg.BTProjectId,
		cfg.DatabaseCfg.BTInstanceId,
		clientOpts...)

	if err != nil {
		logging.Error(context.Background(), err, "Could not create data operations client")
		return nil, err
	}

	Table := bigTableClient.Open(VidaiBtTable)

	return &DBHandler{bigTableClient, Table}, nil
}

func ReadRowFromBT(table *bigtable.Table, filter bigtable.Filter) (model.KeyplayMetadata, error) {
	var keyplayMetadata model.KeyplayMetadata

	err := table.ReadRows(context.Background(), bigtable.RowRange{}, func(row bigtable.Row) bool {
		for _, col := range row[ColumnFamilyName] {
			json.Unmarshal(col.Value, &keyplayMetadata)
		}
		return true
	}, bigtable.RowFilter(filter))

	if err != nil {
		return model.KeyplayMetadata{}, err
	}

	return keyplayMetadata, nil
}

func ReadKeyplayFromBT(table *bigtable.Table, filter bigtable.Filter) (model.KeyplayAttribute, error) {
	var attributeData model.KeyplayAttribute

	err := table.ReadRows(context.Background(), bigtable.RowRange{}, func(row bigtable.Row) bool {
		for _, col := range row[ColumnFamilyName] {
			qualifier := col.Column[strings.IndexByte(col.Column, ':')+1:]

			if qualifier == "keyData" {
				json.Unmarshal(col.Value, &attributeData.KeyData)
			} else if qualifier == "id" {
				attributeData.Id = string(col.Value)
			}
		}
		return true
	}, bigtable.RowFilter(filter))

	if err != nil {
		return model.KeyplayAttribute{}, err
	}

	return attributeData, nil
}

func GetRowKey(tableName *bigtable.Table, filter bigtable.Filter) (string, error) {
	var rowKeys string

	err := tableName.ReadRows(context.Background(), bigtable.RowRange{},
		func(row bigtable.Row) bool {
			rowKeys = row.Key()
			return true
		}, bigtable.RowFilter(filter))

	if err != nil {
		return "", err
	}

	return rowKeys, nil
}

func WriteToBT(tableName *bigtable.Table, rowkey string, mut *bigtable.Mutation) error {
	err := tableName.Apply(context.Background(), rowkey, mut)
	if err != nil {
		return err
	}
	return nil
}
