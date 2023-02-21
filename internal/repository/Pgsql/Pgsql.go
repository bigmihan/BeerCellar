package Pgsql

import (
	"BeerCellar/internal/domain"
	"database/sql"
	"time"
)

type DataBase struct {
	db *sql.DB
}

func NewDataBase(db *sql.DB) *DataBase {
	return &DataBase{db: db}
}

func (cellar *DataBase) InsertDescription(id int, Name string, ABV float32) error {
	rows, err := cellar.db.Query("select * from beer_description where beer_id = $1", id)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		_, err = cellar.db.Exec("insert into beer_description (beer_id, name, abv) values ($1,$2,$3)", id, Name, ABV)

	}

	return err

}

func (cellar *DataBase) InsertBeer(id int, quantity int, dateRegistration time.Time, dateBatch time.Time) error {

	_, err := cellar.db.Exec("insert into cellar (beer_id, quantity, date, batch) values ($1,$2,$3,$4)", id, quantity, dateRegistration, dateBatch)

	return err

}

func (cellar *DataBase) SelectRemains(Name string, DaysAging int, MaxDateBeer time.Time) ([]domain.CellarRemains, error) {

	rows, err := cellar.db.Query("select sum(cellar.quantity) as sum_beer, cellar.batch, cellar.beer_id from cellar left join beer_description on cellar.beer_id=beer_description.beer_id where beer_description.name=$1 And (($2=0) or ($2<>0 AND cellar.batch<=$3)) group by cellar.beer_id, cellar.batch having sum(cellar.quantity)>0 ORDER BY cellar.batch ", Name, DaysAging, MaxDateBeer)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	//var Remains []domain.CellarRemains
	Remains := make([]domain.CellarRemains, 0)
	for rows.Next() {
		var remain domain.CellarRemains
		rows.Scan(&remain.Sum_beer, &remain.Batch, &remain.Beer_id)
		Remains = append(Remains, remain)
	}
	return Remains, nil

}

func (cellar *DataBase) DeletePart(Name string, StartBatch time.Time, EndBatch time.Time) (int64, error) {
	rows, err := cellar.db.Exec("delete from cellar where exists (select 1 from cellar inner join beer_description on cellar.beer_id=beer_description.beer_id where beer_description.name=$1 AND cellar.batch>=$2 AND cellar.batch<=$3)", Name, StartBatch, EndBatch)
	if err != nil {
		return 0, err
	}

	s, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}
	return s, nil
}

func (cellar *DataBase) Write(p []byte) (int, error) {
	_, err := cellar.db.Exec("insert into log (description, date) values ($1,$2)", string(p), time.Now())
	return len(p), err
}
