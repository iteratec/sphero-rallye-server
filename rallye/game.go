package rallye

import "github.com/iteratec/sphero-rallye-server/rallye/actions"

func InitGame() {
	actions.HandleIncomingSpheroActions()
}

