# Golang 基础知识查漏补缺

[toc]

## Golang 简介

### 可见性

1. 声明在函数内部，是函数的本地值，类似private
2. 声明在函数外部，是对当前包可见(包内所有 .go 文件都可见)的全局值，类似protect
3. 声明在函数外部且**首字母大写是所有包可见的全局值,**类似 public

### 声明方式

Go 语言主要有四种声明方式：

> var (声明变量), const (声明常量), type (声明类型), func (声明函数)

Go 程序是保存在多个 .go 文件中的，文件的第一行是 package 声明，用来说明该文件属于哪个包(package)，package 声明下面是 import 声明，再下来是类型，变量，常量或函数的声明。

### 项目构建及编译问题

一个 Go 工程主要包含以下三个目录：

```lua
Go 工程
├── src -- 源代码文件
├── pkg -- 包文件
└── bin -- 相关 bin 文件
```

命令源码文件：包含 main 函数的文件。

库源码文件：不包含 main 函数的文件，主要用于编译成静态文件（.a文件）供其他包调用。

**go run**：编译源码，并且直接执行源码的 main 函数，不会在当前目录留下可执行文件，可执行文件被放在临时文件中被执行，工作目录被设置为当前目录，不能使用 “go run+包” 的方式进行编译，一般用于调试程序。

**go build**：用来测试编译包，对 库源码文件 go build 不会产生文件， 只是测试编译包是否有问题；对 命令源码文件 go build，会在当前执行 go build 命令的目录下产生可执行文件。

**go install**：主要用来生成库和工具。对 库源码文件 go install，直接编译链接整个包，会在 pkg 目录下生成 .a 静态文件， 供其他包调用；对 命令源码文件 go install， 编译+链接+生成可执行文件，会在 bin 目录下生成可执行文件。

## 内置类型

### 值类型

值类型赋值和传参会复制整个数组，开辟一个新空间存储数据，而不是原本类型的指针。因此改变副本的值，不会改变原类型的值。

>bool
>int(32 or 64), int8, int16, int32, int64
>uint(32 or 64), uint8(byte), uint16, uint32, uint64
>float32, float64
>string
>complex64, complex128
>array    -- 固定长度的数组

### 引用类型（指针类型）

指针类型赋值和传参会传递类型指针的值，改变“副本”的值，会改变与原类型的值，两者操作的是同一地址，传参实际传输的是类型的指针。

> slice -- 序列数组（最常用）
> map -- 映射
> chan -- 管道

### 内置函数

Go 语言拥有一些不需要进行导入操作就可以使用的内置函数。它们有时可以针对不同的类型进行操作，例如：len、cap 和 append，或必须用于系统级的操作，例如：panic。因此，它们需要直接获得编译器的支持。

> append          -- 用来追加元素到数组、slice中,返回修改后的数组、slice
> close           -- 主要用来关闭channel
> delete            -- 从map中删除key对应的value
> panic            -- 停止常规的goroutine  （panic和recover：用来做错误处理）
> recover         -- 允许程序定义goroutine的panic动作
> real            -- 返回complex的实部   （complex、real imag：用于创建和操作复数）
> imag            -- 返回complex的虚部
> make            -- 用来分配内存，返回Type本身(只能应用于slice, map, channel)
> new                -- 用来分配内存，主要用来分配值类型，比如int、struct。返回指向Type的指针
> cap                -- capacity是容量的意思，用于返回某个类型的最大容量（只能用于切片和 map）
> copy            -- 用于复制和连接slice，返回复制的数目
> len                -- 来求长度，比如string、array、slice、map、channel ，返回长度
> print、println     -- 底层打印函数，在部署环境中建议使用 fmt 包

### 内置接口 error

只要实现了 Error() 函数，返回值为 String 的都实现了 error 接口。

```go
type error interface{
    Error() String
}
```

## init 函数与 main 函数

### init 函数

Go 语言中 init 函数用于包（package）的初始化，该函数是 Go 语言的一个重要特性。

主要有以下特征：

1. **init函数是用于程序执行前做包的初始化的函数，比如初始化包里的变量等**

2. **每个包可以拥有多个init函数**

3. **包的每个源文件也可以拥有多个init函数**

4. **同一个包中多个init函数的执行顺序go语言没有明确的定义**

5. **不同包的init函数按照包导入的依赖关系决定该初始化函数的执行顺序**

6. **init函数不能被其他函数调用，而是在main函数执行之前，自动被调用**

### main 函数

Go 语言程序的默认入口函数（主函数）: func main()，函数体用 {} 包裹。

### init 函数和 main 函数的异同

**相同点：两个函数在定义时不能有任何参数和返回值，且由 Go 程序自动调用。**

**不同点：init 函数可以应用于任意包中，且可以重复定义多个。main 函数只能用于 main 包中，其只能定义一个。**

对于同一个 Go 文件的 init() 调用顺序从上至下；对于同一个 package 中不同文件是**按文件名字符串比较“从小到大”顺序**调用各文件中的 init() 函数；对于不同的 package  ，如果不相互依赖的话，按照 main 包中**”先 import 的后调用”**的顺序调用其包中的 init()；如果 package 中存在依赖，**则先调用最早被依赖的 package 中的 init()，最后调用 main 函数。**

## 命令

安装了 Golang 环境后，可以在命令行执行 `go` 命令查看相关命令。

```go
$ go
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]
...
```

`go env` 用于打印Go语言的环境信息。

`go run` 命令可以编译并运行命令源码文件。

`go get` 可以根据要求和实际情况从互联网上下载或更新指定的代码包及其依赖包，并对它们进行编译和安装。

`go build` 命令用于编译指定的源码文件或代码包以及它们的依赖包。

`go install` 用于编译并安装指定的代码包及它们的依赖包。

`go clean` 命令会删除掉执行其它命令时产生的一些文件和目录。

`go doc` 命令可以打印附于Go语言程序实体上的文档。可以通过把程序实体的标识符作为该命令的参数来达到查看其文档的目的。

`go test` 命令用于对Go语言编写的程序进行测试。

`go list` 命令的作用是列出指定的代码包的信息。

`go fix` 会把指定代码包的所有Go语言源码文件中的旧版本代码修正为新版本的代码。

`go vet` 是一个用于检查Go语言源码中静态错误的简单工具。

`go tool pprof` 命令来交互式的访问概要文件的内容。

## 运算符

Golang 内置**算术运算符、关系运算符、逻辑运算符、位运算符和赋值运算符**。

**注意：++（自增）和 --（自减）在 Go 语言中是单独的语句，不是运算符。**

### 赋值运算符

| 运算符 |                      描述                      |
| :----: | :--------------------------------------------: |
|   =    | 简单的赋值运算符，将一个表达式的值赋给一个左值 |
|   +=   |                  相加后再赋值                  |
|   -=   |                  相减后再赋值                  |
|   *=   |                  相乘后再赋值                  |
|   /=   |                  相除后再赋值                  |
|   %=   |                  求余后再赋值                  |
|  <<=   |                   左移后赋值                   |
|  >>=   |                   右移后赋值                   |
|   &=   |                  按位与后赋值                  |
|   l=   |                  按位或后赋值                  |
|   ^=   |                 按位异或后赋值                 |

## 下划线

### 下划线在 import 中

`import _`：当导入一个包只需要其下的 init() 都会被执行，不需要将整个包导入时可以使用`import _ 包路径 `，此时无法在通过包名来调用包中的其他函数。

### 下划线在代码中

`n, _ := f.Read(buf) `：占位符，用`_`来表示该值可以丢掉不要，使得编译器可以更好的优化，当方法返回两个结果，实际业务只需要一个结果时，可以使用`_`占位，防止因为某个变量未使用导致编译器报错。

## 变量和常量

### 变量的来历

程序运行过程中的数据都是保存在内存中，想要在代码中操作某个数据时就需要去内存上找到这个变量，在代码中直接通过内存地址操作数据会导致代码可读性非常差，所以**使用变量将数据的内存地址保存起来**，后续通过这个变量就能找到内存上对应的数据了。

### 变量类型

**变量（Variable）的功能是存储数据。**不同的变量保存的数据类型可能会不一样。经过半个多世纪的发展，编程语言已经基本形成了一套固定的类型，常见变量的数据类型有：整型、浮点型、布尔型等。Go语言中的每一个变量都有自己的类型，并且变量必须经过声明才能开始使用。

### 匿名变量

在使用多重赋值时，如果想要忽略某个值，可以使用匿名变量（anonymous variable）。 匿名变量用一个下划线_表示，例如：

```go
func foo() (int, string) {
    return 10, "Q1mi"
}
func main() {
    x, _ := foo()
    _, y := foo()
    fmt.Println("x=", x)
    fmt.Println("y=", y)
}
```

**匿名变量不占用命名空间，不会分配内存**，所以匿名变量之间不存在重复声明。 (在Lua等编程语言里，匿名变量也被叫做哑元变量。)

注意事项：

> 函数外的每个语句都必须以关键字开始（var、const、func等）
>
> := 不能在函数外使用
>
> _ 多用于占位，表示忽略值

### 常量

相对于变量，常量是恒定不变的值，多用于定义程序运行期间不会改变的那些值。 常量的声明和变量声明非常类似，只是把`var`换成了`const`，**常量在定义的时候必须赋值。**

`const`同时声明多个常量时，**如果省略了值则表示和上面一行的值相同。** 例如：

```go
const (
        n1 = 100
        n2
        n3
    )
```

上面示例中，常量`n1、n2、n3`的值都是`100`。

### iota

`iota`是`go`语言的常量计数器，**只能在常量的表达式中使用。**`iota`在`const`关键字出现时将被重置为`0`。`const`中每新增一行常量声明将使`iota`计数一次(`iota`可理解为`const`语句块中的行索引)。 使用`iota`能简化定义，在定义枚举时很有用。

## 基本类型介绍

Golang 更明确的数字类型命名，支持 Unicode，支持常用数据结构。

| 类型          | 长度(字节) | 默认值 | 说明                                      |
| :------------ | :--------- | :----- | :---------------------------------------- |
| bool          | 1          | false  |                                           |
| byte          | 1          | 0      | uint8                                     |
| rune          | 4          | 0      | Unicode Code Point, int32                 |
| int, uint     | 4或8       | 0      | 32 或 64 位                               |
| int8, uint8   | 1          | 0      | -128 ~ 127, 0 ~ 255，byte是uint8 的别名   |
| int16, uint16 | 2          | 0      | -32768 ~ 32767, 0 ~ 65535                 |
| int32, uint32 | 4          | 0      | -21亿~ 21亿, 0 ~ 42亿，rune是int32 的别名 |
| int64, uint64 | 8          | 0      |                                           |
| float32       | 4          | 0.0    |                                           |
| float64       | 8          | 0.0    |                                           |
| complex64     | 8          |        |                                           |
| complex128    | 16         |        |                                           |
| uintptr       | 4或8       |        | 以存储指针的 uint32 或 uint64 整数        |
| array         |            |        | 值类型                                    |
| struct        |            |        | 值类型                                    |
| string        |            | “”     | UTF-8 字符串                              |
| slice         |            | nil    | 引用类型                                  |
| map           |            | nil    | 引用类型                                  |
| channel       |            | nil    | 引用类型                                  |
| interface     |            | nil    | 接口                                      |
| function      |            | nil    | 函数                                      |

