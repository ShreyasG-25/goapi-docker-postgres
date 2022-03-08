package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

type Supplier struct {
	Supplier_id   int `json:"id"`
	Supplier_name string `json:"name"`
	Contact       int    `json:"contact"`
	City          string `json:"city"`
	Address       string `json:"address"`
	Postcode      int    `json:"postcode"`
}

type Category struct {
	Category_id	int `json:"cid"`
	Category_name	string `json:"cname"`
	Description	string `json:"description"`
}

type Product struct {
	Product_id	int `json:"pid"`
	Product_name	string `json:"pname"`
	Supplier_id	int `json:"sid"`
	Category_id	int `json:"cid"`
	Price	int `json:"price"`
	Expirydate string `json:"expdate"`
}

type ViewProdCat struct{
	Product_id	int `json:"pid"`
	Product_name	string `json:"pname"`
	Supplier_name string `json:"sname"`
	Category_name string `json:"cname"`
	Price	int `json:"price"`
}

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []Supplier `json:"data"`
    Message string `json:"message"`
}

// DB set up  
func setupDB() *sql.DB {
	var user=os.Getenv("DB_USERNAME")
	var password=os.Getenv("DB_PASSWORD")
	var dbname=os.Getenv("DB_DB")
	var host=os.Getenv("DB_HOST")
	//var port=os.Getenv("DB_PORT")
	var port=5432
    dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host,port,user,password,dbname)
    db, err := sql.Open("postgres", dbinfo)
    checkErr(err)
    return db
}
// Function for handling errors
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

// Function for handling messages
func printMessage(message string) {
    fmt.Println("")
    fmt.Println(message)
    fmt.Println("")
}

// Init Suppliers var as a slice Supplier struct
var Suppliers []Supplier
var Categories []Category
var Products []Product

// Get all suppliers
func GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting suppliers...")
	// Get all movies from movies table that don't have movieID = "1"
    rows, err := db.Query(`SELECT * FROM "Suppliers"`)
    // check errors
    checkErr(err)
    // var response []JsonResponse
    var Suppliers []Supplier
    // Foreach supplier
    for rows.Next() {
        var s Supplier
        err = rows.Scan(&s.Supplier_id, &s.Supplier_name, &s.Contact,&s.City,&s.Address,&s.Postcode)
        // check errors
        checkErr(err)
        Suppliers= append(Suppliers, Supplier{Supplier_id: s.Supplier_id, Supplier_name: s.Supplier_name,Contact:s.Contact ,City: s.City,Address: s.Address,Postcode: s.Postcode})
    }
    var response = JsonResponse{Type: "success", Data: Suppliers}
    json.NewEncoder(w).Encode(response)
}

// Get single supplier
func GetSupplier(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting supplier...")
	vars := mux.Vars(r)
	idval, _ := strconv.Atoi(vars["id"])
	// Get all supplier having id=2
    rows, err := db.Query(`SELECT * FROM "Suppliers" WHERE supplier_id=$1`,idval)
    // check errors
    checkErr(err)
    // var response []JsonResponse
    var Suppliers []Supplier

    // Foreach supplier
    for rows.Next() {
        var s Supplier
        err = rows.Scan(&s.Supplier_id, &s.Supplier_name, &s.Contact,&s.City,&s.Address,&s.Postcode)
        // check errors
        checkErr(err)
        Suppliers= append(Suppliers, Supplier{Supplier_id: s.Supplier_id, Supplier_name: s.Supplier_name,Contact:s.Contact ,City: s.City,Address: s.Address,Postcode: s.Postcode})
    }
    var response = JsonResponse{Type: "success", Data: Suppliers}
    json.NewEncoder(w).Encode(response)
}

// Add new Supplier
func createSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var s Supplier
	_ = json.NewDecoder(r.Body).Decode(&s)
	db := setupDB()
	printMessage("Creating supplier...")
	res,e := db.Exec(`INSERT INTO "Suppliers"(supplier_name, contact, city, address, postcode) VALUES($1, $2, $3, $4,$5) RETURNING supplier_id`, s.Supplier_name, s.Contact, s.City, s.Address, s.Postcode)
	fmt.Println(res.LastInsertId())
	//fmt.Println(res)
	checkErr(e)
	Suppliers = append(Suppliers, s)
	json.NewEncoder(w).Encode(&s)
}

