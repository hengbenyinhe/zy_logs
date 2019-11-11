# zy_logs

使用golang编写的简单日志库。<br />

# 功能
- zy_log日志库提供6种级别的日志<br />

   1.Debug: 调试程序，日志最详细。但是会影响程序的性能。<br />
   2.Trace: 追踪问题。<br />
   3.Info: 打印日志运行中比较重要的信息，比如状态变化日志。<br />
   4.Warn: 警告，说明程序中出现了潜在的问题。<br />
   5.Error: 错误，程序运行发生了错误。<br />
   6.Access：访问，用户对程序的接口访问日志。<br />
   
- zy_log日志库提供两种日志输出方式<br />

   1.输出到控制台，根据日志等级来显示不同的颜色<br />
   2.输出到文件<br />
   
# 背景
写这个的目的是记录自己学习过程，便于以后自己查阅。在网上找了写资料，然后自己参照别人的教程，自己在添加一些功能。开发的过程记录在不同文件中，分阶段进行讲解。

# 参考资料
[go开发属于自己的日志库](https://juejin.im/user/5814c73f5bbb500059a33944/posts)
