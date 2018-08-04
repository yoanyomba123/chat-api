package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createRoomUserStore() {
	master := RdbStore(p.database).master()
	rdbCreateRoomUserStore(p.ctx, master)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting room users")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertRoomUsers(p.ctx, master, tx, roomUsers, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting room users")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *mysqlProvider) SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRoomUsers(p.ctx, replica, opts...)
}

func (p *mysqlProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRoomUser(p.ctx, replica, roomID, userID)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRoomUserOfOneOnOne(p.ctx, replica, myUserID, opponentUserID)
}

func (p *mysqlProvider) SelectUserIDsOfRoomUser(opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUserIDsOfRoomUser(p.ctx, replica, opts...)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *model.RoomUser) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating room user")
		logger.Error(err.Error())
		return err
	}

	err = rdbUpdateRoomUser(p.ctx, master, tx, roomUser)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while updating room user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *mysqlProvider) DeleteRoomUsers(opts ...DeleteRoomUsersOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting room users")
		logger.Error(err.Error())
		return err
	}

	err = rdbDeleteRoomUsers(p.ctx, master, tx, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting room users")
		logger.Error(err.Error())
		return err
	}

	return nil
}
