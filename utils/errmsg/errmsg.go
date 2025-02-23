package errmsg

const (
	SUCCESS = 200
	ERROR   = 500

	// code = 1000 用户模块的错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_TOKEN_EXIST      = 1004
	ERROR_TOKEN_RUNTIME    = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT    = 1008

	// code = 2000 商家模块的错误
	ERROR_SHOP_USED           = 2001
	ERROR_SHOP_PASSWORD_WRONG = 2002
	ERROR_SHOP_NOT_EXIST      = 2003
	ERROR_SHOP_NO_RIGHT       = 2008
	// code = 3000 商品模块的错误
	ERROR_PRODUCT_CREATE_FAIL = 3001
	ERROR_PRODUCT_NOT_EXIST   = 3002
	// code = 4000 订单模块的错误
	ERROR_ORDER_CREATE_FAIL   = 4001
	ERROR_ORDER_INQUIRE_FAIL  = 4002
	ERROR_ORDER_SHOP_NOT_SAME = 4003
	// code = 5000 商品列表的错误
	// code = 6000 钱包的错误
	ERROR_WALLET_CREATE_FAIL        = 6001
	ERROR_WALLET_NOT_EXIST          = 6002
	ERROR_WALLET_PASSWORD_WRONG     = 6003
	ERROR_WALLET_BALANCE_NOT_ENOUGH = 6004
	ERROR_WALLET_EXIST              = 6005

	// code = 8000
	ERROR_MANAGER_NOT_EXIST = 8001
)

var codemsg = map[int]string{
	SUCCESS:                         "OK",
	ERROR:                           "FAIL",
	ERROR_USERNAME_USED:             "用户已存在",
	ERROR_PASSWORD_WRONG:            "密码错误",
	ERROR_USER_NOT_EXIST:            "用户不存在",
	ERROR_TOKEN_EXIST:               "TOKEN不存在",
	ERROR_TOKEN_RUNTIME:             "TOKEN已过期",
	ERROR_TOKEN_WRONG:               "TOKEN不正确",
	ERROR_TOKEN_TYPE_WRONG:          "TOKEN格式错误",
	ERROR_USER_NO_RIGHT:             "该用户无权限",
	ERROR_SHOP_USED:                 "商家已存在",
	ERROR_SHOP_PASSWORD_WRONG:       "商家密码错误",
	ERROR_SHOP_NOT_EXIST:            "商家不存在",
	ERROR_SHOP_NO_RIGHT:             "该商家无权限",
	ERROR_PRODUCT_CREATE_FAIL:       "创建商品失败",
	ERROR_PRODUCT_NOT_EXIST:         "商品不存在",
	ERROR_ORDER_CREATE_FAIL:         "创建订单失败",
	ERROR_ORDER_INQUIRE_FAIL:        "查询订单失败",
	ERROR_ORDER_SHOP_NOT_SAME:       "商家订单不一致",
	ERROR_WALLET_CREATE_FAIL:        "创建钱包失败",
	ERROR_WALLET_NOT_EXIST:          "钱包不存在",
	ERROR_WALLET_EXIST:              "钱包已存在",
	ERROR_WALLET_PASSWORD_WRONG:     "支付密码错误",
	ERROR_WALLET_BALANCE_NOT_ENOUGH: "余额不足",
	ERROR_MANAGER_NOT_EXIST:         "管理员不存在",
}

func GetErrMsg(code int) string {
	return codemsg[code]
}
