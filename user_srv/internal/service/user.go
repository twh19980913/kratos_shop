package service

import (
	"context"
	"time"
	pb "user_srv/api/helloworld/v1"
	"user_srv/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.UserSrvUsecase
}

func ModelToResponse(user biz.User) pb.UserInfoResponse{
	// in 
	userInfoRsp := pb.UserInfoResponse{
		Id: user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Mobile: user.Mobile,
		Gender: user.Gender,
		Role: int32(user.Role),
	}
	if user.Birthday != nil{
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func NewUserService(uc *biz.UserSrvUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) GetUserLIst(ctx context.Context, req *pb.PageInfo) (*pb.UserListResponse, error) {
	// get user list
	userLIst, err := s.uc.GetUserLIst(ctx, &biz.PageReq{Pn: req.Pn, PSize: req.PSize})
	if err != nil {
		return nil,err
	}
	rsp := &pb.UserListResponse{}
	rsp.Total = userLIst.Total

	for _,user := range userLIst.Data{
		userInfoRsp := ModelToResponse(*user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp,nil
}
func (s *UserService) GetUserMobile(ctx context.Context, req *pb.MobileRequest) (*pb.UserInfoResponse, error) {
	u, err := s.uc.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil,err
	}
	resp := ModelToResponse(*u)
	return &resp, nil
}
func (s *UserService) GetUserById(ctx context.Context, req *pb.IdRequest) (*pb.UserInfoResponse, error) {
	u, err := s.uc.GetUserById(ctx, req.Id)
	if err != nil {
		return nil,err
	}
	resp := ModelToResponse(*u)
	return &resp, nil
}
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserInfo) (*pb.UserInfoResponse, error) {
	user,err := s.uc.CreateUser(ctx,&biz.User{
		NickName: req.NickName,
		Mobile: req.Mobile,
		Password: req.PassWord,
	})
	if err != nil {
		return nil,err
	}
	userInfoRsp := ModelToResponse(*user)
	return &userInfoRsp, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserInfo) (*emptypb.Empty, error) {
	birthDay := time.Unix(int64(req.Birthday),0)
	err := s.uc.UpdateUser(ctx,&biz.User{
		NickName: req.NickName,
		Gender: req.Gender,
		Birthday: &birthDay,
	})
	if err != nil {
		return nil,err
	}
	return &emptypb.Empty{},nil
}
func (s *UserService) CheckPassWord(ctx context.Context, req *pb.PasswordCheckInfo) (*pb.CheckResponse, error) {
	check,err := s.uc.CheckPassWord(ctx,&biz.BizPassWordCheckInfo{
		Password: req.Password,
		EncryptedPassword: req.EncryptedPassword,
	})
	if err != nil {
		return nil,err
	}
	return &pb.CheckResponse{Success: check.Success}, nil
}
