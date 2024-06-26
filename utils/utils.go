package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/razdacoder/mcwale-api/models"
)

var Validate = validator.New()

func ParseJSON(request *http.Request, payload any) error {
	if request.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(request.Body).Decode(payload)
}

func WriteJSON(writer http.ResponseWriter, status int, value any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)

	return json.NewEncoder(writer).Encode(value)
}

func WriteError(writer http.ResponseWriter, status int, err error) {
	WriteJSON(writer, status, map[string]string{"error": err.Error()})
}

func ParseStringToInt(value string, fallback int) int {
	i, err := strconv.ParseInt(value, 10, 16)

	if err != nil {
		return fallback
	}

	return int(i)
}

func ParseProductUpdate(product *models.Product) map[string]interface{} {
	productMap := make(map[string]interface{})
	productMap["is_featured"] = product.IsFeatured
	productMap["discount_percentage"] = product.DiscountPercentage
	if product.Title != "" {
		productMap["title"] = product.Title
	}
	if product.Description != "" {
		productMap["description"] = product.Description
	}
	if product.Price > 0 {
		productMap["price"] = product.Price
	}
	if product.Style != "" {
		productMap["style"] = product.Style
	}
	if product.Images != nil {
		productMap["images"] = product.Images
	}

	return productMap
}
