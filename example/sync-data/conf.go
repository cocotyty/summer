package main

var Conf = struct {
	RedisAddr string
	MysqlDSN  string
}{}

func init() {
	Conf.MysqlDSN = ""
	Conf.RedisAddr = ""
}
