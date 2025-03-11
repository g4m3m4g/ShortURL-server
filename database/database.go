package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"encoding/base64"
	"simpleurl/config"
	"log"
)

// SetupDatabase initializes the database connection with TLS
func SetupDatabase() (*gorm.DB, error) {
	// Load CA certificate
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	// Decode the CA cert from the environment variable
	caCertBase64 := os.Getenv("DB_CA_CERT")
	if caCertBase64 == "" {
		return nil, fmt.Errorf("CA_CERT_BASE64 environment variable is not set")
	}

	// Decode the base64 string into bytes
	caCert, err := base64.StdEncoding.DecodeString(caCertBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CA certificate: %v", err)
	}
	rootCertPool := x509.NewCertPool()
	rootCertPool.AppendCertsFromPEM(caCert)

	// Register TLS config
	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}
	err = mysqlDriver.RegisterTLSConfig("tidb", tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to register TLS config: %v", err)
	}

	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=tidb&parseTime=true", 
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	
	return db, nil
}
