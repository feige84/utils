package utils

import (
	"bytes"
	"fmt"
)

type HashSet struct {
	m map[interface{}]bool
}
type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
}

func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}

//方法Add会返回一个bool类型的结果值，以表示添加元素值的操作是否成功。
//方法Add的声明中的接收者类型是*HashSet。
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] { //当前的m的值中还未包含以e的值为键的键值对
		set.m[e] = true //将键为e(代表的值)、元素为true的键值对添加到m的值当中
		return true     //添加成功
	}
	return false //添加失败
}

//为HashSet中的字段m重新赋值
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

//方法Contains用于判断其值是否包含某个元素值。
//这里判断结果得益于元素类型为bool的字段m
func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

//调用delete内建函数删除HashSet内部支持的字典值
func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e) //第一个参数为目标字典类型，第二个参数为要删除的那个键值对的键
}

//方法Len用于获取HashSet元素值数量
func (set *HashSet) Len() int {
	return len(set.m)
}

//方法Same用来判断两个HashSet类型值是否相同
func (set *HashSet) Same(other *HashSet) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

//这个String方法的签名算是一个惯用法。 //代码包fmt中的打印函数总会使用参数值附带的具有如此签名的String方法的结果值作为该参数值的字符串表示形式。
func (set *HashSet) String() string {
	var buf bytes.Buffer //作为结果值的缓冲区
	buf.WriteString("HashSet{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	//n := 1
	//for key := range set.m {
	// buf.WriteString(fmt.Sprintf("%v", key))
	// if n == len(set.m) {//最后一个元素的后面不添加逗号
	// break;
	// } else {
	// buf.WriteString(",")
	// }
	// n++;
	//}
	buf.WriteString("}")
	return buf.String()
}

//方法Elements用于生成快照
func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m) //获取HashSet中字段m的长度，即m中包含元素的数量
	//初始化一个[]interface{}类型的变量snapshot来存储m的值中的元素值
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	//按照既定顺序将迭代值设置到快照值(变量snapshot的值)的指定元素位置上,这一过程并不会创建任何新值。
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else { //m的值中的元素数量有所增加，使得实际迭代的次数大于先前初始化的快照值的长度
			snapshot = append(snapshot, key) //使用append函数向快照值追加元素值。
		}
		actualLen++ //实际迭代的次数
	}
	//对于已被初始化的[]interface{}类型的切片值来说，未被显示初始化的元素位置上的值均为nil。
	//m的值中的元素数量有所减少，使得实际迭代的次数小于先前初始化的快照值的长度。
	//这样快照值的尾部存在若干个没有任何意义的值为nil的元素，
	//可以通过snapshot = snapshot[:actualLen]将无用的元素值从快照值中去掉。
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

// 判断集合 set 是否是集合 other 的超集
func (set *HashSet) IsSuperset(other *HashSet) bool {
	if other == nil { //如果other为nil，则other不是set的子集
		return false
	}
	setLen := set.Len()                    //获取set的元素值数量
	otherLen := other.Len()                //获取other的元素值数量
	if setLen == 0 || setLen == otherLen { //set的元素值数量等于0或者等于other的元素数量
		return false
	}
	if setLen > 0 && otherLen == 0 { //other为元素数量为0，set元素数量大于0，则set也是other的超集
		return true
	}
	for _, v := range other.Elements() {
		if !set.Contains(v) { //只要set中有一个包含other中的数据，就返回false
			return false
		}
	}
	return true
}

// 生成集合 set 和集合 other 的并集
func (set *HashSet) Union(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的并集为nil
		return nil
	}
	unionedSet := NewHashSet()         //新创建一个HashSet类型值，它的长度为0，即元素数量为0
	for _, v := range set.Elements() { //将set中的元素添加到unionedSet中
		unionedSet.Add(v)
	}
	if other.Len() == 0 {
		return unionedSet
	}
	for _, v := range other.Elements() { //将other中的元素添加到unionedSet中，如果遇到相同，则不添加（在Add方法逻辑中体现）
		unionedSet.Add(v)
	}
	return unionedSet
}

// 生成集合 set 和集合 other 的交集
func (set *HashSet) Intersect(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的交集为nil
		return nil
	}
	intersectedSet := NewHashSet() //新创建一个HashSet类型值，它的长度为0，即元素数量为0
	if other.Len() == 0 {          //other的元素数量为0，直接返回intersectedSet
		return intersectedSet
	}
	if set.Len() < other.Len() { //set的元素数量少于other的元素数量
		for _, v := range set.Elements() { //遍历set
			if other.Contains(v) { //只要将set和other共有的添加到intersectedSet
				intersectedSet.Add(v)
			}
		}
	} else { //set的元素数量多于other的元素数量
		for _, v := range other.Elements() { //遍历other
			if set.Contains(v) { //只要将set和other共有的添加到intersectedSet
				intersectedSet.Add(v)
			}
		}
	}
	return intersectedSet
}

// 生成集合 set 对集合 other 的差集
func (set *HashSet) Difference(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的差集为nil
		return nil
	}
	differencedSet := NewHashSet() //新创建一个HashSet类型值，它的长度为0，即元素数量为0
	if other.Len() == 0 {          // 如果other的元素数量为0
		for _, v := range set.Elements() { //遍历set，并将set中的元素v添加到differencedSet
			differencedSet.Add(v)
		}
		return differencedSet //直接返回differencedSet
	}
	for _, v := range set.Elements() { //other的元素数量不为0，遍历set
		if !other.Contains(v) { //如果other中不包含v，就将v添加到differencedSet中
			differencedSet.Add(v)
		}
	}
	return differencedSet
}

// 生成集合 one 和集合 other 的对称差集
func (set *HashSet) SymmetricDifference(other *HashSet) *HashSet {
	if set == nil || other == nil { // set和other都为nil，则它们的对称差集为nil
		return nil
	}
	diffA := set.Difference(other) //生成集合 set 对集合 other 的差集
	if other.Len() == 0 {          //如果other的元素数量等于0，那么other对集合set的差集为空，则直接返回diffA
		return diffA
	}
	diffB := other.Difference(set) //生成集合 other 对集合 set 的差集
	return diffA.Union(diffB)      //返回集合 diffA 和集合 diffB 的并集
}

// 判断集合 one 是否是集合 other 的超集
// 读者应重点关注IsSuperset与附属于HashSet类型的IsSuperset方法的区别
func IsSuperset(one Set, other Set) bool {
	if one == nil || other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}
