## Go 易错知识点

### 数据定义

#### 1. 函数返回值问题

> 代码是否可以通过编译？

```go
package main

/*
    代码是否编译通过?
*/
func myFunc(x,y int)(sum int,error){
    return x+y,nil
}

func main() {
    num, err := myFunc(1, 2)
    fmt.Println("num = ", num)
}
```

> 不能，首先 main() 中 err 声明了但未使用会报错。再者，**在函数有多个返回值时，只要有一个返回值有指定命名，其他的也必须命名**。 如果有多个返回值必须加上括号；如果只有一个返回值且有命名也需要加上括号； 此处函数第一个返回值有sum名称，第二个未命名，故编译出错。

#### 2. 结构比较问题

> 代码是否可以编译通过？为什么？

```go
package main

import "fmt"

func main() {

	sn1 := struct {
		age  int
		name string
	}{age: 11, name: "qq"}

	sn2 := struct {
		age  int
		name string
	}{age: 11, name: "qq"}

	if sn1 == sn2 {
		fmt.Println("sn1 == sn2") // true
	}

	sm1 := struct {
		age int
		m   map[string]string
	}{age: 11, m: map[string]string{"a": "1"}}

	sm2 := struct {
		age int
		m   map[string]string
	}{age: 11, m: map[string]string{"a": "1"}}

	if sm1 == sm2 {
        fmt.Println("sm1 == sm2") // error 
	}
}
```

> 编译不通过，报错  `invalid operation: sm1 == sm2 (struct containing map[string]string cannot be compared)`
>
> 结构体比较规则注意事项：
>
> 1. 只有相同类型的结构体才可以比较，结构体是否相同不但与属性类型个数有关，**还与属性顺序有关**
>
>    ```go
>    sn1 := struct {
>    	age  int
>    	name string
>    }{age: 11, name: "qq"}
>    
>    sn3:= struct {
>        name string
>        age  int
>    }{age:11, name:"qq"}
>    
>    if sn1 == sn3 {
>        fmt.Println("sn1 == sn3") // error 
>    }
>    ```
>
>    上例代码中的 sn3 与 sn1 就不是相同的结构体，不能进行比较，报错 `invalid operation: sn1 == sn3 (mismatched types struct { age int; name string } and struct { name string; age int })`
>
> 2. 结构体是相同的，但是结构体属性中存在不可以进行比较的类型，如 map、slice，此时结构体不能使用 == 进行比较，如本例所示，存在不可比较的 map 类型。但可以使用reflect.DeepEqual进行比较。
>
>    ```go
>    if reflect.DeepEqual(sm1, sm2) {
>    		fmt.Println("sm1 == sm2") // true
>    } else {
>    		fmt.Println("sm1 != sm2")
>    }
>    ```

#### 3. string 与 nil 类型

> 代码是否能够编译通过？为什么？

```go
package main

import (
    "fmt"
)

func GetValue(m map[int]string, id int) (string, bool) {
	if _, exist := m[id]; exist {
		return "存在数据", true
	}
	return nil, false
}

func main()  {
	intmap:=map[int]string{
		1:"a",
		2:"bb",
		3:"ccc",
	}

	v,err:=GetValue(intmap,3)
	fmt.Println(v,err)
}
```

> 编译不通过，报错: `cannot use nil as type string in return argument`
>
> nil 可以用作 interface、function、pointer、map、slice 和 channel 的“空值”。
>
> 所以将 Get Value 函数改成如下形式就可以了
>
> ```go
> func GetValue(m map[int]string, id int) (string, bool) {
> 	if _, exist := m[id]; exist {
> 		return "存在数据", true
> 	}
>     return "不存在数据", false
> }
> ```

#### 4. 常量

> 此处的函数有什么问题？

```go
package main

const cl = 100

var bl = 123

func main()  {
    println(&bl,bl)
    println(&cl,cl)
}
```

