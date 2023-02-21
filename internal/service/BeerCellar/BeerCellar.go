package BeerCellar

import (
	"BeerCellar/internal/domain"
	"errors"
	"time"
)

type DataBase interface {
	InsertDescription(id int, Name string, ABV float32) error
	InsertBeer(id int, quantity int, dateRegistration time.Time, dateBatch time.Time) error
	SelectRemains(Name string, DaysAging int, MaxDateBeer time.Time) ([]domain.CellarRemains, error)
	DeletePart(Name string, StartBatch time.Time, EndBatch time.Time) (int64, error)
	//WriteLog(description string, dateLog time.Time) error
}

type BeerService struct {
	DataBase DataBase
}

func NewBeerService(dataBase DataBase) *BeerService {
	return &BeerService{DataBase: dataBase}
}

func (beer *BeerService) InsertBeer(Id int, Name string, ABV float32, Quantity int) error {
	if Id == 0 {
		return errors.New("empty part.id")
	}

	if Quantity == 0 {
		return errors.New("empty part.quantity")
	}

	err := beer.DataBase.InsertDescription(Id, Name, ABV)
	if err != nil {
		return err
	}

	dateNow := time.Now()

	return beer.DataBase.InsertBeer(Id, Quantity, dateNow, dateNow)
}

func (beer *BeerService) TakeBeer(Name string, Quantity int, DaysAging int) ([]domain.CellarRemains, error) {
	var s []domain.CellarRemains
	if Name == "" {
		return s, errors.New("empty beer.name")
	}

	if Quantity == 0 {
		return s, errors.New("empty part.quantity")
	}
	MaxDateBeer := time.Now().Add(-1 * time.Hour * 24 * time.Duration(DaysAging))
	remains, err := beer.DataBase.SelectRemains(Name, DaysAging, MaxDateBeer)
	if err != nil {
		return s, err
	}

	PartQuantity := 0
	var PartBatch time.Time

	dateNow := time.Now()
	beerQuantity := Quantity
	beerID := 0
	quantity := 0
	//isBreak := false

	for _, value := range remains {

		PartQuantity = value.Sum_beer
		PartBatch = value.Batch
		beerID = value.Beer_id

		if err != nil {
			return s, err
		}
		if beerQuantity <= PartQuantity {
			quantity = beerQuantity
			//isBreak = true
			beerQuantity = 0
		} else {
			quantity = PartQuantity
			//isBreak = false
			beerQuantity = beerQuantity - PartQuantity

		}
		//_, err = db.Exec("insert into cellar (beer_id, quantity, date, batch) values ($1,$2,$3,$4)", beerID, -1*quantity, dateNow, PartBatch)
		beer.DataBase.InsertBeer(beerID, -1*quantity, dateNow, PartBatch)

		if err != nil {
			return s, err
		}
		s = append(s, domain.CellarRemains{
			Sum_beer: quantity,
			Batch:    PartBatch, //PartBatch.Format("2006-01-02"),
			Beer_id:  beerID,
		})
		if beerQuantity == 0 {
			break
		}
	}
	if beerQuantity == Quantity {
		err = errors.New("No beer")
	}
	return s, err
}

func (beer *BeerService) DeletePart(Name string, BatchString string) (int64, error) {

	var s int64

	startDay, err := time.Parse("2006-01-02", BatchString)
	if err != nil {
		return s, err
	}

	endDay := startDay.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	s, err = beer.DataBase.DeletePart(Name, startDay, endDay)
	if err != nil {
		return s, err
	}
	return s, err
}
