package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/OnlyMD-321/GO-SIMPLE-API/internal/models"
	"github.com/OnlyMD-321/GO-SIMPLE-API/internal/storage"
	"github.com/OnlyMD-321/GO-SIMPLE-API/internal/handlers"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"

)



var mh *storage.MongoHandler


func main() {
	// Load MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Ping the database to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
}
func registerRoutes() http.Handler {
	r := chi.NewRouter()

	// Authentication routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.Register)  // POST /auth/register
		r.Post("/login", handlers.Login)        // POST /auth/login
		r.Post("/logout", handlers.Logout)      // POST /auth/logout (token invalidation not implemented)
	})

	// Contact routes (protected)
	r.Route("/contacts", func(r chi.Router) {
		r.Use(handlers.ValidateTokenMiddleware) // Protect all routes in this block with the middleware
		r.Get("/", getAllContact)               // GET /contacts
		r.Get("/{phonenumber}", getContact)     // GET /contacts/0147344454
		r.Post("/", addContact)                 // POST /contacts
		r.Put("/{phonenumber}", updateContact)  // PUT /contacts/0147344454
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
