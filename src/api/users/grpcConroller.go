package users

import (
	"context"

	"github.com/myrachanto/grpcgateway/pb"
	middle "github.com/myrachanto/grpcgateway/src/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var Usergapi UserServiceServer = &userGapiController{}

type UserServiceServer interface {
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	LoginUser(context.Context, *pb.LoginUserRequest) (*pb.LoginUserResponse, error)
	LogoutUser(context.Context, *pb.LogoutRequest) (*pb.LogoutResponse, error)
	GetOneUser(context.Context, *pb.GetOneRequest) (*pb.GetOneResponse, error)
	GetAllUser(context.Context, *pb.GetAllRequest) (*pb.GetAllResponse, error)
	UpdateUser(context.Context, *pb.UpdateRequest) (*pb.UpdateResponse, error)
	DeleteUser(context.Context, *pb.DeleteRequest) (*pb.DeleteResponse, error)
	pb.UnsafeUserServiceServer
}
type userGapiController struct {
	service UserServiceInterface
	pb.UnimplementedUserServiceServer
}

func NewUserGapiController(ser UserServiceInterface) UserServiceServer {
	return &userGapiController{
		service: ser,
	}
}

func (g *userGapiController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &User{}
	user.Firstname = req.GetFirstname()
	user.Lastname = req.GetLastname()
	user.Phone = req.GetPhone()
	user.Username = req.Username
	user.Address = req.GetAddress()
	user.Email = req.GetEmail()
	user.Password = req.Password
	u, err1 := g.service.Create(user)
	if err1 != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", err1.Code())
	}
	return &pb.CreateUserResponse{
		User: converter(u),
	}, nil
}
func converter(user *User) *pb.User {
	return &pb.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Phone:     user.Phone,
		Username:  user.Username,
		Address:   user.Address,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.Created_At),
		UpdatedAt: timestamppb.New(user.Updated_At),
	}
}

func (g *userGapiController) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user := &LoginUser{}
	user.Email = req.GetEmail()
	user.Password = req.GetPassword()
	auth, err1 := g.service.Login(user)
	if err1 != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", err1.Code())
	}
	return &pb.LoginUserResponse{
		Usercode:            auth.Usercode,
		UserName:            auth.UserName,
		Token:               auth.Token,
		RefleshToken:        auth.RefleshToken,
		Role:                auth.Role,
		Picture:             auth.Picture,
		SessionCode:         auth.SessionCode,
		TokenExpires:        timestamppb.New(auth.TokenExpires),
		RefleshTokenExpires: timestamppb.New(auth.RefleshTokenExpires),
	}, nil
}
func (g *userGapiController) LogoutUser(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	token := req.GetToken()
	_, problem := g.service.Logout(token)
	if problem != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", problem.Code())
	}
	return &pb.LogoutResponse{
		Info: "logout succesifuly",
	}, nil
}
func (g *userGapiController) GetOneUser(ctx context.Context, req *pb.GetOneRequest) (*pb.GetOneResponse, error) {
	_, err := middle.GRPCAuthorization(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", err)
	}
	code := req.GetCode()
	user, problem := g.service.GetOne(code)
	if problem != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", problem.Code())
	}
	return &pb.GetOneResponse{
		User: converter(user),
	}, nil
}
func (g *userGapiController) GetAllUser(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {

	search := req.GetSearch()
	users, problem := g.service.GetAll(search)
	if problem != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", problem.Code())
	}
	us := []*pb.User{}
	for _, u := range users {
		us = append(us, converter(u))
	}
	return &pb.GetAllResponse{
		User: us,
	}, nil
}

func (g *userGapiController) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	user := &User{}
	user.Firstname = req.GetFirstname()
	user.Lastname = req.GetLastname()
	user.Phone = req.GetPhone()
	user.Username = req.Username
	user.Address = req.GetAddress()
	user.Email = req.GetEmail()
	user.Usercode = req.GetUsercode()
	code := user.Usercode
	u, problem := g.service.Update(code, user)
	if problem != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", problem.Code())
	}
	return &pb.UpdateResponse{
		User: converter(u),
	}, nil
}
func (g *userGapiController) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	id := req.GetCode()
	success, failure := g.service.Delete(id)
	if failure != nil {
		return nil, status.Errorf(codes.Internal, "Error : %v", failure.Code())
	}
	return &pb.DeleteResponse{
		Info: success,
	}, nil
}
