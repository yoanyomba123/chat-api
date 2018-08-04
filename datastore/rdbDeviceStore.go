package datastore

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateDeviceStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateDeviceStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	tableMap := dbMap.AddTableWithName(model.Device{}, tableNameDevice)
	tableMap.SetUniqueTogether("user_id", "platform")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" || columnMap.ColumnName == "notification_device_id" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating device table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertDevice(ctx context.Context, dbMap *gorp.DbMap, device *model.Device) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertDevice", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	if err := dbMap.Insert(device); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting device. %v.", err))
		return err
	}

	return nil
}

func rdbSelectDevices(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectDevicesOption) ([]*model.Device, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDevices", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectDevicesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID == "" && opt.platform == scpb.Platform_PlatformNone && opt.token == "" {
		err := errors.New("An error occurred while getting devices. Be sure to specify either userId or platform or token")
		logger.Error(err.Error())
		return nil, err
	}

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE ", tableNameDevice)
	params := map[string]interface{}{}

	if opt.userID != "" {
		query = fmt.Sprintf("%s user_id=:userId AND", query)
		params["userId"] = opt.userID
	}

	if opt.platform != scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s platform=:platform AND", query)
		params["platform"] = opt.platform
	}

	if opt.token != "" {
		query = fmt.Sprintf("%s token=:token AND", query)
		params["token"] = opt.token
	}

	query = query[0 : len(query)-len(" AND")]

	_, err := dbMap.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting devices. %v.", err))
		return nil, err
	}

	return devices, nil
}

func rdbSelectDevice(ctx context.Context, dbMap *gorp.DbMap, userID string, platform scpb.Platform) (*model.Device, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDevice", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND platform=:platform;", tableNameDevice)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err := dbMap.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting device. %v.", err))
		return nil, err
	}

	if len(devices) == 1 {
		return devices[0], nil
	}

	return nil, nil
}

func rdbSelectDevicesByUserID(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]*model.Device, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDevicesByUserID", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId;", tableNameDevice)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting devices by userId. %v.", err))
		return nil, err
	}

	return devices, nil
}

func rdbSelectDevicesByToken(ctx context.Context, dbMap *gorp.DbMap, token string) ([]*model.Device, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDevicesByToken", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=:token;", tableNameDevice)
	params := map[string]interface{}{
		"token": token,
	}
	_, err := dbMap.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting device by token. %v.", err))
		return nil, err
	}

	return devices, nil
}

func rdbUpdateDevice(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, device *model.Device) error {
	span := tracer.Provider(ctx).StartSpan("rdbUpdateDevice", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	deleted := time.Now().Unix()
	err := rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(deleted),
		DeleteSubscriptionsOptionFilterByUserID(device.UserID),
		DeleteSubscriptionsOptionFilterByPlatform(device.Platform),
	)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET token=?, notification_device_id=? WHERE user_id=? AND platform=?;", tableNameDevice)
	_, err = tx.Exec(query, device.Token, device.NotificationDeviceID, device.UserID, device.Platform)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while updating device. %v.", err))
		return err
	}

	return nil
}

func rdbDeleteDevices(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteDevicesOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteDevices", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := deleteDevicesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID == "" && opt.platform == scpb.Platform_PlatformNone {
		err := errors.New("An error occurred while deleting devices. Be sure to specify either userID or platform")
		logger.Error(err.Error())
		return err
	}

	err := rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(opt.logicalDeleted),
		DeleteSubscriptionsOptionFilterByUserID(opt.userID),
		DeleteSubscriptionsOptionFilterByPlatform(opt.platform),
	)
	if err != nil {
		return err
	}

	var query string
	if opt.logicalDeleted != 0 {
		query = fmt.Sprintf("UPDATE %s SET deleted=%d WHERE", tableNameDevice, opt.logicalDeleted)
	} else {
		query = fmt.Sprintf("DELETE FROM %s WHERE", tableNameDevice)
	}

	if opt.userID != "" && opt.platform == scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s user_id=?", query)
		_, err := tx.Exec(query, opt.userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting devices")
			logger.Error(err.Error())
			return err
		}
	}

	if opt.userID != "" && opt.platform != scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s user_id=? AND platform=?", query)
		_, err := tx.Exec(query, opt.userID, opt.platform)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting devices")
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}
