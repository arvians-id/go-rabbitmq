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

var _ = Describe("Category", func() {
	var server *fiber.App
	configuration := config.New("../../.env")
	file, err := os.OpenFile("../../logs/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("There is something wrong with the log file", err)
	}
	server, err = api.NewRoutes(configuration, file)
	if err != nil {
		log.Fatalln("There is something wrong with the server", err)
	}

	AfterEach(func() {
		err = setup.TearDownTest(configuration)
		if err != nil {
			log.Fatalln("There is something wrong with the tear down test", err)
		}
	})

	Describe("Find all categories", func() {
		When("The value of the data category is null", func() {
			It("Should return success with null data", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
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

		When("The value of the data categories isn't null", func() {
			It("Should return a success message upon successfully find all categories", func() {
				bodyRequest := strings.NewReader(`{"name": "Belajar Bahasa Pemrograman"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				req = httptest.NewRequest(http.MethodGet, "/api/categories", nil)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].([]interface{})[0].(map[string]interface{})["name"]).To(Equal("Belajar Bahasa Pemrograman"))
			})
		})
	})

	Describe("Find category by id", func() {
		When("The value of the data category is null", func() {
			It("Should throw an error not found", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/categories/1", nil)
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

		When("The value of the data category isn't null", func() {
			It("Should return a success message upon successfully find category by id", func() {
				bodyRequest := strings.NewReader(`{"name": "Belajar Bahasa Pemrograman"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/categories/%d", id)
				req = httptest.NewRequest(http.MethodGet, target, nil)
				resp, err = server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody = map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("OK"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Belajar Bahasa Pemrograman"))
			})
		})
	})

	Describe("Create category", func() {
		Context("The data in the create category request is invalid", func() {
			When("The field name is blank", func() {
				It("Should throw a validation error in request with message 'bad request' on the 'name' field", func() {
					bodyRequest := strings.NewReader(`{}`)
					req := httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
					req.Header.Add("Content-Type", "application/json")
					resp, err := server.Test(req)
					Expect(err).NotTo(HaveOccurred())

					responseBody := map[string]interface{}{}
					err = json.NewDecoder(resp.Body).Decode(&responseBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusBadRequest))
					Expect(responseBody["status"]).To(Equal("validation error on field:CategoryCreateRequest.Name"))
					Expect(responseBody["data"]).To(BeNil())
				})
			})
		})

		When("The data in the create category request is valid", func() {
			It("Should return a success message upon successfully creating the category", func() {
				bodyRequest := strings.NewReader(`{"name": "Belajar Bahasa Pemrograman"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusCreated))
				Expect(responseBody["status"]).To(Equal("created"))
				Expect(responseBody["data"].(map[string]interface{})["name"]).To(Equal("Belajar Bahasa Pemrograman"))
			})
		})
	})

	Describe("Delete category by id", func() {
		When("The value of the data category is null", func() {
			It("Should throw an error not found", func() {
				req := httptest.NewRequest(http.MethodDelete, "/api/categories/1", nil)
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

		When("The value of the data category isn't null", func() {
			It("Should return a success message upon successfully deleted", func() {
				bodyRequest := strings.NewReader(`{"name": "Belajar Bahasa Pemrograman"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/categories", bodyRequest)
				req.Header.Add("Content-Type", "application/json")
				resp, err := server.Test(req)
				Expect(err).NotTo(HaveOccurred())

				responseBody := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

				id := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				target := fmt.Sprintf("/api/categories/%d", id)
				req = httptest.NewRequest(http.MethodDelete, target, nil)
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
