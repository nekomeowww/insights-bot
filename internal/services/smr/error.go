package smr

import "errors"

var (
	ErrNoLink = errors.New("没有找到链接，可以发送一个有效的链接吗？用法：/smr <链接>")
	ErrParse  = errors.New("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
	ErrScheme = errors.New("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
)
