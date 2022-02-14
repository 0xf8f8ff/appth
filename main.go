package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/0xf8f8ff/appth/appth"
	"github.com/0xf8f8ff/appth/interceptors"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const GitHubUserEndpoint = "https://api.github.com/user"

type grpcMultiplexer struct {
	*grpcweb.WrappedGrpcServer
}

// Handler is used to route requests to either grpc or to regular http
func (m *grpcMultiplexer) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.IsGrpcWebRequest(r) {
			m.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GenerateTLSApi will load TLS certificates and key and create a grpc server with those.
func GenerateTLSApi(pemPath, keyPath string) (*grpc.Server, error) {
	cred, err := credentials.NewServerTLSFromFile(pemPath, keyPath)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(
		grpc.Creds(cred),
		grpc.ChainUnaryInterceptor(
			interceptors.ValidateRequest,
			interceptors.ValidateUser,
			interceptors.AccessControl,
		),
	)
	return s, nil
}

// Server is the Logic handler for the server
// It has to fullfill the GRPC schema generated Interface
type Server struct {
	appth.UnimplementedAuthServer
	Users map[uint32]*appth.User
}

// Ping fullfills the requirement for PingPong Server interface
func (s *Server) GetUser(ctx context.Context, r *appth.UserRequest) (*appth.User, error) {
	fmt.Println("GetUser called")
	recordedUser, ok := s.Users[r.GetId()]
	if !ok {
		return nil, fmt.Errorf("User with ID %d not found", r.GetId())
	}
	return recordedUser, nil
}

// auth

// format [courseID]role
// roles: 1 - teacher, 2 - student
type UserCourses map[uint32]uint32

type ExternalUser struct {
	Email             string
	Name              string
	Login             string
	UserID            string
	AvatarURL         string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
	IDToken           string
}

func login(w http.ResponseWriter, r *http.Request) {

}

func callbackRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authentication with GitHub complete, redirecting to the main page")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	// We Generate a TLS grpc API
	apiserver, err := GenerateTLSApi("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
	// Start listening on a TCP Port
	lis, err := net.Listen("tcp", "127.0.0.1:9990")
	if err != nil {
		log.Fatal(err)
	}
	// We need to tell the code WHAT TO do on each request, ie. The business logic.
	// In GRPC cases, the Server is acutally just an Interface
	// So we need a struct which fulfills the server interface
	// see server.go
	s := &Server{}
	s.Users = make(map[uint32]*appth.User)
	fakeUser := &appth.User{
		Name:     "Test User",
		Username: "test",
		Isadmin:  true,
	}
	s.Users[1] = fakeUser
	// Register the API server
	// The register function is a generated piece by protoc.
	appth.RegisterAuthServer(apiserver, s)
	// Start serving in a goroutine to not block
	go func() {
		log.Fatal(apiserver.Serve(lis))
	}()
	// Wrap the GRPC Server in grpc-web and also host the UI
	grpcWebServer := grpcweb.WrapServer(apiserver)
	// Lets put the wrapped grpc server in our multiplexer struct so
	// it can reach the grpc server in its handler
	multiplex := grpcMultiplexer{
		grpcWebServer,
	}

	// We need a http router
	r := http.NewServeMux()
	// Load the static webpage with a http fileserver
	webapp := http.FileServer(http.Dir("public/appth/build"))

	// OAuth2 flow client: requests code, exchanges code for tokens, renews tokens
	config := oauth2.Config{
		ClientID:     os.Getenv("GITHUB_KEY"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  "/auth/callback",
	}

	// HTTP client with timeouts to fetch user info of the token owner
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	// Host the Web Application at /, and wrap it in the GRPC Multiplexer
	// This allows grpc requests to transfer over HTTP1. then be
	// routed by the multiplexer
	rediredtSecret := "3krj92f8y35tugh3e3fh4"

	r.Handle("/", multiplex.Handler(webapp))
	r.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {

		authRedirectURL := config.AuthCodeURL(rediredtSecret, oauth2.AccessTypeOffline)
		log.Println("Redirecting to ", authRedirectURL)
		http.Redirect(w, r, authRedirectURL, http.StatusOK)
	})
	r.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		// parse request for code and state
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing request on callback: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// validate state
		callbackSecret := r.FormValue("state")
		log.Println("Callback: got state in request: ", callbackSecret)
		if callbackSecret != rediredtSecret {
			log.Printf("Warning: secrets don't match: expected %s, got %s", rediredtSecret, callbackSecret)
			w.WriteHeader(http.StatusMisdirectedRequest)
			return
		}

		// exchange code for token
		code := r.FormValue("code")
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			log.Printf("Error exchanging code for token: %s", err)
		}
		log.Printf("Successfully fetched access token: %s", token.AccessToken)

		// get user info with the token
		req, err := http.NewRequest("GET", GitHubUserEndpoint, nil)
		if err != nil {
			log.Printf("Error creating request to user API: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("Error requesting user info from the user API: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("User API responded with status: %d: %s", resp.StatusCode, resp.Status)
		}

		// decode user info
		user := ExternalUser{
			AccessToken: token.AccessToken,
		}
		respBits, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response bits from user API: ", err.Error())
		}

		tmpUser := struct {
			ID      int    `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Login   string `json:"login"`
			Picture string `json:"avatar_url"`
		}{}

		if err := json.NewDecoder(bytes.NewReader(respBits)).Decode(&tmpUser); err != nil {
			log.Println("Error decoding user response")
		}
		log.Printf("Got external user from user API: %+v", tmpUser)
		user.Name = tmpUser.Name
		user.Login = tmpUser.Login
		user.Email = tmpUser.Email
		user.AvatarURL = tmpUser.Picture
		user.UserID = strconv.Itoa(tmpUser.ID)

		// check if in the db

		// yes? Check if token needs update and update; start session
		//  no? Create user in db, start session
	})
	// Create a HTTP server and bind the router to it, and set wanted address
	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Serve the webapp over TLS
	log.Fatal(srv.ListenAndServeTLS("cert/server.crt", "cert/server.key"))
}
