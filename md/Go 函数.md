# Go 函数

[toc]

## 函数定义

### Golang 函数特点：

+ 无需声明原型
+ 支持不定变参
+ 支持多返回值
+ 支持命名返回参数
+ 支持匿名函数和闭包
+ 函数也是一种类型，函数可以赋值给变量
+ 不支持嵌套，一个包中不能有两个名字一样的函数
+ 不支持重载，即函数名相同，参数列表不同
+ 不支持默认参数

### 函数声明

函数声明包含一个函数名，参数列表， 返回值列表和函数体。如果函数没有返回值，则返回列表可以省略。函数从第一条语句开始执行，直到执行return语句或者执行函数的最后一条语句。

函数可以没有参数或接受多个参数。注意类型在变量名之后 。当两个或多个连续的函数命名参数是同一类型，则除了最后一个类型之外，其他都可以省略。函数可以返回任意数量的返回值。使用关键字 func 定义函数，左大括号依旧不能另起一行。

```go
func test(x, y int, s string) (int, string) {
    // 类型相同的相邻参数，参数类型可合并。 多返回值必须用括号。
    n := x + y          
    return n, fmt.Sprintf(s, n)
}
```

**函数是第一类对象，可作为参数传递。建议将复杂签名定义为函数类型，以便于阅读。**

```go
package main

import "fmt"

func test(fn func() int) int {
    return fn()
}

// 定义函数类型。
type FormatFunc func(s string, x, y int) string 

func format(fn FormatFunc, s string, x, y int) string {
    return fn(s, x, y)
}

func main() {
    s1 := test(func() int { return 100 }) // 直接将匿名函数当参数。

    s2 := format(func(s string, x, y int) string {
        return fmt.Sprintf(s, x, y)
    }, "%d, %d", 10, 20)

    println(s1, s2)
}
```

有返回值的函数，必须有明确的终止语句，否则会引发编译错误。可能会偶尔遇到没有函数体的函数声明，这表示该函数不是以Go实现的。这样的声明定义了函数标识符。

## 函数参数

函数定义时有参数，该变量可称为函数的形参，形参就像是定义在函数体内的局部变量。

但当调用函数，传递过来的变量就是函数的实参，函数可以通过两种方式来传递参数。在默认情况下，Go 语言使用的是值传递，即在调用过程中不会影响到实际参数。

注意1：无论是值传递，还是引用传递，传递给函数的都是变量的副本，不过，**值传递是值的拷贝，引用传递是地址的拷贝**，一般来说，地址拷贝更为高效。而值拷贝取决于拷贝的对象大小，对象越大，则性能越低。

注意2：**map、slice、chan、指针、interface默认以引用的方式传递。**

### 值传递

指在调用函数时将实际参数复制一份传递到函数中，这样在函数中如果对参数进行修改，**将不会影响到实际参数。**

```go
func swap(x, y int) int{
	...
}
```

### 引用传递

是指在调用函数时将实际参数的地址传递到函数中，那么在函数中对参数所进行的修改，**将影响到实际参数。**

```go
package main

import (
    "fmt"
)

/* 定义相互交换值的函数 */
func swap(x, y *int) {
    var temp int

    temp = *x /* 保存 x 的值 */
    *x = *y   /* 将 y 值赋给 x */
    *y = temp /* 将 temp 值赋给 y*/

}

func main() {
    var a, b int = 1, 2
    /*
        调用 swap() 函数
        &a 指向 a 指针，a 变量的地址
        &b 指向 b 指针，b 变量的地址
    */
    swap(&a, &b)

    fmt.Println(a, b)
}
```

### 不定参数传值

即函数的参数数量不是固定的，类型是固定的（可变参数）。

Golang 可变参数本质上就是一个 slice，只能有一个，且必须是最后一个。

在参数赋值时可以不用用一个一个的赋值，**可以直接传递一个数组或者切片，特别注意的是在参数后加上`…`即可。**

```go
func myfunc(args ...int) {    //0个或多个参数
}

func add(a int, args…int) int {    //1个或多个参数
}

func add(a int, b int, args…int) int {    //2个或多个参数
}
```

注意：其中args是一个slice，可以通过arg[index]依次访问所有参数,通过len(arg)来判断传递参数的个数.

### 任意类型的不定参数

函数的参数和每个参数的类型都不是固定的。

用 interface{} 传递任意类型数据是 Go 语言的管理用法，而且 interface{} 是类型安全的。

```go
func myfunc(arg ...interface{}){
	...
}
```

