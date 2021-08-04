package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/QianMason/drone-backend/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalcount = 0

// var client *mongo.Client
// var ctx context.Context

type DroneStruct struct {
	Data []models.DroneSim
}

type DroneDB struct {
	username string
	password string
}

func setupCors(w *http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Header.Get("Origin"))
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CRSF-Token, Authorization")
}

func NewRouter() *mux.Router {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading environment variables")
	}
	mongoPass := os.Getenv("MONGOPASS")
	mongoUser := os.Getenv("MONGOUSER")
	db := &DroneDB{username: mongoUser, password: mongoPass}
	// client, ctx = models.GetClient()
	r := mux.NewRouter()
	// r.HandleFunc("/", middleware.AuthRequired(indexHandler)).Methods("GET")
	r.HandleFunc("/register", db.registerHandler).Methods("POST")
	r.HandleFunc("/drones", db.droneHandler).Methods("GET", "OPTIONS")
	//r.HandleFunc("/thingspeak", tsHandler).Methods("GET")
	//r.HandleFunc("/", indexHandler).Methods("GET")
	//r.HandleFunc("/", postHandler).Methods("POST")
	//fs := http.FileServer(http.Dir("./build"))
	//r.PathPrefix("/build").Handler(http.StripPrefix("/build", fs))
	return r
}

func (db *DroneDB) droneHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in index handler")
	setupCors(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + db.username + ":" + db.password + "@cluster0.14i4y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("DronePlatform").Collection("simData")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(cur)
		fmt.Println("error with cur")
		fmt.Println(err)
		return
	}
	drones := []models.Drone{}

	for cur.Next(ctx) {
		d := models.Drone{}
		err = cur.Decode(&d)
		fmt.Println("d lat:", d.Coordinates.Lat, ":d lng:", d.Coordinates.Lng)
		if err != nil {
			fmt.Println(err)
			return
		}
		drones = append(drones, d)
	}
	cur.Close(ctx)
	if len(drones) == 0 {
		w.WriteHeader(500)
		w.Write([]byte("No data found."))
		return
	}
	jsn, err := json.Marshal(drones)
	fmt.Println("jsn:", jsn)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsn)
}

//
// func (d *DroneDB) indexHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("in index handler")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	dd := DroneStruct{
// 		Data: make([]models.DroneSim, 0),
// 	}
// 	clientOptions := options.Client().
// 		ApplyURI("mongodb+srv://thunderpurtz:" + d.password + "@cluster0.14i4y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err = client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	collection := client.Database("DronePlatform").Collection("simData")
// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		fmt.Println(cur)
// 		fmt.Println("error with cur")
// 		fmt.Println(err)
// 		return
// 	}

// 	for cur.Next(ctx) {
// 		d := models.Drone{}
// 		err = cur.Decode(&d)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		dd.Data = append(dd.Data, d)
// 	}
// 	cur.Close(ctx)
// 	if len(dd.Data) == 0 {
// 		w.WriteHeader(500)
// 		w.Write([]byte("No data found."ctx, cancel := context.WithTimeout(contex
// 	jsn, err := json.Marshal(dd.Data)
// 	fmt.Println("in index handler 2")
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(jsn)
// }

func getUUID() string {
	uid := strings.Replace(uuid.New().String(), "-", "", -1)
	fmt.Println("New UUID:", uid)
	return uid
}

func (d *DroneDB) registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("register handler called")

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + d.username + ":" + d.password + "@cluster0.14i4y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	rBody := string(body)
	register := &models.RegisterDrone{}

	err = json.Unmarshal([]byte(rBody), register)
	if err != nil {
		fmt.Println("in here error unmarshalling:", err)
		return
	}

	dID := getUUID()

	drone := models.Drone{
		ID:      primitive.NewObjectID(),
		DroneID: dID,
		Coordinates: models.Coordinates{
			Lat: register.Lat,
			Lng: register.Lng,
		},
		Address:     register.Address,
		LastUpdated: time.Now().UTC(),
	}

	collection := client.Database("DronePlatform").Collection("simData")

	res, err := collection.InsertOne(ctx, drone)
	if err != nil {
		fmt.Println("error in insert")
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(res.InsertedID)
	fmt.Println("exiting post handler")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(dID))

}

