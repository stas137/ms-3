package input

import "github.com/google/uuid"

type CreateOrderInput struct {
	EngineUUID uuid.UUID
	HullUUID   uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

func (i *CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{i.EngineUUID, i.HullUUID}
	if i.ShieldUUID != nil {
		uuids = append(uuids, *i.ShieldUUID)
	}
	if i.WeaponUUID != nil {
		uuids = append(uuids, *i.WeaponUUID)
	}
	return uuids
}
