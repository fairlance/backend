package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	dbName string
	dbUser string
	dbPass string
)

type Freelancer struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func main() {
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.Parse()

	db, err := gorm.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", dbName, dbUser, dbPass))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	freelancers(db)

	fmt.Println("done.")
}

func freelancers(db *gorm.DB) {
	// statement := "SELECT row_to_json((SELECT d FROM (SELECT id, first_name, last_name, email) d)) FROM freelancers as t"
	statement := `
		SELECT row_to_json(
			(SELECT d FROM (
				SELECT
					t.id AS "id",
					t.first_name AS "firstName",
					t.last_name AS "lastName",
					t.email AS "email"
				) d
			)
		)
		FROM freelancers as t
	`

	rows, err := db.Raw(statement).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var freelancers = []Freelancer{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			fmt.Println(err)
		}
		var freelancer = Freelancer{}
		if err := json.Unmarshal(data, &freelancer); err != nil {
			fmt.Println(err)
		}

		freelancers = append(freelancers, freelancer)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("Found %d freelancers.\n", len(freelancers))

	// opts := MappingOptions{
	// 	Index: "fairlance",
	// 	Type:  "freelancer",
	// }

	// client, err := elastic.NewClient()
	// if err != nil {
	// 	panic(err)
	// }

	// exs := client.IndexExists(opts.Index)
	// if ok, err := exs.Do(); !ok {
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("Creating...")
	// 	if _, err := client.CreateIndex(opts.Index).Do(); err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("done.")
	// } else {
	// 	fmt.Println("Deleting...")
	// 	if _, err := client.DeleteIndex(opts.Index).Do(); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("done. Creating...")
	// 	if _, err := client.CreateIndex(opts.Index).Do(); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("done.")
	// }

	// fmt.Println("Putting mappings...")

	// _, err = client.PutMapping().Index(opts.Index).Type(opts.Type).BodyString(string(mappings[opts.Type])).Do()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("done mapping.")

	// fmt.Println("Adding freelancers...")
	// for _, freelancer := range freelancers {
	// 	f, err := json.Marshal(freelancer)
	// 	fmt.Println("Adding " + string(f) + " ...")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	_, err = client.Index().
	// 		Index(opts.Index).
	// 		Type(opts.Type).
	// 		BodyJson(string(f)).
	// 		Refresh(true).
	// 		Do()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

type MappingOptions struct {
	Index string
	Type  string
}

var mappings = map[string]string{
	"freelancer": `{
		    "freelancer" : {
		        "dynamic": "strict",
		        "properties" : {
	                "id" : {
	                        "type" : "integer",
	                        "index" : "not_analyzed"
	                },
	                "firstName" : {
	                        "type" : "string",
	                        "index" : "not_analyzed"
	                },
	                "lastName" : {
	                        "type" : "string",
	                        "index" : "not_analyzed"
	                },
	                "email" : {
	                        "type" : "string",
	                        "index" : "not_analyzed"
	                }
		        }
		    }
		  }`,
}