```go
package main

import (
    "fmt"
)

func test(s string, n ...int) string {
    var x int
    for _, i := range n {
        x += i
    }

    return fmt.Sprintf(s, x)
}

func main() {
    println(test("sum: %d", 1, 2, 3)) // 输出：sum: 6
}
```

使用 slice 对象做变参时，必须展开，`slice...`

```go
package main

import (
    "fmt"
)

func test(s string, n ...int) string {
    var x int
    for _, i := range n {
        x += i
    }

    return fmt.Sprintf(s, x)
}

func main() {
    s := []int{1, 2, 3}
    res := test("sum: %d", s...)    // slice... 展开slice
    println(res)
}
```

## 返回值

`_`标识符，用来忽略函数的某个返回值，Go 中的返回值可以被命名，并且可以像在函数体开头声明的变量那样使用。返回值的名称应当具有一定的意义，可以作为文档使用。没有参数的 return 语句返回各个返回变量的当前值，这种用法被称作“裸”返回。直接返回语句仅应用在短函数中，在较长的函数中会影响代码的可读性。

```go
package main

import (
    "fmt"
)

func add(a, b int) (c int) {
    c = a + b
    return
}

func calc(a, b int) (sum int, avg int) {
    sum = a + b
    avg = (a + b) / 2

    return
}

func main() {
    var a, b int = 1, 2
    c := add(a, b)
    sum, avg := calc(a, b)
    fmt.Println(a, b, c, sum, avg)
}
```

**Golang返回值不能用容器对象接收多返回值。只能用多个变量，或 `"_"` 忽略。**

```go
package main

func test() (int, int) {
    return 1, 2
}

func main() {
    // s := make([]int, 2)
    // s = test()   // Error: multiple-value test() in single-value context

    x, _ := test()
    println(x)
}
```

多返回值可直接作为其他函数调用实参。

```go
package main

func test() (int, int) {
    return 1, 2
}

func add(x, y int) int {
    return x + y
}

func sum(n ...int) int {
    var x int
    for _, i := range n {
        x += i
    }

    return x
}

func main() {
    println(add(test()))
    println(sum(test()))
}
```

命名返回参数可看做与形参类似的局部变量，最后由 return 隐式返回。

```go
package main

func add(x, y int) (z int) {
    z = x + y
    return
}

func main() {
    println(add(1, 2))
}
```

命名返回参数可被同名局部变量遮蔽，此时需要显式返回。

```go
func add(x, y int) (z int) {
    // 函数体内部块，在大括号外不能声明 z 变量，Go 中一个大括号包含一个作用域
	// var z = 1 // 'z' redeclared in this block
    { // 不能在一个级别，引发 "z redeclared in this block" 错误。
        var z = x + y
        // return   // Error: z is shadowed during return
        return z // 必须显式返回。
    }
}
```

命名返回参数允许 defer 延迟调用通过闭包读取和修改。

```go
package main

func add(x, y int) (z int) {
    defer func() {
        z += 100
    }()

    z = x + y
    return
}

func main() {
    println(add(1, 2)) // 输出：103
}
```

显式 return 返回前，会先修改命名返回参数。

```go
package main

func add(x, y int) (z int) {
    defer func() {
        println(z) // 输出: 203
    }()

    z = x + y
    return z + 200 // 执行顺序: (z = z + 200) -> (call defer) -> (return)
}

func main() {
    println(add(1, 2)) // 输出: 203
}
```

## 匿名函数

匿名函数是指不需要定义函数名的一种函数实现方式。1958年LISP首先采用匿名函数。在Go里面，函数可以像普通变量一样被传递或使用，Go语言支持随时在代码里定义匿名函数。

匿名函数由一个不带函数名的函数声明和函数体组成。匿名函数的优越性在于可以直接使用函数内的变量，不必申明。

```go
package main

import (
    "fmt"
    "math"
)

func main() {
    getSqrt := func(a float64) float64 {
        return math.Sqrt(a)
    }
    fmt.Println(getSqrt(4))
}
```

上面先定义了一个名为getSqrt 的变量，初始化该变量时和之前的变量初始化有些不同，使用了func，func是定义函数的，**可是这个函数和前面说的函数最大不同就是没有函数名，也就是匿名函数。**这里将一个函数当做一个变量一样的操作。**Golang匿名函数可赋值给变量，做为结构字段，或者在 channel 里传送。**

