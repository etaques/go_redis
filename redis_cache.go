package main

import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "fmt"
)

// Redis databases
const (
	AccountsGlobalCache int = 0
	InvoicesGlobalCache = 1
	PaymentMethodsGlobalCache = 2
	ClientsGlobalCache = 3
	SubscriptionsGlobalCache = 4
	AccLastInvoiceGlobalCache = 5
)

type ExampleObject struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Gender string `json:"gender"`
	Where string `json:"where"`
	Is_married bool `json:"is_married"`
}

var ctx = context.Background()

//redis client reference
var redisClient *redis.Client

func main() {
	ctx := context.TODO()
	
	//Connecting in redis database default
	connectRedis(ctx)
	
	//Examples
	var idubersmith int = 1984
	example_object := ExampleObject {"Arthur Dent", 42, "M", "The hitchhiker's guide to the galaxy", false}
	
	//serializing example object
	serializedAccountObject := serializeObject(example_object);
	
	//setting value serialized on database 0
	setToRedis(ctx, idubersmith, serializedAccountObject, AccountsGlobalCache)
	
	//setting value on database 1
	setToRedis(ctx, idubersmith, serializedAccountObject, InvoicesGlobalCache)
	
	//getting value from database 0
	val := getFromRedis(ctx, idubersmith, AccountsGlobalCache)

	fmt.Printf("First value with key id: %s is %s on database: %s\n", idubersmith, val, AccountsGlobalCache)

}

//ConnectRedis connecting to default database
func connectRedis(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr: "192.168.0.32:6379",
		Password: "",
		DB: 0,
	})

	redisClient = client
}

//setToRedis selecting database before, to set data
func setToRedis(ctx context.Context, key, val string, db int) {
	redisClient.Select(ctx, db)
	err := redisClient.Set(ctx, key, val, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

//getToRedis selecting database before, to get data
func getFromRedis(ctx context.Context, key string, db int) string{
	redisClient.Select(ctx, db)
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return val
}

//getAllKeys ... getting all "keys" in a database
func getAllKeys(ctx context.Context, key string, db int) []string{
	redisClient.Select(ctx, db)
	keys := []string{}

	iter := redisClient.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return keys
}

//unserializeObject (from redis) and return an object
func unserializeObject(serializedObject string) interface{
	obj := interface{}
	//convert to bytes array and unserializing an object
	err := json.Unmarshal([]byte(serializedObject), &obj)
	if err != nil {
		fmt.Println("Can't serialize", serializedObject)
	}
	
	return obj
}

//serializeObject and return a string (to save in redis)
func serializeObject(obj interface{}) string{
	//serializing an object
	serializedAccountObject, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Can't serialize", serializedAccountObject)
	}
	//converting to string serialized object
	return string(serializedAccountObject)
}
