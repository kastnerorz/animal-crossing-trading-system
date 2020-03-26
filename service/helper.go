package main

import "go.mongodb.org/mongo-driver/bson"

func Sel(query map[string]string) bson.M {
	result := make(bson.M, len(query))
	for k, v := range query {
		result[k] = v
	}
	return result
}
