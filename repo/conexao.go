package repo

import (
	"fmt"
	"os"

	"github.com/jeansferreira/app2/model"
	"github.com/jeansferreira/app2/tratamento"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS public."Mercado"
(
    cpf_cnpj_comprador character varying(14) COLLATE pg_catalog."default" NOT NULL,
    flg_private boolean,
    flg_incompleto boolean,
    dt_ultima_compra date,
    vl_ticket_medio numeric(10,2),
    vl_ticket_ult_compra numeric(10,2),
    cnpj_loja_freq character varying(14) COLLATE pg_catalog."default" NOT NULL,
    cnpj_loja_ultima character varying(14) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "Mercado_pkey" PRIMARY KEY (cpf_cnpj_comprador)
)
`
var sqlInsert = `INSERT INTO public."Mercado"(cpf_cnpj_comprador, flg_private, flg_incompleto, 
                                                dt_ultima_compra, vl_ticket_medio, vl_ticket_ult_compra, 
                                                cnpj_loja_freq, cnpj_loja_ultima) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Postgres2019!"
	dbname   = "postgres"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//InsereDados metodo que insere os dados dos arquivo na base de dados
func InsereDados(mov []model.Compra) {

	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "example")
	host := getEnv("PG_HOST", "localhost")
	port := getEnv("PG_PORT", "5432")
	database := getEnv("POSTGRES_DB", "postgres")

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user,
		password,
		host,
		port,
		database,
	)

	db, err := sqlx.Connect("postgres", dbinfo)

	//db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password="+password+" sslmode=disable")
	if err != nil {
		fmt.Println("[repo] Erro no conexão com banco de dados. Erro:", err.Error())
	}
	defer db.Close()

	// Cria a tabela na base de dados
	db.MustExec(schema)

	tx := db.MustBegin()

	for i := 0; i < len(mov); i++ {
		//fmt.Println(">>> ", mov[i].Cpf_cnpj_comprador)
		tx.MustExec(sqlInsert, tratamento.RemoveCaracteres(mov[i].Cpf_cnpj_comprador),
			mov[i].GetFlgPrivate(),
			mov[i].GetFlgIncompleto(),
			mov[i].Dt_ultima_compra,
			mov[i].GetVlTicketMedio(),
			mov[i].GetVlTicketUltCompra(),
			tratamento.RemoveCaracteres(mov[i].Cnpj_loja_freq),
			tratamento.RemoveCaracteres(mov[i].Cnpj_loja_ultima))
	}

	tx.Commit()
}
