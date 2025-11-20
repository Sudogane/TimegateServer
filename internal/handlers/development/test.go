package development

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
)

type TestHandler struct {
	gamseServer server.GameServerInterface
}

func NewTestHandler(gs server.GameServerInterface) *TestHandler {
	return &TestHandler{
		gamseServer: gs,
	}
}

func (th *TestHandler) Handle(session *server.PlayerSession, userService *services.UserService) {
	//userID := session.PlayerId
	//test, err := userService.GetAvailableStagesByEpisodeId(1, userID)
	//fmt.Println("tesst", test, err)

}
