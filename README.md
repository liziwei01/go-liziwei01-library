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

### mysql

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
import (
    "context"
    "github.com/liziwei01/go-liziwei01-library/model/mysql"
)

type PaperSearchParams struct {
    StartTime  int64  `json:"start_time"`
    EndTime    int64  `json:"end_time"`
}

func GetPaperSlice(ctx context.Context, params PaperSearchParams)
# init and link the mysql database db_liziwei01
var res []paperModel.PaperInfo
mysql.GetMysqlClient(ctx, "db_liziwei01")
where := map[string]interface{}{
        "_orderby":        "score desc",
        "_limit":          []uint{0, 10},
        "publish_time >=": params.StartTime,
        "publish_time <=": params.EndTime,
    }
columns := []string{"title", "authors"}
# query the 
# table `tb_paper_info` 
# for `columns` 
# with `where` and
# data will be stored in `res` slice
err = client.Query(ctx, "tb_paper_info", where, columns, &res)
```

### ghttp

```bash
import (
    "context"
    "net/http"
    errBase "github.com/liziwei01/go-liziwei01-library/model/error"
    "github.com/liziwei01/go-liziwei01-library/model/ghttp"
)

type PaperSearchParams struct {
    StartTime  int64  `json:"start_time"`
    EndTime    int64  `json:"end_time"`
}

func init() {
    _ = http.ListenAndServe("0.0.0.0:8080", nil)
    http.HandleFunc("/paperList", GetPaperList)
}

func GetPaperList(response http.ResponseWriter, request *http.Request) {
    # initialize ghttp
    g := ghttp.Default(&request, &response)
    # get front end params
    params, err := getPaperInput(ctx, g)
    if err != nil {
        g.Write(params, errBase.ErrorNoClient, err)
    }
    # GetPaperSlice get data from mysql and returns a data slice `res`
    res, err := GetPaperSlice(ctx, g)
    if err != nil {
        g.Write(res, errBase.ErrorNoServer, err)
    }
    # return
    g.Write(res, errBase.ErrorNoSuccess, err)
}

func getPaperInput(ctx context.Context, g ghttp.Ghttp) {
    return PaperSearchParams{
        # also support g.Post()
        EndTime:    g.Get("end_time"),
        Journal:    g.Get("start_time"),
    }
}
```

```bash
# when you 
curl localhost:8080?start_time=0&end_time=100000
# you will get json return like this
{
    "data": # the data slice
    "errno": # error number
    "errmsg": # error message
}
```

## License

MIT License
