package main

import (
	"database/sql"  // Usado para realizar consultas.
	"log"           // Usado para mostrar mensagens no console.
	"net/http"      // Usado para gerenciar URLs e o servidor web.
	"text/template" // Usado para gerenciar templates.

	_ "github.com/go-sql-driver/mysql" // Driver MySQL para Go.
)

// Struct que vai ser usado para exibir os dados no template.
type Names struct {
	Id    int    // Identificador único.
	Name  string // Nome.
	Email string // E-mail.
}

// Função que abre a conexão com o banco de dados.
func dbConn() (db *sql.DB) {
	dbDriver := "mysql" // Nome do SGBD.
	dbUser := "root"    // Usuário.
	dbPass := ""        // Senha.
	dbName := "crud"    // Nome do banco de dados.

	// Abre uma conexão com o banco de dados.
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	// Verifica se houve erros.
	if err != nil {
		// Encerra a aplicação em caso de erros.
		panic(err.Error())
	}
	// Retorna a conexão com o banco de dados.
	return db
}

// tmpl redenriza todos os templates da pasta "tmpl"
// independete da extensão.
var tmpl = template.Must(template.ParseGlob("tmpl/*"))

// Função usada para renderizar o arquivo Index.
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println("Abrindo pagina principal.")
	// Abre a conexão com o banco.
	db := dbConn()
	// Realiza a consulta no banco e trata os erros.
	selDB, err := db.Query("SELECT * FROM `names` ORDER BY `id` DESC")
	if err != nil {
		panic(err.Error())
	}
	// Monta a struct para ser utilizada no template.
	n := Names{}
	// Monta um array para guardar os valores da struct.
	res := []Names{}
	// Percorre todos os dados retornados com um loop.
	for selDB.Next() {
		// Declara as variáveis que vão armazenar os dados.
		var id int
		var name, email string
		// Faz o Scan do SELECT.
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		// Envia os resultados para a struct.
		n.Id = id
		n.Name = name
		n.Email = email
		// Adiciona a struct no array.
		res = append(res, n)
	}
	// Abre a página Index e exibe todos os registros na tela.
	tmpl.ExecuteTemplate(w, "Index", res)
	// Fecha a conexão.
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	log.Println("Abrindo pagina de visualização.")
	db := dbConn()
	// Pega o id do parâmetro da URL.
	nId := r.URL.Query().Get("id")
	// Consulta o registro pelo id e trata os erros.
	selDB, err := db.Query("SELECT * FROM `names` WHERE `id` = ?", nId)
	if err != nil {
		panic(err.Error())
	}
	// Monta a struct para ser utilizada no template.
	n := Names{}
	for selDB.Next() {
		// Declaração das variáveis que vão guardar os dados do selDB.
		var id int
		var name, email string
		// Faz o Scan do SELECT.
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		// Envia os dados para a struct.
		n.Id = id
		n.Name = name
		n.Email = email
	}
	// Mostra o template.
	tmpl.ExecuteTemplate(w, "Show", n)
	// Fecha a conexão.
	defer db.Close()
}

// A função New apenas exibe o formulário para inserir novos dados.
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	log.Println("Abrindo pagina de atualização.")
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM `names` WHERE `id` = ?", nId)
	if err != nil {
		panic(err.Error())
	}
	n := Names{}
	for selDB.Next() {
		var id int
		var name, email string
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		n.Id = id
		n.Name = name
		n.Email = email
	}
	tmpl.ExecuteTemplate(w, "Edit", n)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserindo novos dados.")
	db := dbConn()
	// Verifica qual é o método do formulário passado.
	if r.Method == "POST" {
		// Pega os dados do formulário.
		name := r.FormValue("name")
		email := r.FormValue("email")
		// Prepara a SQL e verifica erros.
		insForm, err := db.Prepare("INSERT INTO `names` (name, email) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		// Insere os dados do formulário na SQL tratada e verifica erros.
		insForm.Exec(name, email)
		// Exibe um log com os valores digitados no formulário.
		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}
	defer db.Close()
	// Retorna para o Index.
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Atualizando dados.")
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		updForm, err := db.Prepare("UPDATE `names` SET `name` = ?, `email` = ? WHERE `id` = ?")
		if err != nil {
			panic(err.Error())
		}
		updForm.Exec(name, email, id)
		log.Println("UPDATE: Name: " + name + " | E-mail: " + email)
	}
	defer db.Close()
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Removendo dados.")
	db := dbConn()
	nId := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM `names` WHERE `id` = ?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(nId)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
	// Informa que o servidor foi iniciado.
	log.Println("Server started on: http://localhost:9000")

	// Gerenciamento das URLs.
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	// Ações.
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	http.ListenAndServe("localhost:8080", nil)
}
