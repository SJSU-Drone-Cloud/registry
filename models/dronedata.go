package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Drone struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	DroneID     string             `json:"droneID" bson:"droneID"`
	Coordinates Coordinates        `json:"coordinates" bson:"coordinates"`
	Address     string             `json:"address" bson:"address"`
	LastUpdated time.Time          `json:"lastUpdated" bson:"lastUpdated"`
}

type RegisterDrone struct {
	Address string  `json:"address" bson:"address"`
	Lat     float64 `json:"lat" bson:"lat"`
	Lng     float64 `json:"lng" bson:"lng"`
}

type Coordinates struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type DroneData struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Data DroneSim           `json:"Data" bson:"Data"`
}

type DroneStruct struct {
	DroneID     string             `json:"droneID,string" bson:"droneID,string"`
	PublishedAt primitive.DateTime `json:"publishedAt" bson:"publishedAt"`
	Battery     uint16             `json:"battery,string" bson:"battery,string"`
	Height      int64              `json:"height,string" bson:"height,string"`
	Temperature int64              `json:"temperature,string" bson:"temperature,string"`
}

type DroneSim struct {
	Barometer Barometer       `json:"barometer" bson:"barometer"`
	Gps       GPS             `json:"gps" bson:"gps"`
	Imu       IMU             `json:"imu" bson:"imu"`
	State     MultiRotorState `json:"state" bson:"state"`
}

type Barometer struct {
	Time_stamp uint64  `json:"time_stamp" bson:"time_stamp"`
	Altitude   float64 `json:"altitude" bson:"altitude"`
	Pressure   float64 `json:"pressure" bson:"pressure"`
	Qnh        float64 `json:"qnh" bson:"qnh"`
}

type GPS struct {
	Gnss       GNSS   `json:"gnss" bson:"gnss"`
	Is_valid   bool   `json:"is_valid" bson:"is_valid"`
	Time_stamp uint64 `json:"time_stamp" bson:"time_stamp"`
}

type IMU struct {
	Angular_velocity    Vector3r    `json:"angular_velocity" bson:"angular_velocity"`
	Linear_acceleration Vector3r    `json:"linear_acceleration" bson:"linear_acceleration"`
	Orientation         Quaternionr `json:"orientation" bson:"orientation"`
	Time_stamp          uint64      `json:"time_stamp" bson:"time_stamp"`
}

type MultiRotorState struct {
	Collision            Collision  `json:"collision" bson:"collision"`
	Gps_location         GeoPoint   `json:"gps_location" bson:"gps_location"`
	Kinematics_estimated Kinematics `json:"kinematics_estimated" bson:"kinematics_estimated"`
	Landed_state         int        `json:"landed_state" bson:"landed_state"`
	Rc_data              RCData     `json:"rc_data" bson:"rc_data"`
}

type RCData struct {
	Is_initialized bool    `json:"is_initialized,bool" bson:"is_initialized,bool"`
	Is_valid       bool    `json:"is_valid,bool" bson:"is_valid,bool"`
	Left_z         float64 `json:"left_z" bson:"left_z"`
	Pitch          float64 `json:"pitch" bson:"pitch"`
	Right_z        float64 `json:"right_z" bson:"right_z"`
	Roll           float64 `json:"roll" bson:"roll"`
	Switches       uint64  `json:"switches" bson:"switches"`
	Throttle       float64 `json:"throttle" bson:"throttle"`
	Timestamp      uint64  `json:"timestamp" bson:"timestamp"`
	Vendor_id      string  `json:"vendor_id" bson:"vendor_id"`
	Yaw            float64 `json:"yaw" bson:"yaw"`
}

func (rd *RegisterDrone) UnmarshalJSON(data []byte) error {
	var a map[string]string
	fmt.Println("registerDrone")
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("Custom regdrone unmarshal error")
		return err
	}
	rd.Address = a["address"]
	lat, err := strconv.ParseFloat(a["lat"], 64)
	if err != nil {
		return err
	}
	lng, err := strconv.ParseFloat(a["lng"], 64)
	if err != nil {
		return err
	}
	rd.Lat = lat
	rd.Lng = lng
	return nil
}

func (co *Coordinates) UnmarshalJSON(data []byte) error {
	var a map[string]float64
	fmt.Println("coordinates")
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("Custom Coordinates unmarshal error")
		return err
	}
	co.Lat = a["lat"]
	co.Lng = a["lng"]
	return nil
}

func (rc *RCData) UnmarshalJSON(data []byte) error {
	var a map[string]interface{}
	fmt.Println("-------------------RCDATA-------------------")
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("Custom RCData unmarshal error")
		return err
	}
	rc.Is_initialized = a["is_initialized"].(bool)
	rc.Is_valid = a["is_valid"].(bool)
	rc.Left_z = a["left_z"].(float64)
	rc.Pitch = a["pitch"].(float64)
	rc.Right_z = a["right_z"].(float64)
	rc.Roll = a["roll"].(float64)
	rc.Switches = uint64(a["switches"].(float64)) //has to be float 64 because JSON unmarshal defaults to float64 for json numbers when unmarshaling into interface value
	rc.Throttle = a["throttle"].(float64)
	rc.Timestamp = uint64(a["timestamp"].(float64))
	rc.Vendor_id = a["vendor_id"].(string)
	rc.Yaw = a["yaw"].(float64)
	return nil
}

