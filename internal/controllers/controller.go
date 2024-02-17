package controllers

import "goshare/pkg/pubsub"

type Controller struct {
	PubSub *pubsub.PubSub
}

func NewController(pubsub *pubsub.PubSub) *Controller {
	return &Controller{
		PubSub: pubsub,
	}
}
