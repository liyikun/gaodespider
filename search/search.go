package search

import (
	"eatinlife.com/gaodespider/insertdb"
	"eatinlife.com/serverutil/logger"
)

func RunSpider() {
	gaoResults := GetGaoResult()
	connConfig := insertdb.TakeCoon()
	db, err := insertdb.ConnectDb(connConfig)
	defer db.Close()
	if err != nil {
		logger.Error.Panicln("db connect error", err.Error())
	}
	for info := range gaoResults {
		err = insertdb.InsertRestaurantInfo(info, db)
		if err != nil {
			logger.Error.Panicln("db insert error", err.Error())
		} else {
			logger.Info.Println("db insert success")
		}
	}
}