type Kinematics struct {
	Angular_acceleration Vector3r    `json:"angular_acceleration" bson:"angular_acceleration"`
	Angular_velocity     Vector3r    `json:"angular_velocity" bson:"angular_velocity"`
	Linear_acceleration  Vector3r    `json:"linear_acceleration" bson:"linear_acceleration"`
	Linear_velocity      Vector3r    `json:"linear_velocity" bson:"linear_velocity"`
	Orientation          Quaternionr `json:"orientation" bson:"orientation"`
	Position             Vector3r    `json:"position" bson:"position"`
}

type Collision struct {
	Has_collided      bool     `json:"has_collided" bson:"has_collided"`
	Impact_point      Vector3r `json:"impact_point" bson:"impact_point"`
	Normal            Vector3r `json:"normal" bson:"normal"`
	Object_id         int64    `json:"object_id" bson:"object_id"`
	Object_name       string   `json:"object_name" bson:"object_name"`
	Penetration_depth float64  `json:"penetration_depth" bson:"penetration_depth"`
	Position          Vector3r `json:"position" bson:"position"`
	Time_stamp        uint64   `json:"time_stamp" bson:"time_stamp"`
}

type Quaternionr struct {
	W_val float64 `json:"w_val" bson:"w_val"`
	X_val float64 `json:"x_val" bson:"x_val"`
	Y_val float64 `json:"y_val" bson:"y_val"`
	Z_val float64 `json:"z_val" bson:"z_val"`
}

func (q *Quaternionr) UnmarshalJSON(data []byte) error {
	//fmt.Println("datadatadata:", string(data))
	var a map[string]float64
	fmt.Println("unmarshalling quaternionr")
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("Custom Quaternionr unmarshal error")
		return err
	}
	q.W_val = a["w_val"]
	q.X_val = a["x_val"]
	q.Y_val = a["y_val"]
	q.Z_val = a["z_val"]
	return nil
}

type GNSS struct {
	Eph       float64  `json:"eph" bson:"eph"`
	Epv       float64  `json:"epv" bson:"epv"`
	Fix_type  int      `json:"fix_type" bson:"fix_type"`
	Geo_point GeoPoint `json:"geo_point" bson:"geo_point"`
	Time_utc  uint64   `json:"time_utc" bson:"time_utc"`
	Velocity  Vector3r `json:"velocity" bson:"velocity"`
}

type GeoPoint struct {
	Altitude  float64 `json:"altitude" bson:"altitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

func (gp *GeoPoint) UnmarshalJSON(data []byte) error {
	var a map[string]float64
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("Custom GeoPoint unmarshal error")
		return err
	}
	gp.Altitude = a["altitude"]
	gp.Latitude = a["latitude"]
	gp.Longitude = a["longitude"]
	return nil
}

type Vector3r struct {
	X_val float64 `json:"x_val" bson:"x_val"`
	Y_val float64 `json:"y_val" bson:"y_val"`
	Z_val float64 `json:"z_val" bson:"z_val"`
}

func (v *Vector3r) UnmarshalJSON(data []byte) error {
	var a map[string]float64
	if err := json.Unmarshal(data, &a); err != nil {
		fmt.Println("data:", string(data))
		fmt.Println("Custom Vector3r unmarshal error")
		return err
	}
	v.X_val = a["x_val"]
	v.Y_val = a["y_val"]
	v.Z_val = a["z_val"]
	return nil
}

// func GetPayloadDefault() *DroneData {
// 	return getPayload(
// 		primitive.NewObjectID(),
// 		"Parrot Anafi",
// 		primitive.NewDateTimeFromTime(time.Now()),
// 		uint16(rand.Intn(101)),
// 		int64(rand.Intn(101)),
// 		int64(rand.Intn(101)),
// 	)
// }

// func getPayload(id primitive.ObjectID, droneID string, published primitive.DateTime,
// 	battery uint16, height int64, temperature int64) *DroneData {

// 	return &DroneData{
// 		id,
// 		droneID,
// 		published,
// 		battery,
// 		height,
// 		temperature,
// 	}
// }

// func GetClient() (*mongo.Client, context.Context) {
// 	clientOptions := options.Client().
// 		ApplyURI("mongodb+srv://thunderpurtz:G2rsb9ae0a64!@cluster0.14i4y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
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
// 	return client, ctx
// }

// func main() {
// 	client, ctx := GetClient()
// 	collection := client.Database("DronePlatform").Collection("droneData")
// 	res, err := collection.InsertOne(ctx, GetPayloadDefault())
// 	if err != nil {
// 		fmt.Println("error in insert")
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(res.InsertedID)
// }