>常量不存在地址，只是一个字面量符号，故会报错，`cannot take the address of cl`
>
>常量不同于变量的在运行期分配内存，常量通常会被编译器在预处理阶段直接展开，作为指令数据使用。
>
>**内存四区概念:** 
>
>+ 数据类型的本质：固定内存大小的别名。
>
>+ 数据类型的作用：编译器预算对象（变量）分配的内存空间大小。
>
>![数据类型的作用](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305205534946.png "数据类型的作用")
>
>+ 内存四区
>
>  流程说明：
>
>  1. 操作系统将物理硬盘代码 load 到内存
>  2. 操作系统将代码分成四个区
>  3. 操作系统找到 main 函数入口执行
>
>  ![内存四区](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305205646905.png "内存四区")
>
>  **栈区(Stack)：**
>
>  **空间较小，要求数据读写性能高，数据存放时间较短暂**。由**编译器自动分配和释放**，存放函数的参数值、函数的调用流程方法地址、局部变量等(**局部变量如果产生逃逸现象，可能会挂在在堆区**)
>
>  **堆区(heap):**
>
>  空间充裕，数据存放时间较久。**一般由开发者分配及释放**，但是在 Golang 中会根据变量的逃逸现象来选择是否分配到栈上或堆上，启动 Golang 的 GC 根据 GC 清除机制自动回收。
>
>  **全局区-静态全局变量区:**
>
>  全局变量的开辟是在程序在 main() 之前就已经放在内存中。而且**对外完全可见**。即作用域在全部代码中，**任何同包代码均可随时使用**，在变量会搞混淆，而且在局部函数中如果同名称变量使用`:=`赋值会出现编译错误。全局变量最终在进程退出时，由操作系统回收。在进行开发的时候，尽量减少使用全局变量。
>
>  **全局区-常量区：**
>
>  常量区也归属于全局区，常量为存放数值字面值单位，即不可修改。或者说的有的常量是直接挂钩字面值的。比如:
>
>  ```
>  const cl = 10
>  ```
>
>  cl 是字面量10的对等符号。所以在 Golang 中，常量无法取出地址，因为字面量符号没有地址。

-------------------------

### 数组与切片

#### 1. 切片的初始化与追加

> 写出程序运行结果

```go
package main

import (
    "fmt"
)

func main(){
    s := make([]int, 10)
    s = append(s, 1, 2, 3) // 切片追加,make() 初始化为0
    fmt.Println(s)
}
// 输出: [0 0 0 0 0 0 0 0 0 0 1 2 3] // 10个0 + 1 2 3
```

#### 2. slice 拼接问题

> 代码是否可以编译通过？

```go
package main

import "fmt"

func main() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	s1 = append(s1, s2)
	fmt.Println(s1)
}
```

> 编译失败，报错: `cannot use s2 (type []int) as type int in append`
>
> ```go
> // The append built-in function appends elements to the end of a slice. If
> // it has sufficient capacity, the destination is resliced to accommodate the
> // new elements. If it does not, a new underlying array will be allocated.
> // Append returns the updated slice. It is therefore necessary to store the
> // result of append, often in the variable holding the slice itself:
> // slice = append(slice, elem1, elem2)
> // slice = append(slice, anotherSlice...)  // 注意此处
> // As a special case, it is legal to append a string to a byte slice, like this:
> // slice = append([]byte("hello "), "world"...)
> func append(slice []Type, elems ...Type) []Type
> ```
>
> append() 方法不支持直接将第二个 slice 传入，需要`s1 = append(s1, s2...)`才成功，即进行 ... 打散在拼接。

#### 3. slice 中 new 的使用

> 代码是否可以编译通过？

```go
package main

import "fmt"

func main() {
	list := new([]int)
	list = append(list, 1)
	fmt.Println(list)
}
```

> 编译失败，报错: `first argument to append must be slice; have *[]int`
>
> 分析：切片指针的解引用，可以使用 list:=make([]int,0)，list类型为切片或使用 list = append(list, 1)，list类型为指针。
>
> new 和 make 的区别：
>
> 二者都是内存的分配（堆上），但是 make 只用于 slice、map、channel 的初始化（非零值）；而 new 用于类型的内存分配，且默认为零值。make 返回的是这三个引用类型本身；而 new 返回的是指向类型的指针。

