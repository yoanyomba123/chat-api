package datastore

import (
	"log"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbRoomUserCreateStore() {
	tableMap := dbMap.AddTableWithName(models.RoomUser{}, TABLE_NAME_ROOM_USER)
	tableMap.SetKeys(true, "id")
	tableMap.SetUniqueTogether("room_id", "user_id")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}
func RdbRoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		result := StoreResult{}

		if err := dbMap.Insert(roomUser); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
		}
		result.Data = roomUser

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUsersInsert(roomId string, roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		trans, err := dbMap.Begin()
		result := StoreResult{}

		defer func() {
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback.", err)
			}
			close(storeChannel)
		}()

		if isDeleteFirst {
			query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
			params := map[string]interface{}{"roomId": roomId}
			_, err = trans.Exec(query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while deleting all room's user item.", err)
				storeChannel <- result
			}
		}

		for _, roomUser := range roomUsers {
			if err := trans.Insert(roomUser); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
				storeChannel <- result
			}
		}

		if result.ProblemDetail == nil {
			if err := trans.Commit(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
			}
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUserUsersSelect(roomId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var users []*models.User
		query := utils.AppendStrings("SELECT u.* ",
			"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
			"LEFT JOIN ", TABLE_NAME_USER, " AS u ",
			"ON ru.user_id = u.user_id ",
			"WHERE room_id = :roomId;")
		params := map[string]interface{}{"roomId": roomId}
		_, err := dbMap.Select(&users, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user list.", err)
		}
		result.Data = users

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUsersSelect(roomId *string, userIds []string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var roomUsers []*models.RoomUser
		var userIdsQuery string
		var userIdsParams map[string]interface{}
		var roomIdParams map[string]interface{}
		var params map[string]interface{}
		if userIds != nil {
			userIdsQuery, userIdsParams = utils.MakePrepareForInExpression(userIds)
		}
		if roomId != nil {
			roomIdParams = map[string]interface{}{"roomId": roomId}
		}
		if userIdsParams == nil {
			params = roomIdParams
		}
		if roomIdParams == nil {
			params = userIdsParams
		}
		if userIdsParams != nil && roomIdParams != nil {
			params = utils.MergeMap(userIdsParams, roomIdParams)
		}

		query := utils.AppendStrings("SELECT * ",
			"FROM ", TABLE_NAME_ROOM_USER,
			" WHERE ")
		if roomId != nil {
			query = utils.AppendStrings(query, " room_id=:roomId")
		}
		if roomId != nil && userIds != nil {
			query = utils.AppendStrings(query, " AND ")
		}
		if userIds != nil {
			query = utils.AppendStrings(query, " user_id IN (", userIdsQuery, ")")
		}
		_, err := dbMap.Select(&roomUsers, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user list.", err)
		}
		result.Data = roomUsers

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUsersSelectUserIds(roomId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var roomUsers []string

		query := utils.AppendStrings("SELECT user_id ",
			"FROM ", TABLE_NAME_ROOM_USER,
			" WHERE room_id=:roomId;")
		params := map[string]interface{}{"roomId": roomId}
		_, err := dbMap.Select(&roomUsers, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user list.", err)
		}
		result.Data = roomUsers

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var roomUserIds []int
		var userIdsQuery string
		var userIdsParams map[string]interface{}
		var roomIdParams map[string]interface{}
		if userIds != nil {
			userIdsQuery, userIdsParams = utils.MakePrepareForInExpression(userIds)
		}
		if roomId != nil {
			roomIdParams = map[string]interface{}{"roomId": roomId}
		}
		params := utils.MergeMap(userIdsParams, roomIdParams)

		query := utils.AppendStrings("SELECT id ",
			"FROM ", TABLE_NAME_ROOM_USER,
			" WHERE ")
		if roomId != nil {
			query = utils.AppendStrings(query, " room_id=:roomId")
		}
		if roomId != nil && userIds != nil {
			query = utils.AppendStrings(query, " AND ")
		}
		if userIds != nil {
			query = utils.AppendStrings(query, " user_id IN (", userIdsQuery, ")")
		}
		_, err := dbMap.Select(&roomUserIds, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user ids.", err)
		}
		result.Data = roomUserIds

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUserSelect(roomId, userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var roomUsers []*models.RoomUser
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id=:userId;")
		params := map[string]interface{}{
			"roomId": roomId,
			"userId": userId,
		}
		if _, err := dbMap.Select(&roomUsers, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user item.", err)
		}
		if len(roomUsers) == 1 {
			result.Data = roomUsers[0]
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUserUpdate(roomUser *models.RoomUser) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		_, err := dbMap.Update(roomUser)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room's user item.", err)
		}
		result.Data = roomUser

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUserDelete(roomId string, userIds []string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var err error
		var query string
		var params map[string]interface{}
		if userIds == nil {
			query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
			params = map[string]interface{}{"roomId": roomId}
		} else {
			var userIdsQuery string
			userIdsQuery, params = utils.MakePrepareForInExpression(userIds)
			params["roomId"] = roomId
			query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
		}
		_, err = dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user item.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUsersDeleteByUserIds(roomId *string, userIds []string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		userIdsQuery, params := utils.MakePrepareForInExpression(userIds)
		if roomId != nil {
			params["roomId"] = roomId
		}
		query := utils.AppendStrings("DELETE ",
			"FROM ", TABLE_NAME_ROOM_USER,
			" WHERE user_id IN (", userIdsQuery, ")")
		if roomId != nil {
			query = utils.AppendStrings(query, " AND room_id=:roomId")
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user list.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=unread_count+1 WHERE room_id=:roomId AND user_id!=:userId;")
		params := map[string]interface{}{
			"roomId": roomId,
			"userId": currentUserId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room's user unread count.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbRoomUserMarkAllAsRead(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=0 WHERE user_id=:userId;")
		params := map[string]interface{}{
			"userId": userId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while mark all as read.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}
