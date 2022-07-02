## Golang 编程设计

### 流？I/O操作？阻塞？epoll?

#### 多线程/多进程解决阻塞场景的方法：

##### 1. 非阻塞、忙轮询

![非阻塞忙轮询](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306130335371.png "非阻塞忙轮询")

```go
while true {
	for i in 流[] {
		if i has 数据 {
			读 或者 其他处理
		}
	}
}
```

非阻塞忙轮询的方式，可以让用户分别与每个快递员取得联系，宏观上来看，是同时可以与多个快递员沟通(并发效果)、 但是快递员在于用户沟通时耽误前进的速度(浪费 CPU )。

##### 2. select

![select](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306130430837.png "select")

开设一个代收网点，让快递员全部送到代收点。这个网点的管理员叫 select。这样用户（CPU）就可以在家休息了，麻烦的事交给 select 就好了。当有快递的时候，select负责给用户打电话，在此期间用户在家休息睡觉就好了。

但 select 代收员比较懒，她记不住快递员的单号，还有快递货物的数量。她只会告诉用户快递到了，但是是谁到的，需要用户挨个快递员问一遍（不清楚是哪个线程的数据）。

##### 3. epoll

![epoll](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306130910868.png "epoll")

epoll 的服务态度要比 select 好很多，在通知用户（CPU）的时候，不仅告诉用户有几个快递到了，还分别告诉用户是谁谁谁。只需要按照epoll给的答复，来询问快递员取快递即可（清楚是哪个线程的数据）。

epoll 与 select，poll 一样，是对 I/O 多路复用的技术，只关心“活跃”的链接，无需遍历全部描述符集合，能够处理大量的链接请求(系统可以打开的文件数目)。

#### epoll API

##### 1. 创建 epoll

```c
/** 
 * @param size 告诉内核监听的数目 
 * 
 * @returns 返回一个 epoll 句柄（即一个文件描述符） 
 */
int epoll_create(int size);
```

```c
int epfd = epoll_create(1000);
```

![创建 epoll](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306131313148.png "创建 epoll")

创建一个epoll句柄，实际上是在内核空间，建立一个root根节点，这个根节点的关系与 epfd 相对应。

##### 2. 控制 epoll

```C
/**
* @param epfd 用 epoll_create 所创建的 epoll 句柄
* @param op 表示对 epoll 监控描述符控制的动作
*
* EPOLL_CTL_ADD(注册新的fd到epfd)
* EPOLL_CTL_MOD(修改已经注册的fd的监听事件)
* EPOLL_CTL_DEL(epfd删除一个fd)
*
* @param fd 需要监听的文件描述符
* @param event 告诉内核需要监听的事件
*
* @returns 成功返回0，失败返回-1, errno查看错误信息
*/
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);


struct epoll_event {
	__uint32_t events; /* epoll 事件 */
	epoll_data_t data; /* 用户传递的数据 */
}

/*
 * events : {EPOLLIN, EPOLLOUT, EPOLLPRI, EPOLLHUP, EPOLLET, EPOLLONESHOT}
 */
typedef union epoll_data {
	void *ptr;
	int fd;
	uint32_t u32;
	uint64_t u64;
} epoll_data_t;
```

```C
struct epoll_event new_event;

new_event.events = EPOLLIN | EPOLLOUT;
new_event.data.fd = 5;

epoll_ctl(epfd, EPOLL_CTL_ADD, 5, &new_event);
```

![控制 epoll](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306132355997.png "控制 epoll")

 创建一个用户态的事件，绑定到某个fd上，然后添加到内核中的epoll红黑树中。

##### 3. 等待 epoll

```C
/**
*
* @param epfd 用epoll_create所创建的epoll句柄
* @param event 从内核得到的事件集合
* @param maxevents 告知内核这个events有多大,
* 注意: 值 不能大于创建epoll_create()时的size.
* @param timeout 超时时间
* -1: 永久阻塞
* 0: 立即返回，非阻塞
* >0: 指定微秒
*
* @returns 成功: 有多少文件描述符就绪,时间到时返回0
* 失败: -1, errno 查看错误
*/
int epoll_wait(int epfd, struct epoll_event *event, int maxevents, int timeout);
```

```C
struct epoll_event my_event[1000];

int event_cnt = epoll_wait(epfd, my_event, 1000, -1);
```

![等待 epoll](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306132548007.png "等待 epoll")

