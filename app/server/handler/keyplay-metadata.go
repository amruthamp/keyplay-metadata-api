package handler

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"cloud.google.com/go/bigtable"
	"github.com/fubotv/keyplay-metadata-api/app/db"
	"github.com/fubotv/keyplay-metadata-api/app/model"
	"github.com/fubotv/keyplay-metadata-api/app/util"
	"goji.io/pat"
)

type ServiceHandler struct {
	DatabaseHandler *db.DBHandler
}

const (
	RowKeyDelimiter = "#"
)

func (s ServiceHandler) GetKeyplayMetadata(w http.ResponseWriter, r *http.Request) {

	programId := pat.Param(r, "programId")
	channelId := pat.Param(r, "channelId")
	id := pat.Param(r, "id")

	filter := bigtable.RowKeyFilter("^" + programId + RowKeyDelimiter + channelId + RowKeyDelimiter + id + RowKeyDelimiter + "metadata" + "$")

	keyplay, err := db.ReadRowFromBT(s.DatabaseHandler.Table, filter)

	if err != nil {
		util.JsonError(context.Background(), w, http.StatusNotFound, err)
	}
	util.Json(context.Background(), w, http.StatusOK, keyplay)
}

func (s ServiceHandler) CreateKeyPlayMetadata(w http.ResponseWriter, r *http.Request) {
	var keyplayAttribute model.KeyplayAttribute
	var keyplayMetadata model.KeyplayMetadata

	programId := pat.Param(r, "programId")
	channelId := pat.Param(r, "channelId")
	id := pat.Param(r, "id")

	keyplayFilter := bigtable.RowKeyFilter("^" + programId + RowKeyDelimiter + channelId + RowKeyDelimiter + id + "$")

	keyplayAttribute, _ = db.ReadKeyplayFromBT(s.DatabaseHandler.Table, keyplayFilter)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(keyplayAttribute.KeyData) >= 0 && len(reqBody) > 0 {
		json.Unmarshal(reqBody, &keyplayMetadata)
		ksAtrributes := reflect.ValueOf(keyplayMetadata.Metadata).MapKeys()

		validData := checkAttributes(keyplayAttribute.KeyData, ksAtrributes)

		if validData {

			rowkey := generateRowkey(programId, channelId, id)

			mut := bigtable.NewMutation()
			binary.Write(new(bytes.Buffer), binary.BigEndian, int64(1))
			keyplayMetadata.Id = id

			keyplayMetadataMarshalled, _ := json.Marshal(keyplayMetadata)

			mut.Set(db.ColumnFamilyName, "keyData", 0, keyplayMetadataMarshalled)

			err := db.WriteToBT(s.DatabaseHandler.Table, rowkey, mut)
			if err != nil {
				util.JsonError(context.Background(), w, http.StatusNotFound, err)
			}

		}

	} else {
		util.JsonError(context.Background(), w, http.StatusMethodNotAllowed, errors.New("please add the keyplay name"))
	}

}

func (s ServiceHandler) DeleteKeyplayMetadata(w http.ResponseWriter, r *http.Request) {

	id := pat.Param(r, "id")

	filter := bigtable.RowKeyFilter(".*" + RowKeyDelimiter + id + RowKeyDelimiter + "metadata" + "$")

	rowkey, _ := db.GetRowKey(s.DatabaseHandler.Table, filter)

	if len(rowkey) > 0 {
		mut := bigtable.NewMutation()
		mut.DeleteRow()

		err := db.WriteToBT(s.DatabaseHandler.Table, rowkey, mut)
		if err != nil {
			util.JsonError(context.Background(), w, http.StatusNotFound, err)
		}

	} else {
		util.JsonError(context.Background(), w, http.StatusMethodNotAllowed, errors.New("keyplay not found"))
	}
}

func checkAttributes(attributes []string, ksAtrributes []reflect.Value) bool {
	for _, attribute := range attributes {
		for _, ks := range ksAtrributes {
			if attribute != fmt.Sprintf("%v", ks) {
				return false
			}
		}
	}
	return true

}

func generateRowkey(programId string, channelId string, id string) string {
	rowkey := programId + RowKeyDelimiter + channelId + RowKeyDelimiter + id + RowKeyDelimiter + "metadata"
	return rowkey
}
