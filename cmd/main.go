package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Go-Simple-API/internal/models"
	"Go-Simple-API/internal/storage"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

var mh *storage.MongoHandler

func main() {
    mongoDbConnection := "mongodb://localhost:27017"
    mh = storage.NewHandler(mongoDbConnection)
    if mh == nil {
        log.Fatal("Failed to initialize MongoHandler")
    }
    r := registerRoutes()
    
    // Add this line to confirm that the server is starting
    fmt.Println("Server is running on http://localhost:3060")
    
    log.Fatal(http.ListenAndServe(":3060", r))
}



func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", getAllContact)                 // GET /contacts
		r.Get("/{phonenumber}", getContact)       // GET /contacts/0147344454
		r.Post("/", addContact)                   // POST /contacts
		r.Put("/{phonenumber}", updateContact)    // PUT /contacts/0147344454
		r.Delete("/{phonenumber}", deleteContact) // DELETE /contacts/0147344454
	})
	return r
}


func getContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusNotFound)
		return
	}
	contact := &contact.Contact{}
	err := mh.GetOne(contact, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, fmt.Sprintf("Contact with phoneNumber: %s not found", phoneNumber), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(contact)
}

func getAllContact(w http.ResponseWriter, r *http.Request) {
	contacts := mh.Get(bson.M{})
	json.NewEncoder(w).Encode(contacts)
}

func addContact(w http.ResponseWriter, r *http.Request) {
	existingContact := &contact.Contact{}
	var contactData contact.Contact

	if err := json.NewDecoder(r.Body).Decode(&contactData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	contactData.CreatedOn = time.Now()
	err := mh.GetOne(existingContact, bson.M{"phoneNumber": contactData.PhoneNumber})
	if err == nil {
		http.Error(w, fmt.Sprintf("Contact with phoneNumber: %s already exists", contactData.PhoneNumber), http.StatusBadRequest)
		return
	}

	_, err = mh.AddOne(&contactData)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Contact created successfully"))
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	existingContact := &contact.Contact{}
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusNotFound)
		return
	}
	err := mh.GetOne(existingContact, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, fmt.Sprintf("Contact with phoneNumber: %s does not exist", phoneNumber), http.StatusBadRequest)
		return
	}
	_, err = mh.RemoveOne(bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Contact deleted"))
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusNotFound)
		return
	}

	var updatedData contact.Contact
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Prepare the update with the `$set` operator
	update := bson.M{"$set": updatedData}

	// Call mh.Update with only the filter and update arguments
	_, err := mh.Update(bson.M{"phoneNumber": phoneNumber}, update)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Contact update successful"))
}