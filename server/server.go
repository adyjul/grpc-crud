// server/server.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"grpc_crud/proto/crud"
)

const (
	port    = ":50051"
	dbHost  = "localhost"
	dbPort  = "3306"
	dbUser  = "root"
	dbPass  = ""
	dbName  = "coba_grpc"
)

type server struct {
	db *sql.DB
	crud.UnimplementedCrudServiceServer // Embed the UnimplementedCrudServiceServer
}

func (s *server) Create(ctx context.Context, req *crud.CreateRequest) (*crud.CreateResponse, error) {
	// Implementasi logika tambah data ke database di sini
	_, err := s.db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", req.Name, req.Age)
    if err != nil {
        return nil, err
    }
	return &crud.CreateResponse{Success: true, Message: "Data created successfully"}, nil
}

func (s *server) Read(ctx context.Context, req *crud.ReadRequest) (*crud.ReadResponse, error) {
	// Implementasi logika membaca data dari database di sini
	  // Membaca data dari tabel users
	  var name string
	  var age int
  
	  err := s.db.QueryRow("SELECT name, age FROM users WHERE id=?", req.Id).Scan(&name, &age)
	  if err != nil {
		  return nil, err
	  }
  
	return &crud.ReadResponse{Name: "John Doe", Age: 25}, nil
}

func (s *server) Update(ctx context.Context, req *crud.UpdateRequest) (*crud.UpdateResponse, error) {
	// Implementasi logika update data di database di sini
	_, err := s.db.Exec("UPDATE users SET name=?, age=? WHERE id=?", req.Name, req.Age, req.Id)
    if err != nil {
        return nil, err
    }

	return &crud.UpdateResponse{Success: true, Message: "Data updated successfully"}, nil
}

func (s *server) Delete(ctx context.Context, req *crud.DeleteRequest) (*crud.DeleteResponse, error) {
	// Implementasi logika hapus data dari database di sini
	_, err := s.db.Exec("DELETE FROM users WHERE id=?", req.Id)
    if err != nil {
        return nil, err
    }

	return &crud.DeleteResponse{Success: true, Message: "Data deleted successfully"}, nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Printf("Server listening on port %s\n", port)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	s := grpc.NewServer()
	crud.RegisterCrudServiceServer(s, &server{db: db})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
