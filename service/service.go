package service

import (
	"fmt"
	w "github.com/MaiaVinicius/go-whatsapp"
	"github.com/MaiaVinicius/wabot/lib"
	"github.com/MaiaVinicius/wabot/model"
	"os"
)

var config = model.GetConfig()

func StartProjects() {
	projects := model.GetProjects()

	for _, element := range projects {

		//prepareReception(element)
		prepareQueue(element)
		prepareReception(element)
	}

	println("===============================")
	println("|     Rotina finalizada.     |")
	println("===============================")
}

func prepareQueue(project model.Project) {
	print(fmt.Sprintf("---->   PROJETO: %s (%s) id: %d \n", project.Phone, project.Label, int(project.SenderID)))

	queue := model.GetQueue(project.SenderID)

	if len(queue) > 0 {
		msg := fmt.Sprintf("----->	Envio iniciado. %d", len(queue))

		model.LogMessage(1, msg, project.ID)

		sendQueue(project.Phone, queue, project.Label, project.ID)
	}
}

func prepareReception(project model.Project) {
	syncResponsesWithServer()

	println(fmt.Sprintf("---->   Recebendo respostas do Celular: %s", project.Phone))

	//var datetimeLastSent string
	//
	//datetimeLastSent = model.GetLastSent(project.ID)
	//
	//t, err := time.Parse("2006-01-02 15:04:05", datetimeLastSent)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}

	var responses []lib.Response

	responses = lib.Receive(project.Phone, config.ReceiveTimeout)

	i := 0

	for _, element := range responses {

		//if element.Timestamp > t.Unix() || datetimeLastSent == "" || true {
		//}
		added := model.InsertResponse(project.ID, element.Phone, element.ID, element.Message, element.Datetime, element.Status, element.FromMe)

		if added {
			i += 1
		}
	}

	msg := fmt.Sprintf("Novas respostas recebidas: %d", i)

	println(msg)
	model.LogMessage(5, msg, project.ID)
}

func sendQueue(senderPhone string, queue []model.Queue, projectName string, projectId int) {
	//connect

	wac, err := lib.Connect(senderPhone)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
	}

	var toRemove []lib.Sent


	for _, element := range queue {
		//element.Message = strings.Replace(element.Message, "[METADATA]", element.Metadata2, -1)

		element.Message = fmt.Sprintf("*%s*. \n\n%s", projectName, element.Message)

		sendMessage(wac, element, projectId, config)

		//add item para remover
		item := lib.Sent{}

		item.LicenseId = element.LicenseId
		item.EventId = element.EventId
		item.AppointmentId = element.AppointmentId

		// append works on nil slices.
		toRemove = append(toRemove, item)
		lib.RemoveQueue(toRemove)
		toRemove = nil
	}

}

func sendMessage(wac *w.Conn, message model.Queue, projectId int, config model.Config) int {

	msg := fmt.Sprintf("Enviando para: %s", message.Phone)
	print(msg)
	model.LogMessage(4, msg, projectId)

	sendMinimumTimeout := config.SendMinimumTimeout
	sendTimeRandom := config.SendTimeRandom
	status := lib.Send(wac, message.Phone, message.Message, sendMinimumTimeout, sendTimeRandom)

	if status == 200 {
		removeFromQueue(message.ID)
	}

	return status
}

func removeFromQueue(messageId int) {
	model.RemoveFromQueue(messageId)
	return
}

func syncResponsesWithServer() {

	var s []lib.ResponseToServer

	responses := model.GetResponsesToSync()

	lastId := 0

	if len(responses) > 0 {
		println("---->	Sincronizando respostas")

		for _, element := range responses {
			item := lib.ResponseToServer{}

			item = element

			// append works on nil slices.
			s = append(s, item)
			if element.AutoId > lastId {
				lastId = element.AutoId
			}
		}

		lib.SendResponsesToServer(s)

		println(lastId)
		model.UpdateSyncedResponses(lastId)
	}
}
