package Tgl

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
)

func GetInvoiceList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		req, err := http.NewRequest("GET", "https://api.mintsoft.co.uk/api/Accounting/Invoice/All", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		invoice_id := r.URL.Query().Get("invoice_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Accounting/Invoice"
		fullURL := fmt.Sprintf("%s/%s/%s", url, invoice_id, "Orders")
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		order_id := r.URL.Query().Get("order_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Order/"
		fullURL := fmt.Sprintf("%s/%s", url, order_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		product_id := r.URL.Query().Get("product_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Product/"
		fullURL := fmt.Sprintf("%s/%s", url, product_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		client_id := r.URL.Query().Get("client_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Product/List?PageNo=1&Limit=100&ClientId="
		fullURL := fmt.Sprintf("%s%s", url, client_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetReturn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		return_id := r.URL.Query().Get("return_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Return/"
		fullURL := fmt.Sprintf("%s/%s", url, return_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetReturns() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		invoice_id := r.URL.Query().Get("invoice_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Accounting/Invoice/"
		fullURL := fmt.Sprintf("%s/%s/%s", url, invoice_id, "Returns")
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetAsns() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		invoice_id := r.URL.Query().Get("invoice_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Accounting/Invoice"
		fullURL := fmt.Sprintf("%s/%s/%s", url, invoice_id, "GoodsIn")
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetAsn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		asn_id := r.URL.Query().Get("asn_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/ASN"
		fullURL := fmt.Sprintf("%s/%s", url, asn_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetExternalInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		invoice_id := r.URL.Query().Get("invoice_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://om.mintsoft.co.uk/Accounts/ShowExternalFulfilmentInvoice/"
		fullURL := fmt.Sprintf("%s/%s", url, invoice_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetStockOutFlow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		product_id := r.URL.Query().Get("product_id")
		fromDate := r.URL.Query().Get("fromDate")
		toDate := r.URL.Query().Get("toDate")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Product"

		// fullURL := fmt.Sprintf("%s/%s/%s%s%s%s%s%s%s", url, product_id, "StockFlow/Filtered?", "FromDate=", fromDate, "&", "ToDate=", toDate, "&WarehouseId=0&&Types=OUT&IncludeOrders=true&IncludeReturns=true&IncludeDetails=false")
		fullURL := fmt.Sprintf("%s/%s/%s%s%s%s%s%s%s", url, product_id, "StockFlow?", "FromDate=", fromDate, "&", "ToDate=", toDate, "&IncludeDetails=true")
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetStockInventory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		product_id := r.URL.Query().Get("product_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Product"

		fullURL := fmt.Sprintf("%s/%s/%s", url, product_id, "Inventory?breakdown=false")

		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetWarehousesData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		// product_id := r.URL.Query().Get("product_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Warehouse"

		// fullURL := fmt.Sprintf("%s/%s/%s", url, product_id, "Inventory?breakdown=false")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetClients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		// product_id := r.URL.Query().Get("product_id")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Client?pageNo=1&limit=100"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		client_id := r.URL.Query().Get("client_id")
		pageno := r.URL.Query().Get("page_no")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Product/List?PageNo="
		fullURL := fmt.Sprintf("%s%s%s%s", url, pageno, "&Limit=100&ClientId=", client_id)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func GetOnhandStockLevel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		product_sku := r.URL.Query().Get("product_sku")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Product/StockLevels?WarehouseId=3&Breakdown=false&ProductId=0&SKU="
		fullURL := fmt.Sprintf("%s%s%s", url, product_sku, "&IncludeSubclients=false")
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func PlaceMintsoftOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Order"
		// fullURL := fmt.Sprintf("%s%s%s", url, product_sku, "&IncludeSubclients=false")
		// fmt.Println(fullURL)

		// bodyBytes, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }
		req, err := http.NewRequest("PUT", url, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func ChangeMintsoftOrderStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ms-apikey")
		order_id := r.URL.Query().Get("order_id")
		status := r.URL.Query().Get("status")

		if apiKey == "" {
			slog.Error("invalid url", "error", "missing api key")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL requested")
			return
		}

		url := "https://api.mintsoft.co.uk/api/Order/"
		fullURL := fmt.Sprintf("%s%s%s%s", url, order_id, "/", status)
		// fmt.Println(fullURL)
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Add("ms-apikey", apiKey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Forward status code, headers, and body
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

// custom APIs

func TglUserSignin(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("TGL Sign in Requested")
		var loginDetails Types.TglSignIn
		readerr := json.NewDecoder(r.Body).Decode(&loginDetails)

		if readerr != nil {
			slog.Error("there was an error in the requested Body", "error", readerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(readerr))
			return
		}

		sessiontoken, _, TglUserAuth, loginerr := storage.TglUserSignin(r.Context(), loginDetails)
		if loginerr != nil {
			slog.Error("error in login ", "error", loginerr)
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(loginerr))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessiontoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		response.WriteJson(w, http.StatusAccepted, TglUserAuth)
	}
}

func TglUploadProduct(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Types.TglProduct
		decerr := json.NewDecoder(r.Body).Decode(&product)

		if decerr != nil {
			slog.Error("Error", "Decode Errror", decerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decerr))
			return
		}

		dberr := storage.TglUploadProduct(r.Context(), product)

		if dberr != nil {
			slog.Error("Error", "Database Errror", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusCreated, "Product Saved Succesfully")

	}
}

func TglGetProducts(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		products, dberr := storage.TglGetProducts(r.Context())

		if dberr != nil {
			slog.Error("Error", "Database Errror", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, products)

	}
}

func TglSaveCustomer(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Customer Types.TglCustomer

		decerr := json.NewDecoder(r.Body).Decode(&Customer)

		if decerr != nil {
			slog.Error("Error", "Decode Error", decerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decerr))
			return
		}

		dberr := storage.TglSaveCustomer(r.Context(), Customer)
		if dberr != nil {
			slog.Error("Error", "Db Error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Customer Saved Succesfully")

	}
}
func TglUpdateCustomer(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Customer Types.TglCustomer

		decerr := json.NewDecoder(r.Body).Decode(&Customer)

		if decerr != nil {
			slog.Error("Error", "Decode Error", decerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decerr))
			return
		}

		dberr := storage.TglUpdateCustomer(r.Context(), Customer)
		if dberr != nil {
			slog.Error("Error", "Db Error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Customer Saved Succesfully")

	}
}

func TglStatusCustomer(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		status := r.URL.Query().Get("status")
		customer_id := r.URL.Query().Get("customer_id")

		if status == "" || customer_id == "" {
			slog.Error("Error", "Decode Error", "invalid parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid parameter")
			return
		}

		boolstatus, err := strconv.ParseBool(status)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(boolstatus)
		if err != nil {
			slog.Error("Error", "Decode Error", "invalid parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid parameter")
			return
		}

		dberr := storage.TglStatusCustomer(r.Context(), boolstatus, customer_id)
		if dberr != nil {
			slog.Error("Error", "Db Error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Customer Status Changedd Succesfully")

	}
}

func TglGetCustomers(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		customers, dberr := storage.TglGetCustomers(r.Context())
		if dberr != nil {
			slog.Error("Error", "Db Error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, customers)

	}
}

func TglPlaceOrder(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var Order []Types.TglOrder

		decerr := json.NewDecoder(r.Body).Decode(&Order)

		if decerr != nil {
			slog.Error("Error", "decode Error", decerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decerr))
			return
		}

		order_id, dberr := storage.TglPlaceOrder(r.Context(), Order)

		if dberr != nil {
			slog.Error("Error", "Db Error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, order_id)

	}
}

func TglAddMintsoftOrderId(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order_id := r.URL.Query().Get("order_id")
		mintsoft_order_id := r.URL.Query().Get("mintsoft_order_id")
		status := r.URL.Query().Get("status")

		if mintsoft_order_id == "" || order_id == "" || status == "" {
			slog.Error("Error", "URL Error", "Invalid URL Request")
			response.WriteJson(w, http.StatusBadRequest, "Invalid Url Request ")
			return
		}

		dberr := storage.TglAddMintsoftOrderId(r.Context(), order_id, mintsoft_order_id, status)

		if dberr != nil {
			slog.Error("Error", "Internal Error", response.GeneralError(dberr))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Mintsoft Id Added succesfully")
	}
}

func TglGetOrders(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		order_id := r.URL.Query().Get("order_id")

		if order_id != "" {

		}

		orders, dberr := storage.TglGetOrders(r.Context())

		if dberr != nil {
			slog.Error("Error", "Db Error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, orders)
	}
}

func TglGetSysClients(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		clientsList, dberr := storage.TglGetSysClients(r.Context())

		if dberr != nil {
			slog.Error("Error", "db Error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, clientsList)
	}
}

func TglAddSysClients(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var TglClient Types.TglClients

		decerr := json.NewDecoder(r.Body).Decode(&TglClient)

		if decerr != nil {
			slog.Error(" Error", "Decode Error", decerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decerr))
			return
		}

		dberr := storage.TglAddSysClients(r.Context(), TglClient)

		if dberr != nil {
			slog.Error(" Error", "Decode Error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusCreated, "Client Added Succesfully")

	}
}

func TglUpdateSysClients(storage Storage.TglApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var TglClient Types.TglClients
		// float64("TglClient.StorageUnitPrice")
		decerr := json.NewDecoder(r.Body).Decode(&TglClient)

		if decerr != nil {
			slog.Error(" Error", "Decode Error", decerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decerr))
			return
		}

		dberr := storage.TglUpdateSysClients(r.Context(), TglClient)

		if dberr != nil {
			slog.Error(" Error", "Decode Error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusCreated, "Client Updated Succesfully")

	}
}
