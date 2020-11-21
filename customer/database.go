package customer

import (
	"database/sql"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

func connection() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.WithError(err).Error("cannot connect to database")
	}
	return db, err
}

func InitServer() error {
	db, err := connection()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(`create table if not exists customers (
    id     serial primary key,
    name   text,
    email  text,
    status text
)`)
	if err != nil {
		return err
	}
	return nil
}

func Index() ([]Customer, error) {
	db, err := connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	stmt, err := db.Prepare("select id, name, email, status from customers")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var customers []Customer
	for rows.Next() {
		var customer Customer
		err = rows.Scan(&customer.Id, &customer.Name, &customer.Email, &customer.Status)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil

}

func Show(id int) (Customer, error) {
	var customer Customer
	db, err := connection()
	if err != nil {
		return customer, err
	}
	defer db.Close()
	stmt, err := db.Prepare("select id, name, email, status from customers where id=$1")
	if err != nil {
		return customer, err
	}
	row := stmt.QueryRow(id)

	err = row.Scan(&customer.Id, &customer.Name, &customer.Email, &customer.Status)
	return customer, err
}

func Update(id int, c Customer) (Customer, error) {
	db, err := connection()
	if err != nil {
		return Customer{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 where id=$1")
	if err != nil {
		return Customer{}, err
	}

	if _, err := stmt.Exec(id, c.Name, c.Email, c.Status); err != nil {
		log.WithError(err).Error("cannot update")
		return Customer{}, err
	}
	return c, nil
}

func Delete(id int) error {
	db, err := connection()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("delete from customers where id=$1")
	if err != nil {
		log.WithError(err).Error("can't prepare statement")
		return err
	}
	if _, err := stmt.Exec(id); err != nil {
		log.WithError(err).Error("delete failed")
		return err
	}
	return nil
}

func Insert(c Customer) (Customer, error) {
	db, err := connection()
	if err != nil {
		return Customer{}, err
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into customers (name,email,status) values ($1,$2,$3) returning id")
	if err != nil {
		log.WithError(err).Error("can't prepare statement")
		return Customer{}, err
	}
	row := stmt.QueryRow(c.Name, c.Email, c.Status)
	err = row.Scan(&c.Id)
	if err != nil {
		return Customer{}, err
	}
	return c, nil

}
