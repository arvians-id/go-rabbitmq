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

var _ = Describe("Todo", func() {
	var server *fiber.App
	var user map[string]interface{}
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

		user = responseBodyUser["data"].(map[string]interface{})
		category = responseBodyCategory["data"].(map[string]interface{})
	})

	AfterEach(func() {
		err = setup.TearDownTest(configuration)
		if err != nil {
			log.Fatalln("There is something wrong with the tear down test", err)
		}
	})

	Describe("Find all todos", func() {
		When("The value of the data todos is null", func() {
			It("Should return success with null data", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
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

		When("The value of the data todos isn't null", func() {
			It("Should return a success message upon successfully find all todos", func() {
				idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
				idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
				bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
				req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				req = httptest.NewRequest(http.MethodGet, "/api/todos", nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].([]interface{})[0].(map[string]interface{})["title"]).To(Equal("This is first todo"))
				Expect(responseBody["data"].([]interface{})[0].(map[string]interface{})["description"]).To(Equal("Lorem ipsum"))
				Expect(int(responseBody["data"].([]interface{})[0].(map[string]interface{})["user_id"].(float64))).To(Equal(int(user["id"].(float64))))
			})
		})
	})

	Describe("Find todo by id", func() {
		When("The value of the data todo is null", func() {
			It("Should throw an error not found", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/todos/1", nil)
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

		When("The value of the data todo isn't null", func() {
			It("Should return a success message upon successfully find todo by id", func() {
				idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
				idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
				bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
				req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/todos/%d", id)
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
				Expect(responseBody["data"].(map[string]interface{})["title"]).To(Equal("This is first todo"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Lorem ipsum"))
				Expect(int(responseBody["data"].(map[string]interface{})["user_id"].(float64))).To(Equal(int(user["id"].(float64))))
			})
		})
	})

	Describe("Create todo", func() {
		Context("The data in the create todo request is invalid", func() {
			When("The field title is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'title' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoCreateRequest.Title"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field description is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'description' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoCreateRequest.Description"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field categories is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'categories' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoCreateRequest.Categories"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field user_id is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'user_id' field", func() {
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoCreateRequest.UserId"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("The data in the create todo request is valid", func() {
			It("Should return a success message upon successfully creating the todo", func() {
				idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
				idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
				bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
				req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusCreated))
				Expect(responseBody["status"]).To(Equal("created"))
				Expect(responseBody["data"].(map[string]interface{})["title"]).To(Equal("This is first todo"))
				Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Lorem ipsum"))
				Expect(int(responseBody["data"].(map[string]interface{})["user_id"].(float64))).To(Equal(int(user["id"].(float64))))
			})
		})
	})

	Describe("Update todo by id", func() {
		When("The value of the data todo is null", func() {
			It("Should throw an error not found", func() {
				idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
				idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
				bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
				req := httptest.NewRequest(http.MethodPatch, "/api/todos/1", bodyRequest)
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

		Context("The data in the update todo request is invalid", func() {
			When("The field title is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'title' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
					target := fmt.Sprintf("/api/todos/%d", id)
					bodyRequest = strings.NewReader(`{"description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoUpdateRequest.Title"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field categories is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'categories' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
					target := fmt.Sprintf("/api/todos/%d", id)
					bodyRequest = strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `}`)
					req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoUpdateRequest.Categories"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field description is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'description' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
					target := fmt.Sprintf("/api/todos/%d", id)
					bodyRequest = strings.NewReader(`{"title": "This is first todo","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoUpdateRequest.Description"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})

			When("The field user_id is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'user_id' field", func() {
					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
					target := fmt.Sprintf("/api/todos/%d", id)
					bodyRequest = strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","categories": ` + idCategory + `}`)
					req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:TodoUpdateRequest.UserId"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		Context("The data in the update todo request is valid", func() {
			When("All fields have been filled in", func() {
				It("Should return a success message upon successfully updating the todo", func() {
					bodyRequest := strings.NewReader(`{"name": "sitampan","email": "sitampan@gmail.com","password": "widdy123"}`)
					req := httptest.NewRequest(http.MethodPost, "/api/users", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBodyUser := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBodyUser)
					Expect(err).NotTo(HaveOccurred())
					idUserNew := fmt.Sprintf("%d", int(responseBodyUser["data"].(map[string]interface{})["id"].(float64)))

					idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
					idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
					bodyRequest = strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
					req = httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					bodyRequest = strings.NewReader(`{"name": "Belajar Backend"}`)
					req = httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBodyCategory := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBodyCategory)
					Expect(err).NotTo(HaveOccurred())

					id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
					idCategory = fmt.Sprintf("[%d]", int(responseBodyCategory["data"].(map[string]interface{})["id"].(float64)))
					target := fmt.Sprintf("/api/todos/%d", id)
					bodyRequest = strings.NewReader(`{"title": "This is first todo edited","description": "Lorem ipsum edited","user_id": ` + idUserNew + `,"is_done": true,"categories": ` + idCategory + `}`)
					req = httptest.NewRequest(http.MethodPatch, target, bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", jwtHeader)
					resp, err = server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody = map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					target = fmt.Sprintf("/api/todos/%d", id)
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
					Expect(responseBody["data"].(map[string]interface{})["title"]).To(Equal("This is first todo edited"))
					Expect(responseBody["data"].(map[string]interface{})["is_done"]).To(Equal(true))
					Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Lorem ipsum edited"))
					Expect(int(responseBody["data"].(map[string]interface{})["user_id"].(float64))).To(Equal(int(responseBodyUser["data"].(map[string]interface{})["id"].(float64))))
				})
			})
		})
	})

	Describe("Delete todo by id", func() {
		When("The value of the data todo is null", func() {
			It("Should throw an error not found", func() {
				req := httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil)
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

		When("The value of the data todo isn't null", func() {
			It("Should return a success message upon successfully deleted", func() {
				idUser := fmt.Sprintf("%d", int(user["id"].(float64)))
				idCategory := fmt.Sprintf("[%d]", int(category["id"].(float64)))
				bodyRequest := strings.NewReader(`{"title": "This is first todo","description": "Lorem ipsum","user_id": ` + idUser + `,"categories": ` + idCategory + `}`)
				req := httptest.NewRequest(http.MethodPost, "/api/todos", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", jwtHeader)
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/todos/%d", id)
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
