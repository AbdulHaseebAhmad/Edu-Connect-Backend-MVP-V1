package Tgl

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

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
