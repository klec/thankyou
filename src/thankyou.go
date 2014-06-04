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
	Time		string
}

type Template struct {
	Body string
	Level int
}

type Page struct {
	Best1, Best2, Best3	string
	Loosers			[8]string
	layout			*template.Template
	SlavesOptions	template.HTML
}

func (a *Application)getReviews() [8]string {
	query := bson.D{
		{"aggregate","reviews"},
		{"pipeline", []bson.M{
			bson.M{"$sort":
				bson.M{ "action": 1},
			},

			bson.M{"$group":
				bson.M{	"_id": "$slave",
					"summ":	bson.M{"$sum": 1},
					"slave": bson.M{"$first":"$slave"},
					"action":bson.M{"$first":"$action"},
	//				"masterip":bson.M{"$first":"$masterip"},
				},
			},
			bson.M{"$sort":
				bson.M{ "summ": -1},
			},
			bson.M{"$limit":8 },
		}},
	}

	answer := struct {
		Result []map[string] interface {}
		Ok     bool
	}{}

	err := a.DbSource.Run(query, &answer)
	if nil!=err {
		fmt.Println( err)
	}
	var Reviews = [8]string{}
	var lastTemplate []string
	for row:= range answer.Result {
		action := answer.Result[row]["action"].(string)

		template := a.GetTemplate(answer.Result[row]["summ"].(int), lastTemplate)
		lastTemplate = append(lastTemplate, template)
		person := a.GetSlave(answer.Result[row]["slave"].(int))
		oneReview:= strings.Replace(template, "Он", person, -1)
		oneReview = strings.Replace(oneReview, "Помог", action, -1)

		Reviews[row]=oneReview
	}
	return Reviews
}

func (a *Application) handll(writer http.ResponseWriter, request *http.Request) {
	if (request.Method == "POST") {
		a.addReview(request)
	}
	loosers:=a.getReviews()

	slaves:=a.GetSlaves()
	page := &Page{loosers[0], loosers[1], loosers[2], loosers, nil, slaves}
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

func (a *Application)GetSlave(id int) string{
	res:=[]Person{}
	iter := a.DbSource.C("persons").Find(bson.M{"id":id})
	err:=iter.All(&res)
	if(err!=nil){fmt.Println(err)}
	p:=res[0]
	return p.Name
}
func (a *Application)GetTemplate(id int, ne []string) string{
	res:=[]Template{}
	iter := a.DbSource.C("templates").Find(bson.M{"level":id, "body":bson.M{"$nin":ne}})
	err:=iter.All(&res)
	if(err!=nil){fmt.Println(err)}
	t:=res[0]
	return t.Body
}

func (a *Application)GetSlaves() template.HTML{
	html := "<option>Имя нашего героя</option>"
	res:=[]Person{}
	iter := a.DbSource.C("persons").Find(nil).Sort("id").Iter()
	err:=iter.All(&res)
	if(err!=nil){fmt.Println(err)}
	for i:=range res{
		person := res[i]
		html+="<option value=\""+strconv.Itoa(person.ID)+"\" >"+person.Name+"</option>"
	}
	return template.HTML(html)
}

func (p *Page)GetLoosers() template.HTML{
	html:=""
	for i:=range p.Loosers{
		if(i>2){
			person := p.Loosers[i]
			html+="<li>"+person+"</li>\n"
		}
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
