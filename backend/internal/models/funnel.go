package models

import "gorm.io/gorm"

type Funnel struct {
	gorm.Model
	Name            string    `json:"name" binding:"required"`
	NextFunnels     []*Funnel `gorm:"many2many:funnel_transitions;joinForeignKey:from_funnel_id;joinReferences:to_funnel_id" json:"next_funnels"`
	PreviousFunnels []*Funnel `gorm:"many2many:funnel_transitions;joinForeignKey:to_funnel_id;joinReferences:from_funnel_id" json:"previous_funnels"`
}
