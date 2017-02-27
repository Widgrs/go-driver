//
// DISCLAIMER
//
// Copyright 2017 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var (
	endpoint string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:8529", "URL used to connect to the database")
}

type Book struct {
	Title   string `json:"title"`
	NoPages int    `json:"no_pages"`
}

func main() {
	flag.Parse()
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{endpoint},
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})

	// Create database
	db, err := c.CreateDatabase(nil, "examples_books", nil)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Create collection
	col, err := db.CreateCollection(nil, "books", nil)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// Create documents
	books := []Book{
		Book{
			Title:   "ArangoDB Getting Started Guide",
			NoPages: 12,
		},
		Book{
			Title:   "ArangoDB Cookbook",
			NoPages: 257,
		},
		Book{
			Title:   "ArangoDB HTTP Manual",
			NoPages: 310,
		},
	}
	metas, errs, err := col.CreateDocuments(nil, books)
	if err != nil {
		log.Fatalf("Failed to create documents: %v", err)
	} else if err := errs.FirstNonNil(); err != nil {
		log.Fatalf("Failed to create documents: first error: %v", err)
	}

	fmt.Printf("Created documents with keys '%s' in collection '%s' in database '%s'\n", strings.Join(metas.Keys(), ","), col.Name(), db.Name())
}
