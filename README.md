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

### env&conf

```bash
# write an app conf file like this
# app.toml
APPName = "go-liziwei01-appui"
[HTTPServer]
Listen="0.0.0.0:{env.LISTEN_PORT|8080}"
```

```bash
import (
    "path/filepath"
    "github.com/liziwei01/go-liziwei01-library/library/conf"
    "github.com/liziwei01/go-liziwei01-library/library/env"
)

const (
    appConfPath = "./conf/app.toml"
)

type Config struct {
    APPName string
    Env env.AppEnv
    HTTPServer struct {
        Listen       string
    }
}

func GetConfig(filePath string) (*Config, error) {
    confPath, err := filepath.Abs(filePath)
    if err != nil {
        return nil, err
    }
    var c *Config
    # here is the main part
    if err := conf.Parse(confPath, &c); err != nil {
        return nil, err
    }
    rootDir := filepath.Dir(filepath.Dir(confPath))
    opt := env.Option{
        AppName: c.APPName,
        RunMode: c.RunMode,
        RootDir: rootDir,
        DataDir: filepath.Join(rootDir, "data"),
        LogDir:  filepath.Join(rootDir, "log"),
        ConfDir: filepath.Join(rootDir, filepath.Base(filepath.Dir(confPath))),
    }
    c.Env = env.New(opt)
    return c, nil
}

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
    # Config.HTTPServer.Listen is 0.0.0.0:8080 in my conf file
    _ = http.ListenAndServe(Config.HTTPServer.Listen, nil)
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
    res, err := GetPaperSlice(ctx, params)
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

when you

```bash
curl localhost:8080?start_time=0&end_time=100000
```

you will get json return like this

```bash
{
    "data": # the data slice
        [{
            "title": "something",
            "author": "somebody",
            "publish_time": "10000000"
        }]
    "errno": # error number
    "errmsg": # error message
}
```

### mysql

```bash
# write the database conf file like this
# ./conf/servicer/db_liziwei01.toml
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

type PaperInfo struct {
    Title           string `db:"title" json:"title"`
    Authors         string `db:"author" json:"authors"`
    PublishTime     int64  `db:"publish_time" json:"publish_time"`
}

func GetPaperSlice(ctx context.Context, params PaperSearchParams) ([]PaperInfo, error) {
    # init and link the mysql database db_liziwei01
    var res []PaperInfo
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
    if err != nil {
        return nil, err
    }
    return res, nil
}
```

## License

MIT License
