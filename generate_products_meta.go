package main

import (
        "fmt"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
)

type Product struct {
        ID              int
        Name            string
        Category        string
}

func dbInit() *sql.DB {
        db, err := sql.Open("mysql", "login:password@/name_of_db")
        if err != nil {
                panic(err)
        }

        return db
}

func getProducts(db *sql.DB) *[]Product {
        results, err := db.Query("SELECT A.product_id, A.name, C.name AS category FROM oc_product_description A JOIN oc_product_to_category B ON A.product_id = B.product_id JOIN oc_category_description C ON B.category_id = C.category_id")
        if err != nil {
                panic(err.Error())
        }
        defer results.Close()

        var products []Product

        for results.Next() {
                var product Product

                err = results.Scan(&product.ID, &product.Name, &product.Category)
                if err != nil {
                    panic(err.Error())
                }

                fmt.Printf("%d %s %s\n", product.ID, product.Name, product.Category)
                products = append(products, product)
        }

        return &products
}

func generateMeta(db *sql.DB, products *[]Product) {
        for _, product := range *products {
                _, err := db.Exec(fmt.Sprintf("UPDATE oc_product_description SET meta_title = \"Buy %s on site ...\", meta_description = \"%s. Fast delivery. Low price\" WHERE product_id = %d", product.Name, product.Category, product.ID))
                if err != nil {
                    panic(err.Error())
                }
        }
}

func main() {
        db := dbInit()
        generateMeta(db, getProducts(db))

        defer db.Close()
}
