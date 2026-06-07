package TglAppStore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Tokens"
)

type TglAppStore struct {
	*Postgress.Postgress
}

func NewTglAppStore(pg *Postgress.Postgress) *TglAppStore {
	return &TglAppStore{pg}
}

func (p *TglAppStore) TglUserSignin(ctx context.Context, TglUserLogin Types.TglSignIn) (sessionTokens string, csrfTokens string, TglUserAuths *Types.TglUserAuthenticated, err error) {
	var TglUserAuth Types.TglUserAuthenticated
	var hashedPassword string
	var sessionToken string
	var csrfToken string

	tx, txErr := p.DB.BeginTx(ctx, nil)
	if txErr != nil {
		return "", "", TglUserAuths, txErr
	}

	qerr := tx.QueryRowContext(ctx, `SELECT 
		user_email,
		role,
		user_status,
		user_id,
		hashed_password,
		user_name
		FROM tgl_users_credentials 		
		WHERE user_email = $1`, TglUserLogin.Email).Scan(&TglUserAuth.Email, &TglUserAuth.Role, &TglUserAuth.Status, &TglUserAuth.Id, &hashedPassword, &TglUserAuth.Name)

	if qerr != nil {
		tx.Rollback()
		return "", "", &TglUserAuth, qerr
	}

	matched, hasherr := HashPassword.Unhashpassword(TglUserLogin.Password, hashedPassword)

	if hasherr != nil {
		tx.Rollback()
		return "", "", &TglUserAuth, hasherr
	}
	if !matched {
		tx.Rollback()
		return "", "", &TglUserAuth, errors.New("invalid credentials")
	}
	TglUserAuth.Authenticated = true
	session_id, sessionerr := Tokens.GenerateToken(10)
	if sessionerr != nil {
		tx.Rollback()
		return "", "", &TglUserAuth, sessionerr
	}
	session_token, stokenerr := Tokens.GenerateToken(10)
	if stokenerr != nil {
		tx.Rollback()
		return "", "", &TglUserAuth, stokenerr
	}
	sessionToken = session_token
	csrf_token, csrftokenerr := Tokens.GenerateToken(10)
	if csrftokenerr != nil {
		tx.Rollback()
		return "", "", &TglUserAuth, csrftokenerr
	}
	csrfToken = csrf_token
	TglUserAuth.CsrfToken = csrfToken
	_, insertqerr := tx.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, session_id,role,credential_id)  VALUES ($1, $2, $3, $4, $5,$6)", session_token, csrf_token, TglUserAuth.Email, session_id, TglUserAuth.Role, TglUserAuth.Id)
	if insertqerr != nil {
		tx.Rollback()
		return "", "", &TglUserAuth, insertqerr
	}
	tx.Commit()

	fmt.Println("st=", sessionToken, "ct=", csrfToken)
	return sessionToken, csrfToken, &TglUserAuth, nil

}

func (p *TglAppStore) TglUploadProduct(ctx context.Context, product Types.TglProduct) error {

	_, dberr := p.DB.ExecContext(ctx, `
		INSERT INTO tgl_products (
			sku, name, ean_barcode, upc_barcode,
			description, discontinued, back_orderable,
			commodity_code, categories, suppliers, cost_price,
			unit_price, box_price, pallet_price, box_quantity,
			pallet_quantity, image_url, has_expiry_date, country_of_manufacture,product_vat
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,$16,$17,$18,$19,$20
		)
	`,
		product.SKU, product.Name, product.EANBarcode, product.UPCBarcode,
		product.Description, product.Discontinued, product.BackOrderable,
		product.CommodityCode, product.Categories, product.Suppliers,
		product.CostPrice, product.UnitPrice, product.BoxPrice, product.PalletPrice,
		product.BoxQuantity, product.PalletQuantity, product.ImageURL, product.HasExpiryDate,
		product.CountryOfManufacture, product.Vat,
	)

	if dberr != nil {
		return dberr
	}

	return nil
}

