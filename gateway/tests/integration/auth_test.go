package integration_test

import (
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

var _ = Describe("Auth", func() {
	var server *fiber.App
	configuration := config.New("../../.env")
	file, err := os.OpenFile("../../logs/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("There is something wrong with the log file", err)
	}

	// Init Rabbit MQ
	_, ch, err := config.InitRabbitMQ(configuration)
	if err != nil {
		log.Fatalln("There is something wrong with the rabbit mq", err)
	}

	// Init Server
	server, err = api.NewRoutes(configuration, file, ch)
	if err != nil {
		log.Fatalln("There is something wrong with the server", err)
	}

	AfterEach(func() {
		err = setup.TearDownTest(configuration)
		if err != nil {
			log.Fatalln("There is something wrong with the tear down test", err)
		}
	})

	Describe("Register", func() {
		Context("The data in the register user request is invalid", func() {
			When("The field name is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'name' field", func() {
					bodyRequest := strings.NewReader(`{"email": "widdy@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:UserCreateRequest.Name"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field email is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'email' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy Tampan","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:UserCreateRequest.Email"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field password is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'password' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy Tampan","email": "widdy@gmail.com"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:UserCreateRequest.Password"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field email is not valid email", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'email' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy Tampan","email": "widdy","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:UserCreateRequest.Email"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field password's length is not enough", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'password' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy Tampan","email": "widdy@gmail.com","password": "wid"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:UserCreateRequest.Password"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("The data in the register user request is valid", func() {
			It("Should return a success message upon successfully creating the user", func() {
				bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusCreated))
				Expect(responseBody["status"]).To(Equal("created"))
				Expect(responseBody["data"].(map[string]interface{})["email"]).To(Equal("widdy@gmail.com"))
			})
		})
	})

	Describe("Login", func() {
		Context("The data in the login user request is invalid", func() {
			When("The field email is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'email' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy", email": "widdy@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					bodyRequest = strings.NewReader(`{"password": "widdy123"}`)
					req = httptest.NewRequest(http.MethodPost, "/login", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:LoginRequest.Email"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field password is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'password' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy", email": "widdy@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					bodyRequest = strings.NewReader(`{"email": "widdy@gmail.com"}`)
					req = httptest.NewRequest(http.MethodPost, "/login", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:LoginRequest.Password"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field email is not valid email", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'email' field", func() {
					bodyRequest := strings.NewReader(`{"name": "Widdy", email": "widdy@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					bodyRequest = strings.NewReader(`{"email": "widdy","password": "widdy123"}`)
					req = httptest.NewRequest(http.MethodPost, "/login", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:LoginRequest.Email"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("The data in the login user request is valid", func() {
			It("Should return a success message upon successfully login", func() {
				bodyRequest := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				bodyRequest = strings.NewReader(`{"email": "widdy@gmail.com","password": "widdy123"}`)
				req = httptest.NewRequest(http.MethodPost, "/login", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody = map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["access_token"]).ToNot(BeNil())
			})
		})
	})
})
