/*
 * @Author: hiLin 123456
 * @Date: 2021-11-08 17:18:56
 * @LastEditors: hiLin 123456
 * @LastEditTime: 2022-12-15 21:59:13
 * @FilePath: /arclinks-go/app/common/util/mongo.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package util

import (
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Find(ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]bson.M, error) {
	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return MongoDataReader(ctx, cur)
}

func MongoDataReader(ctx context.Context, cursor *mongo.Cursor) ([]bson.M, error) {
	defer cursor.Close(ctx)
	records := []bson.M{}

	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		records = append(records, result)
	}
	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return records, nil
}

func Convert(source interface{}, pointer interface{}) error {
	jsonBytes, err := json.Marshal(source)

	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, pointer)

	if err != nil {
		return err
	}

	return nil
}
