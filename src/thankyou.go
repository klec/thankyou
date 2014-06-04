package main

import (
	"net/http"
	"html/template"
	"fmt"
	"labix.org/v2/mgo"
	"bufio"
	"os"
	"strings"
	"labix.org/v2/mgo/bson"
)


func main() {
	var app = new(App)
	//app.init()
	http.Handle("/static/", http.FileServer(http.Dir("./")))
//	http.HandleFunc("/addpers", app.addPersonal)
	http.HandleFunc("/index.html", app.handll)
	http.ListenAndServe(":8080",nil)
}

type App struct {

	layout *template.Template
}

type Person struct {
	ID        int
	Name      string
	Email     string
}

type Template struct {
	Body string
	Level int
}

type Page struct {
	Best1, Best2, Best3	string
	Loosers			map[int]string
}

func (p *Page)init() {
	connection:=GetMongoConnection()


	pipeline:=[]bson.M{
		bson.M{"$group":
			bson.M{
				"_id": "$level",
				"summ":	bson.M{"$sum": 1,},
			},
		},
		bson.M{"$sort":
			bson.M{ "summ": -1},
		},

	}

	query := bson.D{
		{"aggregate","templates"},
		{"pipeline", pipeline},
	}

	answer := struct {
		Result []map[string] int
		Ok     bool
	}{}

	err := connection.Run(query, &answer)
	if nil!=err {
		fmt.Println( err)
	}
	for row:= range answer.Result {
		fmt.Println(answer.Result[row])

	}

}

func (a *App) handll(writer http.ResponseWriter, request *http.Request) {
	if (request.Method == "POST") {

	}
	a.render(writer)
}

func (a *App)render(writer http.ResponseWriter) {
	layout, err := template.ParseFiles("answer.html")
	if (err != nil) {fmt.Println(err)}
	a.layout = layout
	page := &Page{"первый","второй","третий", nil}
	page.init();
	err = a.layout.ExecuteTemplate(writer, "answer", page)
	if nil!=err {fmt.Println( err)}
}

func (a *App)addPersonal(writer http.ResponseWriter, request *http.Request){
	connection:=GetMongoConnection()
	fmt.Println(connection)

	file, _ := os.Open("personal.html")
	scanner := bufio.NewScanner(file)
	i :=0
	for scanner.Scan() {
		i++;
		row:=scanner.Text()
		roww := strings.Split(row, "\t")
		p:=&Person{ID: i, Name: roww[0], Email: roww[1]}
		//connection.C("persons").Insert(p)
		fmt.Println(p)
	}
}

func (p *Page)GetSlaves() template.HTML{
	html := template.HTML("<option>Имя нашего героя</option>")
	connection:=GetMongoConnection()
	res:=[]Person{}
	iter := connection.C("persons").Find(nil).Iter()
	err:=iter.All(&res)
	if(err!=nil){fmt.Println(err)}
	for i:=range res{
		person := res[i]
		//fmt.Println(person.Name)
		html+=template.HTML("</option><option>"+person.Name)
	}
	html+="</option>"
	return html
}

func GetMongoConnection() *mgo.Database {
	connection, err := mgo.Dial("localhost")
	if(err!=nil){fmt.Println(err)}
	db := connection.DB("thanks")

	return db;
}
