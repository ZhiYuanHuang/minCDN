# minCDN
* 借鉴minio源码结构编写的一个基于minio对象存储的实现静态文件（图片、js等）CDN节点服务的golang项目
* 项目工程实现文件二进制缓存，访问后台存储时命名空间锁，http返回cache-control max-age,LastModify,etag等header实现客户端缓存
