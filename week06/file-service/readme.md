健康监测接口
http://localhost:8081/api/health
文件上传接口
http://localhost:8081/files/upload
文件下载接口
http://localhost:8081/files/download
文件预览接口
http://localhost:8081/files/preview

数据库建表语句
CREATE TABLE file_service (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT NOT NULL,        -- 原始文件名
    size      INTEGER CHECK(size > 0), 
    uuid      TEXT UNIQUE NOT NULL, -- 存储系统唯一ID
    type      TEXT NOT NULL        -- MIME类型
);

db.go文件代码包含数据库操作(新建数据库，对数据库增删查改)
main.go文件包含gin框架接口部署及操作方法

在post请求进行upload上传的同时将文件信息插入数据库(包含name,size,uuid,type)
通过指定uuid可进行对应的下载和预览操作

由于postman报错该文件不在您的工作目录中. 与您共享此请求的队友将无法使用此文件. 为了简化协作, 您可以在设置中设置工作目录.导致未能成功运行样例