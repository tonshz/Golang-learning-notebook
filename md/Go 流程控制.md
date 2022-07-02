# Go 流程控制

[toc]

## 条件语句 if

条件语句需要开发者通过指定一个或多个条件，并通过测试条件是否为 true 来决定是否执行指定语句，并在条件为 false 的情况在执行另外的语句。

### if 语句

+ 可省略条件表达式括号
+ 支持初始化语句，可定义代码块局部变量
+ 代码块左括号必须在条件表达式尾部
+ 不支持三元操作符（三目运算符）`a>b ? a:b`

```go
if 布尔表达式 {
    布尔表达式为 true 时执行内容...
}
```

### if...else 语句

```go
if 布尔表达式 {
    布尔表达式为 true 时执行内容...
} else {
    布尔表达式为 false 时执行内容...
}
```

### if 嵌套语句

```go
if 布尔表达式 {
    布尔表达式为 true 时执行内容...
    if 布尔表达式 {
	    布尔表达式为 true 时执行内容...
	}
}
```

## 条件语句 switch

### 语法

switch 语句用于基于不同条件执行不同动作，每一个 case 分支都是唯一的，从上直下逐一测试，直到匹配为止。
Golang switch 分支表达式可以是任意类型，不限于常量。**可省略 break，默认自动终止。**

**Go 里面 switch 默认相当于每个case最后带有break，匹配成功后不会自动向下执行其他case，而是跳出整个switch, 但是可以使用 `fallthrough` 强制执行后面的case代码，`fallthrough`不会判断下一条case的expr结果是否为true。**

```go
switch var1 {
    case val1:
        ...
    case val2:
        ...
    default:
        ...
}
```

变量 var1 可以是任何类型，而 val1 和 val2 则可以是同类型的任意值。类型不被局限于常量或整数，**但必须是相同的类型或者最终结果为相同类型的表达式。**
您可以同时测试多个可能符合条件的值，使用逗号分割它们，例如：**case val1, val2, val3。**

### Type Switch

switch 语句还可以被用于 type-switch 来判断某个 interface 变量中实际存储的变量类型。在type-switch 不能使用 `fallthrough`，会报错`Cannot use 'fallthrough' in the type switch`。

```go
switch x.(type){
    case type:
       statement(s)      
    case type:
       statement(s)
    /* 你可以定义任意个数的case */
    default: /* 可选 */
       statement(s)
}   
```

```go
package main

import "fmt"

func main() {
    var x interface{}
    //写法一：
    switch i := x.(type) { // 带初始化语句
    case nil:
        fmt.Printf("x 的类型: %T\r\n", i)
    case int:
        fmt.Printf("x 是 int 型")
    case float64:
        fmt.Printf("x 是 float64 型")
    case func(int) float64:
        fmt.Printf("x 是 func(int) 型")
    case bool, string:
        fmt.Printf("x 是 bool 或 string 型")
    default:
        fmt.Printf("未知型")
    }
    //写法二
    var j = 0
    switch j {
    case 0:
    case 1:
        fmt.Println("1")
    case 2:
        fmt.Println("2")
    default:
        fmt.Println("def")
    }
    //写法三
    var k = 0
    switch k {
    case 0:
        println("fallthrough")
        fallthrough
        /*
            Go的switch非常灵活，表达式不必是常量或整数，执行的过程从上至下，直到找到匹配项；
            而如果switch没有表达式，它会匹配true。
            Go里面switch默认相当于每个case最后带有break，
            匹配成功后不会自动向下执行其他case，而是跳出整个switch,
            但是可以使用fallthrough强制执行后面的case代码。
        */
    case 1:
        fmt.Println("1")
    case 2:
        fmt.Println("2")
    default:
        fmt.Println("def")
    }
    //写法三
    var m = 0
    switch m {
    case 0, 1:
        fmt.Println("1")
    case 2:
        fmt.Println("2")
    default:
        fmt.Println("def")
    }
    //写法四
    var n = 0
    switch { //省略条件表达式，可当 if...else if...else
    case n > 0 && n < 10:
        fmt.Println("i > 0 and i < 10")
    case n > 10 && n < 20:
        fmt.Println("i > 10 and i < 20")
    default:
        fmt.Println("def")
    }
}  
```

```bash
fallthrough # 注意此处的输出为 println() 方法，输出顺序未按照程序编写顺序。
x 的类型: <nil>
1
1
def
```

`fmt.Println()`与`println()`输出顺序的不同是由于`fmt`中的默认输出到`stdout`，而`println`则是输出到`strerr`中，因此在 IDE 中看到的结果顺序并不是预期的顺序。 

