package cmd

import (
	"log"

	"demo/ch01/internal/sql2struct"
	"github.com/spf13/cobra"
)

// cmd 全局变量，用于结构外部的命令行参数
// 分别对应用户名、密码、主机地址、编码类型、数据库类型、数据库名称和表名称
var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

// 声明 sql 子命令
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql转换和处理",
	Long:  "sql转换和处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

// 声明 sql 子命令的子命令 struct
var sql2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql转换",
	Long:  "sql转换",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo := &sql2struct.DBInfo{
			DBType:   dbType,
			Host:     host,
			UserName: username,
			Password: password,
			Charset:  charset,
		}
		dbModel := sql2struct.NewDBModel(dbInfo)

		// 连接数据库
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel.Connect err: %v", err)
		}
		// 查询 COLUMNS 表信息
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetColumns err: %v", err)
		}

		// 模板对象的组装与渲染
		template := sql2struct.NewStructTemplate()
		templateColumns := template.AssemblyColumns(columns)
		err = template.Generate(tableName, templateColumns)
		if err != nil {
			log.Fatalf("template.Generate err: %v", err)
		}
	},
}

// 进行默认的 cmd 初始化动作和命令行参数的绑定
func init() {
	sqlCmd.AddCommand(sql2structCmd)
	// 绑定子命令以便设置 Mysql 连接参数
	sql2structCmd.Flags().StringVarP(&username, "username", "", "root", "请输入数据库的账号")
	sql2structCmd.Flags().StringVarP(&password, "password", "", "root", "请输入数据库的密码")
	sql2structCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1:3306", "请输入数据库的HOST")
	sql2structCmd.Flags().StringVarP(&charset, "charset", "", "utf8mb4", "请输入数据库的编码")
	sql2structCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "请输入数据库实例类型")
	sql2structCmd.Flags().StringVarP(&dbName, "db", "", "test", "请输入数据库名称")
	sql2structCmd.Flags().StringVarP(&tableName, "table", "", "test", "请输入表名称")
}
