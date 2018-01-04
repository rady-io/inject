package rorm

import ry "rhapsody"

type GormConfig struct {
	ry.Configuration
	App *ry.Application
}
