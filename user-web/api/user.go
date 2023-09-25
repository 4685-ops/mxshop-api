package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// grpc的code 转换为http的状态码
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})

			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Message(),
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": e.Message(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}

		}

		return
	}
}
func HandleValidatorError(c *gin.Context, err error) {
	var errs validator.ValidationErrors
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Translator)),
	})
}

func GetUserList(ctx *gin.Context) {
	ip := global.ServerConfig.UserSrvInfo.Host
	port := global.ServerConfig.UserSrvInfo.Port
	//拨号链接GRPC服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
		return
	}
	//调用接口
	userSrvClient := proto.NewUserClient(userConn)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})

	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表] 失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)

	for _, value := range rsp.Data {
		//data := make(map[string]interface{})
		//data["id"] = value.Id
		//data["name"] = value.NikeName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		user := response.UserResponse{
			Id:       value.Id,
			NikeName: value.NikeName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Gender: value.Gender,
			Mobile: value.Mobile,
		}
		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(c *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	err := c.ShouldBind(&passwordLoginForm)
	if err != nil {
		// 如何返回错误信息
		HandleValidatorError(c, err)
	}
}
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
