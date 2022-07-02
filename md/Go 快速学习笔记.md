## Go 快速学习笔记

### Hello Go!

```go
package main // 程序的包名,表示为主包,与文件名无关,包含main函数的一定是main包

// 导入包
import (
   "fmt"
)

// main函数
func main() { // 函数的( 一定是和函数名在同一行,否则会编译失败
   fmt.Println("Hello Go!")
}
```

----------

### 变量的声明

```go
// 方法一：Go中声明一个变量,不赋值,默认值为0
var i int 
// 方法二：声明一个变量,进行赋值
var b int = 100
// 方法三：初始换变量时,不设置数据类型,Go可通过值自动匹配变量类型
var c = 100
// 方法四：(常用方法)使用 := 为变量直接赋值,不支持作为全局变量进行声明
e := 100

// 多变量声明：单行写法，类型可加可不加
var x,y = 100, 200
var x,y int = 1, 2

// 多变量声明：多行写法,类型可加可不加
var(
		a = 100
		b = 2
	)
```

-------------

### 常量 const 与 iota

```go
package main

import (
	"fmt"
)

// const 定义常量,不允许修改, iota 只能配合 const 使用,使用 var i int = iota+1 会报错
const i int = iota+1

// const + iota 常用来定义枚举类型
const (
	// 在 const() 中可以添加一个关键字 iota ,每行的 iota 会自增,第一行的 iota 默认值为 0
	BEIJING = iota*2 // BEIJING = 0, iota = 0
	SHANGHAI // SHANGHAI = 2, iota = 1
	SHENZHEN // SHENZHEN = 4, iota = 2
	GUANGZHOU = iota*3 // GUANGZHOU = 9, iota = 3
	TIANJING // TIANJING = 12, iota = 4
	CHONGQING // CHONGQING = 15, iota = 5
)

const(
	// iota 第一次使用后,按行数自增,多次显式调用只改变数据自增格式
	a,b = iota+1, iota+2 // iota = 0, a = iota+1 = 1, b = iota+2 = 2
	c,d // iota = 1, c = iota+1 = 2, d = iota+1 = 3,
	e,f = iota*2, iota*3 // iota = 2, e = iota*2 = 4 , f = iota*3 =6,
	g,h // iota = 3, g = iota*2 = 6, h = iota*3 = 9
)

func main() {
	fmt.Println(i)
    fmt.Println(BEIJING, SHANGHAI, SHENZHEN, GUANGZHOU, TIANJING, CHONGQING)
	fmt.Println(a, b, c, d, e, f, g, h)
}
```

```bash
1
0 2 4 9 12 15
1 2 2 3 4 6 6 9
```

--------------

### 函数多返回值

```go
package main

import "fmt"

// 单个返回值不对形参赋值时需要 (),函数会返回零值,没有会报错，同时返回值不能是匿名,即必须有形参
func foo1(a string, b int) int{
   fmt.Println("foo1: a =", a, "b =", b)

   c := 666
   return c
}

// 匿名返回多个返回值
func foo2(a string, b int)(int, int){
   fmt.Println("foo2: a =", a, "b =", b)

   return 1, 2
}

// 有形参名称的返回多个返回值,不对形参进行赋值
func foo3(a string, b int)(x string, y bool, z *int){
   fmt.Println("foo3: a =", a, "b =", b)

   /*
   x,y,z 属于foo3的形参,初始化默认为零值, Go中零值是变量未进行初始化时系统默认设置的值
   string 类型零值为“”, bool 类型零值为 false,数值型零值为 0
   1.var a *int
   2.var a []int
   3.var a map[string] int
   4.var a chan int
   5.var a func(string) int
   6.var a error // error是接口
   以上六种类型零值为nil
    */
   return
}

// 有形参名称的返回多个返回值
func foo4(a string, b int)(x int, y int){
   fmt.Println("foo4: a =", a, "b =", b)

   // 给有名称的返回值赋值
   x, y = 3, 4
   return
}

func main() {
   c := foo1("abc", 888)
   fmt.Println("c =", c)

   ret1, ret2 := foo2("hello", 999)
   fmt.Println("ret1 =", ret1, "ret2 =", ret2)

   ret3, ret4, ret5 := foo3("world", 111)
   fmt.Println("ret3 =", ret3, "ret4 =", ret4, "ret5 =", ret5)

   ret1, ret2 = foo4("golang", 222)
   fmt.Println("ret1 =", ret1, "ret2 =", ret2)
}
```

```bash
foo1: a = abc b = 888
c = 666
foo2: a = hello b = 999
ret1 = 1 ret2 = 2
foo3: a = world b = 111
ret3 =  ret4 = false ret5 = <nil> # 注意ret5的输出结果
foo4: a = golang b = 222
ret1 = 3 ret2 = 4
```

-------------

### import 导包与 init()