## 条件语句 select

select 语句类似于 switch 语句，**但是select会随机执行一个可运行的case，如果没有case可运行，它将阻塞，直到有case可运行。**

select 是Go中的一个控制结构，**类似于用于通信的switch语句，每个case必须是一个通信操作，要么是发送要么是接收。**

select 随机执行一个可运行的case，如果没有case可运行，它将阻塞，直到有case可运行，一个默认的子句应该总是可运行的。

### 语法

```go
select {
    case communication clause  :
       statement(s);      
    case communication clause  :
       statement(s);
    /* 你可以定义任意数量的 case */
    default : /* 可选 */
       statement(s);
} 
```

+ 每个 case 都必须是一个通信(channel 操作)
+ 所有 channel 表达式都会被求值
+ 所有被发送的表达式都会被求值
+ 如果任意某个通信可以进行，它就执行；其他就会被忽略
+ 如果有多个 case 可以运行，select 会随机选出一个执行，其他的不会执行
+ 如果有 default 子句，则执行该语句
+ 如果没有 default 子句，select 将阻塞，直到某个通信可以运行；Go 不会重新对 channel 进行求值。

select 可以监听 channel 的数据流动，select 的用法与 switch 语法非常类似，由 select 开始的一个新的选择块，每个选择条件由 case 语句来描述

与switch语句可以选择任何使用相等比较的条件相比，select由比较多的限制，**其中最大的一条限制就是每个case语句里必须是一个IO操作**

```go
select { //不停的在这里检测
case <-chanl : //检测有没有数据可以读
//如果chanl成功读取到数据，则进行该case处理语句
case chan2 <- 1 : //检测有没有可以写
//如果成功向chan2写入数据，则进行该case处理语句

//假如没有default，那么在以上两个条件都不成立的情况下，就会在此阻塞
//一般default不写在里面，select中的default子句总是可运行的，会很消耗CPU资源
default:
//如果以上都没有符合条件，那么则进行default处理流程
} 
```

在一个select语句中，Go会按顺序从头到尾评估每一个发送和接收的语句。

如果其中的任意一个语句可以继续执行（即没有被阻塞），那么就从那些可以执行的语句中任意选择一条来使用。
如果没有任意一条语句可以执行（即所有的通道都被阻塞），那么有两种可能的情况：
①如果给出了default语句，那么就会执行default的流程，同时程序的执行会从select语句后的语句中恢复。
②如果没有default语句，那么select语句将被阻塞，直到至少有一个case可以进行下去。

### Golang select的使用及典型用法

#### 基本使用

select是Go中的一个控制结构，类似于switch语句，**用于处理异步IO操作**。select会监听case语句中channel的读写操作，当case中channel读写操作为非阻塞状态（即能读写）时，将会触发相应的动作。

```go
package main

import (
   "fmt"
   "time"
)

func main() {
   c1, c2, c3 := make(chan int), make(chan int), make(chan int)
   var i1, i2 int
   go func() {
      for {
         c1 <- 1
         fmt.Println("c2: ", <-c2)
         c3 <- 3
      }
   }()
   
   go func() {
      for{
         select {
         case i1 = <-c1:
            fmt.Printf("received %d from c1\n", i1)
         case c2 <- i2:
            fmt.Printf("sent %d to c2\n", i2)
         case i3, ok := (<-c3):  // same as: i3, ok := <-c3
            if ok {
               fmt.Printf("received %d from c3\n", i3)
            } else {
               fmt.Printf("c3 is closed\n")
            }
         }
      }
   }()
   time.Sleep(time.Second*10)
}
```

#### 典型用法

##### 超时判断

```go
//比如在下面的场景中，使用全局resChan来接受response，如果时间超过3S,resChan中还没有数据返回，则第二条case将执行
var resChan = make(chan int)
// do request
func test() {
    select {
    case data := <-resChan:
        doData(data)
    case <-time.After(time.Second * 3):
        fmt.Println("request time out")
    }
}

func doData(data int) {
    //...
}
```

##### 退出

```go
//主线程（协程）中如下：
var shouldQuit=make(chan struct{})
fun main(){
    {
        //loop
    }
    //...out of the loop
    select {
        case <-c.shouldQuit:
            cleanUp()
            return
        default:
        }
    //...
}

//再另外一个协程中，如果运行遇到非法操作或不可处理的错误，就向shouldQuit发送数据通知程序停止运行
close(shouldQuit)
```

##### 判断channel是否阻塞

