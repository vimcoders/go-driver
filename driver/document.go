// 不允许调用标准库外的包，防止循环引用
package driver

type Document interface {
	DocumentId() string
	DocumentName() string
}