```go
package main

func main() {
    // --- function variable ---
    fn := func() { println("Hello, World!") }
    fn()

    // --- function collection ---
    fns := [](func(x int) int){
        func(x int) int { return x + 1 },
        func(x int) int { return x + 2 },
    }
    println(fns[0](100))

    // --- function as field ---
    d := struct {
        fn func() string
    }{
        fn: func() string { return "Hello, World!" },
    }
    println(d.fn())

    // --- channel of function ---
    fc := make(chan func() string, 2)
    fc <- func() string { return "Hello, World!" }
    println((<-fc)())
}
```

```bash
Hello, World!
101
Hello, World!
Hello, World!
```

## 闭包与递归

### 闭包详解

闭包是由函数及其相关引用环境组合而成的实体(即：闭包=函数+引用环境)。

**“官方”的解释是：所谓“闭包”，指的是一个拥有许多变量和绑定了这些变量的环境的表达式（通常是一个函数），因而这些变量也是该表达式的一部分。**

维基百科讲，闭包（Closure），是引用了自由变量的函数。这个被引用的自由变量将和这个函数一同存在，即使已经离开了创造它的环境也不例外。所以，有另一种说法认为闭包是由函数和与其相关的引用环境组合而成的实体。闭包在运行时可以有多个实例，不同的引用环境和相同的函数组合可以产生不同的实例。

看着上面的描述，会发现闭包和匿名函数似乎有些像。可是可能还是有些云里雾里的。因为跳过闭包的创建过程直接理解闭包的定义是非常困难的。目前在JavaScript、Go、PHP、Scala、Scheme、Common Lisp、Smalltalk、Groovy、Ruby、 Python、Lua、objective c、Swift 以及Java8以上等语言中都能找到对闭包不同程度的支持。通过支持闭包的语法可以发现一个特点，**他们都有垃圾回收(GC)机制。**

```go
package main

import (
    "fmt"
)

// 当函数a()的内部函数b()被函数a()外的一个变量引用的时候，就创建了一个闭包。
func a() func() int {
    i := 0 // a 返回后 i 始终存在
    b := func() int {
        i++
        fmt.Println(i)
        return i
    }
    return b
}

func main() {
    c := a() // c 实际指向了函数 b()
    c()
    c()
    c()

    a() //不会输出i
}
```

```bash
1
2
3
```

闭包复制的是原对象指针，这就很容易解释延迟引用现象。

外部引用函数参数局部变量。

```go
package main

import "fmt"

// 外部引用函数参数局部变量
func add(base int) func(int) int {
   return func(i int) int {
      base += i
      return base
   }
}

func main() {
   tmp1 := add(10)
   fmt.Println(tmp1(1), tmp1(2)) // 同一个 base
   // 此时tmp1和tmp2不是一个实体了
   tmp2 := add(100)
   fmt.Println(tmp2(1), tmp2(2)) // 同一个 base
}
```

```bash
11 13
101 103
```

```go
package main

import "fmt"

// 返回2个函数类型的返回值
func test01(base int) (func(int) int, func(int) int) {
   // 定义2个函数，并返回
   // 相加
   add := func(i int) int {
      base += i
      return base
   }
   // 相减
   sub := func(i int) int {
      base -= i
      return base
   }
   // 返回
   return add, sub
}

func main() {
   f1, f2 := test01(10)
   // base 10
   fmt.Println(f1(1), f2(2)) // base+1-2 = 9
   // base 9
   fmt.Println(f1(3), f2(4)) // base+3-4 = 8
}
```

```bash
11 9
12 8
```

### 递归函数

递归，就是在运行的过程中调用自己。一个函数调用自己，就叫做递归函数。

构成递归需具备的条件：

+ 子问题须与原始问题为同样的事，且更为简单。
+ 不能无限制地调用本身，须有个出口，化简为非递归状况处理。

#### 数组阶乘

```go
package main

import "fmt"

func factorial(i int) int {
    if i <= 1 {
        return 1
    }
    return i * factorial(i-1)
}

func main() {
    var i int = 7
    fmt.Printf("Factorial of %d is %d\n", i, factorial(i))
}
```

#### 斐波那契数列(Fibonacci)

```go
package main

import "fmt"

func fibonacci(i int) int{
	if i == 1{
		return 1
	}
	if i == 2{
		return 1
	}

	return fibonacci(i-1)+fibonacci(i-2)
}

func main() {
	fmt.Println(fibonacci(5)) // 1 1 2 3 5 输出：5
}
```

## Golang 延迟调用

### defer 特性：

+ 关键字 defer 用于注册延迟调用
+ 这些调用知道 return 前才被执行，因此，defer 可以用来做资源清理
+ 多个 defer 语句，按先进后出的方式执行
+ defer 语句中的变量，在 defer 声明时就决定了

