package contacts

import (
	"addressBook/config"
	"addressBook/models/contact"
	"addressBook/repo/contactrepo"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"encoding/csv"
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type ContactController struct {
	contactrepo *contactrepo.ContactRepo
	appCtx      *config.AppCtx
}

func NewContactController(appCtx *config.AppCtx) *ContactController {
	cRepo := contactrepo.NewContactRepo(appCtx.DB)
	return &ContactController{contactrepo: cRepo, appCtx: appCtx}

}

// Create a contact
func (contactCtrl *ContactController) CreateContact(w http.ResponseWriter, r *http.Request) {

	userEmail, ok := r.Context().Value("userEmail").(string)

	if !ok {
		http.Error(w, "User email not found in context", http.StatusInternalServerError)
		return
	}
	var newContact contact.Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		http.Error(w, "Error in decoding body", http.StatusBadRequest)
		return
	}
	err = contactCtrl.contactrepo.CreateContact(&newContact, userEmail)
	if err != nil {

		http.Error(w, "Error in Creating Contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	// ApiResponseHandler(w, http.StatusCreated, newContact)
}

// Get all contacts
func (contactCtrl *ContactController) GetContacts(w http.ResponseWriter, r *http.Request) {

	userEmail, ok := r.Context().Value("userEmail").(string)

	if !ok {
		http.Error(w, "User email not found in context", http.StatusInternalServerError)
		return
	}
	contactList, err := contactCtrl.contactrepo.GetContacts(userEmail)
	if err != nil {

		http.Error(w, "Error in Getting Contacts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contactList)
}

// Update Contact
func (contactCtrl *ContactController) UpdateContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cId := params["cid"]
	contacId, err := strconv.Atoi(cId)
	if err != nil {
		http.Error(w, "invalid Contact ID", http.StatusInternalServerError)
		return
	}

	var updatedContact contact.Contact
	err = json.NewDecoder(r.Body).Decode(&updatedContact)
	if err != nil {
		http.Error(w, "Error in Decoding Contact", http.StatusInternalServerError)
		return
	}

	err = contactCtrl.contactrepo.UpdateContact(uint(contacId), updatedContact)

	if err != nil {
		http.Error(w, "Error in updating Contact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	// ApiResponseHandler(w, http.StatusOK, updatedContact)
}

// Delete Contact
func (contactCtrl *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cId := params["cid"]
	contacId, err := strconv.Atoi(cId)
	if err != nil {
		http.Error(w, "invalid Contact ID", http.StatusInternalServerError)
		return
	}
	err = contactCtrl.contactrepo.DeleteContact(uint(contacId))
	if err != nil {
		http.Error(w, "Error in updating Contact", http.StatusInternalServerError)

		return
	}
}

// Download CSV File
func (contactCtrl *ContactController) DownloadCSV(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	fromEmail := os.Getenv("FROMEMAIL")
	host := os.Getenv("HOST")
	password := os.Getenv("PASSWORD")
	userEmail, ok := r.Context().Value("userEmail").(string)

	if !ok {
		http.Error(w, "User email not found in context", http.StatusInternalServerError)
		return
	}
	contactList, err := contactCtrl.contactrepo.GetContacts(userEmail)
	if err != nil {
		contactCtrl.appCtx.Logger.Error("Error in getting Contact", zap.Error(err))
	}
	filePath := "./csv-store/" + time.Now().String() + userEmail + ".csv"
	if err := checkDirExist("./csv-store/"); err != nil {
		log.Println("error in check directory exist or not", err)
	}
	// fmt.Println("Enter Fields (Comma - seperated) :")
	// var inputFields string

	// fmt.Scan(&inputFields)

	// selectedFields := strings.Split(inputFields, ",")

	if err := json.NewDecoder(r.Body).Decode(&contact.FieldsRequest); err != nil {
		http.Error(w, "Error decoding JSON Body", http.StatusBadRequest)
		return
	}
	selectedFields := contact.FieldsRequest.Fields

	err = createCSVFile(filePath)
	if err != nil {
		contactCtrl.appCtx.Logger.Error("Error in Creating file", zap.Error(err))
		// TODO: Handler error
	}

	fmt.Println("writecontent func call")
	if err := WriteContactToCSV(filePath, contactList, selectedFields); err != nil {
		contactCtrl.appCtx.Logger.Error("Error in writing file", zap.Error(err))
	}

	// defer writer.Flush()

	fmt.Println("Data Written to contact_list.csv successfully")

	err = SendEmail(fromEmail, userEmail, "CSV Report", "Body", host, password, filePath, 587)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		contactCtrl.appCtx.Logger.Error("Error in Sending Email ", zap.Error(err))
		return
	}

}

func createCSVFile(filePath string) error {
	file, err := os.Create(filePath)

	if err != nil {
		return nil
	}
	defer file.Close()

	return nil
}
func checkDirExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0750); err != nil {
			return err
		}
	}
	return nil
}
func WriteContactToCSV(filePath string, contactList []contact.Contact, headers []string) error {

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	if err := writer.Write(headers); err != nil {
		return err
	}
	for _, contact := range contactList {
		data := getStructFieldValues(contact, headers)
		fmt.Println("Data from write content func: ", data)
		if err := writer.Write(data); err != nil {

			return err
		}
	}
	writer.Flush()
	return nil
}

func getStructFieldValues(con contact.Contact, fields []string) []string {
	objValue := reflect.ValueOf(con)
	data := make([]string, len(fields))

	for i, field := range fields {
		fieldValue := objValue.FieldByName(field)
		if !fieldValue.IsValid() {
			fmt.Printf("Field '%s' is not valid\n", field)
			data[i] = "" // or handle this case as per your requirements
		} else {
			data[i] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	}

	return data
}

func SendEmail(fromEmail, toEmail, subject, body, host, password, filePath string, port int) error {
	m := gomail.NewMessage()

	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	b, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error ", err)
	}
	fmt.Println("Data : ", string(b))
	m.Attach(filePath)

	d := gomail.NewDialer(host, 587, fromEmail, password)

	// Send the email

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email sent successfully.")
	// syscall.Exit(0)
	return nil

}

// "smtp.gmail.com"

func ApiResponseHandler(w http.ResponseWriter, httpStatus int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	res := make(map[string]any)
	switch httpStatus {
	case http.StatusOK, http.StatusAccepted, http.StatusCreated:
		res["data"] = data

	case http.StatusInternalServerError:
		res["message"] = "Internal Server Error"
	}
	json.NewEncoder(w).Encode(res)
}