// func userGetHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("userget" + strconv.Itoa(globalcount))
// 	globalcount += 1
// 	session, _ := sessions.Store.Get(r, "session")
// 	untypedUserId := session.Values["user_id"]
// 	currentUserId, ok := untypedUserId.(int64)
// 	fmt.Println(currentUserId)
// 	if !ok {
// 		utils.InternalServerError(w)
// 		return
// 	}
// 	vars := mux.Vars(r) //hashmap of variable names and content passed for that variable
// 	username := vars["username"]
// 	fmt.Println("username", username)

// 	currentPageUserString := strings.TrimLeft(r.URL.Path, "/")
// 	currentPageUser, err := models.GetUserByUsername(currentPageUserString)
// 	if err != nil {
// 		utils.InternalServerError(w)
// 		return
// 	}
// 	currentPageUserID, err := currentPageUser.GetId()
// 	if err != nil {
// 		utils.InternalServerError(w)
// 		return
// 	}
// 	updates, err := models.GetUpdates(currentPageUserID)
// 	if err != nil {
// 		utils.InternalServerError(w)
// 		return
// 	}

// 	utils.ExecuteTemplate(w, "index.html", struct {
// 		Title       string
// 		Updates     []*models.Update
// 		DisplayForm bool
// 	}{
// 		Title:       username,
// 		Updates:     updates,
// 		DisplayForm: currentPageUserID == currentUserId,
// 	})

// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	updates, err := models.GetAllUpdates()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 		return
// 	}
// 	utils.ExecuteTemplate(w, "index.html", struct {
// 		Title       string
// 		Updates     []*models.Update
// 		DisplayForm bool
// 	}{
// 		Title:       "All updates",
// 		Updates:     updates,
// 		DisplayForm: true,
// 	})
// 	fmt.Println("get")
// }

// func postHandlerHelper(w http.ResponseWriter, r *http.Request) error {
// 	session, _ := sessions.Store.Get(r, "session")
// 	untypedUserID := session.Values["user_id"]
// 	userID, ok := untypedUserID.(int64)
// 	if !ok {
// 		return utils.InternalServer
// 	}
// 	currentPageUserString := strings.TrimLeft(r.URL.Path, "/")
// 	currentPageUser, err := models.GetUserByUsername(currentPageUserString)
// 	if err != nil {
// 		return utils.InternalServer
// 	}
// 	currentPageUserID, err := currentPageUser.GetId()
// 	if err != nil {
// 		return utils.InternalServer
// 	}
// 	if currentPageUserID != userID {
// 		return utils.BadPostError
// 	}
// 	r.ParseForm()
// 	body := r.PostForm.Get("adddrone")
// 	fmt.Println(body)
// 	err = models.PostUpdates(userID, body)
// 	if err != nil {
// 		return utils.InternalServer
// 	}
// 	return nil
// }

// func postHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("post handler called")
// 	err := postHandlerHelper(w, r)
// 	if err == utils.InternalServer {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 	}
// 	http.Redirect(w, r, "/", 302)
// }

// func UserPostHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("user post handler called")
// 	fmt.Println(r.URL.Path)
// 	err := postHandlerHelper(w, r)
// 	if err == utils.BadPostError {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Cannot write to another user's page"))
// 	}
// 	http.Redirect(w, r, r.URL.Path, 302)
// }

// func loginGetHandler(w http.ResponseWriter, r *http.Request) {
// 	utils.ExecuteTemplate(w, "login.html", nil)
// }

// func loginPostHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.PostForm.Get("username")
// 	password := r.PostForm.Get("password")

// 	user, err := models.AuthenticateUser(username, password)
// 	if err != nil {
// 		switch err {
// 		case models.InvalidLogin:
// 			utils.ExecuteTemplate(w, "login.html", "User or Pass Incorrect")
// 		default:
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("Internal server error"))
// 		}
// 		return
// 	}
// 	userId, err := user.GetId()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 		return
// 	}
// 	sessions.GetSession(w, r, "session", userId)
// 	http.Redirect(w, r, "/", 302)
// }

// func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
// 	sessions.EndSession(w, r)
// 	http.Redirect(w, r, "/login", 302)
// }

// func registerGetHandler(w http.ResponseWriter, r *http.Request) {
// 	utils.ExecuteTemplate(w, "register.html", nil)
// }

// func registerPostHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.PostForm.Get("username")
// 	password := r.PostForm.Get("password")
// 	err := models.RegisterUser(username, password)
// 	if err == models.UserNameTaken {
// 		utils.ExecuteTemplate(w, "register.html", "username taken")
// 		return
// 	}
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal server error"))
// 		return
// 	}
// 	http.Redirect(w, r, "/login", 302)
// }
