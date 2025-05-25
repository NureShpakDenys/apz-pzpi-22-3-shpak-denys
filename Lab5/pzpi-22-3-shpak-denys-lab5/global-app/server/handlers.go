package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type UsersConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Config   any    `json:"config"`
}

func (uc *UsersConfig) Bind(r *http.Request) error {
	if uc.Username == "" || uc.Password == "" || uc.Config == nil {
		return errors.New("wrong user config provided")
	}
	return nil
}

func SaveCreds(w http.ResponseWriter, r *http.Request) error {
	var req UsersConfig
	if err := render.Bind(r, &req); err != nil {
		return NewAPIError(http.StatusBadRequest, "invalid input: %s", err)
	}

	data, err := os.ReadFile("configs.json")
	if err != nil {
		if os.IsNotExist(err) {
			data = []byte("[]")
		} else {
			return NewAPIError(http.StatusInternalServerError, "could not read configs: %s", err)
		}
	}

	var configs []UsersConfig
	if len(data) != 0 {
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&configs); err != nil {
			return NewAPIError(http.StatusInternalServerError, "could not decode configs: %s", err)
		}
	}

	configs = append(configs, req)

	data, err = json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not marshal configs: %s", err)
	}

	if err = os.WriteFile("configs.json", data, 0644); err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not write configs: %s", err)
	}

	render.Status(r, http.StatusCreated)
	return nil
}

type GetCredsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *GetCredsRequest) Bind(r *http.Request) error {
	if req.Username == "" || req.Password == "" {
		return errors.New("all fields should be provided")
	}
	return nil
}

func GetCreds(w http.ResponseWriter, r *http.Request) error {
	var req GetCredsRequest
	if err := render.Bind(r, &req); err != nil {
		return NewAPIError(http.StatusBadRequest, "invalid input: %s", err)
	}

	data, err := os.ReadFile("configs.json")
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not read config file: %s", err)
	}

	var configs []UsersConfig
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(&configs); err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not decode config list: %s", err)
	}

	config, err := FindUserConfigs(configs, req.Username, req.Password)
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "error finding user config: %s", err)
	}
	if config == nil {
		return NewNoFieldsFoundError("user config")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, config)
	return nil
}

func GetAPK(w http.ResponseWriter, r *http.Request) error {
	const apkPath = "./data/app-debug.apk"

	file, err := os.Open(apkPath)
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not open APK file: %s", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not get APK file info: %s", err)
	}

	w.Header().Set("Content-Type", "application/vnd.android.package-archive")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file); err != nil {
		log.Printf("Error while sending APK: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to send APK")
	}

	return nil
}

func GetServerExe(w http.ResponseWriter, r *http.Request) error {
	const exePath = "./data/main.exe"

	file, err := os.Open(exePath)
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not open EXE file: %s", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not get EXE file info: %s", err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file); err != nil {
		log.Printf("Error while sending EXE: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to send EXE")
	}

	return nil
}

func GetFrontendZip(w http.ResponseWriter, r *http.Request) error {
	const zipPath = "./data/wayra.zip"

	file, err := os.Open(zipPath)
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not open ZIP file: %s", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "could not get ZIP file info: %s", err)
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file); err != nil {
		log.Printf("Error while sending ZIP: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to send ZIP")
	}

	return nil
}

func FindUserConfigs(configs []UsersConfig, username, password string) (*UsersConfig, error) {
	for _, cfg := range configs {
		if cfg.Username == username {
			if err := bcrypt.CompareHashAndPassword([]byte(cfg.Password), []byte(password)); err == nil {
				return &cfg, nil
			}
			return nil, errors.New("wrong creds")
		}
	}
	return nil, nil
}
