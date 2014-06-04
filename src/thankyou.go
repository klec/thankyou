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
	"strconv"
)


func main() {
	var App = new(Application)
	App.GetMongoConnection() //@todo
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	//http.HandleFunc("/addpers", app.addPersonal)
	http.HandleFunc("/index.html", App.handll)
	//http.Post("/addreview", app.addreview)
	http.ListenAndServe(":8080",nil)
}

type Application struct {
	DbSource *mgo.Database
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
	MasterIp 	string
}

type Template struct {
	Body string
	Level int
}

type Page struct {
	Best1, Best2, Best3	string
	Loosers			map[int]string
	layout			*template.Template
	SlavesOptions	template.HTML
}

func (a *Application)getReviews()  {
	pipeline:=[]bson.M{
		bson.M{"$group":
			bson.M{	"_id": "$level","summ":	bson.M{"$sum": 1,},	},
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

	err := a.DbSource.Run(query, &answer)
	if nil!=err {
		fmt.Println( err)
	}
	for row:= range answer.Result {
		fmt.Println(answer.Result[row])

	}

	//return answer.Result
}

func (a *Application) handll(writer http.ResponseWriter, request *http.Request) {
	if (request.Method == "POST") {
//		a.addReview(request)
	}
	//var reviews = []Review{}
	//a.getReviews()

	slaves:=a.GetSlaves()
	page := &Page{"первый","второй","третий", nil, nil, template.HTML("")}
	fmt.Println(slaves)
	page.Render(writer)
}

func (a *Application) addReview(request *http.Request){
	slave, _  := strconv.Atoi(request.FormValue("slave"))
	action := request.FormValue("action")
	master, _ := strconv.Atoi(request.FormValue("master"))
	review:=&Review{Slave:slave, Action:action, Master: master, MasterIp: request.RemoteAddr}
	a.DbSource.C("reviews").Insert(review) //@todo add creation time
	//@todo add slave notification

}


func (a *Application) addPersonal(writer http.ResponseWriter, request *http.Request){
	//fmt.Println(connection)

	file, _ := os.Open("personal.html")
	scanner := bufio.NewScanner(file)
	i :=0
	for scanner.Scan() {
		i++;
		row:=scanner.Text()
		roww := strings.Split(row, "\t")
		p:=&Person{ID: i, Name: roww[0], Email: roww[1]}
//		a.DbSource.C("persons").Insert(p)
		fmt.Println(p)
	}
}

func (a *Application)GetSlaves() template.HTML{
	html := "<option>Имя нашего героя</option>"
	res:=[]Person{}
	iter := a.DbSource.C("persons").Find(nil).Iter()
	err:=iter.All(&res)
	if(err!=nil){fmt.Println(err)}
	for i:=range res{
		person := res[i]
		//fmt.Println(person.Name)
		html+="<option value=\""+strconv.Itoa(i)+"\" >"+person.Name+"</option>"
	}
	return template.HTML(html)
}

func (p *Page)Render(writer http.ResponseWriter) {
	layout, err := template.ParseFiles("answer.html")
	if (err != nil) {fmt.Println(err)}
	p.layout = layout

	err = p.layout.ExecuteTemplate(writer, "answer", p)
	if nil!=err {fmt.Println( err)}
}

func (a *Application) GetMongoConnection() *mgo.Database {
	connection, err := mgo.Dial("localhost")
	if(err!=nil){fmt.Println(err)}
	db := connection.DB("thanks")
	a.DbSource = db
	return db;
}