**Go 语言中不允许将整型转换成布尔类型，布尔类型无法参与数值运算，也无法与其他类型进行转换。**

### 多行字符串

Go 语言中要定义一个多行字符串时，就必须使用**反引号**字符：

```go
package main

import "fmt"

func main() {
   s := `hello
golang
and world` // 此处添加 tab 输出也会增加 tab
   fmt.Println(s)
}
```

```bash
hello
golang
and world
```

**反引号间换行将被作为字符串中的换行，所有的转义字符均无效，文本将会原样输出。**

### 字符串的常用操作

| 方法                                | 介绍           |
| :---------------------------------- | :------------- |
| len(str)                            | 求长度         |
| +或fmt.Sprintf                      | 拼接字符串     |
| strings.Split                       | 分割           |
| strings.Contains                    | 判断是否包含   |
| strings.HasPrefix,strings.HasSuffix | 前缀/后缀判断  |
| strings.Index(),strings.LastIndex() | 子串出现的位置 |
| strings.Join(a[]string, sep string) | join操作       |

### byte 和 rune 类型

组成每个字符串的元素叫做“字符”，**可以通过遍历或者单个获取字符串元素获得字符，字符用单引号（’）包裹**起来，如：

```go
    var a := '中'
    var b := 'x'
```

Go 语言的字符有以下两种：

> uint8类型，或者叫 byte 型，**代表了ASCII码的一个字符。**
>
> rune类型，**代表一个 UTF-8字符。**

当需要处理中文、日文或者其他复合字符时，则需要用到`rune`类型。`rune`类型实际是一个`int32`。**Go 使用了特殊的 `rune` 类型来处理 `Unicode`，让基于 `Unicode`的文本处理更为方便**，也可以使用 `byte` 型进行默认字符串处理，性能和扩展性都有照顾。

```go
package main

import "fmt"

// 遍历字符串
func traversalString() {
   s := "pprof.cn博客"
   for i := 0; i < len(s); i++ { //byte
      fmt.Printf("%v(%c) ", s[i], s[i])
   }
   fmt.Println()
   for _, r := range s { //rune
      fmt.Printf("%v(%c) ", r, r)
   }
   fmt.Println()
}

func main() {
   traversalString()
}
```

```bash
112(p) 112(p) 114(r) 111(o) 102(f) 46(.) 99(c) 110(n) 229(å) 141() 154() 229(å) 174(®) 162(¢) 
112(p) 112(p) 114(r) 111(o) 102(f) 46(.) 99(c) 110(n) 21338(博) 23458(客) 
```

因为UTF8编码下一个中文汉字由`3~4`个字节组成，**不能简单的按照字节去遍历一个包含中文的字符串，否则就会出现输出中第一行的结果。**字符串底层是一个byte数组，所以可以和[]byte类型相互转换。字符串是不能修改的，字符串是由byte字节组成，所以字符串的长度是byte字节的长度。 rune类型用来表示utf8字符，**一个rune字符由一个或多个byte组成。**

### 修改字符串

要修改字符串，**需要先将其转换成`[]rune或[]byte`**，完成后再转换为`string`。**无论哪种转换，都会重新分配内存，并复制字节数组。**

```go
package main

import "fmt"

func changeString() {
	s1 := "hello"
	// 强制类型转换
	byteS1 := []byte(s1)
	byteS1[0] = 'H'
	fmt.Println(string(byteS1))

	s2 := "博客"
	runeS2 := []rune(s2)
	runeS2[0] = '狗'
	fmt.Println(string(runeS2))
}

func main() {
	changeString()
}
```

### 类型转换

**Go 语言中只有强制类型转换，没有隐式类型转换。**强制类型转换的基本语法：`T(表达式)`，其中T表示要转换的类型，表达式包括变量、复杂算子和函数返回值等。

比如计算直角三角形的斜边长时使用math包的Sqrt()函数，该函数接收的是float64类型的参数，而变量a和b都是int类型的，**这个时候就需要将a和b强制类型转换为float64类型。**

```go
package main

import (
   "fmt"
   "math"
)

func sqrtDemo() {
   var a, b = 3, 4
   var c int
   // math.Sqrt()接收的参数是float64类型，需要强制转换
   c = int(math.Sqrt(float64(a*a + b*b)))
   fmt.Println(c)
}

func main() {
   sqrtDemo()
}
```

## 数组

数组定义时，**数组长度必须是常量，一旦定义，长度不能修改**。**数组是值类型**，复制和传参会复制整个数组，而不是指针，因此改变副本的值不会改变原本数组的值。

### 数组初始化

#### 一维数组

```go
package main

import (
    "fmt"
)

var arr0 [5]int = [5]int{1, 2, 3}
var arr1 = [5]int{1, 2, 3, 4, 5}
var arr2 = [...]int{1, 2, 3, 4, 5, 6}
var str = [5]string{3: "hello world", 4: "tom"}

func main() {
    a := [3]int{1, 2}           // 未初始化元素值为 0。
    b := [...]int{1, 2, 3, 4}   // 通过初始化值确定数组长度。
    c := [5]int{2: 100, 4: 200} // 使用引号初始化元素。
    d := [...]struct {
        name string
        age  uint8
    }{
        {"user1", 10}, // 可省略元素类型。
        {"user2", 20}, // 别忘了最后一行的逗号。
    }
    fmt.Println(arr0, arr1, arr2, str)
    fmt.Println(a, b, c, d)
}
```

```bash
[1 2 3 0 0] [1 2 3 4 5] [1 2 3 4 5 6] [   hello world tom]
[1 2 0] [1 2 3 4] [0 0 100 0 200] [{user1 10} {user2 20}]
```

#### 多维数组

```go
package main

import (
    "fmt"
)

var arr0 [5][3]int // 初始化不赋值必须设置行列值，不可使用 ...
var arr1 [2][3]int = [...][3]int{{1, 2, 3}, {7, 8, 9}} // 存在初始化值时，[2][3]int 可省略

func main() {
    a := [2][3]int{{1, 2, 3}, {4, 5, 6}}
    b := [...][2]int{{1, 1}, {2, 2}, {3, 3}} // 第 2 纬度不能用 "..."。
    fmt.Println(arr0, arr1)
    fmt.Println(a, b)
}
```

```bash
[[0 0 0] [0 0 0] [0 0 0] [0 0 0] [0 0 0]] [[1 2 3] [7 8 9]]
[[1 2 3] [4 5 6]] [[1 1] [2 2] [3 3]]
```

值拷贝行为会造成性能问题，通常会建议使用 slice 或数组指针。

```go
package main

import (
    "fmt"
)

func test(x [2]int) {
    fmt.Printf("x: %p\n", &x)
    x[1] = 1000
}

func main() {
    a := [2]int{}
    fmt.Printf("a: %p\n", &a)

    test(a)
    fmt.Println(a)
}
```

```bash
a: 0xc42007c010
x: 0xc42007c030
[0 0]
```

内置函数 len 和 cap 都返回数组长度（元素数量）。

#### 多维数组遍历

```go
package main

import (
    "fmt"
)

func main() {

    var f [2][3]int = [...][3]int{{1, 2, 3}, {7, 8, 9}}

    for k1, v1 := range f {
        for k2, v2 := range v1 {
            fmt.Printf("(%d,%d)=%d ", k1, k2, v2)
        }
        fmt.Println()
    }
}
```

```bash
(0,0)=1 (0,1)=2 (0,2)=3 
(1,0)=7 (1,1)=8 (1,2)=9 
```

### 数组拷贝和传参

```go
package main

import "fmt"

func printArr(arr *[5]int) { // 传入数组指针
    arr[0] = 10
    for i, v := range arr {
        fmt.Println(i, v)
    }
}

func main() {
    var arr1 [5]int
    printArr(&arr1)
    fmt.Println(arr1)
    arr2 := [...]int{2, 4, 6, 8, 10}
    printArr(&arr2)
    fmt.Println(arr2)
}
```

### 数组练习

#### 求数组所有元素之和

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// 求元素和
func sumArr(a [10]int) int {
    var sum int = 0
    for i := 0; i < len(a); i++ {
        sum += a[i]
    }
    return sum
}

func main() {
    // 若想做一个真正的随机数，要种子
    // seed()种子默认是1
    // rand.Seed(1)
    rand.Seed(time.Now().Unix())

    var b [10]int
    for i := 0; i < len(b); i++ {
        // 产生一个0到1000随机数
        b[i] = rand.Intn(1000)
    }
    sum := sumArr(b)
    fmt.Printf("sum=%d\n", sum)
}
```

#### 找出数组中和为某值的两个元素的下标，例如数组[1,3,5,8,7]，找出两个元素之和等于8的下标分别是（0，4）和（1，2）

```go
package main

import "fmt"

// 求元素和，是给定的值
func myTest(a [5]int, target int) { // 传参必须执行数组长度
   // 遍历数组
   for i := 0; i < len(a); i++ {
      other := target - a[i]
      // 继续遍历
      for j := i + 1; j < len(a); j++ {
         if a[j] == other {
            fmt.Printf("(%d,%d)\n", i, j)
         }
      }
   }
}

func main() {
   b := [5]int{1, 3, 5, 8, 7}
   myTest(b, 8)
}
```

## 切片

需要说明，slice 并不是数组或数组指针**。它通过内部指针和相关属性引用数组片段，以实现变长方案。**

>**切片是数组的一个引用，因此切片是引用类型。**但自身是结构体，值拷贝传递。
>
>切片的长度可以改变，因此，切片是一个可变的数组。
>
>切片遍历方式和数组一样，可以用len()求长度。表示可用元素数量，读写操作不能超过该限制。 
>
>**cap可以求出slice最大扩张容量，不能超出数组限制。**0 <= len(slice) <= len(array)，其中array是slice引用的数组。
>
>切片的定义：var 变量名 []类型，比如 var str []string /  var arr []int。
>
>如果 slice == nil，那么 len、cap 结果都等于 0。

### 切片初始化

```go
package main

import (
    "fmt"
)

