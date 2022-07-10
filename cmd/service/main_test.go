package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/dumpsterfireproject/FoodLogiQ/internal/model"
	"github.com/dumpsterfireproject/FoodLogiQ/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Status = "status"
const Body = "body"

var mongoHost = "localhost"
var mongoPort = 27017
var mongoClient *mongo.Client = nil

const dbName = "testDB"
const collectionName = "test"

var testTime = time.Now()
var IDs = []string{
	"0063c3a5e4232e4cd0274ac2", "0163c3a5e4232e4cd0274ac2", "0263c3a5e4232e4cd0274ac2",
	"1063c3a5e4232e4cd0274ac2", "1163c3a5e4232e4cd0274ac2", "1263c3a5e4232e4cd0274ac2",
}
var userIDs = []string{"12345", "98765"}

func setUpRouter() *gin.Engine {
	authenticationService := service.NewAuthenticationService()
	eventHandler := service.NewEventHandlerService(service.WithClient(mongoClient), service.WithCollectionName(collectionName), service.WithDbName(dbName))
	router := setupHandler(authenticationService, eventHandler)
	return router
}

func getStatus(ctx context.Context) int {
	v := ctx.Value(Status)
	switch status := v.(type) {
	case int:
		return status
	default:
		return 0
	}
}

func theServerIsStarted() {}

// This is kind of brute-force and not that clean, but this is a quick prototype
func theSeedDataHasBeenInserted(ctx context.Context) error {
	collection := mongoClient.Database(dbName).Collection(collectionName)
	oids := []primitive.ObjectID{}
	for _, id := range IDs {
		oid, _ := primitive.ObjectIDFromHex(id)
		oids = append(oids, oid)
	}
	events := []interface{}{
		&model.Event{
			Id:        oids[0],
			CreatedAt: &testTime,
			CreatedBy: &userIDs[0],
			IsDeleted: false,
			Type:      model.ShippingType,
			Contents:  []model.Contents{},
		},
		&model.Event{
			Id:        oids[1],
			CreatedAt: &testTime,
			CreatedBy: &userIDs[0],
			IsDeleted: false,
			Type:      model.ShippingType,
			Contents:  []model.Contents{},
		},
		&model.Event{
			Id:        oids[2],
			CreatedAt: &testTime,
			CreatedBy: &userIDs[0],
			IsDeleted: true,
			Type:      model.ShippingType,
			Contents:  []model.Contents{},
		},
		&model.Event{
			Id:        oids[3],
			CreatedAt: &testTime,
			CreatedBy: &userIDs[1],
			IsDeleted: false,
			Type:      model.ShippingType,
			Contents:  []model.Contents{},
		},
		&model.Event{
			Id:        oids[4],
			CreatedAt: &testTime,
			CreatedBy: &userIDs[1],
			IsDeleted: false,
			Type:      model.ShippingType,
			Contents:  []model.Contents{},
		},
		&model.Event{
			Id:        oids[5],
			CreatedAt: &testTime,
			CreatedBy: &userIDs[1],
			IsDeleted: true,
			Type:      model.ShippingType,
			Contents:  []model.Contents{},
		},
	}
	_, err := collection.InsertMany(ctx, events)
	// ids := result.InsertedIDs
	// fmt.Printf("created %v", ids...)
	// cursor, err := collection.Find(ctx, bson.M{})
	// for {
	// 	if !cursor.Next(ctx) {
	// 		break
	// 	}
	// 	x := cursor.Current
	// 	fmt.Println(x)
	// }
	return err
}

func iPerformGETForIDWithToken(ctx context.Context, id string, token string) (context.Context, error) {
	ctx, err := callAPI(ctx, "GET", fmt.Sprintf("/events/%s", id), token, nil)
	return ctx, err
}

func theResponseShouldHaveEvents(ctx context.Context, count int) error {
	return nil
}

func callAPI(ctx context.Context, method string, endpoint string, token string, body io.Reader) (context.Context, error) {
	r := setUpRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	ctx = context.WithValue(ctx, Body, string(responseData))
	ctx = context.WithValue(ctx, Status, w.Code)
	return ctx, nil
}

func iPerformAGETRequestWithBearerToken(ctx context.Context, token string) (context.Context, error) {
	ctx, err := callAPI(ctx, "GET", "/events", token, nil)
	return ctx, err
}

func iPerformAGETRequestWithNoToken(ctx context.Context) (context.Context, error) {
	r := setUpRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", nil)
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	ctx = context.WithValue(ctx, Body, string(responseData))
	ctx = context.WithValue(ctx, Status, w.Code)
	return ctx, nil
}

func theResponseShouldHaveHttpStatus(ctx context.Context, wantStatus int) error {
	gotStatus := getStatus(ctx)
	if gotStatus != wantStatus {
		return fmt.Errorf("wanted status %d but got %d", wantStatus, gotStatus)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the server is started$`, theServerIsStarted)
	ctx.Step(`^the seed data has been inserted$`, theSeedDataHasBeenInserted)

	ctx.Step(`^the response should have http status (\d+)$`, theResponseShouldHaveHttpStatus)
	ctx.Step(`^I perform a GET request with bearer token (\w+)$`, iPerformAGETRequestWithBearerToken)
	ctx.Step(`^I perform a GET request with no token$`, iPerformAGETRequestWithNoToken)

	ctx.Step(`^I perform a GET for ID (\w+) with token (\w+)$`, iPerformGETForIDWithToken)
	ctx.Step(`^the response should have (\d+) events$`, theResponseShouldHaveEvents)
}

func IntializeTestSuite(sc *godog.TestSuiteContext) {
	sc.BeforeSuite(func() {
		mongoURI := fmt.Sprintf("mongodb://%s:%d", mongoHost, mongoPort)
		fmt.Printf("connecting to mongo at %s\n", mongoURI)
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
		if err != nil {
			panic("could not connect to DB!")
		}
		mongoClient = client
	})
	sc.AfterSuite(func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	})

	sc.ScenarioContext().Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return ctx, nil
	})
	sc.ScenarioContext().After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// clear out any data from the mongo test DB - note at some point, we can use test case specific DB or collection
		// to allow for parallel testing
		collection := mongoClient.Database(dbName).Collection(collectionName)
		_, deleteErr := collection.DeleteMany(ctx, bson.M{})
		return ctx, deleteErr
	})

	sc.ScenarioContext().StepContext().Before(func(ctx context.Context, st *godog.Step) (context.Context, error) {
		return ctx, nil
	})
	sc.ScenarioContext().StepContext().After(func(ctx context.Context, st *godog.Step, status godog.StepResultStatus, err error) (context.Context, error) {
		return ctx, nil
	})
}

func TestStartServer(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: IntializeTestSuite,
		Options: &godog.Options{
			Paths:  []string{"_features"},
			Format: "pretty",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func runTests(m *testing.M) int {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.4",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections"),
	}
	mongo, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("error starting mongo container %s\n", err)
		return -1
	}
	port, err := mongo.MappedPort(ctx, "27017")
	if err != nil {
		fmt.Printf("error getting mongo port %s\n", err)
		return -1
	}
	mongoPort = port.Int()
	host, err := mongo.Host(ctx)
	if err != nil {
		fmt.Printf("error getting mongo host %s\n", err)
		return -1
	}
	mongoHost = host

	defer mongo.Terminate(ctx)
	exitCode := m.Run()
	return exitCode
}

func TestMain(m *testing.M) {
	exitCode := runTests(m)
	os.Exit(exitCode)
}
