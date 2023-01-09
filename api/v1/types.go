package v1

import "github.com/google/uuid"

type EntityID uuid.UUID

func ParseEntityID(s string) (EntityID, error) {
	id, err := uuid.Parse(s)
	return EntityID(id), err
}

func (id EntityID) String() string {
	return uuid.UUID(id).String()
}