![go程序执行过程](https://raw.githubusercontent.com/tonshz/test/master/img/go%E7%A8%8B%E5%BA%8F%E6%89%A7%E8%A1%8C%E8%BF%87%E7%A8%8B.png "go程序执行过程")

![项目结构](https://raw.githubusercontent.com/tonshz/test/master/img/%E9%A1%B9%E7%9B%AE%E7%BB%93%E6%9E%84.png "项目结构")

```go
package lib

import "fmt"

// Go对外开放的函数首字母必须大写,小写只能在当前包内使用,libTest()在外部使用会报错
func LibTest(){
	fmt.Println("libTest() ...")
}

func init() {
	fmt.Println("lib.init() ...")
}

```

```go
package lib1

import "fmt"

// 当前lib1包提供的API,函数首字母大写表示为对外开放接口
func Lib1Test(){
   fmt.Println("lib1Test() ...")
}

func init() {
   fmt.Println("lib1.init() ...")
}
```

```go
package main

import (
   "awesomeProject/lib"
   "awesomeProject/lib1"
)

func main() {
   lib.LibTest()
   lib1.Lib1Test()
}
```

```bash
lib.init() ...
lib1.init() ...
libTest() ...
lib1Test() ...
```

--------------

#### import 导包

```go
package main

import (
   // _:匿名,无法使用当前包的方法,只会执行内部的 init() 方法
   _ "awesomeProject/lib"
   //"awesomeProject/lib1"

   // 为包起一个别名,使用方法:mylib.Lib1Test()
   mylib1 "awesomeProject/lib1"
   // .(不推荐使用):包中所有方法可以直接使用 Function(),不需要 lib.Function()
   //. "awesomeProject/lib1"
)

func main() {
   //lib.LibTest()
   //lib1.Lib1Test()
   mylib1.Lib1Test()
   //Lib1Test()
}
```

```bash
lib.init() ...
lib1.init() ...
lib1Test() ...
```

-------------------

### Go 指针(不常使用)

![传值](https://raw.githubusercontent.com/tonshz/test/master/img/%E4%BC%A0%E5%80%BC.png "传值")

![传指针](https://raw.githubusercontent.com/tonshz/test/master/img/%E4%BC%A0%E6%8C%87%E9%92%88.png "传指针")

```go
package main

import "fmt"

func swap(a *int, b *int){
	// a,b 值为地址, *a,*b 中为地址保存的值
	fmt.Println("a =", a, "b =", b)
	fmt.Println("*a =", *a, "*b =", *b)
	var temp int
	temp = *a // temp = main::a
	*a = *b // main::a = mian::b
	*b = temp // main::b = temp
}

func main() {
	var a, b = 10, 20
	swap(&a, &b)
	fmt.Println("a =", a, "b =", b)

	var p *int
	p = &a
	fmt.Println("&a =", &a, "p =", p)

	var pp **int // 二级指针
	pp = &p
	fmt.Println("&p =", &p, "pp =", pp)
}
```

```bash
a = 0xc0000aa058 b = 0xc0000aa070
*a = 10 *b = 20
a = 20 b = 10
&a = 0xc0000aa058 p = 0xc0000aa058
&p = 0xc0000ce020 pp = 0xc0000ce020
```

----------------------

### defer 语句

#### defer 的执行顺序

![defer 入栈出栈流程](https://raw.githubusercontent.com/tonshz/test/master/img/defer%20%E5%85%A5%E6%A0%88%E5%87%BA%E6%A0%88%E6%B5%81%E7%A8%8B.png "defer 入栈出栈流程")

```go
// defer 的执行顺序
package main

import "fmt"

func func1(){
   fmt.Println("A")
}

func func2(){
   fmt.Println("B")
}

func func3(){
   fmt.Println("C")
}

func main() {
   // defer 关键字,类似 Java 中的 finally
   defer fmt.Println("main end1")
   // defer 以压栈的方式执行,先写的 defer 先入栈,先入栈的后执行
   defer fmt.Println("main end2")
   defer func1()
   defer func2()
   defer func3()

   fmt.Println("main::hello go 1")
   fmt.Println("main::hello go 2")
}
```

```bash
main::hello go 1
main::hello go 2
C
B
A
main end2
main end1
```

#### defer 和 return 的执行顺序

```go
// defer 和 return 的执行顺序
package main

import "fmt"

func deferFunc() int{
   fmt.Println("defer func called ...")
   return 0
}

func returnFunc() int{
   fmt.Println("return func called ...")
   return 0
}

func returnAndDefer() int{
	// defer 当前函数生命周期全部结束后才执行
	defer deferFunc()
	return returnFunc()
} // 当程序逻辑执行到此处时才执行 defer 语句

func main() {
   returnAndDefer()
}
```

```bash
return func called ...
defer func called ...
```

#### defer 遇见 panic

遇到panic时，遍历本协程的defer链表，并执行defer。在执行defer过程中，遇到recover则停止panic，返回recover处继续往下执行。如果没有遇到recover，遍历完本协程的defer链表后，向stderr抛出panic信息。

![defer遇见panic](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305173223500.png "defer遇见panic")

> defer 最大的功能是 panic 后依然有效，所以defer可以保证某些资源一定会被关闭，从而避免异常出现。

##### defer遇见panic，但是并不捕获异常的情况

```go
package main

import (
    "fmt"
)

func main() {
    defer_call()

    fmt.Println("main 正常结束")
}

func defer_call() {
    defer func() { fmt.Println("defer: panic 之前1") }()
    defer func() { fmt.Println("defer: panic 之前2") }()

    panic("异常内容")  // 触发 defer 出栈

	defer func() { fmt.Println("defer: panic 之后，永远执行不到") }()
}
```

```bash
defer: panic 之前2
defer: panic 之前1
panic: 异常内容
//... 异常堆栈信息
```

##### defer遇见panic，并捕获异常

```go
package main

import (
    "fmt"
)

func main() {
    defer_call()

    fmt.Println("main 正常结束")
}

func defer_call() {

    defer func() {
        fmt.Println("defer: panic 之前1, 捕获异常")
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()

    defer func() { fmt.Println("defer: panic 之前2, 不捕获") }()

    panic("异常内容")  // 触发 defer 出栈

	defer func() { fmt.Println("defer: panic 之后, 永远执行不到") }()
}
```

```bash
defer: panic 之前2, 不捕获
defer: panic 之前1, 捕获异常
异常内容
main 正常结束
```

#### defer 包含 panic

```go
package main

import (
    "fmt"
)

func main()  {

    defer func() {
       if err := recover(); err != nil{
           fmt.Println(err)
       }else {
           fmt.Println("fatal")
       }
    }()

    defer func() {
        panic("defer panic")
    }()

    panic("panic")
}
// 输出：defer panic
```

panic 仅有最后一个可以被 revover 捕获。触发`panic("panic")`后 defer 顺序出栈执行，第一个被执行的 defer 中有`panic("defer panic")`异常语句，这个异常将会覆盖掉main中的异常`panic("panic")`，最后这个异常被第二个执行的defer捕获到。

#### defer 下的函数包含子函数

```go
package main

import "fmt"

func function(index int, value int) int {

    fmt.Println(index)

    return index
}

func main() {
    defer function(1, function(3, 0))
    defer function(2, function(4, 0))
}
// 输出结果：3 4 2 1
```

这里面有两个defer， 所以 defer 一共会压栈两次，先进栈1，后进栈2。 在压栈 function1 的时候，**需要连同函数地址、函数形参一同进栈**，为了得到 function1 的第二个参数的结果，**需要先执行function3 将第二个参数算出**，故 function3 会被第一个执行。同理压栈 function2，需要执行function4 算出 function2 第二个参数的值。然后函数结束，先出栈 fuction2、再出栈 function1。

> - defer压栈function1，压栈函数地址、形参1、形参2(调用function3) --> 打印3
> - defer压栈function2，压栈函数地址、形参1、形参2(调用function4) --> 打印4
> - defer出栈function2, 调用function2 --> 打印2
> - defer出栈function1, 调用function1--> 打印1

#### defer 面试真题

```go
package main

import "fmt"

/*
将返回值t赋值为传入的i,此时t为1
执行return语句将t赋值给t（等于啥也没做）
执行defer方法,将t + 3 = 4
函数返回 4
因为t的作用域为整个函数所以修改有效。
 */
func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3 // 4
	}()
	return t
}

/*
创建变量t并赋值为1
执行return语句,注意这里是将t赋值给返回值,
此时返回值为1（这个返回值并不是t）
执行defer方法,将t + 3 = 4
函数返回返回值1
此函数等同于
func DeferFunc2(i int) (result int) {
    t := i
    defer func() {
        t += 3
    }()
    return t
}
 */
func DeferFunc2(i int) int {
	t := i
	defer func() {
		t += 3
	}()
	return t // 1
}

/*
首先执行return将返回值t赋值为2
执行defer方法将t + 1
最后返回 3
 */
func DeferFunc3(i int) (t int) {
	defer func() {
		t += i // 3
	}()
	return 2
}

/*
初始化返回值t为零值 0
首先执行defer的第一步,赋值defer中的func入参t为0
执行defer的第二步,将defer压栈,将t赋值为1
执行return语句,将返回值t赋值为2
执行defer的第三步,出栈并执行
因为在入栈时defer执行的func的入参已经赋值了
此时它作为的是一个形式参数,所以打印为0
相对应的因为最后已经将t的值修改为2,所以再打印一个2
 */
func DeferFunc4() (t int) {
	defer func(i int) {
		fmt.Println(i) // 0
		fmt.Println(t) // 2
	}(t)
	t = 1
	return 2
}

func main() {
	fmt.Println(DeferFunc1(1))
	fmt.Println(DeferFunc2(1))
	fmt.Println(DeferFunc3(1))
	DeferFunc4()
}
// 输出结果： 4 1 3 0 2
```

-------------------

### Go 数组与 slice (动态数组)

#### Go 数组

```go
// 数组
package main

import "fmt"

func printArray(myArray [4]int){
   // 值拷贝,在该函数内改变数组值不会影响入参
   for index, value := range myArray{
      fmt.Println("index =", index, "value =", value)
   }
}

func main() {
   // 固定长度的数组
   var myArray1 [10] int
   myArray2 := [10]int{1,2,3,4}
   myArray3 := [4]int{12,22,32,42}

   //for i := 0; i < 10; i++ {
   for i := 0; i < len(myArray1); i++ {
      fmt.Println(myArray1[i])
   }

   // range 根据遍历的集合的不同返回不同的值,遍历数组和切片返回当前元素所在的下标和值
   for index, value := range myArray2{
      fmt.Println("index =", index, "value =", value)
   }

   // 查看数组的数据类型
   fmt.Printf("myArray1 types = %T\n", myArray1)
   fmt.Printf("myArray2 types = %T\n", myArray2)
   fmt.Printf("myArray3 types = %T\n", myArray3)

   // printArray(myArray2) 会报错,固定长度的数组在传参的时候,是严格匹配数组类型的
   printArray(myArray3)
}
```

```bash
0
0
0
0
0
0
0
0
0
0
index = 0 value = 1
index = 1 value = 2
index = 2 value = 3
index = 3 value = 4
index = 4 value = 0
index = 5 value = 0
index = 6 value = 0
index = 7 value = 0
index = 8 value = 0
index = 9 value = 0
myArray1 types = [10]int
myArray2 types = [10]int
myArray3 types = [4]int
index = 0 value = 12
index = 1 value = 22
index = 2 value = 32
index = 3 value = 42
```

#### slice 动态数组

```go
// slice 动态数组
package main

import "fmt"

func printArray(myArray []int){
   // 引用传递,实际传递的是指针,修改值会影响入参
   // _ 表示匿名的变量
   for _, value := range myArray{
      fmt.Println("value =", value)
   }
   myArray[0] = 100
}

func main() {
   // 动态数组, 切片 slice ,使用[]表示动态数组
   myArray := []int{1,2,3,4}
   fmt.Printf("myArray type is: %T\n", myArray)

   printArray(myArray)
   fmt.Println("------")
   for _, value := range myArray{
      fmt.Println("value =", value)
   }
}
```

```bash
myArray type is: []int
value = 1
value = 2
value = 3
value = 4
------
value = 100
value = 2
value = 3
value = 4
```

--------------------

### slice 声明

```go
package main

import "fmt"

func main() {
   // 方法一: 声明 slice1 是一个切片,并进行初始化,值为 1,2,3,4,长度 len 为 4
   slice1 := []int{1,2,3,4}

   // 方法二: 声明 slice2 是一个切片,但是未给 slice2 分配空间, len 为 0,值为 []
   var slice2 []bool
   // 在此处使用 slice2[0] = true 会报错,需要使用 make() 给 slice2 开辟空间才能赋值
   fmt.Printf("slice2: len = %d, slice2 = %v\n", len(slice2), slice2)
   // 开辟三个空间,默认值为初始化类型的零值
   slice2 = make([]bool, 3)
   slice2[0] = true

   // 方法三: 声明 slice3 是一个切片,同时为 slice3 分配空间,初始化值为对应类型的零值
   var slice3 = make([]int, 3)

   // 方法四(常用): 与方法三等价, 通关 := 推导出 slice4 是一个切片
   slice4 := make([]string, 3)
   fmt.Printf("slice4: len = %d, slice4 = %v\n", len(slice4), slice4)
   slice4[0] = "hello go"

   // %v 表示打印出任何类型的详细信息
   fmt.Printf("slice1: len = %d, slice1 = %v\n", len(slice1), slice1)
   fmt.Printf("slice2: len = %d, slice2 = %v\n", len(slice2), slice2)
   fmt.Printf("slice3: len = %d, slice3 = %v\n", len(slice3), slice3)
   fmt.Printf("slice4: len = %d, slice4 = %v\n", len(slice4), slice4)

   var slice5 []int
   // 判断一个 slice 是否为空
   if slice5 == nil{
      fmt.Println("slice5 是一个空切片")
   }else{
      fmt.Println("slice5 不是一个空切片")
   }
}
```

```bash
slice2: len = 0, slice2 = []
slice4: len = 3, slice4 = [  ]
slice1: len = 4, slice1 = [1 2 3 4]
slice2: len = 3, slice2 = [true false false]
slice3: len = 3, slice3 = [0 0 0]
slice4: len = 3, slice4 = [hello go  ]
slice5 是一个空切片
```

---------------

### slice 使用方式

```go
// slice 追加
package main

import "fmt"

func main() {
   // 创建 slice 时, make() 中后两个参数表示长度与容量,当前切片元素数量取决于长度
   var numbers = make([]int, 3, 5)
   fmt.Printf("len = %d, cap = %d, slice = %v\n", len(numbers), cap(numbers), numbers)

   // 向 numbers 切片追加一个元素 1,此时 numbers 长度变为 4,容量不变
   numbers = append(numbers, 1)
   fmt.Printf("len = %d, cap = %d, slice = %v\n", len(numbers), cap(numbers), numbers)

   // 向一个容量已满的 slice 追加元素,当长度超过容量,容量会变成 cap = 2*cap
   numbers = append(numbers, 2, 3)
   fmt.Printf("len = %d, cap = %d, slice = %v\n", len(numbers), cap(numbers), numbers)
}
```

```bash
len = 3, cap = 5, slice = [0 0 0]
len = 4, cap = 5, slice = [0 0 0 1]
len = 6, cap = 10, slice = [0 0 0 1 2 3]
```

```go
// slice 截取
package main

import "fmt"

func main() {
   // len = 3, cap = 3
   s := []int{1, 2, 3}
   // 截取 s 元素, s1: [1, 2],修改 s1 中元素会影响 s,为浅拷贝
   s1 := s[0:2]
   s1[1] = 222
   // s 与 s1 指向相同的地址
   fmt.Println("s =", s, "s1 =", s1)

   // copy() 可以将底层数组的 slice 进行拷贝赋值,为深拷贝
   s2 := make([]int, 4)
   // 将 s 中的值,依次拷贝到 s2 中,超出部分设置为类型零值, s2 与 s类型相同
   copy(s2,s)
   s2[2] = 666
   fmt.Println("s =", s,"s2 =", s2)
}
```

```bash
s = [1 222 3] s1 = [1 222]
s = [1 222 3] s2 = [1 222 666 0]
```

-----------------

### map 声明

```go
package main

import "fmt"

func main() {
   // 方法一: 声明 myMap1 是 map类型,格式: map[key]value
   var myMap1 map[string]string
   if(myMap1 == nil){
      fmt.Println("myMap1 是一个空map")
   }

   // 在使用 map 前,需要先用 make() 为 map 分配空间
   myMap1 = make(map[string]string, 10)

   myMap1["first"] = "Java"
   myMap1["second"] = "Python"
   myMap1["third"] = "Golang"
   fmt.Println("myMap1 =", myMap1)

   // 方法二: 不设置 len
   myMap2 := make(map[int]string)
   myMap2[1] = "Java"
   myMap2[2] = "Python"
   myMap2[3] = "Golang"
   fmt.Println("myMap2 =", myMap2)

   // 方法三: 直接初始化赋值 map
   myMap3 := map[string]string{
      "one": "Java",
      "two": "Python",
      "three": "Golang",
   }
   fmt.Println("myMap3 =", myMap3)
}
```

```bash
myMap1 是一个空map
myMap1 = map[first:Java second:Python third:Golang]
myMap2 = map[1:Java 2:Python 3:Golang]
myMap3 = map[one:Java three:Golang two:Python]
```

------------

### map 使用

```go
// map 增删改查
package main

import "fmt"

// 遍历 map
func printMap(cityMap map[string]string){
   // cityMap 是一个引用传递,修改会影响入参
   for key, value := range cityMap{
      fmt.Println("key =", key, "value =", value)
   }
}

func main() {
   cityMap := make(map[string]string)
   // 添加
   cityMap["China"] = "Beijing"
   cityMap["Japan"] = "Tokyo"
   cityMap["USA"] = "NewYork"
   printMap(cityMap)

   // 删除
   delete(cityMap, "China")
   // 修改
   cityMap["USA"] = "DC"
   fmt.Println("----------修改删除后------------")
   printMap(cityMap)
}
```

```bash
key = China value = Beijing
key = Japan value = Tokyo
key = USA value = NewYork
----------修改删除后------------
key = Japan value = Tokyo
key = USA value = DC
```

-----------------------

### struct 使用

```go
package main

import "fmt"

// type 表示声明一种新的数据类型, myint 是 int 的一个别名
type myint int

// 定义一个结构体, struct 默认传值
type Book struct{
   title string
   auth string
}

func changeBook(book Book){
   // 传递一个 book 的副本, 不改变入参的值
   book.auth = "lisi"
}

// struct 需要修改值时需要传指针
func changeBookPoint(book *Book){
   book.auth = "lisi"
}

func main() {
   var a myint = 10
   fmt.Printf("a = %v, type of a %T\n", a, a)

   var book1 Book
   book1.title = "Golang"
   book1.auth = "zhangsan"
   fmt.Println("book =", book1)
   changeBook(book1)
   fmt.Println("changed book =", book1)
   changeBookPoint(&book1)
   fmt.Println("point changed book =", book1)
}
```

```bash
a = 10, type of a main.myint
book = {Golang zhangsan}
changed book = {Golang zhangsan}
point changed book = {Golang lisi}
```

-------------------

### 类: 封装

```go
/*
类名、属性名、方法名 首字母大写表示对外(其他包)
可以访问，否则只能够在本包内访问
 */
package main

import "fmt"

// Go 语言中的类通过结构体绑定方法实现,类首字母大写,表示类对外开放
type Hero struct{
   // 类属性首字母大写,表示该属性对外开放(public),否则只能类的内部访问(private)
   name string
   ad int
   level int
}

// 同包内无论大小写都可以访问,小写其他包不能访问
func(this Hero) Show(){
   fmt.Println("name =", this.name)
   fmt.Println("ad =", this.ad)
   fmt.Println("level =", this.level)
}

// (this Hero) 表示当前方法绑定到当前文件下结构体 Hero
func(this Hero) GetName() string{
   return this.name
}

//func(this Hero) SetName(name string){
// // this 时调用该方法的对象的一个副本(拷贝),修改不影响入参
// this.name = name
//}

// this *Hero 将 this 换成一个指针,此时修改影响入参
func(this *Hero) SetName(name string){
   this.name = name
}

func main() {
   // 创建一个对象
   hero := Hero{name:"zhangsan", ad:100, level:1}
   hero.Show()
   hero.SetName("lisi")
   hero.Show()
}
```

```bash
name = zhangsan
ad = 100
level = 1
name = lisi
ad = 100
level = 1
```

-------------------------

### 类: 继承

```go
package main

import "fmt"

type Human struct{
   name string
   sex string
}

func(this *Human) Eat(){
   fmt.Println("Human.Eat()...")
}

func(this *Human) Walk(){
   fmt.Println("Human.Walk()...")
}

type SuperMan struct{
   // 将父类的名称写在结构体内,表示 SuperMan 类继承了 Human 类的方法
   Human
   level int
}

// 重定义父类的方法 Eat()
func(this *SuperMan) Eat(){
   fmt.Println("SuperMan.Eat()...")
}

// 子类新方法
func(this *SuperMan) Fly(){
   fmt.Println("SuperMan.Fly()...")
}

func(this *SuperMan) Print(){
   fmt.Println("name =", this.name, "sex =", this.sex, "level =", this.level)
}

func main() {
   h := Human{"zhangsan", "male"}
   h.Eat()
   h.Walk()

   // 定义一个子类对象
   // 方法一: 使用已有 Human 类对象
   superMan := SuperMan{h,  100}
   superMan.Walk()
   superMan.Eat()
   superMan.Fly()
   // 方法二: 初始化的同时赋值
   super := SuperMan{Human{"lisi", "male"},88}
   // 方法三: 先初始化再逐个赋值
   var s SuperMan
   s.name = "wangwu"
   s.sex = "male"
   s.level = 99
   superMan.Print()
   super.Print()
   s.Print()
}
```

```bash
Human.Eat()...
Human.Walk()...
Human.Walk()...
SuperMan.Eat()...
SuperMan.Fly()...
name = zhangsan sex = male level = 100
name = lisi sex = male level = 88
name = wangwu sex = male level = 99
```

-----------------------

### <span id="多态">类: 多态</span>

```go
/*
多态的基本要素:
1.有一个父类(有接口)
2.有子类(实现了父类的全部接口方法)
3.父类类型的变量(指针)指向(引用)子类的具体数据变量(父类指针指向子类对象)
 */
package main

import "fmt"

// 本质是一个指针
type AnimalIF interface {
   Sleep()
   GetColor() string // 获取动物的颜色
   GetType() string // 获取动物的种类
}

// 具体的类
type Cat struct{
   color string
}

func(this *Cat) Sleep(){
   fmt.Println("Cat is Sleep...")
}

func(this *Cat) GetColor() string{
   return this.color
}

func(this *Cat) GetType() string{
   return "Cat"
}

// 具体的类
type Dog struct{
   color string
}

func(this *Dog) Sleep(){
   fmt.Println("Dog is Sleep...")
}

func(this *Dog) GetColor() string{
   return this.color
}

func(this *Dog) GetType() string{
   return "Dog"
}

func showAnimal(animal AnimalIF){
   animal.Sleep() // 多态
   fmt.Println("color =", animal.GetColor())
   fmt.Println("type =", animal.GetType())
}

func main() {
   // 接口的数据类型,父类指针
   var animal AnimalIF
   animal = &Cat{"Green"}
   // 调用 Cat 的 Sleep() 方法,多态的现象
   animal.Sleep()

   animal = &Dog{"Yellow"}
   // 调用 Dog 的 Sleep() 方法,多态的现象
   animal.Sleep()

   fmt.Println("-----多态的直观体现----")
   // 多态的直观体现
   cat := Cat{"Red"}
   dog := Dog{"Black"}

   showAnimal(&cat)
   showAnimal(&dog)
}
```

```bash
Cat is Sleep...
Dog is Sleep...
-----多态的直观体现----
Cat is Sleep...
color = Red
type = Cat
Dog is Sleep...
color = Black
type = Dog
```

------------------

### <span id="interface{}">interface 空接口</span>

```go
package main

import "fmt"

// interface{} 是万能数据类型,类似于 Java 中的 Object
func myFunc(arg interface{}){
   fmt.Println("myFun() is called...")
   fmt.Println(arg)

   // 区分此时引用的底层数据类型: Go 为 interface{} 提供了“类型断言”的机制
   value, ok := arg.(string)
   if !ok {
      fmt.Println("arg is not string type")
   }else {
      fmt.Println("arg is string type, value =", value)
      fmt.Printf("value type is %T\n", value)
   }
}

type Book struct {
   name string
   auth string
}
func main() {
   book := Book{"Golang", "zhangsan"}
   myFunc(book)
   myFunc(100)
   myFunc("abc")
   myFunc(3.14)
}
```

```bash
myFun() is called...
{Golang zhangsan}
arg is not string type
myFun() is called...
100
arg is not string type
myFun() is called...
abc
arg is string type, value = abc
value type is string
myFun() is called...
3.14
arg is not string type
```

------------------

### Go 变量

![Go 变量](https://raw.githubusercontent.com/tonshz/test/master/img/Go%20%E5%8F%98%E9%87%8F.png "Go 变量")

```go
package main

import "fmt"

func main() {
   var a string
   // pair<statictype:string, value:"test>
   a = "test"

   // pair<statictype:string, value:"test>
   var allType interface{}
   allType = a

   s := allType.(string)
   fmt.Println(s)
}
```

```go
package main

import (
   "fmt"
   "io"
   "os"
)

func main() {
   // tty: pair<type:*os.File, value:"/dev/tty"文件描述符>
   tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
   if err != nil{
      fmt.Println("open file error", err)
      return
   }

   // r: pair<type:, value:>
   var r io.Reader
   // r: pair<type:*os.File, value:"/dev/tty"文件描述符>
   r = tty

   // w: pair<type:, value:>
   var w io.Writer
   // w: pair<type:*os.File, value:"/dev/tty"文件描述符>
   w = r.(io.Writer)

   w.Write([]byte("HELLO THIS IS A TEST!!!\n"))
}
```

```go
package main

import "fmt"

type Reader interface {
   ReadBook()
}

type Writer interface {
   WriteBook()
}

//具体类型
type Book struct {
}

func (this *Book) ReadBook() {
   fmt.Println("Read a Book")
}

func (this *Book) WriteBook() {
   fmt.Println("Write a Book")
}

func main() {
   //b: pair<type:Book, value:book{}地址>
   b := &Book{}

   //r: pair<type:, value:>
   var r Reader
   //r: pair<type:Book, value:book{}地址>
   r = b

   r.ReadBook()

   var w Writer
   //r: pair<type:Book, value:book{}地址>
    w = r.(Writer) //此处的断言为什么会成功? 因为w r 具体的type是一致,都是 Book (注意pair中的type)

   w.WriteBook()
}
```

----------------------------

### 反射

```go
package main

import (
   "fmt"
   "reflect"
)

func reflectNum(arg interface{}) {
   fmt.Println("type : ", reflect.TypeOf(arg))
   fmt.Println("value : ", reflect.ValueOf(arg))
}

func main() {
   var num float64 = 1.2345

   reflectNum(num)
}
```

```bash
type :  float64
value :  1.2345
```

```go
package main

import (
   "fmt"
   "reflect"
)

type User struct {
   Id   int
   Name string
   Age  int
}

func (this User) Call() {
   fmt.Println("user is called ..")
   fmt.Printf("%v\n", this)
}

func main() {
   user := User{1, "Aceld", 18}

   DoFiledAndMethod(user)
}

func DoFiledAndMethod(input interface{}) {
   // 获取 input 的 type
   inputType := reflect.TypeOf(input)
   // 直接输出 inputType 为 main.User
   fmt.Println("inputType is :", inputType.Name())

   // 获取 input 的 value
   inputValue := reflect.ValueOf(input)
   fmt.Println("inputValue is:", inputValue)

   /*
   通过type 获取里面的字段
   1. 获取 interface 的 reflect.Type,
   通过 Type 得到 NumField() (类一共有多少个字段),进行遍历
   2. 得到每个 field,数据类型
   3. 通过 filed 有一个 Interface() 方法得到对应的 value
    */
   for i := 0; i < inputType.NumField(); i++ {
      field := inputType.Field(i)
      value := inputValue.Field(i).Interface()

      fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
   }

   // 通过 type 获取类中方法,通过 value 调用类中方法
   for i := 0; i < inputType.NumMethod(); i++ {
      m := inputType.Method(i)
      fmt.Printf("%s: %v\n", m.Name, m.Type)
      // 调用类中方法
      inputValue.Method(i).Call(nil)
   }
}
```

```bash
inputType is : User
inputValue is: {1 Aceld 18}
Id: int = 1
Name: string = Aceld
Age: int = 18
Call: func(main.User)
user is called ..
{1 Aceld 18}
```

-------------------

### 结构体标签 Tag

```go
package main

import (
   "fmt"
   "reflect"
)

type resume struct {
   // 使用 'key:"value"' 在结构体内添加标签,多个时使用 空格 或 , 进行分割
   Name string `info:"name" doc:"我的名字"`
   Sex  string `info:"sex"`
}

// Go 中结构体标签类似于 Java 中的注解
func findTag(str interface{}) {
   // Elem() 表示获取当前结构体全部的元素
   t := reflect.TypeOf(str).Elem()

   for i := 0; i < t.NumField(); i++ {
      taginfo := t.Field(i).Tag.Get("info")
      tagdoc := t.Field(i).Tag.Get("doc")
      fmt.Println("info: ", taginfo, " doc: ", tagdoc)
   }
}

func main() {
   var re resume
   // interface{} 万能类型是指针,需要传地址,即 &re
   findTag(&re)
}
```

```bash
info:  name  doc:  我的名字
info:  sex  doc:  
```

-----------------

### 结构体标签在 JSON 中的使用

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Movie struct {
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Price  int      `json:"rmb"`
	Actors []string `json:"actors"`
}

func main() {
	movie := Movie{"喜剧之王", 2000, 10, []string{"xingye", "zhangbozhi"}}

	// 编码的过程  结构体---> json
	// json.Marshal() 可以将结构体转换成一种JSON格式
	jsonStr, err := json.Marshal(&movie)
	if err != nil {
		fmt.Println("json marshal error", err)
		return
	}
	fmt.Printf("jsonStr: %s\n", jsonStr)

	// 解码的过程 jsonstr ---> 结构体
	myMovie := Movie{}
	// json.Unmarshal() 将 JSON 转换成结构体, 
	err = json.Unmarshal(jsonStr, &myMovie)
	if err != nil {
		fmt.Println("json unmarshal error ", err)
		return
	}
	fmt.Printf("struct: %v\n", myMovie)
}
```

```bash
jsonStr: {"title":"喜剧之王","year":2000,"rmb":10,"actors":["xingye","zhangbozhi"]}
struct: {喜剧之王 2000 10 [xingye zhangbozhi]}
```

------------------

### goroutine 概念

![GMP概念](https://raw.githubusercontent.com/tonshz/test/master/img/GMP%E6%A6%82%E5%BF%B5.png "GMP概念")

![goroutine 调度器](https://raw.githubusercontent.com/tonshz/test/master/img/goroutine%20%E8%B0%83%E5%BA%A6%E5%99%A8.png "goroutine 调度器")

#### 调度器设计策略

<center><b>调度器设计策略</b></center>

|      复用线程      |        利用并行        | 抢占 | 全局G队列 |
| :----------------: | :--------------------: | :--: | :-------: |
| work stealing 机制 | GOMAXPROCS 限定P的个数 |      |           |
|   hand off 机制    |      = CPU核数/2       |      |           |

##### work stealing机制

当本线程(M2)没有可运行的协程G时,尝试从其他内核线程(M1)绑定的本地队列P获取在排队的协程G3，而不是销毁线程.

![work stealing 机制](https://raw.githubusercontent.com/tonshz/test/master/img/work%20stealing%20%E6%9C%BA%E5%88%B6.png "work stealing 机制")

##### hand off机制

当本线程因为G1进行系统调用而阻塞的时候，内核线程(M1)会释放本地队列，将P转移给其他空闲线程(M3)执行,G1与M1绑定,M1会变成睡眠/销毁状态,如果G1还需要继续执行就加到其他队列中.

![hand off 机制](https://raw.githubusercontent.com/tonshz/test/master/img/hand%20off%20%E6%9C%BA%E5%88%B6.png "hand off 机制")

#### go func() 调度流程

![go func() 调度流程](https://raw.githubusercontent.com/tonshz/test/master/img/764f7be119026cc16314e87628e4013f_1920x1080.jpeg "go func() 调度流程")

从上图可得：

1. 通过 go func()来创建一个goroutine；

2. 有两个存储G（goroutine）的队列，一个是局部调度器P的本地队列、一个是全局G队列。**新创建的G会先保存在P的本地队列中**，如果P的本地队列已经满了就会保存在全局的队列中；

3. G只能运行在M中，一个M必须持有一个P，**M与P是1：1的关系**。M会从P的本地队列弹出一个可执行状态的G来执行，如果P的本地队列为空**，首先从全局队列中获取G，全局为空则会从其他的MP组合偷取一个可执行的G来执行。本地队列运行顺序为 先从本地队列查询，全局队列查询，然后再从其他P里偷取，具体源码在 runtime 包的 proc.go中；**

4. 一个M调度G执行的过程是一个循环机制；

5. 当M执行某一个G时如果发生了syscall或者其余阻塞操作，M会阻塞，**如果当前有一些G（在P队列中的G）在执行，runtime 会把这个线程M从P中摘除/分离(detach)**，然后再创建一个新的操作系统的线程(如果有空闲的线程可用就复用空闲线程)来服务于这个P；

6. 当M系统调用(syscall)结束时候，这个G(发生阻塞的 G)会尝试获取一个空闲的P执行，并放入到这个P的本地队列。**如果获取不到P，那么这个线程M变成休眠状态， 加入到空闲线程中，然后这个G会被放入全局队列中。**

#### 调度器生命周期

![调度器生命周期](https://raw.githubusercontent.com/tonshz/test/master/img/b31027eeb493fa86654b41d46f34a98b_439x872.png "调度器生命周期")

##### M0

M0 是启动程序后的编号为0的主线程，这个M对应的实例会在全局变量 runtime.m0 中，不需要在 heap 上分配，**M0负责执行初始化操作和启动第一个G**， 在之后M0 就和其他的M一样了。

##### G0

G0 是每次启动一个M都会第一个创建的gourtine，**G0仅用于负责调度的G，G0不指向任何可执行的函数, 每个M都会有一个自己的G0**。在调度或系统调用时会使用G0的栈空间, 全局变量的G0是M0的G0。

##### 示例代码

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello world")
}
```

针对上面的代码对调度器里面的结构做一个分析，会经历上图所示的过程：

1. runtime创建最初的线程m0和goroutine g0，并把2者关联。
2. 调度器初始化：初始化m0、栈、垃圾回收，以及创建和初始化由 GOMAXPROCS 个 P 构成的 P 列表。
3. 示例代码中的main函数是`main.main`，`runtime`中也有1个main函数——`runtime.main`，代码经过编译后，`runtime.main`会调用`main.main`，程序启动时会为`runtime.main`创建goroutine，**为main goroutine**，然后把main goroutine加入到P的本地队列。
4. 启动m0，m0已经绑定了P，会从P的本地队列获取G，即获取到main goroutine。
5. G拥有栈，M根据G中的栈信息和调度信息设置运行环境
6. M运行G
7. G退出，再次回到M获取可运行的G，这样重复下去，直到`main.main`退出，`runtime.main`执行Defer和Panic处理，或调用`runtime.exit`退出程序。

调度器的生命周期几乎占满了一个Go程序的一生，`runtime.main`的goroutine执行之前都是为调度器做准备工作(流程图中黄色的步骤)，**`runtime.main`的goroutine运行，才是调度器的真正开始，直到`runtime.main`结束而结束。**

#### 可视化 GMP 编程

##### 1. go tool trace

```go
package main

import (
    "os"
    "fmt"
    "runtime/trace"
)

func main() {

    // 创建trace文件
    f, err := os.Create("trace.out")
    if err != nil {
        panic(err)
    }

    defer f.Close()

    //启动trace goroutine
    err = trace.Start(f)
    if err != nil {
        panic(err)
    }
    defer trace.Stop()

    //main
    fmt.Println("Hello World")
}
```

运行程序会得到一个 `trace.out`文件，使用命令 `go tool trace trace.out `可以获得可视化界面的网页。

```bash
2022/03/06 17:08:38 Parsing trace...
2022/03/06 17:08:38 Splitting trace...
2022/03/06 17:08:38 Opening browser. Trace viewer is listening on http://127.0.0.1:5942
```

使用浏览器打开网址，点击view trace 就可以看到可视化的调度流程了。

![img](https://raw.githubusercontent.com/tonshz/test/master/img/eee828bc698d074e439f3e6929be74ef_2724x546.png)

![](https://raw.githubusercontent.com/tonshz/test/master/img/25ede16ec870076f211f8924c2c2bf6f_492x556.png)

点击Goroutines那一行可视化的数据条，可以会看到 G 的一些详细信息。一共有两个G在程序中，一个是特殊的G0，是每个M必须有的一个初始化的G，这个我们不必讨论。其中G1应该就是main goroutine(执行main函数的协程)，在一段时间内处于可运行和运行的状态。

![](https://raw.githubusercontent.com/tonshz/test/master/img/87e3994fbda4e883a8c51bea20dba91a_1168x372.png)

点击Threads那一行可视化的数据条，可以会看到M一些详细信息。 一共有两个M在程序中，一个是特殊的M0，用于初始化使用。

![](https://raw.githubusercontent.com/tonshz/test/master/img/884f6aa775a4596d3c8d4f9451d55e9b_1146x248.png)

G1中调用了`main.main`，创建了`trace goroutine g18`。G1运行在P1上，G18运行在P0上。这里有两个P，因为一个P必须绑定一个M才能调度G。

![](https://raw.githubusercontent.com/tonshz/test/master/img/5e07a8515a023fcd10fbd8cf328b5d64_2736x224.png)

再来看看上方的M信息

![](https://raw.githubusercontent.com/tonshz/test/master/img/6e145c4c49656b77f9f2006733b88859_2142x466.png)

可以发现G18在P0上被运行的时候，确实在Threads行多了一个M的数据，点击查看。多的M2是P0为了执行G18而动态创建的M2。

![](https://raw.githubusercontent.com/tonshz/test/master/img/19602c163ecfe63706e10cdc90e43794_1086x298.png)

##### 2. Debug trace

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    for i := 0; i < 5; i++ {
        time.Sleep(time.Second)
        fmt.Println("Hello World")
    }
}
```

使用 `GODEBUG=schedtrace=1000 ./trace2 `通过Debug方式运行。

```bash
$ GODEBUG=schedtrace=1000 ./trace2 
SCHED 0ms: gomaxprocs=2 idleprocs=0 threads=4 spinningthreads=1 idlethreads=1 runqueue=0 [0 0]
Hello World
SCHED 1003ms: gomaxprocs=2 idleprocs=2 threads=4 spinningthreads=0 idlethreads=2 runqueue=0 [0 0]
Hello World
SCHED 2014ms: gomaxprocs=2 idleprocs=2 threads=4 spinningthreads=0 idlethreads=2 runqueue=0 [0 0]
Hello World
SCHED 3015ms: gomaxprocs=2 idleprocs=2 threads=4 spinningthreads=0 idlethreads=2 runqueue=0 [0 0]
Hello World
SCHED 4023ms: gomaxprocs=2 idleprocs=2 threads=4 spinningthreads=0 idlethreads=2 runqueue=0 [0 0]
Hello World
```

- `SCHED`：调试信息输出标志字符串，代表本行是goroutine调度器的输出；
- `0ms`：即从程序启动到输出这行日志的时间；
- `gomaxprocs`: P的数量，本例有2个P, 因为默认的P的属性是和cpu核心数量默认一致，当然也可以通过GOMAXPROCS来设置；
- `idleprocs`: 处于idle状态的P的数量；通过gomaxprocs和idleprocs的差值，我们就可知道执行go代码的P的数量；
- t`hreads: os threads/M`的数量，包含scheduler使用的m数量，加上runtime自用的类似sysmon这样的thread的数量；
- `spinningthreads`: 处于自旋状态的os thread数量；
- `idlethread`: 处于idle状态的os thread的数量；
- `runqueue=0`： Scheduler全局队列中G的数量；
- `[0 0]`: 分别为2个P的local queue中的G的数量。

-----------------------

### 创建 goroutine

```go
package main

import (
   "fmt"
   "time"
)

//子 goroutine
func newTask() {
   i := 0
   for {
      i++
      fmt.Printf("new Goroutine : i = %d\n", i)
      time.Sleep(1 * time.Second)
   }
}

//主 goroutine, main 退出,其他协程也会结束
func main() {
   //创建一个go程(协程) 去执行 newTask() 流程
   go newTask()
   fmt.Println("main goroutine exit")
   
   //i := 0
   //for {
   // i++
   // fmt.Printf("main goroutine: i = %d\n", i)
   // time.Sleep(1 * time.Second)
   //}
}
```

```bash
main goroutine exit
new Goroutine : i = 1 # 会有两种不同的输出,有时候会没有第二行输出,协程创建速度慢于main()方法结束速度
```

```bash
// goroutine 匿名函数
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	//用 go 创建承载一个形参为空，返回值为空的一个函数
	go func() {
		defer fmt.Println("A.defer")
		// 匿名函数
		func() {
			defer fmt.Println("B.defer")
			// 在子函数中退出当前 goroutine
			runtime.Goexit()
			fmt.Println("B")
		}() // 注意此处的 () 不能省,表示调用匿名函数,可进行传参
		fmt.Println("A")
	}()

	// 不能通过 flag := go func(){...} 获取执行结果,需要使用 channel
	go func(a int, b int) bool {
		fmt.Println("a = ", a, ", b = ", b)
		return true
	}(10, 20) // ()传参

	//死循环
	for {
		time.Sleep(1 * time.Second)
	}
}
```

```bash
B.defer
A.defer
a =  10 , b =  20
```

-----------------------

### channel 定义与使用

```go
package main

import "fmt"

func main() {
	// 定义一个channel: make(chan Type, capacity),未设置 cap 值默认为0,为无缓冲 channel
	c := make(chan int) // cap = 0

	// 主线程和 go程 会同步 channel ,无论哪一方先执行到 channel 处都会等待另一个执行到 channel 处
	go func() {
		defer fmt.Println("goroutine结束")

		fmt.Println("goroutine 正在运行...")

		// channel <- value 发送 value 到 channel
		c <- 666 // 将666写到 c 中
	}()

	// <- channel 接受并将其丢弃,从 channel 中读取数据,但并未使用
	// x := <-channel 从 channel 中接受数据,并赋值给 x, <-与 cahnnel 中间不能加空格
	num := <-c // 从 c 中接受数据，并赋值给 num
	// x, ok := <- channel 功能同上, 同时检查通道是否已关闭或是否为空

	fmt.Println("num = ", num)
	fmt.Println("main goroutine 结束...")
}
```

```bash
goroutine 正在运行...
goroutine结束
num = 666
main goroutine 结束...
```

-----------------

### channel 同步问题

#### 无缓冲 channel

![无缓冲 channel](https://raw.githubusercontent.com/tonshz/test/master/img/%E6%97%A0%E7%BC%93%E5%86%B2%20channel.png "无缓冲 channel")

在第 1 步，两个 goroutine 都到达通道，但都没有开始执行发送或者接收。
在第 2 步，左侧的 goroutine 将它的手伸进了通道，这模拟了向通道发送数据的行为。这时，这个 goroutine 会在通道中被锁住，直到交换完成。
在第 3 步，右侧的 goroutine 将它的手放入通道，这模拟了从通道里接收数据。这个 goroutine 一样也会在通道中被锁住，直到交换完成。
在第 4 步和第 5 步，进行交换，并最终，在第 6 步，两个 goroutine 都将它们的手从通道里里拿出来，这模拟了被锁住的 goroutine 得到释放。两个 goroutine 现在都可以去做其他事情了了。

#### 有缓冲 channel

![有缓冲 channel](https://raw.githubusercontent.com/tonshz/test/master/img/%E6%9C%89%E7%BC%93%E5%86%B2%20channel.png "有缓冲 channel")

在第 1 步，右侧的 goroutine 正在从通道接收一个值。
在第 2 步，右侧的这个 goroutine 独立完成了接收值的动作，而左侧的 goroutine 正在发送一个新值到通道里。
在第 3 步，左侧的 goroutine 还在向通道发送新值，而右侧的 goroutine 正在从通道接收另外一个值。这个步骤里的两个操作既不是同步的，也不会互相阻塞。
最后，在第 4 步，所有的发送和接收都完成，而通道里还有几个值，也有一些空间可以存更多的值。
有缓冲通道特点：当channel已经满，再向里面写数据，就会阻塞；当channel为空，从里面取数据也会阻塞

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 3) // 带有缓冲的 channel

	fmt.Println("len(c) = ", len(c), ", cap(c)", cap(c))

	go func() {
		defer fmt.Println("子go程结束")

		// 此处 i 最大值大于 channel 的 cap,结果输出时会存在多种情况,小于等于时不会
		for i := 0; i < 4; i++ {
			c <- i
			fmt.Println("子go程正在运行, 发送的元素=", i, " len(c)=", len(c), ", cap(c)=", cap(c))
		}
	}()

	time.Sleep(2 * time.Second)

	// 主线程中的 i 最大值若大于 go 程的 channel 会报错
	for i := 0; i < 4; i++ {
		num := <-c // 从 c 中接收数据，并赋值给 num
		fmt.Println("num = ", num)
	}

	fmt.Println("main 结束")
}
```

```bash
# 输出结果存在多种情况,打印数据也会抢占资源
len(c) =  0 , cap(c) 3
子go程正在运行, 发送的元素= 0  len(c)= 1 , cap(c)= 3
子go程正在运行, 发送的元素= 1  len(c)= 2 , cap(c)= 3
子go程正在运行, 发送的元素= 2  len(c)= 3 , cap(c)= 3
num =  0
num =  1
num =  2
num =  3
main 结束 # 最常见

===================================================

len(c) =  0 , cap(c) 3
子go程正在运行, 发送的元素= 0  len(c)= 1 , cap(c)= 3
子go程正在运行, 发送的元素= 1  len(c)= 2 , cap(c)= 3
子go程正在运行, 发送的元素= 2  len(c)= 3 , cap(c)= 3
num =  0
子go程正在运行, 发送的元素= 3  len(c)= 3 , cap(c)= 3
子go程结束
num =  1
num =  2
num =  3
main 结束 # 预期输出
```

-------------------

### 关闭 channel

```go
package main

import "fmt"

func main() {
	// 创建一个没有缓冲的 channel,使用 var c chan Type 时,c 为 nil channel,无论收发都会阻塞
	c := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			c <- i
			// close(c) 此处关闭 channel 会报错: panic: send on closed channel
		}
		// close 可以关闭一个 channel ,只关闭发送端,接收端仍可接收数据
		close(c) // 若将此行注释后会报错: fatal error: all goroutines are asleep - deadlock!
	}()

	for {
		// ok如果为 true 表示 channel 没有关闭，如果为 false 表示 channel 已经关闭
		if data, ok := <-c; ok {
			fmt.Println(data)
		} else {
			break
		}
	}
	fmt.Println("Main Finished..")
}
```

```bash
0
1
2
3
4
Main Finished..
```

channel 不像文件一样需要经常去关闭，只有当没有发送数据或想显式的结束range循环等情况，才去关闭 channel；

关闭 channel 后，无法向 channel 再发送数据(引发 panic 错误后导致接收立即返回零值)；

关闭 channel 后，可以继续从 channel 接收数据；

对于 nil channel (不使用 make() 进行初始化赋值)，无论收发都会被阻塞。

------------------

### channel 与 range

```go
package main

import "fmt"

func main() {
   c := make(chan int)

   go func() {
      for i := 0; i < 5; i++ {
         c <- i
      }

      // close 可以关闭一个 channel
      close(c)
   }()

   // 可以使用 range 来迭代不断操作 channel
   for data := range c {
      fmt.Println(data)
   }

   fmt.Println("Main Finished..")
}
```

-----------------------

### channel 与 select

![channel 与 select](https://raw.githubusercontent.com/tonshz/test/master/img/channel%20%E4%B8%8E%20select.png "channel 与 select")

```go
// 单流程下一个 go 只能监控一个 channel 的状态, select 可以完成监控多个 channel 的状态
package main

import "fmt"

func fibonacci(c, quit chan int) {
   x, y := 1, 1
   // select 上面一层都会嵌套一层 for,来进行循环判断
   for {
      // select 具备监控多路 channel 状态的功能
      select {
      case c <- x:
         // 如果 c 可写，则该 case 就会进来
         x, y = y, x+y
      case <-quit: // 如果 quit 可读数据,则进入该 case 处理语句
         fmt.Println("quit")
         return
      }
   }
}

func main() {
   c := make(chan int)
   quit := make(chan int)

   //sub go
   go func() {
      for i := 0; i < 10; i++ {
         fmt.Println(<-c)
      }
      quit <- 0
   }()

   //main go
   fibonacci(c, quit)
}
```

```bash
1
1
2
3
5
8
13
21
34
55
quit
```

-----------------------------

### GOPATH 工作模式(将被淘汰)

![Go modules](https://raw.githubusercontent.com/tonshz/test/master/img/Go%20modules.png "Go modules")

![GOPATH 的弊端](https://raw.githubusercontent.com/tonshz/test/master/img/GOPATH%20%E7%9A%84%E5%BC%8A%E7%AB%AF.png "GOPATH 的弊端")

-----------------------

### Go Modules 模式

#### go mod 命令(go1.11版本以上)

|       命令       |               作用               |
| :--------------: | :------------------------------: |
| **go mod init**  |         生成 go.mod 文件         |
| go mod download  | 下载 go.mod 文件中指明的所有依赖 |
| **go mode tidy** |          整理现有的依赖          |
|   go mod graph   |        查看现有的依赖结构        |
|   go mod edit    |         编辑 go.mod 文件         |
|  go mod vendor   |  导出项目所有的依赖到vendor目录  |
|  go mod verify   |     校验一个模块是否被篡改过     |
|    go mod why    |     查看为什么需要依赖某模块     |

#### go mod 环境变量

通过 `go env` 命令进行查看。

![go mod环境变量](https://raw.githubusercontent.com/tonshz/test/master/img/go%20mod%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F.png "go mod环境变量")

#### 开启 Go Modules

```bash
go env -w GO111MODULE=on
export GO111MODULE=on // 直接设置系统环境变量
```

--------------------------

### Go Modules 初始化项目

![Go Modules 初始化项目](https://raw.githubusercontent.com/tonshz/test/master/img/Go%20Modules%20%E5%88%9D%E5%A7%8B%E5%8C%96%E9%A1%B9%E7%9B%AE.png "Go Modules 初始化项目")

-----------------------------

### 修改项目模块的版本依赖关系

在命令行输入`go mod edit -replace=zinx@v0.0.0-20200306023939bc416543ae24=zinx@v0.0.0-20200221135252-8a8954e75100`,实际是进行导包依赖重定向的操作， go mod 文件会被修改:

```bash
module github.com/aceld/modules_test
go 1.14
require github.com/aceld/zinx v0.0.0-20200306023939-bc416543ae24 // indirect
replace zinx v0.0.0-20200306023939-bc416543ae24 => zinx v0.0.0-20200221135252-8a8954e75100
```

----------------

### Go 添加第三方包

+ 安装git
+ 在 Goland 中设置 github 账号
+ 国内使用 go get 需要设置代理。 `eg: go env -w GOPROXY=https://goproxy.cn,direct`
+ 不存在的包通过 go get 包名 获取。`eg: go get github.com/gin-gonic/gin`
+ go get -u 包名 // 拉取最新的代码包的最新版本，下载并安装。`eg: go get -u github.com/gin-gonic/gin`
+ 设置`go.mod`文件，执行`go mod download github.com/gin-gonic/gin`命令

-----------------------------

### Golang 中 make 与 new 的区别

指针类型声明空间后不分配空间会报错，`panic: runtime error: invalid memory address or nil pointer dereferenc`。

#### new

```go
// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *Type
```

new() 只接受一个参数，这个参数是一个类型，分配好内存后，返回一个指向该类型内存地址的指针。同时请注意它会把分配的内存置为零，也就是类型的零值。

```go
func main() {
  
   var i *int
   i=new(int)
   fmt.Println(*i) // 此处输出 0 
  
}
```

此处 new() 作用不明显，下列示例中的user类型中的lock字段不需要进行初始化，可以直接使用，不会有无效内存引用异常，因为它被设置为零值了。new  返回的是类型的指针，指向分配类型的内存地址。

```go
package main

import (
    "fmt"
    "sync"
)

type user struct {
    lock sync.Mutex
    name string
    age int
}

func main() {

    u := new(user) // 默认给u分配到内存全部为0

    u.lock.Lock()  // 可以直接使用,因为lock为0,是开锁状态
    u.name = "张三"
    u.lock.Unlock()

    fmt.Println(u)
}
// 输出： &{{0 0} 张三 0}
```

#### make

make 只适用于 chan、map、slice 的内存创建，返回的类型为这三个类型本身，而不是它们的指针类型，因为这三种类型本身就是引用类型，没必要返回指针。同时，因为这三种类型是引用类型，所以必须得初始化，但不是置为零值，与 new 不同。从函数申明中可以看到，返回的还是该类型。

```go
func make(t Type, size ...IntegerType) Type
```

#### make 与 new 的异同

两个都是使用对空间进行分配，但 make 只用于 slice、map、channel 的初始化，这三种类型必须使用 make 进行初始化，然后才可以进行操作；new 用于类型内存分配，默认初始化值为零值，不常用。

>在实际的编码中，new 是不常用的，通常采用短语句声明与结构体的字面量来实现类型效果：
>
>```go
>i : =0
>u := user{}
>```

-------------------------

### 面向对象的编程思维理解interface

#### inerface 接口

接口的最大意义是实现多态的思想，可以根据 interface 类型来设计 API接口。interface 是 Go 语言的基础特性之一。可以理解为一种类型的规范或者约定。与 Java 不同，不需要显示说明实现了某个接口，它没有继承、子类、 "implements" 的关键字，而是通过约定的形式，隐式的实现 interface 中的方法即可。因此，Golang 中的 interface 让编码更灵活、易扩展。

> 理解 Go 语言中的interface ，只需记住以下三点即可：
>
> 1. interface 是方法声明的集合，使用 type 声明。如 `type myInterface interface{}`
> 2. **任何类型的对象实现了在interface 接口中声明的全部方法**，则表明该类型实现了该接口，具体实例参见[类: 多态](#多态)
> 3. interface 可以作为一种数据类型，实现了该接口的任何对象都可以给对应的接口类型变量赋值，即[interface 空接口](#interface{})
>
> 同时需要注意一下两点：
>
> 1. interface 可以被任意对象实现，一个类型/对象也可以实现多个 interface
> 2. Go 中方法不允许重载，如 `eat(), eat(s string)` 不能同时存在，会报错`'eat' redeclared in this package`

#### 面向对象中的开闭原则

> 开闭原则定义：一个软件实体如类、模块和函数应该对扩展开放，对修改关闭。

##### 平铺式的模块设计

interface 数据类型存在的意义是为了满足一些面向对象的编程思想，如开闭原则、依赖倒转原则。

```go
package main

import "fmt"

//我们要写一个类,Banker银行业务员
type Banker struct {
}

//存款业务
func (this *Banker) Save() {
	fmt.Println( "进行了 存款业务...")
}

//转账业务
func (this *Banker) Transfer() {
	fmt.Println( "进行了 转账业务...")
}

//支付业务
func (this *Banker) Pay() {
	fmt.Println( "进行了 支付业务...")
}

func main() {
	banker := &Banker{}

	banker.Save()
	banker.Transfer()
	banker.Pay()
}
```

上例代码表示一个银行业务员可能有很多业务，随着现实业务变复杂，会导致该模块变得越来越臃肿。为 Banker 添加新的业务的时候，需要修改原有的Banker代码，当 Banker 模块的功能会越来越多，出现问题的几率也越来越大。这样会导致模块的耦合度变高，Banker 的职责也不够单一，代码的维护成本也会越来越复杂。

![Bancker 结构体](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305170003767.png "Bancker 结构体")

##### 开闭原则的设计

此时可以通过 interface 抽象出一个 Banker 模块，根据这个抽象模块去实现方法。给系统添加功能，不修改原有代码，而是通过增加新代码，这就是开闭原则的核心思想。

![抽象Banker模块](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305170523508.png "抽象Banker模块")

```go
// 实现架构层(基于抽象层进行业务封装-针对 interface 接口进行封装)
func BankerBusiness(banker AbstractBanker) {
	// 通过接口来向下调用，(多态现象)
	banker.DoBusi()
}
```

```go
//存款的业务员
type SaveBanker struct {
	//AbstractBanker
}

func (sb *SaveBanker) DoBusi() {
	fmt.Println("进行了存款")
}
```

```go
//转账的业务员
type TransferBanker struct {
	//AbstractBanker
}

func (tb *TransferBanker) DoBusi() {
	fmt.Println("进行了转账")
}
```

```go
//支付的业务员
type PayBanker struct {
	//AbstractBanker
}

func (pb *PayBanker) DoBusi() {
	fmt.Println("进行了支付")
}
```

```go
// 在 main() 中实现业务调用
func main() {
	//进行存款
	BankerBusiness(&SaveBanker{})

	//进行存款
	BankerBusiness(&TransferBanker{})

	//进行存款
	BankerBusiness(&PayBanker{})
}
```

#### 面向对象中的依赖倒转原则

> 依赖倒转原则定义: 程序要依赖于抽象接口，不要依赖于具体实现。即抽象不应该依赖于细节，细节应该依赖于抽象。

##### 耦合度极高的模块关系设计

![高耦合模块关系设计](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305171516320.png "高耦合模块关系设计")

##### 面向抽象层的依赖倒转

![抽象层的依赖倒转](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305171647798.png "抽象层的依赖倒转")

##### 依赖倒转小练习

```go
/*
	模拟组装2台电脑
    --- 抽象层 ---
	有显卡Card  方法display
	有内存Memory 方法storage
    有处理器CPU   方法calculate

    --- 实现层层 ---
	有 Intel因特尔公司 、产品有(显卡、内存、CPU)
	有 Kingston 公司， 产品有(内存3)
	有 NVIDIA 公司， 产品有(显卡)

	--- 逻辑层 ---
	1. 组装一台Intel系列的电脑，并运行
    2. 组装一台 Intel CPU  Kingston内存 NVIDIA显卡的电脑，并运行
*/
package main

import "fmt"

//------  抽象层 -----
type Card interface{
	Display()
}

type Memory interface {
	Storage()
}

type CPU interface {
	Calculate()
}

type Computer struct {
	cpu CPU
	mem Memory
	card Card
}

func NewComputer(cpu CPU, mem Memory, card Card) *Computer{
	return &Computer{
		cpu:cpu,
		mem:mem,
		card:card,
	}
}

func (this *Computer) DoWork() {
	this.cpu.Calculate()
	this.mem.Storage()
	this.card.Display()
}

//------  实现层 -----
//intel
type IntelCPU struct {
	CPU	
}

func (this *IntelCPU) Calculate() {
	fmt.Println("Intel CPU 开始计算了...")
}

type IntelMemory struct {
	Memory
}

func (this *IntelMemory) Storage() {
	fmt.Println("Intel Memory 开始存储了...")
}

type IntelCard struct {
	Card
}

func (this *IntelCard) Display() {
	fmt.Println("Intel Card 开始显示了...")
}

//kingston
type KingstonMemory struct {
	Memory
}

func (this *KingstonMemory) Storage() {
	fmt.Println("Kingston memory storage...")
}

//nvidia
type NvidiaCard struct {
	Card
}

func (this *NvidiaCard) Display() {
	fmt.Println("Nvidia card display...")
}



//------  业务逻辑层 -----
func main() {
	//intel系列的电脑
	com1 := NewComputer(&IntelCPU{}, &IntelMemory{}, &IntelCard{})
	com1.DoWork()

	//杂牌子
	com2 := NewComputer(&IntelCPU{}, &KingstonMemory{}, &NvidiaCard{})
	com2.DoWork()
}
```



--------------------------
参考文章: [Golang 修养之路—— 刘丹冰Aceld](https://www.kancloud.cn/aceld/golang/1858955)





