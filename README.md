# go-liziwei01-library
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

This is a school project base written by liziwei from CUHK(SZ)

## Clone & Install Hook

git clone https://github.com/liziwei01/go-liziwei01-library.git &&\
wget https://github.com/liziwei01/hooks/archive/refs/tags/1.0.tar.gz &&\
tar -xzvf 1.0.tar.gz &&\
mv hooks-1.0/commit-msg go-liziwei01-library/.git/hooks &&\
rm -rf hooks-1.0 &&\
rm 1.0.tar.gz &&\
cd go-liziwei01-library

## Run

go run main.go

## Use

# mysql
conf file under ./conf/servicer\
SAMPLE:\
db_liziwei01.toml\

```bash
Username = "root"
Password = "liziwei01"
DbName = "db_liziwei01"
DbDriver = "mysql"
Host = "localhost"
Port = 3306
```

```bash
import "github.com/liziwei01/go-liziwei01-library/model/mysql"
# init and link the database
mysql.GetMysqlClient(ctx, "db_liziwei01")
where := map[string]interface{}{
&emsp;&emsp;"_orderby":        "score desc",
&emsp;&emsp;"_limit":          []uint{intStart, params.PageLength},
&emsp;&emsp;"publish_time >=": params.StartTime,
&emsp;&emsp;"publish_time <=": params.EndTime,
&emsp;}
```

## License

MIT License
