package integration_test

import (
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/api"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/tests/setup"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var _ = Describe("CategoryTodo", func() {
	var server *fiber.App
	var todo map[string]interface{}
	var category map[string]interface{}
	var jwtHeader string
	configuration := config.New("../../.env")
	file, err := os.OpenFile("../../logs/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("There is something wrong with the log file", err)
	}

	// Init Rabbit MQ
	_, ch, err := config.InitRabbitMQ(configuration)
	if err != nil {
		log.Fatalln(err)
	}

	// Init Server
	server, err = api.NewRoutes(configuration, file, ch)
	if err != nil {
		log.Fatalln("There is something wrong with the server", err)
	}

	BeforeEach(func() {
		bodyRequest := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "widdy123"}`)
		req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
		req.Header.Add("Content-Type", "application/json")
		resp, err := server.Test(req)
		Expect(err).NotTo(HaveOccurred())

		responseBodyUser := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&responseBodyUser)
		Expect(err).NotTo(HaveOccurred())
		user := responseBodyUser["data"].(map[string]interface{})

		bodyRequest = strings.NewReader(`{"email": "widdy@gmail.com","password": "widdy123"}`)
		req = httptest.NewRequest(http.MethodPost, "/login", bodyRequest)
		req.Header.Add("Content-Type", "application/json")
		resp, err = server.Test(req)
		Expect(err).NotTo(HaveOccurred())

		responseBody := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		Expect(err).NotTo(HaveOccurred())

		Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
		Expect(responseBody["status"]).To(Equal("OK"))
		Expect(responseBody["data"].(map[string]interface{})["access_token"]).ToNot(BeNil())

		jwtHeader = fmt.Sprintf("Bearer %s", responseBody["data"].(map[string]interface{})["access_token"].(string))

		bodyRequest = strings.NewReader(`{"name": "Belajar Bahasa Pemrograman"}`)
		req = httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", jwtHeader)
		resp, err = server.Test(req)
		Expect(err).NotTo(HaveOccurred())

		responseBodyCategory := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&responseBodyCategory)
		Expect(err).NotTo(HaveOccurred())
		category = responseBodyCategory["data"].(map[string]interface{})

		idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
		idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
		bodyRequest = strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
		req = httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", jwtHeader)
		resp, err = server.Test(req)
		Expect(err).NotTo(HaveOccurred())

		responseBody = map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		Expect(err).NotTo(HaveOccurred())

		todo = responseBody["data"].(map[string]interface{})
	})

	AfterEach(func() {
		err = setup.TearDownTest(configuration)
		if err != nil {
			log.Fatalln("There is something wrong with the tear down test", err)
		}
	})

	Describe("Delete todo by id", func() {
		When("The value of the data todo isn't null", func() {
			It("Should return a success message upon successfully deleted", func() {
				categoryId := int(category["id"].(float64))
				todoId := int(todo["id"].(float64))
				request := fmt.Sprintf(`{"category_id": %d, "todo_id": %d}`, categoryId, todoId)
				bodyRequest := strings.NewReader(request)
				req := httptest.NewRequest(http.MethodDelete, "/api/category-todo", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("deleted"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
