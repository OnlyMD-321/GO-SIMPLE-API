package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"Go-Simple-API/internal/models"
	"github.com/yourusername/Go-Simple-API/internal/storage"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

var mh *storage.MongoHandler

// InitializeMongoHandler initializes the MongoDB handler
func InitializeMongoHandler(connectionString string) {
	mh = storage.NewMongoHandler(connectionString)
}

// Get all contacts
func GetAllContacts(w http.ResponseWriter, r *http.Request) {
	contacts := mh.Get(bson.M{})
	json.NewEncoder(w).Encode(contacts)
}

// Get a specific contact by phone number
func GetContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusNotFound)
		return
	}

	contact := &models.Contact{}
	err := mh.GetOne(contact, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(contact)
}

// Add a new contact
func AddContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	contact.CreatedOn = time.Now()

	// Check if contact already exists
	existingContact := &models.Contact{}
	err := mh.GetOne(existingContact, bson.M{"phoneNumber": contact.PhoneNumber})
	if err == nil {
		http.Error(w, "Contact already exists", http.StatusBadRequest)
		return
	}

	_, err = mh.AddOne(&contact)
	if err != nil {
		http.Error(w, "Failed to add contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Contact added successfully"))
}

// Update an existing contact
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusNotFound)
		return
	}

	var updatedData models.Contact
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": updatedData}
	_, err := mh.Update(bson.M{"phoneNumber": phoneNumber}, update)
	if err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Contact updated successfully"))
}

// Delete a contact
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	phoneNumber := chi.URLParam(r, "phonenumber")
	if phoneNumber == "" {
		http.Error(w, "Phone number not provided", http.StatusNotFound)
		return
	}

	_, err := mh.RemoveOne(bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Contact deleted successfully"))
}
