package services

import (
	"context"
	"database/sql"
	"log"
	"math"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/franklyinc/frankly-lambda-video-s3-reporting/model"
)

func UpdateTableWithFolderData(sqlServerConnString string, folder model.Folder) (error) {
	db, err := sql.Open("sqlserver", sqlServerConnString)
	if err != nil {
		log.Println("ERROR: Cannot connect, ", err)
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = db.Ping()
	if err != nil {
		log.Println("ERROR: Cannot connect, ", err)
		return err
	}
	defer db.Close()

	affiliateName := strings.Replace(folder.Name, "/", "", 1)
	var sizeInMegaBytes float64
	sizeInMegaBytes = Round(folder.SizeInBytes / 1024.0, .1, 2)
	_, err = db.ExecContext(ctx, "p_prod_insert_VideoDailyTotalProxySize",
		sql.Named("AffiliateName", affiliateName),
		sql.Named("Size", sizeInMegaBytes),
	)
	if err != nil {
		return err
	}

	return nil
}

func Round(val float64, roundOn float64, places int) float64 {

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	var round float64
	if val > 0 {
		if div >= roundOn {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
	} else {
		if div >= roundOn {
			round = math.Floor(digit)
		} else {
			round = math.Ceil(digit)
		}
	}

	return round / pow
}