// Update supplier
func updateSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var s Supplier
	_ = json.NewDecoder(r.Body).Decode(&s)
	vars := mux.Vars(r)
	idval, _ := strconv.Atoi(vars["id"])
	db := setupDB()
	printMessage("Updating supplier...")
	res,e := db.Exec(`UPDATE "Suppliers" SET supplier_name=$1, contact=$2, city=$3, address=$4, postcode=$5 WHERE supplier_id=$6`, s.Supplier_name, s.Contact, s.City, s.Address, s.Postcode,idval)
	checkErr(e)
	var v,_=res.RowsAffected()
	if v>0{
		Suppliers = append(Suppliers, s)
		json.NewEncoder(w).Encode(&s)
	}
	//fmt.Println(res)
	
	
}

// Delete supplier
func deleteSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var idval,_=strconv.Atoi(vars["id"])
	db := setupDB()

	rows, err := db.Query(`SELECT * FROM "Suppliers" WHERE supplier_id=$1`,idval)

    // check errors
    checkErr(err)
    // var response []JsonResponse
    var Suppliers []Supplier

    // Foreach supplier
    for rows.Next() {
        var s Supplier
        err = rows.Scan(&s.Supplier_id, &s.Supplier_name, &s.Contact,&s.City,&s.Address,&s.Postcode)
        // check errors
        checkErr(err)

        Suppliers= append(Suppliers, Supplier{Supplier_id: s.Supplier_id, Supplier_name: s.Supplier_name,Contact:s.Contact ,City: s.City,Address: s.Address,Postcode: s.Postcode})
    }
   
	_,e := db.Query(`DELETE FROM "Suppliers" WHERE supplier_id=$1`, idval)
	checkErr(e)
	var response = JsonResponse{Type: "success", Data: Suppliers,Message:"Removed supplier id="+strconv.Itoa(idval)}
	json.NewEncoder(w).Encode(response)
}

//create category
func createCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c Category
	_ = json.NewDecoder(r.Body).Decode(&c)
	db := setupDB()
	printMessage("Creating category...")
	res,e := db.Exec(`INSERT INTO "Categories"(category_name, description) VALUES($1, $2) RETURNING category_id`, c.Category_name, c.Description)
	fmt.Println(res.LastInsertId())
	checkErr(e)
	Categories = append(Categories, c)
	json.NewEncoder(w).Encode(&c)
}

//view categories
func viewCategories(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting categories...")
	// Get all movies from movies table that don't have movieID = "1"
    rows, err := db.Query(`SELECT * FROM "Categories"`)
    // check errors
    checkErr(err)
    // var response []JsonResponse
    var Categories []Category

    // Foreach supplier
    for rows.Next() {
        var s Category

        err = rows.Scan(&s.Category_id,&s.Category_name, &s.Description)

        // check errors
        checkErr(err)

        Categories= append(Categories, Category{Category_id: s.Category_id, Category_name: s.Category_name,Description:s.Description})
    }
    //var response = JsonResponse{Type: "success", Data: Categories}
    json.NewEncoder(w).Encode(Categories)
}

//delete category
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var idval,_=strconv.Atoi(vars["id"])
	db := setupDB()
	res,e := db.Exec(`DELETE FROM "Categories" WHERE category_id=$1`, idval)
	checkErr(e)
	var rows,_=res.RowsAffected() 
	var response JsonResponse
	if rows>0{
		response = JsonResponse{Type: "success",Message:"Removed category id="+strconv.Itoa(idval)}
	}else{
		response = JsonResponse{Type: "failure",Message:"Cant find category id="+strconv.Itoa(idval)}
	}
	json.NewEncoder(w).Encode(response)
}

