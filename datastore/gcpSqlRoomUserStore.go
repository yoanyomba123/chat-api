package datastore

import (
	scpb "github.com/swagchat/protobuf"
)

func (p *gcpSQLProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *gcpSQLProvider) DeleteAndInsertRoomUsers(roomUsers []*scpb.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSQLProvider) InsertRoomUsers(roomUsers []*scpb.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSQLProvider) SelectRoomUser(roomID, userID string) (*scpb.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *gcpSQLProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*scpb.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *gcpSQLProvider) SelectRoomUsersByRoomID(roomID string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *gcpSQLProvider) SelectRoomUsersByUserID(userID string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *gcpSQLProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *gcpSQLProvider) UpdateRoomUser(roomUser *scpb.RoomUser) (*scpb.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *gcpSQLProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
