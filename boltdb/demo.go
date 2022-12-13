package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func main() {
	dbHandler, err := bolt.Open("test.db", 0666, nil)
	if err != nil {
		fmt.Println("open test.db err", err)
		return
	}

	dbHandler.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("bucket1"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("bucket1"))
			if err != nil {
				fmt.Println("create bucket err", err)
				return err
			}
		}

		err := bucket.Put([]byte("name"), []byte("lily"))
		if err != nil {
			fmt.Println("put name failed", err)
			return err
		}

		name1 := bucket.Get([]byte("name"))
		name2 := bucket.Get([]byte("name2"))
		fmt.Println("name1:", string(name1))
		fmt.Println("name2:", string(name2))
		return nil
	})

	defer dbHandler.Close()
}
