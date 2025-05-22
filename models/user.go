package models

type User struct {
    ID        string `json:"id"`
    Nickname  string `json:"nickname"`
    Age       int    `json:"age"`
    Gender    string `json:"gender"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
    Password  string `json:"-"`
    IsOnline  bool   `json:"isOnline"`
}

type LoginRequest struct {
    Identity string `json:"identity"` 
    Password string `json:"password"`
}