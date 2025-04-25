package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	global_infra "github.com/tinello/golang-openapi/core/global/infrastructure"
	sys_actions "github.com/tinello/golang-openapi/core/system/actions"
)

var providerInstance Provider

func GetProviderInstance() Provider {
	if providerInstance == nil {
		providerInstance = NewProvider()
	}
	return providerInstance
}

func NewProvider() *provider {
	p := &provider{}

	/*
		sqlClient, err := p.SqlClient()
		if err != nil {
			log.Fatalln("Failed to get DB connection:", err)
		}
	*/

	p.getServiceInfo = sys_actions.NewGetServiceInfo()

	return p
}

type provider struct {
	getServiceInfo sys_actions.GetServiceInfo
}

func (p *provider) SqlClient() (*global_infra.SqlClient, error) {
	db, err := p.dbConnection()
	if err != nil {
		return nil, err
	}
	return global_infra.NewSqlClient(db), nil
}

func (p *provider) dbConnection() (*sql.DB, error) {
	dbUrl := MustGetEnv("DB_URL")
	dbUser := MustGetEnv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	return p.postgresConnection(dbUrl, dbUser, dbPassword)
}

func (p *provider) postgresConnection(url string, username string, password string) (*sql.DB, error) {
	return sql.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable&connect_timeout=60",
			username, password, url))
}

func (p *provider) GetServiceInfo() sys_actions.GetServiceInfo {
	return p.getServiceInfo
}

func MustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln("Required environment variable:", name)
	}
	return value
}

//go:generate mockgen -package mocks -destination mocks/MockProvider.go . Provider

type Provider interface {
	GetServiceInfo() sys_actions.GetServiceInfo
}
