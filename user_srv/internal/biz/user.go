package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Greeter is a Greeter model.
type User struct {
	ID int32 
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt time.Time
	IsDeleted bool
	Mobile string 
	Password string 
	NickName string 
	Birthday *time.Time
	Gender string 
	Role int 
}

type BizPassWordCheckInfo struct{
	Password string
	EncryptedPassword string
}

type BizCheckResponse struct{
	Success bool
}


type PageReq struct{
	Pn uint32
	PSize uint32
}

type UserLIstRep struct{
	Total int32
	Data []*User
}

// GreeterRepo is a Greater repo.
type UserRepo interface {
	GetUserList(context.Context, *PageReq) (*UserLIstRep,error)
	GetUserByMobile(context.Context,string) (*User,error)
	GetUserById(context.Context,int32) (*User,error)
	CreateUser(context.Context,*User) (*User,error)
	UpdateUser(context.Context,*User)(error)
	CheckPassWord(context.Context,*BizPassWordCheckInfo)(*BizCheckResponse,error)
}

// GreeterUsecase is a Greeter usecase.
type UserSrvUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserSrvUsecase(repo UserRepo, logger log.Logger) *UserSrvUsecase {
	return &UserSrvUsecase{repo: repo, log: log.NewHelper(logger)}
}

// get user list
func (uc *UserSrvUsecase) GetUserLIst(ctx context.Context, req *PageReq) (*UserLIstRep, error) {
	//uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	//return uc.repo.Save(ctx, g)
	return uc.repo.GetUserList(ctx,req)
}

// Search user by mobile
func (uc *UserSrvUsecase) GetUserByMobile(ctx context.Context, mobile string) (*User, error) {
	//uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	//return uc.repo.Save(ctx, g)
	return uc.repo.GetUserByMobile(ctx,mobile)
}

// Search user by ID
func (uc *UserSrvUsecase) GetUserById(ctx context.Context,  Id int32) (*User, error) {
	//uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	//return uc.repo.Save(ctx, g)
	return uc.repo.GetUserById(ctx,Id)
}

func(uc *UserSrvUsecase) CreateUser(ctx context.Context,user *User) (*User,error){
	return uc.repo.CreateUser(ctx,user)
}

func(uc *UserSrvUsecase) UpdateUser(ctx context.Context,user *User)(error){
	return uc.UpdateUser(ctx,user)
}


func(uc *UserSrvUsecase) CheckPassWord (ctx context.Context,checkRequest *BizPassWordCheckInfo) (*BizCheckResponse,error){
	return uc.repo.CheckPassWord(ctx,checkRequest)
}