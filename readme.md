# gomdelcreator使用说明

## 功能介绍

gomodelcreator用于在使用github.com/whencome/gomodel库时一键生成相关的model代码，暂时只支持mysql数据库。可以同时配置多个库并生成相应的model文件。当对应的model文件存在时，智慧刷新自动生成的相关部分代码。

## 配置示例

```toml
# 数据库配置
[[item]]
# 数据库类型
driver = "mysql"
# 数据库连接
dsn = "root:123456@tcp(localhost:3306)/information_schema?charset=utf8"
# 真实数据库名称
database = "test_db"
# 数据库连接名称，不设置则使用database配置,此值用于在model的GetDatabase()方法中返回
# 服务根据此值的返回获取数据库连接，这里不必是真实的数据库名称，应该是获取数据库连接的标识名称
conn_name = "test_db_conn"
# 输出目录
output_dir = "/mnt/d/Code/Go/gomodeltest"
# 包名
package = "gmtest"
```

## 命令介绍

* 根据配置文件生成所有表的model文件

```sh
go run main.go -f config.toml
```

* 生成指定数据库所有表的model文件，这里的数据库是匹配的database参数

```sh
go run main.go -f config.toml -d test_db
```

* 生成指定表的model文件，这里的数据库是匹配的database参数

```sh
go run main.go -f config.toml -d test_db -t test_table
```