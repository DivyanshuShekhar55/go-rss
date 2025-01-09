package rss

import ()

type RSS interface {
	GetFeed() ([]string, error)
}