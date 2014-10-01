// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package usermanager

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/state"
)

var logger = loggo.GetLogger("juju.apiserver.usermanager")

func init() {
	common.RegisterStandardFacade("UserManager", 0, NewUserManagerAPI)
}

// UserManager defines the methods on the usermanager API end point.
type UserManager interface {
	AddUser(args params.AddUsers) (params.AddUserResults, error)
	DeactivateUser(args params.DeactivateUsers) (params.ErrorResults, error)
	SetPassword(args params.EntityPasswords) (params.ErrorResults, error)
	UserInfo(args params.UserInfoRequest) (params.UserInfoResults, error)
}

// UserManagerAPI implements the user manager interface and is the concrete
// implementation of the api end point.
type UserManagerAPI struct {
	state      *state.State
	authorizer common.Authorizer
}

var _ UserManager = (*UserManagerAPI)(nil)

func NewUserManagerAPI(
	st *state.State,
	resources *common.Resources,
	authorizer common.Authorizer,
) (*UserManagerAPI, error) {
	if !authorizer.AuthClient() {
		return nil, common.ErrPerm
	}

	return &UserManagerAPI{
		state:      st,
		authorizer: authorizer,
	}, nil
}

func (api *UserManagerAPI) permissionCheck(user names.UserTag) error {
	// TODO(thumper): Change this permission check when we have real
	// permissions. For now, only the owner of the initial environment is able
	// to add users.
	initialEnv, err := api.state.InitialEnvironment()
	if err != nil {
		return errors.Trace(err)
	}
	if user != initialEnv.Owner() {
		return errors.Trace(common.ErrPerm)
	}
	return nil
}

// AddUser adds a user.
func (api *UserManagerAPI) AddUser(args params.AddUsers) (params.AddUserResults, error) {
	result := params.AddUserResults{
		Results: make([]params.AddUserResult, len(args.Users)),
	}
	if len(args.Users) == 0 {
		return result, nil
	}
	loggedInUser, err := api.getLoggedInUser()
	if err != nil {
		return result, errors.Wrap(err, common.ErrPerm)
	}
	// TODO(thumper): Change this permission check when we have real
	// permissions. For now, only the owner of the initial environment is able
	// to add users.
	if err := api.permissionCheck(loggedInUser); err != nil {
		return result, errors.Trace(err)
	}
	for i, arg := range args.Users {
		user, err := api.state.AddUser(arg.Username, arg.DisplayName, arg.Password, loggedInUser.Id())
		if err != nil {
			err = errors.Annotate(err, "failed to create user")
			result.Results[i].Error = common.ServerError(err)
		} else {
			result.Results[i].Tag = user.Tag().String()
		}
	}
	return result, nil
}

func (api *UserManagerAPI) getUser(tag string) (*state.User, error) {
	userTag, err := names.ParseUserTag(tag)
	if err != nil {
		return nil, errors.Trace(err)
	}
	user, err := api.state.User(userTag)
	if err != nil {
		return nil, errors.Wrap(err, common.ErrPerm)
	}
	return user, nil
}

// DeactivateUser either disables or enables a user based on the params.
func (api *UserManagerAPI) DeactivateUser(args params.DeactivateUsers) (params.ErrorResults, error) {
	result := params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.Users)),
	}
	if len(args.Users) == 0 {
		return result, nil
	}
	loggedInUser, err := api.getLoggedInUser()
	if err != nil {
		return result, errors.Wrap(err, common.ErrPerm)
	}
	// TODO(thumper): Change this permission check when we have real
	// permissions. For now, only the owner of the initial environment is able
	// to add users.
	if err := api.permissionCheck(loggedInUser); err != nil {
		return result, errors.Trace(err)
	}

	for i, arg := range args.Users {
		user, err := api.getUser(arg.Tag)
		if err != nil {
			result.Results[i].Error = common.ServerError(err)
			continue
		}
		if arg.Deactivate {
			err = user.Deactivate()
		} else {
			err = user.Activate()
		}
		if err != nil {
			action := "activate"
			if arg.Deactivate {
				action = "deactivate"
			}
			result.Results[i].Error = common.ServerError(errors.Errorf("failed to %s user: %s", action, err))
		}
	}
	return result, nil
}

// UserInfo returns information on a user.
func (api *UserManagerAPI) UserInfo(args params.UserInfoRequest) (params.UserInfoResults, error) {
	// TODO(thumper): If no specific users are specified
	// we need to return all the users in the database,
	// just showing the enabled ones unless specified.
	results := params.UserInfoResults{
		Results: make([]params.UserInfoResult, len(args.Entities)),
	}

	for i, arg := range args.Entities {
		user, err := api.getUser(arg.Tag)
		if err != nil {
			results.Results[i].Error = common.ServerError(err)
			continue
		}
		results.Results[i] = params.UserInfoResult{
			Result: &params.UserInfo{
				Username:       user.Name(),
				DisplayName:    user.DisplayName(),
				CreatedBy:      user.CreatedBy(),
				DateCreated:    user.DateCreated(),
				LastConnection: user.LastLogin(),
				Deactivated:    user.IsDeactivated(),
			},
		}
	}

	return results, nil
}

func (api *UserManagerAPI) setPassword(loggedInUser names.UserTag, arg params.EntityPassword, permErr error) error {
	user, err := api.getUser(arg.Tag)
	if err != nil {
		return errors.Trace(err)
	}
	if loggedInUser != user.UserTag() {
		// Look to see if the logged in user is admin.
		if permErr != nil {
			return permErr
		}
	}
	if arg.Password == "" {
		return errors.New("can not use an empty password")
	}
	err = user.SetPassword(arg.Password)
	if err != nil {
		return errors.Annotate(err, "failed to set password")
	}
	return nil
}

func (api *UserManagerAPI) SetPassword(args params.EntityPasswords) (params.ErrorResults, error) {
	result := params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.Changes)),
	}
	if len(args.Changes) == 0 {
		return result, nil
	}
	loggedInUser, err := api.getLoggedInUser()
	if err != nil {
		return result, common.ErrPerm
	}
	permErr := api.permissionCheck(loggedInUser)
	for i, arg := range args.Changes {
		if err := api.setPassword(loggedInUser, arg, permErr); err != nil {
			result.Results[i].Error = common.ServerError(err)
		}
	}
	return result, nil
}

func (api *UserManagerAPI) getLoggedInUser() (names.UserTag, error) {
	switch tag := api.authorizer.GetAuthTag().(type) {
	case names.UserTag:
		return tag, nil
	default:
		return names.UserTag{}, errors.New("authorizer not a user")
	}
}
