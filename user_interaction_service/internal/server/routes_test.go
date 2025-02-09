package server

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// type Interaction struct {
// 	Title string `json:"title"`
// 	Type  string `json:"type"`
// }

// // Mock Kafka function
// func MockSendMessage(title string, message string, topic string) error {
// 	if title == "error" {
// 		return fmt.Errorf("Kafka error")
// 	}
// 	return nil
// }

// func TestInteractionsHandler(t *testing.T) {
// 	// Tạo server mock
// 	server := &Server{}
// 	server.SendMessage = MockSendMessage // Đảm bảo server sử dụng mock Kafka function

// 	// 1. Kiểm tra phương thức không hợp lệ
// 	req, err := http.NewRequest("GET", "/interactions", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(server.InteractionsHandler)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

// 	// 2. Kiểm tra lỗi JSON không hợp lệ
// 	invalidJSON := `{"title": "Test Interaction", "type": 123}` // type phải là string
// 	req, err = http.NewRequest("POST", "/interactions", bytes.NewBuffer([]byte(invalidJSON)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)

// 	// 3. Kiểm tra lỗi Kafka
// 	validInteraction := Interaction{
// 		Title: "Test Interaction",
// 		Type:  "video",
// 	}
// 	validJSON, err := json.Marshal(validInteraction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Sử dụng title "error" để trigger Kafka error
// 	req, err = http.NewRequest("POST", "/interactions", bytes.NewBuffer([]byte(`{"title": "error", "type": "video"}`)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)

// 	// 4. Kiểm tra phản hồi thành công
// 	req, err = http.NewRequest("POST", "/interactions", bytes.NewBuffer(validJSON))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr = httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	assert.Equal(t, "OK", rr.Body.String())
// }
