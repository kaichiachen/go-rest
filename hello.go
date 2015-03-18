package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "database/sql"
    
    "github.com/drone/routes"
    _ "github.com/go-sql-driver/mysql"
)

func main() {

    mux := routes.New()
    mux.Post("/soup",SoupList)
    mux.Post("/mainmeal",MainMealList)
    mux.Post("/order",Order)
    mux.Post("/order/add",OrderAdd)

    http.Handle("/", mux)
    http.ListenAndServe(":3340", nil)

    select {}
    fmt.Println("Run Server on port 3340")
}



func SoupList(w http.ResponseWriter, req *http.Request){
	db, err := sql.Open("mysql", "root:yzucse@/order")
	rows, err := db.Query("SELECT * FROM soup")

	type Soup struct {
		ID int
		Name string
		Cost int
	}

	type SoupSlice struct {
		Soups []Soup
	}

	var s SoupSlice

    for rows.Next(){
    	var id ,cost int
    	var name string
    	err = rows.Scan(&id,&name,&cost)
    	s.Soups = append(s.Soups,Soup{ID:id, Name:name, Cost:cost})
    }
    b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprint(w,string(b))
}

func MainMealList(w http.ResponseWriter, req *http.Request){
	db, err := sql.Open("mysql", "root:yzucse@/order")
	rows, err := db.Query("SELECT * FROM main_meal")

	type MainMeal struct {
		ID int
		Name string
		Cost int
	}

	type MainMealSlice struct {
		MainMeals []MainMeal
	}

	var m MainMealSlice

    for rows.Next(){
    	var id ,cost int
    	var name string
    	err = rows.Scan(&id,&name,&cost)
    	m.MainMeals = append(m.MainMeals,MainMeal{ID:id, Name:name, Cost:cost})
    }
    b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprint(w,string(b))
}

func OrderAdd(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	db, err := sql.Open("mysql", "root:yzucse@/order")
	checkErr(err)

	rows,_ := db.Query("SELECT * from eveorder where id = ?",req.Form["ID"][0])
	count := 0
	for rows.Next(){
		count++
	}
	if count == 0 {
		stmp, _ := db.Prepare("INSERT eveorder SET id=?,main_meal=?,soup=?,sCount=?,mCount=?")
		_,error := stmp.Exec(req.Form["ID"][0],req.Form["mainmeal"][0],req.Form["soup"][0],req.Form["sCount"][0],req.Form["mCount"][0])
		checkErr(error)
		fmt.Fprintln(w,"{\"valid\":true}")
	} else {
		stmp, _ := db.Prepare("UPDATE eveorder SET main_meal=?,soup=?,sCount=?,mCount=? where id = ?")
		_,error := stmp.Exec(req.Form["mainmeal"][0],req.Form["soup"][0],req.Form["sCount"][0],req.Form["mCount"][0],req.Form["ID"][0])
		checkErr(error)
		fmt.Fprintln(w,"{\"valid\":true}")
	}

}

func Order(w http.ResponseWriter, req *http.Request){
	db, err := sql.Open("mysql", "root:yzucse@/order")
	checkErr(err)
	rows, err := db.Query("SELECT * FROM eveorder")
	checkErr(err)

	type Order struct {
		ID int
		MainMeal string
		Soup string
		MCount int
		SCount int
		Cost int
	}

	type OrderList struct {
		Orders []Order
	}

	var o OrderList

	for rows.Next(){
    	var id ,main_meal ,soup ,mCount ,sCount int
    	err := rows.Scan(&id,&main_meal,&soup,&mCount,&sCount)
    	checkErr(err)
    	sum := 0
        mstmp, _ := db.Prepare("SELECT * from main_meal where id = ?")
    	var Mname string
    	var Mcost int
    	var mid int
    	err = mstmp.QueryRow(main_meal).Scan(&mid,&Mname,&Mcost)
    	sum += Mcost * mCount

    	sstmp, _ := db.Prepare("SELECT * from soup where id = ?")
    	var Sname string
    	var Scost int
    	var sid int
    	err = sstmp.QueryRow(soup).Scan(&sid,&Sname,&Scost)
    	sum += Scost * sCount

    	o.Orders = append(o.Orders, Order{ID:id, Cost:sum, MCount:mCount, SCount:sCount, MainMeal:Mname, Soup:Sname})
    }
    b, err := json.Marshal(o)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprint(w,string(b))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