### Map

#### 1. Map 的 Value 赋值

> 代码编译会出现什么结果？

```go
package main

import "fmt"

type Student struct {
	Name string
}

var list map[string]Student

func main() {

	list = make(map[string]Student)

	student := Student{"Aceld"}

	list["student"] = student
	list["student"].Name = "LDB"

	fmt.Println(list["student"])
}
```

> 编译失败，报错: `cannot assign to struct field list["student"].Name in map`
>
> map[string]Student 的 value 是一个 Student 结构值，所以当 list["student"] = student，是一个值拷贝过程。而list["student"]则是一个值引用，那么值引用的特点是只读，所以进行类似list["student"].Name = "LDB" 的修改是不被允许的。
>
> **方法一：**
>
> ```go
> tmpStudent := list["student"]
> tmpStudent.Name = "LDB"
> list["student"] = tmpStudent
> ```
>
> 是先做一次值拷贝，获得 tmpStudent 副本，然后修改该副本，然后再次发生一次值拷贝复制回去，list["student"] = tmpStudent ，但是这种会在整体过程中发生2次结构体值拷贝，性能很差。
>
> **方法二：**
>
> 将 map 的类型的 value 由 Student 值，改成 Student 指针。
>
> ```go
> var list map[string]*Student
> ```
>
> 这样实际上每次修改的都是指针所指向的 Student 空间，指针本身是常指针，不能修改，只读属性，但是指向的 Student 是可以随便修改的，而且这里并不需要值拷贝。只是一个指针的赋值。

#### 2. map 的遍历赋值

> 代码有什么问题，说明原因

```go
package main

import (
    "fmt"
)

type student struct {
    Name string
    Age  int
}

func main() {
    // 定义map
    m := make(map[string]*student)

    // 定义student数组
    stus := []student{
        {Name: "zhou", Age: 24},
        {Name: "li", Age: 23},
        {Name: "wang", Age: 22},
    }

    // 将数组依次添加到map中
    for _, stu := range stus {
        m[stu.Name] = &stu
    }

    // 打印map
    for k,v := range m {
        fmt.Println(k ,"=>", v.Name)
    }
}
```

> 遍历结果出现错误，输出结果中三个 key 均指向数组中最后一个结构体。
>
> ```bash
> zhou => wang
> li => wang
> wang => wang
> ```
>
> 分析：
>
> for each 中，stu 是结构体的一个拷贝副本，所以 m[stu.Name]=&stu 实际上一致指向**同一个指针**， 最终该指针的值被覆盖为遍历的最后一个 struct 的值拷贝。
>
> ```bash
> // 输出遍历赋值中的数据,可以发现指向的 value 值为同一个地址,实际的值会在循环中不断覆盖
> &{zhou 24}
> map[zhou:0xc000114060]
> &{li 23}
> map[li:0xc000114060 zhou:0xc000114060]
> &{wang 22}
> map[li:0xc000114060 wang:0xc000114060 zhou:0xc000114060]
> ```
>
> 
>
> ![for each 执行情况](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305214204405.png "for each 执行情况")
>
> 正确写法：
>
> ```go
> // 遍历结构体数组，依次赋值给map
> for i := 0; i < len(stus); i++  {
>     m[stus[i].Name] = &stus[i] 
> }
> ```
>
> ![正确的遍历赋值情况](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305214415640.png "正确的遍历赋值情况")

### interface

#### 1. interface 的赋值问题

> 代码是否能编译成功？为什么？

```go
package main

import (
	"fmt"
)

type People interface {
	Speak(string) string
}

type Stduent struct{}

// 此处修改为 (stu Stduent),不建议使用，使用指针可以修改结构体参数
func (stu *Stduent) Speak(think string) (talk string) {
	if think == "love" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func main() {
	var peo People = Stduent{} // 或者此处修改为 &Stduent{}
	think := "love"
	fmt.Println(peo.Speak(think))
}
```