func (p *TglAppStore) TglGetProducts(ctx context.Context) ([]Types.TglProduct, error) {
	var products []Types.TglProduct

	rows, dberr := p.DB.QueryContext(ctx, `
		Select  
			sku, name, ean_barcode, upc_barcode,
			description, discontinued, back_orderable,
			commodity_code, categories, suppliers, cost_price,
			unit_price, box_price, pallet_price, box_quantity,
			pallet_quantity, image_url, has_expiry_date, country_of_manufacture,product_id,id,product_vat::integer AS product_vat
		from 
		tgl_products
	`)

	if dberr != nil {
		return products, dberr
	}

	defer rows.Close()

	for rows.Next() {
		var product Types.TglProduct
		scerr := rows.Scan(&product.SKU, &product.Name, &product.EANBarcode, &product.UPCBarcode,
			&product.Description, &product.Discontinued, &product.BackOrderable,
			&product.CommodityCode, &product.Categories, &product.Suppliers,
			&product.CostPrice, &product.UnitPrice, &product.BoxPrice, &product.PalletPrice,
			&product.BoxQuantity, &product.PalletQuantity, &product.ImageURL, &product.HasExpiryDate,
			&product.CountryOfManufacture, &product.ProductID, &product.ID, &product.Vat)
		if scerr != nil {
			return products, scerr
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *TglAppStore) TglSaveCustomer(ctx context.Context, customer Types.TglCustomer) error {
	_, dberr := p.DB.ExecContext(ctx, `INSERT into tgl_sales_customers
	 (
	contact_person,customer_phone,business_name,
	 post_code,address,city,town,county,
	 email,notes,is_active) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		customer.ContactPerson, customer.CustomerPhone, customer.BusinessName,
		customer.PostCode, customer.Address, customer.City, customer.Town, customer.County,
		customer.Email, customer.Notes, customer.IsActive,
	)

	if dberr != nil {
		return dberr
	}

	return nil
}
func (p *TglAppStore) TglUpdateCustomer(ctx context.Context, customer Types.TglCustomer) error {
	_, dberr := p.DB.ExecContext(ctx, `UPDATE tgl_sales_customers
SET
    contact_person = $1,
    customer_phone = $2,
    business_name = $3,
    post_code = $4,
    address = $5,
    city = $6,
    town = $7,
    county = $8,
    email = $9,
    notes = $10,
    is_active = $11
WHERE customer_id = $12`,
		customer.ContactPerson, customer.CustomerPhone, customer.BusinessName,
		customer.PostCode, customer.Address, customer.City, customer.Town, customer.County,
		customer.Email, customer.Notes, customer.IsActive, customer.CustomerID,
	)

	if dberr != nil {
		return dberr
	}

	return nil
}

func (p *TglAppStore) TglStatusCustomer(ctx context.Context, status bool, customer_id string) error {
	_, dberr := p.DB.ExecContext(ctx, `UPDATE tgl_sales_customers set is_active = $1 where customer_id = $2`, status, customer_id)

	if dberr != nil {
		return dberr
	}
	return nil
}

func (p *TglAppStore) TglGetCustomers(ctx context.Context) (customers []Types.TglCustomer, err error) {
	var customersList []Types.TglCustomer

	rows, dberr := p.DB.QueryContext(ctx, `Select 
	id,customer_id,contact_person,customer_phone,business_name,
	 post_code,address,city,town,county,
	 email,notes,is_active from tgl_sales_customers`)

	if dberr != nil {
		return nil, dberr
	}

	defer rows.Close()

	for rows.Next() {
		var customer Types.TglCustomer
		scerr := rows.Scan(&customer.ID, &customer.CustomerID, &customer.ContactPerson, &customer.CustomerPhone, &customer.BusinessName,
			&customer.PostCode, &customer.Address, &customer.City, &customer.Town, &customer.County,
			&customer.Email, &customer.Notes, &customer.IsActive)

		if scerr != nil {
			return nil, scerr
		}

		customersList = append(customersList, customer)
	}

	return customersList, nil

}

func (p *TglAppStore) TglPlaceOrder(ctx context.Context, order []Types.TglOrder) (order_return string, err error) {

	var order_id string

	tx, txerr := p.DB.BeginTx(ctx, nil)

	if txerr != nil {
		return order_id, txerr
	}

	err = tx.QueryRowContext(ctx, `INSERT INTO tgl_orders DEFAULT VALUES RETURNING order_id`).Scan(&order_id)

	if err != nil {
		tx.Rollback()
		return order_id, err
	}

	for _, item := range order {

		_, err = tx.ExecContext(ctx, `INSERT INTO tgl_order_items
			(
				order_id,product_id,quantity,
				unit,customer_id,sold_price,
				percentage_discount,price,free
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			order_id, item.ProductID, item.Quantity,
			item.Unit, item.CustomerID, item.SoldPrice,
			item.PercentageDiscount, item.Price, item.Free,
		)

		if err != nil {
			tx.Rollback()
			return order_id, err
		}
	}

	err = tx.Commit()

	if err != nil {
		return order_id, err
	}

	return order_id, nil
}

func (p TglAppStore) TglAddMintsoftOrderId(ctx context.Context, order_id string, mintsoft_order_id string, status string) error {

	_, dberr := p.DB.ExecContext(ctx, `UPDATE tgl_orders SET mintsoft_id = $1 , status = $2 WHERE order_id = $3`, mintsoft_order_id, status, order_id)

	if dberr != nil {
		return dberr
	}

	return nil
}

func (p TglAppStore) TglGetOrders(ctx context.Context) (ordersData []Types.TglOrderReturn, err error) {
	var OrdersData []Types.TglOrderReturn

	rows, dberr := p.DB.QueryContext(ctx, `SELECT 
    o.order_id,
	o.status,
	o.mintsoft_id,
	oi.customer_id,
	cu.business_name,
	TO_CHAR(oi.created_at, 'YYYY-MM-DD'),
    json_agg(
        json_build_object(
			'product_name',pr.name,
			'product_image',pr.image_url,
            'product_id', oi.product_id,
            'quantity', oi.quantity::int,
            'unit', oi.unit,
            'sold_price', oi.sold_price::int,
            'percentage_discount', oi.percentage_discount,
            'price', oi.price::int,
            'free', oi.free,
			'product_vat',pr.product_vat::int
        )
    ) AS items
	FROM tgl_orders o
	LEFT JOIN tgl_order_items oi 
		ON o.order_id = oi.order_id
	LEFT JOIN tgl_sales_customers cu 
		ON oi.customer_id = cu.customer_id
	LEFT JOIN tgl_products pr
		ON oi.product_id = pr.product_id
	GROUP BY o.order_id, oi.customer_id, cu.business_name,oi.created_at,o.status,o.mintsoft_id`)

	if dberr != nil {
		return OrdersData, dberr
	}

	for rows.Next() {
		var Order Types.TglOrderReturn
		var itemsJSON []byte

		scerr := rows.Scan(
			&Order.OrderId, &Order.Status, &Order.MintSoftId, &Order.CustomerID, &Order.BusinessName, &Order.CreatedAt, &itemsJSON,
		)

		if scerr != nil {
			return OrdersData, scerr
		}

		err := json.Unmarshal(itemsJSON, &Order.Items)
		if err != nil {
			return OrdersData, err
		}

		OrdersData = append(OrdersData, Order)

	}
	return OrdersData, nil
}

func (p *TglAppStore) TglGetSysClients(ctx context.Context) (clients []Types.TglClients, err error) {
	var clientsList []Types.TglClients

	rows, dberr := p.DB.QueryContext(ctx, `SELECT id,name,username,password,storage_unit_price,client_id from tgl_clients`)

	if dberr != nil {
		return clientsList, dberr
	}

	for rows.Next() {
		var client Types.TglClients

		scerr := rows.Scan(&client.ID, &client.Name, &client.Username, &client.Password, &client.StorageUnitPrice, &client.ClientId)
		if scerr != nil {
			return clientsList, scerr
		}

		clientsList = append(clientsList, client)
	}

	return clientsList, nil
}

func (p *TglAppStore) TglAddSysClients(ctx context.Context, client Types.TglClients) error {
	_, dberr := p.DB.ExecContext(ctx,
		`INSERT into tgl_clients 
			(id,name,username,password,storage_unit_price) 
			values ($1,$2,$3,$4,$5)`, client.ID, client.Name, client.Username, client.Password, client.StorageUnitPrice)

	if dberr != nil {
		return dberr
	}
	return nil
}
func (p *TglAppStore) TglUpdateSysClients(ctx context.Context, client Types.TglClients) error {
	_, dberr := p.DB.ExecContext(ctx,
		`UPDATE  tgl_clients 
			SET id = $1, name = $2, username =$3, password =$4, storage_unit_price = $5 
		WHERE client_id = $6`, client.ID, client.Name, client.Username, client.Password, client.StorageUnitPrice, client.ClientId)

	if dberr != nil {
		return dberr
	}
	return nil
}
