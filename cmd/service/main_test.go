package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/dumpsterfireproject/FoodLogiQ/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Status = "status"
const Body = "body"

var mongoHost = "127.0.0.1"
var mongoPort = 27017
var mongoClient *mongo.Client = nil

func setUpRouter() *gin.Engine {
	authenticationService := service.NewAuthenticationService()
	eventHandler := service.NewEventHandlerService(service.WithClient(mongoClient), service.WithCollectionName("test"), service.WithDbName("testdb"))
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

func theServerIsStarted(ctx context.Context) {}

func iPerformAGETRequestWithBearerToken(ctx context.Context, token string) (context.Context, error) {
	r := setUpRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	ctx = context.WithValue(ctx, Body, string(responseData))
	ctx = context.WithValue(ctx, Status, w.Code)
	return ctx, nil
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
	ctx.Step(`^the response should have http status (\d+)$`, theResponseShouldHaveHttpStatus)
	ctx.Step(`^I perform a GET request with bearer token (\w+)$`, iPerformAGETRequestWithBearerToken)
	ctx.Step(`^I perform a GET request with no token$`, iPerformAGETRequestWithNoToken)
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
		return ctx, nil
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
