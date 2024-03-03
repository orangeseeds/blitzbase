package core

import (
	"github.com/orangeseeds/blitzbase/utils"
)

/*
Hook struct specifies and event and a collection of handlers for the event,

Example:

Let's say there is a onAppServe hook, which says
when AppServe event occurs, trigger a specific handler to handle the event
But if we have Apps which can serve on multiple different ports,
then we will need multiple different handlers to handle the event.

So we can ad multiple different handlers to handle serves on different port.
Inside app.Serve() we will call Hook.Trigger(AppServe)
*/
type HandlerFunc[T any] func(e T) error

type Handler[T any] struct {
	id          string
	handlerFunc HandlerFunc[T]
}

type Hook[T any] struct {
	event    T
	handlers []*Handler[T]
}

func (h *Hook[T]) Add(handler HandlerFunc[T]) string {
	id := utils.RandStr(8)
	h.handlers = append(h.handlers, &Handler[T]{
		id:          id,
		handlerFunc: handler,
	})
	return id
}

func (h *Hook[T]) Remove(id string) {
	for i := len(h.handlers) - 1; i >= 0; i-- {
		if h.handlers[i].id == id {
			h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
			return
		}
	}
}

func (h *Hook[T]) RemoveAll() {
	h.handlers = make([]*Handler[T], 0)
}

func (h *Hook[T]) Trigger(data T) error {
	for _, handler := range h.handlers {
		err := handler.handlerFunc(data)
		if err != nil {
			return err
		}
	}
	return nil
}
