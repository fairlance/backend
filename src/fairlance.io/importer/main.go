package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/olivere/elastic.v3"
)

var (
	dbName string
	dbUser string
	dbPass string
)

type FreelancerRow struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

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

	statement := "SELECT row_to_json((SELECT d FROM (SELECT id, first_name, last_name, email) d)) FROM freelancers as t"
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
		var freelancerRow = FreelancerRow{}
		if err := json.Unmarshal(data, &freelancerRow); err != nil {
			fmt.Println(err)
		}

		freelancers = append(freelancers, Freelancer{
			freelancerRow.Id,
			freelancerRow.FirstName,
			freelancerRow.LastName,
			freelancerRow.Email,
		})
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("Found %d freelancers.\n", len(freelancers))

	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	opts := struct {
		Index string
		Type  string
	}{
		"fairlance",
		"freelancer",
	}

	exs := client.IndexExists(opts.Index)
	if ok, err := exs.Do(); !ok {
		if err != nil {
			panic(err)
		}

		fmt.Println("Creating...")
		if _, err := client.CreateIndex(opts.Index).Do(); err != nil {
			panic(err)
		}

		fmt.Println("done.")
	} else {
		fmt.Println("Deleting...")
		if _, err := client.DeleteIndex(opts.Index).Do(); err != nil {
			panic(err)
		}
		fmt.Println("done. Creating...")
		if _, err := client.CreateIndex(opts.Index).Do(); err != nil {
			panic(err)
		}
		fmt.Println("done.")
	}

	fmt.Println("Putting mappings...")

	mapping := `{
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
		  }`

	_, err = client.PutMapping().Index(opts.Index).Type(opts.Type).BodyString(string(mapping)).Do()
	if err != nil {
		panic(err)
	}

	fmt.Println("done.")

	fmt.Println("Adding freelancers...")
	for _, freelancer := range freelancers {
		f, err := json.Marshal(freelancer)
		fmt.Println("Adding " + string(f) + " ...")
		if err != nil {
			panic(err)
		}
		_, err = client.Index().
			Index(opts.Index).
			Type(opts.Type).
			BodyJson(string(f)).
			Refresh(true).
			Do()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("done.")
}
