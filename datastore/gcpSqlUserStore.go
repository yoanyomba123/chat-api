package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (provider GcpSqlProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (provider GcpSqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) StoreResult {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (provider GcpSqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (provider GcpSqlProvider) SelectUsers() StoreResult {
	return RdbSelectUsers()
}

func (provider GcpSqlProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (provider GcpSqlProvider) UpdateUser(user *models.User) StoreResult {
	return RdbUpdateUser(user)
}

func (provider GcpSqlProvider) UpdateUserDeleted(userId string) StoreResult {
	return RdbUpdateUserDeleted(userId)
}

func (provider GcpSqlProvider) SelectContacts(userId string) StoreResult {
	return RdbSelectContacts(userId)
}
