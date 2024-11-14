package users

import (
	"addressBook/config"
	jwttoken "addressBook/helpers/jwtToken"
	"addressBook/helpers/validation"
	"addressBook/models/user"
	repo "addressBook/repo/userrepo"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type UserController struct {
	userrepo *repo.UserRepo
	appCtx   *config.AppCtx
}

func NewUserController(appCtx *config.AppCtx) *UserController {

	repo := repo.NewUserRepo(appCtx.DB)

	return &UserController{userrepo: repo, appCtx: appCtx}
}
func (usrCtrl *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		logger.Fatal("Failed to create a logger : ", zap.Error(err))
	}
	var newUser user.User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to Decode a body : ", zap.Error(err))
		return
	}
	err = validation.ValidateUserInput(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("Validation error ", zap.Error(err))
		return
	}
	passwordBytes := []byte(newUser.Password)
	hash := md5.Sum(passwordBytes)
	hashString := hex.EncodeToString(hash[:])
	err = usrCtrl.userrepo.RegisterUser(newUser.Email, hashString)
	if err != nil {
		// w.WriteHeader(http.StatusConflict)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}
func (usrCtrl *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {

	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("Failed to create a logger : ", zap.Error(err))
	}

	var logUser user.User
	err = json.NewDecoder(r.Body).Decode(&logUser)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("Failed to Decode a body : ", zap.Error(err))
		return
	}

	isValid, err := usrCtrl.userrepo.LoginUser(logUser.Email, logUser.Password)
	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("failed to validate user credentials : ", zap.Error(err))
		return
	}
	if !isValid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		// w.WriteHeader(http.StatusUnauthorized)
		logger.Error("Invalid Credentials")
		return
	}

	token, err := jwttoken.GenerateJWT(logUser.Email)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to generate JWT token ", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
	// json.NewEncoder(w).Encode(logUser)
}

func (usrCtrl *UserController) LogoutUser(w http.ResponseWriter, r *http.Request) {

}