```go
//在某些情况下是存在不希望channel缓存满了的需求的，可以用如下方法判断
ch := make (chan int, 5)
//...
data：=0
select {
case ch <- data:
default:
    //做相应操作，比如丢弃data。视需求而定
} 
```

## 循环语句 for

for循环是一个循环控制结构，可以执行指定次数的循环。

### 语法

Go语言的For循环有3中形式，只有其中的一种使用分号。

```go
for init; condition; post { }
for condition { }
for { }  
```

`init`： 一般为赋值表达式，给控制变量赋初值；
condition： 关系表达式或逻辑表达式，循环控制条件；
post： 一般为赋值表达式，给控制变量增量或减量。

for语句执行过程如下：

+ 先对表达式 `init` 赋初值；
+ 判别赋值表达式 `init` 是否满足给定 condition 条件，若其值为真，满足循环条件，则执行循环体内语句，然后执行 post，进入第二次循环，再判断 condition；
+ 否则判断 condition 的值为假，不满足条件，就终止for循环，执行循环体外语句。 

```go
s := "abc"

for i, n := 0, len(s); i < n; i++ { // 常见的 for 循环，支持初始化语句。
    println(s[i])
}

n := len(s)
for n > 0 {                // 替代 while (n > 0) {}
    n-- 
    println(s[n])        // 替代 for (; n > 0;) {}
}

for {                    // 替代 while (true) {}
    println(s)            // 替代 for (;;) {}
}  
```

### 循环嵌套

```go
for [condition |  ( init; condition; increment ) | Range]
{
   for [condition |  ( init; condition; increment ) | Range]
   {
      statement(s)
   }
   statement(s)
}  
```

### 无限循环

如过循环中条件语句永远不为 false 则会进行无限循环，可以通过 for 循环语句中只设置一个条件表达式来执行无限循环：

```go
for{ // 或者 for true{},两者等价
    无限循环操作
}
```

## 循环语句 range

Golang range类似迭代器操作，返回 (索引, 值) 或 (键, 值)。

for 循环的 range 格式可以对 slice、map、数组、字符串等进行迭代循环。格式如下：

```go
for key, value := range oldMap {
    newMap[key] = value
} 
```

|             | 1st value | 2nd value |              |
| ----------- | --------- | --------- | ------------ |
| string      | index     | s[index]  | unicode,rune |
| array/slice | index     | s[index]  |              |
| map         | key       | m[key]    |              |
| channel     | element   |           |              |

### range 示例

通过`_`这个特殊变量可忽略不想要的返回值。

```go
func main() {
    s := "abc"
    // 忽略 2nd value，支持 string/array/slice/map。
    for i := range s {
        println(s[i])
    }
    // 忽略 index。
    for _, c := range s {
        println(c)
    }
    // 忽略全部返回值，仅迭代。
    for range s {

    }

    m := map[string]int{"a": 1, "b": 2}
    // 返回 (key, value)。
    for k, v := range m {
        println(k, v)
    }
} 
```

```go
package main

import "fmt"

func main() {
    a := [3]int{0, 1, 2}

    for i, v := range a { // index、value 都是从复制品中取出。

        if i == 0 { // 在修改前，我们先修改原数组。
            a[1], a[2] = 999, 999
            fmt.Println(a) // 确认修改有效，输出 [0, 999, 999]。
        }

        a[i] = v + 100 // 使用复制品中取出的 value 修改原数组。

    }

    fmt.Println(a) // 输出 [100, 101, 102]。
}  
```

建议改用引用类型，其底层数据不会被复制。

```go
package main

func main() {
    s := []int{1, 2, 3, 4, 5}

    for i, v := range s { // 复制 struct slice { pointer, len, cap }。

        if i == 0 {
            s = s[:3]  // 对 slice 的修改，不会影响 range。
            s[2] = 100 // 对底层数据的修改。
        }

        println(i, v)
    }
}   
```

map、channel 是指针包装，而不像 slice 是 struct。

### for 和 for range 有什么区别

两者的使用场景不同。

+ for 可以遍历 array 和 slice、遍历 key 为整型递增的 map、遍历 string
+ for range 可以完成所有 for 能做的事情，也能做到 for 不能做的。包括遍历 key 为 string 类型的 map 并同时获取 key 和 value、遍历 channel。 

## 循环控制 Goto、Break、Continue

+ 三个语句都可以配合标签(label)使用
+ 标签名区分大小写，设置以后若不是用会造成编译错误
+ continue、break 配合标签(label)可用于多层循环跳出
+ goto 是调整执行位置，与 continue、break 配合标签(label)的结果并不相同

