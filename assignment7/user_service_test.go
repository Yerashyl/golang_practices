package assignment7_test

import (
	"errors"
	"testing"

	"assignment7"
	"assignment7/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := assignment7.NewUserService(mockRepo)

	tests := []struct {
		name    string
		user    *assignment7.User
		email   string
		setup   func()
		wantErr bool
		errString string
	}{
		{
			name:  "User already exists",
			user:  &assignment7.User{ID: 1, Name: "Test", Email: "test@example.com"},
			email: "test@example.com",
			setup: func() {
				mockRepo.EXPECT().GetByEmail("test@example.com").Return(&assignment7.User{}, nil)
			},
			wantErr: true,
			errString: "user with this email already exists",
		},
		{
			name:  "New User -> Success",
			user:  &assignment7.User{ID: 2, Name: "New", Email: "new@example.com"},
			email: "new@example.com",
			setup: func() {
				mockRepo.EXPECT().GetByEmail("new@example.com").Return(nil, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Repository error on CreateUser",
			user:  &assignment7.User{ID: 3, Name: "Err", Email: "err@example.com"},
			email: "err@example.com",
			setup: func() {
				mockRepo.EXPECT().GetByEmail("err@example.com").Return(nil, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			errString: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.RegisterUser(tt.user, tt.email)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := assignment7.NewUserService(mockRepo)

	tests := []struct {
		name    string
		id      int
		newName string
		setup   func()
		wantErr bool
	}{
		{
			name:    "Empty name",
			id:      1,
			newName: "",
			setup:   func() {},
			wantErr: true,
		},
		{
			name:    "User not found/repo error",
			id:      999,
			newName: "NewName",
			setup: func() {
				mockRepo.EXPECT().GetUserByID(999).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:    "Successful update",
			id:      1,
			newName: "NewName",
			setup: func() {
				user := &assignment7.User{ID: 1, Name: "OldName"}
				mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *assignment7.User) error {
					assert.Equal(t, "NewName", u.Name)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name:    "UpdateUser Fails",
			id:      1,
			newName: "NewName",
			setup: func() {
				user := &assignment7.User{ID: 1, Name: "OldName"}
				mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.UpdateUserName(tt.id, tt.newName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := assignment7.NewUserService(mockRepo)

	tests := []struct {
		name    string
		id      int
		setup   func()
		wantErr bool
	}{
		{
			name:    "Attempt to delete admin",
			id:      1,
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "Successfull delete",
			id:   2,
			setup: func() {
				mockRepo.EXPECT().DeleteUser(2).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Repository Error",
			id:   3,
			setup: func() {
				mockRepo.EXPECT().DeleteUser(3).Return(errors.New("deletion failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.DeleteUser(tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
