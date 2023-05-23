package folder

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	blazem_folder "blazem/pkg/domain/folder"
	"blazem/pkg/domain/middleware"

	"blazem/pkg/domain/global"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func AddFolder(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank != global.MASTER {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Should be master",
			})
			return
		}
		var folder endpoint.Folder
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Cannot read request body {" + err.Error() + "}",
			})
			return
		}

		err = json.Unmarshal(body, &folder)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "Cannot unmarshal JSON request {" + err.Error() + "}",
			})
			return
		}

		c, err := req.Cookie("token")
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  "No user logged in",
			})
			return
		}

		jwtVal := c.Value
		userId, err := middleware.GetCurrentUserId(jwtVal)

		folder.DateCreated = time.Now().Format("2006-01-02T15:04:05")
		folder.CreatedBy = userId

		folderMap, err := blazem_folder.FolderToMap(folder)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		err = e.DataStore.Store(folder.Key, folder.Folder, folderMap)
		if err != nil {
			e.Responder.Respond(w, types.EndpointResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		global.DataChanged = true
		e.Responder.Respond(w, types.EndpointResponse{
			Code: 200,
			Msg:  "Added folder successfully",
		})
	}
}