//update category
func updateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var s Category
	_ = json.NewDecoder(r.Body).Decode(&s)
	vars := mux.Vars(r)
	idval, _ := strconv.Atoi(vars["id"])
	db := setupDB()
	printMessage("Updating category...")
	res,e := db.Exec(`UPDATE "Categories" SET category_name=$1, description=$2 WHERE category_id=$3`, s.Category_name, s.Description,idval)
	checkErr(e)
	var v,_=res.RowsAffected()
	if v>0{
		s.Category_id=idval
		Categories = append(Categories, s)
		json.NewEncoder(w).Encode(&s)
	}	
}

//add product
func addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c Product
	_ = json.NewDecoder(r.Body).Decode(&c)
	db := setupDB()
	printMessage("Adding product...")
	_,e := db.Exec(`INSERT INTO "Product"(product_name, supplier_id, category_id, price, expirydate) VALUES($1, $2, $3, $4, $5) RETURNING product_id`, c.Product_name, c.Supplier_id, c.Category_id, c.Price, c.Expirydate)
	//fmt.Println(res.LastInsertId())
	checkErr(e)
	Products = append(Products, c)
	json.NewEncoder(w).Encode(&c)
}

//view products
func viewProduct(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting products...")
    rows, err := db.Query(`SELECT * FROM "Product" order by category_id`)
    checkErr(err)
    var Products []Product
    // Foreach supplier
    for rows.Next() {
        var s Product
        err = rows.Scan(&s.Product_id,&s.Product_name, &s.Supplier_id, &s.Category_id, &s.Price, &s.Expirydate)
        // check errors
        checkErr(err)
        Products= append(Products, Product{Product_id: s.Product_id, Product_name: s.Product_name,Supplier_id:s.Supplier_id,Category_id:s.Category_id,Price:s.Price,Expirydate:s.Expirydate})
    }
    //var response = JsonResponse{Type: "success", Data: Categories}
    json.NewEncoder(w).Encode(Products)
}

//view products
func viewProductDetails(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting products Details...")
    rows, err := db.Query(`SELECT "Product".product_id,"Product".product_name,"Suppliers".supplier_name,"Categories".category_name,"Product".price
	FROM (("Product" INNER JOIN "Suppliers" ON "Product".supplier_id="Suppliers".supplier_id) INNER JOIN "Categories" ON "Product".category_id="Categories".category_id)`)
    checkErr(err)
    var ProdCats []ViewProdCat
    // Foreach supplier
    for rows.Next() {
        var s ViewProdCat
        err = rows.Scan(&s.Product_id,&s.Product_name, &s.Supplier_name, &s.Category_name, &s.Price)
        checkErr(err)
        ProdCats= append(ProdCats, ViewProdCat{Product_id: s.Product_id, Product_name: s.Product_name,Supplier_name:s.Supplier_name,Category_name:s.Category_name,Price:s.Price})
    }
    json.NewEncoder(w).Encode(ProdCats)
}


//view products
func viewProd(w http.ResponseWriter, r *http.Request) {
	//db := setupDB()
	printMessage("Getting products Details...")
	var pagecount,_=strconv.Atoi(r.URL.Query().Get("page"))
	if pagecount==0{
		pagecount=1
	}
	perpage:=9
	db := setupDB()
	rows, err := db.Query(`SELECT "Product".product_id,"Product".product_name,"Suppliers".supplier_name,"Categories".category_name,"Product".price
	FROM (("Product" INNER JOIN "Suppliers" ON "Product".supplier_id="Suppliers".supplier_id) INNER JOIN "Categories" ON "Product".category_id="Categories".category_id) ORDER BY "Product".product_id LIMIT $1 OFFSET $2`,perpage,(pagecount-1)*perpage)
    checkErr(err)
    var ProdCats []ViewProdCat
    // Foreach supplier
    for rows.Next() {
        var s ViewProdCat
        err = rows.Scan(&s.Product_id,&s.Product_name, &s.Supplier_name, &s.Category_name, &s.Price)
        checkErr(err)
        ProdCats= append(ProdCats, ViewProdCat{Product_id: s.Product_id, Product_name: s.Product_name,Supplier_name:s.Supplier_name,Category_name:s.Category_name,Price:s.Price})
    }
    json.NewEncoder(w).Encode(ProdCats)
}

