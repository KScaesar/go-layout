package dataflow

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"unsafe"
)

var (
	ErrNotFoundSubject = errors.New("not found subject mux")
)

// middleware 只會用在 handler 和 defaultHandler
// handle message 執行順序, 依照號碼 0~3
type paramHandler struct {
	// any
	middlewares []Middleware

	// 0
	transform HandleFunc

	// 1
	handler     HandleFunc
	handlerName string

	// 2
	defaultHandler     HandleFunc
	defaultHandlerName string

	// 3
	notFoundHandler HandleFunc
}

func (param *paramHandler) register(leafNode *trie, path []Middleware) error {
	if param.middlewares != nil {
		leafNode.middlewares = append(leafNode.middlewares, param.middlewares...)
	}

	if param.transform != nil {
		if leafNode.transform != nil {
			return errors.New("assign duplicated transform")
		}
		leafNode.transform = param.transform
	}

	if param.handler != nil {
		if leafNode.handler != nil {
			return errors.New("assign duplicated handler")
		}
		leafNode.handler = Link(param.handler, path...)

		if param.handlerName == "" {
			leafNode.handlerName = functionName(param.handler)
		} else {
			leafNode.handlerName = param.handlerName
		}
	}

	if param.defaultHandler != nil {
		if leafNode.defaultHandler != nil {
			return errors.New("assign duplicated defaultHandler")
		}
		leafNode.defaultHandler = Link(param.defaultHandler, path...)

		if param.defaultHandlerName == "" {
			leafNode.defaultHandlerName = functionName(param.defaultHandler)
		} else {
			leafNode.defaultHandlerName = param.defaultHandlerName
		}
	}

	if param.notFoundHandler != nil {
		if leafNode.notFoundHandler != nil {
			return errors.New("assign duplicated notFoundHandler")
		}
		leafNode.notFoundHandler = param.notFoundHandler
	}

	return nil
}

func functionName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func newTrie(delimiter string) *trie {
	return &trie{
		staticChild: make(map[byte]*trie),
		delimiter:   delimiter,
	}
}

type trie struct {
	staticChild map[byte]*trie // key : value => char : child

	wildcardChildWord string
	wildcardChild     *trie

	delimiter   string
	fullSubject string
	paramHandler
}

func (node *trie) addRoute(subject string, cursor int, param *paramHandler, path []Middleware) *trie {
	if node.middlewares != nil {
		path = append(path, node.middlewares...)
	}

	if len(subject) == cursor {
		if param == nil { // for Mux.Group
			param = &paramHandler{
				middlewares: path,
			}
		}

		leafNode := node
		err := param.register(leafNode, path)
		if err != nil {
			Err := fmt.Errorf("subject=%q: %w", subject, err)
			panic(Err)
		}
		return leafNode
	}

	char := subject[cursor]
	if char != '{' {
		child, exist := node.staticChild[char]
		if !exist {
			child = newTrie(node.delimiter)
			child.fullSubject = node.fullSubject + string(char)
			node.staticChild[char] = child
		}
		return child.addRoute(subject, cursor+1, param, path)
	}

	if node.delimiter == "" {
		err := fmt.Errorf("subject=%q: route delimiter is empty: not support wildcard", subject)
		panic(err)
	}

	idx := cursor
	for idx < len(subject) && subject[idx] != '}' {
		idx++
	}

	if subject[idx] != '}' {
		err := fmt.Errorf("subject=%q: lack wildcard '}'", subject)
		panic(err)
	}

	if node.wildcardChild != nil {
		if node.wildcardChildWord != subject[cursor+1:idx] {
			err := fmt.Errorf("subject=%q: assign duplicated wildcard: %q", node.wildcardChild.fullSubject, subject)
			panic(err)
		}
		return node.wildcardChild.addRoute(subject, idx+1, param, path)
	}

	child := newTrie(node.delimiter)
	child.fullSubject = node.fullSubject + subject[cursor:idx+1] // {word}, include {}
	node.wildcardChildWord = subject[cursor+1 : idx]             // word, exclude {}
	node.wildcardChild = child
	return child.addRoute(subject, idx+1, param, path)
}

func (node *trie) handleMessage(cursor int, message *Message, dep any) error {
	current := node

	var defaultHandler, notFoundHandler HandleFunc

	const notWildcard = -1
	wildcardStart := notWildcard
	var wildcardParent *trie

	for cursor <= len(message.Subject) {
		if current.transform != nil {
			err := current.transform(message, dep)
			if err != nil {
				return err
			}
		}

		if current.defaultHandler != nil {
			defaultHandler = current.defaultHandler
		}

		if current.notFoundHandler != nil {
			notFoundHandler = current.notFoundHandler
		}

		if current.wildcardChild != nil {
			wildcardStart = cursor
			wildcardParent = current
		}

		if cursor == len(message.Subject) {
			break
		}

		child, exist := current.staticChild[message.Subject[cursor]]
		if !exist {
			break
		}
		cursor++
		current = child
	}

	// for static route
	if current.handler != nil && cursor == len(message.Subject) {
		return current.handler(message, dep)
	}
	if wildcardParent == nil {
		if defaultHandler != nil {
			return defaultHandler(message, dep)
		}
		if notFoundHandler != nil {
			return notFoundHandler(message, dep)
		}
		return ErrNotFoundSubject
	}

	// for wildcard route
	wildcardFinish := wildcardStart
	for wildcardFinish < len(message.Subject) && message.Subject[wildcardFinish] != current.delimiter[0] {
		wildcardFinish++
	}

	bytes := unsafe.Slice(unsafe.StringData(message.Subject), len(message.Subject))
	value := bytes[wildcardStart:wildcardFinish]
	str := unsafe.String(unsafe.SliceData(value), wildcardFinish-wildcardStart)
	message.RouteParam.Set(wildcardParent.wildcardChildWord, str)

	err := wildcardParent.wildcardChild.handleMessage(wildcardFinish, message, dep)
	if err != nil && errors.Is(err, ErrNotFoundSubject) {
		if defaultHandler != nil {
			return defaultHandler(message, dep)
		}
		if notFoundHandler != nil {
			return notFoundHandler(message, dep)
		}
		return ErrNotFoundSubject
	}
	return err
}

// pair = [subject, function]
func (node *trie) endpoint() (pairs [][2]string) {
	pairs = make([][2]string, 0)
	node._endpoint_(&pairs)

	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})
	return
}

func (node *trie) _endpoint_(paris *[][2]string) {
	if node.handler != nil {
		*paris = append(*paris, [2]string{node.fullSubject, node.handlerName})
	}
	if node.defaultHandler != nil {
		*paris = append(*paris, [2]string{node.fullSubject + ".*", node.defaultHandlerName})
	}

	for _, next := range node.staticChild {
		next._endpoint_(paris)
	}

	if node.wildcardChild == nil {
		return
	}
	node.wildcardChild._endpoint_(paris)
}