epoll_wait 是一个阻塞的状态，如果内核检测到IO的读写响应，会抛给上层的 epoll_wait, 返回给用户态一个已经触发的事件队列，同时阻塞返回。开发者可以从队列中取出事件来处理，其中事件里就有绑定的对应 `fd `是哪个(之前添加 epoll 事件时已绑定)。

##### 4.  使用epoll编程主流程骨架

```C
int epfd = epoll_crete(1000);

//将 listen_fd 添加进 epoll 中
epoll_ctl(epfd, EPOLL_CTL_ADD, listen_fd,&listen_event);

while (1) {
	//阻塞等待 epoll 中 的fd 触发
	int active_cnt = epoll_wait(epfd, events, 1000, -1);

	for (i = 0 ; i < active_cnt; i++) {
		if (evnets[i].data.fd == listen_fd) {
			//accept. 并且将新accept 的fd 加进epoll中.
		}
		else if (events[i].events & EPOLLIN) {
			//对此fd 进行读操作
		}
		else if (events[i].events & EPOLLOUT) {
			//对此fd 进行写操作
		}
	}
}
```

#### epoll 的触发模式

##### 1. 水平触发

![水平触发](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306151955361.png "水平触发")

水平触发的主要特点是，如果用户在监听`epoll`事件，当内核有事件的时候，会拷贝给用户态事件，但是**如果用户只处理了一次，那么剩下没有处理的会在下一次epoll_wait 再次返回该事件**。这样如果用户永远不处理这个事件，就导致每次都会有该事件从内核到用户的拷贝，耗费性能，但是水平触发相对安全，最起码事件不会丢掉，除非用户处理完毕。

##### 2. 边缘触发

![边缘触发](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306152027545.png "边缘触发")

边缘触发，相对跟水平触发相反，当内核有事件到达， 只会通知用户一次，至于用户处理还是不处理，以后将不会再通知。这样减少了拷贝过程，增加了性能，但是相对来说，如果用户忘记处理，将会产生事件丢失的情况。

#### 简单的 epoll 服务器(C 语言)

略

----------------------------

### 对于操作系统而言进程、线程以及 Goroutine 协程的区别

#### 1. 进程内存

进程，可执行程序运行中形成一个独立的内存体，这个内存体**有自己独立的地址空间(Linux会给每个进程分配一个虚拟内存空间32位操作系统为4G, 64位为很多T)，有自己的堆**，上级挂靠单位是操作系统。操作系统会以进程为单位，分配系统资源（CPU时间片、内存等资源），**进程是资源分配的最小单位**。

![进程结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306133100554.png "进程结构")

#### 2. 线程内存

