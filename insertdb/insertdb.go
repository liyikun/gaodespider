package insertdb

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	mathcers "eatinlife.com/gaodespider/matchers"
)

func InsertRestaurantInfo(info mathcers.RestaurantInfo, db *sql.DB) error {
	stmt, err := db.Prepare(`insert into Restaurant_Gao
	(name,types,address,location,tel,pname,cityname,adname,photos)
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	on conflict (name) do update
	   set tel=$3,
	   photos=$9
	RETURNING id   
	`)
	if err != nil {
		return err
	}
	photos := info.TakePhotos()
	tel := fmt.Sprintf("%v", info.Tel)
	_, err = stmt.Exec(
		info.Name,
		info.Types,
		info.Address,
		info.Location,
		tel,
		info.Pname,
		info.Cityname,
		info.Adname,
		pq.Array(photos))
	return err
}