//view products
func viewProdByCategory(w http.ResponseWriter, r *http.Request) {
	//db := setupDB()
	printMessage("Getting products Details...")
	var pagecount,_=strconv.Atoi(r.URL.Query().Get("page"))
	var cat_id,_=strconv.Atoi(r.URL.Query().Get("category_id"))
	if pagecount==0{
		pagecount=1
	}
	perpage:=9
	fmt.Println(pagecount)
	db := setupDB()
	rows, err := db.Query(`SELECT "Product".product_id,"Product".product_name,"Suppliers".supplier_name,"Categories".category_name,"Product".price
	FROM (("Product" INNER JOIN "Suppliers" ON "Product".supplier_id="Suppliers".supplier_id) INNER JOIN "Categories" ON "Product".category_id="Categories".category_id) 
	WHERE "Product".category_id=$1 LIMIT $2 OFFSET $3`,cat_id,perpage,(pagecount-1)*perpage)
    checkErr(err)
    var ProdCats []ViewProdCat
    // Foreach supplier
    for rows.Next() {
        var s ViewProdCat
        err = rows.Scan(&s.Product_id,&s.Product_name, &s.Supplier_name, &s.Category_name, &s.Price)
        checkErr(err)
        ProdCats= append(ProdCats, ViewProdCat{Product_id: s.Product_id, Product_name: s.Product_name,Supplier_name:s.Supplier_name,Category_name:s.Category_name,Price:s.Price})
    }
    json.NewEncoder(w).Encode(ProdCats)
}

//update product
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var s Product
	_ = json.NewDecoder(r.Body).Decode(&s)
	vars := mux.Vars(r)
	idval, _ := strconv.Atoi(vars["id"])
	db := setupDB()
	printMessage("Updating product...")
	res,e := db.Exec(`UPDATE "Product" SET product_name=$1,supplier_id=$2,category_id=$3, price=$4, expirydate=$5 WHERE product_id=$6`, s.Product_name, s.Supplier_id,s.Category_id,s.Price,s.Expirydate,idval)
	checkErr(e)
	var v,_=res.RowsAffected()
	if v>0{
		s.Product_id=idval
		Products = append(Products, s)
		json.NewEncoder(w).Encode(&s)
	}	
}

//delete product
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var idval,_=strconv.Atoi(vars["id"])
	db := setupDB()
	res,e := db.Exec(`DELETE FROM "Product" WHERE product_id=$1`, idval)
	checkErr(e)
	var rows,_=res.RowsAffected() 
	var response JsonResponse
	if rows>0{
		response = JsonResponse{Type: "success",Message:"Removed product id="+strconv.Itoa(idval)}
	}else{
		response = JsonResponse{Type: "failure",Message:"Cant find product id="+strconv.Itoa(idval)}
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Init router
	r := mux.NewRouter()
	// Route handles & endpoints
	//suppliers CRUD
	r.HandleFunc("/suppliers", GetAllSuppliers).Methods("GET")
	r.HandleFunc("/suppliers/{id}", GetSupplier).Methods("GET")
	r.HandleFunc("/suppliers", createSupplier).Methods("POST")
	r.HandleFunc("/suppliers/{id}", updateSupplier).Methods("PUT")
	r.HandleFunc("/suppliers/{id}", deleteSupplier).Methods("DELETE")
	//categories CRUD
	r.HandleFunc("/categories", createCategory).Methods("POST")
	r.HandleFunc("/categories", viewCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")
	//products
	r.HandleFunc("/product", addProduct).Methods("POST")
	//raw view
	r.HandleFunc("/product", viewProduct).Methods("GET")
	//view with foreign relation column values of 3 tables
	r.HandleFunc("/product/details", viewProductDetails).Methods("GET")
	//localhost:4000\products?page=2
	//List all with pagination
	r.HandleFunc("/product", viewProd).Methods("GET")
	//localhost:4000\products\category?category_id=2
	//List by category and pagination
	r.HandleFunc("/product/category", viewProdByCategory).Methods("GET")
	r.HandleFunc("/product/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")


	// Start server
	log.Println("API is running!")
	log.Fatal(http.ListenAndServe(":4000", r))
	
}
