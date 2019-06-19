package onfido_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getground/go-onfido"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUploadDocument_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	docReq := onfido.DocumentRequest{
		File: bytes.NewReader([]byte("test")),
		Type: onfido.DocumentTypeIDCard,
		Side: onfido.DocumentSideFront,
	}

	_, err := client.UploadDocument(context.Background(), "", docReq)
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestUploadDocument_DocumentUploaded(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.Document{
		ID:           "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86/download",
		FileName:     "localfile.png",
		FileType:     "png",
		FileSize:     282123,
		Type:         onfido.DocumentTypePassport,
		Side:         onfido.DocumentSideBack,
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("POST")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	d, err := client.UploadDocument(context.Background(), applicantID, onfido.DocumentRequest{
		File: bytes.NewReader([]byte("test")),
		Type: expected.Type,
		Side: expected.Side,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, d.ID)
	assert.Equal(t, expected.Href, d.Href)
	assert.Equal(t, expected.DownloadHref, d.DownloadHref)
	assert.Equal(t, expected.FileName, d.FileName)
	assert.Equal(t, expected.FileType, d.FileType)
	assert.Equal(t, expected.FileSize, d.FileSize)
	assert.Equal(t, expected.Type, d.Type)
	assert.Equal(t, expected.Side, d.Side)
}

func TestGetDocument_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	_, err := client.GetDocument(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected server to return non ok response, got successful response")
	}
}

func TestGetDocument_DocumentRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.Document{
		ID:           "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86/download",
		FileName:     "localfile.png",
		FileType:     "png",
		FileSize:     282123,
		Type:         onfido.DocumentTypePassport,
		Side:         onfido.DocumentSideBack,
	}
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents/{documentId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])
		assert.Equal(t, expected.ID, vars["documentId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	d, err := client.GetDocument(context.Background(), applicantID, expected.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.ID, d.ID)
	assert.Equal(t, expected.Href, d.Href)
	assert.Equal(t, expected.DownloadHref, d.DownloadHref)
	assert.Equal(t, expected.FileName, d.FileName)
	assert.Equal(t, expected.FileType, d.FileType)
	assert.Equal(t, expected.FileSize, d.FileSize)
	assert.Equal(t, expected.Type, d.Type)
	assert.Equal(t, expected.Side, d.Side)
}

func TestListDocuments_NonOKResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"error\": \"things went bad\"}"))
	}))
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListDocuments("")
	if it.Next(context.Background()) == true {
		t.Fatal("expected iterator not to return next item, got next item")
	}
	if it.Err() == nil {
		t.Fatal("expected iterator to return error message, got nil")
	}
}

func TestListDocuments_DocumentsRetrieved(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	expected := onfido.Document{
		ID:           "ce62d838-56f8-4ea5-98be-e7166d1dc33d",
		Href:         "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86",
		DownloadHref: "/v2/live_photos/7410A943-8F00-43D8-98DE-36A774196D86/download",
		FileName:     "localfile.png",
		FileType:     "png",
		FileSize:     282123,
		Type:         onfido.DocumentTypePassport,
		Side:         onfido.DocumentSideBack,
	}
	expectedJson, err := json.Marshal(onfido.Documents{
		Documents: []*onfido.Document{&expected},
	})
	if err != nil {
		t.Fatal(err)
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(expectedJson)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	it := client.ListDocuments(applicantID)
	for it.Next(context.Background()) {
		d := it.Document()

		assert.Equal(t, expected.ID, d.ID)
		assert.Equal(t, expected.Href, d.Href)
		assert.Equal(t, expected.DownloadHref, d.DownloadHref)
		assert.Equal(t, expected.FileName, d.FileName)
		assert.Equal(t, expected.FileType, d.FileType)
		assert.Equal(t, expected.FileSize, d.FileSize)
		assert.Equal(t, expected.Type, d.Type)
		assert.Equal(t, expected.Side, d.Side)
	}
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
}

func TestDownloadDocument(t *testing.T) {
	applicantID := "541d040b-89f8-444b-8921-16b1333bf1c6"
	documentID := "ce62d838-56f8-4ea5-98be-e7166d1dc33d"

	dummy_content := []byte("hello world")

	expected := onfido.DocumentDownload{
		Size:    len(dummy_content),
		Content: dummy_content,
	}

	m := mux.NewRouter()
	m.HandleFunc("/applicants/{applicantId}/documents/{documentId}/download", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, applicantID, vars["applicantId"])
		assert.Equal(t, documentID, vars["documentId"])

		w.WriteHeader(http.StatusOK)
		w.Write(dummy_content)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	dd, err := client.DownloadDocument(context.Background(), applicantID, documentID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Size, dd.Size)
	assert.Equal(t, expected.Content, dd.Content)
}

func TestDownloadDocument(t *testing.T) {
	livePhotoID := "ce62d838-56f8-4ea5-98be-e7166d1dc33d"

	dummy_live_photo := []byte("hi pretty")

	expected := onfido.DocumentDownload{
		Size:    len(dummy_live_photo),
		Content: dummy_live_photo,
	}

	m := mux.NewRouter()
	m.HandleFunc("/live_photos/{livePhotoID}/download", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assert.Equal(t, livePhotoID, vars["livePhotoID"])

		w.WriteHeader(http.StatusOK)
		w.Write(dummy_live_photo)
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	client := onfido.NewClient("123")
	client.Endpoint = srv.URL

	lp, err := client.GetLivePhotoDownload(context.Background(), livePhotoID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Size, lp.Size)
	assert.Equal(t, expected.Content, lp.Content)
}