> 编译失败，报错: `cannot use Stduent{} (type Stduent) as type People in assignment:
> 	Stduent does not implement People (Speak method has pointer receiver)`
>
> 多态发生的要素有：
>
> 1. 有interface接口，并且有接口定义的方法。
>
> 2. 有子类重写了interface的接口。
>
> 3. 有父类指针指向子类的具体对象
>
> 满足上述条件后，父类指针就可以调用子类的具体方法。
>
> 所以上述代码报错的地方在 var peo People = Student{} 这条语句，此时需要将父类指针指向子类对象，即 var peo People = &Student{}

#### 2. interface 的内部构造(非空接口 iface 情况)

> 代码会打印什么内容？为什么？

```go
package main

import (
	"fmt"
)

type People interface {
	Show()
}

type Student struct{}

func (stu *Student) Show() {

}

func live() People {
	var stu *Student
	return stu
}

func main() {
	if live() == nil {
		fmt.Println("AAAAAAA")
	} else {
		fmt.Println("BBBBBBB")
	}
}
// 输出: BBBBBBB
```

> 分析:
>
> 此处需要了解 interface 的内部结构，才能理解题目的含义。
>
> interface 在使用的过程中，共有两种表现形式
>
> 一种为空接口(empty interface)，定义如下
>
> ```go
> var MyInterface interface{}
> ```
>
> 另一种为非空接口，定义如下
>
> ```go
> type MyInterface interface {
> 		function()
> }
> ```
>
> 这两种 interface 类型分别用两种 struct 表示，空接口为 eface, 非空接口为 iface。
>
> ![两种 interface 类型](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305222100734.png "两种 interface 类型")
>
> **空接口 eface**
>
> 空接口 eface 结构，由两个属性构成，一个是类型信息 _type，一个是数据信息。其数据结构声明如下：
>
> ```go
> type eface struct {      // 空接口
>  _type *_type         // 类型信息
>  data  unsafe.Pointer // 指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
> }
> ```
>
> **_type属性**：是GO语言中所有类型的公共描述，Go语言几乎所有的数据结构都可以抽象成  _type，是所有类型的公共描述，**type负责决定data应该如何解释和操作，**type的结构代码如下:
>
> ```go
> type _type struct {
>  size       uintptr  // 类型大小
>  ptrdata    uintptr  // 前缀持有所有指针的内存大小
>  hash       uint32   // 数据hash值
>  tflag      tflag
>  align      uint8    // 对齐
>  fieldalign uint8    // 嵌入结构体时的对齐
>  kind       uint8    // kind 有些枚举值kind等于0是无效的
>  alg        *typeAlg // 函数指针数组，类型实现的所有方法
>  gcdata    *byte
>  str       nameOff
>  ptrToThis typeOff
> }
> ```
>
> **data属性:** 表示指向具体的实例数据的指针，他是一个 unsafe.Pointer 类型，相当于一个C的万能指针 void* 
>
> ![空指针 eface 结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305222445878.png "空指针 eface 结构")
>
> **非空接口 iface**
>
> iface 表示 non-empty interface 的数据结构，非空接口初始化的过程就是初始化一个iface类型的结构，其中data 的作用同 eface 中的相同。
>
> ```go
> type iface struct {
> tab  *itab
> data unsafe.Pointer
> }
> ```
>
> iface 结构中最重要的是 itab 结构（结构如下），每一个 itab 占据32字节的空间。itab 可以理解为 pair<interface type, concrete type>。itab 中包含了 interface 的一些关键信息，比如 method 的具体实现。
>
> ```go
> type itab struct {
> inter  *interfacetype   // 接口自身的元信息
> _type  *_type           // 具体类型的元信息
> link   *itab
> bad    int32
> hash   int32            // _type里也有一个同样的hash，此处多放一个是为了方便运行接口断言
> fun    [1]uintptr       // 函数指针，指向具体类型所实现的方法
> }
> ```
>
> 其中值得注意的字段如下：
>
> 1. `interface type`包含了一些关于interface本身的信息，比如`package path`，包含的`method`。这里的interface type是定义interface的一种抽象表示。
> 2. `type`表示具体化的类型，与eface的 *type类型相同。*
> 3. `hash`字段其实是对`_type.hash`的拷贝，它会在interface的实例化时，用于快速判断目标类型和接口中的类型是否一致。另，Go 的interface的 Duck-typing 机制也是依赖这个字段来实现。
> 4. `fun`字段其实是一个动态大小的数组，虽然声明时是固定大小为1，但在使用时会直接通过fun指针获取其中的数据，并且不会检查数组的边界，所以该数组中保存的元素数量是不确定的。
>
> ![非空接口 iface 结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305223005074.png "非空接口 iface 结构")
>
> 所以，People 拥有一个 Show 方法的，属于非空接口，People的内部定义应该是一个 iface 结构体
>
> ```go
> type People interface {
>     Show()  
> }  
> ```
>
> ![People 接口结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305223131355.png "People 接口结构")
>
> ```go
> func live() People {
>     var stu *Student
>     return stu      
> }     
> ```
>
> stu 是一个指向 nil 的空指针，但是最后 return stu  会触发 匿名变量 People = stu 的值拷贝动作，所以最后 live() 返回给上层的是一个 People interface{} 类型，也就是一个 iface struct{} 类型。 stu 为 nil，只是 iface 中的data 为 nil 而已。 **但是 iface struct{} 本身并不为 nil。**
>
> ![代码执行情况](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305223859123.png "代码执行情况")
>
> 故如下代码的判断结果为 BBBBBBB:
>
> ```go
> func main() {   
>     if live() == nil {  
>         fmt.Println("AAAAAAA")      
>     } else {
>         fmt.Println("BBBBBBB")
>     }
> }
> ```
>
> 当使用 == 直接将一个 interface 与 nil 进行比较的时候，Golang 会对 interface 的类型和值分别进行判断。如果两者都为 nil，在与 nil 直接比较时才会返回 true，否则直接返回 false。所以上面代码中 interface 与 nil 进行比较时返回的是 false，**因为此时 interface 变量的值是nil，但是它的类型不是 nil**，已经有了明确的实现类型，即 *Student。
>
> 故在实际开发过程中，**当 interface 类型的返回值已经明确为 nil 时，应该直接返回 nil**，而不是具体实现结构的未赋值空指针。修改原 live() 如下。
>
> ```go
> func live() People{
> 	var stu *Student
> 	if(stu == nil){
> 		return nil
> 	}
> 	return stu
> }
> // 输出: AAAAAAA
> ```

