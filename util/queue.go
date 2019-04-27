package util

type QNode struct {
	Prev *QNode
	Next *QNode
	Value interface{}
}

type Queue struct {
	Top *QNode
	End *QNode
}

func (q *Queue) Append(value interface{}) {
	qn := &QNode{
		Value:value,
	}
	if q.Top == nil {
		q.Top = qn
		q.End = qn
		return
	}
	q.End.Next = qn
	qn.Prev = q.End
	qn.Next = nil
	q.End = qn
}

func (q *Queue) Pop() ( n *QNode, err error) {
	if q.Top == nil {
		err := QueueOutFlowError{
			Message:"Out of Queue",
			Level:"0001",
		}
		return nil, err
	}
	n = q.End
	if n.Prev == nil {
		q.End = nil
		q.Top = nil
		return
	}
	q.End = n.Prev
	q.End.Next = nil
	n.Prev = nil
	return
}

func (q *Queue) Len() int{
	var (
		cur = &QNode{}
		length int
	)
	length = 0
	cur = q.Top
	for {
		if cur == nil {
			break
		}
		length ++
		cur = cur.Next
	}
	return length
}
