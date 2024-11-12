package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados
func InitDB() error {
	var err error
	DB, err = sql.Open("mysql", "caiordgs:HokagE123!@tcp(127.0.0.1:3306)/dynastytracker")
	if err != nil {
		return err
	}

	// Verifica a conexão
	err = DB.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Conexão com o banco de dados estabelecida!")
	return nil
}
