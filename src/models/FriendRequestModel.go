package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/enhuizhu/gps-tracking-go-backend/src/helpers"
)

// AddFriendRequest insert frined request into the table
func AddFriendRequest(from int, toID int) bool {
	// check if the request is already there
	if IsRequestExist(from, toID) {
		return false
	}

	traceDb.Query("insert into friend_request (from_id, to_id, request_status) values (?, ?, '0')", from, toID)
	return true
}

// IsRequestExist check if there is existing request
func IsRequestExist(from int, toID int) bool {
	var number int
	err := traceDb.QueryRow("select count(*) from friend_request where from_id=? and to_id=? and request_status='0'", from, toID).Scan(&number)

	fmt.Println(err)
	if err != nil {
		return false
	}

	return number > 0
}

// IsIdExist check if the id already exist in the table or not
func IsIdExist(userID int) bool {
	var number int

	err := traceDb.QueryRow("select count(*) from user_login where userId = ?", userID).Scan(&number)

	if err != nil {
		panic(err.Error())
	}

	return number > 0
}

// IsFriendsRecordExist check if the friends records is inside the table
func IsFriendsRecordExist(userID int) bool {
	var number int

	err := traceDb.QueryRow("select count(*) from friends where userId = ?", userID).Scan(&number)

	if err != nil {
		panic(err.Error())
	}

	return number > 0
}

// GetFriendIDs get friendIDs base on user id
func GetFriendIDs(userID int) ([]int, error) {
	var friendIds []int
	var friendsStr string
	err := traceDb.QueryRow("select friends from friends where userId = ?", userID).Scan(&friendsStr)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(friendsStr), friendIds)

	if err != nil {
		return nil, err
	}

	return friendIds, nil
}

// AreTheyFriends check if they are already friends.
func AreTheyFriends(id1 int, id2 int) bool {
	if IsIdExist(id1) && IsIdExist(id2) {
		friendIds, err := GetFriendIDs(id1)

		if err != nil {
			return false
		}

		return helpers.ArrayContain(friendIds, id2)
	}

	return false
}

// GetUserIDAndFriendIDBaseOnRequestID get the user id base on the reqeust id
func GetUserIDAndFriendIDBaseOnRequestID(requestID int) (int, int) {
	var userID int
	var friendID int
	err := traceDb.QueryRow("select from_id, to_id from friend_request where id=?", requestID).Scan(&userID, &friendID)

	if err != nil {
		panic(err.Error())
	}

	return userID, friendID
}

// AcceptFriendRequest accept the request sent by user's friend
func AcceptFriendRequest(requestID int) error {
	// get user id base on request id
	userID, friendID := GetUserIDAndFriendIDBaseOnRequestID(requestID)
	recordExist := IsFriendsRecordExist(userID)

	var ctx context.Context

	db := traceDb.CreateCon()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return err
	}

	if !recordExist {
		friendIDs := []int{friendID}
		friendIDsStr, err := helpers.JSONStringify(friendIDs)

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "insert into friends (userId, friends) values (?, ?)", userID, friendIDsStr)

		if err != nil {
			return err
		}
	} else {
		friendIDs, err := GetFriendIDs(userID)

		if !helpers.ArrayContain(friendIDs, friendID) {
			friendIDs[len(friendIDs)] = friendID
		}

		if err != nil {
			return err
		}

		friendIDsStr, err := helpers.JSONStringify(friendIDs)

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "update friends set friends='?' where userId=?", friendIDsStr, userID)

		if err != nil {
			return err
		}
	}

	_, err = db.ExecContext(ctx, "update friend_request set request_status='1' where id=?", requestID)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
