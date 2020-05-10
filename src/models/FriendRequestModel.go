package models

type FriendRequest struct {
	friendIds []string
}

func AddFriendRequest(from int, to []int)  {
	for _, toId := range to {
		traceDb.Query("insert into friend_request (from_id, to_id, request_status) values (?, ?, '0')", from, toId)
	}
}