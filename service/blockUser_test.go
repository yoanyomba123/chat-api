package service

import (
	"log"
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameCreateBlockUsers = "create block users test"
	TestNameGetBlockUsers    = "get block users test"
	TestNameGetBlockedUsers  = "get blocked users test"
	TestNameDeleteBlockUsers = "delete block users test"
)

func TestBlockUser(t *testing.T) {
	t.Run(TestNameCreateBlockUsers, func(t *testing.T) {
		req := &model.CreateBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := CreateBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsers)
		}
	})

	t.Run(TestNameGetBlockUsers, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
		if len(blockUsers.BlockUserIDs) != 3 {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}

		req.ResponseType = scpb.ResponseType_UserList
		blockUsers, errRes = GetBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
		if len(blockUsers.BlockUsers) != 3 {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
	})

	t.Run(TestNameGetBlockedUsers, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "service-user-id-0002"
		blockedUsers, errRes := GetBlockedUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}
		log.Printf("%#v\n", blockedUsers)
		if len(blockedUsers.BlockedUserIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}

		req.ResponseType = scpb.ResponseType_UserList
		blockedUsers, errRes = GetBlockedUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}
		if len(blockedUsers.BlockedUsers) != 1 {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}
	})

	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
		req := &model.DeleteBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := DeleteBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
	})
}