**线程，有时被称为轻量级进程(Lightweight Process，LWP），是操作系统调度（CPU调度）执行的最小单位**。

![线程结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306133145324.png "线程结构")

**多个线程共同“寄生”在一个进程上，除了拥有各自的栈空间，其他的内存空间都是一起共享。**由于这个特性，使得线程之间的内存关联性很大，互相通信很简单(堆区、全局区等数据都共享，需要加锁机制即可完成同步通信)，但是同时也让线程之间生命体联系较大，比如一个线程出问题，导致进程出问题，也会导致其他线程出问题。

#### 3. 执行单元

对于Linux来讲，不区进程还是线程，他们都是一个单独的执行单位，CPU一视同仁，均分配时间片。所以，如果一个进程想更大程度的与其他进程抢占CPU的资源，那么多开线程是一个好的办法。

![CPU 分配时间片](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306133415648.png "CPU 分配时间片")

如上图，进程A没有开线程，那么默认就是`1个线程`，对于内核来讲，它只有1个`执行单元`，进程B开了`3个线程`，那么在内核中，该进程就占有3个`执行单元`。CPU的视野是只能看见内核的，它不知道谁是进程，谁是线程，哪个线程属于哪个进程，时间片轮询平均调度分配。此时进程B拥有的3个单元就有了资源供给的优势。

#### 4. 切换问题与协程

通过上述的描述，可以知道，线程越多，进程利用(或者)抢占的 CPU 资源就越高。

![时间片轮转](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306133814185.png "时间片轮转")

那么是不是线程可以无限制的多呢？答案当然不是的，当 CPU 在内核态切换一个执行单元的时候，会有一个时间成本和性能开销。其中性能开销至少会有两个开销：切换内核栈、切换硬件上下文。这两个切换会带来的后果和影响是：

+ 保存寄存器中的内容：将之前执行流程的状态保存
+ CPU 高速缓存失效：页表查找是一个很慢的过程，因此通常使用 Cache 来缓存常用的地址映射，这样可以加速页表查找，这个 cache 就是 TLB（ 转译后备缓冲器）。当进程切换后页表也要进行切换，页表切换后 TLB 就失效了，cache失效导致命中率降低，那么虚拟地址转换为物理地址就会变慢，**表现出来的就是程序运行会变慢**。

![CPU 浪费的成本](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306133938693.png "CPU 浪费的成本")

综上，不能够大量的开辟，因为线程执行流程越多，CPU 在切换的时间成本越大。很多编程语言就想了办法，既然不能左右和优化 CPU 切换线程的开销，那么能否让 CPU 内核态不切换执行单元， 而是在用户态切换执行流程呢？

很显然，用户是没权限修改操作系统内核机制的，那么只能在用户态再来一个伪执行单元，这就是协程。

![协程示意图](C:\Users\zyc\AppData\Roaming\Typora\typora-user-images\image-20220306134544116.png "协程示意图")

#### 5. 协程的切换成本

协程切换比线程切换快主要有两点：

（1）协程切换**完全在用户空间进行**，线程切换涉及**特权模式切换，需要在内核空间完成**；

（2）协程切换相比线程切换**做的事情更少**，线程需要有内核和用户态的切换，系统调用过程。

##### 协程切换成本

协程切换非常简单，将**当前协程的 CPU 寄存器状态保存起来，然后将需要切换进来的协程的 CPU 寄存器状态加载的 CPU 寄存器上**就行了。而且**完全在用户态进行**，一般一次协程上下文切换最多是**几十 ns** 量级。

##### 线程切换成本

系统内核调度的对象是线程，因为线程是调度的基本单元（进程是资源拥有的基本单元，进程的切换需要做的事情更多，这里暂时不讨论进程切换），**线程的调度只有拥有最高权限的内核空间才可以完成**，所以线程的切换涉及到**用户空间和内核空间的切换**，也就是特权模式切换，然后需要操作系统调度模块完成线程调度，而且除了和协程相同基本的 CPU 上下文，还有线程的私有栈和寄存器等，简而言之就是上下文比协程多一些，比较下 `task_strcut `和 任何一个协程库的 coroutine 的 struct 结构体大小就能明显区分出来。而且特权模式切换的开销比协程切换开销大很多。

##### 内存占用

进程占用 4g 内存，线程根据不同的操作系统版本占用内存存在差异，但线程基本都是维持 Mb 的量级单位，一般是4~64Mb不等， 多数维持约10M上下，而协程基本维持 kb 的量级单位。

------------------------

### Go是否可以无限go？ 如何限定数量？

#### 1. 不控制goroutine数量引发的问题

goroutine 具备如下两个特点

- 体积轻量
- 优质的GMP调度

那么 goroutine 是否可以无限开辟呢，如果做一个服务器或者一些高业务的场景，能否随意的开辟goroutine并且放养不管呢？让他们自生自灭，毕竟有强大的GC和优质的调度算法支撑？

```go
package main

import (
    "fmt"
    "math"
    "runtime"
)

func main() {
    //模拟用户需求业务的数量
    task_cnt := math.MaxInt64

    for i := 0; i < task_cnt; i++ {
        go func(i int) {
            //... do some busi...

            fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
        }(i)
    }
}
```

```bash
panic: too many concurrent operations on a single file or socket (max 1048575)
```

所以，在迅速开辟 goroutine (**不控制并发的 goroutine 数量** )会在短时间内占据操作系统的资源(CPU、内存、文件描述符等)。

- CPU 使用率浮动上涨
- Memory 占用不断上涨。
- 主进程崩溃（被杀掉了）

这些资源实际上是所有用户态程序共享的资源，所以大批的 goroutine 最终引发的灾难不仅仅是自身，还会关联其他运行的程序。**所以在编写逻辑业务的时候，限制goroutine是必须要重视的问题。**

#### 2. 一些简单方法控制goroutines数量

##### 1. 用有buffer的channel来限制

```go
package main

import (
    "fmt"
    "math"
    "runtime"
)

func busi(ch chan bool, i int) {

    fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
    <-ch
}

func main() {
    // 模拟用户需求业务的数量
    task_cnt := math.MaxInt64
    // task_cnt := 10

    ch := make(chan bool, 3)

    for i := 0; i < task_cnt; i++ { // 循环速度

        ch <- true

        go busi(ch, i)
    }

}
```

从结果看，程序并没有出现崩溃，而是按部就班的顺序执行，并且 go 程的数量控制在了3（4是因为还有一个主线程 main goroutine），从数字上看是不是在跑的 goroutines 有几十万个呢？

![程序内部执行情况](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306140529515.png "程序内部执行情况")

这里使用了 buffer 为3的 channel，在写的过程中，实际上是限制了循环速度。这个速度决定了 go 的创建速度，而 go 的结束速度取决于 `busi()`函数的执行速度。 这样就能够保证了，同一时间内运行的goroutine的数量与buffer的数量一致，从而达到了限定效果。但是这段代码有一个小问题，就是如果将 `go_cnt`的数量变的小一些，会出现打出的结果不正确。

```go
package main

import (
    "fmt"
    //"math"
    "runtime"
)

func busi(ch chan bool, i int) {

    fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
    <-ch
}

func main() {
    //模拟用户需求业务的数量
    //task_cnt := math.MaxInt64
    task_cnt := 10

    ch := make(chan bool, 3)

    for i := 0; i < task_cnt; i++ {

        ch <- true

        go busi(ch, i)
    }
    // main 开辟完 go 程就退出了
}
```

```bash
go func  2  goroutine count =  4
go func  3  goroutine count =  4
go func  4  goroutine count =  4
go func  5  goroutine count =  4
go func  6  goroutine count =  4
go func  1  goroutine count =  4
go func  8  goroutine count =  4
```

是因为`main`将全部的 go 开辟完之后，就立刻退出进程了。所以想全部 go 都执行，需要在 main 的最后进行阻塞操作。

##### 2. 只使用 sync 同步机制

```go
package main

import (
    "fmt"
    "math"
    "sync"
    "runtime"
)

var wg = sync.WaitGroup{}

func busi(i int) {

    fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())
    wg.Done()
}

func main() {
    //模拟用户需求业务的数量
    task_cnt := math.MaxInt64


    for i := 0; i < task_cnt; i++ {
		wg.Add(1)
        go busi(i)
    }

	  wg.Wait()
}
```

```bash
panic: too many concurrent operations on a single file or socket (max 1048575)
```

单纯使用 sync 无法达到控制 goroutine 的数量，所以最终结果依然是崩溃。

##### 3. channel 与 sync 同步组合方式

```go
package main

import (
    "fmt"
    "math"
    "sync"
    "runtime"
)

var wg = sync.WaitGroup{}

func busi(ch chan bool, i int) {

    fmt.Println("go func ", i, " goroutine count = ", runtime.NumGoroutine())

    <-ch // 从 channel 中读取数据,但不使用

    wg.Done()
}

func main() {
    // 模拟用户需求go业务的数量
    task_cnt := math.MaxInt64

    ch := make(chan bool, 3)

    for i := 0; i < task_cnt; i++ {
		wg.Add(1)

        ch <- true // 向 channel 中写入数据

        go busi(ch, i)
    }

	  wg.Wait()
}
```

此时程序不会再造成资源爆炸而崩溃，并且运行 go 的数量控制住了在 buffer 为 3 的这个范围内。

##### 4. 利用无缓冲 channel 与任务发送/执行分离方式

```go
package main

import (
    "fmt"
    "math"
    "sync"
    "runtime"
)

var wg = sync.WaitGroup{}

// 任务执行
func busi(ch chan int) {

    // 只要 channel 有数据就会读取
    for t := range ch {
        fmt.Println("go task = ", t, ", goroutine count = ", runtime.NumGoroutine())
        wg.Done()
    }
}

// 任务发送
func sendTask(task int, ch chan int) {
    wg.Add(1)
    // 向 channel 写数据
    ch <- task
}

func main() {

    ch := make(chan int)   // 无buffer channel

    goCnt := 3              // 启动goroutine的数量
    for i := 0; i < goCnt; i++ {
        //启动go
        go busi(ch)
    }

    taskCnt := math.MaxInt64 // 模拟用户需求业务的数量
    for t := 0; t < taskCnt; t++ {
        // 不断发送任务
        sendTask(t, ch)
    }

	  wg.Wait()
}
```

执行流程大致如下，实际是将任务的发送和执行做了业务上的分离。使得消息出去，输入 `SendTask` 的频率可设置、执行Goroutine 的数量也可设置。也就是既控制输入(生产)，又控制输出(消费)。使得可控更加灵活。这也是很多Go 框架的 Worker 工作池的最初设计思想理念。

![无缓冲 channel 与任务发送/执行分离](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220306142850640.png "无缓冲 channel 与任务发送/执行分离")

-------------------

### 动态保活 Worker 工作池设计

#### 如何知道一个Goroutine已经死亡？

实际上，Go 语言并没有暴露如何知道一个 Goroutine 是否存在的接口，如果要证明一个Go是否存在，可以在子 Goroutine 的业务中，**定期写向一个keep live的Channel**，然后主 Goroutine 来发现当前子 Go的状态。Go 语言在对于 Go 和 Go 之间没有像进程和线程一样有强烈的父子、兄弟等关系，每个 Go 实际上对于调度器都是一个独立的，平等的执行流程。

>PS: 如果是监控子线程、子进程的死亡状态，就没有这么简单了，这里也要感谢 Go 的调度器提供的方便，既然用 Go，就要基于 Go 的调度器来实现该模式。

那么如何知道一个 Goroutine 已经死亡了呢？

##### 子Goroutine

可以通过给一个被监控的 Goroutine 添加一个 `defer` ，然后`recover()` 捕获到当前 Goroutine 的异常状态，最后给主 Goroutine 发送一个死亡信号，通过`Channel`。

##### 主Goroutine

在`主Goroutine`上，从这个`Channel`读取内容，当读到内容时，就重启这个`子Goroutine`，当然`主Goroutine`需要记录`子Goroutine`的`ID`，这样也就可以针对性的启动了。

#### 代码实现

我们这里以一个工作池的场景来对上述方式进行实现。`WorkerManager`作为`主Goroutine`, `worker`作为子`Goroutine`。

```go
type WorkerManager struct {
   // 用来监控 Worker 是否已经死亡的缓冲 Channel
   workerChan chan *worker
   // 一共要监控的 worker 数量
   nWorkers int
}

// 创建一个 WorkerManager 对象
func NewWorkerManager(nworkers int) *WorkerManager {
   return &WorkerManager{
      nWorkers:nworkers,
      workerChan: make(chan *worker, nworkers),
   }
}

//启动 worker 池，并为每个 Worker 分配一个ID，让每个 Worker 进行工作
func (wm *WorkerManager)StartWorkerPool() {
   // 开启一定数量的 Worker
   for i := 0; i < wm.nWorkers; i++ {
      i := i
      wk := &worker{id: i}
      go wk.work(wm.workerChan)
   }

  // 启动保活监控
   wm.KeepLiveWorkers()
}

// 保活监控workers
func (wm *WorkerManager) KeepLiveWorkers() {
   // 如果有worker已经死亡 workChan会得到具体死亡的worker然后 打出异常，然后重启
   for wk := range wm.workerChan {
      // log the error
      fmt.Printf("Worker %d stopped with err: [%v] \n", wk.id, wk.err)
      // reset err
      wk.err = nil
      // 当前这个wk已经死亡了，需要重新启动他的业务
      go wk.work(wm.workerChan)
   }
}
```

```go
type worker struct {
   id  int
   err error
}

func (wk *worker) work(workerChan chan<- *worker) (err error) {
   // 任何Goroutine只要异常退出或者正常退出 都会调用defer 函数，所以在defer中想WorkerManager的WorkChan发送通知
   defer func() {
      // 捕获异常信息，防止panic直接退出
      if r := recover(); r != nil {
         if err, ok := r.(error); ok {
            wk.err = err
         } else {
            wk.err = fmt.Errorf("Panic happened with [%v]", r)
         }
      } else {
         wk.err = err
      }
 
     // 通知 主 Goroutine，当前子Goroutine已经死亡
      workerChan <- wk
   }()

   // do something
   fmt.Println("Start Worker...ID = ", wk.id)

   // 每个worker睡眠一定时间之后，panic退出或者 Goexit()退出
   for i := 0; i < 5; i++ {
      time.Sleep(time.Second*1)
   }

   panic("worker panic..")
   // runtime.Goexit()

   return err
}
```

#### 测试

```go
func main() {
   wm := NewWorkerManager(10)
   wm.StartWorkerPool()
}
```

可以发现无论子 Goroutine 是因为 panic() 异常退出，还是`Goexit()`退出，都会被主 Goroutine 监听到并且重启，这样就能起到保活的功能了。当然如果线程死亡？进程死亡？又该如何保证？使用 Go开发实际上是基于 Go 的调度器来开发的，进程、线程级别的死亡，会导致调度器死亡，那么全部基础框架都将会塌陷。那么就要看线程、进程如何保活了，不在 Go 开发的范畴之内了。

