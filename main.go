package main

import (
	"fmt"
    "log"
	"net/http"
	"encoding/json"	
	"github.com/gorilla/mux"
	"io/ioutil"
	"os"

)
//OMDB responselarındaki data modelleri
type Movie struct{
	Title string `json:"Title"`
    Year string `json:"Year"`
	imdbID int `json:"imdbID"`
	Type string `json:"Type"`
    Poster string `json:"Poster"`

}
type Movies struct{
	Search []Movie `json:"Search"`
	totalResults int `json:"totalResults"`
	Response string `json:"Response"`
	SearchQuery string `json:"SearchQuery"`

}
	//Anasayfa
func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Film aramak için localhost:1000/search/{search input} kullanabilirsiniz")
    fmt.Println("Endpoint Hit: Anasayfa")
}
	//Search Sayfası Controller ı
func searchPage(w http.ResponseWriter, r *http.Request)  {

	//Geri göndereceğimiz response JSON olacak
	w.Header().Set("Content-Type", "application/json")
	
	//Kullanıcınıın url de kullandığı {search_query}'i mux kullanarak query değişkenine atıyoruz
	vars:=mux.Vars(r)
	query:=vars["search_query"]

	//"query" anahtarı ile bir cache var mı kontrolü
	cache_data := Movies{}
	if(fileExists("cache.json")){
		file, _ := ioutil.ReadFile("cache.json")
		_ = json.Unmarshal([]byte(file), &cache_data)
	}
	if(cache_data.SearchQuery==query){
			json.NewEncoder(w).Encode(cache_data)
			fmt.Println("Cache bulundu:")
	}else{
		//OMDB API'dan sonuçları "resp" değişkenine istiyoruz, herhangi bir error varsa loglanıyor
		resp, err := http.Get("http://www.omdbapi.com/?s="+query+"&apikey=db626cbb")
		if err != nil {
			log.Fatalln(err)
		}
		
		defer resp.Body.Close()

		// m adında Movies türünde bir değişkene OMDBden aldığımız dataları atıyoruz.
		var m Movies
		json.NewDecoder(resp.Body).Decode(&m)
		//Consoleda da görmek için print
		fmt.Println("Cache bulunamadı")
		m.SearchQuery=query
		//w Reponsewriterına kendi değişkenimiz m üzerinden kullanıcıya serve
		json.NewEncoder(w).Encode(m)
		file, _ := json.MarshalIndent(m, "", " ")
 
		_ = ioutil.WriteFile("cache.json", file, 0644)

	}
	
}
//Endpointleri controllerlara dağıtmak için router
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/search/{search_query}",searchPage)
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}


//Main
func main()  {
    handleRequests()
}