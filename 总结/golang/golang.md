# golang

## 1、new 和 make 的区别

* new和make都用于申请内存
* new 要求传入一个类型，它会申请该类型的内存大小，<font color =red>初始化为0值</font>，并返回一个指针指向该内存空间。
* make 也用于分配内存，但是它只用与分配引用类型（chan、slice、map），<font color=red>返回的是类型本身</font>。

## 2、值传递和指针传递有什么区别

* 值传递是值拷贝
* 指针传递创建相同内存地址的副本
* 如果函数内部返回指针，会产生内存逃逸

## 3、内存逃逸分析

1. 什么是内存逃逸
    * 在程序中，原本应该存储到栈的变量由于太大或者当做返回值引用到函数外部时就会发生内存逃逸到堆内存中。简单来说就是<font color=red>局部变量通过堆分配或者回收</font>，就叫内存逃逸
2. 内存逃逸的危害
    * 变量在堆上分配和回收都比栈的开销大的多
    * 增加了GC的压力
    * 容易造成内存碎片
3. 如何分析程序是否发生内存逃逸
    * build时添加-gcflags=-m 选项可分析内存逃逸情况,比如输出./main.go:3:6: moved to heap: x 表示局部变量x逃逸到了堆上。
4. 内存逃逸发生的时机
    * 向channel发送指针数据。因为在编译时，不知道channel会被那个groutine接收，因此编译器没法知道变量何时会被释放，因此只能放入到堆中。

    ```go
    package main
    func main() {
        ch := make(chan int, 1)
        x := 5
        ch <- x  // x不发生逃逸，因为只是复制的值
        ch1 := make(chan *int, 1)
        y := 5
        py := &y
        ch1 <- py  // y逃逸，因为y地址传入了chan中，编译时无法确定什么时候会被接收，所以也无法在函数返回后回收y
    }
    ```

    * 局部变量在函数调用结束后还被其他地方使用，比如函数返回局部变量的指针或者<font color=red>闭包中引用包外的值</font>。因为变量的生命周期大于函数周期，因此只能放入堆中。

    ```go
    package main

    func addN(num int)func(p *int)int{
        n := num
        return func(p *int)int{
            *p += num
            return *p
        }
    }

    func main(){
        f := addN(10)
        i := 11
        f(&i)
        fmt.Println("i=",i)
    }
    ```

    * 在slice或者map 中存储指针。

    ```go
    func test2() {
        i := 10
        var m []*int
        m = append(m, &i)
    }
    ```

    * 切片扩容后长度太大，导致栈空间不足，逃逸到堆上

    ```go
    // test3 内存逃逸3
    func test3() {

        s := make([]int, 10000, 10000)
        for idx, _ := range s {
        s[idx] = idx
        }
    }
    ```

   * 使用接口调用函数。因为函数真正的实现只有在使用的时候知道。

    ```go
    type itf interface{ Get() }
    type me struct{}

    func (me me) Get() {}
    func test4() {
        var i itf = me{}
        i.Get()
    }
    ```

5. 避免内存逃逸的办法
    * 对于小型数据，使用传值而不是传指针，防止内存逃逸。
    * 避免使用长度不固定的切片，在编译期无法确定长度，只能将切片使用堆分配。
    * <font color=red>interface调用方法会发生内存逃逸，热点代码谨慎使用</font>。

## 4、golang的内存管理

  golang的内存管理本质上就是一个<font color=red>内存池</font>，只不过做了很多优化。比如自动伸缩内存池大小，合理的切割块等等。

* ### 内存池 mheap

    golang程序在启动之除，会从操作系统中申请一大块内存作为内存池，这块内存会放在结构体mheap中进行管理，mheap负责将<font color=red>这一块内存切割成不同区域</font>，并将其中一部分的内存切割成合适的大小，分配给用户使用。
  * page页：内存页，一块8k大小的内存空间，<font color=red>Go 与操作系统之间的内存申请和是否都是以页为单位的</font>。
  * span：内存块，<font color=red>一个或者多个连续的page组成span内存块</font>，如果吧page看成工人，则span就是队伍。工人被分为不同的队伍，不同的队伍干不同的活。
  * sizeclass：空间规格，每个span都带有一个sizeclass，<font color=red>标记着该span中的page该如何使用</font>。比如sizeclass标志着span干着什么样的活。
  * object：对象，用来存储一个变量数据的内存空间，一个span在初始化时，会被切割成一堆等大的object。假设object的大小是16B，span的大小是8K，那么就会把span切割成8K/16B=512个object，<font color=red>所谓内存分配就是分配一个Object出去</font>。

* ### mcentral

* ### mcache

## 5、线程有几种模型？

## 6、Goroutine 的原理

## 7、在GPM调度模型，goroutine 有哪几种状态？线程呢？

## 8、如果 goroutine 一直占用资源怎么办，GMP模型怎么解决这个问题

## 9、如果若干个线程发生OOM，会发生什么？Goroutine中内存泄漏的发现与排查？项目出现过OOM吗，怎么解决

## 10、Go数据竞争怎么解决

## 11、goroutine的锁机制了解过吗？Mutex有哪几种模式？Mutex 锁底层如何实现

## 12、Go的垃圾回收算法

## 13、go项目基于alpine 构建docker镜像，启动容器的时候一直报错 standard_init_linux.go:211: exec user process caused "no such file or directory"

问题找到了，一个是静态链接，一个是动态链接，动态链接的在微型镜像alpine上不支持。

总结
默认go使用静态链接，在docker的golang环境中默认是使用动态编译。
如果想使用docker编译+alpine部署，可以通过禁用cgo set CGO_ENABLED=0来解决。
如果要使用cgo可以通过go build --ldflags "-extldflags -static" 来让gcc使用静态编译。
