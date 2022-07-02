## Go调度器调度场景过程全解析

### 场景1

P拥有G1，M1获取P后开始运行G1，G1使用`go func()`创建了G2，**为了局部性G2优先加入到P1的本地队列**。
![场景1](https://raw.githubusercontent.com/tonshz/test/master/img/2debce43683adca1acb5ca5210057232_1074x900.png "场景1")

------

### 场景2

G1运行完成后(函数：`goexit`)，M上运行的goroutine切换为G0，**G0负责调度时协程的切换**（函数：`schedule`）。从P的本地队列取G2，从G0切换到G2，并开始运行G2(函数：`execute`)。**实现了线程M1的复用。**
![场景2](https://raw.githubusercontent.com/tonshz/test/master/img/93658da22081d52ed1caf32f42145e5a_1624x984.png "场景2")

------

### 场景3

假设每个P的本地队列只能存3个G。G2创建了6个G，前3个G（G3, G4, G5）已经加入p1的本地队列，p1本地队列满了。
![场景3](https://raw.githubusercontent.com/tonshz/test/master/img/6415bfab3595fc22090595acc7c1b4b1_1104x1030.png "场景3")

------

### 场景4

G2在创建G7的时候，发现P1的本地队列已满，需要执行**负载均衡**，即将P1中本地队列中前一半的G，还有新创建G**转移**到全局队列。实现中并不一定是新的G，如果G是G2之后就执行的，会被保存在本地队列，利用某个老的G（即在队列前面的G）替换新G加入全局队列。这些G被转移到全局队列时，会被打乱顺序。此时G3,G4,G7被转移到全局队列。
![场景4](https://raw.githubusercontent.com/tonshz/test/master/img/d12776bfd5cd10f8c1979c61d467499c_1120x1068.png "场景4")

------

### 场景5

G2创建G8时，P1的本地队列未满，**所以G8会被加入到P1的本地队列**。 G8加入到P1本地队列的原因是因为P1此时在与M1绑定，而G2此时是M1在执行。**所以G2创建的新的G会优先放置到自己的M绑定的P上。**

![场景5](https://raw.githubusercontent.com/tonshz/test/master/img/7a01ac7a3a4fd14493224827409f77f8_1036x1048.png "场景5")

------

### 场景6

规定：**在创建G时，运行的G会尝试唤醒其他空闲的P和M组合去执行**。假定G2唤醒了M2，M2绑定了P2，并运行G0，但P2本地队列没有G，M2此时为自旋线程**（没有G但为运行状态的线程，不断寻找G）**。

![场景6](https://raw.githubusercontent.com/tonshz/test/master/img/606acb1e4bf6b85c352b213744771601_1976x1098.png "场景6")

------

### 场景7

M2尝试从全局队列(简称“GQ”)取一批G放到P2的本地队列（函数：`findrunnable()`）。M2从全局队列取的G数量符合下面的公式：

```
n =  min(len(GQ) / GOMAXPROCS +  1,  cap(LQ) / 2 )
```

相关源码参考:

```
// 从全局队列中偷取，调用时必须锁住调度器
func globrunqget(_p_ *p, max int32) *g {
	// 如果全局队列中没有 g 直接返回
	if sched.runqsize == 0 {
		return nil
	}

	// per-P 的部分，如果只有一个 P 的全部取，即 gomaxprocs = 1
	n := sched.runqsize/gomaxprocs + 1
	if n > sched.runqsize {
		n = sched.runqsize
	}

	// 不能超过取的最大个数
	if max > 0 && n > max {
		n = max
	}

	// 计算能不能在本地队列中放下 n 个
	if n > int32(len(_p_.runq))/2 {
		n = int32(len(_p_.runq)) / 2
	}

	// 修改本地队列的剩余空间
	sched.runqsize -= n
	// 拿到全局队列队头 g
	gp := sched.runq.pop()
	// 计数
	n--

	// 继续取剩下的 n-1 个全局队列放入本地队列
	for ; n > 0; n-- {
		gp1 := sched.runq.pop()
		runqput(_p_, gp1, false)
	}
	return gp
}
```

至少从全局队列取1个g，但每次不要从全局队列移动太多的g到p本地队列，给其他p留点。这是**从全局队列到P本地队列的负载均衡**。

![场景7](https://raw.githubusercontent.com/tonshz/test/master/img/d646d4d213d7b603f211cd74ba0dd391_1920x1080.jpeg "场景7")

假定场景中一共有4个P（GOMAXPROCS设置为4，最多允许4个P来供M使用）。所以此时M2只从能从全局队列取1个G（即G3）移动P2本地队列，然后完成从G0到G3的切换，运行G3。

------

### 场景8

假设G2一直在M1上运行，经过2轮后，M2已经把G7、G4从全局队列获取到了P2的本地队列并完成运行，全局队列和P2的本地队列都空了,如场景8图的左半部分。**全局队列已经没有G，那M2就要执行work stealing(偷取)，从其他有G的P哪里偷取一半G过来，放到自己的P本地队列**。P2从P1的本地队列尾部取一半的G，本例中一半则只有1个G8，放到P2的本地队列并执行。

![场景8](https://raw.githubusercontent.com/tonshz/test/master/img/8abf6b47b0871011b6f18f55365c1774_1632x1130.png "场景8")

------

### 场景9

G1本地队列G5、G6已经被其他M偷走并运行完成，当前M1和M2分别在运行G2和G8，M3和M4没有goroutine可以运行，M3和M4处于**自旋状态**，它们不断寻找goroutine。为什么要让m3和m4自旋？自旋本质是在运行，线程在运行却没有执行G，就变成了浪费CPU。为什么不销毁线程来节约CPU资源？因为创建和销毁CPU也会浪费时间，我们**希望当有新goroutine创建时，立刻能有M运行它**，如果销毁再新建就增加了时延，降低了效率。当然也考虑了过多的自旋线程是浪费CPU，所以系统中最多有`GOMAXPROCS`个自旋的线程(当前例子中的`GOMAXPROCS`=4，所以一共4个P)，空闲的M线程会让他们休眠。

![场景9](https://raw.githubusercontent.com/tonshz/test/master/img/10b49d04c42c3d688986ff41005ee63b_1682x1084.png "场景9")

------

### 场景10

假定当前除了M3和M4为自旋线程，还有M5和M6为空闲的线程(没有得到P的绑定，注意此时设定最多只能存在4个P, 大部分都是M在抢占需要运行的P)，G8创建了G9，G8进行了**阻塞的系统调用**，M2和P2立即解绑，P2会执行以下判断：如果P2本地队列有G、全局队列有G或有空闲的M，P2都会立马唤醒1个M和它绑定，否则P2则会加入到空闲P列表，等待M来获取可用的P。本场景中，P2本地队列有G9，可以和其他空闲的线程M5绑定。
![场景10](https://raw.githubusercontent.com/tonshz/test/master/img/549aa458fab4c3e7cac5086d9326c1d2_2642x1494.png "场景10")

------

### 场景11

G8创建了G9，假如G8进行了**非阻塞系统调用**。M2和P2会解绑，但M2会记住P2，然后G8和M2进入**系统调用**状态。当G8和M2退出系统调用时，会尝试获取P2，如果无法获取，则获取空闲的P，如果依然没有，**G8会被记为可运行状态，并加入到全局队列，M2因为没有P的绑定而变成休眠状态**，长时间休眠等待GC回收销毁。

![场景11](https://raw.githubusercontent.com/tonshz/test/master/img/246c03cdf1eb8307b70865d0debaa1f0_2678x1466.png "场景11")

### 总结

Go调度器很轻量也很简单，足以撑起 goroutine 的调度工作，并且让 Go 具有了原生（强大）并发的能力。**Go 调度本质是把大量的 goroutine 分配到少量线程上去执行，并利用多核并行，实现更强大的并发。**