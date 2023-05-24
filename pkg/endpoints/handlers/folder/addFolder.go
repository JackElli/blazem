package folder

import (
	"blazem/pkg/domain/endpoint"
	types "blazem/pkg/domain/endpoint"
	"blazem/pkg/domain/endpoint_manager"
	blazem_folder "blazem/pkg/domain/folder"
	"errors"

	"blazem/pkg/domain/global"
	"encoding/json"
	"net/http"
	"time"
)

// We want to add a document to Blazem, we check if it's a POST, unmarshal the data
// coming in, write to disk and add to the map
func AddFolder(e *endpoint_manager.EndpointManager) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if e.Node.Rank != global.MASTER {
			e.Responder.Error(w, 500, errors.New("Should be master"))
			return
		}

		c, err := req.Cookie("token")
		if err != nil {
			e.Responder.Error(w, 500, errors.New("No user logged in"))
			return
		}

		var folder endpoint.Folder
		err = json.NewDecoder(req.Body).Decode(&folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		err = validate(folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		jwtVal := c.Value
		userId, err := e.GetCurrentUserId(jwtVal)
		folder.DateCreated = time.Now().Format("2006-01-02T15:04:05")
		folder.CreatedBy = userId

		folderMap, err := blazem_folder.FolderToMap(folder)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		err = e.DataStore.Store(folder.Key, folder.Folder, folderMap)
		if err != nil {
			e.Responder.Error(w, 500, err)
			return
		}

		global.DataChanged = true
		e.Responder.Respond(w, 200, types.EndpointResponse{
			Msg: "Added folder successfully",
		})
	}
}

// validate checks whether the incoming folder is a valid
// folder structure
func validate(folder endpoint.Folder) error {
	if folder.Name == "" ||
		folder.Key == "" {
		return errors.New("Folder is invalid")
	}

	return nil
}
