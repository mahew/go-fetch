package main

//https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	UpdatedAt   time.Time `json:"created_at"`
	CreatedAt   time.Time `json:"updated_at"`
}

type comment struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type allPosts []post

var posts = allPosts{
	{
		ID:          "0",
		Title:       "First Post",
		Description: "My first post to my blog.",
		Content:     "Content Here.",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

// create

func createPost(w http.ResponseWriter, r *http.Request) {
	var newPost post

	requestBody, e := ioutil.ReadAll(r.Body)

	if e != nil {
		fmt.Fprintf(w, "Make sure post has correct information.")
	}

	now := time.Now()
	newPost.CreatedAt = now
	newPost.UpdatedAt = now

	json.Unmarshal(requestBody, &newPost)
	posts = append(posts, newPost)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newPost)
}

// read

func getOnePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	for _, singlePost := range posts {
		if singlePost.ID == postID {
			json.NewEncoder(w).Encode(singlePost)
			break
		}
	}
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(posts)
}

// update

func updatePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	var updatedPost post

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &updatedPost)

	for i, post := range posts {
		if post.ID == postID {
			post.Title = updatedPost.Title
			post.Description = updatedPost.Description
			post.Content = updatedPost.Content
			now := time.Now()
			post.UpdatedAt = now
			posts = append(posts[:i], post)
			json.NewEncoder(w).Encode(post)
			break
		}
	}
}

// delete

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	for i, post := range posts {
		if post.ID == postID {
			posts = append(posts[:i], posts[i+1:]...)
			fmt.Fprintf(w, "Post ID %v deleted.", postID)
			break
		}
	}
}

// page

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	//initPosts()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/post", createPost).Methods("POST")
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", getOnePost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PATCH")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
