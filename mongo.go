package main

import (
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"
)

var (
	db  *mgo.Session
	err error
)

func GetConnUrl() string {
	url := "mongodb://" + mongoUsername + ":" + mongoPassword + "@"
	hostSlice := make([]string, 0, 5)
	for _, v := range mongoAddrList {
		h := v.Host + ":" + v.Port
		hostSlice = append(hostSlice, h)
	}
	url += strings.Join(hostSlice, ",") + "/" + mongoDatabase
	return url
}

func getConn() (*mgo.Session, error) {
	if db != nil && db.Ping() == nil {
		return db.Clone(), nil
	}
	fmt.Println(GetConnUrl())
	db, err = mgo.Dial(GetConnUrl())
	if err != nil {
		return db, err
	}
	return db.Clone(), nil
}
