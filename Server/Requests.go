package Server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func GetAchivments(w http.ResponseWriter, r *http.Request) {
	var achivments []Achivment
	if err := DB.Find(&achivments).Error; err != nil {
		http.Error(w, "Failed to retrieve achievments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(achivments)
	if err != nil {
		return
	}
}

func CreateAchivment(w http.ResponseWriter, r *http.Request) {
	var achivment Achivment
	if err := json.NewDecoder(r.Body).Decode(&achivment); err != nil {
		http.Error(w, "Failed to create achivment", http.StatusInternalServerError)
		return
	}
	DB.Create(&achivment)
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(achivment)
	if err != nil {
		return
	}
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/achievments", GetAchivments).Methods(http.MethodGet)
	r.HandleFunc("/achievments", CreateAchivment).Methods(http.MethodPost)
	r.HandleFunc("/winner", PostWinnerHandler).Methods(http.MethodPost)
}

type ReceivedData struct {
	Name     string
	Kills    int
	Killable bool
}

func PostWinnerHandler(w http.ResponseWriter, r *http.Request) {

	// Чтение данных из тела запроса

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	var d ReceivedData
	query, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Unable to parse query")
		return
	}

	// Извлечение значений полей из запроса
	killable := query.Get("killable")
	name := query.Get("name")
	kills, _ := strconv.Atoi(query.Get("kills"))
	d.Killable, _ = strconv.ParseBool(killable)
	d.Name = name
	d.Kills = kills

	fmt.Println("Received winner:", d.Name, "Kills IS ", d.Kills, d.Killable)

	// Ответ клиенту
	w.Write([]byte("Winner received"))

	achiev, champion := AddCharacterToAchievment(&d)
	if err := DB.Where("name = ?", achiev.Name).First(&achiev).Error; err != nil {
		fmt.Println("Error loading achievment:", err)
		http.Error(w, "Achievment not found", http.StatusNotFound)
		return
	}
	if !contains(achiev.Heroes, d.Name) {
		achiev.Heroes = append(achiev.Heroes, d.Name)
	}

	fmt.Println("Achievment:", achiev, "Heroes:", champion)
	if err := DB.Model(&achiev).Updates(Achivment{Heroes: achiev.Heroes}).Error; err != nil {
		fmt.Println("Error updating achievment:", err)
		http.Error(w, "Unable to update achievment", http.StatusInternalServerError)
		return
	}
	//DB.Save(&achiev)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func AddCharacterToAchievment(data *ReceivedData) (*Achivment, []string) {

	var names []string
	names = append(names, data.Name)
	if !contains(names, data.Name) {
		names = append(names, data.Name)
	}
	if !data.Killable {
		switch data.Kills {
		case 1:
			fmt.Println("FB")
			return &Achivment{
				Name:   "First blood",
				Heroes: names,
			}, names
		case 2:
			fmt.Println("DOUBLE KILL")
			return &Achivment{
				Name:   "Double kill",
				Heroes: names,
			}, names
		case 3:
			fmt.Println("TRIPPLE")
			return &Achivment{

				Name:   "Tripple kill",
				Heroes: names,
			}, names

		}
	} else {
		return &Achivment{
			Name:   "First unluck",
			Heroes: append(names, data.Name),
		}, names
	}
	return &Achivment{
		Name:   "Unstoppable",
		Heroes: append(names, data.Name),
	}, names

}
