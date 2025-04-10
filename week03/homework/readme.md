此目录存放本周课后作业，可以在此文件添加作业设计思路和流程图等


结构设计：

  使用Word-Translation-Phrase三级结构精确映射JSON层级

  omitempty标签处理可选字段

性能优化：

  事务批量提交（相比自动提交模式，速度提升约10-100倍）

  预处理语句复用（减少SQL解析开销）

健壮性处理：

  defer tx.Rollback()确保事务安全

  INSERT OR IGNORE避免唯一约束冲突

  详细的错误上下文信息（包含出错单词）

数据序列化：

  显式处理JSON序列化错误

  动态处理可选phrases字段

资源管理：

  defer db.Close()保证连接释放

  统一删除旧数据库文件（避免残留数据干扰）