### defer 的用途

+ 关闭文件句柄
+ 锁资源释放
+ 数据库连接释放

### defer

go 语言的defer功能强大，对于资源管理非常方便，但是如果没用好，也会有陷阱。d**efer 是先进后出，后面的语句会依赖前面的资源，如果先前面的资源先释放了，后面的语句就没法执行了。**

```go
package main

import "fmt"

func main() {
    var whatever [5]struct{}
    for i := range whatever {
        defer func() { fmt.Println(i) }()
    }
} 
```

```bash
4 # 注意输出都为4
4
4
4
4
```

函数正常执行，由于闭包用到的变量 i 在执行的时候已经变成4，所以输出全都是4。

### defer f.Close

```go
package main

import "fmt"

type Test struct {
    name string
}

func (t *Test) Close() {
    fmt.Println(t.name, " closed")
}
func main() {
    ts := []Test{{"a"}, {"b"}, {"c"}}
    for _, t := range ts {
        defer t.Close() // Possible resource leak, 'defer' is called in the 'for' loop 
    }
} 
```

```bash
c  closed
c  closed
c  closed
```

输出的不是预期的输出c b a,而是输出c c c。这是由于执行的时候 t 的值已经变成了 c，与上方的原因相同。

**defer后面的语句在执行的时候，函数调用的参数会被保存起来，但是不执行，也就是复制了一份。**但是并没有说struct这里的this指针如何处理，通过这个例子可以看出go语言并没有把这个明确写出来的this指针当作参数来看待。

### 多个 defer

多个 defer 注册，按 FILO 次序执行 ( 先进后出 )。哪怕函数或某个延迟调用发生错误，这些调用依旧会被执行。

```go
package main

func test(x int) {
    defer println("a")
    defer println("b")

    defer func() {
        println(100 / x) // div0 异常未被捕获，逐步往外传递，最终终止进程。
    }()

    defer println("c")
}

func main() {
    test(0)
} 
```

```bash
c
b
a
panic: runtime error: integer divide by zero
```

### defer 延迟读取

延迟调用参数在注册时求值或复制，**可用指针或闭包 “延迟” 读取。**

```go
package main

func test() {
    x, y := 10, 20

    defer func(i int) {
        println("defer:", i, y) // y 闭包引用
    }(x) // x 被复制, 此后 x 的值被修改不影响函数执行

    x += 10
    y += 100
    println("x =", x, "y =", y)
}

func main() {
    test()
}  
```

```bash
x = 20 y = 120
defer: 10 120 # 注意此处 x 的值为10
```

### 滥用 defer

```go
package main

import (
   "fmt"
   "sync"
   "time"
)

var lock sync.Mutex

func test() {
   lock.Lock()
   lock.Unlock()
}

func testdefer() {
   lock.Lock()
   defer lock.Unlock()
}

func main() {
   func() {
      t1 := time.Now()

      for i := 0; i < 100000; i++ {
         test()
      }
      elapsed := time.Since(t1)
      fmt.Println("test elapsed: ", elapsed)
   }()
   func() {
      t1 := time.Now()

      for i := 0; i < 100000; i++ {
         testdefer()
      }
      elapsed := time.Since(t1)
      fmt.Println("testdefer elapsed: ", elapsed)
   }()

}
```

```bash
test elapsed:  553.2µs
testdefer elapsed:  1.6264ms
```

### defer 陷阱

#### defer 与 closure（闭包）

```go
package main

import (
   "errors"
   "fmt"
)

func foo(a, b int) (i int, err error) {
   defer fmt.Printf("first defer err %v\n", err) // 此处方法与2相同，都实时向函数传递参数
   defer func(err error) { fmt.Printf("second defer err %v\n", err) }(err)
   defer func() { fmt.Printf("third defer err %v\n", err) }()  // 使用了闭包
   if b == 0 {
      err = errors.New("divided by zero!")
      return
   }

   i = a / b
   return
}

func main() {
   foo(2, 0)
}
```

```bash
third defer err divided by zero! # 只有此处获取了最新的 err 值，其余均为空
second defer err <nil>
first defer err <nil>
```

如果 defer 后面跟的不是一个闭包（closure）最后执行的时候我们得到的并不是最新的值。

#### defer 与 return

```go
package main

import "fmt"

func foo() (i int) {

    i = 0
    defer func() {
        fmt.Println(i) // 2
    }()

    return 2
}

func main() {
    foo()
}
```

