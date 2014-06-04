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
	"net"
	"strconv"
)


func main() {
	var app = new(App)
	//app.init()
	http.Handle("/static/", http.FileServer(http.Dir("./")))
//	http.HandleFunc("/addpers", app.addPersonal)
	http.HandleFunc("/index.html", app.handll)
	//http.Post("/addreview", app.addreview)
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

type Review struct {
	Slave		int
	Action		string
	Master		int
	MasterIp 	net.IP
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
		a.addReview(request)


	}
	a.render(writer)
}

func (a *App) addReview(request *http.Request){
	slave, _  := strconv.Atoi(request.FormValue("slave"))
	action := request.FormValue("action")
	master, _ := strconv.Atoi(request.FormValue("master"))
	connection:=GetMongoConnection()
	review:=&Review{Slave:slave, Action:action, Master: master, MasterIp: net.IPv4(127,0,0,1)}
	connection.C("reviews").Insert(review)
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
	html := "<option>Имя нашего героя</option>"
	connection:=GetMongoConnection()
	res:=[]Person{}
	iter := connection.C("persons").Find(nil).Iter()
	err:=iter.All(&res)
	if(err!=nil){fmt.Println(err)}
	for i:=range res{
		person := res[i]
		//fmt.Println(person.Name)
		html+="<option value="+strconv.Itoa(i)+">"+person.Name+"</option>"
	}
	html+="</option>"
	return template.HTML(html)
}

func GetMongoConnection() *mgo.Database {
	connection, err := mgo.Dial("localhost")
	if(err!=nil){fmt.Println(err)}
	db := connection.DB("thanks")

	return db;
}
