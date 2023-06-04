# ecodepost
## 论坛样式
![image](./docs/statics/img.png)
## 支持功能
* 文章
  * 流式布局
  * 卡片布局
  * 列表布局
* 问答
* 专栏
* 友链

## 如何启动
### 下载release包
### 创建数据库
```sql
CREATE DATABASE `post_main`;
CREATE DATABASE `post_user`;
```
### 修改配置
```toml
# mysql换成你的ip
# redis换成你的ip
```
### 使用指令
```bash
export EGO_DEBUG=true && ./ecodepost --job=install
```
### 初始化数据
```bash
export EGO_DEBUG=true && ./ecodepost --job=init
```
### 启动网站
```bash
export EGO_DEBUG=true && ./ecodepost 
```

### 访问
访问 https://localhost:9002

## todo
* 合并成一个数据库
* 优化配置
* docker-compose
* 自动生成一个管理员账号
* 接口全部变成/api

