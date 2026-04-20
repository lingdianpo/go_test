package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go_test/web/forms"
	"go_test/web/global"
	"go_test/web/global/reponse"
	"go_test/web/middlewares"
	"go_test/web/models"
	"go_test/web/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": "err.Error",
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Tarns)),
	})
}

func GetUserList(ctx *gin.Context) {
	// 从注册中心获取到用户服务的信息

	//claims, _ := ctx.Get("claims")
	//currentUser := claims.(*models.CustomClaims)
	//zap.S().Infof("访问用户：%d", currentUser.ID)
	//zap.S().Info("用户列表页")
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("p_size", "2")
	pSizeInt, _ := strconv.Atoi(pSize)
	zap.S().Info("pnInt:", pnInt)
	zap.S().Info("pSizeInt:", pSizeInt)

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			Nickname: value.NickName,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(c *gin.Context) {
	// 表单验证
	passWordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passWordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	zap.S().Info("验证码" + passWordLoginForm.Captcha)
	zap.S().Info("验证码id" + passWordLoginForm.CaptchaId)
	if !store.Verify(passWordLoginForm.CaptchaId, passWordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	//登录的逻辑
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		// 只是查询到用户了而已，并没有检查密码
		if passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passWordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"password": "登录失败",
			})
		} else {
			zap.S().Info("位置1")
			if passRsp.Success {
				// 生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),                          //签名的生效时间
						ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), //30天过期
						Issuer:    "imooc",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"msg":        "登录成功",
					"id":         rsp.Id,
					"token":      token,
					"nick_name":  rsp.NickName,
					"expired_at": time.Now().Add(time.Hour*24*30).Unix() * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误",
				})
			}

		}
	}

}
func RedisTest(c *gin.Context) {
	// 1. 将数据存储到 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	err := rdb.Set(context.Background(), "15922060532", "123456", time.Second*30).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to set data in Redis",
		})
		return
	}

	// 2. 读取 Redis 数据
	val, err := rdb.Get(context.Background(), "15922060532").Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Key does not exist",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get data from Redis",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": val,
	})

}
func Register(c *gin.Context) {
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	//验证码校验
	//rdb := redis.NewClient(&redis.Options{
	//	Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	//})
	//val, err := rdb.Get(context.Background(), "15922060532").Result()
	//if err != nil {
	//	if err == redis.Nil {
	//		c.JSON(http.StatusNotFound, gin.H{
	//			"msg": "验证码错误",
	//		})
	//	}
	//	return
	//}
	//if val != registerForm.Code {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"msg": "验证码错误",
	//	})
	//}

	//拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error())
	}
	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	user, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   registerForm.Mobile,
		Password: registerForm.PassWord,
		NickName: registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 创建 【新建用户失败】失败 %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.Id,
		"nick_name": user.NickName,
	})

}

func GoroutineTest(c *gin.Context) {

}
