// package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// 	routers "github.com/sofc-t/task_manager/task8/router"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Couldn't load .env file")
// 	}

// 	MongoURI := os.Getenv("db_mongo_uri")
// 	DbMongoName := os.Getenv("db_mongo_name")
// 	if MongoURI == "" || DbMongoName == "" {
// 		log.Fatal("Couldn't find MongoDB URI or DbMongoName in .env")
// 	}

// 	clientOptions := options.Client().ApplyURI(MongoURI)
// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal("Couldn't create MongoDB client")
// 	}
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to ping MongoDB: %v", err)
// 	}

// 	r := routers.SetUpRouter(1000*time.Second, *client.Database(DbMongoName))

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		log.Fatal("Port not Specified")
// 	}
// 	r.Run(":" + port)
// }



package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	
	"os/exec"
	"path/filepath"
	"strings"
)

// Info about each Go file
type FileInfo struct {
	FileName   string
	Functions  []string
	Structs    []string
	RoutePaths []string
}

func main() {
	root := "./" // The root directory to start the search
	var fileInfos []FileInfo

	// Walk through the directory tree
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			fileInfo := parseGoFile(path)
			fileInfos = append(fileInfos, fileInfo)
			runGoFile(path) // Optionally run the file
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Output collected information
	for _, fileInfo := range fileInfos {
		fmt.Printf("File: %s\n", fileInfo.FileName)
		fmt.Printf("Functions: %v\n", fileInfo.Functions)
		fmt.Printf("Structs: %v\n", fileInfo.Structs)
		fmt.Printf("Route Paths: %v\n\n", fileInfo.RoutePaths)
	}
}

// Parse the Go file to extract functions, structs, and routes
func parseGoFile(filePath string) FileInfo {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file %s: %v", filePath, err)
	}

	var fileInfo FileInfo
	fileInfo.FileName = filePath

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			fileInfo.Functions = append(fileInfo.Functions, x.Name.Name)
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {
				fileInfo.Structs = append(fileInfo.Structs, x.Name.Name)
				// Optionally parse struct fields here if needed
			}
		case *ast.CallExpr:
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				// Looking for Gin or similar route definitions
				if fun.Sel.Name == "GET" || fun.Sel.Name == "POST" || fun.Sel.Name == "PUT" || fun.Sel.Name == "DELETE" {
					if len(x.Args) > 0 {
						if lit, ok := x.Args[0].(*ast.BasicLit); ok {
							fileInfo.RoutePaths = append(fileInfo.RoutePaths, lit.Value)
						}
					}
				}
			}
		}
		return true
	})

	return fileInfo
}

// Run a Go file and print its output (optional)
func runGoFile(filePath string) {
	cmd := exec.Command("go", "run", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to run file %s: %v\n", filePath, err)
	}
	fmt.Printf("Output of %s:\n%s\n", filePath, string(output))
}
