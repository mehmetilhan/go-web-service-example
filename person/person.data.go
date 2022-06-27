package person

import (
	"database/sql"
	"log"
	"mehmet.com/database"
)

func getPersonList() ([]Person, error) {

	var personList []Person

	results, err := database.DbConn.Query("select * from persons")

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	var person Person
	for results.Next() {
		results.Scan(&person.Id, &person.Name, &person.Age, &person.Address)
		personList = append(personList, person)
	}
	return personList, err
}

func getPerson(id string) (*Person, error) {
	result := database.DbConn.QueryRow("select * from persons where id = ?", id)
	person := &Person{}
	err := result.Scan(&person.Id, &person.Name, &person.Age, &person.Address)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return person, nil
}

func addPerson(person *Person) error {
	_, err := database.DbConn.Exec("insert into persons (name,age,address) values (?,?,?)",
		&person.Name, &person.Age, &person.Address)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func updatePerson(person *Person, id string) error {
	_, err := database.DbConn.Exec("update persons set name = ?,age=?,address=? where id=?",
		&person.Name, &person.Age, &person.Address, id)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func deletePerson(id string) (int, error) {
	rows, err := database.DbConn.Exec("delete from persons where id = ?", id)

	count, countError := rows.RowsAffected()
	if count < 1 {
		return 0, countError
	}

	if err != nil {
		return 0, err
	}
	return int(count), nil
}
