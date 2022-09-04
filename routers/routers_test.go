package routers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ernesto27/docs/db"
	"github.com/ernesto27/docs/interfaces"
	"github.com/ernesto27/docs/structs"
	"github.com/gin-gonic/gin"
)

func TestCreateDocs(t *testing.T) {
	type args struct {
		db         interfaces.DocDB
		formParams map[string]string
	}

	type result struct {
		httpStatus int
		response   structs.ResponseApi
	}

	tests := []struct {
		name     string
		args     args
		expected result
	}{
		{
			name: "Error create new doc if empty title or body",
			args: args{
				db: &db.Mock{},
				formParams: map[string]string{
					"title": "",
					"body":  "",
				},
			},
			expected: result{
				httpStatus: http.StatusOK,
				response: structs.ResponseApi{
					Status:  "error",
					Message: "title or body is empty",
				},
			},
		},
		{
			name: "Error create new doc DB error",
			args: args{
				db: &MockDBError{},
				formParams: map[string]string{
					"title": "title",
					"body":  "body",
				},
			},
			expected: result{
				httpStatus: http.StatusOK,
				response: structs.ResponseApi{
					Status:  "error",
					Message: "error creating doc",
				},
			},
		},
		{
			name: "Success create new doc",
			args: args{
				db: &db.Mock{},
				formParams: map[string]string{
					"title": "title",
					"body":  "body",
				},
			},
			expected: result{
				httpStatus: http.StatusOK,
				response: structs.ResponseApi{
					Status:  "success",
					Message: "success created doc",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			buf := new(bytes.Buffer)
			mw := multipart.NewWriter(buf)

			for k, v := range tt.args.formParams {
				mw.WriteField(k, v)
			}
			mw.Close()

			c.Request, _ = http.NewRequest(http.MethodPost, "docs/create", buf)
			c.Request.Header.Set("Content-Type", mw.FormDataContentType())

			CreateDoc(tt.args.db, c)

			if w.Code != tt.expected.httpStatus {
				t.Errorf("got %v, want %v", w.Code, tt.expected.httpStatus)
			}

			var resBody structs.ResponseApi
			json.Unmarshal(w.Body.Bytes(), &resBody)
			if resBody.Status != tt.expected.response.Status {
				t.Errorf("Expected status %v, got %v", tt.expected.response.Status, resBody.Status)
			}

			if resBody.Message != tt.expected.response.Message {
				t.Errorf("Expected status %v, got %v", tt.expected.response.Message, resBody.Message)
			}

		})
	}

}

type MockDBError struct {
	db.Mock
}

func (m *MockDBError) CreateDoc(doc structs.Doc) error {
	return errors.New("error creating doc")
}