#### 3. interface 内部构造(空接口 eface 情况)

>  代码的执行结果是什么？为什么？

```go
func Foo(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}
func main() {
	var p *int = nil
	Foo(p)
}
// 结果：non-empty interface
```

> 分析：
>
> Foo() 的形参 x interface{} 是一个空接口类型 eface struct{}。
>
> ![x interface{} 内部结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305224715580.png "x interface{} 内部结构")
>
> 在执行 Foo(p) 的时候，触发 x interface{} = p 语句，所以此时 x 的内部结构为:
>
> ![触发后 x 的内部结构](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305224829177.png "触发后 x 的内部结构")
>
> 所以 x 结构体本身不为 nil，而是 data 指针指向的 p 为 nil。**即类型不为空，值为空，在与 nil 的直接比较时会返回 false。**

#### 4. interface{} 与 *interface{}

> ABCD 哪一行存在错误？

```go
type S struct {
}

func f(x interface{}) {
}

func g(x *interface{}) {
}

func main() {
	s := S{}
	p := &s
	f(s) //A
	g(s) //B
	f(p) //C
	g(p) //D
}
```

>B、D两行错误
>B错误为： `cannot use s (type S) as type *interface {} in argument to g:
>	*interface {} is pointer to interface, not interface`
>D错误为：`cannot use p (type *S) as type *interface {} in argument to g:
>	*interface {} is pointer to interface, not interface`
>
>看到这道题需要第一时间想到的是 Golang 是强类型语言，**interface 是所有 Golang 类型的父类**，函数中 func f(x interface{})  的 interface{} 可以支持传入 Golang 的任何类型，包括指针。但是函数 func g(x *interface{}) 只能接受 *interface{}。

