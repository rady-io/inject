package rorm

import ry "rady"

type GormConfig struct {
	ry.Configuration
	App *ry.Application
}
