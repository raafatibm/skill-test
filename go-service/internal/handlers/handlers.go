package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reporting/internal/httpclient"
	"reporting/internal/models"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/phpdave11/gofpdf"
)

type Handler struct {
	cfg models.Config

	httpClient *httpclient.HTTPClient
	csrfToken  string
}

func NewHandler(cfg models.Config) (*Handler, error) {
	var err error

	handler := Handler{
		cfg: cfg,
	}

	handler.httpClient, err = httpclient.NewHttpClient()
	if err != nil {
		return nil, err
	}

	err = handler.backendLogin()
	if err != nil {
		return nil, err
	}

	return &handler, nil
}

func (handler *Handler) GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	value := chi.URLParam(r, "student_id")
	studentID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid student ID")
		return
	}

	studentDetails, err := handler.backendGetStudentDetails(studentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// start creating the PDF file
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.WriteAligned(0, 20, "Student Details Report", "C")

	var lineYStep float64 = 12
	var lineY float64 = 40

	fnWriteDetails := func(label, data string) {
		pdf.SetXY(10, lineY)
		pdf.WriteAligned(50, 10, label+"  ", "R")
		pdf.Cell(60, 10, data)
		lineY += lineYStep
	}

	fnWriteDetails("Sudent ID:", fmt.Sprint(studentDetails.ID))
	fnWriteDetails("First Name:", studentDetails.FirstName)
	fnWriteDetails("Last Name:", studentDetails.LastName)
	fnWriteDetails("Gender:", studentDetails.Gender)
	fnWriteDetails("Birth Date:", studentDetails.BirthDate.Format("2006-01-02"))
	fnWriteDetails("Class:", studentDetails.Class)
	fnWriteDetails("Enrollment Date:", studentDetails.EnrollmentDate.Format("2006-01-02"))
	fnWriteDetails("Status:", studentDetails.Status)

	pdfFileName := fmt.Sprintf("%v_%v_%v_%v.pdf", time.Now().UTC().Unix(), studentDetails.ID, studentDetails.FirstName, studentDetails.LastName)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, pdfFileName))

	err = pdf.Output(w)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, errorText string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errorText)))
}

////////////////////////////////////////////////////////////////////////////
// NodeJs backend API consumptions
////////////////////////////////////////////////////////////////////////////

type studentDetails struct {
	ID             uint64    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Gender         string    `json:"gender"`
	BirthDate      time.Time `json:"birth_date"`
	Class          string    `json:"class"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
}

func (handler *Handler) backendLogin() error {

	statusCode, _, err := handler.httpClient.PostRequest(
		handler.cfg.BackendBaseUrl+"/api/v1/auth/login",
		map[string]interface{}{
			"username": handler.cfg.BackendUserName,
			"password": handler.cfg.BackendUserPassword,
		},
		map[string]string{
			"Content-Type": "application/json",
		},
	)

	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to login to nodeJs backend with status code [%v]", statusCode)
	}

	handler.csrfToken, err = handler.httpClient.GetCookieValue(handler.cfg.BackendBaseUrl, "csrfToken")
	if err != nil {
		return err
	}

	return nil
}

func (handler *Handler) backendRefresh() error {
	statusCode, _, err := handler.httpClient.GetRequest(
		handler.cfg.BackendBaseUrl+"/api/v1/auth/refresh",
		map[string]string{},
	)

	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("failed to refresh authentication from nodeJs backend with status code [%v]", statusCode)
	}

	handler.csrfToken, err = handler.httpClient.GetCookieValue(handler.cfg.BackendBaseUrl, "csrfToken")
	if err != nil {
		return err
	}

	return nil
}

func (handler *Handler) backendGetStudentDetails(studentID uint64) (*studentDetails, error) {

	fnCallAPI := func() (int, []byte, error) {
		return handler.httpClient.GetRequest(
			fmt.Sprintf("%v/api/v1/students/%v", handler.cfg.BackendBaseUrl, studentID),
			map[string]string{
				"x-csrf-token": handler.csrfToken,
			},
		)
	}

	statusCode, body, err := fnCallAPI()
	if err != nil {
		return nil, err
	}

	if statusCode == http.StatusUnauthorized {
		err = handler.backendRefresh()
		if err != nil {
			return nil, err
		}

		statusCode, body, err = fnCallAPI()
		if err != nil {
			return nil, err
		}
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("failed to retrieve student details from nodeJs backend with status code [%v]", statusCode)
	}

	studentDtls := studentDetails{}

	err = json.Unmarshal(body, &studentDtls)
	if err != nil {
		return nil, err
	}

	return &studentDtls, nil
}