### channel

#### channel 读写特性(15子口诀)

channel有以下特性：

- 给一个 nil channel 发送数据，造成永远阻塞
- 从一个 nil channel 接收数据，造成永远阻塞
- 给一个已经关闭的 channel 发送数据，引起 panic
- 从一个已经关闭的 channel 接收数据，如果缓冲区中为空，则返回一个零值
- 无缓冲的channel是同步的，而有缓冲的channel是非同步的

可以通过口诀来记忆：**“空读写阻塞，写关闭异常，读关闭空零”**。

> 代码执行的结果是什么？

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1000)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("close")
				return
			}
			fmt.Println("a: ", a)
		}
	}()
	close(ch)
	fmt.Println("ok")
    time.Sleep(time.Second * 100) // 将此处移至 close(ch) 前即可
}
```

> 代码向已关闭的 channel 写入数据会产生 panic。main() 再开辟完两个 goroutine 后就立即管理了ch，下方为执行结果。

```bash
ok
panic: send on closed channel
```

### WaitGroup

#### WaitGroup 与 goroutine的竞速问题

> 编译并运行代码会发生什么？

```go
package main

import (
	"sync"
)

const N = 10

// WaitGroup 能够阻塞主线程的执行，直到所有的goroutine执行完成。
var wg = &sync.WaitGroup{}

func main() {

	for i := 0; i < N; i++ {
		go func(i int) {
			wg.Add(1) // 添加或者减少等待goroutine的数量
			println(i)
            // 相当于Add(-1),wg.Done()最好使用defer注册一下，避免函数内部出错执行不到
			defer wg.Done() 
		}(i)
	}
    // Wait:执行阻塞，直到所有的WaitGroup数量变成0
	wg.Wait()
    
}
```

> 结果不唯一，代码存在风险，所有的 go 程未必能执行到。多次执行上面的代码会发现输出都会不同甚至会出现报错的问题。
>
> ```bash
> panic: sync: WaitGroup is reused before previous Wait has returned
> ```
>
> 只要计数器的值始于0又归为0，就可以被视为一个计数周期。在一个此类值的生命周期中，它可以经历任意多个计数周期。但是，只有在它走完当前的计数周期之后，才能够开始下一个计数周期。
>
> ![WaitGroup 计数值变化](https://raw.githubusercontent.com/tonshz/test/master/img/image-20220305233926038.png "WaitGroup 计数值变化")
>
> **不要把增加计数器值的操作和调用Wait方法的代码，放在不同的 `goroutine` 中执行。换句话说，要杜绝对同一个`WaitGroup` 值的两种操作的并发执行。**
>
> 最好用 **先统一 `Add` ，再并发 `Done` ，最后 `Wait`** 这种标准方式，来使用 `WaitGroup` 值。 尤其不要在调用 `Wait` 方法的同时，并发地通过调用 `Add` 方法去增加其计数器的值，因为这也有可能引发 `panic` 。
>
> 这是因为 go 执行太快了，导致 wg.Add(1) 还没有执行 main 函数就执行完毕了。修改为如下代码即可。
>
> ```go
> package main
> 
> import (
> 	"sync"
> )
> 
> const N = 10
> 
> var wg = &sync.WaitGroup{}
> 
> func main() {
> 
>     for i:= 0; i< N; i++ {
>         wg.Add(1) // Add放在Goroutine外
>         go func(i int) {
>             println(i)
>             defer wg.Done() // Done放在Goroutine中，逻辑复杂时建议用defer保证调用
>         }(i)
>     }
> 
>     wg.Wait()
> }
> ```