在有具名返回值的函数中（这里具名返回值为 i），**执行 return 2 的时候实际上已经将 i 的值重新赋值为 2。**所以defer closure 输出结果为 2 而不是 1。

#### defer nil 函数

```go
package main

import (
    "fmt"
)

func test() {
    var run func() = nil
    defer run()
    fmt.Println("runs")
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()
    test()
} 
```

```bash
runs
runtime error: invalid memory address or nil pointer dereference
```

名为 test 的函数一直运行至结束，然后 defer 函数会被执行且会因为值为 nil 而产生 panic 异常。然而值得注意的是，run() 的声明是没有问题，因为在test函数运行完成后它才会被调用。

#### 在错误位置使用 defer 

```go
package main

import "net/http"

func do() error {
    res, err := http.Get("http://www.google.com") // 当 http.Get 失败时会抛出异常。
    defer res.Body.Close()
    if err != nil {
        return err
    }

    // ..code...

    return nil
}

func main() {
    do()
} 
```

```bash
panic: runtime error: invalid memory address or nil pointer dereference
```

此处没有检查请求是否成功执行，当请求失败后，访问 Body 中的空变量 res，会抛出异常。故可以在一次成功的资源分配下面使用 defer ，对于这种情况来说意味着：当且仅当 http.Get 成功执行时才使用 defer。

```go
package main

import "net/http"

func do() error {
    res, err := http.Get("http://xxxxxxxxxx")
    if res != nil {
        defer res.Body.Close()
    }

    if err != nil {
        return err
    }

    // ..code...

    return nil
}

func main() {
    do()
} 
```

#### 不检查错误

在这里，f.Close() 可能会返回一个错误，但这个错误会被忽略掉。

```go
package main

import "os"

func do() error {
    f, err := os.Open("book.txt")
    if err != nil {
        return err
    }

    if f != nil {
        defer f.Close()
    }

    // ..code...

    return nil
}

func main() {
    do()
}  
```

```go
// 解决方案
package main

import "os"

func do() (err error) {
    f, err := os.Open("book.txt")
    if err != nil {
        return err
    }

    if f != nil {
        defer func() {
            if ferr := f.Close(); ferr != nil {
                err = ferr // 通过命名的返回变量来返回 defer 内的错误。
            }
        }()
    }

    // ..code...

    return nil
}

func main() {
    do()
} 
```

#### 释放相同资源

使用相同变量释放不同的资源，会导致操作无法正常执行。当延迟函数执行时，只有最后一个变量会被用到，因此，f 变量会成为最后那个资源 (another-book.txt)。而且两个 defer 都会将这个资源作为最后的资源来关闭。

```go
package main

import (
    "fmt"
    "os"
)

func do() error {
    f, err := os.Open("book.txt")
    if err != nil {
        return err
    }
    if f != nil {
        defer func() {
            if err := f.Close(); err != nil {
                fmt.Printf("defer close book.txt err %v\n", err) // 此处报错
            }
        }()
    }

    // ..code...

    f, err = os.Open("another-book.txt") // 与第一个变量命名相同，都为 f
    if err != nil {
        return err
    }
    if f != nil {
        defer func() {
            if err := f.Close(); err != nil {
                fmt.Printf("defer close another-book.txt err %v\n", err)
            }
        }()
    }

    return nil
}

func main() {
    do()
} 
```

```bash
defer close book.txt err close test/p/w.go: file already closed
```

```go
// 解决方案
package main

import (
    "fmt"
    "io"
    "os"
)

func do() error {
    f, err := os.Open("book.txt")
    if err != nil {
        return err
    }
    if f != nil {
        defer func(f io.Closer) {
            if err := f.Close(); err != nil {
                fmt.Printf("defer close book.txt err %v\n", err)
            }
        }(f) // 将 f 作为参数传入 或者修改两者变量名为不相同
    }

    // ..code...

    f, err = os.Open("another-book.txt")
    if err != nil {
        return err
    }
    if f != nil {
        defer func(f io.Closer) {
            if err := f.Close(); err != nil {
                fmt.Printf("defer close another-book.txt err %v\n", err)
            }
        }(f)
    }

    return nil
}

func main() {
    do()
} 
```

## 异常处理

Golang 没有结构化异常，使用 panic 抛出错误，recover 捕获错误。异常的使用场景简单描述：Go中可以抛出一个panic的异常，然后在defer中通过recover捕获这个异常，然后正常处理。

panic:

1. 内置函数
2. 假如函数F中书写了panic语句，会终止其后要执行的代码，在panic所在函数F内如果存在要执行的defer函数列表，按照defer的逆序执行



























