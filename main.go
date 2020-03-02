package main

import(
	"log"
	"net/http"
	"html/template"
	"io/ioutil"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
)


type Data struct {
	A,B int
}

type Result struct{
	Result interface{}
}

func main(){
	router := httprouter.New()
	router.POST("/",calculate)
	router.GET("/calculate", myFunc)
	log.Fatal(http.ListenAndServe(":8989", router))
}

func myFunc(w http.ResponseWriter, req *http.Request, _ httprouter.Params){
	t,err := template.ParseFiles("index.html")
	if err!=nil {
		panic(err)
	}
	t.Execute(w,nil)
	
}

func calculate(w http.ResponseWriter, req *http.Request, _ httprouter.Params){
		var d Data
		data,_ := ioutil.ReadAll(req.Body)
		err := json.Unmarshal(data, &d)
		if err!=nil{
			panic(err)
		}
		if !validate(d){
			result := Result{
				Result : "Incorrect Input",
			}
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusBadRequest)
			resultJson, err := json.Marshal(result)
			if err!=nil{
				panic(err)
			}
			w.Write(resultJson)
			return
		}
		ch := make(chan int)
		go func(){
			ch <- fact(d.A)
		}()
		go func(){
			ch <- fact(d.B)
		}()
		ans := <-ch * <- ch
		result := Result{
			Result : ans,
		}
		resultJson, err := json.Marshal(result)
		if err!=nil{
			panic(err)
		}
		
		w.Header().Set("Content-Type","application/json")
		w.Write(resultJson)
} 

func fact(n int)int{
	if n == 0 {
		return 1
	}
	return n*fact(n-1)
}

func validate(d Data)bool{
	if d.A <=0 || d.B <= 0 {
		return false
	}
	return true
}