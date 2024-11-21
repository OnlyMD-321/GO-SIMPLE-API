package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	contact "github.com/OnlyMD-321/GO-SIMPLE-API/internal/models" // Alias import
	"github.com/OnlyMD-321/GO-SIMPLE-API/internal/storage"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

var mh *storage.MongoHandler


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

	contact := &contact.Contact{} // Use the alias
	err := mh.GetOne(contact, bson.M{"phoneNumber": phoneNumber})
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(contact)
}

// Add a new contact
func AddContact(w http.ResponseWriter, r *http.Request) {
	var newContact contact.Contact // Use the alias
	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	newContact.CreatedOn = time.Now()

	// Check if contact already exists
	existingContact := &contact.Contact{}
	err := mh.GetOne(existingContact, bson.M{"phoneNumber": newContact.PhoneNumber})
	if err == nil {
		http.Error(w, "Contact already exists", http.StatusBadRequest)
		return
	}

	_, err = mh.AddOne(&newContact)
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

	var updatedData contact.Contact // Use the alias
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
