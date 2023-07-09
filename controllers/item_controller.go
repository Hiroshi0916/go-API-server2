package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Item represents a simple key-value pair
type Item struct {
	ID    string `gorm:"column:id" json:"id"`
	Value string `gorm:"column:value" json:"value"`
}

// Create item
func CreateItemHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// この関数内で createItemHandler のコードを書く
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println("Failed to decode item:", err)
		http.Error(w, "error creating item", http.StatusBadRequest)
		return
	}

	result := db.Create(&item)
	if result.Error != nil {
		log.Println("Failed to store item:", result.Error)
		http.Error(w, "error storing item", http.StatusInternalServerError)
		return
	}
	// _, err = db.Exec("INSERT INTO items (id, value) VALUES ($1, $2)", item.ID, item.Value)
	// if err != nil {
	// 	log.Println("Failed to store item:", err)
	// 	http.Error(w, "error storing item", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Read item
func GetItemHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Item型の変数を宣言
	var item Item
	result := db.First(&item, "id = ?", id)
	if result.Error != nil {
		log.Println("Failed to retrieve item:", result.Error)
		http.Error(w, "error retrieving item", http.StatusInternalServerError)
		return
	}

	// // SELECT文を使用して、データベースからアイテムを取得
	// row := db.QueryRow("SELECT id, value FROM items WHERE id = $1", id)

	// err := row.Scan(&item.ID, &item.Value)
	// if err != nil {
	// 	log.Println("Failed to retrieve item:", err)
	// 	http.Error(w, "error retrieving item", http.StatusInternalServerError)
	// 	return
	// }

	// JSONとしてレスポンスにエンコード
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Update item
func UpdateItemHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println("Failed to decode the request body for item update:", err)
		http.Error(w, "failed to decode the request body for item update", http.StatusBadRequest)
		return
	}

	result := db.Model(&item).Where("id = ?", id).Updates(Item{Value: item.Value})
	if result.Error != nil {
		log.Println("Failed to update item:", result.Error)
		http.Error(w, "error updating item", http.StatusInternalServerError)
		return
	}

	// _, err = db.Exec("UPDATE items SET value = $1 WHERE id = $2", item.Value, id)
	// if err != nil {
	// 	log.Println("Failed to update item:", err)
	// 	http.Error(w, "error updating item", http.StatusInternalServerError)
	// 	return
	// }

	// Setting the Content-Type for the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// item.ID = id
	json.NewEncoder(w).Encode(item)
}

// Delete item
func DeleteItemHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item Item
	result := db.Where("id = ?", id).Delete(&item)
	if result.Error != nil {
		log.Println("Failed to delete item:", result.Error)
		http.Error(w, "error deleting item", http.StatusInternalServerError)
		return
	} else if result.RowsAffected == 0 {
		// レコードが存在しない場合の処理
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item successfully deleted"))
}
