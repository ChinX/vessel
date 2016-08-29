package point

import "github.com/containerops/vessel/models"

func CheckCondition(pointVsn *models.PointVersion, readyMap map[string]bool) (bool, bool) {
	meet := true
	ended := false
	if pointVsn.Kind == models.StartPoint {
		return meet, ended
	}
	for _, condition := range pointVsn.Conditions {
		if meet := readyMap[condition]; !meet{
			return meet, ended
		}
	}
	ended = pointVsn.Kind == models.EndPoint
	return meet, ended
}
