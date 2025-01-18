package checker

import "github.com/mmtaee/go-oc-utils/models"

func GetIds(users []models.OcUser) []uint {
	ids := make([]uint, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}
