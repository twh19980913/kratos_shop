package data

import (
	"context"
	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
	"time"
	"strings"
	"fmt"
	"user_srv/internal/biz"
	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type BaseModel struct{
	ID int32 `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `grom:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct{
	BaseModel
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male表示男'"`
	Role int `gorm:"column:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func Paginate(page,pageSize int) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB {
		if page == 0{
			page = 1
		}
		switch{
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func UserToBizUser (user User) biz.User {
	bizUser := biz.User{
		ID: user.ID,
		Password: user.Password,
		Mobile: user.Mobile,
		NickName: user.NickName,
		Gender: user.Gender,
		Role: user.Role,
		Birthday: user.Birthday,
	}
	return bizUser
}

func(r *userRepo) GetUserList(ctx context.Context,req *biz.PageReq) (*biz.UserLIstRep,error){
	// get user list
	var users []User
	result := r.data.db.Find(&users)
	if result.Error != nil{
		return nil,result.Error
	}
	rsp := &biz.UserLIstRep{}
	rsp.Total = int32(result.RowsAffected)
	
	r.data.db.Scopes(Paginate(int(req.Pn),int(req.PSize))).Find(&users)

	for _,user := range users{
		userInfoRsp := UserToBizUser(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp,nil
}

func (r *userRepo)GetUserByMobile(ctx context.Context,mobile string) (*biz.User,error)  {
	// sear user by mobile
	var user User
	result := r.data.db.Where(&User{Mobile: mobile}).First(&user)
	if result.RowsAffected == 0{
		// return nil,status.Errorf(codes.NotFound,"user not found")
		return nil,errors.New(404,"USER_NOT_FOUND","user not found")
	}
	if result.Error != nil{
		return nil,result.Error
	}

	userInfoRsp := UserToBizUser(user)
	return &userInfoRsp,nil
}

func (r *userRepo)GetUserById(ctx context.Context,Id int32) (*biz.User,error)  {
	// sear user by ID
	var user User
	result := r.data.db.First(&user,Id)
	if result.RowsAffected == 0{
		// return nil,status.Errorf(codes.NotFound,"user not found")
		return nil,errors.New(404,"USER_NOT_FOUND","user not found")
	}
	if result.Error != nil{
		return nil,result.Error
	}

	userInfoRsp := UserToBizUser(user)
	return &userInfoRsp,nil
}

func (r *userRepo) CreateUser(ctx context.Context,req *biz.User) (*biz.User,error){
	// create user
	var user User
	result := r.data.db.Where(&User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1{
		return nil,errors.New(500,"USER_IS_ALREADY_EXISTS","user is alreay exists")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	// password secret
	options := &password.Options{SaltLen: 16,Iterations: 100,KeyLen: 32,HashFunction: sha512.New}
	salt,encodedPwd := password.Encode(req.Password,options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
	user.Password = newPassword

	result = r.data.db.Create(&user)
	if result.Error != nil{
		return nil,result.Error
	}

	userInfoRsp := UserToBizUser(user)
	return &userInfoRsp,nil
}

func (r *userRepo) UpdateUser (ctx context.Context,req *biz.User) (error){
	// owner user update
	var user User
	result := r.data.db.First(&user,req.ID)
	if result.RowsAffected == 0{
		return errors.New(404,"USER_NOT_FOUND","user not found")
	}

	user.NickName = req.NickName
	user.Birthday = req.Birthday
	user.Gender = req.Gender

	result = r.data.db.Save(&user)
	if result.Error != nil{
		return errors.New(500,result.Error.Error(),result.Error.Error())
	}
	return nil
}

func (r *userRepo) CheckPassWord(ctx context.Context,req *biz.BizPassWordCheckInfo) (*biz.BizCheckResponse,error){
	// vertify password
	options := &password.Options{SaltLen: 16,Iterations: 100,KeyLen: 32,HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword,"$")
	check:= password.Verify(req.Password,passwordInfo[2],passwordInfo[3],options)
	return &biz.BizCheckResponse{Success: check},nil
}