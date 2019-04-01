# grpc连接创建方式性能测试

多年来我一直没弄懂怎样才是grpc连接创建的最佳方式，今天来试试我想象中的3种方式各自的优劣

测试方法：
通过以下3种连接方式发送一定数量的请求，判断效率
- 每次创建新的连接（其实我也知道这个方式肯定最慢）
- 复用同一个conn连接
- 使用同一个stream反复传输

代码直接使用了```github.com/razeencheng/demo-go/tree/master/grpc/demo3```中的server与client，做了简单修改以适应测试