// 全局
var arr = [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var slice0 []int = arr[2:8]
var slice1 []int = arr[0:6]        //可以简写为 var slice []int = arr[:end]
var slice2 []int = arr[5:10]       //可以简写为 var slice []int = arr[start:]
var slice3 []int = arr[0:len(arr)] //var slice []int = arr[:]
var slice4 = arr[:len(arr)-1]      //去掉切片的最后一个元素
func main() {
    fmt.Printf("全局变量：arr %v\n", arr)
    fmt.Printf("全局变量：slice0 %v\n", slice0)
    fmt.Printf("全局变量：slice1 %v\n", slice1)
    fmt.Printf("全局变量：slice2 %v\n", slice2)
    fmt.Printf("全局变量：slice3 %v\n", slice3)
    fmt.Printf("全局变量：slice4 %v\n", slice4)
    fmt.Printf("-----------------------------------\n")
    // 局部
    arr2 := [...]int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
    slice5 := arr[2:8]
    slice6 := arr[0:6]         //可以简写为 slice := arr[:end]
    slice7 := arr[5:10]        //可以简写为 slice := arr[start:]
    slice8 := arr[0:len(arr)]  //slice := arr[:]
    slice9 := arr[:len(arr)-1] //去掉切片的最后一个元素
    fmt.Printf("局部变量： arr2 %v\n", arr2)
    fmt.Printf("局部变量： slice5 %v\n", slice5)
    fmt.Printf("局部变量： slice6 %v\n", slice6)
    fmt.Printf("局部变量： slice7 %v\n", slice7)
    fmt.Printf("局部变量： slice8 %v\n", slice8)
    fmt.Printf("局部变量： slice9 %v\n", slice9)
}
```

```bash
全局变量：arr [0 1 2 3 4 5 6 7 8 9]
全局变量：slice0 [2 3 4 5 6 7]
全局变量：slice1 [0 1 2 3 4 5]
全局变量：slice2 [5 6 7 8 9]
全局变量：slice3 [0 1 2 3 4 5 6 7 8 9]
全局变量：slice4 [0 1 2 3 4 5 6 7 8]
-----------------------------------
局部变量： arr2 [9 8 7 6 5 4 3 2 1 0]
局部变量： slice5 [2 3 4 5 6 7]
局部变量： slice6 [0 1 2 3 4 5]
局部变量： slice7 [5 6 7 8 9]
局部变量： slice8 [0 1 2 3 4 5 6 7 8 9]
局部变量： slice9 [0 1 2 3 4 5 6 7 8]
```

### 通过 make 来创建切片

```go
var slice []type = make([]type, len) // var 后的 []type 可省略
var slice []type = make([]type, len， cap)
slice := make([]type, len)
slice := make([]type, len, cap)
```

```go
package main

import (
   "fmt"
)

var slice0 []int = make([]int, 10)
var slice1 = make([]int, 10)
var slice2 = make([]int, 10, 10)

func main() {
   fmt.Printf("make全局slice0 ：%v\n", slice0)
   fmt.Printf("make全局slice1 ：%v\n", slice1)
   fmt.Printf("make全局slice2 ：%v\n", slice2)
   fmt.Println("--------------------------------------")
   slice3 := make([]int, 10)
   slice4 := make([]int, 10)
   slice5 := make([]int, 10, 10)
   fmt.Printf("make局部slice3 ：%v\n", slice3)
   fmt.Printf("make局部slice4 ：%v\n", slice4)
   fmt.Printf("make局部slice5 ：%v\n", slice5)
}
```

```bash
make全局slice0 ：[0 0 0 0 0 0 0 0 0 0]
make全局slice1 ：[0 0 0 0 0 0 0 0 0 0]
make全局slice2 ：[0 0 0 0 0 0 0 0 0 0]
--------------------------------------
make局部slice3 ：[0 0 0 0 0 0 0 0 0 0]
make局部slice4 ：[0 0 0 0 0 0 0 0 0 0]
make局部slice5 ：[0 0 0 0 0 0 0 0 0 0]
```

slice 实际进行读写操作的是底层数组，需要注意在 slice 中与数组中索引号的差别。可直接创建 slice 对象，自动分配底层数组。

```go
package main

import "fmt"

func main() {
	// 通过初始化表达式构造，可使用索引号，即 8:100 表示索引号为 8 的值为 100
	s1 := []int{0, 1, 2, 3, 8:100}
	// 此处若还需添加数据，需使用 s1 = append(s1, 999)，此操作后，len=10，cap=18
	fmt.Println(s1, len(s1), cap(s1))

	s2 := make([]int, 6, 8) // 使用 make 创建，指定 len 和 cap 值。
	fmt.Println(s2, len(s2), cap(s2))

	s3 := make([]int, 6) // 省略 cap，相当于 cap = len。
	fmt.Println(s3, len(s3), cap(s3))
}
```

```bash
[0 1 2 3 0 0 0 0 100] 9 9 # 注意此处索引号为8的值为 100，且长度和容量为9
[0 0 0 0 0 0] 6 8
[0 0 0 0 0 0] 6 6
```

使用 make 动态创建slice，避免了数组必须用常量做长度的麻烦。还可用指针直接访问底层数组，退化成普通数组操作。

### 用 append 内置函数操作切片（切片追加）

```go
package main

import (
   "fmt"
)

func main() {

   var a = []int{1, 2, 3}
   fmt.Printf("slice a : %v\n", a)
   var b = []int{4, 5, 6}
   fmt.Printf("slice b : %v\n", b)
   c := append(a, b...)
   fmt.Printf("slice c : %v\n", c)
   d := append(c, 7)
   fmt.Printf("slice d : %v\n", d)
   e := append(d, 8, 9, 10)
   fmt.Printf("slice e : %v\n", e)

}
```

```bash
slice a : [1 2 3]
slice b : [4 5 6]
slice c : [1 2 3 4 5 6]
slice d : [1 2 3 4 5 6 7]
slice e : [1 2 3 4 5 6 7 8 9 10]
```

#### append：向 slice 尾部添加数据，返回新的 slice 对象。

```go
package main

import (
    "fmt"
)

func main() {

    s1 := make([]int, 0, 5)
    fmt.Printf("%p\n", &s1)

    s2 := append(s1, 1)
    fmt.Printf("%p\n", &s2)

    fmt.Println(s1, s2)

}
```

```bash
0xc000096060
0xc000096078 # 与最初的而地址不同
[] [1]
```

#### append：超出原 slice.cap 限制，会重新分配底层数组，即便原数组并未填满。

```go
package main

import (
    "fmt"
)

func main() {

    data := [...]int{0, 1, 2, 3, 4, 10: 0}
    // 截取data前两位，并设置切片容量为3，此处不设置 cap(s)=11,即data容量
    s := data[:2:3] // [左开位:右避位:容量]

    s = append(s, 100, 200) // 一次 append 两个值，超出 s.cap 限制。

    fmt.Println(s, data)         // 重新分配底层数组，与原数组无关。
    fmt.Println(&s[0], &data[0]) // 比对底层数组起始指针。
}
```

```bash
[0 1 100 200] [0 1 2 3 4 0 0 0 0 0 0]
0xc00000c3f0 0xc00006a060 # 进行append()后未超出容量限制，两者地址应该相同
```

从输出结果可以看出，append 后的 s 重新分配了底层数组，并复制数据。**如果只追加一个值，则不会超过 s.cap 限制，也就不会重新分配。**
**通常以 2 倍容量重新分配底层数组。**在大批量添加数据时，建议一次性分配足够大的空间，以减少内存分配和数据复制开销。或初始化足够长的 len 属性，改用索引号进行操作。及时释放不再使用的 slice 对象，避免持有过期数组，造成 GC 无法回收。

### 切片拷贝

```go
package main

import (
    "fmt"
)

func main() {

    s1 := []int{1, 2, 3, 4, 5}
    fmt.Printf("slice s1 : %v\n", s1)
    s2 := make([]int, 10)
    fmt.Printf("slice s2 : %v\n", s2)
    copy(s2, s1) // 将s1赋值给s2
    fmt.Printf("copied slice s1 : %v\n", s1)
    fmt.Printf("copied slice s2 : %v\n", s2)
    s3 := []int{1, 2, 3}
    fmt.Printf("slice s3 : %v\n", s3)
    s3 = append(s3, s2...)
    fmt.Printf("appended slice s3 : %v\n", s3)
    s3 = append(s3, 4, 5, 6)
    fmt.Printf("last slice s3 : %v\n", s3)

}
```

```bash
slice s1 : [1 2 3 4 5]
slice s2 : [0 0 0 0 0 0 0 0 0 0]
copied slice s1 : [1 2 3 4 5]
copied slice s2 : [1 2 3 4 5 0 0 0 0 0]
slice s3 : [1 2 3]
appended slice s3 : [1 2 3 1 2 3 4 5 0 0 0 0 0]
last slice s3 : [1 2 3 1 2 3 4 5 0 0 0 0 0 4 5 6]
```

copy：函数 copy 在两个 slice 间复制数据，**复制长度以 len 小的为准**（如例，若copy(s1,s2)，则超过s1长度的值不会进行复制）。两个 slice 可指向同一底层数组，允许元素区间重叠。

```go
package main

import (
    "fmt"
)

func main() {

    data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    fmt.Println("array data : ", data)
    s1 := data[8:]
    s2 := data[:5]
    fmt.Printf("slice s1 : %v\n", s1)
    fmt.Printf("slice s2 : %v\n", s2)
    copy(s2, s1)
    fmt.Printf("copied slice s1 : %v\n", s1)
    fmt.Printf("copied slice s2 : %v\n", s2)
    fmt.Println("last array data : ", data)

}
```

```bash
array data :  [0 1 2 3 4 5 6 7 8 9]
slice s1 : [8 9]
slice s2 : [0 1 2 3 4]
copied slice s1 : [8 9]！
copied slice s2 : [8 9 2 3 4]
last array data :  [8 9 2 3 4 5 6 7 8 9] # 原数组也被修改了！
```

### slice 遍历

```go
package main

import (
   "fmt"
)

func main() {

   data := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
   slice := data[:]
   for index, value := range slice { // 索引，值
      fmt.Printf("inde : %v , value : %v\n", index, value)
   }

}
```

### 切片 resize（调整大小）

```go
package main

import (
	"fmt"
)

func main() {
	var a = []int{1, 3, 4, 5}
	fmt.Printf("slice a : %v , len(a) : %v, cap(a) : %v\n", a, len(a), cap(a))
	b := a[1:2] // 此处a[0:4]不会报错，b中值与a中相同
	fmt.Printf("slice b : %v , len(b) : %v, cap(b) : %v\n", b, len(b), cap(b))
	c := b[0:3] // 此处为b中索引号，而不是a，故b[0,4]会报错，超出容量
	fmt.Printf("slice c : %v , len(c) : %v, cap(c) : %v\n", c, len(c), cap(c))
}
```

```bash
slice a : [1 3 4 5] , len(a) : 4, cap(a) : 4
slice b : [3] , len(b) : 1, cap(b) : 3
slice c : [3 4 5] , len(c) : 3, cap(c) : 3
```

**a[x:y:z] 切片内容 [x:y] 切片长度: y-x 切片容量:z-x**

```go
package main

import (
    "fmt"
)

func main() {
    slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    d1 := slice[6:8]
    fmt.Println(d1, len(d1), cap(d1))
    d2 := slice[:6:8]
    fmt.Println(d2, len(d2), cap(d2))
}
```

```go
[6 7] 2 4
[0 1 2 3 4 5] 6 8
```

### 数组和切片的内存布局

![数组和切片的内存布局](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220312133626518.png "数组和切片的内存布局")

### 字符串和切片

**string底层就是一个byte的数组，因此，也可以进行切片操作。**

```go
package main

import (
    "fmt"
)

func main() {
    str := "hello world"
    s1 := str[0:5]
    fmt.Println(s1) // 输出 hello

    s2 := str[6:]
    fmt.Println(s2) // 输出 world
}
```

string本身是不可变的，因此要改变string中字符。需要如下操作：

#### 英文字符串

```go
package main

import (
    "fmt"
)

func main() {
    str := "Hello world"
    s := []byte(str) // 中文字符需要用[]rune(str)
    s[6] = 'G'
    s = s[:8]
    s = append(s, '!')
    str = string(s)
    fmt.Println(str) // 输出 Hello Go!
}
```

#### 含中文字符串

```go
package main

import (
    "fmt"
)

func main() {
    str := "你好，世界！hello world！"
    s := []rune(str) 
    s[3] = '够'
    s[4] = '浪'
    s[12] = 'g'
    s = s[:14]
    str = string(s)
    fmt.Println(str) // 输出 你好，够浪！hello go
}
```

#### 数组/切片转字符串

```go
strings.Replace(strings.Trim(fmt.Sprint(array_or_slice), "[]"), " ", ",", -1)
```

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(slice)
	// fmt.Sprint() 返回值类型为 string，fmt.Println() 返回值类型为 (int, error)
	str := strings.Replace(strings.Trim(fmt.Sprint(slice), "[]"), " ", ",", -1)
	fmt.Println(str)
}
```

```bash
[0 1 2 3 4 5 6 7 8 9]
0,1,2,3,4,5,6,7,8,9
```

## slice 底层实现

切片是 Go 中的一种基本的数据结构，使用这种结构可以用来管理数据集合。切片的设计想法是由动态数组概念而来，为了开发者可以更加方便的使一个数据结构可以自动增加和减少。**但是切片本身并不是动态数组或者数组指针。**切片常见的操作有 reslice、append、copy。与此同时，切片还具有可索引，可迭代的优秀特性。源码实现在 `slice.go`中。

在 Go 中，Go 数组是值类型，**赋值和函数传参操作都会复制整个数组数据。**

```go
package main

import (
	"fmt"
)

func main() {
	arrayA := [2]int{100, 200}
	var arrayB [2]int

	arrayB = arrayA

	fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA)
	fmt.Printf("arrayB : %p , %v\n", &arrayB, arrayB)

	testArray(arrayA)
}

func testArray(x [2]int) {
	fmt.Printf("func Array : %p , %v\n", &x, x)
}
```

```bash
arrayA : 0xc00000a0c0 , [100 200]
arrayB : 0xc00000a0d0 , [100 200]
func Array : 0xc00000a120 , [100 200]
```

可以看到，**三个内存地址都不同**，这也就验证了 Go 中数组赋值和函数传参都是值复制的。此时当数组较大，由于每次传参都必须复制一遍，会消耗掉大量内存。此时有人想到，函数传参时传数组的指针。

```go
package main

import (
   "fmt"
)

func main() {
   arrayA := []int{100, 200}
   testArrayPoint(&arrayA)   // 1.传数组指针
   arrayB := arrayA[:]
   testArrayPoint(&arrayB)   // 2.传切片
   fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA)
}

func testArrayPoint(x *[]int) {
   fmt.Printf("func Array : %p , %v\n", x, *x)
   (*x)[1] += 100
}
```

```bash
func Array : 0xc000004078 , [100 200]
func Array : 0xc0000040a8 , [100 300]
arrayA : 0xc000004078 , [100 400] # 地址与传入数组相同
```

此时就算传入大数组，也只需在栈上分配一个8字节的内存给指针就可以了。能更有效的利用内存，且性能相较之前有较大提升。但传指针会有一个弊端，从打印结果可以看到，第一行和第三行指针地址都是同一个，**当原数组的指针指向更改了，函数里面的指针指向也会随之更改。**

此时切片的优势就表现出来了。用切片传数组参数，既可以达到节约内存的目的，也可以达到合理处理好共享内存的问题。打印结果第二行就是切片，**切片的指针和原来数组的指针是不同的。**

由此可以得出结论：把第一个大数组传递给函数会消耗很多内存，采用切片的方式传参可以避免上述问题。**切片是引用传递，所以它们不需要使用额外的内存并且比使用数组更有效率。**

但存在反例：

```go
// 测试文件命名需要加上 _test 后缀，eg: filename_test.go
package main

import "testing"

func array() [1024]int {
    var x [1024]int
    for i := 0; i < len(x); i++ {
        x[i] = i
    }
    return x
}

func slice() []int {
    x := make([]int, 1024)
    for i := 0; i < len(x); i++ {
        x[i] = i
    }
    return x
}

func BenchmarkArray(b *testing.B) { // Benchmark 基准测试
    for i := 0; i < b.N; i++ {
        array()
    }
}

func BenchmarkSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        slice()
    }
}
```

进行性能测试，并禁用内联和优化，来观察切片的堆上内存分配的情况。

```bash
> go test -bench . -benchmem -gcflags "-N -l" # 进行性能测试的执行，在此文件夹下需要有测试文件
goos: windows
goarch: amd64
pkg: awesomeProject/test
cpu: AMD Ryzen 7 3700X 8-Core Processor
BenchmarkArray-16         659865              1809 ns/op               0 B/op          0 allocs/op
BenchmarkSlice-16         461538              2594 ns/op            8192 B/op          1 allocs/op
PASS
ok      awesomeProject/test     2.572s
```

在测试 Array 的时候，用的是16核，循环次数是659865，平均每次执行时间是1809ns，每次执行堆上分配内存总量是0，分配次数是0 。

而切片的结果就“差”一点，用的是16核，循环次数是461538，平均每次执行时间是2594ns，**但是每次执行一次堆上分配内存总量是8192，分配次数是1 。**

通过以上对比，并非所有时候都适合用切片代替数组，**因为切片底层数组可能会在堆上分配内存，而且小数组在栈上拷贝的消耗也未必比 make 分配内存的消耗大。**

### 切片的数据结构

切片本身并不是动态数组或者数组指针。**它内部实现的数据结构通过指针引用底层数组，设定相关属性将数据读写操作限定在指定的区域内。切片本身是一个只读对象，其工作机制类似于对数组指针的一种封装。**

切片（slice）是对数组一个连续片段的引用，所以切片是一个引用类型（因此更类似于 C/C++ 中的数组类型，或者 Python 中的 list 类型）。这个片段可以是整个数组，或者是由起始和终止索引标识的一些项的子集。需要注意的是，终止索引标识的项不包括在切片内。**切片提供了一个与指向数组的动态窗口。**

给定项的切片索引可能比相关数组的相同元素的索引小。和数组不同的是，**切片的长度可以在运行时修改，最小为 0 最大为相关数组的长度，切片是一个长度可变的数组。**

slice 的数据结构定义如下：

```go
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```

![slice 数据结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220312142244154.png "slice 数据结构")

切片的结构体由3部分构成，Pointer 是指向一个数组的指针，len 代表当前切片的长度，cap 是当前切片的容量。cap 总是大于等于 len 的。

![slice 内存空间](https://raw.githubusercontent.com/tonshz/test/master/img/202203121430876.png)

如果想从 slice 中得到一块内存地址，可以这样做：

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	s := make([]byte, 200)
    ptr := unsafe.Pointer(&s) // 必须使用unsafe.Pointer()转换类型，否则无法获得内存地址
	ptr0 := unsafe.Pointer(&s[0])
	ptr1 := unsafe.Pointer(&s[1])
	fmt.Println(ptr, ptr0, ptr1)
}
```

```bash
0xc000096060 0xc0000d6000 0xc0000d6001
```

从 Go 的内存地址中构造一个 slice,可以这样做：

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	length := 10
	var ptr unsafe.Pointer
	var s1 = struct {
		addr uintptr
		len int
		cap int
    }{uintptr(ptr), length, length} // 原文此处无 uintptr(ptr) 的类型转换
	s := *(*[]byte)(unsafe.Pointer(&s1))
	fmt.Println(s)
}
```

构造一个虚拟的结构体，把 slice 的数据结构拼出来。

当然还有更加直接的方法，**在 Go 的反射中就存在一个与之对应的数据结构 SliceHeader**，可以用它来构造一个 slice。

```go
package main

import (
   "reflect"
   "unsafe"
)
func main() {
   var length = 10
   var ptr unsafe.Pointer
   var o []byte
   sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&o)))
   sliceHeader.Cap = length
   sliceHeader.Len = length
   sliceHeader.Data = uintptr(ptr)
}
```

### 创建切片(makeslice)

make 函数允许在运行期动态指定数组长度，绕开了数组类型必须使用编译期常量的限制。创建切片有两种形式，make 创建切片，空切片。

#### make 与切面字面量

![make 创建切片](https://raw.githubusercontent.com/tonshz/test/master/img/202203121652076.png "make 创建切片")

上图是用 make 函数创建的一个 len = 4， cap = 6 的切片。内存空间申请了6个 int 类型的内存大小。由于 len = 4，所以后面2个暂时访问不到，但是容量还是在的。这时候数组里面每个变量都是0 。

除了 make 函数可以创建切片以外，字面量也可以创建切片。

```go
slice := []int{10,20,30,40,50,60}
```

![字面量创建切片](https://raw.githubusercontent.com/tonshz/test/master/img/202203121654364.png "字面量创建切片")

这里是用字面量创建的一个 len = 6，cap = 6 的切片，这时候数组里面每个元素的值都初始化完成了。**需要注意的是 [ ] 里面不要写数组的容量，因为如果写了个数以后就是数组了，而不是切片了。**

![同一数组创建多个切片](https://raw.githubusercontent.com/tonshz/test/master/img/202203121655826.png "同一数组创建多个切片")

还有一种简单的字面量创建切片的方法，如上图。上图就 Slice A 创建出了一个 len = 3，cap = 3 的切片。从原数组的第二位元素(0是第一位)开始切，一直切到第四位为止(不包括第五位)。同理，Slice B 创建出了一个 len = 2，cap = 4 的切片。

#### nil 切片和空切片

nil 切片和空切片也是常用的。

```go
var slice []int
```

![nil 切片](https://raw.githubusercontent.com/tonshz/test/master/img/202203121656572.png "nil 切片")

nil 切片被用在很多标准库和内置函数中，描述一个不存在的切片的时候，就需要用到 nil 切片。比如函数在发生异常的时候，返回的切片就是 nil 切片，nil 切片的指针指向 nil。

空切片一般会用来表示一个空的集合。比如数据库查询，一条结果也没有查到，那么就可以返回一个空切片。

```go
slice := make([]int, 0)
slice := []int{}
```

![空切片](https://raw.githubusercontent.com/tonshz/test/master/img/%E7%A9%BA%E5%88%87%E7%89%87.png "空切片")

空切片和 nil 切片的区别在于，**空切片指向的地址不是nil，指向的是一个内存地址**，但是它没有分配任何内存空间，即底层元素包含0个元素。最后需要说明的一点是。不管是使用 nil 切片还是空切片，**对其调用内置函数 append，len 和 cap 的效果都是一样的。**

### 切片扩容(growslice)

切片扩容主要需要关注的有两点，一个是扩容时候的策略，还有一个就是扩容时生成全新的内存地址还是在原来的地址后追加。

#### 扩容策略

![扩容策略](https://raw.githubusercontent.com/tonshz/test/master/img/202203121708534.png "扩容策略")

从图上可以看出新的切片和之前的切片已经不同了，因为新的切片更改了一个值，并没有影响到原来的数组，**新切片指向的数组是一个全新的数组**。并且 cap 容量也发生了变化。

Go 中切片扩容的策略是这样的：

**如果切片的容量小于 1024 个元素，于是扩容的时候就翻倍增加容量。**一旦元素个数超过 1024 个元素，**那么增长因子就变成 1.25** ，即每次增加原来容量的四分之一。**扩容扩大的容量都是针对原来的容量而言的，而不是针对原来数组的长度而言的。**

#### 新地址？老地址？

扩容之后的切片不一定是新的。

##### 情况一

![扩容后数组内存地址](https://raw.githubusercontent.com/tonshz/test/master/img/202203121710533.png "扩容后数组内存地址")

可以看到，此处情况下扩容以后并没有新建一个新的数组，扩容前后的数组都是同一个，这也就导致了新的切片修改了一个值，也影响到了老的切片了。并且 append() 操作也改变了原来数组里面的值。一个 append() 操作影响了这么多地方，如果原数组上有多个切片，那么这些切片都会被影响！无意间就产生了莫名的 bug！

这种情况，**由于原数组还有容量可以扩容**，所以执行 append() 操作以后，**会在原数组上直接操作**，所以这种情况下，扩容以后的数组还是指向原来的数组。

这种情况也极容易出现在字面量创建切片时候，第三个参数 cap 传值的时候，如果用字面量创建切片，cap 并不等于指向数组的总容量，那么这种情况就会发生。

```go
slice := array[1:2:3]
```

上面这种情况非常危险，极度容易产生 bug 。建议用字面量创建切片的时候，避免共享原数组导致的 bug。

##### 情况二

情况二其实就是在扩容策略里面举的例子，在那个例子中之所以生成了新的切片，是因为**原来数组的容量已经达到了最大值**，再想扩容， Go 默认会先开一片内存区域，把原来的值拷贝过来，然后再执行 append() 操作。这种情况丝毫不影响原数组。

所以建议尽量避免情况一，尽量使用情况二，避免 bug 产生。

### 切片拷贝(slicecopy)

在 slicecopy 方法中会把源切片值(即 fm Slice )中的元素复制到目标切片(即 to Slice )中，并返回被复制的元素个数，copy 的两个类型必须一致。slicecopy 方法最终的复制结果取决于较短的那个切片，**当较短的切片复制完成，整个复制过程就全部完成了。**

![切片拷贝](https://raw.githubusercontent.com/tonshz/test/master/img/202203121723655.png "切片拷贝")

在切片拷贝中有一个需要注意的问题，在使用 range 方式去遍历一个切片，**拿到的 Value 其实是切片里面的值拷贝，每次打印 Value 的地址都不变。**由于 Value 是值拷贝的，并非引用传递，所以直接改 Value 是达不到更改原切片值的目的的，**需要通过 &slice[index] 获取真实的地址。**

```go
package main

import (
   "fmt"
)

func main() {
   slice := []int{10, 20, 30, 40}
   for index, value := range slice {
      fmt.Printf("value = %d , value-addr = %x , slice-addr = %x\n", value, &value, &slice[index])
   }
}
```

```bash
value = 10 , value-addr = c00000a0a8 , slice-addr = c00000e280
value = 20 , value-addr = c00000a0a8 , slice-addr = c00000e288
value = 30 , value-addr = c00000a0a8 , slice-addr = c00000e290
value = 40 , value-addr = c00000a0a8 , slice-addr = c00000e298
```

![slice 与 range](https://raw.githubusercontent.com/tonshz/test/master/img/202203121725054.png "slice 与 range")

## 指针

Go语言中的函数传参都是值拷贝，当我们想要修改某个变量的时候，我们可以创建一个指向该变量地址的指针变量。传递数据使用指针，而无须拷贝数据。**类型指针不能进行偏移和运算。**Go语言中的指针操作非常简单，只需要记住两个符号：`&`（取地址）和`*`（根据地址取值）。

### 指针地址和指针类型

每个变量在运行时都拥有一个地址，这个地址代表变量在内存中的位置。Go语言中使用&字符放在变量前面对变量进行“取地址”操作。 Go语言中的值类型`（int、float、bool、string、array、struct）`都有对应的指针类型，如：`*int、*int64、*string`等。

取变量指针的语法：`ptr := &v`，其中 v 代表被取地址的变量，类型为 T, ptr 用于接收地址的变量，ptr 的类型为 *T，称为 T 的指针类型，\* 代表指针。下图表示`b := &a`。

![b := &a](https://raw.githubusercontent.com/tonshz/test/master/img/202203122015539.png "b := &a")

### 指针取值

在对普通变量使用 & 操作获取地址后可以获得这个变量的指针，然后可以对指针使用 \* 操作，进行指针取值。取地址操作符 & 与取值操作符 \* 是一对互补操作符，& 取出地址， \* 根据地址取出地址指向的值。

变量、指针地址、指针变量、取地址、取值的相互关系和特性如下：

> 1. 对变量进行取地址（&）操作，可以获得这个变量的指针变量。
> 2. 指针变量的值时指针地址。
> 3. 对指针变量进行取值（\*）操作，可以获得指针变量指向的原变量的值。

### 空指针

当一个指针被定义后没有分配到任何变量时，它的值为 nil。

```go
package main

import "fmt"

func main() {
   var p *string
   fmt.Println(p)
   fmt.Printf("p的值是%v\n", p)
   if p != nil {
      fmt.Println("非空")
   } else {
      fmt.Println("空值")
   }
}
```

```bash
<nil>
p的值是<nil>
空值
```

### new 和 make

```go
package main

import "fmt"

func main() {
	var a *int // 未使用 new() 分配空间
	*a = 100
	fmt.Println(*a)

	var b map[string]int // 未使用 make() 分配空间
	b["测试"] = 100
	fmt.Println(b)
}
```

执行上面的代码会引发 panic，**在Go语言中对于引用类型的变量，在使用的时候不仅要声明它，还要为它分配内存空间，否则值就没办法存储。**而对于值类型的声明不需要分配内存空间，是因为它们在声明的时候已经默认分配好了内存空间。要分配内存，就需要使用new和make。 **Go语言中new和make是内建的两个函数，主要用来分配内存。**

#### new

new 是一个内置函数，函数签名如下：

`func new(Type) *Type`

其中，Type 表示类型，new 函数只接受一个参数，这个参数是一个类型；*Type 表示类型指针，new 函数返回一个指向该类型内存地址的指针。

new 函数不太常用，使用 new 函数得到的是一个类型的指针，并且该指针对应的值为该类型的零值。

示例代码中`var a *int`只是声明了一个指针变量a但是没有初始化，**指针作为引用类型需要初始化后才会拥有内存空间**，才可以给它赋值。应该使用内置的new函数对a进行初始化之后才可以正常对其赋值。

```go
var a = new(int)
*a = 100
fmt.Println(*a)
```

#### make

make也是用于内存分配的，区别于new，**它只用于slice、map以及chan的内存创建，而且它返回的类型就是这三个类型本身，而不是他们的指针类型，**因为这三种类型就是引用类型，所以就没有必要返回他们的指针了。make函数的函数签名如下：

```go
func make(t Type, size ...IntegerType) Type
```

make函数是无可替代的，在使用slice、map以及channel的时候，都需要使用make进行初始化，然后才可以对它们进行操作。

示例代码中`var b map[string]int`只是声明变量b是一个map类型的变量，需要使用make函数进行初始化操作之后，才能对其进行键值对赋值。

```go
var b = make(map[string]int) // 使用 make() 初始化存储空间
b["测试"] = 100
fmt.Println(b)
```

#### new 与 make 的区别

1. 两者都是用来做内存分配的。
2. make 只用于slice、map以及channel的初始化，**返回的还是这三个引用类型本身。**
3. new 用于类型的内存分配，并且内存对应的值为类型零值，**返回的是指向类型的指针。**

### 指针练习

程序定义一个 int 变量 num 的地址并打印，将 num 的地址赋值给指针 ptr，并通过 ptr 去修改 num 的值。

```go
package main

import "fmt"

func main() {
	// 定义一个 int 变量 num 的地址并打印
	num := 10
	fmt.Println(&num)
	// 将 num 的地址赋值给指针 ptr
	ptr = &num 
	// 通过 ptr 去修改 num 的值
	*ptr = 222
	fmt.Println(num)
}
```

### 普通指针、unsafe.Pointer、uintptr 区别

- *类型：普通指针类型，用于传递对象地址，不能进行指针运算，最常用。
- unsafe.Pointer：通用指针类型，用于转换不同类型的指针，不能进行指针运算，不能读取内存存储的值（**必须转换到某一类型的普通指针**）。
- uintptr：用于指针运算，GC 不把 uintptr 当指针，uintptr 无法持有对象，uintptr 类型的目标会被回收。

**unsafe.Pointer 是桥梁，可以让任意类型的指针实现相互转换，也可以将任意类型的指针转换为 uintptr 进行指针运算。**
**unsafe.Pointer 不能参与指针运算**，比如要在某个指针地址上加上一个偏移量，Pointer是不能做这个运算的，而 uintptr类型可以，只需要将Pointer类型转换成uintptr类型，做完运算后，转换成Pointer，通过*操作，取值，修改值，都可以。

 **总结：`unsafe.Pointer` 可以让变量在不同的普通指针类型转来转去，也就是表示为任意可寻址的指针类型。而 `uintptr` 常与 `unsafe.Pointer` 配合，用于做指针运算。**

### unsafe.Pointer

unsafe.Pointer称为通用指针，官方文档对该类型有四个重要描述：
（1）任何类型的指针都可以被转化为Pointer
（2）Pointer可以被转化为任何类型的指针
（3）uintptr可以被转化为Pointer
（4）Pointer可以被转化为uintptr
unsafe.Pointer是特别定义的一种指针类型（译注：类似C语言中的void类型的指针），在Golang中是用于各种指针相互转换的桥梁，**它可以包含任意类型变量的地址。**
**但是不可以直接通过 *p 来获取 unsafe.Pointer 指针指向的真实变量的值，因为不知道变量的具体类型。**
**和普通指针一样，unsafe.Pointer指针也是可以比较的，并且支持和nil常量比较判断是否为空指针。**

### uintptr

uintptr是一个整数类型。**即使uintptr变量仍然有效，由uintptr变量表示的地址处的数据也可能被GC回收。**

```go
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
```

### usafe 包

unsafe包只有两个类型，三个函数（当时的 unsafe包），但是功能很强大。

```go
// 编写笔记时 unsafe 包中包含内容，三个类型，五个函数
type ArbitraryType int
type IntegerType int // 当时包中未包含
type Pointer *ArbitraryType
func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
func Add(ptr Pointer, len IntegerType) Pointer // 同上
func Slice(ptr *ArbitraryType, len IntegerType) []ArbitraryType // 同上
```

### unsafe.pointer 用于普通指针类型转换

```go
package main

import (
   "fmt"
   "reflect"
   "unsafe"
)

func main() {

   v1 := uint(12)
   v2 := int(13)

   fmt.Println(reflect.TypeOf(v1)) //uint
   fmt.Println(reflect.TypeOf(v2)) //int

   fmt.Println(reflect.TypeOf(&v1)) //*uint
   fmt.Println(reflect.TypeOf(&v2)) //*int

   p := &v1
   p = (*uint)(unsafe.Pointer(&v2)) //使用unsafe.Pointer进行类型的转换

   fmt.Println(reflect.TypeOf(p)) // *unit
   fmt.Println(*p) //13
}
```

### unsafe.pointer 用于访问操作结构体的私有变量

文件路径如下：

```lua
awesomeProject
└── test 
	├── main.go
	└── p
		└── v.go
```

```go
// v.go
package p

import (
   "fmt"
)

type V struct {
   i int32
   j int64
}

func (this V) PutI() {
   fmt.Printf("i=%d\n", this.i)
}

func (this V) PutJ() {
   fmt.Printf("j=%d\n", this.j)
}
```

```go
// main.go
package main

import (
   "awesomeProject/test/p"
   "unsafe"
)

func main() {
   var v = new(p.V)
   var i = (*int32)(unsafe.Pointer(v))
   *i = int32(98)
   var j = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + unsafe.Sizeof(int32(0))))
   *j = int64(763)
   v.PutI()
   v.PutJ()
}
```

进行修改操作时需要知道结构体V的成员布局，要修改的成员大小以及成员的偏移量。核心思想是：结构体的成员在内存中的分配是一段连续的内存，结构体中第一个成员的地址就是这个结构体的地址，可以认为是相对于这个结构体偏移了0。相同的，这个结构体中的任一成员都可以相对于这个结构体的偏移来计算出它在内存中的绝对地址。

具体来讲解下main方法的实现：

```go
var v = new(p.V)
```

new是Golang的内置方法，用来分配一段内存(会按类型的零值来清零)，并返回一个指针。所以v就是类型为p.V的一个指针。

将指针v转成通用指针，再转成int32指针。这里就看到了unsafe.Pointer的作用了，不能直接将v转成int32类型的指针，那样将会 panic。刚才说了v的地址其实就是它的第一个成员的地址，所以这个i就很显然指向了v的成员i，通过给i赋值就相当于给v.i赋值了，同时由于i只是个指针，需要赋值得解引用。

```go
var i = (*int32)(unsafe.Pointer(v))
*i = int32(98)
```

现在已经成功的改变了v的私有成员i的值，但是对于v.j来说，怎么来得到它在内存中的地址呢？可以获取它相对于v的偏移量(unsafe.Sizeof可以做这个事)。

已经知道v是有两个成员的，包括i和j，并且在定义中，i位于j的前面，而i是int32类型，也就是说i占4个字节。所以j是相对于v偏移了4个字节。可以用uintptr(4)或uintptr(unsafe.Sizeof(int32(0)))来做这个事。unsafe.Sizeof方法用来得到一个值应该占用多少个字节空间。之所以转成uintptr类型是因为需要做指针运算。v的地址加上j相对于v的偏移地址，也就得到了v.j在内存中的绝对地址，别忘了j的类型是int64，所以现在的j就是一个指向v.j的指针，接下来给它赋值：

```go
var j = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + unsafe.Sizeof(int32(0))))
*j = int64(763)
```

在p目录下新建w.go文件，代码如下：

```go
package p

import (
   "fmt"
   "unsafe"
)

type W struct {
   b byte
   i int32
   j int64
}

func init() {
   var w *W = new(W)
   fmt.Printf("size=%d\n", unsafe.Sizeof(*w))
}
```

w.go里定义了一个特殊方法init，它会在导入p包时自动执行。每个包都可定义多个init方法，**它们会在包被导入时自动执行**，在执行main方法前被执行，通常用于初始化工作，但是，最好在一个包中只定义一个init方法，否则将会很难预期它的行为。

上述代码输出结果为`size=16`，之所以是这样的结果是因为发生了对齐。在struct中，它的对齐值是它的成员中的最大对齐值。每个成员类型都有它的对齐值，可以用unsafe.Alignof()方法来计算，比如unsafe.Alignof(w.b)就可以得到b在w中的对齐值。同理，我们可以计算出w.b的对齐值是1，w.i的对齐值是4，w.j的对齐值也是4。如果您认为w.j的对齐值是8那就错了，所以前面的代码能正确执行(试想一下，如果w.j的对齐值是8，那前面的赋值代码就有问题了。也就是说前面的赋值中，如果v.j的对齐值是8，那么v.i跟v.j之间应该有4个字节的填充。所以得到正确的对齐值是很重要的)。对齐值最小是1，这是因为存储单元是以字节为单位。所以b就在w的首地址，而i的对齐值是4，它的存储地址必须是4的倍数，因此，**在b和i的中间有3个填充**，同理j也需要对齐**，但因为i和j之间不需要填充，所以w的Sizeof值应该是13+3=16**。如果要通过unsafe来对w的三个私有成员赋值，b的赋值同前，而i的赋值则需要跳过3个字节，也就是计算偏移量的时候多跳过3个字节，同理j的偏移可以通过简单的数学运算就能得到。

## Map

map是一种无序的基于key-value的数据结构，Go语言中的map是引用类型，必须初始化才能使用。

### map 定义

Go语言中 map的定义语法如下：

```go
map[KeyType]ValueType
```

其中，KeyType 表示键的类型；ValueType 表示键对应的值的类型。

map 类型的变量默认初始值为 nil，需要使用 make() 函数来分配内存。

```go
make(map[KeyType]ValueType, [cap])
```

其中cap表示map的容量，该参数虽然不是必须的，但是在初始化map的时候最好为其指定一个合适的容量。

### map 基本使用

map支持在声明的时候填充元素，例如：

```go
package main

import "fmt"

func main() {
   userInfo := map[string]string{
      "username": "pprof.cn",
      "password": "123456", // 注意此处有 ,
   }
   fmt.Println(userInfo)
}
```

### 判断某个键是否存在

Go语言中有个判断map中键是否存在的特殊写法，格式如下:

```go
value, ok := map[key]
```

ok 表示 key 的存在与否，不存在为 false。

### map 的遍历

 Go语言中使用for range遍历map，遍历map时的元素顺序与添加键值对的顺序无关。

```go
package main

import "fmt"

func main() {
   scoreMap := make(map[string]int)
   scoreMap["张三"] = 90
   scoreMap["小明"] = 100
   scoreMap["王五"] = 60
   for k, v := range scoreMap { // 只遍历 k 时可在此处去掉 v 即可，不可只遍历 v
      fmt.Println(k, v)
   }
}
```

### 使用 delete() 删除键值对

使用delete()内建函数从map中删除一组键值对，delete()函数的格式如下：

```go
func delete(m map[Type]Type1, key Type)
```

其中，map 表示要删除键值对的 map；key 表示要删除键值对的键。

```go
package main

import "fmt"

func main(){
   scoreMap := make(map[string]int)
   scoreMap["张三"] = 90
   scoreMap["小明"] = 100
   scoreMap["王五"] = 60
   delete(scoreMap, "小明") // 将 小明:100 从 map 中删除
   for k,v := range scoreMap{
      fmt.Println(k, v)
   }
}
```

### 按照指定顺序遍历 map

```go
package main

import (
   "fmt"
   "math/rand"
   "sort"
   "time"
)

func main() {
   rand.Seed(time.Now().UnixNano()) //初始化随机数种子
   var scoreMap = make(map[string]int, 200)

   for i := 0; i < 100; i++ {
      key := fmt.Sprintf("stu%02d", i) // 生成stu开头的字符串
      value := rand.Intn(100)          // 生成0~99的随机整数
      scoreMap[key] = value
   }
   // 取出map中的所有key存入切片keys
   var keys = make([]string, 0, 200)
   for key := range scoreMap {
      keys = append(keys, key)
   }
   // 对切片进行排序
   sort.Strings(keys)
   // 按照排序后的key遍历map
   for _, key := range keys {
      fmt.Println(key, scoreMap[key])
   }
}
```

### 元素为map类型的切片

```go
package main

import (
   "fmt"
)

func main() {
   // 创建 map 切片
   var mapSlice = make([]map[string]string, 3)
   for index, value := range mapSlice {
      fmt.Printf("index:%d value:%v\n", index, value)
   }
   fmt.Println("after init")
   // 对切片中的map元素进行初始化
   mapSlice[0] = make(map[string]string, 10)
   mapSlice[0]["name"] = "王五"
   mapSlice[0]["password"] = "123456"
   mapSlice[0]["address"] = "红旗大街"
   for index, value := range mapSlice {
      fmt.Printf("index:%d value:%v\n", index, value)
   }
}
```

### 值为切片类型的map

```go
package main

import (
   "fmt"
)

func main() {
   var sliceMap = make(map[string][]string, 3) // map[KeyType][表示为切片]ValueType
   fmt.Println(sliceMap)
   fmt.Println("after init")
   key := "中国"
   value, ok := sliceMap[key]
   // key 不存在
   if !ok {
      value = make([]string, 0, 2)
   }
   value = append(value, "北京", "上海")
   sliceMap[key] = value
   fmt.Println(sliceMap)
}
```

## Map 实现原理

### key，value 存储

Map 是一种通过 key 来获取 value 的一个数据结构，**其底层存储方式为数组**，在存储时key 不能重复，当 key 重复时，value 会进行覆盖，通过 key 进行 hash 运算（可以简单理解为把 key 转化为一个整形数字）然后对数组的长度取余，得到 key 存储在数组的哪个下标位置，最后将 key 和 value 组装为一个结构体，放入数组下标处，如下图所示：

![key，value 存储](https://raw.githubusercontent.com/tonshz/test/master/img/202203122116964.png "key，value 存储")

### hash 冲突

如上图所示，数组一个下标处只能存储一个元素，也就是说一个数组下标只能存储一对key，value, `hashkey(xiaoming)=4`占用了下标0的位置，假设遇到另一个key，hashkey(xiaowang) 也是4，这就是hash冲突（不同的 key 经过 hash 之后得到的值一样）。

### hash 冲突的常见解决方法

#### 开放定址法

存储一个key，value时，发现 hashkey(key) 的下标已经被别的 key 占用，那可以在这个数组空间中重新找一个没被占用的存储这个冲突的 key。常见的有线性探测法，线性补偿探测法，随机探测法，这里主要介绍以下线性探测法。

线性探测，字面意思就是按照顺序来，从冲突的下标处开始往后探测，到达数组末尾时，从数组开始处探测，直到找到一个空位置存储这个key，当整个数组遍历完都找不到的情况下会进行扩容（事实上当数组容量快满的时候就会扩容了）；查找某一个key的时候，找到key对应的下标，比较key是否相等，如果相等直接取出来，否则按照顺序探测直到碰到一个空位置，说明此处key不存在。如下图：首先存储key=xiaoming在下标0处，当存储key=xiaowang时，hash冲突了，按照线性探测，存储在下标1处，（红色的线是冲突或者下标已经被占用了） 再者key=xiaozhao存储在下标4处，当存储key=xiaoliu是，hash冲突了，按照线性探测，从头开始，存储在下标2处 （黄色的是冲突或者下标已经被占用了）

![开放定址法](https://raw.githubusercontent.com/tonshz/test/master/img/202203122124320.png "开放定址法(线性探测)")

#### 拉链法

简单理解为链表，当key发生hash冲突时，在冲突位置的元素上形成一个链表，通过指针互连接，当查找时，发现key冲突，顺着链表一直往下找，直到链表的尾节点，找不到则返回空，如下图：

![拉链法](https://raw.githubusercontent.com/tonshz/test/master/img/202203122130965.png "拉链法")

#### 开放定址（线性探测）和拉链的优缺点

- 由上面可以看出拉链法比线性探测处理简单
- 线性探测查找相比拉链法更消耗时间
- 线性探测更容易导致扩容，而拉链不会
- 拉链存储了指针，所以空间上会比线性探测占用多一点
- 拉链是动态申请存储空间的，更适合链长不确定的

### Go 中 Map 的使用

下列为 map 的创建、初始化、增删改查等操作。

```go
package main

import "fmt"

func main() {
   //直接创建初始化一个mao
   var mapInit = map[string]string {"xiaoli":"湖南", "xiaoliu":"天津"}
   fmt.Println(mapInit)
   //声明一个map类型变量,
   //map的key的类型是string，value的类型是string
   var mapTemp map[string]string
   //使用make函数初始化这个变量,并指定大小(也可以不指定)
   mapTemp = make(map[string]string,10)
   //存储key ，value
   mapTemp["xiaoming"] = "北京"
   mapTemp["xiaowang"]= "河北"
   //根据key获取value,
   //如果key存在，则ok是true，否则是flase
   //v1用来接收key对应的value,当ok是false时，v1是nil
   v1,ok := mapTemp["xiaoming"]
   fmt.Println(ok,v1)
   //当key=xiaowang存在时打印value
   if v2,ok := mapTemp["xiaowang"]; ok{
      fmt.Println(v2)
   }
   //遍历map,打印key和value
   for k,v := range mapTemp{
      fmt.Println(k,v)
   }
   //删除map中的key
   delete(mapTemp,"xiaoming")
   //获取map的大小
   l := len(mapTemp)
   fmt.Println(l)
}
```

### Go中Map的实现原理

map的源码位于 src/runtime/map.go 中map同样也是数组存储的的，每个数组下标处存储的是一个 bucket，每个bucket中可以存储 8个kv键值对，**当每个bucket存储的kv对到达8个之后，会通过 overflow 指针指向一个新的 bucket，从而形成一个链表，**事实上，kv结构体和overflow指针并没有显示定义，而是通过指针运算进行访问的。

```go
// A bucket for a Go map. bucket结构体定义
type bmap struct {
   // tophash generally contains the top byte of the hash value
   // for each key in this bucket. If tophash[0] < minTopHash,
   // tophash[0] is a bucket evacuation state instead.
   // top hash通常包含该bucket中每个键的hash值的高八位。
   // 如果tophash[0]小于mintophash，则tophash[0]为桶疏散状态，bucketCnt 的初始值是8
   tophash [bucketCnt]uint8
   // Followed by bucketCnt keys and then bucketCnt elems.
   // 接下来是bucketcnt键，然后是bucketcnt值。
   // NOTE: packing all the keys together and then all the elems together makes the
   // code a bit more complicated than alternating key/elem/key/elem/... but it allows
   // us to eliminate padding which would be needed for, e.g., map[int64]int8.
   // Followed by an overflow pointer.
   // 注意：将所有键打包在一起，然后将所有值打包在一起    
   // 使得代码比交替键/值/键/值/更复杂。但它允许我们消除可能需要的填充，    
   // 例如map[int64]int8 后面跟一个溢出指针
}
```

看上面代码以及注释，能得到bucket中存储的kv，**tophash用来快速查找key值是否在该bucket中**，从而不不需要每次都通过真值进行比较；还有kv的存放，为什么不是k1v1，k2v2….. 而是k1k2…v1v2…，上面的注释说的 map[int64]int8,key是int64（8个字节），value是int8（一个字节），kv的长度不同，如果按照kv格式存放，则考虑内存对齐v也会占用int64，而按照后者存储时，8个v刚好占用一个int64。

![bucket 结构体](https://raw.githubusercontent.com/tonshz/test/master/img/202203122144055.png "bucket 结构体")

最后分析一下go的整体内存结构，如下图所示，当往map中存储一个kv对时，通过k获取hash值，hash值的低八位和bucket数组长度取余，定位到在数组中的那个下标，hash值的高八位存储在bucket中的tophash中，用来快速判断key是否存在，key和value的具体值则通过指针运算存储，当一个bucket满时，通过overfolw指针链接到下一个bucket。

![map 内存结构](https://raw.githubusercontent.com/tonshz/test/master/img/202203122145062.png "map 内存结构")

## 结构体

Go语言中没有“类”的概念，也不支持“类”的继承等面向对象的概念。Go语言中通过结构体的内嵌再配合接口比面向对象具有更高的扩展性和灵活性。

Go语言中的基础数据类型可以表示一些事物的基本属性，但当想表达一个事物的全部或部分属性时，这时候再用单一的基本数据类型明显就无法满足需求了，Go语言提供了一种自定义数据类型，可以封装多个基本数据类型，这种数据类型叫结构体，英文名称struct。 也就是可以通过struct来定义自己的类型了。**Go语言中通过struct来实现面向对象。**

### 自定义类型

在Go语言中有一些基本的数据类型，如string、整型、浮点型、布尔等数据类型，Go语言中可以使用type关键字来定义自定义类型。

自定义类型是定义了一个全新的类型。可以基于内置的基本类型定义，也可以通过struct定义。例如：

```go
type MyInt int // 将MyInt定义为int类型
```

通过Type关键字的定义，MyInt就是一种新的类型，它具有int的特性。

### 类型别名

类型别名规定：TypeAlias只是Type的别名**，本质上TypeAlias与Type是同一个类型**。就像一个孩子小时候有小名、乳名，上学后用学名，英语老师又会给他起英文名，但这些名字都指的是他本人。

```go
type TypeAlias = Type
```

Go 中的 `rune`与`byte`就是类型别名：

```go
type byte = uint8
type rune = int32
```

### 自定义类型与类型别名的区别

```go
package main

import "fmt"
// 类型定义
type NewInt int

// 类型别名
type MyInt = int

func main() {
   var a NewInt
   var b MyInt

   fmt.Printf("type of a:%T\n", a) //type of a:main.NewInt
   fmt.Printf("type of b:%T\n", b) //type of b:int
} 
```

```bash
type of a:main.NewInt # 注意两者类型的差别
type of b:int
```

结果显示a的类型是main.NewInt，表示main包下定义的NewInt类型。b的类型是int。**MyInt类型只会在代码中存在，编译完成时并不会有MyInt类型。**

### 结构体的定义

使用type和struct关键字来定义结构体，具体代码格式如下：

```go
type TypeName struct{
    FiledName FiledType
    FiledName FiledType
    ...
}
```

其中，TypeName 表示自定义结构体的名称，在一个包内不能重复；FiledName 表示结构体字段名，结构体中字段名必须唯一；FiledType 表示结构体字段的具体类型。

例如简单定义一个 Person 结构体：

```go
type Person struct{
    name string
    city string
    age int8
}
```

相同类型的字段也可以写在同一行，如上例中的 name 与 city 可以写在同一行`name, city string`。

这样就拥有了一个person的自定义类型，它有name、city、age三个字段，分别表示姓名、城市和年龄。可以使用这个Person结构体很方便的在程序中表示和存储人信息了。

语言内置的基础数据类型是用来描述一个值的，而结构体是用来描述一组值的。比如一个人有名字、年龄和居住城市等，**本质上是一种聚合型的数据类型。**

### 结构体实例化

只有当结构体实例化时，才会真正地分配内存，也就是必须实例化后才能使用结构体的字段。结构体本身也是一种类型，可以像声明内置类型一样使用var关键字声明结构体类型。

```go
package main

import "fmt"

type Person struct{
   name, city string
   age int
}
func main() {
   p := Person{
      name: "zhangsan",
      city: "beijing",
      age:  0,
   }
   fmt.Println(p)
}
```

可以通过`.`来访问结构体的字段（成员变量），例如`p.name`等。

### 匿名结构体

在定义一些临时数据结构等场景下还可以使用匿名结构体。

```go
package main

import (
	"fmt"
)

func main() {
	// var
	var user struct{Name string; Age int}
	user.Name = "pprof.cn"
	user.Age = 18
    fmt.Printf("%#v\n", user) // %#v：先输出结构体名字值,再输出结构体(字段名字+字段的值) 

	// :=
	user1 := struct{
		name string
		city string
	}{
		name: "zhangsan",
		city: "changsha",
	} // 注意此处还有个大括号，可在此内赋值
	fmt.Printf("%#v\n", user1)
}
```

```bash
struct { Name string; Age int }{Name:"pprof.cn", Age:18}
struct { name string; city string }{name:"zhangsan", city:"changsha"}
```

### 创建指针类型结构体

可以通过使用new关键字对结构体进行实例化，得到的是结构体的地址。 格式如下：

```go
var p2 = new(person)
fmt.Printf("%T\n", p2)     // *main.person
fmt.Printf("p2=%#v\n", p2) // p2=&main.person{name:"", city:"", age:0}  
```

从打印的结果中可以看出p2是一个结构体指针。需要注意的是在Go语言中支持对结构体指针直接使用`.`来访问结构体的成员。

```go
var p2 = new(person)
p2.name = "测试"
p2.age = 18
p2.city = "北京"
fmt.Printf("p2=%#v\n", p2) //p2=&main.person{name:"测试", city:"北京", age:18} 
```

### 取结构体的地址实例化

使用&对结构体进行取地址操作相当于对该结构体类型进行了一次new实例化操作。`p3.name = “博客”`其实在底层是`(*p3).name = “博客”`，**这是Go语言实现的语法糖。**

```go
p3 := &person{} // 注意有个大括号
fmt.Printf("%T\n", p3)     // *main.person
fmt.Printf("p3=%#v\n", p3) // p3=&main.person{name:"", city:"", age:0}
p3.name = "博客"
p3.age = 30
p3.city = "成都"
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"博客", city:"成都", age:30}
```

### 使用值的列表初始化

初始化结构体的时候可以简写，也就是初始化的时候不写键，直接写值：

```go
package main

import "fmt"

type Person struct{
   name, city string
   age int
}

func main() {
   p := Person{
      "zhangsan",
      "changsha",
      11,
   }
   fmt.Printf("%#v", p // main.Person{name:"zhangsan", city:"changsha", age:11}
}
```

使用这种格式初始化时，需要注意：

+ **必须初始化结构体的所有字段。**
+ **初始值的填充顺序必须与字段在结构体中的声明顺序一致。**
+ **该方式不能和键值初始化方式混用。** 

### 构造函数

Go语言的结构体没有构造函数，但可以自己实现。 例如，下方的代码就实现了一个person的构造函数。 因为struct是值类型，如果结构体比较复杂的话，值拷贝性能开销会比较大，所以该构造函数返回的是结构体指针类型。

```go
func newPerson(name, city string, age int8) *person {
    return &person{
        name: name,
        city: city,
        age:  age,
    }
}
```

调用构造函数

```go
p := newPerson("pprof.cn", "测试", 90)
fmt.Printf("%#v\n", p9) 
```

### 方法和接收者

**Go语言中的方法（Method）是一种作用于特定类型变量的函数，这种特定类型变量叫做接收者（Receiver）。**接收者的概念就类似于其他语言中的 this 或者 self。

方法的定义格式如下：

```go
func(接收者变量 接收者类型) 方法名 (参数列表)(返回参数列表){
    函数体
}
```

其中，接收者中的参数变量名在命名时，**官方建议使用接收者类型名的第一个小写字母**，而不是self、this之类的命名，例如，Person类型的接收者变量应该命名为 p，Connector类型的接收者变量应该命名为c等；**接收者类型和参数类似，可以是指针类型和非指针类型**；方法名、参数列表、返回参数，具体格式与函数定义相同。

方法与函数的区别是，**函数不属于任何类型，方法属于特定的类型。**

```go
package main

import "fmt"

//Person 结构体
type Person struct {
   name string
   age  int8
}

//NewPerson 构造函数
func NewPerson(name string, age int8) *Person {
   return &Person{
      name: name,
      age:  age,
   }
}

// Dream Person做梦的方法
func (p Person) Dream() {
   fmt.Printf("%s的梦想是学好Go语言！\n", p.name)
}

func main() {
   p1 := NewPerson("测试", 25)
   p1.Dream() // p1可以调用 Dream()
}
```

### 指针类型的接收者

指针类型的接收者由一个结构体的指针组成，由于指针的特性，**调用方法时修改接收者指针的任意成员变量**，在方法结束后，修改都是有效的。这种方式就十分接近于其他语言中面向对象中的this或者self。 如下为Person添加一个SetAge方法，来修改实例变量的年龄。

```go
package main

import "fmt"

//Person 结构体
type Person struct {
   name string
   age  int8
}

//NewPerson 构造函数
func NewPerson(name string, age int8) *Person {
   return &Person{
      name: name,
      age:  age,
   }
}

//Dream Person做梦的方法
func (p Person) Dream() {
   fmt.Printf("%s的梦想是学好Go语言！\n", p.name)
}

// SetAge 设置p的年龄
// 使用指针接收者：*Person
func (p *Person) SetAge(newAge int8) {
   p.age = newAge
}

func main() {
   p1 := NewPerson("测试", 25)
   p1.Dream()
   fmt.Println(p1.age) // 25
   p1.SetAge(30)
   fmt.Println(p1.age) // 30
}
```

使用指针类型接收者的情况：

+ 需要修改接受者中的值。
+ 接受这是拷贝代价比较大的对象。
+ **保证一致性，如果有某个方法使用了指针接收者，那么其他的方法也应该使用指针接收者。**

### 值类型的接收者

当方法作用于值类型接收者时，Go语言会在代码运行时将接收者的值复制一份。在值类型接收者的方法中可以获取接收者的成员值，**但修改操作只是针对副本，无法修改接收者变量本身。**

### 任意类型添加方法

在Go语言中，接收者的类型可以是任何类型，不仅仅是结构体，任何类型都可以拥有方法。 举个例子，基于内置的int类型使用type关键字可以定义新的自定义类型，然后为自定义类型添加方法。**非本地类型不能定义方法，也就是说不能给别的包的类型定义方法。**

```go
package main

import "fmt"

//MyInt 将int定义为自定义MyInt类型
type MyInt int

//SayHello 为MyInt添加一个SayHello的方法
func (m MyInt) SayHello() {
   fmt.Println("Hello, 我是一个int。")
}
func main() {
   var m1 MyInt
   m1.SayHello() //Hello, 我是一个int。
   m1 = 100
   fmt.Printf("%#v  %T\n", m1, m1) //100  main.MyInt
} 
```

### 结构体的匿名字段

结构体允许其成员字段在声明时没有字段名而只有类型，这种没有名字的字段就称为匿名字段。**匿名字段默认采用类型名作为字段名**，结构体要求字段名称必须唯一，**因此一个结构体中同种类型的匿名字段只能有一个。**

```go
package main

import "fmt"

//Person 结构体Person类型
type Person struct {
   string // 不能再出现第二个 string 类型
   int
}

func main() {
   p1 := Person{
      "pprof.cn",
      18,
   }
   fmt.Printf("%#v\n", p1)        //main.Person{string:"pprof.cn", int:18}
   fmt.Println(p1.string, p1.int) //pprof.cn 18
}
```

### 嵌套结构体

一个结构体中可以嵌套包含另一个结构体或结构体指针，**注意内部嵌套结构体的初始化。**

```go
package main

import "fmt"

//Address 地址结构体
type Address struct {
   Province string
   City     string
}

//User 用户结构体
type User struct {
   Name    string
   Gender  string
   Address Address
}

func main() {
   user1 := User{
      Name:   "pprof",
      Gender: "女",
      Address: Address{
         Province: "黑龙江",
         City:     "哈尔滨",
      }, // 注意此处的初始化
   }
   fmt.Printf("user1=%#v\n", user1) // user1=main.User{Name:"pprof", Gender:"女", Address:main.Address{Province:"黑龙江", City:"哈尔滨"}}
}
```

### 嵌套匿名结构体

**当访问结构体成员时会先在结构体中查找该字段，找不到再去匿名结构体中查找。**

```go
package main

import "fmt"

//Address 地址结构体
type Address struct {
   Province string
   City     string
}

//User 用户结构体
type User struct {
   Name    string
   Gender  string
   Address //匿名结构体
}

func main() {
   var user2 User
   user2.Name = "pprof"
   user2.Gender = "女"
   user2.Address.Province = "黑龙江"    // 通过匿名结构体.字段名访问
   user2.City = "哈尔滨"                // 直接访问匿名结构体的字段名
   fmt.Printf("user2=%#v\n", user2) //user2=main.User{Name:"pprof", Gender:"女", Address:main.Address{Province:"黑龙江", City:"哈尔滨"}}
} 
```

### 嵌套结构体的字段名冲突

嵌套结构体内部可能存在相同的字段名，**这个时候为了避免歧义需要指定具体的内嵌结构体的字段。**

```go
package main

//Address 地址结构体
type Address struct {
   Province   string
   City       string
   CreateTime string
}

//Email 邮箱结构体
type Email struct {
   Account    string
   CreateTime string // 同名字段
}

//User 用户结构体
type User struct {
   Name   string
   Gender string
   Address
   Email
}

func main() {
   var user3 User
   user3.Name = "pprof"
   user3.Gender = "女"
   // user3.CreateTime = "2019" //ambiguous selector user3.CreateTime
   user3.Address.CreateTime = "2000" // 指定Address结构体中的CreateTime
   user3.Email.CreateTime = "2000"   // 指定Email结构体中的CreateTime
}
```

### 结构体的“继承”

Go语言中使用结构体也可以实现其他编程语言中面向对象的继承。

```go
package main

import "fmt"

//Animal 动物
type Animal struct {
   name string
}

func (a *Animal) move() {
   fmt.Printf("%s会动！\n", a.name)
}

//Dog 狗
type Dog struct {
   Feet    int8
   *Animal //通过嵌套匿名结构体实现继承
}

func (d *Dog) wang() {
   fmt.Printf("%s会汪汪汪~\n", d.name)
}

func main() {
   d1 := &Dog{
      Feet: 4,
      Animal: &Animal{ //注意嵌套的是结构体指针
         name: "乐乐",
      },
   }
   d1.wang() //乐乐会汪汪汪~
   d1.move() //乐乐会动！
}
```

### 结构体字段的可见性

结构体中字段大写开头表示可公开访问，小写表示私有（仅在定义当前结构体的包中可访问）。

### 结构体与 JSON 序列化

JSON 是一种轻量级的数据交换格式。易于人阅读和编写。同时也易于机器解析和生成。JSON键值对是用来保存JS对象的一种方式，键/值对组合中的键名写在前面并用双引号””包裹，使用冒号:分隔，然后紧接着值；多个键值之间使用英文,分隔。Go 中 json 序列化可使用`json.Marshal(v interface{})`，反序列化可使用`json.Unmarshal(data []byte, v interface{})`。

```go
package main

import (
   "encoding/json"
   "fmt"
)

//Student 学生
type Student struct {
   ID     int
   Gender string
   Name   string
}

//Class 班级
type Class struct {
   Title    string
   Students []*Student
}

func main() {
   c := &Class{
      Title:    "101",
      Students: make([]*Student, 0, 200),
   }
   for i := 0; i < 10; i++ {
      stu := &Student{
         Name:   fmt.Sprintf("stu%02d", i),
         Gender: "男",
         ID:     i,
      }
      c.Students = append(c.Students, stu)
   }
   //JSON序列化：结构体-->JSON格式的字符串
   data, err := json.Marshal(c)
   if err != nil {
      fmt.Println("json marshal failed")
      return
   }
   fmt.Printf("json:%s\n", data)
   //JSON反序列化：JSON格式的字符串-->结构体
   str := `{"Title":"101","Students":[{"ID":0,"Gender":"男","Name":"stu00"},{"ID":1,"Gender":"男","Name":"stu01"},{"ID":2,"Gender":"男","Name":"stu02"},{"ID":3,"Gender":"男","Name":"stu03"},{"ID":4,"Gender":"男","Name":"stu04"},{"ID":5,"Gender":"男","Name":"stu05"},{"ID":6,"Gender":"男","Name":"stu06"},{"ID":7,"Gender":"男","Name":"stu07"},{"ID":8,"Gender":"男","Name":"stu08"},{"ID":9,"Gender":"男","Name":"stu09"}]}`
   c1 := &Class{}
   err = json.Unmarshal([]byte(str), c1)
   if err != nil {
      fmt.Println("json unmarshal failed!")
      return
   }
   fmt.Printf("%#v\n", c1)
} 
```

### 结构体标签（Tag）

Tag是结构体的元信息，可以在运行的时候通过反射的机制读取出来。Tag在结构体字段的后方定义，由一对反引号包裹起来，具体的格式如下：

```go
`key1:"value1" key2:"value2"`  
```

结构体标签由一个或多个键值对组成。键与值使用冒号分隔，**值用双引号括起来**。**键值对之间使用一个空格分隔**。 结构体编写Tag时，必须严格遵守键值对的规则。结构体标签的解析代码的容错能力很差，一旦格式写错，**编译和运行时都不会提示任何错误，通过反射也无法正确取值。例如不要在key和value之间添加空格。**

例如为Student结构体的每个字段定义json序列化时使用的Tag：

```go
package main

import (
   "encoding/json"
   "fmt"
)

//Student 学生
type Student struct {
   ID     int    `json:"id"` // 通过指定tag实现json序列化该字段时的key
   Gender string // json 序列化是默认使用字段名作为key
   name   string // name 首字母小写表示私有不能被json包访问
}

func main() {
   s1 := Student{
      ID:     1,
      Gender: "女",
      name:   "pprof",
   }
   data, err := json.Marshal(s1) // json 序列化
   if err != nil {
      fmt.Println("json marshal failed!")
      return
   }
   fmt.Printf("json str:%s\n", data) //json str:{"id":1,"Gender":"女"}
} 
```

### 小练习

```go
package main

import "fmt"

type student struct {
    id   int
    name string
    age  int
}

func demo(ce []student) {
    // 切片是引用传递，是可以改变值的
    ce[1].age = 999
    // ce = append(ce, student{3, "xiaowang", 56})
    // return ce
}
func main() {
    var ce []student  //定义一个切片类型的结构体
    ce = []student{
        student{1, "xiaoming", 22},
        student{2, "xiaozhang", 33},
    }
    fmt.Println(ce)
    demo(ce)
    fmt.Println(ce)
}
```

代码运行结果：

```bash
[{1 xiaoming 22} {2 xiaozhang 33}]
[{1 xiaoming 22} {2 xiaozhang 999}]
```

### 删除 map 类型中的结构体

```go
package main

import "fmt"

type student struct {
    id   int
    name string
    age  int
}

func main() {
    ce := make(map[int]student)
    ce[1] = student{1, "xiaolizi", 22}
    ce[2] = student{2, "wang", 23}
    fmt.Println(ce)
    delete(ce, 2)
    fmt.Println(ce)
}
```

### 实现 map 有序输出

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    map1 := make(map[int]string, 5)
    map1[1] = "www.topgoer.com"
    map1[2] = "rpc.topgoer.com"
    map1[5] = "ceshi"
    map1[3] = "xiaohong"
    map1[4] = "xiaohuang"
    // map 直接输出顺序不定
    sli := []int{}
    for k, _ := range map1 {
        sli = append(sli, k)
    }
    sort.Ints(sli) // 将 key 进行排序
    for i := 0; i < len(map1); i++ {
        fmt.Println(map1[sli[i]]) // 按 key 的顺序输出
    }
}
```

### 小案例

[采用切片类型的结构体接受查询数据库信息返回的参数](https://github.com/lu569368/struct)

---------------------------

[参考文章](https://www.topgoer.cn/docs/golang/chapter02)

