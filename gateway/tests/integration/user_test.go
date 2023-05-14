package integration_test

import (
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/api"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
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

var _ = Describe("User", func() {
	var server *fiber.App
	var jwtHeader string
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

	BeforeEach(func() {
		bodyRequest := strings.NewReader(`{"name": "Widdy","email": "widdy@gmail.com","password": "widdy123"}`)
		req := httptest.NewRequest(http.MethodPost, "/register", bodyRequest)
		req.Header.Add("Content-Type", "application/json")
		resp, err := server.Test(req)
		Expect(err).NotTo(HaveOccurred())

		responseBodyUser := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&responseBodyUser)
		Expect(err).NotTo(HaveOccurred())

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

		id := int(responseBodyUser["data"].(map[string]interface{})["id"].(float64))
		target := fmt.Sprintf("/api/users/%d", id)
		req = httptest.NewRequest(http.MethodDelete, target, nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", jwtHeader)
		resp, err = server.Test(req)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err = setup.TearDownTest(configuration)
		if err != nil {
			log.Fatalln("There is something wrong with the tear down test", err)
		}
	})

	Describe("Find all users", func() {
		When("The value of the data users is null", func() {
			It("Should return success with null data", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("The value of the data users isn't null", func() {
			It("Should return a success message upon successfully find all users", func() {
				bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				req = httptest.NewRequest(http.MethodGet, "/api/users", nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].([]interface{})[0].(map[string]interface{})["email"]).To(Equal("widdy@gmail.com"))
			})
		})
	})

	Describe("Find user by id", func() {
		When("The value of the data user is null", func() {
			It("Should throw an error not found", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal(response.GrpcErrorNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("The value of the data user isn't null", func() {
			It("Should return a success message upon successfully find user by id", func() {
				bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/users/%d", id)
				req = httptest.NewRequest(http.MethodGet, target, nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody = map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["email"]).To(Equal("widdy@gmail.com"))
			})
		})
	})

	Describe("Create user", func() {
		Context("The data in the create user request is invalid", func() {
			When("The field name is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'name' field", func() {
					bodyRequest := strings.NewReader(`{"email": "widdy@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
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
					req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
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
					req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
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
					req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
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
		})

		When("The data in the create user request is valid", func() {
			It("Should return a success message upon successfully creating the user", func() {
				bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
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

	Describe("Update user by id", func() {
		When("The value of the data user is null", func() {
			It("Should throw an error not found", func() {
				bodyRequest := strings.NewReader(`{"name": "widdtampan"}`)
				req := httptest.NewRequest(http.MethodPatch, "/api/users/1", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal(response.GrpcErrorNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		Context("The data in the update user request is invalid", func() {
			When("The field name is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'name' field", func() {
					bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
					target := fmt.Sprintf("/api/users/%d", id)
					bodyRequest = strings.NewReader(`{}`)
					req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:UserUpdateRequest.Name"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("The data in the update user request is valid", func() {
			It("Should return a success message upon successfully updating the user", func() {
				bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/users/%d", id)
				bodyRequest = strings.NewReader(`{"name": "widdtampan","password": "tampan123"}`)
				req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody = map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("updated"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("widdtampan"))
			})
		})
	})

	Describe("Delete user by id", func() {
		When("The value of the data user is null", func() {
			It("Should throw an error not found", func() {
				req := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusNotFound))
				Expect(responseBody["status"]).To(Equal(response.GrpcErrorNotFound))
				Expect(responseBody["data"]).To(BeNil())
			})
		})

		When("The value of the data user isn't null", func() {
			It("Should return a success message upon successfully deleted", func() {
				bodyRequest := strings.NewReader(`{"name": "widdy","email": "widdy@gmail.com","password": "widdy123"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/users/%d", id)
				req = httptest.NewRequest(http.MethodDelete, target, nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody = map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("deleted"))
				Expect(responseBody["data"]).To(BeNil())
			})
		})
	})
})
